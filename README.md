# Mailer

Mailer is a simple email sending microservice and library built with Go.

## Installation

To install and start Mailer, run the following commands:

```bash
git clone https://github.com/wizzldev/mailer.git

# Change to the Mailer directory and start the service
cd mailer
make run
```

## Client

To use the Mailer client, run the following command:

```bash
go get github.com/wizzldev/mailer/client
```

## Usage

To send an email using the Mailer client, run the following command:

```go
package main

import (
	"fmt"
	mailer "github.com/wizzldev/mailer/client"
)

func main() {
	client := mailer.NewClient("localhost:8080")
	err := client.SendText("recipient@example.com", "Hello", "Hello, world!")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email sent!")
}
```
