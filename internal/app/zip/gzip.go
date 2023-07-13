package zip

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func newGzipWriter(w http.ResponseWriter) *gzipWriter {
	return &gzipWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (c *gzipWriter) Header() http.Header {
	return c.w.Header()
}

func (c *gzipWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		contentType := c.Header().Get("Content-Type")
		if headerCheck(contentType, "application/json") || headerCheck(contentType, "text/html") {
			c.w.Header().Set("Content-Encoding", "gzip")
		}
	}
	c.w.WriteHeader(statusCode)
}

func (c *gzipWriter) Write(p []byte) (int, error) {
	if headerCheck(c.Header().Get("Content-Encoding"), "gzip") {
		return c.zw.Write(p)
	}
	return c.w.Write(p)
}

func (c *gzipWriter) Close() error {
	return c.zw.Close()
}

type gzipReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newGzipReader(r io.ReadCloser) (*gzipReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &gzipReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c gzipReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *gzipReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

func GzipMiddleware(h http.Handler) http.Handler {
	compFn := func(w http.ResponseWriter, r *http.Request) {
		if headerCheck(r.Header.Get("Content-Encoding"), "gzip") {
			cr, err := newGzipReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close()
		}
		ow := w
		if headerCheck(r.Header.Get("Accept-Encoding"), "gzip") {
			cw := newGzipWriter(w)
			ow = cw
			defer cw.Close()
		}

		h.ServeHTTP(ow, r)

	}
	return http.HandlerFunc(compFn)
}

func headerCheck(head, par string) bool {
	options := strings.Split(head, ",")
	for _, option := range options {
		option = strings.TrimSpace(option)
		if option == par {
			return true
		}
	}
	return false
}
