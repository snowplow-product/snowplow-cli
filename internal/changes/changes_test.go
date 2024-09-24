package changes

import (
	. "github.com/snowplow-product/snowplow-cli/internal/console"
	. "github.com/snowplow-product/snowplow-cli/internal/model"
	"testing"
)

func Test_GetChangesCreate(t *testing.T) {
	local := DataStructure{
		Meta: DataStructureMeta{Hidden: true, SchemaType: "entity"},
		Data: map[string]any{
			"self": map[string]any{
				"vendor":  "string",
				"name":    "string",
				"format":  "string",
				"version": "1-0-1",
			},
			"schema": "string"},
	}

	res, err := GetChanges(map[string]DataStructure{"file": local}, []ListResponse{}, "DEV")

	if err != nil {
		t.Fatalf("Can't calcuate changes %s", err)
	}

	if len(res.ToCreate) != 1 || len(res.ToUpdateMeta) != 0 || len(res.ToUpdatePatch) != 0 || len(res.ToUpdateNewVersion) != 0 {
		t.Fatalf("Unexpected result, expecting one data structre to be created, got %+v", res)
	}

}

func Test_GetChangesUpdateAndMeta(t *testing.T) {
	local := DataStructure{
		Meta: DataStructureMeta{Hidden: true, SchemaType: "entity"},
		Data: map[string]any{
			"self": map[string]any{
				"vendor":  "string",
				"name":    "string",
				"format":  "string",
				"version": "1-0-1",
			},
			"schema": "string"},
	}
	remote := ListResponse{
		Hash:   "different",
		Vendor: "string",
		Name:   "string",
		Meta:   DataStructureMeta{Hidden: false, SchemaType: "entity"},
		Format: "string",
		Deployments: []Deployment{
			{
				Version:     "1-0-0",
				Env:         "DEV",
				ContentHash: "different",
			},
		},
	}

	res, err := GetChanges(map[string]DataStructure{"file": local}, []ListResponse{remote}, "DEV")

	if err != nil {
		t.Fatalf("Can't calcuate changes %s", err)
	}

	if len(res.ToCreate) != 0 || len(res.ToUpdateMeta) != 1 || len(res.ToUpdatePatch) != 0 || len(res.ToUpdateNewVersion) != 1 {
		t.Fatalf("Unexpected result, expecting one data structre to update metadata and new version, got %+v", res)
	}

}

func Test_GetChangesPatch(t *testing.T) {
	local := DataStructure{
		Meta: DataStructureMeta{Hidden: true, SchemaType: "entity"},
		Data: map[string]any{
			"self": map[string]any{
				"vendor":  "string",
				"name":    "string",
				"format":  "string",
				"version": "1-0-0",
			},
			"schema": "string"},
	}
	remote := ListResponse{
		Hash:   "different",
		Vendor: "string",
		Name:   "string",
		Meta:   DataStructureMeta{Hidden: true, SchemaType: "entity"},
		Format: "string",
		Deployments: []Deployment{
			{
				Version:     "1-0-0",
				Env:         "DEV",
				ContentHash: "different",
			},
		},
	}

	res, err := GetChanges(map[string]DataStructure{"file": local}, []ListResponse{remote}, "DEV")

	if err != nil {
		t.Fatalf("Can't calcuate changes %s", err)
	}

	if len(res.ToCreate) != 0 || len(res.ToUpdateMeta) != 0 || len(res.ToUpdatePatch) != 1 || len(res.ToUpdateNewVersion) != 0 {
		t.Fatalf("Unexpected result, expecting one data structre to update patch, got %+v", res)
	}

}

func Test_GetChangesNoChange(t *testing.T) {
	local := DataStructure{
		Meta: DataStructureMeta{Hidden: true, SchemaType: "entity"},
		Data: map[string]any{
			"self": map[string]any{
				"vendor":  "string",
				"name":    "string",
				"format":  "string",
				"version": "1-0-0",
			},
			"schema": "string"},
	}
	remote := ListResponse{
		Hash:   "different",
		Vendor: "string",
		Name:   "string",
		Meta:   DataStructureMeta{Hidden: true, SchemaType: "entity"},
		Format: "string",
		Deployments: []Deployment{
			{
				Version:     "1-0-0",
				Env:         "DEV",
				ContentHash: "d11f7d148af26aa88594bfe37e94b34dd5d933bbd7a19be50a0cde5feb532e2a",
			},
		},
	}

	res, err := GetChanges(map[string]DataStructure{"file": local}, []ListResponse{remote}, "DEV")

	if err != nil {
		t.Fatalf("Can't calcuate changes %s", err)
	}

	if len(res.ToCreate) != 0 || len(res.ToUpdateMeta) != 0 || len(res.ToUpdatePatch) != 0 || len(res.ToUpdateNewVersion) != 0 {
		t.Fatalf("Unexpected result, expecting no changes, got %+v", res)
	}

}

func Test_GetChangesProdDeploy(t *testing.T) {
	local := DataStructure{
		Meta: DataStructureMeta{Hidden: true, SchemaType: "entity"},
		Data: map[string]any{
			"self": map[string]any{
				"vendor":  "string",
				"name":    "string",
				"format":  "string",
				"version": "1-0-0",
			},
			"schema": "string"},
	}
	remote := ListResponse{
		Hash:   "different",
		Vendor: "string",
		Name:   "string",
		Meta:   DataStructureMeta{Hidden: true, SchemaType: "entity"},
		Format: "string",
		Deployments: []Deployment{
			{
				Version:     "1-0-0",
				Env:         "DEV",
				ContentHash: "different",
			},
		},
	}

	res, err := GetChanges(map[string]DataStructure{"file": local}, []ListResponse{remote}, "PROD")

	if err != nil {
		t.Fatalf("Can't calcuate changes %s", err)
	}

	if len(res.ToCreate) != 0 || len(res.ToUpdateMeta) != 0 || len(res.ToUpdatePatch) != 0 || len(res.ToUpdateNewVersion) != 1 {
		t.Fatalf("Unexpected result, expecting one data structre to update, got %+v", res)
	}

}
