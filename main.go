package main

import (
    "bytes"
    "flag"
    "fmt"
    "path/filepath"
    "html/template"
    "image"
    "image/draw"
    _ "image/png"
    _ "image/gif"
    "image/jpeg"
    "launchpad.net/goamz/aws"
    "launchpad.net/goamz/s3"
    "log"
     "net/http"
)

var (
    port = flag.String("port", ":8080", "port")
    homeTemplate = template.Must(template.ParseFiles("templates/home.html"))
    auth = aws.Auth{ AccessKey: S3_ACCESS_KEY, SecretKey: S3_SECRET_KEY }
    s3Connection = s3.New(auth, aws.USEast)
    bucket = s3.Bucket{ s3Connection, S3_BUCKET_NAME }
    ACCEPTED_EXTENSIONS = []string{".jpg", ".jpeg", ".png", ".gif"}
)

func main() {
    handlers := map[string] func(http.ResponseWriter, *http.Request) () {
        "/"         : homeHandler,
        "/upload"   : photoUploadHandler,
    }
    for route, handler := range handlers {
        http.HandleFunc(route, handler)
    }

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
    log.Printf("+++ started photo upload")

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
    if checkFileExtension(filepath.Ext(filename)) == false {
        http.Error(response, "invalid file extension", 500)
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

func resizeAndUploadPhotos(filename string, photo *image.Image) {
    /* For each desired size for the photo, start a goroutine
    to scale and upload the image
    */
    for _, size := range []int{0, 50, 100} {
        go resizeAndUpload(filename, *photo, size)
    }
}

func getBounds(sourceImage *image.Image) (image.Rectangle, image.Point) {
    /* TODO this should be called once, regardless of the number of additional sizes
    */
    var (
        canvasSize = (*sourceImage).Bounds()
        height = canvasSize.Max.Y
        width = canvasSize.Max.X
        topLeftPoint = image.Point{0, 0}
    )
    log.Printf("bounds: %v; height, width: %d, %d", canvasSize, height, width)
    if height > width {
        canvasSize = image.Rect(0, 0, width, width)
        topLeftPoint = image.Point{0, (height-width)/2}
    } else if width > height {
        canvasSize = image.Rect(0, 0, height, height)
        topLeftPoint = image.Point{(width-height)/2, 0}
    }
    return canvasSize, topLeftPoint
}

func getSquare(sourceImage *image.Image) (image.Image) {
    squareSize, topLeftPoint := getBounds(sourceImage)
    squareImage := image.NewRGBA(squareSize)
    draw.Draw(squareImage, squareSize, *sourceImage, topLeftPoint, draw.Src)
    return squareImage
}

func resize(photo *image.Image, size int) (image.Image) {
    /* Get the desired rectangle we want to crop, given the size, 
    and then call the exported Resize function from resize.go
    */ 
    squareImage := getSquare(photo)
    return Resize(squareImage, squareImage.Bounds(), size, size)
}

func resizeAndUpload(filename string, photo image.Image, size int) {
    /* Step 1: create an image.Image that is a scaled version of :photo
    Step 2: convert :photo to a []byte
    Step 3: upload the rescaled file to s3
    */

    // resize the image
    var croppedImage image.Image
    switch size {
    case 100:
        croppedImage = resize(&photo, size)
    case 50:
        croppedImage = resize(&photo, size)
    case 0:
        croppedImage = photo
    default:
        log.Printf("invalid size: %s", size)
        return
    }

    // convert the image to a []byte
    var photoBytes bytes.Buffer
    options := jpeg.Options{Quality: 100}
    if err := jpeg.Encode(&photoBytes, croppedImage, &options); err != nil {
        log.Printf("couldn't jpeg encode: %s", err) 
    }
    newFilename := fmt.Sprintf("%s_%d", filename, size)

    // upload the image to s3
    uploadToS3(newFilename, photoBytes.Bytes())
}

func uploadToS3(filename string, photoBytes []byte) {
    /* PUT :photo in the s3 bucket, with the name :filename
    */
    log.Printf("uploading: %s", filename)
    bucket.Put(filename, photoBytes, "image/jpeg", s3.PublicRead)
    log.Printf("done uploading: %s", filename)
}

func checkFileExtension(extension string) (bool) {
    /* return true if the extenion passed in is one of the 
    accepted extensions, which are image file extensions.
    return false otherwise
    */
    for _, acceptedExtension := range ACCEPTED_EXTENSIONS {
        if extension == acceptedExtension {
            return true
        }
    }
    return false
}
