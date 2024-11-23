package services

import (
	"io"
	"net/http"
	"time"
)

func ProxyRequest(w http.ResponseWriter, r *http.Request) {
	targetService := r.URL.Query().Get("target")
	if targetService == "" {
		http.Error(w, "Target service not specified", http.StatusBadRequest)
		return
	}

	req, err := http.NewRequest(r.Method, targetService+r.URL.Path, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer req.Body.Close()
	req.Header = r.Header
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
