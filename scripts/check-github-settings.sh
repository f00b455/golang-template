#!/bin/bash

# Script to check and configure GitHub repository settings
# Usage: ./scripts/check-github-settings.sh [--fix]

set -uo pipefail

REPO=$(gh repo view --json nameWithOwner -q '.nameWithOwner')
FIX_MODE=false
NEEDS_FIX=false
REPO_NEEDS_FIX=false
SECURITY_NEEDS_FIX=false
LABELS_NEEDS_FIX=false

if [[ "${1:-}" == "--fix" ]]; then
    FIX_MODE=true
fi

echo "🔍 Checking GitHub repository settings for: $REPO"
echo "================================================"

# Function to check a setting
check_setting() {
    local name="$1"
    local current="$2"
    local expected="$3"
    local status="✅"

    if [[ "$current" != "$expected" ]]; then
        status="❌"
    fi

    printf "%-50s %s (current: %s, expected: %s)\n" "$name:" "$status" "$current" "$expected"

    if [[ "$current" != "$expected" ]]; then
        return 1
    fi
    return 0
}

# 1. Check Actions permissions
echo -e "\n📋 GitHub Actions Permissions:"
echo "------------------------------"

# Get current Actions settings - using the correct endpoint
ACTIONS_SETTINGS=$(gh api repos/$REPO/actions/permissions/workflow)
WORKFLOW_PERMISSIONS=$(echo "$ACTIONS_SETTINGS" | jq -r '.default_workflow_permissions // "read"')
CAN_APPROVE_PR=$(echo "$ACTIONS_SETTINGS" | jq -r '.can_approve_pull_request_reviews // false')

check_setting "Workflow permissions" "$WORKFLOW_PERMISSIONS" "write" || NEEDS_FIX=true
check_setting "Can create/approve PRs" "$CAN_APPROVE_PR" "true" || NEEDS_FIX=true

if [[ "$NEEDS_FIX" == "true" ]] && [[ "$FIX_MODE" == "true" ]]; then
    echo -e "\n🔧 Fixing Actions permissions..."
    gh api -X PUT repos/$REPO/actions/permissions/workflow \
        --field default_workflow_permissions=write \
        --field can_approve_pull_request_reviews=true
    echo "✅ Actions permissions updated!"
fi

# 2. Check repository settings
echo -e "\n📋 Repository Settings:"
echo "----------------------"

REPO_SETTINGS=$(gh api repos/$REPO)
AUTO_MERGE=$(echo "$REPO_SETTINGS" | jq -r '.allow_auto_merge')
DELETE_BRANCH=$(echo "$REPO_SETTINGS" | jq -r '.delete_branch_on_merge')
SQUASH_MERGE=$(echo "$REPO_SETTINGS" | jq -r '.allow_squash_merge')
MERGE_COMMIT=$(echo "$REPO_SETTINGS" | jq -r '.allow_merge_commit')
REBASE_MERGE=$(echo "$REPO_SETTINGS" | jq -r '.allow_rebase_merge')
ISSUES_ENABLED=$(echo "$REPO_SETTINGS" | jq -r '.has_issues')
WIKI_ENABLED=$(echo "$REPO_SETTINGS" | jq -r '.has_wiki')
PROJECTS_ENABLED=$(echo "$REPO_SETTINGS" | jq -r '.has_projects')

check_setting "Issues enabled" "$ISSUES_ENABLED" "true" || REPO_NEEDS_FIX=true
check_setting "Wiki enabled" "$WIKI_ENABLED" "true" || REPO_NEEDS_FIX=true
check_setting "Projects enabled" "$PROJECTS_ENABLED" "true" || REPO_NEEDS_FIX=true
check_setting "Auto-merge allowed" "$AUTO_MERGE" "false"
check_setting "Delete branch on merge" "$DELETE_BRANCH" "true" || REPO_NEEDS_FIX=true
check_setting "Squash merge allowed" "$SQUASH_MERGE" "true"
check_setting "Merge commit allowed" "$MERGE_COMMIT" "true"
check_setting "Rebase merge allowed" "$REBASE_MERGE" "true"

