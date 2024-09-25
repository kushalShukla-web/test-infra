// Copyright 2024 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
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

func uploadTSDBFiles(bucket objstore.Bucket, tsdbPath string) error {
	err := filepath.Walk(tsdbPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		objectKey := filepath.Base(filePath)

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

func (c *Store) fun(*kingpin.ParseContext) error {
	config := c.ObjectConfig
	configBytes, err := os.ReadFile(config)
	if err != nil {
		panic(fmt.Errorf("failed to read config file: %v", err))
	}

	bucket, err := client.NewBucket(log.NewNopLogger(), configBytes, "example")
	if err != nil {
		panic(err)
	}
	exists, err := bucket.Exists(context.Background(), "example")
	if err != nil {
		panic(err)
	}
	fmt.Println(exists)

	err = uploadTSDBFiles(bucket, c.TsdbPath)

	if err != nil {
		fmt.Printf("Failed to upload TSDB files: %v\n", err)
	} else {
		fmt.Println("TSDB files successfully uploaded to object store.")
	}

	return err
}

func main() {
	app := kingpin.New(filepath.Base(os.Args[0]), "Tool for storing TSDB data to object storage")
	app.HelpFlag.Short('h')

	s := newstore()
	Objstore := app.Command("block-sync", `Using Thanos to store the data`)
	Objstore.Flag("path", "Path for The TSDB data in prometheus").Required().StringVar(&s.TsdbPath)
	Objstore.Flag("objstore.config-file", "Path for The Config file").Required().StringVar(&s.ObjectConfig)
	Objstore.Command("upload", "Uploading data").Action(s.fun)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
