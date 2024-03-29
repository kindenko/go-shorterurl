package zip

import (
	"compress/gzip"
	"log"
	"net/http"
)

func MiddlewareCompressGzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(`Content-Encoding`) == `gzip` {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				log.Printf("Gzip middleware: error while reading body")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			r.Body = gz
			defer gz.Close()
		}

		next.ServeHTTP(w, r)
	})
}
