package main

import (
	"flag"
	"fmt"
	"html/template"
	"image"
	"log"
	"net/http"
)

var (
	port               = flag.String("port", ":8080", "port")
	homeTemplate       = template.Must(template.ParseFiles("templates/home.html"))
	acceptedContentTypes = []string{"image/jpeg", "image/png", "image/gif"}
)

func main() {
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
	filename := fileHeader.Filename
    contentType := (fileHeader.Header).Get("Content-Type")
	if checkContentType(contentType) == false {
		http.Error(response, "invalid content type", 500)
	}

	photo, _, err := image.Decode(file)
	if err != nil {
		http.Error(response, fmt.Sprintf("%v", err), 500)
		return
	}

	// asynchronously resize and upload the photo to S3, returning
	// a response to the http client
	go resizeAndUploadPhotos(filename, &photo)
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
