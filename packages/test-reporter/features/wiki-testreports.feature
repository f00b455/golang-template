# Issue: #5
# URL: https://github.com/f00b455/ts-pure-template/issues/5
@pkg(test-reporter) @issue-5
Feature: Testberichte in GitHub Wiki veröffentlichen
  Als Entwickler
  möchte ich alle von der CI/CD-Pipeline erzeugten Testberichte automatisch auf der GitHub-Wiki-Seite sehen
  damit ich Build-Qualität und Testergebnisse zentral, versionsbezogen und ohne lokale Tools nachvollziehen kann

  Background:
    Given das Repository hat die GitHub Wiki aktiviert
    And die Pipeline erzeugt Testberichte als Artefakte

  @happy-path
  Scenario: Erfolgreiche Veröffentlichung eines Laufes
    Given ein erfolgreicher Pipeline-Lauf auf Branch "main" mit Commit-SHA "abc123"
    When die Publish-Stage ausgeführt wird
    Then werden HTML-Reports unter "wiki/reports/main/<run-id>/" abgelegt
    And die Wiki-Startseite verlinkt "Letzter Lauf: <run-id> (abc123)"
    And die Seite ist öffentlich innerhalb des Projekts einsehbar

  @multiple-branches
  Scenario: Getrennte Verzeichnisse pro Branch
    Given ein Lauf auf Branch "feature-x"
    When die Publish-Stage ausgeführt wird
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
    When die Publish-Stage ausgeführt wird
    Then wird ein Eintrag "Keine Reports verfügbar" mit Zeitstempel erzeugt
    And die Pipeline schlägt nicht allein wegen des Wiki-Schritts fehl

  @retention
  Scenario: Aufräumen alter Reports
    Given mehr als 20 Läufe existieren
    When die Publish-Stage läuft
    Then werden ältere Verzeichnisse entfernt, sodass die Wiki-Größe kontrolliert bleibt