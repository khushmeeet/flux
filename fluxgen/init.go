package fluxgen

import (
	"fmt"
	"log"
	"os"
	"path"
)

func InitProject(projectName string) {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Unable to retrieve current working directory")
	}

	fmt.Println("Creating root folder: " + projectName)
	if err := os.Mkdir(path.Join(currentDir, projectName), 0777); err != nil {
		if mkdirErr, ok := err.(*os.PathError); ok {
			log.Fatal("Unable to create folder at path: " + mkdirErr.Path)
		}
	}

	for _, folderName := range []string{TemplatesFolder, StaticFolder, PagesFolder, SiteFolder} {
		fmt.Println("Creating sub-folder: " + folderName)
		if err := os.Mkdir(path.Join(currentDir, projectName, folderName), 0777); err != nil {
			if mkdirErr, ok := err.(*os.PathError); ok {
				log.Fatal("Unable to create folder at path: " + mkdirErr.Path)
			}
		}
	}

	fmt.Println("Project scaffolding complete!")
}
