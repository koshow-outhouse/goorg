package logger

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type (
	Response struct {
		StatusCode int
		Header     http.Header
	}
	responseWriter struct {
		http.ResponseWriter
		resp Response
	}
)

func (w *responseWriter) Write(buf []byte) (int, error) {
	n, e := w.ResponseWriter.Write(buf)
	if w.resp.StatusCode == 0 {
		w.resp.StatusCode = http.StatusOK
	}
	return n, e
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.resp.StatusCode = statusCode
}

func (w *responseWriter) CloseNotify() <-chan bool {
	if closeNotifier, ok := w.ResponseWriter.(http.CloseNotifier); ok {
		return closeNotifier.CloseNotify()
	}
	return nil
}

func (w *responseWriter) Flush() {
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

func WrapResponseWriter(w http.ResponseWriter) (http.ResponseWriter, *Response) {
	rw := responseWriter{
		ResponseWriter: w,
		resp: Response{
			Header: w.Header(),
		},
	}

	return &rw, &rw.resp
}

func Log(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		rAddr := r.RemoteAddr
		method := r.Method
		path := r.URL.Path
		log.Printf("Remote: %s [%s] %s", rAddr, method, path)
		w, resp := WrapResponseWriter(w)
		if r.Method == http.MethodOptions {
			return
		}
		h(w, r, p)
		log.Println(fmt.Sprintf("Status: %v", resp.StatusCode))
	}
}

func LogHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rAddr := r.RemoteAddr
		method := r.Method
		path := r.URL.Path
		log.Printf("Remote: %s [%s] %s", rAddr, method, path)
		w, resp := WrapResponseWriter(w)
		if r.Method == http.MethodOptions {
			return
		}
		h.ServeHTTP(w, r)
		log.Println(fmt.Sprintf("Status: %v", resp.StatusCode))
	})
}

func LogHandlerFunc(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rAddr := r.RemoteAddr
		method := r.Method
		path := r.URL.Path
		log.Printf("Remote: %s [%s] %s", rAddr, method, path)
		w, resp := WrapResponseWriter(w)
		if r.Method == http.MethodOptions {
			return
		}
		h.ServeHTTP(w, r)
		log.Println(fmt.Sprintf("Status: %v", resp.StatusCode))
	}
}