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

type Noncer interface {
	Nonce() string
}

type Stamper interface {
	Timestamp() string
}

type Nonce struct{}

// Returns nonce base64 encoded random 32 byte string
func (n Nonce) Nonce() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

type Stamp struct{}

//Timestamp returns the Unix epoch seconds string
func (s Stamp) Timestamp() string {
	return fmt.Sprintf("%v", time.Now().Unix())
}

type OAuth struct {
	*Keys
	Noncer
	Stamper
}

func NewOAuth(keys *Keys, noncer Noncer, stamper Stamper) *OAuth {

	if noncer == nil {
		noncer = new(Nonce)
	}

	if stamper == nil {
		stamper = new(Stamp)
	}

	return &OAuth{keys, noncer, stamper}
}

type Keys struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

func (o *OAuth) AuthHeader(method string, urlStr string) (map[string][]string, error) {
	queries := strings.Split(urlStr, "?")
	baseUrl := queries[0]
	//DI
	headerParams := NewHeaderParams(baseUrl, o.Keys, o.Nonce(), o.Timestamp())

	baseString := method + "&" + url.QueryEscape(baseUrl) + "&"
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
	req, err := http.NewRequest(method, baseUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", authHeader)

	return req.Header, nil
}

//Returns HMAC-SHA1 base64 encoded from key and message provided
func HMACSHA1(key string, message string) string {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
