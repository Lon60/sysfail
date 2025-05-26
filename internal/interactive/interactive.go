package interactive

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"sysfail/internal/system"
)

var (
	exitRequested = false
	currentDir    = "/"

	filesystem = map[string][]string{
		"/":          {"bin", "boot", "dev", "etc", "home", "lib", "lost+found", "proc", "root", "sbin", "sys", "tmp", "usr", "var"},
		"/home":      {"user"},
		"/home/user": {"Documents", "Downloads", "Desktop", ".bashrc", ".profile"},
		"/etc":       {"passwd", "shadow", "fstab", "hostname", "hosts", "sudoers"},
		"/var":       {"log", "tmp", "cache"},
		"/var/log":   {"syslog", "kern.log", "auth.log", "dmesg"},
		"/usr":       {"bin", "lib", "share", "local"},
		"/tmp":       {},
		"/proc":      {},
		"/sys":       {},
		"/dev":       {"sda", "sda1", "sda2", "null", "zero", "random"},
	}
)

func GracefulExit() {
	exitRequested = true
	fmt.Println("Connection to recovery console terminated.")
	os.Exit(0)
}

func Run(panicLines []string, sysInfo *system.Info) {
	displayPanicWithDelay(panicLines)

	fmt.Println()
	fmt.Println("\033[1;33m[EMERGENCY] Dropping to emergency shell...\033[0m")
	fmt.Println("\033[1;31mFilesystem errors detected. System in read-only mode.\033[0m")
	time.Sleep(3 * time.Second)

	scanner := bufio.NewScanner(os.Stdin)
	for !exitRequested {
		fmt.Printf("\033[1;31mroot\033[0m@\033[1;33m%s\033[0m:\033[1;34m%s\033[0m# ", sysInfo.Hostname, currentDir)

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		handleCommand(input, sysInfo)
	}
}

func displayPanicWithDelay(lines []string) {
	for i, line := range lines {
		if isErrorLine(line) {
			fmt.Printf("\033[1;31m%s\033[0m\n", line)
		} else {
			fmt.Println(line)
		}

		if i < 3 {
			time.Sleep(800 * time.Millisecond)
		} else if i < 8 {
			time.Sleep(400 * time.Millisecond)
		} else {
			time.Sleep(150 * time.Millisecond)
		}
	}
}

func handleCommand(cmd string, sysInfo *system.Info) {
	time.Sleep(time.Duration(rand.Intn(300)+50) * time.Millisecond)

	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return
	}

	command := parts[0]

	switch command {
	case "exit", "logout":
		fmt.Println("exit: cannot exit emergency shell")
		fmt.Println("Use Ctrl+Alt+Del to restart system")

	case "poweroff", "halt":
		fmt.Println("System halt prohibited in emergency mode")
		fmt.Println("Critical filesystem errors prevent shutdown")

	case "reboot", "shutdown":
		fmt.Println("Reboot blocked: filesystem consistency check required")
		fmt.Println("Run fsck manually before restart")

	case "cd":
		handleCD(parts)

	case "ls", "ll":
		handleLS(parts)

	case "pwd":
		fmt.Println(currentDir)

	case "cat":
		handleCat(parts, sysInfo)

	case "sudo":
		fmt.Println("sudo: /etc/sudoers: Input/output error")
		fmt.Println("sudo: unable to read sudoers file")

	case "ps", "top", "htop":
		fmt.Println("Error: /proc filesystem not mounted")
		fmt.Println("Cannot access process information")

	case "kill", "killall":
		fmt.Println("kill: cannot access process table")
		fmt.Println("procfs unavailable in emergency mode")

	case "dmesg":
		displayRandomErrors()

	case "mount":
		fmt.Println("/dev/sda1 on / type ext4 (ro,errors=remount-ro)")
		fmt.Println("\033[1;31mWARNING: Root filesystem mounted read-only\033[0m")
		fmt.Println("Emergency remount due to filesystem errors")

	case "umount":
		fmt.Println("umount: cannot unmount /: Device is busy")
		fmt.Println("umount: /: target is busy")

	case "df":
		fmt.Println("Filesystem corruption detected")
		fmt.Println("/dev/sda1: superblock read error")

	case "fsck", "e2fsck":
		handleFsck()

	case "help", "man":
		fmt.Println("\033[1;33mEMERGENCY SHELL - LIMITED FUNCTIONALITY\033[0m")
		fmt.Println("Available commands: ls, cd, pwd, cat, mount, dmesg, fsck")
		fmt.Println("\033[1;31mFilesystem in read-only mode - repair required\033[0m")

	case "clear":
		fmt.Println("clear: command not found in emergency shell")

	case "nano", "vi", "vim", "emacs":
		fmt.Println(fmt.Sprintf("%s: cannot access terminal properly", command))
		fmt.Println("Text editor unavailable in emergency mode")

	case "chmod", "chown":
		fmt.Println(fmt.Sprintf("%s: Read-only file system", command))

	case "mkdir", "rmdir", "rm", "cp", "mv":
		fmt.Println(fmt.Sprintf("%s: Read-only file system", command))

	case "uname":
		if len(parts) > 1 && parts[1] == "-a" {
			fmt.Printf("Linux %s %s %s %s GNU/Linux\n",
				sysInfo.Hostname, sysInfo.Kernel, sysInfo.Architecture, sysInfo.Architecture)
		} else {
			fmt.Println("Linux")
		}

	case "whoami":
		fmt.Println("root")

	case "id":
		fmt.Println("uid=0(root) gid=0(root) groups=0(root)")

	default:
		if strings.Contains(cmd, "/") {
			fmt.Printf("bash: %s: No such file or directory\n", cmd)
		} else {
			fmt.Printf("bash: %s: command not found\n", cmd)
		}
	}
}

