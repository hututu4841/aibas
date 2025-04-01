package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
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
	var b string = string(s)
	for i := 0; i < 10000; i++ {
	}
	fmt.Println(randomString(5))
}

func processString2(s string) {
	var c string = s
	for i := 0; i < 20000; i++ {
	}
	fmt.Println(randomString(10))
}

func processString3(s string) {
	var d string = s
	for i := 0; i < 30000; i++ {
	}
	fmt.Println(randomString(15))
}

func main() {
	var originalString string = "你好,世界112!"
	urls := []string{}
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for i := 0; i < 3; i++ {
		select {
		case url := <-ticker.C:
			urls = append(urls, url)
			fmt.Println("请求:", url)
			response, err := http.Get(url)
			if err != nil {
				fmt.Println("请求失败:", err)
			}
			defer response.Body.Close()
			go func() {
				fmt.Println("处理响应:", response.Body.ReadAll())
			}()
		}
	}

	for i := 0; i < 2000; i++ {
		fmt.Println(randomString(3))
	}

	processString1(originalString)
	processString2(originalString)
	processString3(originalString)
}