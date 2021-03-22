package fluxgen

import (
	"log"
	"os"
	"path"
)

func InitProject(projectName string) {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Unable to retrieve current working directory")
	}

	if err := os.Mkdir(path.Join(currentDir, projectName), 744); err != nil {
		if mkdirErr, ok := err.(*os.PathError); ok {
			log.Fatal("Unable to create folder at path: " + mkdirErr.Path)
		}
	}

	if err := os.Chdir(projectName); err != nil {
		if chdirErr, ok := err.(*os.PathError); ok {
			log.Fatal("Unable to create folder at path: " + chdirErr.Path)
		}
	}

	for _, folderName := range []string{"templates", "static", "pages", "_site"} {
		if err := os.Mkdir(folderName, 744); err != nil {
			if mkdirErr, ok := err.(*os.PathError); ok {
				log.Fatal("Unable to create folder at path: " + mkdirErr.Path)
			}
		}
	}
}
