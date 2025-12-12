// main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	startTime = time.Now()
)

type Info struct {
	Service   string            `json:"service"`
	Version   string            `json:"version"`
	StartTime string            `json:"start_time"`
	Env       map[string]string `json:"env"`
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/healthz", healthHandler)
	http.HandleFunc("/readyz", readyHandler)
	http.HandleFunc("/info", infoHandler)
	http.Handle("/metrics", promhttp.Handler())

	port := getEnv("PORT", "8080")
	log.Printf("starting traffic-simulate on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "ok")
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	// simple readiness: always ready if running
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "ready")
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	info := Info{
		Service:   "traffic-simulate",
		Version:   getEnv("APP_VERSION", "dev"),
		StartTime: startTime.Format(time.RFC3339),
		Env: map[string]string{
			"ENVIRONMENT": getEnv("ENVIRONMENT", "local"),
			"FEATURE_X":   getEnv("FEATURE_X", "off"),
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
