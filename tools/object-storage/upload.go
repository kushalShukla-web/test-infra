package main

import (
	"context"
	"fmt"

	// "io"
	"os"
	"path/filepath"

	"github.com/go-kit/log"
	"github.com/thanos-io/objstore"
	"github.com/thanos-io/objstore/client"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Store struct {
	TsdbPath     string
	ObjectKey    string
	ObjectConfig string
}

func newstore() *Store {
	return &Store{
		TsdbPath:     "",
		ObjectKey:    "",
		ObjectConfig: "",
	}
}

func uploadTSDBFiles(bucket objstore.Bucket, tsdbPath string, commitDir string) error {
	err := filepath.Walk(tsdbPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		objectKey := filepath.Join(commitDir, filepath.Base(filePath))

		fmt.Println(objectKey, "Hi from objectKey")

		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %v", filePath, err)
		}
		defer file.Close()

		if err := bucket.Upload(context.Background(), objectKey, file); err != nil {
			return fmt.Errorf("failed to upload file %s: %v", objectKey, err)
		}

		fmt.Printf("Uploaded file: %s\n", objectKey)
		return nil
	})

	return err
}

func (c *Store) Upload(*kingpin.ParseContext) error {
	config := c.ObjectConfig
	configBytes, err := os.ReadFile(config)
	if err != nil {
		panic(fmt.Errorf("failed to read config file: %v", err))
	}
	commitDir := c.ObjectKey

	bucket, err := client.NewBucket(log.NewNopLogger(), configBytes, "example")
	if err != nil {
		panic(err)
	}
	exists, err := bucket.Exists(context.Background(), "example")
	if err != nil {
		panic(err)
	}
	fmt.Println(exists)

	err = uploadTSDBFiles(bucket, c.TsdbPath, commitDir)

	if err != nil {
		fmt.Printf("Failed to upload TSDB files: %v\n", err)
	} else {
		fmt.Println("TSDB files successfully uploaded to object store.")
	}

	return err
}