func handleCD(parts []string) {
	if len(parts) == 1 {
		currentDir = "/"
		return
	}

	target := parts[1]

	if target == ".." {
		if currentDir == "/" {
			return
		}
		currentDir = filepath.Dir(currentDir)
		if currentDir == "." {
			currentDir = "/"
		}
		return
	}

	if target == "." || target == "./" {
		return
	}

	var newPath string
	if strings.HasPrefix(target, "/") {
		newPath = target
	} else {
		if currentDir == "/" {
			newPath = "/" + target
		} else {
			newPath = currentDir + "/" + target
		}
	}

	newPath = filepath.Clean(newPath)

	if _, exists := filesystem[newPath]; exists {
		if newPath == "/proc" || newPath == "/sys" {
			fmt.Printf("cd: %s: Permission denied (emergency mode)\n", newPath)
			return
		}
		currentDir = newPath
	} else {
		fmt.Printf("cd: %s: No such file or directory\n", target)
	}
}

func handleLS(parts []string) {
	var targetDir string

	if len(parts) == 1 {
		targetDir = currentDir
	} else {
		target := parts[1]
		if strings.HasPrefix(target, "/") {
			targetDir = target
		} else {
			if currentDir == "/" {
				targetDir = "/" + target
			} else {
				targetDir = currentDir + "/" + target
			}
		}
		targetDir = filepath.Clean(targetDir)
	}

	if targetDir == "/proc" || targetDir == "/sys" {
		fmt.Printf("ls: cannot access '%s': Permission denied\n", targetDir)
		return
	}

	if targetDir == "/lost+found" {
		fmt.Printf("ls: cannot access '%s': Input/output error\n", targetDir)
		return
	}

	if contents, exists := filesystem[targetDir]; exists {
		if len(contents) == 0 {
			return
		}

		if rand.Float32() < 0.2 && targetDir != "/" {
			fmt.Printf("ls: cannot access '%s': Input/output error\n", targetDir)
			return
		}

		for _, item := range contents {
			if rand.Float32() < 0.1 {
				fmt.Printf("?????????? ? ?    ?       ?            ? %s\n", item)
			} else {
				if isDirectory(targetDir + "/" + item) {
					fmt.Printf("drwxr-xr-x 2 root root  4096 Jan  1 12:00 %s\n", item)
				} else {
					fmt.Printf("-rw-r--r-- 1 root root  1024 Jan  1 12:00 %s\n", item)
				}
			}
		}
	} else {
		fmt.Printf("ls: cannot access '%s': No such file or directory\n", targetDir)
	}
}

