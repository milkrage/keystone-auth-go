# keystone-auth-go

## Features

1. **Token issue** - supports password authentication with the following scopes: Domain, Project, System, and Unscoped.

2. **Transparent token injection** - `keystoneauth.Transport` lets you use a plain `http.Client` without worrying about 
attaching the token yourself. The provided `http.RoundTripper` automatically adds the token to every request and 
refreshes it if a `401` response is received.

## Quick Start

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/milkrage/keystone-auth-go"
)

func main() {
	ctx := context.Background()

	credentials := keystoneauth.Credentials{
		Method: keystoneauth.MethodPassword,
		User: keystoneauth.User{
			Name:     "username",
			Password: "password",
			Domain: keystoneauth.Domain{
				Name: "Default",
			},
		},
		Scope: keystoneauth.Scope{
			Project: keystoneauth.Project{
				Name: "admin",
				Domain: keystoneauth.Domain{
					Name: "Default",
				},
			},
		},
	}

	authenticator, err := keystoneauth.NewAuthenticator(ctx, "http://127.0.0.1:5000", credentials)
	if err != nil {
		slog.Error("failed to create authenticator", "error", err)
		os.Exit(1)
	}

	// You can get the token where it is needed.
	token, err := authenticator.GetToken(ctx)
	if err != nil {
		slog.Error("failed to get token", "error", err)
		os.Exit(1)
	}
	fmt.Println("keystone token:", token)

	// You can also use http.Client with Transport
	// that automatically adds the token and refreshes it if necessary.
	client := http.Client{
		Transport: &keystoneauth.Transport{
			Authenticator: authenticator,
		},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://127.0.0.1:5000/v3/regions", nil)
	if err != nil {
		slog.Error("failed to create request", "error", err)
		os.Exit(1)
	}

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("failed to send request", "error", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	type Region struct {
		ID string `json:"id"`
	}

	type Response struct {
		Regions []Region `json:"regions"`
	}

	body := Response{}

	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		slog.Error("failed to decode response", "error", err)
		os.Exit(1)
	}

	for _, region := range body.Regions {
		fmt.Println(region.ID)
	}
}


```
