package mkm

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
	v1 Endpoint = iota
	v2
)

func (e Endpoint) Str() string {
	if e == 0 {
		return "/v1.1"
	} else {
		return "/v2.0"
	}
}
