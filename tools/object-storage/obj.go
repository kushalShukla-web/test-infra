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
	"os"
	"path/filepath"

	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	app := kingpin.New(filepath.Base(os.Args[0]), "Tool for storing TSDB data to object storage")
	app.HelpFlag.Short('h')

	s := newstore()
	Objstore := app.Command("block-sync", `Using Thanos to store the data`)
	Objstore.Flag("path", "Path for The TSDB data in prometheus").Required().StringVar(&s.TsdbPath)
	Objstore.Flag("objstore.config-file", "Path for The Config file").Required().StringVar(&s.ObjectConfig)
	Objstore.Flag("key", "Path for the Key where to store block data").Required().StringVar(&s.ObjectKey)
	Objstore.Command("upload", "Uploading data").Action(s.upload)
	Objstore.Command("download", "Downloading data").Action(s.download)
	kingpin.MustParse(app.Parse(os.Args[1:]))
}
