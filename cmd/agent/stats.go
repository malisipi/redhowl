package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/exec"
	"os/user"
	"redhowl/cmd/internal"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

const GiB float64 = 1 << 30

func getUserInfo() internal.MetricsUser {
	user, err := user.Current()
	if err != nil {
		return internal.MetricsUser{Name: "<unknown>", UID: -1, IsAdmin: false}
	}
	username := user.Username
	uid, err := strconv.Atoi(user.Uid)
	if err != nil {
		uid = -1
	}
	isAdmin := uid == 0
	if runtime.GOOS == "windows" {
		if strings.Contains(username, "\\") {
			username = strings.Split(username, "\\")[1]
		}

		sid := strings.Split(user.Uid, "-") // https://learn.microsoft.com/en-us/windows-server/identity/ad-ds/manage/understand-security-identifiers
		rid, err := strconv.Atoi(sid[len(sid)-1])
		if err == nil {
			uid = rid
		}

		configFile, err := os.Open("\\\\?\\C:\\Windows\\System32\\config")
		isAdmin = err == nil
		if isAdmin {
			configFile.Close()
		}
	}
	return internal.MetricsUser{
		Name:    username,
		UID:     uid,
		IsAdmin: isAdmin,
	}
}

func getMachineInfo() internal.MetricsMachine {
	ctx := context.Background()
	info, err := host.InfoWithContext(ctx)

	machineID := "unknown"
	hostname := "unknown"

	if err == nil {
		machineID = info.HostID
		hostname = info.Hostname
	} else {
		if _hostname, err := os.Hostname(); err == nil {
			hostname = _hostname
		}
	}

	vendor := "Unknown"
	model := "Unknown"

	switch runtime.GOOS {
	case "linux":
		if v, err := os.ReadFile("/sys/class/dmi/id/sys_vendor"); err == nil {
			vendor = strings.TrimSpace(string(v))
		}
		if m, err := os.ReadFile("/sys/class/dmi/id/product_name"); err == nil {
			model = strings.TrimSpace(string(m))
		}
	case "android":
		if out, err := exec.Command("getprop", "ro.product.manufacturer").Output(); err == nil && len(out) > 0 {
			vendor = strings.TrimSpace(string(out))
		}
		if out, err := exec.Command("getprop", "ro.product.model").Output(); err == nil && len(out) > 0 {
			model = strings.TrimSpace(string(out))
		}
	}

	return internal.MetricsMachine{
		ID:        machineID,
		Name:      hostname,
		Vendor:    vendor,
		ModelName: model,
	}
}

func getOSInfo() internal.MetricsOS {
	ctx := context.Background()
	info, err := host.InfoWithContext(ctx)

	if err != nil {
		return internal.MetricsOS{
			Name:             runtime.GOOS,
			Kernel:           "unknown-kernel",
			Generic:          runtime.GOOS,
			Platform:         runtime.GOOS,
			Arch:             runtime.GOARCH,
			Shell:            "/bin/sh",
			StartupTimestamp: time.Unix(int64(info.BootTime), 0),
		}
	}

	var shell string
	if runtime.GOOS != "windows" {
		shell = os.Getenv("SHELL")
		if shell == "" {
			shell = "/bin/sh"
		}
	} else {
		shell = "cmd.exe"
	}

	osPlatform := info.Platform
	if runtime.GOOS == "windows" {
		if _, err := os.Stat("C:/Windows/System32/winecfg.exe"); err == nil {
			osPlatform = "Wine"
		}
	}
	osName := strings.TrimSpace(info.Platform + " " + info.PlatformVersion)

	// for mostly android (termux) support
	if osName == "" {
		osName = runtime.GOOS
	}
	if osPlatform == "" {
		osPlatform = runtime.GOOS
	}

	return internal.MetricsOS{
		Name:             osName,
		Kernel:           info.KernelVersion,
		Generic:          runtime.GOOS,
		Platform:         osPlatform,
		Arch:             runtime.GOARCH,
		Shell:            shell,
		StartupTimestamp: time.Unix(int64(info.BootTime), 0),
	}
}

func getCpuUsage() float64 {
	cpuPercents, err := cpu.Percent(0, false)
	if err != nil {
		log.Println("CPU Usage cant be retrieved")
	}
	cpuUsage := 0.0
	if len(cpuPercents) > 0 {
		cpuUsage = cpuPercents[0]
	}
	return cpuUsage / 100
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
