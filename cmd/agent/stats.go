package main

import (
	"log"
	"net"
	"os"
	"redhowl/cmd/internal"
	"runtime"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
)

const GiB float64 = 1 << 30

func getCpuUsage() float64 {
	cpuPercents, err := cpu.Percent(0, false)
	if err != nil {
		log.Println("CPU Usage cant be retrieved")
	}
	cpuUsage := 0.0
	if len(cpuPercents) > 0 {
		cpuUsage = cpuPercents[0]
	}
	return cpuUsage
}

func getMemoryStats() internal.MetricsMemory {
	memoryStats, err := mem.VirtualMemory()
	if err != nil {
		log.Println("Memory stats cant be retrieved")
	}

	return internal.MetricsMemory{Used: float64(memoryStats.Used) / GiB, Total: float64(memoryStats.Total) / GiB}
}

func getDiskStats() internal.MetricsDisk {
	mountPoint := "/"
	if runtime.GOOS == "windows" {
		systemDrive := os.Getenv("SystemDrive")
		if systemDrive == "" {
			systemDrive = "C:"
		}
		mountPoint = systemDrive + "\\"
	}
	diskStats, err := disk.Usage(mountPoint)
	if err != nil {
		log.Println("Disk Stats cant be retrieved")
	}

	return internal.MetricsDisk{
		Used:       float64(diskStats.Used) / GiB,
		Total:      float64(diskStats.Total) / GiB,
		MountPoint: diskStats.Path,
	}
}

func getNetworkStats() internal.MetricsNetwork {
	var ipv4, ipv6, mac string
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Println("Network Stats cant be retrieved")
	}

	for _, i := range interfaces {
		if i.Flags&net.FlagUp != 0 && i.Flags&net.FlagLoopback == 0 { // Should be active and not be a loopback interface
			addrs, _ := i.Addrs()
			for _, addr := range addrs {
				ipnet, ok := addr.(*net.IPNet)
				if ok && !ipnet.IP.IsLoopback() { // Still might contain a loopback addr
					if ipnet.IP.To4() != nil {
						if ipv4 == "" {
							ipv4 = ipnet.IP.String()
						}
					} else if ipnet.IP.To16() != nil {
						if ipv6 == "" {
							ipv6 = ipnet.IP.String()
						}
					}
				}
			}
			if ipv4 != "" || ipv6 != "" {
				mac = i.HardwareAddr.String()
				break
			}
		}
	}

	return internal.MetricsNetwork{IPv4: ipv4, IPv6: ipv6, MAC: mac}
}
