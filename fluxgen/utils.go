package fluxgen

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func copyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}

func getMapValue(m map[string]interface{}, k string) string {
	if val, ok := m[k]; ok {
		return val.(string)
	}
	return ""
}

func createFileWritePath(fileName string, filePath string) string {
	fileWritePath := ""
	if fileName == "index.html" {
		fileWritePath = SiteFolder
	} else {
		fileWritePath = filepath.Join(SiteFolder, filePath)
	}
	return fileWritePath
}

func createFileWriteDir(filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		_ = os.MkdirAll(filePath, 0744)
	} else {
		fmt.Println("Directory already exists")
	}
}
