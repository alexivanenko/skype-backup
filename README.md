# Skype History Backup
## Simple golang driven tool for backup Skype history file in GoogleDrive 

### Requirements

skype-backup version 0.1 requires Go >= 1.6

##### Installation

```sh
$ go get github.com/alexivanenko/skype-backup/...
$ cd skype-backup
$ make
```

1. Go to the https://console.developers.google.com/start/api?id=drive
2. Use this wizard to create or select a project in the Google Developers Console and automatically turn on the API. Click Continue, then Go to credentials.
3. At the top of the page, select the OAuth consent screen tab. Select an Email address, enter a Product name if not already set, and click the Save button.
4. Select the Credentials tab, click the Create credentials button and select OAuth client ID.
5. Select the application type Other, enter the name "Skype-Backup", and click the Create button.
6. Click OK to dismiss the resulting dialog.
7. Click the file_download (Download JSON) button to the right of the client ID.
8. Move this file to your working directory(where the bin file located) and rename it client_secret.json

#### cache.json - store path to Skype dir for quicker main.db files search, please run `chmod 0777 cache.json` in your working dir
