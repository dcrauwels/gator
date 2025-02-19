package main

import (
	"fmt"
	"os"
)

func main() {
	homePath, _ := os.UserHomeDir()
	fmt.Println(homePath + "/.haha")
}
