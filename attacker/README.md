# attacker-service
Attacker strategy json rpc service.

# how to use
1. clone and build use go1.20
2. run `./build/bin/attacker-server` to start server

# how to call rpc
## attackclient
This is a example code to call `time_echo` rpc with attackclient.
```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tsinghua-cel/attacker-service/attackclient"
)

func main() {
	client, err := attackclient.Dial("http://localhost:10000")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	data := "Hello, World!"
	response, err := client.Echo(context.Background(), data)
	if err != nil {
		log.Fatalf("Failed to call time_echo: %v", err)
	}

	fmt.Printf("Response from time_echo: %s\n", response)
}
```

## use curl
```bash
curl -X POST -H "Content-Type: application/json" --data '{"jsonrpc":"2.0","method":"time_echo","params":["Hello, World!"],"id":1}' http://localhost:10000
```