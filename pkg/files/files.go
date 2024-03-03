package files

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func WriteHTML(data string, filepath string) error {
	// Write the JSON data to a file
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	_, err = file.Write([]byte(data))
	if err != nil {
		return fmt.Errorf("error writing HTML to file: %v", err)
	}

	return nil
}

func CreateFolderStorage() {
	basePath := "./storage"
	err := os.MkdirAll(basePath, 0755)
	if err != nil {
		return
	}
}

func RemoveFolders() {
	basePath := "./storage"
	err := os.RemoveAll(basePath)
	if err != nil {
		return
	}
}

func CreateZipArchive(sourceFolder string, zipFilePath string) error {
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return fmt.Errorf("error creating the archive: %w", err)
	}
	defer func() {
		if derr := zipFile.Close(); derr != nil {
			fmt.Println("Error closing zip file:", derr)
		}
	}()

	zipWriter := zip.NewWriter(zipFile)
	defer func() {
		if derr := zipWriter.Close(); derr != nil {
			fmt.Println("Error closing zip writer:", derr)
		}
	}()

	err = filepath.Walk(sourceFolder, func(filePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error while crawling files: %w", err)
		}

		if fileInfo.IsDir() {
			return nil
		}

		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("error opening file %s: %w", filePath, err)
		}
		defer func() {
			if derr := file.Close(); derr != nil {
				fmt.Println("Error closing file:", derr)
			}
		}()

		zipEntry, err := zipWriter.Create(filepath.Base(filePath))
		if err != nil {
			return fmt.Errorf("error creating a file in the archive: %w", err)
		}

		if _, err = io.Copy(zipEntry, file); err != nil {
			return fmt.Errorf("error copying file to an archive: %w", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
