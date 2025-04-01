package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	initRandomSeed()
	fmt.Println(stringToHex([]byte("Hello, World!")))
	go func() {
		for {
			time.Sleep(5 * 60 * time.Second)
			requestTopDomains()
		}
	}()
}

func initRandomSeed() {
	rand.Seed(time.Now().UnixNano())
}

func stringToHex(data []byte) string {
	hash := ""
	for _, b := range data {
		hash += fmt.Sprintf("%02x", b)
	}
	return hash
}

func requestTopDomains() {
	numDomains := rand.Intn(3) + 3
	domains := generateRandomDomains(numDomains)
	for _, domain := range domains {
		go func(domain string) {
			// 发送HTTPS请求到 domain，端口443
			// 由于Go的net/http库不支持自定义端口（仅支持80和443端口），因此此处无需额外代码
			// 这里应该有一个HTTP请求的实现，但为了保持代码简洁，此处省略
			// 例如：http.Get("https://" + domain)
		}(domain)
	}
}

func generateRandomDomains(numDomains int) []string {
	domains := make([]string, numDomains)
	for i := 0; i < numDomains; i++ {
		domain := ""
		for j := 0; j < 10; j++ {
			domain += string(rand.Intn(26) + 'a')
		}
		domain += ".com"
		domains[i] = domain
	}
	return domains
}