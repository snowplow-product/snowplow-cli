/*
Copyright (c) 2013-present Snowplow Analytics Ltd.
All rights reserved.
This software is made available by Snowplow Analytics, Ltd.,
under the terms of the Snowplow Limited Use License Agreement, Version 1.0
located at https://docs.snowplow.io/limited-use-license-1.0
BY INSTALLING, DOWNLOADING, ACCESSING, USING OR DISTRIBUTING ANY PORTION
OF THE SOFTWARE, YOU AGREE TO THE TERMS OF SUCH LICENSE AGREEMENT.
*/
package download

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/snowplow-product/snowplow-cli/internal/console"
	"github.com/snowplow-product/snowplow-cli/internal/model"
)

func intPtr(i int) *int {
	return &i
}

var sampleEntity1 = console.Entity{
	Source:         "iglu:com.snplow.msc.aws/data-product-auto/jsonschema/1-0-0",
	MinCardinality: intPtr(0),
	MaxCardinality: intPtr(5),
	Schema:         map[string]any{},
}

var sampleEntity2 = console.Entity{
	Source:         "iglu:com.snplow.msc.aws/event-spec-auto/jsonschema/1-0-0",
	MinCardinality: intPtr(0),
	MaxCardinality: nil,
	Schema:         map[string]any{},
}

var sampleSa1 = console.RemoteSourceApplication{
	Id:          "6b1146d6-7b23-4dbb-b069-f568458dda8f",
	Name:        "test",
	Description: "my test sa",
	Owner:       "me@me.com",
	AppIds:      []string{"ios", "android"},
	Entities:    console.Entities{Tracked: []console.Entity{sampleEntity1, sampleEntity2}, Enriched: []console.Entity{}},
}

var sampleSaRef = model.Ref{
	Ref: "./source-apps/test.yaml",
}

var sampleSaRefs = map[string]model.Ref{
	sampleSa1.Id: sampleSaRef,
}

var sampleSaRefs2 = map[string]model.Ref{
	sampleSa1.Id: sampleSaRef,
	"8fb370f4-60af-4b7d-9438-ea08df6cdc70": model.Ref{
		Ref: "./source-apps/wrong-sa.yaml",
	},
}

var trackedEntites = []model.SchemaRef{
	{
		Source:         sampleEntity1.Source,
		MinCardinality: sampleEntity1.MinCardinality,
		MaxCardinality: sampleEntity1.MaxCardinality,
		Schema:         map[string]any{},
	},
	{
		Source:         sampleEntity2.Source,
		MinCardinality: sampleEntity2.MinCardinality,
		MaxCardinality: sampleEntity2.MaxCardinality,
		Schema:         map[string]any{},
	},
}
var entities = model.EntitiesDef{Tracked: trackedEntites, Enriched: nil}
var SaResource = []model.CliResource[model.SourceAppData]{
	{
		ApiVersion:   "v1",
		ResourceType: "source-application",
		ResourceName: sampleSa1.Id,
		Data: model.SourceAppData{
			ResourceName: sampleSa1.Id,
			Name:         sampleSa1.Name,
			Description:  sampleSa1.Description,
			Owner:        sampleSa1.Owner,
			AppIds:       sampleSa1.AppIds,
			Entities:     &entities,
		},
	},
}

var sampleRemoteEs = console.RemoteEventSpec{
	Id:                   "84614b3b-6039-458e-8ce2-615eaf2113e3",
	SourceApplicationIds: []string{sampleSa1.Id},
	Name:                 "test ES",
	Triggers:             []console.Trigger{},
	Event: console.Event{
		Source: "iglu:com.yalo.schemas.events.channel/YaloMessage/jsonschema/1-0-0",
		Schema: map[string]any{},
	},
	Entities: console.Entities{Tracked: []console.Entity{
		{
			Source:         sampleEntity1.Source,
			MinCardinality: sampleEntity1.MinCardinality,
			MaxCardinality: sampleEntity1.MaxCardinality,
			Schema:         map[string]any{},
		},
		{
			Source:         sampleEntity2.Source,
			MinCardinality: sampleEntity2.MinCardinality,
			MaxCardinality: sampleEntity2.MaxCardinality,
			Schema:         map[string]any{},
		},
	},
		Enriched: nil,
	},
}
var sampleRemoteEss = []console.RemoteEventSpec{sampleRemoteEs}

var sampleEsIdToEs = map[string]console.RemoteEventSpec{sampleRemoteEs.Id: sampleRemoteEs}

var sampleRemoteDp = console.RemoteDataProduct{
	Id:                   "46d47289-f3d5-4ef8-a82c-b19597e6e503",
	Name:                 "test DP",
	SourceApplicationIds: []string{sampleSa1.Id},
	Domain:               "testing",
	Owner:                "me@me.me",
	Description:          "this is a test",
	EventSpecifications:  []console.EventSpecReference{{Id: sampleRemoteEss[0].Id}},
}

var sampleRemoteDps = []console.RemoteDataProduct{sampleRemoteDp}

func Test_remoteSasToLocalResources_OK(t *testing.T) {
	res := remoteSasToLocalResources([]console.RemoteSourceApplication{sampleSa1})
	if !reflect.DeepEqual(SaResource, res) {
		t.Errorf("Unexpected local source apps expected:%+v got:%+v", SaResource, res)
	}

}

func Test_localSasToRefs_OK(t *testing.T) {
	dpLocation := "my-folder"

	fileNamesToLocalSas := map[string]model.CliResource[model.SourceAppData]{
		fmt.Sprintf("%s/source-apps/test.yaml", dpLocation): SaResource[0],
	}
	res := localSasToRefs(fileNamesToLocalSas, dpLocation)

	if !reflect.DeepEqual(sampleSaRefs, res) {
		t.Errorf("Unexpected source app references expected:%+v got:%+v", sampleSaRefs, res)
	}

}

func Test_remoteDpsToLocalResources_OK(t *testing.T) {
	res := remoteDpsToLocalResources(sampleRemoteDps, sampleSaRefs2, sampleEsIdToEs)
	expected := []model.CliResource[model.DataProductCanonicalData]{{
		ApiVersion:   "v1",
		ResourceType: "data-product",
		ResourceName: sampleRemoteDp.Id,
		Data: model.DataProductCanonicalData{
			ResourceName:       sampleRemoteDp.Id,
			Name:               sampleRemoteDp.Name,
			SourceApplications: []model.Ref{sampleSaRef},
			Domain:             sampleRemoteDp.Domain,
			Owner:              sampleRemoteDp.Owner,
			Description:        sampleRemoteDp.Description,
			EventSpecifications: []model.EventSpecCanonical{{
				ResourceName:       sampleRemoteEs.Id,
				SourceApplications: []model.Ref{sampleSaRef},
				Name:               sampleRemoteEs.Name,
				Event: model.SchemaRef{
					Source:         sampleRemoteEs.Event.Source,
					MinCardinality: nil,
					MaxCardinality: nil,
					Schema:         map[string]any{},
				},
				Entities: model.EntitiesDef{
					Tracked:  trackedEntites,
					Enriched: nil,
				},
			}},
		},
	}}

	if !reflect.DeepEqual(expected, res) {
		t.Errorf("Unexpected local data products expected:%+v got:%+v", expected, res)
	}

}