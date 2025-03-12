# Kubernetes Networking: IP-Adressierung

## Grundlagen der Kubernetes-Netzwerkarchitektur

Kubernetes verwendet ein Netzwerkmodell, das auf vier grundlegenden Anforderungen basiert:

1. **Pods auf einem Node können mit allen Pods auf allen Nodes kommunizieren ohne NAT**
2. **Agents auf einem Node können mit allen Pods auf diesem Node kommunizieren**
3. **Pods mit Host-Netzwerk können mit allen Pods auf allen Nodes kommunizieren ohne NAT**
4. **Die IP, die ein Pod "sieht", ist die gleiche IP, die andere Komponenten sehen**

Diese Grundsätze führen zu einem flachen Netzwerkmodell, das die Kommunikation vereinfacht und Entwicklern ein konsistentes Erlebnis bietet.

## IP-Adressierung in Kubernetes

### Pod IP-Adressierung

Die IP-Adressierung in Kubernetes folgt einem klaren Konzept:

- **Jeder Pod bekommt eine eigene IP-Adresse:** Im Gegensatz zu traditionellen Anwendungsarchitekturen, bei denen mehrere Dienste auf einem Host laufen und Ports teilen, erhält jeder Kubernetes-Pod eine eigene, dedizierte IP-Adresse.

- **Cluster-weite Eindeutigkeit:** Diese IP-Adresse ist **innerhalb des gesamten Clusters eindeutig** - es gibt keine IP-Adresskonflikte zwischen Pods, unabhängig davon, auf welchen Nodes sie laufen.

- **Direkte Routbarkeit:** Alle Pod-IPs sind **im gesamten Cluster routbar**, was bedeutet, dass jeder Pod direkt mit jedem anderen Pod kommunizieren kann, unabhängig von ihrer physischen Platzierung im Cluster.

- **Keine NAT zwischen Pods:** Pods können sich gegenseitig direkt über ihre IP-Adressen erreichen, **ohne NAT (Network Address Translation)**. Dies vereinfacht die Netzwerkkommunikation erheblich und vermeidet Probleme, die typischerweise mit NAT verbunden sind.

```
┌─────────────────────────────────────────────────────┐
│                      Kubernetes Cluster             │
│                                                     │
│  ┌─────────────┐           ┌─────────────┐          │
│  │    Node 1   │           │    Node 2   │          │
│  │             │           │             │          │
│  │ ┌─────────┐ │           │ ┌─────────┐ │          │
│  │ │ Pod A   │ │           │ │ Pod C   │ │          │
│  │ │10.1.1.2 │ │           │ │10.1.2.2 │ │          │
│  │ └─────────┘ │◄─────────►│ └─────────┘ │          │
│  │             │  Direkte  │             │          │
│  │ ┌─────────┐ │Kommunika- │ ┌─────────┐ │          │
│  │ │ Pod B   │ │   tion    │ │ Pod D   │ │          │
│  │ │10.1.1.3 │ │ohne NAT   │ │10.1.2.3 │ │          │
│  │ └─────────┘ │           │ └─────────┘ │          │
│  └─────────────┘           └─────────────┘          │
└─────────────────────────────────────────────────────┘
```

Dieses Modell bietet mehrere Vorteile:
- **Vereinfachte Service Discovery:** Jeder Pod ist über eine eindeutige IP-Adresse erreichbar
- **Konsistentes Kommunikationsmodell:** Ein Pod kommuniziert mit anderen Pods immer auf die gleiche Weise
- **Portabilität von Anwendungen:** Keine Notwendigkeit, Anwendungen für NAT-Traversal zu modifizieren

### CIDR-Blöcke und IP-Pools

Kubernetes verwendet CIDR (Classless Inter-Domain Routing) Blöcke, um IP-Adressbereiche zu definieren:

- **Cluster CIDR:** Definiert den gesamten IP-Adressbereich für Pods im Cluster
  - Wird mit `--cluster-cidr` beim Start des Kubernetes Controller Managers konfiguriert
  - Typische Größe: /16 (65.536 IP-Adressen)
  - Beispiel: `10.32.0.0/16`

