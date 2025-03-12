# Kubernetes Networking: Network Address Translation (NAT)
## Grundlagen von NAT

Network Address Translation (NAT) ist eine Technik, die IP-Adressen eines Netzwerks in IP-Adressen eines anderen Netzwerks übersetzt. Dies wird typischerweise eingesetzt, um:

- Private IP-Adressen hinter einer öffentlichen IP-Adresse zu verstecken
- IPv4-Adressknappheit zu bewältigen
- Netzwerksicherheit zu erhöhen, indem interne Netzwerktopologien verborgen werden

### Source NAT (SNAT)

Source NAT (SNAT) ändert die Quell-IP-Adresse eines Pakets:

- **Funktionsweise:** Die Quell-IP und oft auch der Quell-Port werden ersetzt, bevor das Paket das Gateway verlässt
- **Verwendung:** Wird eingesetzt, wenn ein Gerät mit privater IP-Adresse mit dem Internet kommunizieren möchte
- **Rückweg:** Das NAT-Gateway führt Buch über die Übersetzungen, um ankommende Antwortpakete korrekt zurückzuleiten

```
Vor SNAT:
[Pod: 10.1.1.2] ---> [Ziel: 8.8.8.8]

Nach SNAT am Node:
[Node: 192.168.1.5] ---> [Ziel: 8.8.8.8]
```

### Destination NAT (DNAT)

Destination NAT (DNAT) ändert die Ziel-IP-Adresse eines Pakets:

- **Funktionsweise:** Die Ziel-IP und oft auch der Ziel-Port werden ersetzt
- **Verwendung:** Wird für Port-Forwarding und Load-Balancing eingesetzt
- **Kubernetes-Kontext:** Service-IPs werden mittels DNAT auf Pod-IPs umgesetzt

```
Vor DNAT:
[Client] ---> [Service: 10.96.1.5:80]

Nach DNAT durch kube-proxy:
[Client] ---> [Pod: 10.1.1.2:8080]
```

## NAT in Kubernetes

Kubernetes verwendet NAT selektiv in seinem Netzwerkmodell.

### NAT-freie Pod-zu-Pod Kommunikation

Ein Grundprinzip des Kubernetes-Netzwerkmodells ist, dass:

- **Jeder Pod bekommt eine eigene IP-Adresse**
- Diese IP-Adresse ist **innerhalb des gesamten Clusters eindeutig und routbar**
- Pods können sich gegenseitig direkt über ihre IP-Adressen erreichen, **ohne NAT (Network Address Translation)**

Diese NAT-freie Kommunikation zwischen Pods ist einer der wichtigsten Aspekte des Kubernetes-Netzwerkmodells und vermeidet viele der typischen Probleme, die mit NAT verbunden sind:

- Keine Portkonflikte zwischen Anwendungen
- Vereinfachte Service Discovery
- Konsistentes Kommunikationsmodell für Anwendungen

### Wann wird NAT in Kubernetes verwendet?

Obwohl die Pod-zu-Pod-Kommunikation innerhalb des Clusters NAT-frei ist, wird NAT in Kubernetes in folgenden Szenarien eingesetzt:

1. **Ausgehender Traffic ins Internet (SNAT):**
   - Wenn Pods mit externen Diensten außerhalb des Clusters kommunizieren
   - Die Quell-IP wird zur Node-IP umgeschrieben

2. **Service-Implementierung (DNAT):**
   - Übersetzung von Service-IPs auf die tatsächlichen Pod-IPs
   - Wird von kube-proxy implementiert

3. **NodePort und LoadBalancer Services (SNAT in bestimmten Fällen):**
   - Bei Traffic-Weiterleitungen zwischen Nodes

## Source NAT in Kubernetes

### NodePort und LoadBalancer Services

Bei NodePort und LoadBalancer Services wird SNAT unterschiedlich angewendet:

- **NodePort Services:**
  - **Direkter Traffic:** Wenn der Traffic an den Node gerichtet ist, auf dem der Ziel-Pod läuft: **Kein SNAT**
  - **Cross-Node Traffic:** Wenn der Traffic an einen Node gerichtet ist, aber zu einem Pod auf einem anderen Node weitergeleitet werden muss: **SNAT wird angewendet**

```
┌──────────────────────────────────────────────────────────┐
│                     Kubernetes Cluster                   │
│                                                          │
│  ┌─────────────┐                       ┌─────────────┐   │
│  │   Node 1    │                       │   Node 2    │   │
│  │             │                       │             │   │
│  │ NodePort:30080                      │             │   │
│  │             │                       │ ┌─────────┐ │   │
│  │             │     Traffic mit SNAT  │ │ Pod     │ │   │
│  │             │─────────────────────► │ │10.1.2.3 │ │   │
│  │             │  Quell-IP wird zu     │ └─────────┘ │   │
│  │             │  Node 1 IP geändert   │             │   │
│  └─────────────┘                       └─────────────┘   │
└──────────────────────────────────────────────────────────┘
```

