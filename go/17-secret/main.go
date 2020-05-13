package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/ramin0/live/go/secret/secret"
)

func main() {
	var (
		flagEncodingKey = flag.String("key", "", "Enconding key.")
		flagKeysPath    = flag.String("path", "vault.enc", "Path to vault file.")
	)
	flag.Parse()

	v := secret.FileVault(*flagEncodingKey, *flagKeysPath)
	switch cmd := flag.Arg(0); cmd {
	case "set":
		keyName, keyValue := flag.Arg(1), flag.Arg(2)
		if keyName == "" || keyValue == "" {
			log.Fatalf("Missing key name or value.")
		}
		if err := v.Set(keyName, keyValue); err != nil {
			log.Fatalf("Failed to set %q to %q: %v", keyName, keyValue, err)
		}
		fmt.Println("Value set!")
	case "get":
		keyName := flag.Arg(1)
		if keyName == "" {
			log.Fatalf("Missing key name.")
		}
		value, err := v.Get(keyName)
		if err != nil {
			// if err == secret.ErrKeyNotFound {
			// 	log.Fatalf("Key not found: %s", keyName)
			// }
			log.Fatalf("Failed to get %q: %v", keyName, err)
		}
		fmt.Println(value)
	case "list":
		keyNames, err := v.List()
		if err != nil {
			log.Fatalf("Failed to list: %v", err)
		}
		for _, keyName := range keyNames {
			fmt.Println(keyName)
		}
	case "delete":
		keyName := flag.Arg(1)
		if keyName == "" {
			log.Fatalf("Missing key name.")
		}
		if err := v.Delete(keyName); err != nil {
			log.Fatalf("Failed to delete %q: %v", keyName, err)
		}
		fmt.Println("Key deleted!")
	default:
		log.Fatalf("Unknown command %q, please use set or get.", cmd)
	}
}
