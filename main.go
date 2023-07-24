package main

import (
	"fmt"
	"net"
	"net/http"
	_ "os"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

func main() {
	http.HandleFunc("/time", timeHandler)
	http.HandleFunc("/disk", diskHandler)
	http.HandleFunc("/version", versionHandler)
	http.HandleFunc("/network", networkHandler)
	http.HandleFunc("/ram", ramHandler)

	http.ListenAndServe(":8080", nil)
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Current time: %s", time.Now().Format(time.RFC1123))
}

func diskHandler(w http.ResponseWriter, r *http.Request) {
	diskStat, _ := disk.Usage("/")
	fmt.Fprintf(w, "Disk Usage: Total: %v, Free: %v, Used: %v", diskStat.Total, diskStat.Free, diskStat.UsedPercent)
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Go version: %s\nOS: %s\n", runtime.Version(), runtime.GOOS)
}

func networkHandler(w http.ResponseWriter, r *http.Request) {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)

	fmt.Fprintf(w, "Local IP : %s\n", localAddr.IP)

	gateways, _ := net.InterfaceAddrs()
	for _, address := range gateways {
		fmt.Fprintf(w, "Gateway: %s\n", address.String())
	}
}

func ramHandler(w http.ResponseWriter, r *http.Request) {
	v, _ := mem.VirtualMemory()
	fmt.Fprintf(w, "RAM: Total: %v, Free: %v, UsedPercent: %f%%\n", v.Total, v.Free, v.UsedPercent)
}
