# Arten für Ressourcen Kontrolle
### Resource Quotas

Resource Quotas begrenzen die Gesamtmenge an Ressourcen, die innerhalb eines Namespaces verwendet werden können.

- **Funktion**: Legen Obergrenzen für CPU, Arbeitsspeicher, Speicherplatz, Pod-Anzahl etc. fest
- **Durchsetzung**: Anfragen, die Quotas überschreiten, werden vom API-Server abgelehnt
- **Granularität**: Wird auf Namespace-Ebene angewendet

Beispiel einer ResourceQuota:

```yaml
apiVersion: v1
kind: ResourceQuota
metadata:
  name: compute-resources
  namespace: team-alpha
spec:
  hard:
    pods: "10"
    requests.cpu: "4"
    requests.memory: 8Gi
    limits.cpu: "8"
    limits.memory: 16Gi
```

### Limit Ranges

Limit Ranges definieren Default-, Minimal- und Maximalwerte für Ressourcenanforderungen und -limits von Containern in einem Namespace.

- **Hauptzweck**: Sicherstellung, dass alle Pods angemessene Ressourcenlimits haben, auch wenn diese nicht explizit angegeben sind
- **Anwendung**: Automatisches Setzen von Default-Werten für Ressourcenanforderungen und -limits
- **Überprüfung**: Validierung, dass Ressourcenanforderungen und -limits innerhalb akzeptabler Bereiche liegen

Beispiel eines LimitRange:

```yaml
apiVersion: v1
kind: LimitRange
metadata:
  name: default-limits
  namespace: team-alpha
spec:
  limits:
  - default:
      cpu: 500m
      memory: 512Mi
    defaultRequest:
      cpu: 200m
      memory: 256Mi
    max:
      cpu: 1000m
      memory: 1Gi
    min:
      cpu: 100m
      memory: 128Mi
    type: Container
```