- **LoadBalancer Services:**
  - Verhalten ähnlich wie bei NodePort
  - Externe Load Balancer leiten Traffic an Nodes weiter
  - SNAT-Verhalten hängt von der Cloud-Provider-Implementation ab

### Cluster-externe Kommunikation

Wenn Pods mit externen Diensten kommunizieren:

1. Ein Pod mit IP `10.1.1.2` sendet ein Paket an externen Dienst `203.0.113.1`
2. Der Paket erreicht den Node-Gateway
3. Source NAT wird angewendet: Die Quell-IP wird zur Node-IP `192.168.1.5` geändert
4. Der externe Dienst sieht die Anfrage von `192.168.1.5` kommend
5. Die Antwort geht zurück an `192.168.1.5`
6. Der Node empfängt die Antwort und leitet sie basierend auf der NAT-Tabelle zurück an Pod `10.1.1.2`

```
┌────────────────────────────────────────┐   ┌───────────────┐
│             Kubernetes Cluster         │   │               │
│                                        │   │               │
│  ┌─────────────┐      ┌─────────────┐  │   │    Internet   │
│  │    Node     │      │   Gateway   │  │   │               │
│  │             │      │             │  │   │               │
│  │ ┌─────────┐ │      │             │  │   │ ┌─────────────┐
│  │ │  Pod    │ │      │   SNAT:     │  │   │ │ Externer    │
│  │ │10.1.1.2 │─┼───►  │10.1.1.2 -->─┼──┼───┼►│ Dienst      │
│  │ └─────────┘ │      │192.168.1.5  │  │   │ │203.0.113.1  │
│  │             │      │             │  │   │ └─────────────┘
│  └─────────────┘      └─────────────┘  │   │               │
└────────────────────────────────────────┘   └───────────────┘
```

### Externaltraffic Policy

Kubernetes Services bieten eine Option namens `externalTrafficPolicy`, die das SNAT-Verhalten beeinflusst:

- **Cluster (Standard):** 
  - Traffic wird an alle Pods verteilt, unabhängig vom Node
  - SNAT wird angewendet für Cross-Node Traffic
  - Die ursprüngliche Client-IP geht verloren

- **Local:** 
  - Traffic wird nur an Pods auf dem empfangenden Node weitergeleitet
  - Kein SNAT, die Client-IP bleibt erhalten
  - Kann zu ungleichmäßiger Lastverteilung führen, wenn Pods nicht gleichmäßig auf Nodes verteilt sind

```yaml
apiVersion: v1
kind: Service
metadata:
  name: web-service
spec:
  externalTrafficPolicy: Local  # Verhindert SNAT, erhält Client-IP
  ports:
  - port: 80
    targetPort: 8080
  selector:
    app: web
  type: LoadBalancer
```

## Destination NAT in Kubernetes

### Service-IP zu Pod-IP

Ein Kubernetes Service erhält eine virtuelle IP-Adresse (Cluster-IP), die nicht mit einem bestimmten Pod verbunden ist. Wenn Traffic an diese Service-IP gesendet wird, muss Kubernetes den Traffic auf die tatsächlichen Pod-IPs umleiten:

1. Client sendet Traffic an Service-IP (z.B. `10.96.1.5`)
2. kube-proxy (oder ähnliche Komponente) führt DNAT durch
3. Die Ziel-IP wird zu einer der Pod-IPs geändert (z.B. `10.1.1.2`)
4. Die Antwort wird zurück durch den Service geleitet

```
Vor DNAT:
[Client] ---> [Service: 10.96.1.5:80]

Nach DNAT:
[Client] ---> [Pod: 10.1.1.2:8080]
```

### Implementierung in kube-proxy

kube-proxy ist die Komponente, die für die Service-Abstraktion in Kubernetes verantwortlich ist und DNAT implementiert. Sie kann in verschiedenen Modi betrieben werden:

1. **userspace-Modus:**
   - Der älteste Modus, in dem kube-proxy selbst als Proxy fungiert
   - Ineffizient, da Traffic durch den Userspace geleitet wird

2. **iptables-Modus (Standard):**
   - Verwendet Linux iptables für DNAT
   - Keine Proxy-Prozesse im Userspace
   - Bessere Performance, aber keine echte Load-Balancing-Logik (rein zufällige Auswahl)

3. **IPVS-Modus:**
   - Verwendet Linux IP Virtual Server
   - Bessere Performance und erweiterte Load-Balancing-Algorithmen
   - Komplexer zu konfigurieren

Beispiel für iptables-Regeln, die von kube-proxy erstellt werden:

