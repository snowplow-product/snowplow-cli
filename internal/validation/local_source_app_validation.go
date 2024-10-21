/**
 * Copyright (c) 2013-present Snowplow Analytics Ltd.
 * All rights reserved.
 * This software is made available by Snowplow Analytics, Ltd.,
 * under the terms of the Snowplow Limited Use License Agreement, Version 1.0
 * located at https://docs.snowplow.io/limited-use-license-1.0
 * BY INSTALLING, DOWNLOADING, ACCESSING, USING OR DISTRIBUTING ANY PORTION
 * OF THE SOFTWARE, YOU AGREE TO THE TERMS OF SUCH LICENSE AGREEMENT.
 */

package validation

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/snowplow-product/snowplow-cli/internal/model"
)

func ValidateSAMinimum(sa model.SourceApp) DPValidations {
	errors := []string{}

	_, err := uuid.Parse(sa.ResourceName)
	if err != nil {
		errors = append(errors, "resourceName must be a valid uuid")
	}

	if len(sa.Data.Name) == 0 {
		errors = append(errors, "data.name required")
	}

	return DPValidations{errors, []string{}, []string{}, []string{}}
}

func ValidateSAAppIds(sa model.SourceApp) DPValidations {
	errors := []string{}

	for i, a := range sa.Data.AppIds {
		if len(a) == 0 {
			errors = append(errors, fmt.Sprintf("data.appIds[%d] can't be empty", i))
		}
	}

	return DPValidations{errors, []string{}, []string{}, []string{}}
}

func ValidateSAEntitiesSources(sa model.SourceApp) DPValidations {
	errors := []string{}

	if sa.Data.Entities == nil {
		return DPValidations{}
	}

	for i, e := range sa.Data.Entities.Tracked {
		if len(e.Source) == 0 {
			errors = append(errors, fmt.Sprintf("data.entities.tracked[%d].source required", i))
		}
	}

	for i, e := range sa.Data.Entities.Enriched {
		if len(e.Source) == 0 {
			errors = append(errors, fmt.Sprintf("data.entities.enriched[%d].source required", i))
		}
	}

	return DPValidations{errors, []string{}, []string{}, []string{}}
}

func cardinalityCheck(key string, i int, s model.SchemaRef) []string {
	errors := []string{}
	if s.MinCardinality != nil {
		if *s.MinCardinality < 0 {
			errors = append(errors, fmt.Sprintf("data.entities.%s[%d].minCardinality must be > 0", key, i))
		}
		if s.MaxCardinality != nil {
			if *s.MaxCardinality < *s.MinCardinality {
				errors = append(errors, fmt.Sprintf("data.entities.%s[%d].maxCardinality must be > minCardinality: %d", key, i, *s.MinCardinality))
			}
		}
	} else {
		if s.MaxCardinality != nil {
			errors = append(errors, fmt.Sprintf("data.entities.%s[%d].maxCardinality without minCardinality", key, i))
		}
	}
	return errors
}

func ValidateSAEntitiesCardinalities(sa model.SourceApp) DPValidations {
	errors := []string{}

	if sa.Data.Entities == nil {
		return DPValidations{}
	}

	for i, e := range sa.Data.Entities.Tracked {
		errors = append(errors, cardinalityCheck("tracked", i, e)...)
	}

	for i, e := range sa.Data.Entities.Enriched {
		errors = append(errors, cardinalityCheck("enriched", i, e)...)
	}

	return DPValidations{errors, []string{}, []string{}, []string{}}
}

func ValidateSAEntitiesHaveNoRules(sa model.SourceApp) DPValidations {
	errors := []string{}

	if sa.Data.Entities == nil {
		return DPValidations{}
	}

	for i, e := range sa.Data.Entities.Tracked {
		if e.Schema != nil {
			errors = append(
				errors,
				fmt.Sprintf("data.entities.tracked[%d].schema property rules unsupported for source applications", i),
			)
		}
	}

	for i, e := range sa.Data.Entities.Enriched {
		if e.Schema != nil {
			errors = append(
				errors,
				fmt.Sprintf("data.entities.enriched[%d].schema property rules unsupported for source applications", i),
			)
		}
	}

	return DPValidations{errors, []string{}, []string{}, []string{}}
}