// Copyright (c) 2022 Red Hat, Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"service-provider-integration-scm-file-retriever-server/websocket"
	"time"

	"github.com/gorilla/mux"
	"github.com/redhat-appstudio/service-provider-integration-scm-file-retriever/gitfile"
)

func OkHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func GetFileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoUrl := vars["repoUrl"]
	if repoUrl == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid repoUrl")
		return
	}
	filepath := vars["filepath"]
	if filepath == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid filepath")
		return
	}

	ref := vars["ref"]
	if ref == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid ref")
		return
	}
	namespace := vars["namespace"]
	if namespace == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid namespace")
		return
	}
	pageId := r.Header.Get("X-WebSocket-pageId")
	if pageId == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid header \"X-WebSocket-pageId\"")
		return
	}

	if !pool.IsClientKnown(pageId) {
		respondWithError(w, http.StatusBadRequest, "No registered websockets connected. Please reload the page.")
		return
	}

	ctx := context.TODO()
	content, err := gitfile.Default().GetFileContents(ctx, namespace, repoUrl, filepath, ref, func(ctx context.Context, url string) {
		message := websocket.Message{Type: 777, Body: url, ClientID: pageId}
		pool.SendMessage <- message
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_, err = io.Copy(w, content)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	message := websocket.Message{Type: 999, Body: "close if open", ClientID: pageId}
	pool.SendMessage <- message
}
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing middleware", r.Method)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers:", "Origin, Content-Type, X-Auth-Token, Authorization")
		next.ServeHTTP(w, r)
		log.Println("Executing middleware again")
	})
}
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func SendIndexHtml(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}
func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	log.Println("WebSocket Endpoint Hit")
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}
	vars := mux.Vars(r)
	log.Println("WebSocket vars[pageId]" + vars["pageId"])
	client := &websocket.Client{
		Conn: conn,
		Pool: pool,
		ID:   vars["pageId"],
	}

	pool.Register <- client
	client.Read()
}

// initialize with default fetcher
var pool = websocket.NewPool()

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	go pool.Start()

	router := mux.NewRouter()
	router.HandleFunc("/health", OkHandler).Methods("GET")
	router.HandleFunc("/ready", OkHandler).Methods("GET")
	router.HandleFunc("/gitfile", GetFileHandler).Queries("repoUrl", "{repoUrl}").Queries("filepath", "{filepath}").Queries("ref", "{ref}").Queries("namespace", "{namespace}").Methods("GET")
	router.HandleFunc("/ws/{pageId}", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	srv := &http.Server{
		Addr: "0.0.0.0:8000",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		//Handler:      corsMiddleware(router), // UI testing
		Handler: loggingMiddleware(router), // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)

}
