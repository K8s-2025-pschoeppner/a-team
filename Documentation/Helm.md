# Helm
### Warum Helm verwenden?

- **Vereinfachte Anwendungsverwaltung**: Helm fasst mehrere Kubernetes-Ressourcen in einem einzigen "Chart" zusammen
- **Saubere Installation und Deinstallation**: Verhindert Ressourcenkonflikte und gewährleistet vollständige Entfernung aller zugehörigen Komponenten
- **Versionierung**: Ermöglicht einfaches Rollback auf frühere Versionen
- **Abhängigkeitsverwaltung**: Verwaltet automatisch Abhängigkeiten zwischen verschiedenen Charts
- **Sicherheitsprüfungen**: Weigert sich, Änderungen durchzuführen, die bestehende Ressourcen beschädigen könnten

### Praktisches Beispiel: ingress-nginx

Die Installation des NGINX Ingress Controllers wird mit Helm erheblich vereinfacht:

```bash
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm install ingress-nginx ingress-nginx/ingress-nginx
```

Ohne Helm müssten mehrere YAML-Dateien manuell angewendet werden, und die Deinstallation wäre komplizierter.

### Helm-Befehle im Überblick

| Befehl | Beschreibung |
|--------|--------------|
| `helm install` | Installiert ein Chart |
| `helm upgrade` | Aktualisiert eine Release |
| `helm rollback` | Setzt eine Release auf eine frühere Version zurück |
| `helm uninstall` | Entfernt eine Release |
| `helm list` | Listet alle installierten Releases auf |
| `helm repo add` | Fügt ein Chart-Repository hinzu |