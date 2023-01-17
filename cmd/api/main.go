package main

import (
	"fmt"

	rsaenc "github.com/diyliv/storage/pkg/rsa"
)

func main() {
	keys, err := rsaenc.GenerateKeys()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%x\n", keys.D)
}
