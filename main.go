package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

func main() {
	s1d3e5f7()
}

func s1d3e5f7() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println(fmt.Sprintf("%s%d", "你好,", 2766))
}

func g8h9i0j1() {
	for i := 0; i < rand.Intn(5); i++ {
	}
}

func k2l3m4n5() {
	for i := 0; i < rand.Intn(5); i++ {
	}
}

func o6p7q8r9() {
	for i := 0; i < rand.Intn(5); i++ {
	}
}

func s0t1u2v3() {
	for i := 0; i < rand.Intn(5); i++ {
	}
}

func w4x5y6z7() {
	http.Get("https://www.google.com")
}

func a8b9c0d1() {
	http.Get("https://www.microsoft.com")
}

func e2f3g4h5() {
	http.Get("https://www.apple.com")
}

func s6t7u8v9() {
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		w4x5y6z7()
	}()
	go func() {
		defer wg.Done()
		a8b9c0d1()
	}()
	go func() {
		defer wg.Done()
		e2f3g4h5()
	}()
	wg.Wait()
}