package debugserver

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/julienschmidt/httprouter"
)

var (
	mu      sync.Mutex
	storage = make(map[string][]Request)
	records = &List{}
)

var emptyResponse = []byte("[]")

func API() http.Handler {
	router := httprouter.New()

	router.DELETE("/bucket/:id", record)
	router.GET("/bucket/:id", record)
	router.PATCH("/bucket/:id", record)
	router.POST("/bucket/:id", record)

	router.DELETE("/report/:id", show)
	router.GET("/report/:id", show)

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

	mu.Lock()
	requests := storage[id]
	storage[id] = append(requests, Request{
		URL:           r.URL.Path,
		Method:        r.Method,
		Headers:       r.Header,
		Body:          base64.StdEncoding.EncodeToString(body.Bytes()),
		ContentLength: r.ContentLength,
	})
	records.Add(id)
	mu.Unlock()
}

func show(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	mu.Lock()
	requests := storage[id]

	if requests == nil {
		mu.Unlock()
		w.Write(emptyResponse)
		return
	}
	mu.Unlock()

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")

	if err := encoder.Encode(requests); err != nil {
		log.Printf("Failed to serialize payload with ID: %q, error: %s", id, err)
	}
}
