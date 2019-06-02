package main

import (
	"fmt"
	"log"
)

var (
	keysMap = make(map[int64]int)
	cards   = []int{
		1, 4, 4, 4, // 2,3,4,5
		4, 4, 4, 4, // 6,7,8,9
		4, 4, 4, 4, 3, // 10, j,q,k,ace
	}
)

func dumpKeysTag() {
	for k, v := range keyTags {
		fmt.Printf("agariTable[0x%x]=0x%x\n", k, v)
	}
}

func main() {

	log.Println("keysMap len:", len(keysMap))

	genAllTag()
	dumpKeysTag()
}
