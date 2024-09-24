/*
Copyright (c) 2013-present Snowplow Analytics Ltd.
All rights reserved.
This software is made available by Snowplow Analytics, Ltd.,
under the terms of the Snowplow Limited Use License Agreement, Version 1.0
located at https://docs.snowplow.io/limited-use-license-1.0
BY INSTALLING, DOWNLOADING, ACCESSING, USING OR DISTRIBUTING ANY PORTION
OF THE SOFTWARE, YOU AGREE TO THE TERMS OF SUCH LICENSE AGREEMENT.
*/

package util

import (
	"fmt"
	. "github.com/snowplow-product/snowplow-cli/internal/model"
	"os"
	"path/filepath"
	"testing"
)

func TestCreatesDataStructuresFolderWithFiles(t *testing.T) {
	extension := "yaml"
	vendor1 := "test.my.vendor"
	name1 := "my-test-ds"
	ds1 := DataStructure{
		Meta: DataStructureMeta{Hidden: true, SchemaType: "entity", CustomData: map[string]string{
			"additionalProp1": "string",
			"additionalProp2": "string",
			"additionalProp3": "string",
		},
		},
		Data: map[string]any{
			"self": map[string]any{
				"vendor":  vendor1,
				"name":    name1,
				"format":  "string",
				"version": "1-2-0",
			},
			"schema": "string"},
	}
	vendor2 := "com.test.vendor"
	name2 := "ds2"
	ds2 := DataStructure{
		Meta: DataStructureMeta{Hidden: true, SchemaType: "entity", CustomData: map[string]string{
			"additionalProp1": "string",
			"additionalProp2": "string",
			"additionalProp3": "string",
		},
		},
		Data: map[string]any{
			"self": map[string]any{
				"vendor":  vendor2,
				"name":    name2,
				"format":  "string",
				"version": "1-0-1",
			},
			"schema": "string"},
	}

	dir := filepath.Join("../..", "out", "test-ds2")
	files := Files{DataStructuresLocation: dir, ExtentionPreference: extension}
	err := files.CreateDataStructures([]DataStructure{ds1, ds2})

	if err != nil {
		t.Fatalf("Can't create directory %s", err)
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Fatalf("%s does not exists", dir)
	}

	filePath1 := filepath.Join(dir, vendor1, fmt.Sprintf("%s.%s", name1, extension))
	if _, err := os.Stat(filePath1); os.IsNotExist(err) {
		t.Fatalf("%s does not exists", filePath1)
	}

	filePath2 := filepath.Join(dir, vendor2, fmt.Sprintf("%s.%s", name2, extension))
	if _, err := os.Stat(filePath2); os.IsNotExist(err) {
		t.Fatalf("%s does not exists", filePath2)
	}

}

func TestCreatesDataStructuresFolderWithFilesJson(t *testing.T) {
	extension := "json"
	vendor1 := "test.my.vendor"
	name1 := "my-test-ds"
	ds1 := DataStructure{
		Meta: DataStructureMeta{Hidden: true, SchemaType: "entity", CustomData: map[string]string{
			"additionalProp1": "string",
			"additionalProp2": "string",
			"additionalProp3": "string",
		},
		},
		Data: map[string]any{
			"self": map[string]any{
				"vendor":  vendor1,
				"name":    name1,
				"format":  "string",
				"version": "1-2-0",
			},
			"schema": "string"},
	}
	vendor2 := "com.test.vendor"
	name2 := "ds2"
	ds2 := DataStructure{
		Meta: DataStructureMeta{Hidden: true, SchemaType: "entity", CustomData: map[string]string{
			"additionalProp1": "string",
			"additionalProp2": "string",
			"additionalProp3": "string",
		},
		},
		Data: map[string]any{
			"self": map[string]any{
				"vendor":  vendor2,
				"name":    name2,
				"format":  "string",
				"version": "1-0-1",
			},
			"schema": "string"},
	}

	dir := filepath.Join("../..", "out", "test-ds2")
	files := Files{DataStructuresLocation: dir, ExtentionPreference: extension}
	err := files.CreateDataStructures([]DataStructure{ds1, ds2})

	if err != nil {
		t.Fatalf("Can't create directory %s", err)
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Fatalf("%s does not exists", dir)
	}

	filePath1 := filepath.Join(dir, vendor1, fmt.Sprintf("%s.%s", name1, extension))
	if _, err := os.Stat(filePath1); os.IsNotExist(err) {
		t.Fatalf("%s does not exists", filePath1)
	}

	filePath2 := filepath.Join(dir, vendor2, fmt.Sprintf("%s.%s", name2, extension))
	if _, err := os.Stat(filePath2); os.IsNotExist(err) {
		t.Fatalf("%s does not exists", filePath2)
	}

}