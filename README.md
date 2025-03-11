# ctf

Capture the flag client and server for Kubernetes

## Connect to K8s cluster

``` bash
mkdir -p ~/.kube
nano ~/.kube/config1

export KUBECONFIG=~/.kube/config1

kubectl cluster-info
```

## Create pods/services

``` bash
kubectl delete -f ./pod.yaml

kubectl create -f ./pod.yaml

```

## Apply changes to pods/services

``` bash
kubectl delete pod github-action-pod -n a-team
kubectl apply -f ./pod.yaml
```

## Check the cluster

``` bash
kubectl get pods -n <your_namespace>

kubectl get svc -n <your_namespace>

kubectl get ing -n <your_namespace>

kubectl logs <pod_name> -n <your_namespace>

kubectl describe pod <pod_name> -n <your_namespace>
```

## Add a serviceAccount to K8s

You find all related files in `roles` directory.

1. Create a service account
2. Create a role
3. Create a rolebinding
4. Add to the pod

``` bash
spec:
  serviceAccountName: ctf-serviceaccount
```

## Deploy an application with alreaady running ingress controller

1. Create a deployment
2. Create a service
3. Create an ingress

## Create different workload resources

Find all related files in `workload` directory.

- Deployment: Manages stateless applications by running a set of identical pods with automatic scaling, updates, and rollbacks.
- DaemonSet: Ensures a single pod runs on each (or selected) node, typically used for logging, monitoring, or node-specific services.
- CronJob: Executes pods periodically on a defined schedule, ideal for recurring tasks like backups or scheduled data processing.
- Job: Runs pods to completion for one-time, batch-processing tasks, such as database migrations or batch computations.

## Create a job disruption

Find all related files in `disruption` directory.

This ensure a minimun of X pods are running at all times.
In the PodDisruptionBudget you configure the minimum number of pods and the application this applies to.




# Kubernetes Networking

## Netzwerk innerhalb eines Nodes

### Pod Netzwerk-Namespaces
- Jeder Pod hat einen eigenen Netzwerk-Namespace
- Container innerhalb desselben Pods **teilen** sich diesen Netzwerk-Namespace
- Diese Architektur ermöglicht, dass alle Container eines Pods über **localhost** miteinander kommunizieren können

### Loop Back Interface
- Alle Container innerhalb eines Pods können über das Loop Back Interface (`localhost`) kommunizieren
- Dies ist eine fundamentale Eigenschaft des Pod-Konzepts in Kubernetes

### Port-Sharing
- Da alle Container eines Pods denselben Netzwerk-Namespace teilen, teilen sie sich auch den Port-Bereich
- Konsequenz: Container innerhalb eines Pods können **nicht** denselben Port nutzen
- Beispiel: Wenn Container A bereits Port 80 verwendet, kann Container B **nicht** ebenfalls Port 80 im selben Pod verwenden

### Container Network Interface (CNI)
- Kubernetes verwendet das Container Network Interface (CNI) als Standard für die Netzwerkkonfiguration
- CNI ist eine Schnittstelle, die von verschiedenen Netzwerk-Plugins implementiert werden kann
- Bekannte CNI-Implementierungen: Calico, Flannel, Weave, Cilium
- CNI-Plugins sind verantwortlich für:
  - Zuweisung von IP-Adressen an Pods
  - Einrichtung von Routing-Regeln
  - Implementierung von Network Policies

---

## Netzwerk innerhalb des Clusters

### Pod-zu-Pod-Kommunikation
- Kubernetes garantiert, dass jeder Pod mit jedem anderen Pod im Cluster kommunizieren kann
- Diese Kommunikation erfolgt unabhängig davon, auf welchen Nodes die Pods ausgeführt werden
- Jeder Pod erhält eine **eindeutige IP-Adresse** innerhalb des Clusters

### IP-Adressierung
- Jeder Pod bekommt eine eigene IP-Adresse
- Diese IP-Adresse ist **innerhalb des gesamten Clusters** eindeutig und routbar
- Pods können sich gegenseitig direkt über ihre IP-Adressen erreichen, ohne NAT (Network Address Translation)

