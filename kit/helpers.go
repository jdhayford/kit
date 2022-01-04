package kit

import (
	"fmt"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func clearLastLine() {
	fmt.Println("\033[1A\x1b[K\033[1A")
}

func clearLastNLines(n int) {
	nString := fmt.Sprint(n)
	fmt.Println("\033[" + nString + "A\x1b[K\033[" + nString + "A")
}
