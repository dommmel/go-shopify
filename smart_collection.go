package shopify

import (
	"bytes"

	"encoding/json"

	"fmt"

	"time"
)

type SmartCollection struct {
	BodyHtml string `json:"body_html"`

	Disjunctive bool `json:"disjunctive"`

	Handle string `json:"handle"`

	Id int64 `json:"id"`

	PublishedAt time.Time `json:"published_at"`

	PublishedScope string `json:"published_scope"`

	SortOrder string `json:"sort_order"`

	TemplateSuffix string `json:"template_suffix"`

	Title string `json:"title"`

	UpdatedAt time.Time `json:"updated_at"`

	Rules []Rule `json:"rules"`

	api *API
}

func (api *API) SmartCollections() ([]SmartCollection, error) {
	res, status, err := api.request("/admin/smart_collections.json?limit=250", "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("Status returned: %d", status)
	}

	r := &map[string][]SmartCollection{}
	err = json.NewDecoder(res).Decode(r)

	result := (*r)["smart_collections"]

	if err != nil {
		return nil, err
	}

	for _, v := range result {
		v.api = api
	}

	return result, nil
}

func (api *API) SmartCollection(id int64) (*SmartCollection, error) {
	endpoint := fmt.Sprintf("/admin/smart_collections/%d.json", id)

	res, status, err := api.request(endpoint, "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("Status returned: %d", status)
	}

	r := map[string]SmartCollection{}
	err = json.NewDecoder(res).Decode(&r)

	result := r["smart_collection"]

	if err != nil {
		return nil, err
	}

	result.api = api

	return &result, nil
}

func (api *API) NewSmartCollection() *SmartCollection {
	return &SmartCollection{api: api}
}

func (obj *SmartCollection) Save() error {
	endpoint := fmt.Sprintf("/admin/smart_collections/%d.json", obj.Id)
	method := "PUT"
	expectedStatus := 201

	if obj.Id == 0 {
		endpoint = fmt.Sprintf("/admin/smart_collections.json")
		method = "POST"
		expectedStatus = 201
	}

	body := map[string]*SmartCollection{}
	body["smart_collection"] = obj

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

	r := map[string]SmartCollection{}
	err = json.NewDecoder(res).Decode(&r)

	if err != nil {
		return err
	}

	*obj = r["smart_collection"]

	return nil
}
