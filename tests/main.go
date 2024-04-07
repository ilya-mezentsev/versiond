// Package main is for running test server that increments fake version each second request.

package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

var (
	requestsCounter int64 = 0
	versionNumber   int64 = 1
)

func handler(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt64(&requestsCounter, 1)

	if requestsCounter%2 == 0 {
		atomic.AddInt64(&versionNumber, 1)
	}

	_, _ = fmt.Fprintf(w, "{\"newest_version\": \"v%d\"}", atomic.LoadInt64(&versionNumber))
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8888", nil))
}
