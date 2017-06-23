package bootstrap

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

func Init() {
	WritePid()
	SetGOMAXPROCS()
	LoadTimeZone()
}

func WritePid() {
	tmpDir := os.TempDir()
	exe := path.Base(os.Args[0])
	tmpFileName := tmpDir + "/" + exe + ".pid"
	fh, err := os.OpenFile(tmpFileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer fh.Close()
	b := make([]byte, 1024)
	n, _ := fh.Read(b)
	if n > 0 {
		pid := string(b[:n])
		fmt.Println("Last pid " + pid)

		lastPid, err := strconv.Atoi(pid)
		if err == nil {
			p, err := os.FindProcess(lastPid)
			if err == nil {
				err := p.Signal(syscall.Signal(0))
				if err == nil {
					fmt.Println("The pid already using: " + pid)
					os.Exit(1)
				}
			}

		} else {
			fmt.Println(err)
		}
	}
	err = syscall.Flock(int(fh.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		fmt.Println("App already running: ", err)
		os.Exit(1)
	}

	err = fh.Truncate(0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fh.Seek(0, 0)

	pid := os.Getpid()
	_, err = fh.WriteString(strconv.Itoa(pid))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Wrote pid %v to file %v\n", pid, tmpFileName)
}

func SetGOMAXPROCS() {
	cpu := runtime.NumCPU()
	log.Println("CPU:", cpu)
	runtime.GOMAXPROCS(cpu)
}

func LoadTimeZone() {
	l, err := time.LoadLocation("Etc/GMT-3")
	if err != nil {
		panic(err)
	}

	time.Local = l
	log.Println("Time:", time.Now())
}
