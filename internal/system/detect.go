package system

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Info struct {
	Hostname     string
	Kernel       string
	Distro       string
	Username     string
	Architecture string
	Hardware     string
}

func DetectSystem() *Info {
	info := &Info{
		Hostname:     getHostname(),
		Kernel:       getKernel(),
		Distro:       getDistro(),
		Username:     getUsername(),
		Architecture: getArchitecture(),
		Hardware:     getHardware(),
	}
	return info
}

func getHostname() string {
	if hostname, err := os.Hostname(); err == nil {
		return hostname
	}
	return "localhost"
}

func getKernel() string {
	if output, err := exec.Command("uname", "-r").Output(); err == nil {
		return strings.TrimSpace(string(output))
	}
	return "5.15.0-generic"
}

func getDistro() string {
	if data, err := os.ReadFile("/etc/os-release"); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "PRETTY_NAME=") {
				name := strings.TrimPrefix(line, "PRETTY_NAME=")
				name = strings.Trim(name, "\"")
				return name
			}
		}
	}

	if _, err := os.ReadFile("/etc/arch-release"); err == nil {
		return "Arch Linux"
	}

	if data, err := os.ReadFile("/etc/fedora-release"); err == nil {
		return strings.TrimSpace(string(data))
	}

	if data, err := os.ReadFile("/etc/redhat-release"); err == nil {
		return strings.TrimSpace(string(data))
	}

	return "Linux"
}

func getUsername() string {
	if user := os.Getenv("USER"); user != "" {
		return user
	}
	if user := os.Getenv("USERNAME"); user != "" {
		return user
	}
	return "user"
}

func getArchitecture() string {
	if output, err := exec.Command("uname", "-m").Output(); err == nil {
		return strings.TrimSpace(string(output))
	}
	return "x86_64"
}

func getHardware() string {
	if data, err := os.ReadFile("/sys/class/dmi/id/product_name"); err == nil {
		product := strings.TrimSpace(string(data))
		if vendor, err := os.ReadFile("/sys/class/dmi/id/sys_vendor"); err == nil {
			vendorStr := strings.TrimSpace(string(vendor))
			if product != "" && vendorStr != "" {
				return fmt.Sprintf("%s %s", vendorStr, product)
			}
		}
		if product != "" {
			return product
		}
	}

	if output, err := exec.Command("lscpu").Output(); err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "Model name:") {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					return strings.TrimSpace(parts[1])
				}
			}
		}
	}

	return "Generic PC"
}
