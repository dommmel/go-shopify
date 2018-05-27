package shopify

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type ApplicationCredit struct {
	Id          int64  `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
	Amount      string `json:"amount,omitempty"`
	Test        *bool  `json:"test,omitempty"`
	api         *API
}

func (api *API) ApplicationCredits() ([]ApplicationCredit, error) {
	res, status, err := api.request("/admin/application_credits.json", "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("Status returned: %d", status)
	}

	r := &map[string][]ApplicationCredit{}
	err = json.NewDecoder(res).Decode(r)

	fmt.Printf("things are: %v\n\n", *r)

	result := (*r)["application_credits"]

	if err != nil {
		return nil, err
	}

	for _, v := range result {
		v.api = api
	}

	return result, nil
}

func (api *API) ApplicationCredit(id int64) (*ApplicationCredit, error) {
	endpoint := fmt.Sprintf("/admin/application_credits/%d.json", id)

	res, status, err := api.request(endpoint, "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("Status returned: %d", status)
	}

	r := map[string]ApplicationCredit{}
	err = json.NewDecoder(res).Decode(&r)

	fmt.Printf("things are: %v\n\n", r)

	result := r["application_credit"]

	if err != nil {
		return nil, err
	}

	result.api = api

	return &result, nil
}

func (api *API) NewApplicationCredit() *ApplicationCredit {
	return &ApplicationCredit{api: api}
}

func (obj *ApplicationCredit) Save() error {
	endpoint := fmt.Sprintf("/admin/application_credits/%d.json", obj.Id)
	method := "PUT"
	expectedStatus := 201

	if obj.Id == 0 {
		endpoint = fmt.Sprintf("/admin/application_credits.json")
		method = "POST"
		expectedStatus = 201
	}

	body := map[string]*ApplicationCredit{}
	body["application_credit"] = obj

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

	r := map[string]ApplicationCredit{}
	err = json.NewDecoder(res).Decode(&r)

	if err != nil {
		return err
	}

	fmt.Printf("things are: %v\n\n", r)

	*obj = r["application_credit"]

	fmt.Printf("things are: %v\n\n", res)

	return nil
}
