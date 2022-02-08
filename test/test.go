package main

import (
	"fmt"
	"os"
	"path"
)

func getKitDir() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Something went wrong retrieving the user home directory")
	}

	kitDirPath := path.Join(home, "/.kit")
	fmt.Println(kitDirPath)
}

func main() {
	getKitDir()
}
