package custom

import (
	"context"
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

	// Files
	entries, err := os.ReadDir("/tmp/report")
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		log.Println(e.Name())
		// Open local file.
		f, err := os.Open("/tmp/report/" + e.Name())
		if err != nil {
			log.Fatalf("os.Open: %v", err)
		}
		defer f.Close()

		ctx, cancel := context.WithTimeout(ctx, time.Second*50)
		defer cancel()

		o := client.Bucket(bucket).Object(e.Name())

		// Optional: set a generation-match precondition to avoid potential race
		// conditions and data corruptions. The request to upload is aborted if the
		// object's generation number does not match your precondition.
		// For an object that does not yet exist, set the DoesNotExist precondition.
		// o = o.If(storage.Conditions{DoesNotExist: true})
		// If the live object already exists in your bucket, set instead a
		// generation-match precondition using the live object's generation number.
		// attrs, err := o.Attrs(ctx)
		// if err != nil {
		// 	log.Fatalf("object.Attrs: %v", err)
		// }
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
	return nil
}
