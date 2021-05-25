package fluxgen

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
)

var errLogger *log.Logger

func init() {
	errLogger = log.New(os.Stderr, "ERROR: ", 0)
}

func FluxInit(projectName string) {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Unable to retrieve current working directory")
	}

	fc := fluxConfig{
		"site": fluxConfig{
			"title":   projectName,
			"email":   "hello@flux.com",
			"twitter": fmt.Sprintf("@%v", projectName),
			"github":  projectName,
		},
		"minify_css": false,
		//"minify_html": false,
		//"rss_feed": false,
	}

	fmt.Println("Creating root folder: " + projectName)
	if err := os.Mkdir(path.Join(currentDir, projectName), 0777); err != nil {
		if mkdirErr, ok := err.(*os.PathError); ok {
			log.Fatal("Unable to create folder at path: " + mkdirErr.Path)
		}
	}

	fmt.Println("Creating project config file: config.json")
	fcJson, err := json.Marshal(fc)
	err = ioutil.WriteFile(path.Join(currentDir, projectName, ConfigFile), fcJson, 0644)
	if err != nil {
		log.Fatal("Unable to write data to Config File")
	}

	for _, folderName := range []string{TemplatesDir, CSSDir, AssetsDir, PostsDir, SiteDir} {
		fmt.Println("Creating sub-folder: " + folderName)
		if err := os.Mkdir(path.Join(currentDir, projectName, folderName), 0777); err != nil {
			if mkdirErr, ok := err.(*os.PathError); ok {
				log.Fatal("Unable to create folder at path: " + mkdirErr.Path)
			}
		}
	}

	fmt.Println("Project scaffolding complete!")
}
