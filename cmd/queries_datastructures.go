package cmd

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

type msgResponse struct {
	Message string
}

type pubResponse struct {
	Success  bool
	Errors   []string
	Warnings []string
	Info     []string
	msgResponse
}

type pubError struct {
	Messages []string
}

type PublishResponse = pubResponse
type ValidateResponse struct {
	pubResponse
	Valid bool
}

type PublishError = pubError

func (e *pubError) Error() string {
	return strings.Join(e.Messages, "\n")
}

type dataStructureEnv string

const (
	DEV       dataStructureEnv = "DEV"
	PROD      dataStructureEnv = "PROD"
	VALIDATED dataStructureEnv = "VALIDATED"
)

type publishRequest struct {
	Format  string           `json:"format"`
	Message string           `json:"message"`
	Name    string           `json:"name"`
	Source  dataStructureEnv `json:"source"`
	Target  dataStructureEnv `json:"target"`
	Vendor  string           `json:"vendor"`
	Version string           `json:"version"`
}

func Validate(cnx context.Context, client *ApiClient, ds DataStructure) (*ValidateResponse, error) {

	body, err := json.Marshal(ds)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(cnx, "POST", fmt.Sprintf("%s/data-structures/v1/validation-requests", client.BaseUrl), bytes.NewBuffer(body))
	auth := fmt.Sprintf("Bearer %s", client.Jwt)
	req.Header.Add("authorization", auth)

	if err != nil {
		return nil, err
	}
	resp, err := client.Http.Do(req)
	if err != nil {
		return nil, err
	}
	rbody, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	var vresp ValidateResponse
	err = json.Unmarshal(rbody, &vresp)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, errors.New(vresp.Message)
	}

	vresp.Valid = vresp.Success

	return &vresp, nil
}

func PublishDev(cnx context.Context, client *ApiClient, ds DataStructure, isPatch bool) (*PublishResponse, error) {
	return publish(cnx, client, VALIDATED, DEV, ds, isPatch)
}

func PublishProd(cnx context.Context, client *ApiClient, ds DataStructure) (*PublishResponse, error) {
	return publish(cnx, client, DEV, PROD, ds, false)
}

func publish(cnx context.Context, client *ApiClient, from dataStructureEnv, to dataStructureEnv, ds DataStructure, isPatch bool) (*PublishResponse, error) {

	dsData, err := ds.parseData()
	if err != nil {
		return nil, err
	}

	pr := &publishRequest{
		Message: "",
		Source:  from,
		Target:  to,
		Vendor:  dsData.Self.Vendor,
		Name:    dsData.Self.Name,
		Format:  dsData.Self.Format,
		Version: dsData.Self.Version,
	}

	body, err := json.Marshal(pr)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(cnx, "POST", fmt.Sprintf("%s/data-structures/v1/deployment-requests", client.BaseUrl), bytes.NewBuffer(body))
	auth := fmt.Sprintf("Bearer %s", client.Jwt)
	req.Header.Add("authorization", auth)
	if isPatch {
		q := req.URL.Query()
		q.Add("patch", "true")
		req.URL.RawQuery = q.Encode()
	}

	if err != nil {
		return nil, err
	}
	resp, err := client.Http.Do(req)
	if err != nil {
		return nil, err
	}
	rbody, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	var dresp PublishResponse
	err = json.Unmarshal(rbody, &dresp)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, errors.New(dresp.Message)
	}

	if !dresp.Success {
		return nil, &PublishError{Messages: dresp.Errors}
	}

	return &dresp, nil
}

type Deployment struct {
	Version     string           `json:"version"`
	Env         dataStructureEnv `json:"env"`
	ContentHash string           `json:"contentHash"`
}

type ListResponse struct {
	Hash        string            `json:"hash"`
	Vendor      string            `json:"vendor"`
	Format      string            `json:"format"`
	Name        string            `json:"name"`
	Meta        DataStructureMeta `json:"meta"`
	Deployments []Deployment      `json:"deployments"`
}

func GetDataStructureListing(cnx context.Context, client *ApiClient) ([]ListResponse, error) {
	req, err := http.NewRequestWithContext(cnx, "GET", fmt.Sprintf("%s/data-structures/v1", client.BaseUrl), nil)
	auth := fmt.Sprintf("Bearer %s", client.Jwt)
	req.Header.Add("authorization", auth)

	if err != nil {
		return nil, err
	}
	resp, err := client.Http.Do(req)
	if err != nil {
		return nil, err
	}
	rbody, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	var listResp []ListResponse
	err = json.Unmarshal(rbody, &listResp)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("not expected response code %d", resp.StatusCode)
	}
	return listResp, nil
}

func GetAllDataStructures(cnx context.Context, client *ApiClient) ([]DataStructure, error) {

	req, err := http.NewRequestWithContext(cnx, "GET", fmt.Sprintf("%s/data-structures/v1", client.BaseUrl), nil)
	auth := fmt.Sprintf("Bearer %s", client.Jwt)
	req.Header.Add("authorization", auth)

	if err != nil {
		return nil, err
	}

	listResp, err := GetDataStructureListing(cnx, client)
	if err != nil {
		return nil, err
	}

	var res []DataStructure

	for _, dsResp := range listResp {
		for _, deployment := range dsResp.Deployments {
			if deployment.Env == DEV {
				req, err := http.NewRequestWithContext(cnx, "GET", fmt.Sprintf("%s/data-structures/v1/%s/versions/%s", client.BaseUrl, dsResp.Hash, deployment.Version), nil)
				auth := fmt.Sprintf("Bearer %s", client.Jwt)
				req.Header.Add("authorization", auth)
				slog.Info("fetching data structure", "uri", fmt.Sprintf("iglu:%s/%s/%s/%s", dsResp.Vendor, dsResp.Name, dsResp.Format, deployment.Version))

				if err != nil {
					return nil, err
				}
				resp, err := client.Http.Do(req)
				if err != nil {
					return nil, err
				}
				rbody, err := io.ReadAll(resp.Body)
				defer resp.Body.Close()
				if err != nil {
					return nil, err
				}

				var ds map[string]any
				err = json.Unmarshal(rbody, &ds)
				if err != nil {
					return nil, err
				}

				if resp.StatusCode == http.StatusNotFound {
					continue
				}

				if resp.StatusCode != http.StatusOK {
					return nil, fmt.Errorf("not expected response code %d", resp.StatusCode)
				}

				dataStructure := DataStructure{"v1", "data-structure", dsResp.Meta, ds}
				res = append(res, dataStructure)
			}
		}
	}

	return res, nil
}

func MetadateUpdate(cnx context.Context, client *ApiClient, ds *DataStructure) error {

	data, err := ds.parseData()
	if err != nil {
		return err
	}

	toHash := fmt.Sprintf("%s-%s-%s-%s", client.OrgId, data.Self.Vendor, data.Self.Name, data.Self.Format)
	dsHash := sha256.Sum256([]byte(toHash))

	body, err := json.Marshal(ds.Meta)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/data-structures/v1/%x/meta", client.BaseUrl, dsHash)
	req, err := http.NewRequestWithContext(cnx, "PUT", url, bytes.NewBuffer(body))
	auth := fmt.Sprintf("Bearer %s", client.Jwt)
	req.Header.Add("authorization", auth)

	if err != nil {
		return err
	}

	resp, err := client.Http.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {

		rbody, err := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			return err
		}

		var dresp msgResponse
		err = json.Unmarshal(rbody, &dresp)
		if err != nil {
			return errors.Join(err, errors.New("bad response with no message"))
		}

		return fmt.Errorf("bad response: %s", dresp.Message)
	}

	return nil
}
