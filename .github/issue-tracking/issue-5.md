# Issue #5: TestReports

**Issue URL**: https://github.com/f00b455/ts-pure-template/issues/5
**Created by**: @f00b455
**Auto-processed**: Wed Sep 24 15:39:55 UTC 2025

## Description
@claude

User Story — TestReports

Als Entwickler möchte ich alle von der CI/CD-Pipeline erzeugten Testberichte (z. B. Unit, Integration, E2E) automatisch auf der GitHub-Wiki-Seite sehen, damit ich Build-Qualität und Testergebnisse zentral, versionsbezogen und ohne lokale Tools nachvollziehen kann.

Business Value

Transparente Qualität, weniger „wo ist der Report?“, schnellere Fehleranalyse, bessere Team-Kommunikation.

⸻

BDD-Ansatz (erst Feature-File & Tests, dann Implementierung)
	1.	Zuerst: Feature-File in Gherkin anlegen (features/wiki_testreports.feature).
	2.	Dann: Step-Definitions & Pipeline-Tests (z. B. Skript-Unit-Tests für Index-Generierung, Smoke-Test für Wiki-Push).
	3.	Iterativ implementieren, bis alle Szenarien grün sind.

Gherkin (Feature-File Skeleton)

Feature: Testberichte in GitHub Wiki veröffentlichen
  Um Testergebnisse transparent zu machen,
  möchte ich, dass die CI/CD-Pipeline Reports nach jedem Lauf in die Wiki pusht.

  Background:
    Given das Repository hat die GitHub Wiki aktiviert
    And die Pipeline erzeugt Testberichte als Artefakte (z. B. JUnit XML, HTML)

  @happy-path
  Scenario: Erfolgreiche Veröffentlichung eines Laufes
    Given ein erfolgreicher Pipeline-Lauf auf Branch "main" mit Commit-SHA <sha>
    When die Publish-Stage ausgeführt wird
    Then werden HTML-Reports unter "wiki/reports/main/<run-id>/" abgelegt
    And die Wiki-Startseite verlinkt "Letzter Lauf: <run-id> (<sha>)"
    And die Seite ist öffentlich innerhalb des Projekts einsehbar

  @multiple-branches
  Scenario: Getrennte Verzeichnisse pro Branch
    Given ein Lauf auf Branch "feature-x"
    Then werden Reports unter "wiki/reports/feature-x/<run-id>/" gespeichert
    And ein Branch-Index existiert und ist verlinkt

  @index
  Scenario: Automatische Index-Aktualisierung
    Given ein neuer Report wurde abgelegt
    When das Index-Skript läuft
    Then werden die letzten N Läufe pro Branch gelistet, absteigend sortiert nach Zeit

  @fallback
  Scenario: Kein Report vorhanden
    Given die Tests erzeugen keine Artefakte oder sind fehlgeschlagen
    Then wird ein Eintrag "Keine Reports verfügbar" mit Zeitstempel erzeugt
    And die Pipeline schlägt nicht allein wegen des Wiki-Schritts fehl (warnend)

  @retention
  Scenario: Aufräumen alter Reports
    Given mehr als N Läufe existieren
    When die Publish-Stage läuft
    Then werden ältere Verzeichnisse entfernt, sodass die Wiki-Größe kontrolliert bleibt


⸻

Akzeptanzkriterien (aus den Szenarien)
	•	Automatisch pro Lauf: Nach jedem Pipeline-Run steht ein verlinkter Report im Wiki.
	•	Struktur: reports/<branch>/<run-id>/ mit index.html (oder Markdown-Seiten) je Lauf.
	•	Index: Startseite zeigt „Letzter Lauf“ + Liste der letzten N (z. B. 20) Läufe je Branch.
	•	Formate: Mindestens JUnit-XML → HTML konvertiert; weitere Formate optional (Allure, Jest, Pytest, Cypress).
	•	Robustheit: Wiki-Push darf die Build-Stage nicht „rot“ machen, wenn nur das Publizieren scheitert (best effort + Warnung).
	•	Sicherheit: Nur benötigte Rechte; kein Leaken von Secrets in Wiki.
	•	Performance: Veröffentlichungs-Schritt ≤ 60 s bei typischen Report-Größen; Repo-Bloat wird begrenzt (Retention).

Nicht-funktionale Kriterien
	•	Rechteverwaltung „least privilege“ (Token mit contents:write).
	•	Reproduzierbarkeit: deterministisches Index-Skript (stabile Sortierung, konsistente Pfade).
	•	Lesbarkeit: klare Navigation, Datum/Uhrzeit in Europe/Berlin.

⸻

Technische Skizze (GitHub Actions – Beispiel)

Minimal lauffähig, später ausbauen (Konverter, Allure, Pruning, schöne Templates).

name: CI

