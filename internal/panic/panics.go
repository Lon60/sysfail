package panic

var Panics = map[string][]string{
	"panic": {
		"Kernel panic - not syncing: Fatal sysfail error",
		"You can look but you cannot leave.",
	},
	"oops": {
		"Kernel panic - not syncing: Oops: Exception in kernel code",
		"CPU: 0 PID: 1234 Comm: troll_task Not tainted 6.2.0-arch #1 SMP PREEMPT",
		"Call Trace:",
		" dump_stack_lvl+0x5f/0x7a",
		" panic+0xe8/0x24f",
		" sys_fail_oops+0x3a/0x50",
		" RIP: 0010:sys_fail_oops+0x3a/0x50",
		" RSP: 0000:ffffc90008a6bdb8 EFLAGS: 00010002",
	},
}
