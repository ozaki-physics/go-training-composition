package main

import (
	"fmt"
)

type IPAddr [4]byte

func (i IPAddr) String() string {
	return fmt.Sprintf("%v.%v.%v.%v", i[0], i[1], i[2], i[3])
}

func main() {
	hosts := map[string]IPAddr{
		"loopback":  {127, 0, 0, 1},
		"googleDNS": {8, 8, 8, 8},
	}

	for name, ip := range hosts {
		// 本来なら loopback: [127 0 0 1] とただの配列として出力される
		fmt.Printf("%v: %v\n", name, ip)
		// String() を IPAddr 型に実装したことで
		// loopback: 127.0.0.1 と出力されるようになった
	}
}
