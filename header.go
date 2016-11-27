package mkm

import (
	"net/url"
	"sort"
	"strings"
)

type HeaderParams map[string]string

func NewHeaderParams(urlStr string, keys *Keys, nonce string, timestamp string) HeaderParams {
	return HeaderParams{
		"realm":                  urlStr,
		"oauth_consumer_key":     keys.ConsumerKey,
		"oauth_token":            keys.AccessToken,
		"oauth_nonce":            nonce,
		"oauth_timestamp":        timestamp,
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_version":          "1.0",
	}
}

func (h HeaderParams) OrderKeys() []string {
	var keys []string
	for k := range h {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (h HeaderParams) EscapeParams() string {
	escape := make(HeaderParams)
	for key, value := range h {
		if key != "realm" {
			escape[url.QueryEscape(key)] = url.QueryEscape(value)
		}
	}
	params := make([]string, len(h)-1, len(h)-1)
	okeys := escape.OrderKeys()
	for i := 0; i < len(okeys); i++ {
		key := okeys[i]
		value := escape[key]
		params[i] = key + "=" + value
	}
	return url.QueryEscape(strings.Join(params, "&"))
}
