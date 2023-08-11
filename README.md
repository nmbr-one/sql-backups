# Wie nutze ich die Software?
## Windows
1. Downloade die neueste Version der Software über den folgenden Link:
   
   https://github.com/nmbr-one/sql-backups/releases/latest/download/sql-backups-windows-amd64.zip

2. Entpacke die Archivdatei, sodass die sql-backups.exe direkt in einem Ordner liegt.
3. Erstelle eine Datei namens config.properties mit den Werten aus der [Beispiel-Konfigurationsdatei](https://github.com/nmbr-one/sql-backups/blob/main/example.config.properties) im gleichen Ordner wie die sql-backups.exe.
4. Passe die Werte in der config.properties so an, dass die Anmeldedaten für deinen MySQL-/MariaDB-Server, der Pfad zu mysqldump und die Webhook-URL für Discord stimmen.
5. Doppelklicke die sql-backups.exe. Alternativ kannst du auch eine cmd oder powershell öffnen, den Befehl ``.\sql-backups.exe`` dort eingeben und mit Enter bestätigen. Letzteres macht besonders dann Sinn, wenn ein Doppelklick offenbar keine Wirkung hat, da in einem Terminal die Fehlermeldungen besser lesbar sind.
## Linux
1. Downloade die neueste Version der Software über den folgenden Link:
   
   https://github.com/nmbr-one/sql-backups/releases/latest/download/sql-backups-linux-amd64.tar.gz

   Beispiel per wget-Befehl:
   ```
   wget https://github.com/nmbr-one/sql-backups/releases/latest/download/sql-backups-linux-amd64.tar.gz
   ```
2. Entpacke die Archivdatei, sodass die Datei sql-backups direkt in einem Ordner liegt.
   
   Beispiel per tar-Befehl:
   ```
   tar xvf sql-backups-linux-amd64.tar.gz
   ```
3. Erstelle eine Datei namens config.properties mit den Werten aus der [Beispiel-Konfigurationsdatei](https://github.com/nmbr-one/sql-backups/blob/main/example.config.properties) im gleichen Ordner wie die Datei sql-backups.
4. Passe die Werte in der config.properties so an, dass die Anmeldedaten für deinen MySQL-/MariaDB-Server und die Webhook-URL für Discord stimmen.
5. Führe die Datei sql-backups in einer screen-Session aus.
  
   (Ein Tutorial über die Installation und Handhabung von screen findest du [hier](https://linuxhint.com/screen-linux/).)

   Beispiel:
   ```
   screen -dmS sql-backups ./sql-backups
   ```
6. Nutze folgenden Befehl, um die laufenden screen-Sessions anzuzeigen. So kannst du überprüfen, ob die Session erfolgreich angelaufen ist:
   ```
   screen -ls
   ```
   Sollte keine Session mit dem Namen sql-backups in der Liste vorkommen, führe die Datei sql-backups direkt aus, um die Fehlermeldung lesen zu können:
   ```
   ./sql-backups
   ```
# Was bedeuten die Werte in der Konfigurationsdatei? (config.properties)
## mysqldump-path
Hierbei handelt es sich um den absoluten Pfad zu deiner mysqldump Datei. Diese wird bei MySQL- und MariaDB-Installationen automatisch beigeliefert und ist nötig, um die Datenbank-Backups durchzuführen. Ich rate also dazu, diese Software auf dem selben Server auszuführen, auf dem auch der MySQL-/MariaDB-Server installiert ist.

Bei einer sauberen MariaDB-Installation unter Windows sieht der Pfad beispielsweise so aus:

``C:\Program Files\MariaDB 11.0\bin\mysqldump.exe``

Bei einer sauberen MySQL-Installation unter Windows sieht der Pfad beispielsweise so aus:

``C:\Program Files\MySQL\MySQL Server 8.0\bin\mysqldump.exe``

Bei einer Installation von XAMPP unter Windows sieht der Pfad beispielsweise so aus:

``C:\xampp\mysql\bin\mysqldump.exe``

Unter Linux wird der Pfad zu mysqldump in der Regel automatisch gefunden. Es kann also einfach folgendes eingetragen werden:

``mysqldump``
## webhook-url
Hierbei handelt es sich um die URL, die du bei der Erstellung eines Webhooks für einen Channel in Discord direkt kopieren kannst.

(Ein Tutorial über Discord Webhooks findest du [hier](https://support.discord.com/hc/en-us/articles/228383668-Intro-to-Webhooks).)

So eine URL sieht beispielsweise so aus:

``https://discord.com/api/webhooks/9939128298523341122/u7kJk8QBG62HbRCIkaU4SkAympwis-EVXqq5Q8gQymnCOoN-ib7aCwolGx0_-Eb_2MY0``
## backup-interval
Hier kannst du den Intervall einstellen, in dem Backups durchgeführt werden.

Beispielsweise ``24h`` für ein Backup am Tag oder ``12h`` für ein Backup alle 12 Stunden.
## initial-backup-time
Hier kannst du die Tageszeit einstellen, zu der das erste Backup nach Start der Software durchgeführt wird.

Trägst du hier beispielsweise ``06:00:00`` (also 6 Uhr morgens) ein und startest die Software um 3 Uhr morgens, wird das erste Backup um 6 Uhr ausgeführt und das nächste nach Ablauf von ``backup-interval`` und so weiter.

Startest du die Software aber erst um 7 Uhr morgens, wird ein ``backup-interval`` auf die eingetragene Zeit von 6 Uhr aufgerechnet. Hast du also einen ``backup-interval`` von 12 Stunden konfiguriert, wird das erste Backup um 18 Uhr durchgeführt und jedes weitere dann wieder nach Ablauf von ``backup-interval``.
## databases
Hier trägst du die Namen aller Datenbanken ein, die gesichert werden sollen. Diese werden durch Komma (``,``) getrennt. Bitte beachte, dass hier keine Leerzeichen erwünscht sind und diese zu Fehlern führen können!

Wenn du lediglich eine Datenbank sichern möchtest, sind keine Kommas nötig.
## host, user, password
Hier kannst du die Zugangsdaten zu deinem MySQL-/MariaDB-Server konfigurieren.

``host`` kann dabei entweder eine IP-Adresse oder ein Hostname sein. In den meisten Fällen sollte die Software aber ja ohnehin auf dem gleichen Computer wie der MySQL-/MariaDB-Server ausgeführt werden, wofür der Wert ``localhost`` vollkommen ausreicht.

``user`` muss genau wie ``host`` in jedem Fall angegeben werden und ``password`` kann leer gelassen werden, falls ein Benutzer ohne Passwort verwendet wird.
