package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	t1 := time.Now()
	go func() {
		for range time.Tick(5 * time.Minute) {
			d := rand.Intn(100)
			fmt.Println(requestDomain(d))
		}
	}()
	fmt.Println("你好,世界1112!")
}

func requestDomain(d int) string {
	domains := []string{
		"https://www.baidu.com",
		"https://www.qq.com",
		"https://www.taobao.com",
	}
	num := rand.Intn(len(domains))
	return domains[num]
}