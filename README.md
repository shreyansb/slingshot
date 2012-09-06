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

You'll need a file `settings.go` that looks like this

```
package main

const (
	S3_BUCKET_NAME = "your_bucket_name"
    S3_ACCESS_KEY = "your_access_key"
    S3_SECRET_KEY = "your_secret_key"
)
```

`go run *.go` then `http://localhost:8080`