### Underlay vs. Overlay Network
- **Underlay Network**: Das physische Netzwerk, in dem die Nodes existieren
  - Beispiel: Das Datencenter-Netzwerk oder VPC in der Cloud
- **Overlay Network**: Ein virtuelles Netzwerk, das über dem Underlay Network aufgebaut wird
  - Wird für die Kommunikation zwischen Pods verwendet
  - Ermöglicht eine flexible Zuweisung von IP-Adressen unabhängig vom physischen Netzwerk

### Wie das Overlay Network funktioniert
- Nodes erhalten IP-Adressen im Underlay Network (physisches Netzwerk)
- Darüber wird ein virtuelles Overlay Network aufgebaut, in dem die Pods ihre IPs erhalten
- Das Overlay Network existiert nur innerhalb des Kubernetes-Clusters
- Es ist in der Regel nicht direkt von außen erreichbar
- Netzwerkpakete zwischen Pods werden in Pakete des Underlay Networks gekapselt (Encapsulation)

### Vorteile des Overlay Netzwerks
- Effiziente Nutzung des IP-Adressraums
- Flexibilität bei der Zuweisung von IP-Bereichen
- Isolation vom Underlay Network
- Einfachere Migration von Pods zwischen Nodes

---

## Service API

### Grundkonzept
- Services sind eine Abstraktionsschicht über Pods
- Sie bieten eine stabile Netzwerkadresse für eine Gruppe von Pods
- Services verwenden **Labels und Selektoren**, um zu bestimmen, welche Pods Teil des Services sind

### Wie Services funktionieren
- Ein Service definiert einen Selektor (z.B. `app: webserver`)
- Kubernetes identifiziert alle Pods mit den entsprechenden Labels
- Der Service verwaltet die Liste der IP-Adressen dieser Pods
- Anfragen an den Service werden mittels Load Balancing an die Pods verteilt (standardmäßig Round Robin)

### Service-Typen

#### ClusterIP (Standard)
- Der Service erhält eine stabile, interne IP-Adresse
- Nur innerhalb des Clusters erreichbar
- Ideal für interne Kommunikation zwischen Anwendungen

#### NodePort
- Erweitert ClusterIP
- Öffnet einen spezifischen Port auf **allen Nodes**
- Dieser Port wird auf allen Nodes geöffnet, unabhängig davon, ob der Pod auf dieser Node läuft
- Ermöglicht Zugriff von außerhalb des Clusters über `<Node-IP>:<NodePort>`
- Der Port-Bereich liegt standardmäßig zwischen 30000-32767

#### LoadBalancer
- Erweitert NodePort
- Fordert zusätzlich einen externen Load Balancer von der Cloud-Plattform an
- Der Cloud Controller Manager erkennt diese Anfrage und erstellt einen externen Load Balancer
- Der externe Load Balancer leitet den Verkehr an den NodePort aller Nodes weiter
- Wichtig zu verstehen: Der Traffic kann an jede Node im Cluster gelangen, auch an Nodes, auf denen kein relevanter Pod läuft
- Wenn der Traffic an einer Node ankommt:
	1. kube-proxy auf dieser Node erkennt den Traffic für den Service
	2. kube-proxy verwendet iptables-Regeln, um zu entscheiden, an welchen Pod der Traffic weitergeleitet wird
	3. Die Ziel-Pod-Auswahl erfolgt nach dem Round-Robin-Prinzip über alle verfügbaren Pods des Service
	4. Falls der ausgewählte Pod auf einer anderen Node läuft, wird der Traffic durch das Cluster-Netzwerk an diesen weitergeleitet
- Dieser Mechanismus ermöglicht, dass jede Node im Cluster als "Eintrittspunkt" für den Service fungieren kann

### kube-proxy
- Implementiert die Service-Funktionalität auf jeder Node
- Verwaltet die Netzwerkregeln (meist iptables-Regeln), um Traffic an die richtigen Pods weiterzuleiten
- Hat mehrere Betriebsmodi:
  - **iptables-Modus** (Standard): Verwendet Linux iptables für Paketweiterleitung
  - **IPVS-Modus**: Verwendet Linux IP Virtual Server für verbesserte Leistung bei vielen Services
  - **userspace-Modus** (veraltet): Weiterleitung über einen Proxy-Prozess

