// +build linux darwin

package main

import (
    "math/rand"
    "runtime"
    "syscall"
    "time"
)

func init() {
    rand.Seed(time.Now().UnixNano())

    var rlimit syscall.Rlimit
    syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlimit)
    rlimit.Cur = rlimit.Max
    syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rlimit)

    runtime.GOMAXPROCS(runtime.NumCPU()*2 + 1)
}
