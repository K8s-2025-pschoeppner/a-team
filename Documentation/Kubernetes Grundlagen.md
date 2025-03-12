# Kubernetes Grundlagen
## Helm Package Manager

[[Helm]] ist ein Paketmanager für Kubernetes, der die Verwaltung von Anwendungen erheblich vereinfacht. Es löst mehrere Probleme, die bei der direkten Verwendung von `kubectl` auftreten können:


---

## Resilienz erhöhen

Eine der Hauptstärken von Kubernetes ist die Fähigkeit, die Resilienz von Anwendungen zu verbessern. Das System leitet Traffic nur zu Containern, die als "ready" markiert sind, was sicherstellt, dass Anfragen nur an funktionsfähige Instanzen gesendet werden.

### Probes

Kubernetes verwendet verschiedene Arten von [[Probes]] (Sonden), um den Zustand eines Containers zu überwachen


---

## Ressourcenkontrolle

Die Kontrolle über die Ressourcennutzung ist ein entscheidender Aspekt des Kubernetes-Clusterbetriebs.

## Arten Ressourcen zu Kontrollieren
[[Arten für Ressourcen Kontrolle]]

### Warum Ressourcen kontrollieren?

- **Kosteneffizienz**: Optimierung der Ressourcennutzung in Cloud-Umgebungen
- **Fairness**: Verhinderung der Monopolisierung von Cluster-Ressourcen durch einzelne Anwendungen oder Teams
- **Stabilität**: Schutz vor Ressourcenknappheit und daraus resultierenden Ausfällen
- **Vorhersagbarkeit**: Gewährleistung konsistenter Leistung für alle Anwendungen


### Best Practices für Ressourcenmanagement

- **Realistische Anforderungen**: Setzen Sie Ressourcenanforderungen basierend auf tatsächlichem Bedarf
- **Monitoring**: Überwachen Sie die tatsächliche Ressourcennutzung und passen Sie die Limits entsprechend an
- **Namespace-Organisation**: Gruppieren Sie zusammengehörige Workloads in gemeinsamen Namespaces
- **Hierarchisches Management**: Kombinieren Sie Cluster-weite Policies mit Namespace-spezifischen Quotas


---

## Updates ohne Ausfälle

Die Durchführung von Updates ohne Dienstunterbrechung ist ein wesentliches Merkmal moderner Kubernetes-Anwendungen.

### Warum sind unterbrechungsfreie Updates wichtig?

- **Benutzererfahrung**: Nutzer können die Anwendung ohne Unterbrechung weiter verwenden
- **Geschäftskontinuität**: Keine Umsatzeinbußen durch Ausfallzeiten
- **Häufigere Updates**: Ermöglicht schnellere Fehlerbehebungen und Feature-Rollouts
- **Risikominimierung**: Probleme können früh erkannt und Rollbacks durchgeführt werden, bevor alle Nutzer betroffen sind

### Deployment Strategien

Kubernetes unterstützt verschiedene Strategien für unterbrechungsfreie Updates: [[Deploymentstrategien]]

### Vergleich der Deployment-Strategien

| Strategie      | Vorteile                                                 | Nachteile                                    | Anwendungsfälle                                  |
| -------------- | -------------------------------------------------------- | -------------------------------------------- | ------------------------------------------------ |
| Blue-Green     | Vollständige Validierung vor Switch, sofortiger Rollback | Höherer Ressourcenbedarf                     | Kritische Anwendungen, komplexe Updates          |
| Canary         | Risikominimierung, echte Nutzerfeedbacks                 | Komplexere Konfiguration, längere Updatezeit | Öffentliche Dienste mit hoher Nutzerbasis        |
| Rolling Update | Einfach zu konfigurieren, standardmäßig in Kubernetes    | Begrenzte Kontrolle, gemischte Versionen     | Standard-Updates, backward-kompatible Änderungen |


---

## Storage in Kubernetes

Die Speicherverwaltung ist ein zentraler Aspekt bei der Bereitstellung zustandsbehafteter Anwendungen in Kubernetes.

#### Storage Arten
- [[Storage Arten]]

#### PersistentVolumes und PersistentVolumeClaims

PersistentVolumes (PVs) stellen Speicherressourcen auf Cluster-Ebene dar, während PersistentVolumeClaims (PVCs) Anforderungen für Speicher von Anwendungen sind.

- **PersistentVolume (PV)**:
  - Cluster-weite Speicherressource
  - Von Administratoren oder dynamisch bereitgestellt
  - Unabhängig vom Pod-Lebenszyklus

- **PersistentVolumeClaim (PVC)**:
  - Anforderung eines Pods für Speicher
  - Spezifiziert benötigte Größe und Zugriffsart
  - Wird an passendes PV gebunden

**Funktionsweise**:
1. Administrator erstellt PVs oder richtet dynamische Bereitstellung ein
2. Entwickler erstellt PVC mit Spezifikation des benötigten Speichers
3. Kubernetes bindet PVC an passendes PV
4. Pod verwendet PVC als Volume

**Vorteile**:
- Trennung von Bereitstellung und Nutzung
- Abstraktion der zugrundeliegenden Speicherinfrastruktur
- Dynamische Bereitstellung möglich

Beispiel für einen PersistentVolumeClaim:

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: database-data
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
  storageClassName: standard
```

Verwendung des PVC in einem Pod:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: database
spec:
  containers:
  - name: postgres
    image: postgres:13
    volumeMounts:
    - mountPath: "/var/lib/postgresql/data"
      name: database-storage
  volumes:
  - name: database-storage
    persistentVolumeClaim:
      claimName: database-data
```

#### StorageClasses

StorageClasses definieren verschiedene "Klassen" von Speicher mit unterschiedlichen Leistungsmerkmalen, Bereitstellungsmethoden oder Kostenstrukturen.

- **Funktionen**:
  - Definition von Speichertypen
  - Automatische Bereitstellung von PVs
  - Standardeinstellungen für Speicheranforderungen

Beispiel einer StorageClass:

```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: fast
provisioner: kubernetes.io/aws-ebs
parameters:
  type: gp2
  fsType: ext4
reclaimPolicy: Delete
allowVolumeExpansion: true
```

### Tipps für Speichernutzung in Kubernetes

- **Ressourcenplanung**: Berücksichtigen Sie Wachstumsraten bei der Speicherdimensionierung
- **Backup-Strategien**: Implementieren Sie regelmäßige Backups für persistente Daten
- **Performance-Tests**: Testen Sie verschiedene StorageClasses für Ihre spezifischen Anwendungsfälle
- **Monitoring**: Überwachen Sie Speichernutzung und -leistung
- **Zustandslose Designs bevorzugen**: Wo möglich, entwerfen Sie Anwendungen so, dass sie minimal auf persistenten Speicher angewiesen sind