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
	"sync"
)

// Use https://github.com/gobuffalo/plush for templating

func FluxBuild() {
	wg := sync.WaitGroup{}
	wg.Add(2)

	FluxClean()

	fluxConfig, err := parseFluxConfig(ConfigFile)
	if err != nil {
		errLogger.Fatalf("FluxConfig parsing failed\n\t[%v]", err)
	}

	resources, err := loadResources(CSSDir, AssetsDir)
	if err != nil {
		errLogger.Fatalf("Resource loading Failed\n\t[%v]", err)
	}

	pageList, postList, err := parsePages(&fluxConfig, &resources)
	if err != nil {
		errLogger.Fatalf("Pages parsing failed\n\t[%v]", err)
	}

	By(descendingOrderByDate).Sort(postList)

	err = parseHTMLTemplates(pageList, postList)
	if err != nil {
		errLogger.Fatalf("HTML parsing failed\n\t[%v]", err)
	}

	err = processPageAssets(PostsDir)
	if err != nil {
		errLogger.Fatalf("Processing page Assets failed\n\t[%v]", err)
	}

	go func() {
		err := processStaticFolders(CSSDir, &fluxConfig, &wg)
		if err != nil {
			errLogger.Fatalf("Processing CSS Files failed\n\t[%v]", err)
		}
	}()

	go func() {
		err := processStaticFolders(AssetsDir, &fluxConfig, &wg)
		if err != nil {
			errLogger.Fatalf("Processing Static Files failed\n\t[%v]", err)
		}
	}()

	wg.Wait()
	printMsg("Done", "party")
	//if fluxConfig["minify_html"] == true {
	//	minifyHtml()
	//}
}

func loadResources(path ...string) (resources, error) {
	resources := resources{}
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
				return resources, err
			}
		}
	}
	return resources, nil
}

func parseFluxConfig(path string) (fluxConfig, error) {
	fluxConfig := make(fluxConfig)
	configFile, err := ioutil.ReadFile(path)
	if err != nil {
		return fluxConfig, err
	}
	err = json.Unmarshal(configFile, &fluxConfig)
	if err != nil {
		return fluxConfig, err
	}
	printMsg("Parsed Config File", "tick")
	return fluxConfig, nil
}

func parsePages(config *fluxConfig, resources *resources) (pages, pages, error) {
	var allPages pages
	var postList pages
	if _, err := os.Stat(PostsDir); err == nil {
		err := filepath.WalkDir(PostsDir, func(path string, d fs.DirEntry, err error) error {
			if !d.IsDir() && filepath.Ext(path) == ".md" {
				mdPage, _ := parseMarkdown(path, config, resources)
				allPages = append(allPages, mdPage)
				postList = append(postList, mdPage)
			}
			return nil
		})
		if err != nil {
			return allPages, postList, err
		}

		dirContent, err := ioutil.ReadDir(".")
		if err != nil {
			return allPages, postList, err
		}

		for _, f := range dirContent {
			if !f.IsDir() && filepath.Ext(f.Name()) == ".html" {
				htmlPage, _ := parseHTML(f.Name(), config, resources)
				allPages = append(allPages, htmlPage)
			}
		}
		printMsg("Parsed pages", "tick")
	}

	return allPages, postList, nil
}

func parseHTMLTemplates(allPages, posts pages) error {
	for _, p := range allPages {
		p.PostsList = &posts
		buffer, err := p.applyTemplate()
		if err != nil {
			return err
		}

		fileWritePath := createFileWritePath(p.filename, p.Href)
		err = createFileWriteDir(fileWritePath)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(filepath.Join(fileWritePath, "index.html"), []byte(buffer), 0744)
		if err != nil {
			return err
		}
	}
	return nil
}

func processPageAssets(dir string) error {
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
			return err
		}
		printMsg("Processed page Assets", "tick")
	}
	return nil
}

func processStaticFolders(filePath string, fc *fluxConfig, wg *sync.WaitGroup) error {
	defer wg.Done()
	if _, err := os.Stat(filePath); err == nil {
		err = filepath.WalkDir(filePath, func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				err := os.MkdirAll(filepath.Join(SiteDir, path), 0744)
				if err != nil {
					return err
				}
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
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
		printMsg(fmt.Sprintf("Processed %s/", filePath), "tick")
	}
	return nil
}

//func minifyHtml() {
//	m := minify.New()
//	m.AddFunc("text/html", html.Minify)
//	err := filepath.WalkDir(SiteDir, func(path string, d fs.DirEntry, err error) error {
//		if !d.IsDir() && filepath.Ext(path) == ".html" {
//			fr, err := os.OpenFile(path, os.O_RDWR, 0755)
//			if err != nil {
//				return err
//			}
//			defer fr.Close()
//
//			fw, err := os.OpenFile(path, os.O_WRONLY, 0755)
//			if err != nil {
//				return err
//			}
//			defer fw.Close()
//
//			if err := m.Minify("text/html", fw, fr); err != nil {
//				return err
//			}
//		}
//		return nil
//	})
//	if err != nil {
//		log.Fatalf("error minifying html! - %v", err)
//	}
//}
