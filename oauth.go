package mkm

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

//Noncer define capacity to create a nonce
type Noncer interface {
	Nonce() string
}

//Stamper define capacity to create a timestamp
type Stamper interface {
	Timestamp() string
}

//Nonce implements Noncer
type Nonce struct{}

//Nonce returns nonce base64 encoded random 32 byte string
func (n Nonce) Nonce() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

//Stamp implements Timestamp
type Stamp struct{}

//Timestamp returns the Unix epoch seconds string
func (s Stamp) Timestamp() string {
	return fmt.Sprintf("%v", time.Now().Unix())
}

//OAuth struct groups information to create a oath header
type OAuth struct {
	*Keys
	Noncer
	Stamper
}

//NewOAuth create a new OAuth
func NewOAuth(keys *Keys, noncer Noncer, stamper Stamper) *OAuth {
	if noncer == nil {
		noncer = new(Nonce)
	}

	if stamper == nil {
		stamper = new(Stamp)
	}

	return &OAuth{keys, noncer, stamper}
}

//Keys to access services
type Keys struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

//AuthHeader return a header parameters encodended and signed
func (o *OAuth) AuthHeader(method string, urlStr string) (map[string][]string, error) {
	queries := strings.Split(urlStr, "?")
	baseURL := queries[0]
	headerParams := NewHeaderParams(baseURL, o.Keys, o.Nonce(), o.Timestamp())
	baseString := method + "&" + url.QueryEscape(baseURL) + "&"
	if len(queries) > 1 {
		multiqry := strings.Split(queries[1], "&")
		for _, s := range multiqry {
			query := strings.Split(s, "=")
			if len(query) > 1 {
				key := query[0]
				value := query[1]
				headerParams[key] = value
			}
		}
	}

	param := headerParams.EscapeParams()
	baseString += param
	signatureKey := url.QueryEscape(o.ConsumerSecret) + "&" + url.QueryEscape(o.AccessTokenSecret)
	oAuthSignature := HMACSHA1(signatureKey, baseString)
	headerParams["oauth_signature"] = oAuthSignature
	headerParamStrings := make([]string, len(headerParams), len(headerParams))

	i := 0
	for key, value := range headerParams {
		headerParamStrings[i] = key + "=\"" + value + "\""
		i++
	}

	authHeader := "OAuth " + strings.Join(headerParamStrings, ", ")
	req, err := http.NewRequest(method, baseURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", authHeader)

	return req.Header, nil
}

//HMACSHA1 return base64 encoded from key and message provided
func HMACSHA1(key string, message string) string {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
