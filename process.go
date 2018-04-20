// process
package main

import (
	"log"
	"strings"

	"github.com/AstroProfundis/gopsutil/cpu"
	"github.com/AstroProfundis/gopsutil/process"
)

type ProcessStat struct {
	Name     string                  `json:"name"`
	Pid      int32                   `json:"pid"`
	Exec     string                  `json:"exec"`
	Cmdline  string                  `json:"cmd"`
	Status   string                  `json:"status"`
	CpuTimes *cpu.TimesStat          `json:"cpu_times"`
	Memory   *process.MemoryInfoStat `json:"memory"`
	Rlimit   []RlimitUsage           `json:"resource_limit"`
}

type RlimitUsage struct {
	Resource string `json:"resource"`
	Soft     int64  `json:"soft"`
	Hard     int64  `json:"hard"`
	Used     uint64 `json:"used"`
}

func GetProcStats() []ProcessStat {
	stats := make([]ProcessStat, 0)
	ti_servers := []string{"pd-server", "tikv-server", "tidb-server"}
	for _, proc_name := range ti_servers {
		proc, err := getProcessesByName(proc_name)
		var stat ProcessStat
		if err != nil {
			log.Fatal(err)
		}
		stat.getProcessStat(proc)
		stats = append(stats, stat)
	}
	return stats
}

func getRlimitUsage(proc *process.Process) []RlimitUsage {
	resources := map[int32]string{
		// Resource limit constants are from:
		// /usr/include/x86_64-linux-gnu/bits/resource.h
		// from libc6-dev package in Ubuntu 16.10

		// Per-process CPU limit, in seconds.
		0: "cpu",

		// Largest file that can be created, in bytes.
		1: "fsize",

		// Maximum size of data segment, in bytes.
		2: "data",

		// Maximum size of stack segment, in bytes.
		3: "stack",

		// Largest core file that can be created, in bytes.
		4: "core",

		// Largest resident set size, in bytes.
		// This affects swapping; processes that are exceeding their
		// resident set size will be more likely to have physical memory
		// taken from them.
		5: "rss",

		// Number of processes.
		6: "nproc",

		// Number of open files.
		7: "nofile",

		// Locked-in-memory address space.
		8: "memlock",

		// Address space limit.
		9: "as",

		// Maximum number of file locks.
		10: "locks",

		// Maximum number of pending signals.
		11: "sigpending",

		// Maximum bytes in POSIX message queues.
		12: "msgqueue",

		// Maximum nice priority allowed to raise to.
		// Nice levels 19 .. -20 correspond to 0 .. 39
		// values of this resource limit.
		13: "nice",

		// Maximum realtime priority allowed for non-priviledged
		// processes.
		14: "rtprio",

		// Maximum CPU time in µs that a process scheduled under a real-time
		// scheduling policy may consume without making a blocking system
		// call before being forcibly descheduled.
		15: "rttime",
	}

	result := make([]RlimitUsage, 0)
	rlimit, _ := proc.RlimitUsage(true)
	for _, res := range rlimit {
		var usage RlimitUsage
		usage.Resource = resources[res.Resource]
		usage.Soft = int64(res.Soft)
		usage.Hard = int64(res.Hard)
		usage.Used = res.Used
		result = append(result, usage)
	}
	return result
}

func (proc_stat *ProcessStat) getProcessStat(proc *process.Process) {
	proc_stat.Pid = proc.Pid
	proc_stat.Name, _ = proc.Name()
	proc_stat.Exec, _ = proc.Exe()
	proc_stat.Cmdline, _ = proc.Cmdline()
	proc_stat.Status, _ = proc.Status()
	proc_stat.CpuTimes, _ = proc.Times()
	proc_stat.Memory, _ = proc.MemoryInfo()
	proc_stat.Rlimit = getRlimitUsage(proc)
}

func getProcessesByName(search_name string) (*process.Process, error) {
	proc_list, err := process.Processes()
	if err != nil {
		return nil, err
	}
	for _, proc := range proc_list {
		proc_name, err := proc.Name()
		if err != nil {
			return nil, err
		}
		// TODO: return multiple processes that match the search
		if strings.Contains(proc_name, search_name) {
			return proc, err
		}
	}
	return nil, err
}
