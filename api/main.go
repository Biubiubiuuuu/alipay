package main

import (
	"fmt"

	"github.com/Biubiubiuuuu/alipay"
)

func main() {
	str, _ := alipay.GetAuth()
	fmt.Println(str)
}
