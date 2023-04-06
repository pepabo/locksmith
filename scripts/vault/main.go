package main

import (
	"fmt"
	"time"
	"net/http"
	"encoding/json"
	"github.com/hashicorp/vault/api"
)

const (
	vaultAddr = "REPLACE ME"
	token = "REPLACE ME"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func main() {

	client, err := api.NewClient(&api.Config{Address: vaultAddr, HttpClient: httpClient})
	if err != nil {
			panic(err)
	}

	client.SetToken(token)
	data, err := client.Logical().Read("REPLACE ME")
	if err != nil {
			panic(err)
	}

	b, _ := json.Marshal(data.Data)
	fmt.Println(string(b))
}
