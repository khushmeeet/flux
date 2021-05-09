package fluxgen

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func FluxBuild() {
	FluxClean()
	fluxConfig := parseFluxConfig(ConfigFile)
	resources := loadResources(CSSDir, AssetsDir)
	pageList, postList := parsePages(&fluxConfig, &resources)
	By(descendingOrderByDate).Sort(postList)
	parseHTMLTemplates(pageList, postList)
	processPageAssets(PostsDir)
	processStaticFolders(CSSDir, &fluxConfig)
	processStaticFolders(AssetsDir, &fluxConfig)
}

func loadResources(path ...string) Resources {
	resources := Resources{}
	for _, p := range path {
		if _, err := os.Stat(p); err == nil {
			err := filepath.WalkDir(p, func(path string, d fs.DirEntry, err error) error {
				if !d.IsDir() {
					if filepath.Ext(d.Name()) == ".scss" {
						fileNameNoExt := strings.TrimSuffix(path, ".scss")
						resources[filepath.Base(d.Name())] = filepath.Join("/", fileNameNoExt+".css")
					} else {
						resources[filepath.Base(d.Name())] = filepath.Join("/", path)
					}
				}
				return nil
			})
			if err != nil {
				log.Fatalf("Error walking (%v) - %v", p, err)
			}
		}
	}
	return resources
}

func parseFluxConfig(path string) FluxConfig {
	fluxConfig := make(FluxConfig)
	configFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("[Error Reading (%v)] - %v", path, err)
	}
	err = json.Unmarshal(configFile, &fluxConfig)
	if err != nil {
		log.Fatalf("[Error Unmarshalling (%v)] - %v", path, err)
	}
	return fluxConfig
}

func parsePages(config *FluxConfig, resources *Resources) (Pages, Pages) {
	var pageList Pages
	var postList Pages
	if _, err := os.Stat(PostsDir); err == nil {
		err := filepath.WalkDir(PostsDir, func(path string, d fs.DirEntry, err error) error {
			if !d.IsDir() && filepath.Ext(path) == ".md" {
				mdPage := parseMarkdown(path, config, resources)
				pageList = append(pageList, mdPage)
				postList = append(postList, mdPage)
			}
			return nil
		})
		if err != nil {
			log.Fatalf("[Error Walking (%v)] - %v", PostsDir, err)
		}

		dirContent, err := ioutil.ReadDir(".")
		if err != nil {
			log.Fatalf("[Error Reading (%v)] - %v", ".", err)
		}

		for _, f := range dirContent {
			if !f.IsDir() && filepath.Ext(f.Name()) == ".html" {
				htmlPage := parseHTML(f.Name(), config, resources)
				pageList = append(pageList, htmlPage)
			}
		}
	}

	return pageList, postList
}

func parseHTMLTemplates(pages Pages, posts Pages) {
	for _, p := range pages {
		p.postsList = &posts
		buffer, err := p.applyTemplate()
		if err != nil {
			log.Fatalf("[Error Applying template to Page] - %v", err)
		}

		fileWritePath := createFileWritePath(p.filename, p.href)
		createFileWriteDir(fileWritePath)

		err = ioutil.WriteFile(filepath.Join(fileWritePath, "index.html"), []byte(buffer), 0744)
		if err != nil {
			log.Fatalf("[Error Writing File (%v)] - %v", p.href, err)
		}
		fmt.Printf("Writing File: %v\n", p.href+p.oldExtention)
	}
}

func processPageAssets(dir string) {
	if _, err := os.Stat(dir); err == nil {
		err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if !d.IsDir() && filepath.Ext(path) != ".md" {
				err := copyFile(path, filepath.Join(SiteDir, path))
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			log.Fatalf("error copying file [%v]", err)
		}
	}
}

func processStaticFolders(filePath string, fc *FluxConfig) {
	if _, err := os.Stat(filePath); err == nil {
		err = filepath.WalkDir(filePath, func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				err := os.MkdirAll(filepath.Join(SiteDir, path), 0744)
				if err != nil {
					return err
				}
				fmt.Printf("Creating Folder: %v\n", path)
			} else {
				dstFilePath := filepath.Join(SiteDir, path)

				if filepath.Ext(d.Name()) == ".scss" {
					src, err := os.Open(path)
					if err != nil {
						log.Fatalf("error opening file: %v", err)
					}
					defer src.Close()

					dst, err := os.Create(filepath.Join(SiteDir, strings.TrimSuffix(path, filepath.Ext(path))+".css"))
					if err != nil {
						log.Fatalf("error opening file: %v", err)
					}
					defer dst.Close()

					c := createSassCompiler(src, dst, fc)
					c.compileSass()
				} else {
					err := copyFile(path, dstFilePath)
					if err != nil {
						return err
					}
					fmt.Printf("Copying File: %v\n", path)
				}
			}
			return nil
		})
		if err != nil {
			log.Fatalf("[Error Walking (%v)] - %v", filePath, err)
		}
	} else {
		fmt.Printf("\"%v\" does not exists... Skipping", filePath)
	}
}
