package mkm

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Method int

const (
	Get Method = iota
	Post
	Put
	Delete
)

func (m Method) Str() string {
	switch m {
	case 0:
		return "GET"
	case 1:
		return "POST"
	case 2:
		return "PUT"
	case 3:
		return "DELETE"
	default:
		return "GET"
	}
}

type OutputFormat int

const (
	XML OutputFormat = iota
	JSON
)

func (f OutputFormat) Str() string {
	if f == 1 {
		return "/output.json"
	} else {
		return ""
	}
}

type Endpoint int

const (
	Sandbox Endpoint = iota
	Prodution
)

func (e Endpoint) Str() string {
	if e == 0 {
		return "https://sandbox.mkmapi.eu/ws"
	} else {
		return "https://www.mkmapi.eu/ws"
	}
}

type Version int

const (
	V1 Version = iota
	V2
)

func (e Version) Str() string {
	if e == 0 {
		return "/v1.1"
	} else {
		return "/v2.0"
	}
}

type Client struct {
	client    *http.Client
	oauth     *OAuth
	endpoint  string
	version   string
	outputfmt string
}

//return a new client that with timeout after 10 seconds
func NewClient(keys *Keys, endpoint Endpoint, version Version, output OutputFormat) *Client {
	return &Client{
		client:    &http.Client{Timeout: time.Second * 10},
		oauth:     NewOAuth(keys, nil, nil),
		endpoint:  endpoint.Str(),
		version:   version.Str(),
		outputfmt: output.Str(),
	}
}

func (c *Client) do(req *http.Request) ([]byte, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		err = errors.New(strconv.Itoa(resp.StatusCode))
		return nil, err
	}
	return body, nil
}

func (c *Client) Request(method Method, resource string, data []byte) ([]byte, error) {
	url := fmt.Sprint(c.endpoint, c.version, c.outputfmt, resource)
	req, err := http.NewRequest(method.Str(), url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	hdr, err := c.oauth.AuthHeader(req.Method, url)
	if err != nil {
		return nil, err
	}
	req.Header = hdr
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	return c.do(req)
}