- **Node CIDR:** Jeder Node erhält einen Subnetz-Bereich aus dem Cluster CIDR
  - Typische Größe: /24 (256 IP-Adressen pro Node)
  - Beispiel: Node 1 bekommt `10.32.1.0/24`, Node 2 bekommt `10.32.2.0/24`, etc.

- **Service CIDR:** Separater IP-Bereich für Services (nicht für Pods)
  - Wird mit `--service-cluster-ip-range` beim Start des API-Servers konfiguriert
  - Beispiel: `10.96.0.0/16`

### IP-Adresszuweisung

Die IP-Adresszuweisung erfolgt dynamisch:

1. **Bei Pod-Erstellung:** Wenn ein Pod erstellt wird, weist das CNI-Plugin (Container Network Interface) ihm eine IP-Adresse aus dem entsprechenden Node-CIDR-Block zu.

2. **Lebensdauer:** Die IP-Adresse bleibt für die gesamte Lebensdauer des Pods bestehen.

3. **Wiederverwendung:** Wenn ein Pod gelöscht wird, wird seine IP-Adresse wieder in den Pool verfügbarer Adressen zurückgeführt.

4. **Persistenz:** Bei Neustarts eines Pods (ohne Löschung) bleibt die IP-Adresse in der Regel erhalten.

## Pod-zu-Pod Kommunikation

### Direktes Routing

Innerhalb des Kubernetes-Clusters erfolgt die Pod-zu-Pod-Kommunikation ohne NAT:

1. **Auf demselben Node:**
   - Direktes Routing über die lokale Bridge des Node
   - Keine Übersetzung der IP-Adressen

2. **Zwischen verschiedenen Nodes:**
   - Overlay-Netzwerk oder direkte Routing-Tabellen ermöglichen die Kommunikation
   - Die ursprüngliche Pod-IP bleibt als Quell-IP erhalten
   - Der Ziel-Pod sieht die tatsächliche IP-Adresse des Quell-Pods

```
┌────────────────────────────────────────────────────────────┐
│                      Kubernetes Cluster                    │
│                                                            │
│  ┌─────────────┐                       ┌─────────────┐     │
│  │    Node 1   │                       │    Node 2   │     │
│  │             │                       │             │     │
│  │ ┌─────────┐ │  Direkte Kommunikation│ ┌─────────┐ │     │
│  │ │ Pod A   │ │     Quell-IP:         │ │ Pod B   │ │     │
│  │ │10.1.1.2 │─┼─────────────────────► │ │10.1.2.3 │ │     │
│  │ └─────────┘ │    bleibt 10.1.1.2    │ └─────────┘ │     │
│  │             │                       │             │     │
│  └─────────────┘                       └─────────────┘     │
└────────────────────────────────────────────────────────────┘
```

### Cross-Node Kommunikation

Die direkte Pod-zu-Pod-Kommunikation über Node-Grenzen hinweg wird durch verschiedene Netzwerk-Mechanismen ermöglicht:

1. **Overlay-Netzwerke:**
   - Encapsulieren Pod-Traffic in Pakete, die zwischen Nodes transportiert werden können
   - Beispiele: VXLAN, Geneve, IPinIP

2. **BGP-basierte Lösungen:**
   - Werben Node-CIDRs über BGP an das umgebende Netzwerk an
   - Ermöglichen direktes Routing ohne Encapsulation

3. **Kernel-Routing-Tabellen:**
   - Statische Routen zu Pod-CIDRs auf jedem Node
   - Erfordern ein Netzwerk, das diese Routen unterstützt

Die tatsächliche Implementation hängt vom verwendeten CNI-Plugin ab.

## CNI (Container Network Interface)

### Populäre CNI-Plugins

Verschiedene CNI-Plugins implementieren das Kubernetes-Netzwerkmodell auf unterschiedliche Weise:

1. **Calico:**
   - Verwendet BGP für Routing ohne Encapsulation (in Standard-Konfiguration)
   - Bietet auch VXLAN-Modus für Netzwerke, die kein BGP unterstützen
   - Bekannt für starke Netzwerkpolicy-Unterstützung

