package custom

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"cloud.google.com/go/storage"
)

// uploadFile uploads an object.
func UploadFile() error {
	bucket := "gotrade.esol.top"
	// object := "object-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("storage.NewClient: %v", err)
	}
	defer client.Close()

	uploadDir(client, ctx, "swd", bucket)
	if time.Now().Weekday() == time.Friday {
		uploadDir(client, ctx, "week", bucket)
	}
	uploadDir(client, ctx, "", bucket)
	return nil
}

func uploadDir(client *storage.Client, ctx context.Context, dir string, bucket string) {
	var directory = ""
	if dir != "" {
		directory = fmt.Sprintf("%s/", dir)
	}
	// Files
	entriesSwd, err := os.ReadDir(fmt.Sprintf("/tmp/report/%s", directory))
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entriesSwd {
		if e.IsDir() {
			continue
		}
		log.Println(e.Name())
		// Open local file.
		f, err := os.Open(fmt.Sprintf("/tmp/report/%s", directory) + e.Name())
		if err != nil {
			log.Fatalf("os.Open: %v", err)
		}
		defer f.Close()

		ctx, cancel := context.WithTimeout(ctx, time.Second*50)
		defer cancel()

		o := client.Bucket(bucket).Object(directory + e.Name())

		// o = o.If(storage.Conditions{DoesNotExist: true})
		// o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

		// Upload an object with storage.Writer.
		wc := o.NewWriter(ctx)
		if _, err = io.Copy(wc, f); err != nil {
			log.Fatalf("io.Copy: %v", err)
		}
		if err := wc.Close(); err != nil {
			log.Fatalf("Writer.Close: %v", err)
		}
	}
}
