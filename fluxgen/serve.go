package fluxgen

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func WatchAndServe(port string, watch bool) {
	serve(port, "Running http server at")

	if watch {
		done := make(chan bool)
		w, err := fsnotify.NewWatcher()
		if err != nil {
			panic(err)
		}
		defer w.Close()

		err = filepath.WalkDir(SiteDir, func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				return w.Add(path)
			}
			return nil
		})
		if err != nil {
			return
		}

		go func() {
			for {
				select {
				case e, ok := <-w.Events:
					if !ok {
						return
					}
					if e.Op == fsnotify.Write {
						fmt.Println("File changed: ", strings.TrimSuffix(e.Name, "~"))
						FluxBuild()
						serve(port, "Restarting http server at")
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
	}
}

func serve(port, message string) {
	http.Handle("/", http.FileServer(http.Dir(SiteDir)))
	fmt.Printf("%s :%s...\n", message, port)
	err := http.ListenAndServe(fmt.Sprintf(":%v", port), Logger(os.Stderr, http.DefaultServeMux))
	if err != nil {
		log.Fatal("Unable to start http server!")
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
