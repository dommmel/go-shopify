package shopify

import (
	"bytes"

	"encoding/json"

	"fmt"

	"time"
)

type Asset struct {
	CreatedAt time.Time `json:"created_at"`

	Id          int64     `json:"id"`
	ContentType string    `json:"content_type"`
	Attachment  string    `json:"attachment"`
	Key         string    `json:"key"`
	PublicUrl   string    `json:"public_url"`
	ThemeId     int64     `json:"theme_id"`
	UpdatedAt   time.Time `json:"updated_at"`
	api         *API
}

// func (api *API) Assets(themeId int64) ([]Asset, error) {
// 	endpoint := fmt.Sprintf("/admin/themes/%d/assets.json", themeId)
// 	res, status, err := api.request(endpoint, "GET", nil, nil)

// 	if err != nil {
// 		return nil, err
// 	}

// 	if status != 200 {
// 		return nil, fmt.Errorf("Status returned: %d", status)
// 	}

// 	r := &map[string][]Asset{}
// 	err = json.NewDecoder(res).Decode(r)

// 	result := (*r)["asset"]

// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, v := range result {
// 		v.api = api
// 	}

// 	return result, nil
// }

// func (api *API) Asset(themeId int64) (*Asset, error) {
// 	endpoint := fmt.Sprintf("/admin/themes/%d/assets.json", themeId)

// 	res, status, err := api.request(endpoint, "GET", nil, nil)

// 	if err != nil {
// 		return nil, err
// 	}

// 	if status != 200 {
// 		return nil, fmt.Errorf("Status returned: %d", status)
// 	}

// 	r := map[string]Asset{}
// 	err = json.NewDecoder(res).Decode(&r)

// 	result := r["asset"]

// 	if err != nil {
// 		return nil, err
// 	}

// 	result.api = api

// 	return &result, nil
// }

func (api *API) NewAsset() *Asset {
	return &Asset{api: api}
}

func (obj *Asset) Save() error {
	endpoint := fmt.Sprintf("/admin/themes/%d/assets.json", obj.ThemeId)
	method := "PUT"
	expectedStatus := 200

	body := map[string]*Asset{}
	body["asset"] = obj

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

	r := map[string]Asset{}
	err = json.NewDecoder(res).Decode(&r)

	if err != nil {
		return err
	}

	*obj = r["asset"]

	return nil
}
