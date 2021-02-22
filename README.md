# SF-Downloader
![GitHub](https://img.shields.io/github/license/omski/SF-Downloader?style=for-the-badge) ![GitHub all releases](https://img.shields.io/github/downloads/omski/SF-Downloader/total?style=for-the-badge) ![GitHub release (latest by date including pre-releases)](https://img.shields.io/github/v/release/omski/SF-Downloader?include_prereleases&style=for-the-badge) ![GitHub Workflow Status](https://img.shields.io/github/workflow/status/omski/SF-Downloader/Go?style=for-the-badge) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/omski/SF-Downloader?style=for-the-badge)

SF-Downloader dient für den vereinfachten Umgang mit vielen Dateien in Schoolfox - FoxDrive

## Beschreibung

FD-Downloader ist ein Kommandozeilenprogramm geschrieben in Go <https://golang.org/>.
Bei der Veröffentlichung ([Releases](<https://github.com/omski/SF-Downloader/releases>)) neuer Versionen werden ausführbare Versionen für die Betriebssysteme Windows, MacOS und Linux erstellt.

## Anleitung

* Nach dem Start des Programm fordert es zur Eingabe des Schoolfox Benutzernamens und Passworts auf.</br>
![Login](./assets/login.png)
* Anschließend wird eine Liste des 'Inventars' dargestellt.
Das Programm fordert zur Auswahl eines Inventareintrags aus.</br>
![Select Pupil](./assets/select_pupil.png)
Es werden nun alle Einträge des FoxDrive Hauptverzeichnisses (root) dargestellt.
Man kann nun in Unterverzeichnisse navigieren oder das aktuell dargestellte Verzeichnis auswählen.</br>
![Select Folder](./assets/select_folder.png)
* Nachdem ein Verzeichnis ausgewählt wurde besteht die Option:</br>
1 - die Dateien des ausgewählten Verzeichnisses herunterzuladen</br>
2 - die Dateien des ausgewählten Verzeichnisses und aller Unterverzeichnisse herunterzuladen</br>
3 - die Dateien des ausgewählten Verzeichnisses herunterzuladen und nach erfolgreichem Download auf Foxdrive zu löschen</br>
4 - die Dateien des ausgewählten Verzeichnisses und aller Unterverzeichnisse herunterzuladen und nach erfolgreichem Download auf Foxdrive zu löschen</br>
![Select Command](./assets/select_command.png)
* Nach Ausführung des Kommandos wird das Programm beendet.</br>
![End](./assets/end.png)
## geplante Features

Zyklisches Ausführen des gewählten Kommandos um das gewählte Verzeichnis über einen längeren Zeitraum automatisch zu synchronisieren.
