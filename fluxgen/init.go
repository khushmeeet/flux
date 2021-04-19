package fluxgen

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type fluxConfig struct {
	SiteTitle       string `json:"site_title"`
	Email           string `json:"email"`
	TwitterUsername string `json:"twitter_username"`
	GithubUsername  string `json:"github_username"`
}

func FluxInit(projectName string) {
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

	fmt.Println("Creating project config file: config.json")
	fcJson, err := json.Marshal(fluxConfig{
		SiteTitle:       projectName,
		Email:           "hello@flux.com",
		TwitterUsername: "@" + projectName,
		GithubUsername:  projectName,
	})
	err = ioutil.WriteFile(path.Join(currentDir, projectName, ConfigFile), fcJson, 0644)
	if err != nil {
		log.Fatal("Unable to write data to Config File")
	}

	for _, folderName := range []string{TemplatesFolder, CSSFolder, AssetsFolder, PagesFolder, SiteFolder} {
		fmt.Println("Creating sub-folder: " + folderName)
		if err := os.Mkdir(path.Join(currentDir, projectName, folderName), 0777); err != nil {
			if mkdirErr, ok := err.(*os.PathError); ok {
				log.Fatal("Unable to create folder at path: " + mkdirErr.Path)
			}
		}
	}

	fmt.Println("Project scaffolding complete!")
}
