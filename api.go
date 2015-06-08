package shopify

import (
  "log"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"appengine"
	"appengine/urlfetch"
)

type API struct {
	URI    string
	Token  string
	Secret string
	client *http.Client
	Context appengine.Context
}

type errorResponse struct {
	Errors map[string]interface{} `json:"errors"`
}

func (api *API) request(endpoint string, method string, params map[string]interface{}, body io.Reader) (result *bytes.Buffer, status int, err error) {
	if api.client == nil {
		api.client = urlfetch.Client(api.Context)
	}

	uri := fmt.Sprintf("%s/%s", api.URI, endpoint)
	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		return
	}
	log.Printf(api.Token)
	req.Header.Add("X-Shopify-Access-Token", api.Token)
  req.Header.Add("Content-Type", "application/json")

	resp, err := api.client.Do(req)
	fmt.Printf("resp %v err %v", resp, err)
	if err != nil {
		return
	}

	status = resp.StatusCode

	result = &bytes.Buffer{}
	if _, err = io.Copy(result, resp.Body); err != nil {
		return
	}

	return
}
