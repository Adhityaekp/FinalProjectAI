package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"a21hc3NpZ25tZW50/service"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

var fileService = &service.FileService{}
var aiService = &service.AIService{Client: &http.Client{}}
var store = sessions.NewCookieStore([]byte("my-key"))

func getSession(r *http.Request) *sessions.Session {
	session, _ := store.Get(r, "chat-session")
	return session
}

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Retrieve the Hugging Face token from the environment variables
	token := os.Getenv("HUGGINGFACE_TOKEN")
	if token == "" {
		log.Fatal("HUGGINGFACE_TOKEN is not set in the .env file")
	}

	// Set up the router
	router := mux.NewRouter()

	router.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		token := os.Getenv("HUGGINGFACE_TOKEN")
		if token == "" {
			http.Error(w, "Token not provided", http.StatusInternalServerError)
			return
		}

		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Failed to read file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		fileContent, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, "Failed to read file content", http.StatusInternalServerError)
			return
		}

		table, err := fileService.ProcessFile(string(fileContent))
		if err != nil {
			http.Error(w, "Failed to process file", http.StatusInternalServerError)
			return
		}

		query := "Find the least and most electricity usage" // Contoh query
		response, err := aiService.AnalyzeData(table, query, token)
		if err != nil {
			http.Error(w, "Failed to analyze data", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"status": "success",
			"answer": response,
		})
	}).Methods("POST")

	router.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		token := os.Getenv("HUGGINGFACE_TOKEN")
		if token == "" {
			http.Error(w, "Token not provided", http.StatusInternalServerError)
			return
		}

		var request struct {
			Context string `json:"context"`
			Query   string `json:"query"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		response, err := aiService.ChatWithAI(request.Context, request.Query, token)
		if err != nil {
			http.Error(w, "Failed to communicate with AI", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"status": "success",
			"answer": response.GeneratedText,
		})
	}).Methods("POST")

	// Enable CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"}, // Allow your React app's origin
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}).Handler(router)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler))
}
