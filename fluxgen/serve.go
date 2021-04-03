package fluxgen

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type responseObserver struct {
	http.ResponseWriter
	status      int
	written     int64
	wroteHeader bool
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
		logger.Printf("%s - [%s] %q",
			addr,
			time.Now().Format("02/Jan/2006:15:04:05 -0700"),
			fmt.Sprintf("%s %s %s", r.Method, r.URL, r.Proto),
		)
	})
}

func FluxServe() {
	http.Handle("/", http.FileServer(http.Dir(SiteFolder)))
	if err := http.ListenAndServe(":5050", Logger(os.Stderr, http.DefaultServeMux)); err != nil {
		log.Fatal("Unable to start http static server!")
	}
}
