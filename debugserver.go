package debugserver

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var storage = NewStorage()
var emptyResponse = []byte("[]\n")

func API() http.Handler {
	router := httprouter.New()

	router.DELETE("/buckets/:id", record)
	router.GET("/buckets/:id", record)
	router.PATCH("/buckets/:id", record)
	router.PUT("/buckets/:id", record)
	router.POST("/buckets/:id", record)

	router.DELETE("/reports/:id", del)
	router.GET("/reports/:id", show)

	return router
}

type Request struct {
	ContentLength int64               `json:"content_length"`
	Body          string              `json:"body"`
	Method        string              `json:"method"`
	URL           string              `json:"url"`
	Headers       map[string][]string `json:"headers"`
}

func record(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	bodyReader := http.MaxBytesReader(w, r.Body, 1<<20)
	defer bodyReader.Close()
	var body bytes.Buffer

	if _, err := io.Copy(&body, bodyReader); err != nil {
		http.Error(w, "unable to read body", http.StatusInternalServerError)
		return
	}

	storage.Add(id, Request{
		URL:           r.URL.Path,
		Method:        r.Method,
		Headers:       r.Header,
		Body:          base64.StdEncoding.EncodeToString(body.Bytes()),
		ContentLength: r.ContentLength,
	})
}

func show(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	requests := storage.Get(p.ByName("id"))

	if requests == nil {
		w.Write(emptyResponse)
		return
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")

	if err := encoder.Encode(requests); err != nil {
		log.Printf("Failed to serialize payload: %s", err)
	}
}

func del(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	storage.Del(p.ByName("id"))
}
