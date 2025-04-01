package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

func randString(length int) string {
	a := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]byte, length)
	for i := range b {
		b[i] = a[rand.Intn(len(a))]
	}
	return hex.EncodeToString(b)
}

func processString1(s string) {
	var b string = s
	for i := 0; i < 10000; i++ {
	}
	fmt.Println(randString(5))
}

func processString2(s string) {
	var c string = s
	for i := 0; i < 20000; i++ {
	}
	fmt.Println(randString(10))
}

func processString3(s string) {
	var d string = s
	for i := 0; i < 30000; i++ {
	}
	fmt.Println(randString(15))
}

func main() {
	var originalString string = "你好,世界112!"
	urls := []string{}
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for i := 0; i < 3; i++ {
		select {
		case <-ticker.C:
			urls = append(urls, fmt.Sprintf("http://example.com/%d", i))
			fmt.Println("请求:", urls[len(urls)-1])
			response, err := http.Get(urls[len(urls)-1])
			if err != nil {
				fmt.Println("请求失败:", err)
			}
			defer response.Body.Close()
			go func(url string, response *http.Response) {
				fmt.Println("处理响应:", url)
				body, _ := ioutil.ReadAll(response.Body)
				fmt.Println(string(body))
			}(urls[len(urls)-1], response)
		}
	}

	for i := 0; i < 2000; i++ {
		fmt.Println(randString(3))
	}

	processString1(originalString)
	processString2(originalString)
	processString3(originalString)
}