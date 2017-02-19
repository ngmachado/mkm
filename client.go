//Package mkm provides a client implementation to the magiccardmarket.eu api
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

//Method type represent http verbs allowed
type Method int

const (
	Get Method = iota
	Post
	Put
	Delete
)

//String return a string representation of type Method
func (m Method) String() string {
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

//OutputFormat represent the response format from the services
type OutputFormat int

const (
	XML OutputFormat = iota
	JSON
)

//String return a string representation of type OutputFormat
func (f OutputFormat) String() string {
	if f == 1 {
		return "/output.json"
	}

	return ""
}

//Endpoint define which environment to make requests
type Endpoint int

const (
	Sandbox Endpoint = iota
	Prodution
)

//String return a string representation of type Endpoint
func (e Endpoint) String() string {
	if e == 0 {
		return "https://sandbox.mkmapi.eu/ws"
	}

	return "https://www.mkmapi.eu/ws"
}

//Version define the API version to use
type Version int

const (
	V1 Version = iota
	V2
)

//String return a string representation of type Version
func (e Version) String() string {
	if e == 0 {
		return "/v1.1"
	} else {
		return "/v2.0"
	}
}

//Client is the starting point to making request to services
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
		endpoint:  endpoint.String(),
		version:   version.String(),
		outputfmt: output.String(),
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

//Request method will create the request, validate with oath header and sended to service
//Return the response and error if exists
func (c *Client) Request(method Method, resource string, data []byte) ([]byte, error) {
	url := fmt.Sprint(c.endpoint, c.version, c.outputfmt, resource)
	req, err := http.NewRequest(method.String(), url, bytes.NewReader(data))
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
