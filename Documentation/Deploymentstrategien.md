# Deploymentstrategien
#### Blue-Green Deployment

Blue-Green Deployment ist eine Technik, bei der zwei identische Umgebungen parallel betrieben werden: eine aktive (Blue/Production) und eine inaktive (Green/Development).

- **Funktionsweise**:
  1. Blue-Umgebung ist aktiv und bedient Produktionstraffic
  2. Updates werden in der Green-Umgebung getestet und validiert
  3. Nach erfolgreicher Validierung wird der Traffic von Blue zu Green umgeleitet
  4. Die ehemalige Blue-Umgebung wird zur neuen Development-Umgebung

- **Vorteile**:
  - Vollständige Testmöglichkeit vor dem Switch
  - Sofortige Rollback-Möglichkeit durch Zurückschalten auf die alte Umgebung
  - Keine gemischten Versionen während des Updates

- **Implementierung in Kubernetes**:
  - Verwendung von Labels zur Unterscheidung zwischen Blue und Green
  - Änderung der Service-Selector, um den Traffic umzuleiten
  - Resource Quotas sicherstellen, um ausreichend Kapazität für beide Umgebungen zu haben

Beispiel für Blue-Green Deployment mit Services:

```yaml
# Blue Service (aktuell in Produktion)
apiVersion: v1
kind: Service
metadata:
  name: frontend
spec:
  selector:
    app: frontend
    version: v1  # Blue-Umgebung
  ports:
  - port: 80
    targetPort: 8080

# Green Deployment (neue Version)
apiVersion: apps/v1
kind: Deployment
metadata:
  name: frontend-v2
spec:
  replicas: 3
  selector:
    matchLabels:
      app: frontend
      version: v2  # Green-Umgebung
  template:
    metadata:
      labels:
        app: frontend
        version: v2
    spec:
      containers:
      - name: frontend
        image: frontend:v2
```

#### Canary Deployment

Canary Deployment ermöglicht die schrittweise Einführung einer neuen Version, indem zunächst nur ein kleiner Prozentsatz des Traffics zur neuen Version geleitet wird.

- **Funktionsweise**:
  1. Neue Version wird mit minimaler Anzahl von Replicas deployt
  2. Ein kleiner Prozentsatz des Traffics wird zur neuen Version geleitet
  3. Bei erfolgreicher Validierung wird der Prozentsatz schrittweise erhöht
  4. Nach vollständiger Validierung ersetzt die neue Version die alte vollständig

- **Vorteile**:
  - Risikominimierung durch graduelle Einführung
  - Früherkennung von Problemen mit echtem Traffic
  - Möglichkeit zur A/B-Testung neuer Features

- **Implementierung**:
  - Direkt in Kubernetes nicht nativ unterstützt
  - Erfordert zusätzliche Komponenten wie Ingress-Controller (z.B. ingress-nginx) oder Service Mesh
  - Gateway API als moderne Alternative zu Ingress bietet bessere Unterstützung

Beispiel für Canary Deployment mit NGINX Ingress Controller:

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: frontend-canary
  annotations:
    nginx.ingress.kubernetes.io/canary: "true"
    nginx.ingress.kubernetes.io/canary-weight: "20"
spec:
  rules:
  - host: myapp.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: frontend-v2
            port:
              number: 80
```