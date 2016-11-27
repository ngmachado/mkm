package mkm

import (
	"encoding/base64"
	"strconv"
	"strings"
	"testing"
)

type TNonce struct {
	value string
}

func (t *TNonce) Nonce() string {
	return t.value
}

type TStamper struct {
	value string
}

func (t *TStamper) Timestamp() string {
	return t.value
}

func TestEncondedParam(t *testing.T) {
	nonce := "LyNoDV8MBnIrpBXtoTGefh9Vvcng5FQ1OCqB3kE5Ryk="
	timestamp := "1479036986"
	headerParams := NewHeaderParams("test.local", &Keys{}, nonce, timestamp)

	result := headerParams.EscapeParams()
	expect := "oauth_consumer_key%3D%26oauth_nonce%3DLyNoDV8MBnIrpBXtoTGefh9Vvcng5FQ1OCqB3kE5Ryk%253D%26oauth_signature_method%3DHMAC-SHA1%26oauth_timestamp%3D1479036986%26oauth_token%3D%26oauth_version%3D1.0"

	if expect != result {
		t.Error("expect ", expect, " got ", result)
	}

}

func TestNonce(t *testing.T) {
	n := &Nonce{}
	b64 := n.Nonce()
	str, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		t.Error("Nonce is not base64")
	}
	if len([]byte(str)) != 32 {
		t.Error("Nonce is not 32 bytes")
	}
}

func TestTimestamp(t *testing.T) {
	s := &Stamp{}
	unixTime := s.Timestamp()
	_, err := strconv.ParseInt(unixTime, 10, 64)
	if err != nil {
		t.Error("Unix Timestamp is not valid")
	}
}

func TestHMACSHA1(t *testing.T) {
	key := "ThisIs&MyKey"
	message := "Hello World"
	expect := "CSEPkzhx25VfVqU5YlTwhAR5Mk0="

	result := HMACSHA1(key, message)

	if expect != result {
		t.Error("expect ", expect, " got ", result)
	}
}

func TestOAuthSign(t *testing.T) {
	var result string
	expect := "oauth_signature=\"TRqCoignxlN/4FNNYClTMU/DChw=\""
	keys := &Keys{"ThisIs&MyKey1", "ThisIs&MyKey2", "ThisIs&MyKey3", "ThisIs&MyKey4"}
	nonce := &TNonce{"LyNoDV8MBnIrpBXtoTGefh9Vvcng5FQ1OCqB3kE5Ryk="}
	timestamp := &TStamper{"1479036986"}

	oauth := NewOAuth(keys, nonce, timestamp)
	header, err := oauth.AuthHeader("GET", "local.test")
	if err != nil {
		t.Error("Auth Header is not valid")
	}
	auth := header["Authorization"]
	tsx := strings.Split(auth[0], ",")
	for _, v := range tsx {
		if strings.Contains(v, "oauth_signature=") {
			result = strings.Trim(v, " ")
		}
	}
	if expect != result {
		t.Error("expect ", expect, " got ", result)
	}
}
