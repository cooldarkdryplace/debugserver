package debugserver

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

var storage = NewStorage()
var emptyResponse = []byte("[]\n")
var errNoID = errors.New("no ID provided")

func API(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/buckets/") {
		record(w, r)
	}

	if strings.HasPrefix(r.URL.Path, "/reports/") {
		if r.Method == http.MethodDelete {
			del(w, r)
		}
		if r.Method == http.MethodGet {
			show(w, r)
		}
	}
}

type Request struct {
	ContentLength int64               `json:"content_length"`
	Body          string              `json:"body"`
	Method        string              `json:"method"`
	URL           string              `json:"url"`
	QueryParams   string              `json:"query_params"`
	RemoteIP      string              `json:"remote_ip"`
	Headers       map[string][]string `json:"headers"`
}

func requestID(path string) (string, error) {
	path = strings.TrimPrefix(path, "/")
	chunks := strings.Split(path, "/")
	if len(chunks) < 2 {
		return "", errNoID
	}
	id := strings.TrimSuffix(chunks[1], "/")
	id = strings.TrimSpace(id)

	if id == "" {
		return "", errNoID
	}

	return id, nil
}

func record(w http.ResponseWriter, r *http.Request) {
	id, err := requestID(r.URL.Path)
	if err != nil {
		http.Error(w, "no ID provided", http.StatusBadRequest)
		return
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Printf("Failed to sptil remote address to IP and Port: %s", err)
	}

	bodyReader := http.MaxBytesReader(w, r.Body, 1<<20)
	defer bodyReader.Close()
	var body bytes.Buffer

	if _, err := io.Copy(&body, bodyReader); err != nil {
		http.Error(w, "unable to read body", http.StatusInternalServerError)
		return
	}

	storage.Add(id, Request{
		RemoteIP:      ip,
		URL:           r.URL.Path,
		QueryParams:   r.URL.RawQuery,
		Method:        r.Method,
		Headers:       r.Header,
		Body:          base64.StdEncoding.EncodeToString(body.Bytes()),
		ContentLength: r.ContentLength,
	})
}

func show(w http.ResponseWriter, r *http.Request) {
	id, err := requestID(r.URL.Path)
	if err != nil {
		http.Error(w, "no ID provided", http.StatusBadRequest)
		return
	}

	requests := storage.Get(id)
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

func del(w http.ResponseWriter, r *http.Request) {
	id, err := requestID(r.URL.Path)
	if err != nil {
		http.Error(w, "no ID provided", http.StatusBadRequest)
		return
	}

	storage.Del(id)
}
