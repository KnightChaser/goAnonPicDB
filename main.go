package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Image struct {
	gorm.Model
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

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Display the form to upload images
	tmpl, err := template.New("index.html").Parse(`
<!DOCTYPE html>
<html>
<head>
	<title>Image Uploader</title>
</head>
<body>
	<form action="/upload" method="post" enctype="multipart/form-data">
		<input type="file" name="image" />
		<input type="submit" value="Upload" />
	</form>
</body>
</html>
`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the form data to get the uploaded file
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Save the file to the server
	savePath := filepath.Join("static", handler.Filename)
	out, err := os.Create(savePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Save file information to the database
	img := Image{Filename: handler.Filename}
	db.Create(&img)

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
