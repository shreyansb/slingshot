Receives a photo via HTTP POST, resizes and crops the photo into 50x50 and 100x100 squares, and uploads all three to S3.

```
[shreyans ~/Code/slingshot (master)]$ go run *.go
2012/09/06 18:50:13 [main] starting server on localhost:8080
2012/09/06 18:50:13 [receivePhotos] locking this goroutine to an OS thread
2012/09/06 18:50:23 [photoUploadHandler] receiving photo: burntedges
2012/09/06 18:50:23 [photoUploadHandler] returning HTTP response to client
2012/09/06 18:50:23 [sendPhotoDetails] sending photoDetails over chan
2012/09/06 18:50:23 [receivePhotos] received photo over chan: burntedges
2012/09/06 18:50:23 [uploadToS3] starting upload of burntedges/50
2012/09/06 18:50:23 [uploadToS3] starting upload of burntedges/100
2012/09/06 18:50:23 [uploadToS3] starting upload of burntedges/full
2012/09/06 18:50:24 [uploadToS3] done uploading burntedges/50
2012/09/06 18:50:24 [uploadToS3] done uploading burntedges/100
2012/09/06 18:50:27 [uploadToS3] done uploading burntedges/full
```

You'll need to define the environment variables: `S3_BUCKET_NAME`, `S3_ACCESS_KEY`, and `S3_SECRET_KEY`

Run it with `go run *.go` or `go build -o slingshot && ./slingshot` 