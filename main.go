package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

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
	// Register subdirectories for serving files inside those directories
	http.Handle("/asset/", http.StripPrefix("/asset", http.FileServer(http.Dir("asset"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/frontend/", http.StripPrefix("/frontend/", http.FileServer(http.Dir("frontend"))))

	// Specify GET/POST action occuring position
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", HomeHandler).Methods("GET")
	muxRouter.HandleFunc("/upload", UploadHandler).Methods("POST")

	// Default
	http.Handle("/", muxRouter)

	// Specifying ports
	listeningPort := ":8080"
	fmt.Printf("Server is running on %s\n", listeningPort)
	http.ListenAndServe(listeningPort, nil)

}

func HomeHandler(responseWriter http.ResponseWriter, request *http.Request) {
	// Count the total number of images in the database
	var imageQty int64
	db.Model(&Image{}).Count(&imageQty)

	// Grab the recent 6 images with both Filename and CreatedAt fields
	var images []Image
	db.Select("id, uploader, filename, created_at").Order("created_at desc").Limit(6).Find(&images)

	// Read HTML content from file
	htmlContent, err := os.ReadFile("frontend/index.html")
	if err != nil {
		http.Error(responseWriter, "Unable to read HTML file", http.StatusInternalServerError)
		return
	}

	// Create a template using html/template package for more dynamic content
	tmpl, err := template.New("index").Parse(string(htmlContent))
	if err != nil {
		http.Error(responseWriter, "Unable to parse HTML template", http.StatusInternalServerError)
		return
	}

	// Send the HTML content with both the image count and the recent images as the response
	responseWriter.Header().Set("Content-Type", "text/html")
	responseWriter.WriteHeader(http.StatusOK)

	// Create a struct to hold both the image count and the images
	data := struct {
		ImageQty int64
		Images   []Image
	}{
		ImageQty: imageQty,
		Images:   images,
	}

	tmpl.Execute(responseWriter, data)
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

	// Save the file to the server. The image will be saved at dir "/${imageSaveDirectoryName}"
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

	// Create picture saving path
	savePath := filepath.Join("static", handler.Filename)
	out, err := os.Create(savePath)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	defer out.Close()

	// Save image to the created path
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
