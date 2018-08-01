package shopify

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"

	"appengine"
	"appengine/urlfetch"
)

type App struct {
	AppId           string
	APIKey          string
	APISecret       string
	RedirectURI     string
	IgnoreSignature bool
}

func (s *App) ValidateShopifyHostName(host string) bool {
	host = strings.Trim(host, " ")

	re, _ := regexp.Compile(`^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*(myshopify\.com)$`)
	if re.MatchString(host) {
		return true
	}
	return false
}

func (s *App) AuthorizeURL(shop string, scopes string, state string) string {
	var u url.URL
	u.Scheme = "https"
	u.Host = shop
	u.Path = "/admin/oauth/authorize"
	q := u.Query()
	q.Set("client_id", s.APIKey)
	q.Set("scope", scopes)
	q.Set("redirect_uri", s.RedirectURI)
	q.Set("state", state)
	u.RawQuery = q.Encode()

	return u.String()
}

// Verify a message against a message HMAC
func (s *App) VerifyMessage(message, messageMAC string) bool {
	mac := hmac.New(sha256.New, []byte(s.APISecret))
	mac.Write([]byte(message))
	expectedMAC := mac.Sum(nil)

	// shopify HMAC is in hex so it needs to be decoded
	actualMac, _ := hex.DecodeString(messageMAC)

	return hmac.Equal(actualMac, expectedMAC)
}

// Verifying URL callback parameters.
func (s *App) VerifyAuthorizationURL(u *url.URL, nonce string) bool {
	q := u.Query()
	messageMAC := q.Get("hmac")
	state := q.Get("state")

	// Check nonce according to https://help.shopify.com/en/api/getting-started/authentication/oauth
	if nonce != "" && state != nonce {
		return false
	}
	// Check hostname according to to https://help.shopify.com/en/api/getting-started/authentication/oauth
	shop := q.Get("shop")
	if !s.ValidateShopifyHostName(shop) {
		return false
	}
	// Remove hmac and signature and leave the rest of the parameters alone.
	q.Del("hmac")
	q.Del("signature")

	message, _ := url.QueryUnescape(q.Encode())

	return s.VerifyMessage(message, messageMAC)
}

func (s *App) AdminSignatureOk(u *url.URL, nonce string) bool {
	return s.VerifyAuthorizationURL(u, nonce)
}

func (s *App) AppProxySignatureOk(u *url.URL) bool {
	if s.IgnoreSignature {
		return true
	}

	params := u.Query()
	signature := params["signature"]
	if signature == nil || len(signature) != 1 {
		return false
	}

	mac := hmac.New(sha256.New, []byte(s.APISecret))
	mac.Write([]byte(s.signatureString(u, false)))
	calculated := hex.EncodeToString(mac.Sum(nil))

	return 1 == subtle.ConstantTimeCompare([]byte(signature[0]), []byte(calculated))
}

func (s *App) signatureString(u *url.URL, prependSig bool) string {
	params := u.Query()

	keys := []string{}
	for k, _ := range params {
		if k != "signature" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	input := ""
	if prependSig {
		input = s.APISecret
	}
	for _, k := range keys {
		input = fmt.Sprintf("%s%s=%s", input, k, params[k][0])
	}
	return input
}

func (s *App) AccessToken(context appengine.Context, shop string, code string) (string, error) {
	url := fmt.Sprintf("https://%s/admin/oauth/access_token.json", shop)

	data := map[string]string{
		"client_id":     s.APIKey,
		"client_secret": s.APISecret,
		"code":          code,
	}

	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(data)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	client := urlfetch.Client(context)
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	token := map[string]string{}
	err = json.NewDecoder(response.Body).Decode(&token)

	if err != nil {
		return "", err
	}

	if _, ok := token["error"]; ok {
		return "", fmt.Errorf("%s", token["error"])
	}

	if _, ok := token["access_token"]; !ok {
		return "", fmt.Errorf("access_token not found in response")
	}

	return token["access_token"], nil
}