---

## CoreDNS und Service Discovery

### DNS in Kubernetes
- Jeder Pod und jeder Service erhält einen DNS-Namen
- Erleichtert die Kommunikation zwischen Anwendungen, da statt IP-Adressen Namen verwendet werden können
- Die DNS-Auflösung wird von CoreDNS bereitgestellt

### Namenskonventionen
- **Services**: `<service-name>.<namespace>.svc.<cluster-domain>`
  - Beispiel: `mysql.database.svc.cluster.local`
- **Pods**: `<pod-ip-mit-bindestrichen-statt-punkten>.<namespace>.pod.<cluster-domain>`
  - Beispiel: `10-244-2-5.default.pod.cluster.local`

### CoreDNS
- In Go geschriebener DNS-Server, der in Kubernetes als Standard-DNS-Lösung dient
- Hochgradig erweiterbar durch Plugins
- Konfigurierbar über ConfigMaps in Kubernetes
- Unterstützt DNS-basierte Service Discovery
- Die Standard-Cluster-Domain ist `cluster.local`

### Service Discovery
- Anwendungen können andere Dienste über ihre DNS-Namen finden
- Beispiel: Eine Web-App verbindet sich mit einer Datenbank über `mysql.database.svc.cluster.local`
- Bei Änderungen der Pod-IPs (z.B. nach Neustart) bleibt der DNS-Name stabil
- Dies ermöglicht lose Kopplung zwischen Anwendungskomponenten

---

## Kommunikation mit externen Netzwerken

### Von innen nach außen (Egress)
- Pods können in der Regel mit externen Netzwerken kommunizieren
- Dabei wird **Source NAT (SNAT)** verwendet:
  - Die Quell-IP-Adresse des Pods wird durch die IP-Adresse der Node ersetzt
  - Dies ermöglicht, dass Antwortpakete korrekt zurückgeleitet werden können
  - Der Prozess ist für die Anwendung transparent

### Von außen nach innen (Ingress)
- Verschiedene Methoden, um externen Zugriff auf Pods zu ermöglichen:
  - **NodePort Services**
  - **LoadBalancer Services**
  - **Ingress-Ressourcen**

### Cloud Controller Manager
- Integriert Kubernetes mit der Cloud-Infrastruktur
- Verantwortlich für die Erstellung externer Load Balancer bei Verwendung des Service-Typs LoadBalancer
- Überwacht Kubernetes-Ressourcen und aktualisiert die entsprechenden Cloud-Ressourcen

---

## Host Network

### Konzept
- Pods können optional im Netzwerk-Namespace der Host-Node ausgeführt werden
- Dies wird durch Setzen von `hostNetwork: true` in der Pod-Spezifikation aktiviert

### Eigenschaften
- Der Pod verwendet direkt die Netzwerkschnittstellen der Node
- Der Pod erhält die IP-Adresse der Node
- Der Pod kann alle Ports der Node direkt nutzen (auch privilegierte Ports unter 1024)
- Es gibt keine Netzwerkisolation zwischen dem Pod und der Node

### Anwendungsfälle
- Netzwerk-Monitoring-Tools, die direkten Zugriff auf die Netzwerkschnittstellen benötigen
- Netzwerk-Plugins und CNI-Implementierungen
- Spezielle Infrastrukturkomponenten wie kube-proxy selbst
- Performance-kritische Anwendungen, die die Overhead des Overlay-Netzwerks vermeiden müssen

### Sicherheitsaspekte
- Weniger isoliert als reguläre Pods
- Höheres Sicherheitsrisiko, da Pods mit direktem Zugriff auf Node-Netzwerkressourcen laufen
- Sollte nur für spezifische Anwendungsfälle verwendet werden, bei denen es absolut notwendig ist

---

## Ingress API

### Grundkonzept
- Ingress definiert Regeln für den externen Zugriff auf Services im Cluster
- Fungiert als Layer-7-Load-Balancer (HTTP/HTTPS)
- Ermöglicht URL-basiertes Routing, SSL-Terminierung und mehr

