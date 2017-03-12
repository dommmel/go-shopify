package shopify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type UsageCharge struct {
	Id                           int64      `json:"id,omitempty"`
	ApiClientId                  int64      `json:"api_client_id,omitempty"`
	UpdatedAt                    *time.Time `json:"updated_at,omitempty"`
	CreatedAt                    *time.Time `json:"created_at,omitempty"`
	RecurringApplicationChargeId int64      `json:"recurring_application_charge_id,omitempty"`
	Price                        string     `json:"price,omitempty"`
	CappedAmount                 string     `json:"capped_amount,omitempty"`
	Description                  string     `json:"description,omitempty"`
	api                          *API
}

func (api *API) UsageCharges(recurringApplicationChargeId int64) ([]UsageCharge, error) {
	endpoint := fmt.Sprintf("/admin/recurring_application_charges/%d/usage_charges.json", recurringApplicationChargeId)
	res, status, err := api.request(endpoint, "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("Status returned: %d", status)
	}

	//api.Context.Infof("RESPONSE---- %s", res)
	r := &map[string][]UsageCharge{}
	err = json.NewDecoder(res).Decode(r)
	if err != nil {
		return nil, err
	}

	result := (*r)["usage_charges"]
	//api.Context.Infof("CHARGES ARE: %s\n\n", r)

	for _, v := range result {
		v.api = api
	}

	return result, nil
}

func (api *API) NewUsageCharge() *UsageCharge {
	return &UsageCharge{api: api}
}

func (obj *UsageCharge) Save() error {
	endpoint := fmt.Sprintf("/admin/recurring_application_charges/%d/usage_charges.json", obj.RecurringApplicationChargeId)
	//obj.api.Context.Infof("Endpoint Usage Charge: %s", endpoint)
	method := "POST"
	expectedStatus := 201

	body := map[string]*UsageCharge{}
	body["usage_charge"] = obj

	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(body)

	if err != nil {
		return err
	}

	//obj.api.Context.Infof("REQUEST BODY: %v", buf)

	res, status, err := obj.api.request(endpoint, method, nil, buf)

	json.NewDecoder(res)
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
	//obj.api.Context.Infof("RESPONSE BODY: %v", res)
	r := map[string]UsageCharge{}
	err = json.NewDecoder(res).Decode(&r)

	if err != nil {
		return err
	}

	*obj = r["usage_charge"]

	return nil
}
