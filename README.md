# mkm
API Client to Magiccardmarket.eu

### How to Use
- Getting the code :

```
go get github.com/ngmachado/mkm
```

- import into your project
```
package main 

import (
...
"github.com/ngmachado/mkm"
)
...
```

# Working with mkm (example)

```
package main

import (
  "fmt"
  "log"

  "github.com/ngmachado/mkm"
)

func main() {
  //set your keys
  keys := &mkm.Keys{
    ConsumerKey:       "key1",
    ConsumerSecret:    "Key2",
    AccessToken:       "Key3",
    AccessTokenSecret: "Key4",
  }
  //create a client with some configuration. Client will timeout after 10 seconds
  //mkm.NewClient(keys, endpoint, version, output)
  //keys : your magiccardmarket keys
  //endpoint : (environment) Sandbox or Production
  //version : which version of the API to use. v1 or v2
  //output : response format. JSON or XML
  cli := mkm.NewClient(keys, mkm.Sandbox, mkm.V2, mkm.JSON)

  //make a request to magiccardmarket
  //cli.Request(method, resource, data)
  //method : HTTP method Get,Post,Put,Delete
  //resource : resource you want to work
  //data : if you want to send data to resource
  resp, err := cli.Request(mkm.Get, "/games", nil)

  if err != nil {
    //log.Fatal(err)
    log.Fatal("Request Error")
  }

  fmt.Printf("%+v", resp)

}

```