if [[ "$REPO_NEEDS_FIX" == "true" ]] && [[ "$FIX_MODE" == "true" ]]; then
    echo -e "\n🔧 Fixing repository settings..."

    # Fix delete branch on merge if needed
    if [[ "$DELETE_BRANCH" != "true" ]]; then
        gh api -X PATCH repos/$REPO \
            --field delete_branch_on_merge=true
        echo "✅ Delete branch on merge enabled!"
    fi

    # Fix issues if needed
    if [[ "$ISSUES_ENABLED" != "true" ]]; then
        gh api -X PATCH repos/$REPO \
            --field has_issues=true
        echo "✅ Issues enabled!"
    fi

    # Enable wiki if needed
    if [[ "$WIKI_ENABLED" != "true" ]]; then
        gh api -X PATCH repos/$REPO \
            --field has_wiki=true
        echo "✅ Wiki enabled!"
    fi

    # Enable projects if needed
    if [[ "$PROJECTS_ENABLED" != "true" ]]; then
        gh api -X PATCH repos/$REPO \
            --field has_projects=true
        echo "✅ Projects enabled!"
    fi
fi

# 3. Check branch protection (if main branch exists)
echo -e "\n📋 Branch Protection (main):"
echo "----------------------------"

# Check if main branch has protection
PROTECTION=$(gh api repos/$REPO/branches/main/protection 2>/dev/null || echo '{"message":"Branch not protected"}')

if [[ "$PROTECTION" == *"Branch not protected"* ]] || [[ "$PROTECTION" == "{}" ]] || [[ -z "$PROTECTION" ]]; then
    echo "❌ No branch protection rules for 'main' branch"
    echo "   Direct pushes to main are allowed!"

    if [[ "$FIX_MODE" == "true" ]]; then
        echo -e "\n🔧 Setting up branch protection for main..."
        echo "   - Requiring pull requests before merging"
        echo "   - Disabling force pushes"
        echo "   - Protecting branch from deletion"

        # Set up branch protection: require PR, no direct pushes
        # Create JSON payload for branch protection
        cat > /tmp/branch-protection.json << EOF
{
  "required_status_checks": {
    "strict": false,
    "contexts": []
  },
  "enforce_admins": false,
  "required_pull_request_reviews": {
    "dismiss_stale_reviews": false,
    "require_code_owner_reviews": false,
    "required_approving_review_count": 0
  },
  "restrictions": null,
  "allow_force_pushes": false,
  "allow_deletions": false
}
EOF
        gh api -X PUT repos/$REPO/branches/main/protection \
            --input /tmp/branch-protection.json 2>/dev/null && echo "✅ Branch protection enabled!" || echo "⚠️  Note: Branch protection requires admin access"
        rm -f /tmp/branch-protection.json
    fi
