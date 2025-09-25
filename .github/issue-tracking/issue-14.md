# Issue #14: User Story Hello-World-CLI

**Issue URL**: https://github.com/f00b455/ts-pure-template/issues/14
**Created by**: @f00b455
**Auto-processed**: 2025-09-24T18:36:00Z

## Description
Als Nutzer:in einer Hello-World-CLI
möchte ich per hello-cli --name <Name> eine farbenfrohe Ausgabe erhalten,
damit ich sofort sehe, dass das Tool läuft und „Hello" sagt – inkl. Spinner (5 s) und einer Dummy-Progressbar, die in 3 s auf 100 % läuft.
Dabei nutzt die CLI die bestehende Funktion fooGreet aus lib-foo, um die eigentliche Grußnachricht zu generieren.

## Akzeptanzkriterien
1. **Greet über lib**: Die finale Ausgabe basiert auf fooGreet(name) aus lib-foo.
2. **Spinner**: Nach Start erscheint ein Spinner für ~5 Sekunden und signalisiert danach „Bereit!" o. ä.
3. **Progressbar**: Direkt im Anschluss läuft eine Progressbar in ~3 Sekunden von 0 % auf 100 %.
4. **Bunte Ausgabe**: Die finale Grußnachricht erscheint **farbig/gradient** und eingerahmt (eye-candy).
5. **CLI-Ergonomie**: --name ist optional (Default vorhanden). Fehler führen zu klarer, farbiger Ausgabe und Exit Code ≠ 0.

## Empfohlene Libraries
- **Args**: cac
- **Farben/Style**: picocolors, gradient-string, boxen
- **Spinner**: nanospinner
- **Progressbar**: cli-progress

## Work Log
- Branch created: issue-14-user-story-hello-world-cli (manually created due to automation gap)
- [ ] Implementation
- [ ] Tests
- [ ] Documentation