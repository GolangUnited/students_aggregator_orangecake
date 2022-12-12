package main

import "fmt"

func main() {
	cfg, err := NewServiceConfig()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(cfg)
}