else
    echo "✅ Branch protection is configured"

    REQUIRE_PR=$(echo "$PROTECTION" | jq -r '.required_pull_request_reviews // null')
    ALLOW_FORCE=$(echo "$PROTECTION" | jq -r '.allow_force_pushes.enabled // false')
    ALLOW_DELETIONS=$(echo "$PROTECTION" | jq -r '.allow_deletions.enabled // false')

    PROTECTION_NEEDS_FIX=false

    if [[ "$REQUIRE_PR" != "null" ]]; then
        echo "   ✅ Pull request required for merging"
    else
        echo "   ⚠️  Direct pushes to main allowed (consider requiring PR)"
    fi

    if [[ "$ALLOW_FORCE" == "false" ]]; then
        echo "   ✅ Force pushes disabled"
    else
        echo "   ❌ Force pushes allowed (should be disabled)"
        PROTECTION_NEEDS_FIX=true
    fi

    if [[ "$ALLOW_DELETIONS" == "false" ]]; then
        echo "   ✅ Branch deletion protection enabled"
    else
        echo "   ❌ Branch can be deleted (should be protected)"
        PROTECTION_NEEDS_FIX=true
    fi

    # Fix existing branch protection if needed
    if [[ "$PROTECTION_NEEDS_FIX" == "true" ]] && [[ "$FIX_MODE" == "true" ]]; then
        echo -e "\n🔧 Updating branch protection settings..."

        # Get current settings to preserve them
        CURRENT_STATUS_CHECKS=$(echo "$PROTECTION" | jq '.required_status_checks // null')
        CURRENT_PR_REVIEWS=$(echo "$PROTECTION" | jq '.required_pull_request_reviews // null')
        CURRENT_RESTRICTIONS=$(echo "$PROTECTION" | jq '.restrictions // null')
        CURRENT_ENFORCE_ADMINS=$(echo "$PROTECTION" | jq -r '.enforce_admins.enabled // false')

        # Create JSON payload for updating branch protection
        cat > /tmp/branch-protection-update.json << EOF
{
  "required_status_checks": $CURRENT_STATUS_CHECKS,
  "enforce_admins": $CURRENT_ENFORCE_ADMINS,
  "required_pull_request_reviews": $CURRENT_PR_REVIEWS,
  "restrictions": $CURRENT_RESTRICTIONS,
  "allow_force_pushes": false,
  "allow_deletions": false
}
EOF
        gh api -X PUT repos/$REPO/branches/main/protection \
            --input /tmp/branch-protection-update.json 2>/dev/null && echo "✅ Branch protection updated!" || echo "⚠️  Note: Branch protection requires admin access"
        rm -f /tmp/branch-protection-update.json
    fi
fi

# 4. Check workflows
echo -e "\n📋 GitHub Actions Workflows:"
echo "----------------------------"

