package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
)


func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	req, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Request got an Error: %s", req)
		return
	} else {
		log.Printf("Request at this point: %s", req)
	}

	switch r.Method {
	case http.MethodGet:
		handleGet(w, r, a)
	case http.MethodPut: 
		handlePut(w, r, a)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return	
	}	
}


func handleGet(w http.ResponseWriter, r *http.Request, a *App) {
	// parse the url to get the key and the ttl in seconds 
	key := strings.Split(r.URL.Path, "/")[1]
	log.Printf("Key: %s", key)

	itemVal , err := a.handleGet(key)
	if err != nil {
		log.Printf("Response: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(itemVal)
}

func handlePut(w http.ResponseWriter, r *http.Request, a *App) {
	key := strings.Split(r.URL.Path, "/")[1]
	ttlInSeconds , err := strconv.Atoi(r.URL.Query().Get("ttl"))
	if err != nil {
		log.Printf("TTL is not an integer: %s", r.URL.Query().Get("ttl"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var buff []byte
	buff = make([]byte, r.ContentLength)
	r.Body.Read(buff)

	err = a.handlePut(key, &ttlInSeconds, buff)
	if err != nil {
		log.Printf("Response: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	log.Printf("Key: %s", key)
	log.Printf("TTL: %d", ttlInSeconds)
	log.Printf("Body: %s", buff)
}

