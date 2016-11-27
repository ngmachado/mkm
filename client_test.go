package mkm

import (
	"fmt"
	"testing"
)

func TestClient(t *testing.T) {
	keys := &Keys{
		ConsumerKey:       "key1",
		ConsumerSecret:    "Key2",
		AccessToken:       "Key3",
		AccessTokenSecret: "Key4",
	}
	cli := NewClient(keys, Prodution, v1, JSON)

	resp, err := cli.request(Get, "/games", nil)
	if err != nil {
		t.Error("Request Error")
	}

	fmt.Printf("%+v", resp)
}