2. **Flannel:**
   - Einfache Setup mit VXLAN-Overlay-Netzwerk
   - Minimal, aber effektiv für grundlegende Szenarien

3. **Cilium:**
   - eBPF-basiert für hohe Performance
   - Erweiterte Netzwerk- und Sicherheitsfunktionen
   - Anwendungs-Layer-Awareness (L7)

4. **Weave Net:**
   - Mesh-Netzwerk zwischen Nodes
   - Funktioniert in fast jeder Umgebung, auch bei NAT

### Einfluss auf Routing

Die Wahl des CNI-Plugins beeinflusst, wie Routing implementiert wird:

- **Overlay vs. Native Routing:**
  - Overlay-Netzwerke (wie VXLAN) haben höheren Overhead, bieten aber mehr Kompatibilität
  - Native Routing (wie BGP) bietet bessere Performance, erfordert aber Netzwerkuntersützung

- **Netzwerkpolicies:**
  - Die Implementierung von Kubernetes NetworkPolicies variiert je nach CNI-Plugin
  - Einige bieten erweiterte Funktionen über den Kubernetes-Standard hinaus

## Debugging und Troubleshooting

Häufige Probleme und Debugging-Methoden für IP-Adressierung:

1. **Überprüfen der Pod-IP-Adressierung:**
   ```bash
   kubectl get pods -o wide
   ```

2. **Netzwerkverbindung zwischen Pods testen:**
   ```bash
   # Exec in einen Pod und pinge einen anderen
   kubectl exec -it mypod -- ping 10.1.2.3
   
   # Überprüfe Konnektivität mit curl
   kubectl exec -it mypod -- curl http://service-name
   ```

3. **Überprüfen des Routing auf Node-Ebene:**
   ```bash
   # Auf dem Node
   ip route
   
   # Traceroute zu einer Pod-IP
   traceroute 10.1.2.3
   ```

4. **CNI-Plugin-Konfiguration überprüfen:**
   ```bash
   # CNI-Konfigurationsdateien anzeigen
   ls -la /etc/cni/net.d/
   cat /etc/cni/net.d/10-calico.conflist
   ```

5. **IP-Adressbereiche überprüfen:**
   ```bash
   # Cluster-CIDR überprüfen
   kubectl cluster-info dump | grep -m 1 cluster-cidr
   
   # Node-CIDRs überprüfen
   kubectl get nodes -o jsonpath='{.items[*].spec.podCIDR}'
   ```

## Best Practices

1. **IP-Adressplanung:**
   - Ausreichend große CIDR-Blöcke für zukünftiges Wachstum reservieren
   - Node-CIDRs groß genug wählen, um maximale Pod-Dichte zu unterstützen
   - Klare Trennung zwischen Cluster-CIDR und Service-CIDR

2. **Dual-Stack-Konfiguration (IPv4/IPv6):**
   - Falls erforderlich, frühzeitig planen
   - Kompatibilität mit CNI-Plugin sicherstellen

3. **Vermeidung von IP-Überlappungen:**
   - Keine Überlappung mit anderen Netzwerken (VPCs, On-Premises)
   - Dokumentation der verwendeten IP-Bereiche

4. **Monitoring der IP-Auslastung:**
   - Überwachung der IP-Adresszuweisung
   - Alarme bei hoher Auslastung des IP-Pools

5. **CNI-Plugin-Auswahl:**
   - Plugin basierend auf spezifischen Anforderungen wählen (Performance, Sicherheit, Einfachheit)
   - Netzwerkinfrastruktur bei der Auswahl berücksichtigen (BGP-Unterstützung, MTU, etc.)

6. **Network Policy definieren:**
   - Default-Deny-Policies als Basis für Zero-Trust-Netzwerksicherheit einsetzen
   - Explizite Policies für gewünschte Kommunikation definieren

7. **Dokumentation:**
   - IP-Bereiche dokumentieren
   - Netzwerktopologie visualisieren
   - Änderungen am Netzwerk protokollieren
