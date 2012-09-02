package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
	"log"
)

var (
	auth         = aws.Auth{AccessKey: S3_ACCESS_KEY, SecretKey: S3_SECRET_KEY}
	s3Connection = s3.New(auth, aws.USEast)
	bucket       = s3.Bucket{s3Connection, S3_BUCKET_NAME}
)

func uploadPhoto(filename string, photo image.Image, size int) {
	/* convert the incoming image.Image to a jpeg encoded byte array,
	   and upload it to S3 
	*/

	// convert the image to a []byte
	var photoBytes bytes.Buffer
	options := jpeg.Options{Quality: 100}
	if err := jpeg.Encode(&photoBytes, photo, &options); err != nil {
		log.Printf("couldn't jpeg encode: %s", err)
	}
	newFilename := fmt.Sprintf("%s_%d", filename, size)

	// upload the image to s3
	uploadToS3(newFilename, photoBytes.Bytes())
}

func uploadToS3(filename string, photoBytes []byte) {
	/* upload the given byte array :photoButes, to S3 */
	bucket.Put(filename, photoBytes, "image/jpeg", s3.PublicRead)
}
