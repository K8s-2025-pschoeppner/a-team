# Probes
## Verschiedene Arten
#### Readiness Probes

Readiness Probes bestimmen, ob ein Container bereit ist, Traffic zu empfangen.

- **Funktion**: Kubernetes leitet Traffic nur zu Pods, deren Readiness Probes erfolgreich sind
- **Implementierung**: Führt Befehle aus und überprüft den Exit-Code (0 = erfolgreich)
- **Erfolgsmetrik**: Typischerweise gilt ein Container als bereit, wenn eine bestimmte Anzahl von Checks innerhalb eines definierten Zeitraums erfolgreich ist (z.B. 5 von 10 in 10 Sekunden)

Beispiel einer Readiness Probe in einem Deployment:

```yaml
readinessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 10
  successThreshold: 1
  failureThreshold: 3
```

#### Liveness Probes

Liveness Probes überprüfen, ob ein Container noch funktionsfähig ist oder ob er neu gestartet werden sollte.

- **Hauptzweck**: Erkennung von Deadlocks und anderen schwerwiegenden Zustandsfehlern
- **Wann wichtig**: Besonders am Anfang, bevor Readiness Probes aktiv werden
- **Konsequenz bei Fehlschlag**: Kubernetes startet den Container neu

Beispiel einer Liveness Probe:

```yaml
livenessProbe:
  exec:
    command:
    - cat
    - /tmp/healthy
  initialDelaySeconds: 15
  periodSeconds: 20
```

#### Startup Probes

Startup Probes sind speziell für Container mit langen Startzeiten konzipiert.

- **Zweck**: Verhindern, dass Liveness Probes Container während des Startvorgangs beenden
- **Funktionsweise**: Sobald die Startup Probe erfolgreich ist, übernehmen die Liveness Probes
- **Anwendungsfall**: Ideal für Legacy-Anwendungen oder solche mit langen Initialisierungszeiten

Beispiel einer Startup Probe:

```yaml
startupProbe:
  httpGet:
    path: /healthz
    port: liveness-port
  failureThreshold: 30
  periodSeconds: 10
```

### Best Practices für Probes

- **Leichtgewichtige Checks**: Probes sollten minimale Ressourcen verbrauchen
- **Angemessene Timeouts**: Setzen Sie realistische Timeout-Werte basierend auf der erwarteten Antwortzeit
- **Spezifische Endpunkte**: Implementieren Sie dedizierte Health-Check-Endpunkte in Ihrer Anwendung
- **Kombinieren von Probes**: Verwenden Sie verschiedene Probe-Typen gemeinsam für optimale Resilienz