on:
  push:
    branches: [ "**" ]
  pull_request:

permissions:
  contents: write   # für Wiki push
  actions: read

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      # ⬇️ Stelle deine Tests hier her (Beispiele):
      # - name: Run tests (JUnit XML)
      #   run: |
      #     npm ci
      #     npm test -- --reporters=junit --outputFile=reports/junit.xml

      # Demo: Reports-Ordner bereitstellen
      - name: Prepare demo report
        run: |
          mkdir -p reports
          echo "<html><body><h1>Demo Test Report</h1><p>Run $GITHUB_RUN_ID</p></body></html>" > reports/index.html

      - name: Clone Wiki
        run: |
          git config --global user.email "bot@users.noreply.github.com"
          git config --global user.name "ci-bot"
          git clone "https://x-access-token:${GITHUB_TOKEN}@github.com/${GITHUB_REPOSITORY}.wiki.git" wiki
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

      - name: Publish reports to Wiki
        run: |
          BRANCH="${GITHUB_HEAD_REF:-${GITHUB_REF_NAME}}"
          TARGET="wiki/reports/${BRANCH}/${GITHUB_RUN_ID}"
          mkdir -p "$TARGET"
          cp -r reports/* "$TARGET/"

          # Index aktualisieren (Startseite)
          START="wiki/Home.md"
          DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
          mkdir -p wiki/reports/${BRANCH}
          # Einfache Indexzeile
          printf "## Letzter Lauf (%s)\n\n- **Branch:** %s  \n- **Run:** [%s](reports/%s/%s/index.html)  \n- **Commit:** %s  \n- **Zeit:** %s (UTC)\n\n" \
            "${GITHUB_RUN_ID}" "${BRANCH}" "${GITHUB_RUN_ID}" "${BRANCH}" "${GITHUB_RUN_ID}" "${GITHUB_SHA}" "${DATE}" > wiki/_latest.md

          # Home.md zusammensetzen (ultra simpel, später templaten)
          {
            echo "# Test Reports"
            echo
            cat wiki/_latest.md
            echo "### Historie"
            find wiki/reports -type d -mindepth 2 -maxdepth 2 | sort -r | head -n 50 | while read d; do
              RUNID=$(basename "$d")
              BR=$(basename $(dirname "$d"))
              echo "- $BR / $RUNID → [Report](reports/$BR/$RUNID/index.html)"
            done
          } > "$START"

          # Optional: Retention (halte nur letzte 20 pro Branch)
          for BR in $(ls wiki/reports); do
            KEEP=20
            ls -1dt wiki/reports/$BR/* | tail -n +$((KEEP+1)) | xargs -r rm -rf
          done

          pushd wiki
          git add .
          git commit -m "Publish reports: ${GITHUB_REF_NAME} run ${GITHUB_RUN_ID}"
          git push
          popd

⚠️ Hinweise:
	•	Wiki aktivieren: Repo-Settings → „Wikis“ einschalten.
	•	GITHUB_TOKEN hat i. d. R. genug Rechte für das Wiki (gleiche „contents“-Berechtigung). Falls Policies das verhindern: PAT in secrets.WIKI_TOKEN mit repo-Scope nutzen und im Clone-Step einsetzen.
	•	Für JUnit→HTML: z. B. junit2html oder native Reporter deiner Testumgebung. Für schicke Dashboards: Allure (später), statisch in Wiki veröffentlichen.
	•	PR-Runs: Optional nur bei push auf main/release/* veröffentlichen, PRs lediglich als Artefakt beifügen.

⸻

Inkrementeller Plan (Test-First)
	1.	Feature-File & Steps erstellen → rote Tests.
	2.	Index-Skript als eigenständiges Modul + Unit-Tests (Pfadlogik, Sortierung, Retention).
	3.	Pipeline-Job minimal: HTML kopieren, Wiki pushen → Szenario @happy-path grün.
	4.	Mehr Formate (JUnit, Cypress, Jest) + Konvertierung → @multiple-branches, @index grün.
	5.	Robustheit (Retries, Warnungen statt Fail) + Retention → @fallback, @retention grün.
	6.	Schöne Templates (Markdown/HTML, Logos, Badges) und A11y.

Definition of Ready
	•	Wiki aktiviert, benötigte Secrets/Permissions vorhanden.
	•	Testreport-Pfad(e) festgelegt (Konvention).
	•	Retention-Policy und Ziel-Branches entschieden.

Definition of Done
	•	Alle BDD-Szenarien grün.
	•	Reports erscheinen zuverlässig nach jedem (gewählten) Lauf.
	•	Index nachvollziehbar, saubere Navigation, begrenzte Repo-Größe.

## Work Log
- Branch created: issue-5-testreports (via cron job)
- [ ] Implementation
- [ ] Tests
- [ ] Documentation
