// +build ignore

package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"Geekbux.com/api"
)

// #####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####
//	Main Function to start the Server!
// #####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####

func main() {
	httpPort := 8080
	server := api.NewServer()
	log.Printf("Server running on Port: %v\n", httpPort)
	err := http.ListenAndServe(fmt.Sprintf(":%d", httpPort), logRequest(server))
	if err != nil {
		log.Fatal(err)
	}
}

// #####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####
//	Usefull Function to Log Requests
// #####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####-#####

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", requestGetRemoteAddress(r), r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func ipAddrFromRemoteAddr(s string) string {
	idx := strings.LastIndex(s, ":")
	if idx == -1 {
		return s
	}
	return s[:idx]
}

func requestGetRemoteAddress(r *http.Request) string {
	hdr := r.Header
	hdrRealIP := hdr.Get("X-Real-Ip")
	hdrForwardedFor := hdr.Get("X-Forwarded-For")
	if hdrRealIP == "" && hdrForwardedFor == "" {
		return ipAddrFromRemoteAddr(r.RemoteAddr)
	}
	if hdrForwardedFor != "" {
		parts := strings.Split(hdrForwardedFor, ",")
		for i, p := range parts {
			parts[i] = strings.TrimSpace(p)
		}

		return parts[0]
	}
	return hdrRealIP
}
