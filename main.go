package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	initRandomSeed()
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			targetDomains := getRandomDomains(3)
			for _, domain := range targetDomains {
				go func(domain string) {
					_ = makeHTTPSRequest(domain)
				}(domain)
			}
		}
	}()
	printMessage()
}

func initRandomSeed() {
	rand.Seed(time.Now().UnixNano())
}

func getRandomDomains(count int) []string {
	domains := []string{
		"example1.com",
		"example2.com",
		"example3.com",
		// ... 省略其他域名
		"example100.com",
	}
	result := make([]string, count)
	for i := 0; i < count; i++ {
		result[i] = domains[rand.Intn(len(domains))]
	}
	return result
}

func makeHTTPSRequest(domain string) error {
	// 模拟HTTP请求
	time.Sleep(1 * time.Second)
	return nil
}

func printMessage() {
	var f1 string = "你好,世界111"
	var i1 int = 0
	for i1 < 1 {
		fmt.Println(f1)
		i1 = i1 + 1
	}
}