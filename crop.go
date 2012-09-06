package main

import (
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"runtime"
)

var (
	uploadSizes = []int{0, 50, 100}
)

func receivePhotos() {
	/* This function handles the CPU intensize work of cropping
	   and resizing uploaded photos.

	   Since CPU bound work blocks the go scheduler, this function is 
	   run as a goroutine and locked to an OS thread, and thus core(?), 
	   restricting all CPU heavy work to one core, leaving the remaining cores
	   free to handle incoming requests.

	   This function waits for incoming photo data on a channel, resizes and crops the photo,
	   and kicks off some goroutines to upload the photos to S3
	*/
	log.Printf("[receivePhotos] locking this goroutine to an OS thread")
	runtime.LockOSThread()
	for {
		select {
		case photoDetails := <-resizerChan:
			log.Printf("[receivePhotos] received photo over chan: %s", photoDetails.filename)
			resizeAndUploadPhotos(photoDetails.filename, photoDetails.photo)
		}
	}
}

func resizeAndUploadPhotos(filename string, photo *image.Image) {
	/* crop and resize the incoming :photo into a few sizes and 
	   kick off the upload
	*/

	// get one square crop, to be used for the various resizings
	squareImage := getSquareCrop(photo)

	for _, size := range uploadSizes {
		var imageToUpload image.Image
		switch size {
		case 0:
			imageToUpload = *photo
		default:
			imageToUpload = resize(squareImage, size)
		}
		go uploadPhoto(filename, imageToUpload, size)
	}
}

func resize(photoToResize image.Image, size int) image.Image {
	/* a helper to call Resize in resize.go */
	return Resize(photoToResize, photoToResize.Bounds(), size, size)
}

func getSquareCrop(sourceImage *image.Image) image.Image {
	/* create and return a new image.Image that holds the middle square of 
	   the :sourceImage
	*/
	squareSize, topLeftPoint := getSquareBounds(sourceImage)
	squareImage := image.NewRGBA(squareSize)
	draw.Draw(squareImage, squareSize, *sourceImage, topLeftPoint, draw.Src)
	return squareImage
}
func getSquareBounds(sourceImage *image.Image) (image.Rectangle, image.Point) {
	/* return the dimensions and top-left-point required to create a square crop
	   of the :sourceImage
	*/
	var (
		canvasSize   = (*sourceImage).Bounds()
		height       = canvasSize.Max.Y
		width        = canvasSize.Max.X
		topLeftPoint = image.Point{0, 0}
	)
	if height > width {
		canvasSize = image.Rect(0, 0, width, width)
		topLeftPoint = image.Point{0, (height - width) / 2}
	} else if width > height {
		canvasSize = image.Rect(0, 0, height, height)
		topLeftPoint = image.Point{(width - height) / 2, 0}
	}
	return canvasSize, topLeftPoint
}