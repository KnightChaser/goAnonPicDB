package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Image struct
type Image struct {
	gorm.Model
	Uploader string
	Filename string
}

var db *gorm.DB

func init() {
	// Open a database connection
	var err error
	db, err = gorm.Open(sqlite.Open("imageUploader.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database")
	}

	// AutoMigrate creates the table if it doesn't exist
	db.AutoMigrate(&Image{})
}

func main() {
	muxRouter := mux.NewRouter()

	muxRouter.HandleFunc("/", HomeHandler).Methods("GET")
	muxRouter.HandleFunc("/upload", UploadHandler).Methods("POST")

	// Serve static files (images)
	muxRouter.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	http.Handle("/", muxRouter)

	listeningPort := ":8080"
	fmt.Printf("Server is running on %s", listeningPort)
	http.ListenAndServe(listeningPort, nil)
}

func HomeHandler(responseWriter http.ResponseWriter, request *http.Request) {
	// Read HTML content from file
	htmlContent, err := os.ReadFile("upload_form.html")
	if err != nil {
		http.Error(responseWriter, "Unable to read HTML file", http.StatusInternalServerError)
		return
	}

	// Send the HTML content as the response
	responseWriter.Header().Set("Content-Type", "text/html")
	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Write(htmlContent)
}

func UploadHandler(responseWriter http.ResponseWriter, request *http.Request) {
	// Parse the form data to get the uploaded file
	// Up to 25 MiB limit
	err := request.ParseMultipartForm(25 << 20)
	if err != nil {
		log.Fatal("Error during handling parsing POST request in UploadHandler")
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	file, handler, err := request.FormFile("image")
	uploader := request.FormValue("uploader")
	if err != nil {
		log.Fatal("Error during fetching data from \"image\" field in the request body.")
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Save the file to the server
	imageSaveDirectoryName := "static"
	_, err = os.Stat(imageSaveDirectoryName)
	if os.IsNotExist(err) {
		// Create the directory if it doesn't exist
		err := os.Mkdir(imageSaveDirectoryName, 0755)
		if err != nil {
			panic(fmt.Sprintf("Error creating directory: %v\n", err))
		}
		fmt.Printf("Directory \"%s\" created because it didn't exist.\n", imageSaveDirectoryName)
	} else if err != nil {
		// Handle other errors that may occur during stat
		panic(fmt.Sprintf("Error checking directory: %v\n", err))
	}

	savePath := filepath.Join("static", handler.Filename)
	out, err := os.Create(savePath)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	// Save file information to the database
	img := Image{Filename: handler.Filename, Uploader: uploader}
	db.Create(&img)

	// Redirect to home page
	http.Redirect(responseWriter, request, "/", http.StatusSeeOther)
}
