package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sketch-bridge/arduino-compile-server/web"
)

func main() {
	ctx := context.Background()

	port := "8080"
	h := func(w http.ResponseWriter, r *http.Request) {
		handleRequest(w, r, ctx)
	}
	http.HandleFunc("/build", h)
	log.Printf("Arduino Sketch Build Server is running on port %s.\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

// Handles the HTTP request.
func handleRequest(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	log.Printf("%s %s %s\n", r.Method, r.URL, r.Proto)

	params, err := web.ParseQueryParameters(r)
	if err != nil {
		log.Printf("[ERROR] %s\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())
		return
	}

	fmt.Printf("Project ID: %s\n", params.ProjectId)

	// For CORS
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "authorization,content-type")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	cmd := exec.Command("/usr/local/bin/arduino-cli", "compile", "--output-dir", "/app/build", "--fqbn", "arduino:avr:uno", "/app/sketches/blink")
	cmd.Dir = "/app"
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	stdoutString := stdout.String()
	if err != nil {
		log.Println("Building failed")
		stderrString := stderr.String()
		log.Printf("[ERROR] %s\n", err.Error())
		log.Printf("[ERROR] %s\n", stderrString)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "Building failed\n")
		io.WriteString(w, stderrString)
		return
	}

	// For CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Building succeeded\n")
	io.WriteString(w, stdoutString)

	files, err := os.ReadDir("/app/build")
	if err != nil {
		log.Printf("[ERROR] %s\n", err.Error())
		return
	}
	for _, file := range files {
		fmt.Println(file.Name())
	}
}
