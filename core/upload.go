package core

import (
	"log"

	"mime"
	"path/filepath"
	"time"

	"regexp"

	"golang.org/x/net/context"
	"google.golang.org/api/drive/v3"
)

func UploadAll(srv *drive.Service, ctx context.Context, stdLogger *log.Logger) {
	backupDirId, err := createBackupRootDir(srv, ctx)

	if err != nil {
		log.Fatalf("Unable to create Backup Dir: %v", err)
	}

	stdLogger.Printf("[LOG] %v | Searching for Skype main.db history files \n", time.Now().Format("01/02/2006 - 15:04:05"))

	dbFiles := searchForSkypeDbs()

	for _, fullPath := range dbFiles {
		re := regexp.MustCompile("Skype/(.*?)/main.db")
		rm := re.FindStringSubmatch(fullPath)

		userDirId, err := createSubDir(srv, ctx, rm[1], backupDirId)
		if err != nil {
			continue
		}

		stdLogger.Printf("[LOG] %v | Uploading.... \n", time.Now().Format("01/02/2006 - 15:04:05"))
		f := uploadDbFile(srv, ctx, fullPath, userDirId)
		stdLogger.Printf("[LOG] %v | Uploaded %s - %s, Download link: %s\n", time.Now().Format("01/02/2006 - 15:04:05"), f.Id, f.Name, f.WebContentLink)
	}
}

func uploadDbFile(srv *drive.Service, ctx context.Context, filePath string, parentDirId string) *drive.File {
	srcFile, srcFileInfo, err := OpenFile(filePath)
	if err != nil {
		log.Fatalf("Unable to open file: %v", err)
	}

	// Close file on function exit
	defer srcFile.Close()

	//Create date dir
	dateDirId, err := createSubDir(srv, ctx, time.Now().Format("01022006"), parentDirId)
	if err != nil {
		log.Fatalf("Unable to create Date Dir: %v", err)
	}

	// Initiate Backup file
	dbFile := &drive.File{Name: srcFileInfo.Name()}
	dbFile.MimeType = mime.TypeByExtension(filepath.Ext(dbFile.Name))
	dbFile.Parents = []string{dateDirId}

	f, err := srv.Files.Create(dbFile).Fields("id", "name", "size", "webContentLink").Context(ctx).Media(srcFile).Do()
	if err != nil {
		log.Fatalf("Unable to upload file: %v", err)
	}

	return f
}

//Create SkypeBackup dir in Drive or get existing dir ID
func createBackupRootDir(srv *drive.Service, ctx context.Context) (string, error) {
	r, err := srv.Files.List().PageSize(1).Fields("files(id, name)").Q("mimeType='application/vnd.google-apps.folder' and name='SkypeBackups' and trashed = false").Do()
	if err != nil {
		log.Fatalf("Unable to get GDrive folders list: %v", err)
		return "", err
	}

	var id string

	if len(r.Files) > 0 {
		id = r.Files[0].Id
		err = nil
	} else {
		// Initiate Root Dir
		backupDir := &drive.File{Name: "SkypeBackups", MimeType: "application/vnd.google-apps.folder"}
		id, err = createDir(srv, ctx, backupDir)
	}

	return id, err
}

//Create subdir in Drive or get existing subdir ID
func createSubDir(srv *drive.Service, ctx context.Context, name string, parentId string) (string, error) {
	r, err := srv.Files.List().PageSize(1).Fields("files(id, name)").Q("mimeType='application/vnd.google-apps.folder' and name='" + name + "' and '" + parentId + "' in parents and trashed = false").Do()
	if err != nil {
		log.Fatalf("Unable to get GDrive folders list: %v", err)
		return "", err
	}

	var id string

	if len(r.Files) > 0 {
		id = r.Files[0].Id
		err = nil
	} else {
		subDir := &drive.File{Name: name, MimeType: "application/vnd.google-apps.folder"}
		subDir.Parents = []string{parentId}
		id, err = createDir(srv, ctx, subDir)
	}

	return id, err
}

func createDir(srv *drive.Service, ctx context.Context, dir *drive.File) (string, error) {
	r, err := srv.Files.Create(dir).Fields("id").Context(ctx).Do()
	if err != nil {
		log.Fatalf("Unable to create Backup Dir: %v", err)
		return "", err
	} else {
		return r.Id, nil
	}
}
