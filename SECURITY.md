# Security Policy

## Supported Versions

Currently supporting security updates for:

| Version | Supported          |
| ------- | ------------------ |
| main    | :white_check_mark: |

## Reporting a Vulnerability

We take the security of this project seriously. If you discover a security vulnerability, please follow these steps:

### 1. Do NOT Create a Public Issue
Security vulnerabilities should **never** be reported through public GitHub issues.

### 2. Report Privately
Please report security vulnerabilities by sending an email to the repository maintainers or through GitHub's private vulnerability reporting:

1. Navigate to the **Security** tab of this repository
2. Click on **Report a vulnerability**
3. Fill out the form with details about the vulnerability

### 3. Include the Following Information
When reporting, please include:

- Type of vulnerability (e.g., SQL injection, XSS, authentication bypass)
- Full paths of source file(s) related to the vulnerability
- Location of the affected source code (tag/branch/commit or direct URL)
- Step-by-step instructions to reproduce the issue
- Proof-of-concept or exploit code (if possible)
- Impact of the vulnerability and potential attack scenarios

### 4. Response Timeline
- **Initial Response**: Within 48 hours
- **Status Update**: Within 5 business days
- **Resolution Timeline**: Depends on severity
  - Critical: 7 days
  - High: 14 days
  - Medium: 30 days
  - Low: 60 days

### 5. Disclosure Process
1. Security report received and acknowledged
2. Vulnerability confirmed and impact assessed
3. Fix developed and tested
4. Security advisory prepared
5. Fix released and advisory published
6. Credit given to reporter (unless anonymity requested)

## Security Best Practices for Contributors

When contributing to this project:

1. **Never commit secrets**: API keys, passwords, tokens
2. **Validate all inputs**: Prevent injection attacks
3. **Use parameterized queries**: Avoid SQL injection
4. **Sanitize outputs**: Prevent XSS attacks
5. **Check dependencies**: Keep them updated and scan for vulnerabilities
6. **Follow secure coding guidelines**: OWASP Top 10

## Security Features

This repository has the following security features enabled:

- ✅ Secret scanning
- ✅ Secret scanning push protection
- ✅ Dependabot vulnerability alerts
- ✅ Dependabot security updates
- ✅ Branch protection on main
- ✅ Required PR reviews

## Resources

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [Go Security Best Practices](https://golang.org/doc/security)
- [GitHub Security Features](https://docs.github.com/en/code-security)