### Ingress vs. Service
- **Services** arbeiten auf Layer 4 (TCP/UDP)
- **Ingress** arbeitet auf Layer 7 (HTTP/HTTPS)
- Ingress kann mehrere Services über eine einzige IP-Adresse und URL-Pfade bereitstellen

### Ingress Controller
- Kubernetes definiert nur die Ingress-API, aber keinen Controller
- Ein Ingress Controller muss separat installiert werden
- Beliebte Ingress-Controller:
  - NGINX Ingress Controller
  - Traefik
  - HAProxy Ingress
  - Kong Ingress
  - AWS ALB Ingress Controller

### Funktionen
- **Path-basiertes Routing**: Weiterleitung von Anfragen basierend auf URL-Pfaden
  - Beispiel: `/api/*` → API-Service, `/admin/*` → Admin-Service
- **Host-basiertes Routing**: Weiterleitung basierend auf dem Host-Header
  - Beispiel: `api.example.com` → API-Service, `www.example.com` → Web-Service
- **TLS/SSL-Terminierung**: Verwaltet TLS-Zertifikate für verschlüsselte Verbindungen
- **Rewriting**: Ändern von URL-Pfaden vor der Weiterleitung an Services

### IngressClass
- Ermöglicht die Verwendung mehrerer Ingress-Controller im selben Cluster
- Jeder Ingress kann einer spezifischen IngressClass zugeordnet werden
- Nützlich für verschiedene Umgebungen (z.B. Staging vs. Production) oder für spezialisierte Ingress-Controller

---

## Gateway API

### Überblick
- Die Gateway API ist der Nachfolger der Ingress API
- Bietet ein flexibleres und ausdrucksstärkeres Modell für das Routing von Netzwerkverkehr
- Aktuell im Beta-Stadium, aber bereits in vielen Clustern nutzbar

### Vorteile gegenüber Ingress
- Unterstützt mehrere Protokolle (nicht nur HTTP/HTTPS)
- Bessere Trennung der Verantwortlichkeiten
- Detailliertere Konfigurationsmöglichkeiten
- Bessere Multi-Tenant-Unterstützung

### Gateway API-Ressourcen
- **GatewayClass**: Definiert die Implementierung (ähnlich wie IngressClass)
- **Gateway**: Repräsentiert einen Load Balancer-Endpunkt
- **HTTPRoute**: Definiert HTTP-Routing-Regeln
- **TCPRoute**, **UDPRoute**, **TLSRoute**: Definieren Routing für andere Protokolle
- **ReferenceGrant**: Erlaubt Cross-Namespace-Referenzen (für Multi-Tenant-Szenarien)

### Status
- Die Gateway API wird aktiv entwickelt
- Viele Ingress-Controller-Anbieter implementieren bereits Unterstützung
- Langfristig wird die Gateway API die Ingress API ersetzen

---

## Network Policies

### Grundkonzept
- Network Policies definieren Regeln für den Netzwerkverkehr zu und von Pods
- Sie funktionieren wie eine Firewall auf Pod-Ebene
- Standardmäßig sind, ohne definierte Network Policies, alle Pods für Verbindungen offen

### Funktionsweise
- Basieren auf Labels und Selektoren, ähnlich wie Services
- Können sowohl eingehenden als auch ausgehenden Verkehr steuern
- Regeln können auf IP-Bereiche, Ports, Protokolle und Namespaces angewendet werden

### Beispiele für Network Policies
- Beschränkung des Zugriffs auf eine Datenbank nur für bestimmte Anwendungs-Pods
- Beschränkung des ausgehenden Verkehrs auf bestimmte externe IP-Bereiche
- Blockierung aller Verbindungen außer HTTP/HTTPS

### Unterstützung durch CNI-Plugins
- Nicht alle CNI-Plugins unterstützen Network Policies
- Unterstützende CNI-Plugins:
  - Calico
  - Cilium
  - Weave Net
  - Antrea
- Flannel unterstützt Network Policies nicht nativ, kann aber mit Calico kombiniert werden

### Best Practices
- Beginnen mit restriktiven Default-Policies
- Spezifische Policies für jede Anwendung definieren
- Policies regelmäßig überprüfen und aktualisieren
- Testen von Policies in einer nicht-produktiven Umgebung