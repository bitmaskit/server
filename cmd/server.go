package main

import (
	"flag"
	"log"
	"net/http"
	"time"
)

var (
	port      string
	directory string
)

func main() {
	flag.StringVar(&port, "p", "8080", "port to serve on")
	flag.StringVar(&directory, "d", ".", "the directory of static file to host")
	flag.Parse()

	http.Handle("/", noCache(http.FileServer(http.Dir(directory))))

	log.Printf("Serving %s on http://localhost:%s/\n", directory, port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

var noCacheHeaders = map[string]string{
	"Expires":         time.Unix(0, 0).Format(time.RFC1123),
	"Cache-Control":   "no-cache, private, max-age=0",
	"Pragma":          "no-cache",
	"X-Accel-Expires": "0",
}

var etagHeaders = []string{"ETag", "If-Modified-Since", "If-Match", "If-None-Match", "If-Range", "If-Unmodified-Since"}

// noCache middleware sets no-caching headers
func noCache(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Delete any ETag headers that may have been set
		for _, v := range etagHeaders {
			if r.Header.Get(v) != "" {
				r.Header.Del(v)
			}
		}

		// Set our NoCache headers
		for k, v := range noCacheHeaders {
			w.Header().Set(k, v)
		}

		handler.ServeHTTP(w, r)
	})
}
