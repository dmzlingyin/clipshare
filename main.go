package main

import (
	"clipshare/cmd"
	"fmt"
)

func main() {
	fmt.Println(cmd.Version())
	fmt.Println(cmd.Logo())
}
