package panic

import (
	"fmt"
	"math/rand"
	"sysfail/internal/system"
)

func GetPanic(panicType string, sysInfo *system.Info) ([]string, bool) {
	switch panicType {
	case "panic":
		return generateKernelPanic(sysInfo), true
	case "oops":
		return generateKernelOops(sysInfo), true
	case "segfault":
		return generateSegfault(sysInfo), true
	default:
		return nil, false
	}
}

func generateKernelPanic(sysInfo *system.Info) []string {
	return []string{
		"kernel panic - not syncing: Fatal exception in interrupt",
		fmt.Sprintf("Pid: 0, comm: swapper Not tainted %s", sysInfo.Kernel),
		fmt.Sprintf("Hardware name: %s", sysInfo.Hardware),
		"Call Trace:",
		" <IRQ>",
		fmt.Sprintf(" [<ffffffff%08x>] panic+0x101/0x1f0", rand.Intn(0x99999999)),
		fmt.Sprintf(" [<ffffffff%08x>] do_IRQ+0x62/0x130", rand.Intn(0x99999999)),
		fmt.Sprintf(" [<ffffffff%08x>] common_interrupt+0x6a/0x6a", rand.Intn(0x99999999)),
		" <EOI>",
		"---[ end Kernel panic - not syncing: Fatal exception in interrupt",
		"",
		"System halted.",
	}
}

func generateKernelOops(sysInfo *system.Info) []string {
	return []string{
		"BUG: unable to handle kernel NULL pointer dereference at 0000000000000000",
		fmt.Sprintf("IP: [<ffffffff%08x>] __wake_up_common+0x2f/0x90", rand.Intn(0x99999999)),
		"PGD 0 ",
		"Oops: 0002 [#1] SMP ",
		"Modules linked in: nvidia(PO) drm nvidia_modeset(PO) nvidia_uvm(PO)",
		fmt.Sprintf("CPU: %d PID: %d Comm: systemd Tainted: P           O    %s",
			rand.Intn(8), rand.Intn(9999)+1000, sysInfo.Kernel),
		fmt.Sprintf("Hardware name: %s", sysInfo.Hardware),
		fmt.Sprintf("task: ffff%012x ti: ffff%012x task.ti: ffff%012x",
			rand.Int63n(0x999999999999), rand.Int63n(0x999999999999), rand.Int63n(0x999999999999)),
		fmt.Sprintf("RIP: 0010:[<ffffffff%08x>]  [<ffffffff%08x>] __wake_up_common+0x2f/0x90",
			rand.Intn(0x99999999), rand.Intn(0x99999999)),
		fmt.Sprintf("RSP: 0018:ffff%012x  EFLAGS: %08x",
			rand.Int63n(0x999999999999), rand.Intn(0x99999999)),
		fmt.Sprintf("RAX: %016x RBX: ffff%012x RCX: %016x",
			rand.Int63(), rand.Int63n(0x999999999999), rand.Int63()),
		fmt.Sprintf("RDX: %016x RSI: %016x RDI: ffff%012x",
			rand.Int63(), rand.Int63(), rand.Int63n(0x999999999999)),
		"Call Trace:",
		fmt.Sprintf(" [<ffffffff%08x>] __wake_up+0x48/0x60", rand.Intn(0x99999999)),
		fmt.Sprintf(" [<ffffffff%08x>] wake_up_process+0x15/0x20", rand.Intn(0x99999999)),
		fmt.Sprintf(" [<ffffffff%08x>] wake_up_worker+0x20/0x30", rand.Intn(0x99999999)),
		fmt.Sprintf("---[ end trace %016x ]---", rand.Int63()),
		"Kernel panic - not syncing: Fatal exception",
	}
}

func generateSegfault(sysInfo *system.Info) []string {
	return []string{
		fmt.Sprintf("segfault at 0 ip %016x sp %016x error 4 in libc-2.31.so[%x+25000]",
			rand.Int63(), rand.Int63(), rand.Intn(0x7f9b8c456)),
		fmt.Sprintf("Code: %02x %02x %02x %02x %02x %02x %02x %02x %02x %02x %02x %02x %02x %02x %02x %02x",
			rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256),
			rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256),
			rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256),
			rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256)),
		"general protection fault: 0000 [#1] SMP PTI",
		fmt.Sprintf("CPU: %d PID: %d Comm: chrome Not tainted %s",
			rand.Intn(8), rand.Intn(9999)+2000, sysInfo.Kernel),
		fmt.Sprintf("Hardware name: %s", sysInfo.Hardware),
		"RIP: 0010:memcpy+0x12/0x20",
		fmt.Sprintf("Code: %02x %02x %02x %02x %02x %02x %02x %02x %02x %02x %02x %02x %02x %02x %02x %02x",
			rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256),
			rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256),
			rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256),
			rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256)),
		fmt.Sprintf("RSP: 0018:ffffb2e4c0a7bd80 EFLAGS: %08x", rand.Intn(0x99999999)),
		fmt.Sprintf("RAX: %016x RBX: %016x RCX: %016x",
			rand.Int63(), rand.Int63(), rand.Int63()),
		fmt.Sprintf("RDX: %016x RSI: %016x RDI: %016x",
			rand.Int63(), rand.Int63(), rand.Int63()),
		"Segmentation fault (core dumped)",
	}
}
