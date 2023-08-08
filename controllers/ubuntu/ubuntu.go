// Package ubuntu предоставляет функции для получения информации о системе Ubuntu.
package ubuntu

import (
	"fmt"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"net"
	"runtime"
	"time"
)

// Ubuntu представляет структуру для работы с функциональностью Ubuntu.
type Ubuntu struct {
}

// Name возвращает имя операционной системы.
func (o *Ubuntu) Name() string {
	return fmt.Sprintf("Ubuntu")
}

// Time возвращает текущее системное время.
func (o *Ubuntu) Time() string {
	fmt.Println("time using Ubuntu controllers...")
	return fmt.Sprintf("Current time: %s", time.Now().Format(time.RFC1123))
}

// Disk возвращает информацию о использовании дискового пространства.
func (o *Ubuntu) Disk() string {
	fmt.Println("disk using Ubuntu controllers...")
	diskStat, _ := disk.Usage("/")
	return fmt.Sprintf("Disk Usage: Total: %v, Free: %v, Used: %v", diskStat.Total, diskStat.Free, diskStat.UsedPercent)
}

// Version возвращает версию Go и имя операционной системы.
func (o *Ubuntu) Version() string {
	fmt.Println("version using Ubuntu controllers...")
	return fmt.Sprintf(fmt.Sprintf("Go version: %s\nOS: %s\n", runtime.Version(), runtime.GOOS))
}

// Network возвращает информацию о сетевых настройках системы.
func (o *Ubuntu) Network() string {
	fmt.Println("network using Ubuntu controllers...")
	conn, _ := net.Dial("udp", "8.8.8.8:80")

	defer func() {
		if err := conn.Close(); err != nil {
			fmt.Printf("Failed to close connection: %v", err)
		}
	}()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	gateways, _ := net.InterfaceAddrs()

	var gateway string
	for _, address := range gateways {
		gateway = address.String()
	}

	return fmt.Sprintf("Local IP : %s\nGateway: %s", localAddr.IP, gateway)
}

// Ram возвращает информацию о использовании оперативной памяти.
func (o *Ubuntu) Ram() string {
	fmt.Println("ram using Ubuntu controllers...")
	v, _ := mem.VirtualMemory()
	return fmt.Sprintf("RAM: Total: %v, Free: %v, UsedPercent: %f%%\n", v.Total, v.Free, v.UsedPercent)
}
