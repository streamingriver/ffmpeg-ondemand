package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

func New(url []string, path string) *FFmpeg {
	cmd := exec.Command("ffmpeg", url...)
	return &FFmpeg{
		cmd,
		&sync.RWMutex{},
		false,
		time.Now(),
		path,
	}
}

type FFmpeg struct {
	cmd     *exec.Cmd
	mu      *sync.RWMutex
	running bool
	hit     time.Time
	path    string
}

func (f *FFmpeg) Start() error {
	if f.running {
		log.Printf("running.....")
		return nil
	}

	f.cmd.Stdout = os.Stdout
	f.cmd.Stderr = os.Stderr
	f.cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid:   true,
		Pdeathsig: syscall.SIGKILL,
	}
	f.cmd.Dir = f.path
	f.running = true
	f.Hit()
	return errors.New("FFmpeg process started")
}

func (f *FFmpeg) Stop() error {
	f.running = false
	return f.cmd.Process.Signal(syscall.SIGKILL)
}

func (f *FFmpeg) Wait() error {
	return f.cmd.Wait()
}

func (f *FFmpeg) IsRunning() bool {
	return f.running
}

func (f *FFmpeg) Hit() {
	f.hit = time.Now()
}

func (f *FFmpeg) HitExpired() bool {
	// return time.Now().Sub(f.hit).Minutes() > 0
	return time.Now().Sub(f.hit).Seconds() > 14 && f.running
}
