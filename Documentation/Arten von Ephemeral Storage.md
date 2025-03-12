# Arten von Ephemeral Storage

1. **ConfigMaps und Secrets**
   - **Verwendung**: Konfigurationsdaten und sensible Informationen
   - **Eigenschaften**: Schreibgeschützt für Container
   - **Zugriff**: Als Umgebungsvariablen oder als gemountete Dateien

   Beispiel für ConfigMap-Mounting:
   ```yaml
   volumes:
   - name: config-volume
     configMap:
       name: app-config
   ```

2. **Downward API**
   - **Verwendung**: Zugriff auf Pod- und Container-Metadaten
   - **Eigenschaften**: Schreibgeschützt, dynamisch aktualisiert
   - **Beispieldaten**: Pod-Name, Namespace, Labels, Ressourcenlimits

3. **emptyDir**
   - **Verwendung**: Temporärer Speicher für die Lebensdauer eines Pods
   - **Eigenschaften**: Überlebt Container-Neustarts, aber nicht Pod-Neustarts
   - **Anwendungsfälle**: Zwischenspeicher, Scratch Space, Datenaustausch zwischen Containern

   Beispiel für emptyDir:
   ```yaml
   volumes:
   - name: cache-volume
     emptyDir: {}
   ```