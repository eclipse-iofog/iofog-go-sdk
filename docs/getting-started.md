# Call the Controller from Go

You will add the module, open an authenticated client against Controller REST v3, and read Controller status. By the end you have a program that prints the Controller version string.

## What you need

- Go **1.26.9** or newer (`go.mod` sets this minimum).
- A running ioFog Controller that speaks REST **v3**.
- A user access token, or an email, password, and TOTP that the Controller accepts.

The base URL always includes the API prefix:

```text
http://127.0.0.1:51121/api/v3
```

Default Controller port is `51121` (`client.ControllerPort`).

## Step 1: Add the module

```bash
go get github.com/eclipse-iofog/iofog-go-sdk/v3
```

Import the neutral upstream path. Do not import `github.com/datasance/iofog-go-sdk/v3`. Datasance and eclipse-iofog publish the same module path.

## Step 2: Create a client with a token

Services and CI should pass a bearer token. `NewWithToken` does not read the terminal.

```go
package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/eclipse-iofog/iofog-go-sdk/v3/pkg/client"
)

func main() {
	base, err := url.Parse("http://127.0.0.1:51121/api/v3")
	if err != nil {
		log.Fatal(err)
	}
	clt, err := client.NewWithToken(client.Options{BaseURL: base}, "ACCESS_TOKEN")
	if err != nil {
		log.Fatal(err)
	}
	status, err := clt.GetStatus()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(status.Status, status.Versions.Controller)
}
```

`New` calls `GET /status` during construction and caches the Controller version. A Controller that is down still returns a client; `GetVersion` is empty until a later successful status call. `GetStatus` itself does not require a token. `BaseURL` must be non-nil.

## Step 3: Confirm the call

Run the program against your Controller. A healthy Controller prints a status string and a version such as `3.9.0`.

If the process hangs or prints a password prompt, you called `NewAndLogin` or `Login` with an empty `Totp`. Those paths read an OTP from the terminal. Use `NewWithToken` for non-interactive programs. Details are in [Use the Controller client](howto/controller-client.md).

## What you have

A `*client.Client` bound to `/api/v3`. Next:

- [Deploy an application from YAML](howto/deploy-yaml.md)
- [Client method reference](reference/client.md)
- [Package map](explanation/packages.md), if you are choosing between `client`, `apps`, and `microservices`
