Crops and resized a picture to 100x100 and 50x50 squares, and uploads them and the original to S3.
You'll need a file `settings.go` that looks like this

```
package main

const (
	S3_BUCKET_NAME = "your_bucket_name"
    S3_ACCESS_KEY = "your_access_key"
    S3_SECRET_KEY = "your_secret_key"
)
```

To run the program:

```go run *.go```

and visit:

```http://localhost:8080```