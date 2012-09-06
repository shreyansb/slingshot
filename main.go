package main

import (
	"flag"
	"fmt"
	"html/template"
	"image"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"runtime"
)

var (
	port                 = flag.String("port", ":8080", "port")
	homeTemplate         = template.Must(template.ParseFiles("home.html"))
	acceptedContentTypes = []string{"image/jpeg", "image/png", "image/gif"}
	resizerChan          chan PhotoDetails
)

type PhotoDetails struct {
	photo    *image.Image
	filename string
}

func main() {
	// use all the CPU cores available
	runtime.GOMAXPROCS(runtime.NumCPU())

	// start a goroutine to handle photo resizing on a separate core and
	// initialize a chan to send data to the goroutine
	resizerChan = make(chan PhotoDetails)
	go receivePhotos()

	// set up handlers
	handlers := map[string]func(http.ResponseWriter, *http.Request){
		"/":       homeHandler,
		"/upload": photoUploadHandler,
	}
	for route, handler := range handlers {
		http.HandleFunc(route, handler)
	}

	// start server
	flag.Parse()
	log.Printf("[main] starting server on localhost%s", *port)
	if err := http.ListenAndServe(*port, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func homeHandler(response http.ResponseWriter, request *http.Request) {
	/* serves the HTML upload form */
	homeTemplate.Execute(response, request.Host)
}

func photoUploadHandler(response http.ResponseWriter, request *http.Request) {
	/* extracts the uploaded file from the HTTP request, and after performing some 
	   basic checks, launches a goroutine to resize and upload the files.
	*/
	// if the request method isn't a POST, redirect to the homepage
	if request.Method != "POST" {
		homeHandler(response, request)
		return
	}

	// extract the file from the request
	file, fileHeader, err := request.FormFile("photo")
	if err != nil {
		http.Error(response, fmt.Sprintf("%v", err), 500)
		return
	}

	// close the file after this handler has run
	defer file.Close()

	// check file types, make sure we received an image
	filename := getFilename(request, fileHeader)
	contentType := (fileHeader.Header).Get("Content-Type")
	if checkContentType(contentType) == false {
		http.Error(response, "invalid content type", 500)
	}
	log.Printf("[photoUploadHandler] receiving photo: %s", filename)

	photo, _, err := image.Decode(file)
	if err != nil {
		http.Error(response, fmt.Sprintf("%v", err), 500)
		return
	}

	// send the details of this photo over a channel to the 
	// resizer goroutine, which is locked to a core
	photoDetails := PhotoDetails{&photo, filename}
	go sendPhotoDetails(photoDetails)

	log.Printf("[photoUploadHandler] returning HTTP response to client")
}

func sendPhotoDetails(photoDetails PhotoDetails) {
	log.Printf("[sendPhotoDetails] sending photoDetails over chan")
	resizerChan <- photoDetails
}

func checkContentType(contentType string) bool {
	/* return true if the extenion passed in is one of the 
	   accepted extensions, which are image file extensions.
	   return false otherwise
	*/
	for _, acceptedContentType := range acceptedContentTypes {
		if contentType == acceptedContentType {
			return true
		}
	}
	return false
}

func getFilename(request *http.Request, fileHeader *multipart.FileHeader) string {
	var filename string
	filename = request.FormValue("filename")
	if filename == "" {
		baseFilename := fileHeader.Filename
		extension := filepath.Ext(baseFilename)
		filename = baseFilename[:len(baseFilename)-len(extension)]
	}
	return filename
}
