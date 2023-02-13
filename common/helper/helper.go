package helper

import (
	"os"
	"path/filepath"

	"github.com/vault-thirteen/auxie/zipper"
)

const (
	DataFolder = "data"
	DbFolder   = "db"
	TempFolder = "tmp"
)

func UnpackDbFile(archivePath string) (dbFilePath string, err error) {
	err = createTemporaryDbDataFolder()
	if err != nil {
		return "", err
	}

	tmpFolderPath := filepath.Join(DataFolder, DbFolder, TempFolder)

	return zipper.UnpackZipFile(archivePath, tmpFolderPath)
}

func createTemporaryDbDataFolder() (err error) {
	folderPath := filepath.Join(DataFolder, DbFolder, TempFolder)
	return os.MkdirAll(folderPath, 0777)
}

func DeleteTemporaryDataFolders() (err error) {
	err = deleteTemporaryDbDataFolder()
	if err != nil {
		return err
	}

	return nil
}

func deleteTemporaryDbDataFolder() (err error) {
	folderPath := filepath.Join(DataFolder, DbFolder, TempFolder)
	return os.RemoveAll(folderPath)
}
