package main

import (
	"bytes"
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/storage"
	"fmt"
	"google.golang.org/api/option"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sketch-bridge/arduino-compile-server/common"
	"sketch-bridge/arduino-compile-server/database"
	"sketch-bridge/arduino-compile-server/web"
	"strings"
)

type BuildResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func main() {
	ctx := context.Background()

	app := createFirebaseApp(ctx)
	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer firestoreClient.Close()
	storageClient, err := app.Storage(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	port := "8080"
	h := func(w http.ResponseWriter, r *http.Request) {
		handleRequest(w, r, ctx, firestoreClient, storageClient, app)
	}
	http.HandleFunc("/build", h)
	log.Printf("Arduino Sketch Build Server is running on port %s.\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func createFirebaseApp(ctx context.Context) *firebase.App {
	sa := option.WithCredentialsFile("sketch-bridge-c8804059e16c.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	//app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatalln(err)
	}
	return app
}

// Handles the HTTP request.
func handleRequest(w http.ResponseWriter, r *http.Request, ctx context.Context, firestoreClient *firestore.Client, storageClient *storage.Client, app *firebase.App) {
	log.Printf("%s %s %s\n", r.Method, r.URL, r.Proto)

	// For CORS
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "authorization,content-type")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	uid, err := authenticateUser(r, app)
	if err != nil {
		sendFailureResponse(ctx, w, err, http.StatusUnauthorized)
		return
	}
	fmt.Printf("uid: %s\n", uid)

	params, err := web.ParseParameters(r)
	if err != nil {
		sendFailureResponse(ctx, w, err, http.StatusBadRequest)
		return
	}

	fmt.Printf("Project ID: %s\n", params.ProjectId)

	project, err := database.GetProject(ctx, firestoreClient, params.ProjectId)
	if err != nil {
		sendFailureResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}

	buildDirectoryPath, err := deleteProjectBuildDirectory(project.Id)
	if err != nil {
		sendFailureResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}
	_, err = deleteProjectSketchDirectory(project.Id)
	if err != nil {
		sendFailureResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}
	sketchDirectoryPath, err := createProjectSketch(project)
	if err != nil {
		sendFailureResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}

	cmd := exec.Command("/usr/local/bin/arduino-cli", "compile", "--output-dir", buildDirectoryPath, "--fqbn", "arduino:avr:uno", sketchDirectoryPath)
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
		sendSuccessfulResponse(ctx, w, false, stderrString)
		return
	}

	err = database.UploadHexFile(ctx, storageClient, project, uid)
	if err != nil {
		sendFailureResponse(ctx, w, err, http.StatusInternalServerError)
		return
	}

	sendSuccessfulResponse(ctx, w, true, stdoutString)
}

func sendSuccessfulResponse(ctx context.Context, w http.ResponseWriter, success bool, body string) {
	// For CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	response := BuildResponse{
		Success: success,
		Message: body,
	}
	json.NewEncoder(w).Encode(response)
}

func sendFailureResponse(ctx context.Context, w http.ResponseWriter, cause error, status int) {
	// For CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")

	log.Printf("[ERROR] %s\n", cause.Error())
	w.WriteHeader(status)
	io.WriteString(w, cause.Error())
}

func deleteProjectBuildDirectory(projectId string) (string, error) {
	dirPath := filepath.Join("/app/build", projectId)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return dirPath, nil
	}
	err := os.RemoveAll(dirPath)
	if err != nil {
		return "", err
	}
	return dirPath, nil
}

func deleteProjectSketchDirectory(projectId string) (string, error) {
	dirPath := filepath.Join("/app/sketches", projectId)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return dirPath, nil
	}
	err := os.RemoveAll(dirPath)
	if err != nil {
		return "", err
	}
	return dirPath, nil
}

func createProjectSketch(project *common.Project) (string, error) {
	dirPath := filepath.Join("/app/sketches", project.Id)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory: %v", err)
	}
	filePath := filepath.Join(dirPath, project.Id+".ino")
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()
	if _, err := file.WriteString(project.Code); err != nil {
		return "", fmt.Errorf("failed to write to file: %v", err)
	}
	return dirPath, nil
}

func authenticateUser(r *http.Request, app *firebase.App) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is missing")
	}

	idToken := strings.TrimPrefix(authHeader, "Bearer ")
	if idToken == authHeader {
		return "", fmt.Errorf("invalid authorization header format")
	}

	client, err := app.Auth(context.Background())
	if err != nil {
		return "", fmt.Errorf("error getting Auth client: %v", err)
	}

	token, err := client.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		return "", fmt.Errorf("error verifying ID token: %v", err)
	}

	return token.UID, nil
}
