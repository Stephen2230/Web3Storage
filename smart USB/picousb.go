package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var storageProviders = []string{
	"f01149094",
	"f0455466",
	"f0440182",
	"f01137729",
	"f0440208",
	"f0440208",
	"f01149094",
	"f01137729",
	"f0455466",
	"f0440182",
	"f0440208",
	"f01137729",
	"f0440182",
	"f0455466",
	"f01149094",
	"f0455466",
	"f01149094",
	"f01137729",
	"f0440208",
	"f0440182",
}

const usbCapacityGB = 32

func main() {
	// Call the MinIO integration function
	err := integrateWithMinIO()
	if err != nil {
		log.Fatal(err)
	}

}

func integrateWithMinIO() error {
	serviceAddress := "buckets.storage.io"
	accessKey, secKey := "", ""

	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("can't get PWD: %s", err)
	}

	myDocument1 := filepath.Join(pwd, "document1.pdf")
	myDocument2 := filepath.Join(pwd, "document2.pdf")

	// Initiate a client
	client, err := minio.New(serviceAddress, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secKey, ""),
		Region: "us-east-1",
		Secure: true,
	})
	if err != nil {
		return err
	}
	ctx := context.Background()

	bucketName := "test-bucket"
	// Create a bucket
	err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		return fmt.Errorf("error creating bucket \"%s\": %s", bucketName, err.Error())
	}

	// Check if the USB is 90% full
	if isUSB90PercentFull() {
		for _, provider := range storageProviders {
			err := requestStorageFromProvider(client, ctx, bucketName, provider)
			if err != nil {
				log.Printf("Error requesting storage from provider %s: %s", provider, err.Error())
			}
		}
	}

	return nil
}

func isUSB90PercentFull() bool {

	usedSpaceGB := 25 // Replace with the actual used space on your USB drive in GB

	// Calculate the threshold for 90% full
	threshold := float64(usbCapacityGB) * 0.9

	// Check if used space is greater than or equal to 90% threshold
	return float64(usedSpaceGB) >= threshold
}

func requestStorageFromProvider(client *minio.Client, ctx context.Context, bucketName, provider string) error {
	ctx := context.Background()

	// bucket name that we created
	bucketName := "test-bucket"

	// present working directory
	pwd, err := os.Getwd()
	if err != nil {
		log.Panicf("can't get PWD: %s", err)
	}

	// document paths
	myDocument1 := filepath.Join(pwd, "document1.pdf")

	// S3 path keys
	document1Key := "letter1/document.pdf"

	// upload objects to the bucket
	_, err = client.FPutObject(
		ctx, bucketName, document1Key, myDocument1, minio.PutObjectOptions{
			DisableMultipart: true,
		})
	if err != nil {
		log.Panicf("error putting object \"%s\" to bucket: %s", myDocument1, err.Error())
	}
	log.Printf("Requesting storage from provider %s", provider)
	return nil
}
