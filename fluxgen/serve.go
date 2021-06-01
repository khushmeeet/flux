package fluxgen

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/muesli/termenv"
)

var orange = termenv.ColorProfile().Color("#f9b208")

func WatchAndServe(port string, watch bool) {
	if watch {
		done := make(chan bool)
		w, err := fsnotify.NewWatcher()
		if err != nil {
			panic(err)
		}
		defer w.Close()

		err = filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() && d.Name() != SiteDir && d.Name() != ".git" {
				return w.Add(path)
			}
			return nil
		})
		if err != nil {
			return
		}

		FluxBuild()
		serve("HTTP Server running at", port, true)

		go func() {
			for {
				select {
				case e, ok := <-w.Events:
					if !ok {
						return
					}
					if e.Op == fsnotify.Write {
						fmt.Println("...")
						changedFile := termenv.String("File changed: " + strings.TrimSuffix(e.Name, "~")).Foreground(orange).String()
						fmt.Println(changedFile)
						FluxBuild()
						fmt.Printf("HTTP Server running at :%s\n", port)
					}
				case err, ok := <-w.Errors:
					if !ok {
						return
					}
					fmt.Println("err:", err)
				}
			}
		}()
		<-done
	} else {
		FluxBuild()
		serve("Running http server at", port, false)
	}
}

func serve(message, port string, async bool) {
	http.Handle("/", http.FileServer(http.Dir(SiteDir)))
	fmt.Printf("%s :%s...\n", message, port)
	if async == true {
		go func() {
			err := http.ListenAndServe(fmt.Sprintf(":%v", port), Logger(os.Stderr, http.DefaultServeMux))
			if err != nil {
				log.Fatalf("Unable to start http server %v", err)
			}
		}()
	} else {
		err := http.ListenAndServe(fmt.Sprintf(":%v", port), Logger(os.Stderr, http.DefaultServeMux))
		if err != nil {
			log.Fatalf("Unable to start http server %v", err)
		}
	}
}

func Logger(out io.Writer, h http.Handler) http.Handler {
	logger := log.New(out, "", 0)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		o := &responseObserver{ResponseWriter: w}
		h.ServeHTTP(o, r)
		addr := r.RemoteAddr
		if i := strings.LastIndex(addr, ":"); i != -1 {
			addr = addr[:i]
		}
		logger.Printf("[%s] %q",
			time.Now().Format("02/Jan/2006:15:04:05 -0700"),
			fmt.Sprintf("%s %s %s", r.Method, r.URL, r.Proto),
		)
	})
}