```bash
# DNAT-Regel für einen Service
-A KUBE-SERVICES -d 10.96.1.5/32 -p tcp -m tcp --dport 80 \
  -j KUBE-SVC-XGLOBAL56NCVAMC3

# Load-Balancing zwischen Pods (zufällige Auswahl)
-A KUBE-SVC-XGLOBAL56NCVAMC3 -m statistic --mode random --probability 0.33332999982 \
  -j KUBE-SEP-IQRD2WJI3OHJNAC7
-A KUBE-SVC-XGLOBAL56NCVAMC3 -m statistic --mode random --probability 0.5 \
  -j KUBE-SEP-LH4WIXS23MLEGDUC
-A KUBE-SVC-XGLOBAL56NCVAMC3 \
  -j KUBE-SEP-OCDEKMK7SVHOTCX3

# DNAT zur konkreten Pod-IP
-A KUBE-SEP-IQRD2WJI3OHJNAC7 -p tcp -m tcp \
  -j DNAT --to-destination 10.1.1.2:8080
```

## CNI-Plugin Einfluss auf NAT

Die Wahl des CNI-Plugins beeinflusst, wie NAT in Kubernetes implementiert wird:

### NAT-Konfigurationen in CNI-Plugins

Verschiedene CNI-Plugins bieten unterschiedliche Konfigurationsmöglichkeiten für NAT:

1. **Calico:**
   - Erlaubt detaillierte Kontrolle über NAT-Verhalten mit `natOutgoing` Flag
   - Kann SNAT für bestimmte IP-Pools konfigurieren
   - Unterstützt sowohl vollständiges als auch selektives NAT

   ```yaml
   apiVersion: projectcalico.org/v3
   kind: IPPool
   metadata:
     name: default-ipv4-ippool
   spec:
     cidr: 10.1.0.0/16
     natOutgoing: true  # Aktiviert SNAT für ausgehenden Traffic
   ```

2. **Flannel:**
   - Weniger Konfigurationsoptionen
   - Implementiert standardmäßig SNAT für externe Kommunikation

3. **Cilium:**
   - Fortschrittliche NAT-Kontrolle mit BPF
   - Kann NAT auf Anwendungsebene (L7) durchführen
   - Unterstützt identitätsbasierte Sicherheit

4. **AWS VPC CNI:**
   - Verwendet EC2-Network-Interfaces, was NAT für Pod-zu-Pod-Kommunikation innerhalb des VPC vermeidet
   - Kann bei VPC-übergreifender Kommunikation SNAT benötigen

## Debugging und Troubleshooting

NAT-Probleme können komplex sein. Hier sind einige Debugging-Strategien:

1. **Überprüfen von iptables-Regeln:**
   ```bash
   # Alle NAT-Regeln anzeigen
   sudo iptables -t nat -L -n -v
   
   # Spezifisch nach Kubernetes-Regeln filtern
   sudo iptables -t nat -L KUBE-SERVICES -n -v
   ```

2. **Paketerfassung mit tcpdump:**
   ```bash
   # Traffic auf einem Node überwachen
   sudo tcpdump -n -i any port 80
   
   # Spezifischen Traffic zu/von einer Pod-IP überwachen
   sudo tcpdump -n host 10.1.1.2
   ```

3. **Verbindungs-Tracking überprüfen:**
   ```bash
   # Aktive NAT-Verbindungen anzeigen
   sudo conntrack -L
   
   # Nach bestimmten Verbindungen filtern
   sudo conntrack -L | grep 10.1.1.2
   ```

4. **Service-Konfiguration prüfen:**
   ```bash
   # ExternalTrafficPolicy überprüfen
   kubectl get svc my-service -o jsonpath='{.spec.externalTrafficPolicy}'
   ```

5. **Temporäre Test-Pods erstellen:**
   ```bash
   # Einen Debugging-Pod erstellen
   kubectl run netshoot --rm -it --image nicolaka/netshoot -- /bin/bash
   
   # Im Pod: Quell-IP bei externen Anfragen überprüfen
   curl -v ifconfig.me
   ```

## Best Practices

1. **Source NAT bewusst einsetzen:**
   - Für Services, bei denen die Client-IP erhalten werden muss, `externalTrafficPolicy: Local` verwenden
   - Abwägen zwischen Client-IP-Erhaltung und gleichmäßiger Lastverteilung

2. **Cluster-Egress kontrollieren:**
   - Egress-Gateway oder Proxy für ausgehenden Traffic konfigurieren
   - Ermöglicht Kontrolle und Auditing des ausgehenden Traffics

3. **NAT-Tabellen überwachen:**
   - Größe der Verbindungstabellen im Auge behalten
   - `nf_conntrack_max` und verwandte Parameter bei Bedarf anpassen

4. **Dokumentation:**
   - NAT-Verhalten im Cluster dokumentieren
   - Besonders wichtig bei benutzerdefinierten CNI-Konfigurationen

5. **Überprüfung der Sicherheitsimplikationen:**
   - NAT kann Sicherheitsboundaries beeinflussen
   - Netzwerkpolicies entsprechend anpassen

6. **Performance-Optimierung:**
   - IPVS-Modus von kube-proxy für bessere Performance bei vielen Services verwenden
   - NAT-bezogene Kernel-Parameter optimieren