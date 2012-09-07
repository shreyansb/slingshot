package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
	"log"
	"os"
)

var (
	auth         aws.Auth
	s3Connection *s3.S3
	bucket       s3.Bucket
)

func setupS3Connection() {
	if *s3_bucket_name == "" || *s3_access_key == "" || *s3_secret_key == "" {
		log.Printf("[init] missing S3 params")
		os.Exit(1)
	}
	auth = aws.Auth{*s3_access_key, *s3_secret_key}
	s3Connection = s3.New(auth, aws.USEast)
	bucket = s3.Bucket{s3Connection, *s3_bucket_name}
}

// convert the incoming image.Image to a jpeg encoded byte array,
// and upload it to S3 
func uploadPhoto(filename string, photo image.Image, size int) {
	// convert the image to a []byte
	var photoBytes bytes.Buffer
	options := jpeg.Options{Quality: 100}
	if err := jpeg.Encode(&photoBytes, photo, &options); err != nil {
		log.Printf("[uploadPhoto] couldn't jpeg encode: %s", err)
	}

	var newFilename string
	if size == 0 {
		newFilename = fmt.Sprintf("%s/full", filename)
	} else {
		newFilename = fmt.Sprintf("%s/%d", filename, size)
	}

	// upload the image to s3
	uploadToS3(newFilename, photoBytes.Bytes())
}

// upload the given byte array :photoButes, to S3
func uploadToS3(filename string, photoBytes []byte) {
	log.Printf("[uploadToS3] starting upload of %s", filename)
	err := bucket.Put(filename, photoBytes, "image/jpeg", s3.PublicRead)
	if err != nil {
		log.Printf("[uploadToS3] error uploading: %s", err)
	}
	log.Printf("[uploadToS3] done uploading %s", filename)
}
