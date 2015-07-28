package shopify


import (
  
    "bytes"
  
    "encoding/json"
  
    "fmt"
  
    "time"
  
)


type CustomCollection struct {
  
    BodyHtml string `json:"body_html"`
  
    Handle string `json:"handle"`
  
    Id int64 `json:"id"`
  
    PublishedAt time.Time `json:"published_at.omitempty"`
  
    PublishedScope string `json:"published_scope"`
  
    SortOrder string `json:"sort_order"`
  
    TemplateSuffix string `json:"template_suffix"`
  
    Title string `json:"title"`
  
    UpdatedAt time.Time `json:"updated_at,omitempty"`
  

  
    api *API
  
}


func (api *API) CustomCollections() ([]CustomCollection, error) {
  res, status, err := api.request("/admin/custom_collections.json?limit=250", "GET", nil, nil)

  if err != nil {
    return nil, err
  }

  if status != 200 {
    return nil, fmt.Errorf("Status returned: %d", status)
  }

  r := &map[string][]CustomCollection{}
  err = json.NewDecoder(res).Decode(r)

  fmt.Printf("things are: %v\n\n", *r)

  result := (*r)["custom_collections"]

	if err != nil {
		return nil, err
  }

  for _, v := range result {
    v.api = api
  }

  return result, nil
}




func (api *API) CustomCollection(id int64) (*CustomCollection, error) {
  endpoint := fmt.Sprintf("/admin/custom_collections/%d.json", id)

  res, status, err := api.request(endpoint, "GET", nil, nil)

  if err != nil {
    return nil, err
  }

  if status != 200 {
    return nil, fmt.Errorf("Status returned: %d", status)
  }

  r := map[string]CustomCollection{}
  err = json.NewDecoder(res).Decode(&r)

  fmt.Printf("things are: %v\n\n", r)

  result := r["custom_collection"]

	if err != nil {
		return nil, err
  }

  result.api = api

  return &result, nil
}


func (api *API) NewCustomCollection() *CustomCollection {
  return &CustomCollection{api: api}
}


func (obj *CustomCollection) Save() (error) {
  endpoint := fmt.Sprintf("/admin/custom_collections/%d.json", obj.Id)
  method := "PUT"
  expectedStatus := 201

  if obj.Id == 0 {
    endpoint = fmt.Sprintf("/admin/custom_collections.json")
    method = "POST"
    expectedStatus = 201
  }

  body := map[string]*CustomCollection{}
  body["custom_collection"] = obj

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

  r := map[string]CustomCollection{}
  err = json.NewDecoder(res).Decode(&r)

	if err != nil {
		return err
  }

  fmt.Printf("things are: %v\n\n", r)

  *obj = r["custom_collection"]

  fmt.Printf("things are: %v\n\n", res)

  return nil
}





