package util

import (
	"net/http"
	"os"
)

func SetupHttpHandlerDummy() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "5000"
	}
	go func() {
		http.HandleFunc("/", HttpRedirectToApp)
		http.ListenAndServe(":"+port, nil)
	}()
}

// HttpRedirectToApp - Provides a HTTP redirect to the pganalyze app
func HttpRedirectToApp(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://app.pganalyze.com/", http.StatusFound)
}