if [[ -d ".github/workflows" ]]; then
    WORKFLOWS=$(ls -1 .github/workflows/*.yml .github/workflows/*.yaml 2>/dev/null | wc -l | xargs)
    echo "Workflow files found: $WORKFLOWS"
    for workflow in .github/workflows/*.yml .github/workflows/*.yaml; do
        [[ -f "$workflow" ]] || continue
        basename "$workflow" | sed 's/^/  - /'
    done
else
    echo "⚠️  No .github/workflows directory found"
fi

# 5. Check Security Features
echo -e "\n📋 Security Features:"
echo "--------------------"

SECURITY_SETTINGS=$(gh api repos/$REPO | jq -r '.security_and_analysis')
VULN_ALERTS=$(gh api repos/$REPO/vulnerability-alerts --silent 2>/dev/null && echo "enabled" || echo "disabled")
DEPENDABOT_UPDATES=$(echo "$SECURITY_SETTINGS" | jq -r '.dependabot_security_updates.status // "disabled"')
SECRET_VALIDITY=$(echo "$SECURITY_SETTINGS" | jq -r '.secret_scanning_validity_checks.status // "disabled"')

check_setting "Vulnerability alerts" "$VULN_ALERTS" "enabled" || SECURITY_NEEDS_FIX=true
check_setting "Dependabot security updates" "$DEPENDABOT_UPDATES" "enabled" || SECURITY_NEEDS_FIX=true
check_setting "Secret scanning validity checks" "$SECRET_VALIDITY" "enabled" || SECURITY_NEEDS_FIX=true

if [[ "$SECURITY_NEEDS_FIX" == "true" ]] && [[ "$FIX_MODE" == "true" ]]; then
    echo -e "\n🔧 Fixing security settings..."

    # Enable vulnerability alerts
    if [[ "$VULN_ALERTS" == "disabled" ]]; then
        gh api -X PUT repos/$REPO/vulnerability-alerts 2>/dev/null && echo "✅ Vulnerability alerts enabled!" || echo "⚠️  Could not enable vulnerability alerts"
    fi

    # Enable Dependabot security updates
    if [[ "$DEPENDABOT_UPDATES" == "disabled" ]]; then
        gh api -X PUT repos/$REPO/automated-security-fixes 2>/dev/null && echo "✅ Dependabot security updates enabled!" || echo "⚠️  Could not enable Dependabot"
    fi
fi

# 6. Check and create standard labels
echo -e "\n📋 Repository Labels:"
echo "--------------------"

# Define standard labels (name:color:description)
STANDARD_LABELS=(
    "bug:d73a4a:Something isn't working"
    "enhancement:a2eeef:New feature or request"
    "documentation:0075ca:Improvements or additions to documentation"
    "user-story:1d76db:User story for feature development"
    "triage:d876e3:Needs review and categorization"
    "work-in-progress:fbca04:Pull request is still being worked on"
    "good-first-issue:7057ff:Good for newcomers"
    "help-wanted:008672:Extra attention is needed"
    "wontfix:ffffff:This will not be worked on"
    "duplicate:cfd3d7:This issue or pull request already exists"
    "question:d876e3:Further information is requested"
    "breaking-change:d73a4a:Introduces breaking changes"
    "performance:f9d0c4:Performance improvements"
    "security:ee0000:Security vulnerability or improvement"
    "testing:bfd4f2:Related to tests and testing"
    "ci/cd:0e8a16:Continuous Integration/Deployment"
)

# Get existing labels
EXISTING_LABELS=$(gh label list --json name,color,description --limit 100)

LABELS_CREATED=0
LABELS_EXIST=0

for label_def in "${STANDARD_LABELS[@]}"; do
    IFS=':' read -r name color description <<< "$label_def"

    # Check if label exists
    EXISTS=$(echo "$EXISTING_LABELS" | jq -r --arg name "$name" '.[] | select(.name == $name) | .name')

    if [[ -z "$EXISTS" ]]; then
        echo "❌ Missing label: $name"
        LABELS_NEEDS_FIX=true

        if [[ "$FIX_MODE" == "true" ]]; then
            gh label create "$name" --color "$color" --description "$description" 2>/dev/null && \
                echo "  ✅ Created label: $name" && \
                ((LABELS_CREATED++)) || \
                echo "  ⚠️  Could not create label: $name"
        fi
    else
        echo "✅ Label exists: $name"
        ((LABELS_EXIST++))
    fi
done

if [[ "$LABELS_CREATED" -gt 0 ]]; then
    echo -e "\n✅ Created $LABELS_CREATED new labels"
fi

# 7. Check Claude Code Review setup
echo -e "\n📋 Claude Code Review Setup:"
echo "-----------------------------"

# Check if Claude Code Review workflow exists
if [[ -f ".github/workflows/claude-code-review.yml" ]]; then
    echo "✅ Claude Code Review workflow exists"
else
    echo "❌ Claude Code Review workflow not found"
fi

# Check if Claude interactive workflow exists
if [[ -f ".github/workflows/claude.yml" ]]; then
    echo "✅ Claude interactive workflow exists"
else
    echo "❌ Claude interactive workflow not found"
fi

# Check if Claude iterative review workflow exists
if [[ -f ".github/workflows/claude-iterative-review.yml" ]]; then
    echo "✅ Claude Iterative Review workflow exists"
else
    echo "❌ Claude Iterative Review workflow not found"
    echo "   This workflow automatically fixes code review suggestions with max 3 iterations"
fi

# Check if CLAUDE_CODE_OAUTH_TOKEN secret exists
SECRET_EXISTS=$(gh secret list | grep -c "CLAUDE_CODE_OAUTH_TOKEN" || true)

if [[ "$SECRET_EXISTS" -gt 0 ]]; then
    echo "✅ CLAUDE_CODE_OAUTH_TOKEN secret is configured"
else
    echo "❌ CLAUDE_CODE_OAUTH_TOKEN secret is missing"
    echo "   To fix: Set up OAuth token using 'claude setup-token' and add to secrets"
    echo "   gh secret set CLAUDE_CODE_OAUTH_TOKEN"

    if [[ "$FIX_MODE" == "true" ]]; then
        echo ""
        echo "⚠️  Cannot automatically set CLAUDE_CODE_OAUTH_TOKEN - manual setup required:"
        echo "   1. Run 'claude setup-token' in your terminal"
        echo "   2. Copy the generated token"
        echo "   3. Run: gh secret set CLAUDE_CODE_OAUTH_TOKEN"
        echo "   4. Paste the token and press Enter"
    fi
fi

# Check workflow authentication configuration
WORKFLOWS_CONFIGURED=0
for workflow in "claude-code-review.yml" "claude.yml" "claude-iterative-review.yml"; do
    if [[ -f ".github/workflows/$workflow" ]]; then
        if grep -q "claude_code_oauth_token:" ".github/workflows/$workflow"; then
            echo "✅ $workflow configured to use OAuth token"
            ((WORKFLOWS_CONFIGURED++))
        else
            echo "❌ $workflow not using OAuth token authentication"
        fi
    fi
done

if [[ "$WORKFLOWS_CONFIGURED" -eq 0 ]] && [[ "$FIX_MODE" == "true" ]]; then
    echo -e "\n🔧 Creating missing Claude workflows..."
    mkdir -p .github/workflows

    # Create basic Claude Code Review workflow if missing
    if [[ ! -f ".github/workflows/claude-code-review.yml" ]]; then
        cat > .github/workflows/claude-code-review.yml << 'EOF'
name: Claude Code Review

on:
  pull_request:
    types: [opened, synchronize]

jobs:
  claude-review:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      pull-requests: write
      issues: read
      id-token: write

    steps:
      - name: Checkout repository
        uses: actions/checkout@v4
        with:
          fetch-depth: 1

      - name: Run Claude Code Review
        id: claude-review
        uses: anthropics/claude-code-action@v1
        with:
          claude_code_oauth_token: ${{ secrets.CLAUDE_CODE_OAUTH_TOKEN }}
          prompt: |
            Please review this pull request and provide feedback on:
            - Code quality and best practices
            - Potential bugs or issues
            - Performance considerations
            - Security concerns
            - Test coverage

            Use the repository's CLAUDE.md for guidance on style and conventions. Be constructive and helpful in your feedback.

            Use `gh pr comment` with your Bash tool to leave your review as a comment on the PR.

          claude_args: '--allowed-tools "Bash(gh issue view:*),Bash(gh search:*),Bash(gh issue list:*),Bash(gh pr comment:*),Bash(gh pr diff:*),Bash(gh pr view:*),Bash(gh pr list:*)"'
EOF
        echo "✅ Created basic Claude Code Review workflow"
    fi

    echo "   Remember to:"
    echo "   1. Run 'claude setup-token' to get OAuth token"
    echo "   2. Add token as secret: gh secret set CLAUDE_CODE_OAUTH_TOKEN"
    echo "   3. Commit and push the workflow files"
fi

# 8. Check open issues and PRs
echo -e "\n📋 Issues and Pull Requests:"
echo "-----------------------------"

OPEN_ISSUES=$(gh issue list --state open --json number | jq 'length')
OPEN_PRS=$(gh pr list --state open --json number | jq 'length')

echo "Open issues: $OPEN_ISSUES"
echo "Open PRs: $OPEN_PRS"

# 9. Summary
echo -e "\n================================================"
if [[ "$FIX_MODE" == "false" ]]; then
    echo "ℹ️  Run with --fix flag to automatically fix issues"
else
    echo "✅ Fix mode completed"
fi

echo -e "\n📝 Quick Commands Reference:"
echo "----------------------------"
echo "# Enable Actions write permissions:"
echo "gh api -X PUT repos/$REPO/actions/permissions/workflow \\"
echo "  --field default_workflow_permissions=write \\"
echo "  --field can_approve_pull_request_reviews=true"
echo ""
echo "# Enable auto-delete branches on merge:"
echo "gh api -X PATCH repos/$REPO \\"
echo "  --field delete_branch_on_merge=true"
echo ""
echo "# Set up branch protection:"
echo "gh api -X PUT repos/$REPO/branches/main/protection \\"
echo "  --field required_pull_request_reviews='{\"required_approving_review_count\":0}' \\"
echo "  --field allow_force_pushes=false \\"
echo "  --field allow_deletions=false"