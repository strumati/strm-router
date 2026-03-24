package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

// Routing map
//var routes = map[string]string{
//	"s3": "http://localhost:3901",
//	"j2": "http://localhost:3902",
//}

var routes map[string]interface{}

// Extract AWS service from Authorization header
func extractService(auth string) string {

	parts := strings.Split(auth, "Credential=")
	if len(parts) < 2 {
		return ""
	}

	credPart := strings.Split(parts[1], ",")[0]
	segments := strings.Split(credPart, "/")

	if len(segments) < 4 {
		return ""
	}

	return segments[3]
}

// Create reverse proxy
func newProxy(target string) *httputil.ReverseProxy {
	u, err := url.Parse(target)
	if err != nil {
		log.Println("Invalid target:", target)
		return nil
	}

	proxy := httputil.NewSingleHostReverseProxy(u)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)

		// Preserve original host if needed
		req.Header.Set("X-Forwarded-Host", req.Host)
		req.Header.Set("X-Real-IP", req.RemoteAddr)
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Println("Proxy error:", err)
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
	}

	return proxy
}

// Main handler
func handler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")

	service := extractService(auth)
	log.Println("Detected service:", service)

	target, ok := routes[service]
	if !ok {
		http.Error(w, "Unknown service", http.StatusBadRequest)
		return
	}

	proxy := newProxy(target.(string))
	if proxy == nil {
		http.Error(w, "Proxy setup failed", http.StatusInternalServerError)
		return
	}

	log.Printf("Routing %s → %s\n", service, target)
	proxy.ServeHTTP(w, r)
}

func main() {
	raw, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(raw, &routes); err != nil {
		panic(err)
	}

	http.HandleFunc("/", handler)

	log.Println("Router running on :3900")
	log.Fatal(http.ListenAndServe(":3900", nil))
}
