package shopify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"google.golang.org/appengine/log"
	"net/url"
	"time"
)

type Metafield struct {
	CreatedAt time.Time `json:"created_at,omitempty"`

	Description string `json:"description,omitempty"`

	Id int64 `json:"id,omitempty"`

	Key string `json:"key"`

	Namespace string `json:"namespace,omitempty"`

	OwnerId int64 `json:"owner_id,omitempty"`

	UpdatedAt time.Time `json:"updated_at,omitempty"`

	Value string `json:"value"`

	ValueType string `json:"value_type"`

	OwnerResource string `json:"owner_resource,omitempty"`

	api *API
}

func (api *API) Metafields(params ...url.Values) ([]Metafield, error) {
	return api.ResourceMetafields("", params[0])
}

func (api *API) ResourceMetafields(resourceName string, params url.Values) ([]Metafield, error) {
	if resourceName != "" {
		resourceName = resourceName + "/"
	}
	encodedParams := ""
	if params != nil {
		encodedParams = params.Encode()
	}
	endpoint := fmt.Sprintf("/admin/%smetafields.json?%s", resourceName, encodedParams)
	log.Errorf(api.Context, endpoint)
	res, status, err := api.request(endpoint, "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("Status returned: %d", status)
	}

	r := &map[string][]Metafield{}
	err = json.NewDecoder(res).Decode(r)

	result := (*r)["metafields"]

	if err != nil {
		return nil, err
	}

	for _, v := range result {
		v.api = api
	}

	return result, nil
}

func (api *API) Metafield(id int64) (*Metafield, error) {
	endpoint := fmt.Sprintf("/admin/metafields/%d.json", id)

	res, status, err := api.request(endpoint, "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("Status returned: %d", status)
	}

	r := map[string]Metafield{}
	err = json.NewDecoder(res).Decode(&r)

	result := r["metafield"]

	if err != nil {
		return nil, err
	}

	result.api = api

	return &result, nil
}

func (api *API) NewMetafield() *Metafield {
	return &Metafield{api: api}
}

func (obj *Metafield) Save() error {
	endpoint := fmt.Sprintf("/admin/metafields/%d.json", obj.Id)
	method := "PUT"
	expectedStatus := 200

	if obj.Id == 0 {
		endpoint = fmt.Sprintf("/admin/metafields.json")
		method = "POST"
		expectedStatus = 201
	}

	body := map[string]*Metafield{}
	body["metafield"] = obj

	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(body)

	if err != nil {
		return err
	}

	res, status, err := obj.api.request(endpoint, method, nil, buf)

	if err != nil {
		return err
	}

	if status != expectedStatus {
		r := errorResponse{}
		err = json.NewDecoder(res).Decode(&r)
		if err == nil {
			return fmt.Errorf("Status %d: %v", status, r.Errors)
		} else {
			return fmt.Errorf("Status %d, and error parsing body: %s", status, err)
		}
	}

	r := map[string]Metafield{}
	err = json.NewDecoder(res).Decode(&r)

	if err != nil {
		return err
	}

	*obj = r["metafield"]

	return nil
}
