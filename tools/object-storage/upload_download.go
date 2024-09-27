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

func (c *Store) upload(*kingpin.ParseContext) error {
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

	err = objstore.UploadDir(context.Background(), log.NewNopLogger(), bucket, c.TsdbPath, commitDir)
	if err != nil {
		panic(err)
	}
	return nil
}

func (c *Store) download(*kingpin.ParseContext) error {
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
		return fmt.Errorf("failed to create bucket: %v", err)
	}
	fmt.Println(exists)

	err = objstore.DownloadDir(context.Background(), log.NewNopLogger(), bucket, "dir/", commitDir, c.TsdbPath)
	if err != nil {
		return fmt.Errorf("failed to download directory: %v", err)
	}
	return nil
}