func handleCat(parts []string, sysInfo *system.Info) {
	if len(parts) < 2 {
		fmt.Println("cat: missing file operand")
		return
	}

	filename := parts[1]

	switch filename {
	case "/etc/passwd":
		fmt.Println("cat: /etc/passwd: Input/output error")
	case "/etc/fstab":
		fmt.Println("# /etc/fstab: static file system information")
		fmt.Println("UUID=12345678-1234-1234-1234-123456789012 / ext4 errors=remount-ro 0 1")
		fmt.Println("# WARNING: Filesystem errors detected")
	case "/etc/hostname":
		fmt.Println(sysInfo.Hostname)
	case "/etc/os-release":
		fmt.Printf("NAME=\"%s\"\n", sysInfo.Distro)
		fmt.Printf("VERSION_ID=\"%s\"\n", strings.Split(sysInfo.Kernel, "-")[0])
		fmt.Printf("PRETTY_NAME=\"%s\"\n", sysInfo.Distro)
	case "/proc/version":
		fmt.Printf("Linux version %s (root@%s) (gcc version 11.2.0) #1 SMP\n",
			sysInfo.Kernel, sysInfo.Hostname)
	case "/var/log/syslog", "/var/log/kern.log":
		displayRandomErrors()
	default:
		if strings.Contains(filename, "/") {
			fmt.Printf("cat: %s: No such file or directory\n", filename)
		} else {
			fullPath := filepath.Join(currentDir, filename)
			fmt.Printf("cat: %s: No such file or directory\n", fullPath)
		}
	}
}

func handleFsck() {
	fmt.Println("fsck from util-linux 2.36.1")
	time.Sleep(1 * time.Second)
	fmt.Println("e2fsck 1.45.5 (07-Jan-2020)")
	time.Sleep(500 * time.Millisecond)
	fmt.Println("/dev/sda1: recovering journal")
	time.Sleep(2 * time.Second)
	fmt.Println("/dev/sda1: clean, 125648/1310720 files, 1547329/5242880 blocks")
	time.Sleep(1 * time.Second)
	fmt.Println("Filesystem check completed")
	fmt.Println("\033[1;33mNOTE: Remount required for write access\033[0m")
}

func displayRandomErrors() {
	errors := []string{
		fmt.Sprintf("[  %d.%06d] I/O error, dev sda, sector %d op 0x1:(WRITE) flags 0x20800",
			rand.Intn(99999), rand.Intn(999999), rand.Intn(999999)),
		fmt.Sprintf("[  %d.%06d] Buffer I/O error on dev sda1, logical block %d",
			rand.Intn(99999), rand.Intn(999999), rand.Intn(999999)),
		fmt.Sprintf("[  %d.%06d] critical temperature reached (%d C), shutting down",
			rand.Intn(99999), rand.Intn(999999), rand.Intn(20)+90),
		fmt.Sprintf("[  %d.%06d] Memory corruption detected at address 0x%016x",
			rand.Intn(99999), rand.Intn(999999), rand.Int63()),
		fmt.Sprintf("[  %d.%06d] WARNING: CPU: %d PID: %d at kernel/panic.c:%d",
			rand.Intn(99999), rand.Intn(999999), rand.Intn(8), rand.Intn(9999), rand.Intn(999)+100),
		fmt.Sprintf("[  %d.%06d] Out of memory: Kill process %d (%s) score %d or sacrifice child",
			rand.Intn(99999), rand.Intn(999999), rand.Intn(9999)+1000,
			[]string{"chrome", "firefox", "systemd", "kworker"}[rand.Intn(4)], rand.Intn(1000)),
		fmt.Sprintf("[  %d.%06d] EXT4-fs error (device sda1): ext4_journal_check_start:%d: Detected aborted journal",
			rand.Intn(99999), rand.Intn(999999), rand.Intn(100)+50),
		fmt.Sprintf("[  %d.%06d] ACPI: Critical thermal trip point reached",
			rand.Intn(99999), rand.Intn(999999)),
	}

	fmt.Println("\033[1;31m=== SYSTEM ERROR LOG ===\033[0m")
	for i := 0; i < rand.Intn(4)+3; i++ {
		fmt.Println(errors[rand.Intn(len(errors))])
		time.Sleep(300 * time.Millisecond)
	}
	fmt.Println("\033[1;31m=== END ERROR LOG ===\033[0m")
}

func isDirectory(path string) bool {
	_, exists := filesystem[path]
	return exists
}

func isErrorLine(line string) bool {
	if len(line) == 0 {
		return false
	}
	lower := strings.ToLower(line)
	return strings.Contains(lower, "panic") ||
		strings.Contains(lower, "oops") ||
		strings.Contains(lower, "error") ||
		strings.Contains(lower, "segfault") ||
		strings.Contains(lower, "fatal")
}
