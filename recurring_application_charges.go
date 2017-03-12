package shopify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type RecurringApplicationCharge struct {
	Id          int64  `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	ApiClientId int64  `json:"api_client_id,omitempty"`

	UpdatedAt          *time.Time `json:"updated_at,omitempty"`
	CreatedAt          *time.Time `json:"created_at,omitempty"`
	ActivatedOn        *IsoDate   `json:"activated_on,omitempty"`
	BillingOn          *IsoDate   `json:"billing_on,omitempty"`
	CancelledOn        *IsoDate   `json:"cancelled_on,omitempty"`
	TrialDays          int64      `json:"trial_days,omitempty"`
	TrialEndsOn        *IsoDate   `json:"trial_ends_on,omitempty"`
	DecoratedReturnUrl string     `json:"decorated_return_url,omitempty"`
	ConfirmationUrl    *string    `json:"confirmation_url,omitempty"`
	ReturnUrl          string     `json:"return_url,omitempty"`
	Price              string     `json:"price,omitempty"`
	CappedAmount       string     `json:"capped_amount,omitempty"`
	Terms              string     `json:"terms,omitempty"`
	Test               *bool      `json:"test,omitempty"`
	Status             string     `json:"status,omitempty"`
	api                *API
}

func (api *API) RecurringApplicationCharges() ([]RecurringApplicationCharge, error) {
	res, status, err := api.request("/admin/recurring_application_charges.json", "GET", nil, nil)

	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("Status returned: %d", status)
	}

	//api.Context.Infof("RESPONSE---- %s", res)
	r := &map[string][]RecurringApplicationCharge{}
	err = json.NewDecoder(res).Decode(r)
	if err != nil {
		return nil, err
	}

	result := (*r)["recurring_application_charges"]
	//api.Context.Infof("CHARGES ARE: %s\n\n", r)

	for _, v := range result {
		v.api = api
	}

	return result, nil
}

func (api *API) RecurringApplicationCharge(id string) (*RecurringApplicationCharge, error) {
	endpoint := fmt.Sprintf("/admin/recurring_application_charges/%s.json", id)

	res, status, err := api.request(endpoint, "GET", nil, nil)
	if err != nil {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("Status returned: %d", status)
	}

	r := map[string]RecurringApplicationCharge{}
	err = json.NewDecoder(res).Decode(&r)

	result := r["recurring_application_charge"]

	if err != nil {
		return nil, err
	}

	result.api = api

	return &result, nil
}

func (api *API) NewRecurringApplicationCharge() *RecurringApplicationCharge {
	return &RecurringApplicationCharge{api: api}
}

func (api *API) DeleteRecurringApplicationCharge(chargeId int64) error {
	endpoint := fmt.Sprintf("/admin/recurring_application_charges/%d.json", chargeId)
	//return fmt.Errorf("OBJ %v", &obj.api)
	method := "DELETE"
	expectedStatus := 200

	res, status, err := api.request(endpoint, method, nil, nil)

	if status != expectedStatus {
		r := errorResponse{}
		err = json.NewDecoder(res).Decode(&r)
		if err == nil {
			return fmt.Errorf("Status %d: %v", status, r.Errors)
		} else {
			return fmt.Errorf("Status %d, and error parsing body: %s", status, err)
		}
	}
	return nil
}

func (obj *RecurringApplicationCharge) Activate() error {
	endpoint := fmt.Sprintf("/admin/recurring_application_charges/%d/activate.json", obj.Id)
	method := "POST"
	expectedStatus := 201

	res, status, err := obj.api.request(endpoint, method, nil, nil)

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

	r := map[string]RecurringApplicationCharge{}
	err = json.NewDecoder(res).Decode(&r)

	if err != nil {
		return err
	}

	*obj = r["recurring_application_charge"]

	return nil
}

func (obj *RecurringApplicationCharge) Save() error {
	endpoint := fmt.Sprintf("/admin/recurring_application_charges/%d.json", obj.Id)
	method := "PUT"
	expectedStatus := 201

	if obj.Id == 0 {
		endpoint = fmt.Sprintf("/admin/recurring_application_charges.json")
		method = "POST"
		expectedStatus = 201
	}

	body := map[string]*RecurringApplicationCharge{}
	body["recurring_application_charge"] = obj

	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(body)

	if err != nil {
		return err
	}

	obj.api.Context.Infof("REQUEST BODY: %v", buf)

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

	r := map[string]RecurringApplicationCharge{}
	err = json.NewDecoder(res).Decode(&r)

	if err != nil {
		return err
	}

	*obj = r["recurring_application_charge"]

	return nil
}
