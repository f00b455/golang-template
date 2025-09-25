# Issue: #7
# URL: https://github.com/f00b455/ts-pure-template/issues/7
@pkg(web) @issue-7
Feature: Header zeigt Top 5 SPIEGEL-Schlagzeilen inkl. Datum
  Damit ich mehrere aktuelle News schnell erfassen kann,
  möchte ich im Header bis zu 5 Headlines mit Datum sehen.

  Background:
    Given die Anwendung ist geöffnet
    And der Header ist sichtbar

  @happy-path
  Scenario: Anzeige von bis zu 5 neuesten Schlagzeilen mit Datum
    Given die API liefert mindestens 5 Einträge
    When die Header-News-Komponente initialisiert
    Then werden genau 5 Headlines angezeigt
    And jede Headline zeigt den Titel und das Veröffentlichungsdatum
    And die Liste ist nach Datum absteigend sortiert (neueste zuerst)
    And ein Klick auf eine Headline öffnet den Artikel in einem neuen Tab

  @less-than-5
  Scenario: Weniger als 5 Einträge
    Given die API liefert 3 Einträge
    When die Komponente lädt
    Then werden genau 3 Headlines angezeigt (keine Platzhalter für fehlende)

  @date-format
  Scenario: Datumsformat und Zeitzone
    Given ein Eintrag hat das UTC-Datum "2025-09-24T08:05:00Z"
    When die Komponente rendert
    Then wird das Datum in Europe/Berlin als "24.09.2025 10:05" angezeigt

  @refresh
  Scenario: Automatisches Aktualisieren
    Given die Seite ist länger geöffnet
    When 5 Minuten vergangen sind
    Then wird die Liste via API aktualisiert
    And die Reihenfolge und Anzahl bleiben konsistent

  @error
  Scenario: Fehler/Fallback
    Given der API-Aufruf scheitert oder ist leer
    When die Komponente rendert
    Then erscheint eine dezente Fallback-Nachricht
    And es gibt keine Layoutsprünge