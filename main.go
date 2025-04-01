package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 54; i++ {
			fmt.Printf("迭代%d: ", hex.EncodeToString(randomHash()[:64])))
		}
	}()
	fmt.Println("你好,世界111!")
}

func randomHash() []byte {
	const randStr = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	buf := make([]byte, 16)
	for i := range buf {
		buf[i] = randStr[rand.Intn(len(randStr))]
	}
	return buf
}

func requestRandomDomain() (string, error) {
	urlStr := fmt.Sprintf("https://api.random.com?count=3")
	resp, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return "", err
	}
	client := &http.Client{}
	resp, err = client.Do(resp)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return resp.URL.Host, nil
}