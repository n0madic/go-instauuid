# go-instauuid

[![Go Reference](https://pkg.go.dev/badge/github.com/n0madic/go-instauuid.svg)](https://pkg.go.dev/github.com/n0madic/go-instauuid)

`go-instauuid` is a Go package that provides a simple and efficient way to generate Instagram-like unique identifiers (UUIDs) for your applications.

Generates 8-byte UUID that consists of:

* 41 bits for time in milliseconds (gives us 41 years of IDs Instagrams custom epoch)
* 13 bits for additional information - Instagram used it to store the logical shard ID
* 10 bits that represent an auto-incrementing sequence, modulus 1024

## Installation

To install `go-instauuid`, use the following command:

```shell
go get github.com/n0madic/go-instauuid
```

## Usage

Here's an example of how to use `go-instauuid`:

```go
package main

import (
    "fmt"
    "github.com/n0madic/go-instauuid"
)

func main() {
    shardID := 1
    myEpoch := 1609459200
    generator := instauuid.NewGenerator(shardID, myEpoch)
	fmt.Println("ID:", generator.GenerateID())
	fmt.Println("Base64:", generator.GenerateBase64())
	fmt.Println("Hex:", generator.GenerateHex())
	fmt.Println("Buffer:", generator.GenerateBuffer())
	fmt.Println("BufferBE:", generator.GenerateBufferBE())
}
```

Output example:

```shell
ID: 14440563324670182400
Base64: yGcv1ZuABAE
Hex: c8672fd59b800402
Buffer: [3 4 128 155 213 47 103 200]
BufferBE: [200 103 47 213 155 128 4 4]
```

If the epoch value is not provided (set to zero), the default Instagram epoch value is `1387263000` (Wed Aug 24 2011 21:07:01 UTC).


## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request.

## License

This package is licensed under the MIT License. See the [LICENSE](LICENSE) file for more information.
