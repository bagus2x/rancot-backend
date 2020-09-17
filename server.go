package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	h "github.com/bagus2x/rancot-backend/helpers"
)

// CORSMiddleware -
func CORSMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Origin, Authorization, Accept, Client-Security-Token, Accept-Encoding, x-access-token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" {
			fmt.Println("OPTIONS")
			http.Error(w, "", 200)
		} else {
			h.ServeHTTP(w, r)
		}
	})

}

func main() {
	gm := new(h.GeMux)
	gm.Use(CORSMiddleware)
	gm.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, mamank"))
	})

	gm.HandleFunc("/api/ws", h.WS)
	port := os.Getenv("PORT")
	svr := http.Server{
		Addr:    ":" + port,
		Handler: gm,
	}
	log.Println("server running on port", port)
	log.Fatal(svr.ListenAndServe())
}
