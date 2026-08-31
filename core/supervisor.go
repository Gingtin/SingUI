package core

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os/exec"
	"sync"
	"time"
)

type SupervisorStatus struct {
	IsRunning bool      `json:"is_running"`
	PID       int       `json:"pid"`
	StartTime time.Time `json:"start_time"`
	LastError string    `json:"last_error"`
	Version   string    `json:"version"`
}

type Supervisor struct {
	binPath    string
	configPath string
	cmd        *exec.Cmd
	ctx        context.Context
	cancel     context.CancelFunc
	mu         sync.Mutex

	status SupervisorStatus

	logMu       sync.RWMutex
	logBuffer   []string
	maxLogLines int
	subscribers map[chan string]bool
}

var Instance *Supervisor

func InitSupervisor(binPath, configPath string) *Supervisor {
	Instance = &Supervisor{
		binPath:     binPath,
		configPath:  configPath,
		maxLogLines: 500,
		logBuffer:   make([]string, 0, 500),
		subscribers: make(map[chan string]bool),
	}
	Instance.fetchVersion()
	return Instance
}

func (s *Supervisor) fetchVersion() {
	out, err := exec.Command(s.binPath, "version").Output()
	if err == nil {
		s.status.Version = string(out)
	} else {
		s.status.Version = "Sing-box (Executable not found or not in PATH)"
	}
}

func (s *Supervisor) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status.IsRunning {
		return nil
	}

	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.cmd = exec.CommandContext(s.ctx, s.binPath, "run", "-c", s.configPath)

	stdout, err := s.cmd.StdoutPipe()
	if err != nil {
		s.status.LastError = err.Error()
		return err
	}
	stderr, err := s.cmd.StderrPipe()
	if err != nil {
		s.status.LastError = err.Error()
		return err
	}

	if err := s.cmd.Start(); err != nil {
		s.status.LastError = err.Error()
		s.appendLog(fmt.Sprintf("[Supervisor] Failed to start sing-box: %v", err))
		return err
	}

	s.status.IsRunning = true
	s.status.PID = s.cmd.Process.Pid
	s.status.StartTime = time.Now()
	s.status.LastError = ""
	s.appendLog(fmt.Sprintf("[Supervisor] Sing-box started successfully with PID: %d", s.status.PID))

	go s.pipeLogs(stdout)
	go s.pipeLogs(stderr)

	go func() {
		err := s.cmd.Wait()
		s.mu.Lock()
		s.status.IsRunning = false
		s.status.PID = 0
		if err != nil {
			s.status.LastError = err.Error()
			s.appendLog(fmt.Sprintf("[Supervisor] Sing-box process exited: %v", err))
		} else {
			s.appendLog("[Supervisor] Sing-box process stopped normally.")
		}
		s.mu.Unlock()
	}()

	return nil
}

func (s *Supervisor) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.status.IsRunning || s.cmd == nil {
		return nil
	}

	s.appendLog("[Supervisor] Stopping Sing-box...")
	if s.cancel != nil {
		s.cancel()
	}

	s.status.IsRunning = false
	s.status.PID = 0
	return nil
}

func (s *Supervisor) Restart() error {
	_ = s.Stop()
	time.Sleep(500 * time.Millisecond)
	return s.Start()
}

func (s *Supervisor) GetStatus() SupervisorStatus {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.status
}

func (s *Supervisor) pipeLogs(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		s.appendLog(line)
	}
}

func (s *Supervisor) appendLog(line string) {
	formatted := fmt.Sprintf("[%s] %s", time.Now().Format("2006-01-02 15:04:05"), line)
	log.Println(line)

	s.logMu.Lock()
	if len(s.logBuffer) >= s.maxLogLines {
		s.logBuffer = s.logBuffer[1:]
	}
	s.logBuffer = append(s.logBuffer, formatted)

	// Broadcast to active WebSocket channels
	for ch := range s.subscribers {
		select {
		case ch <- formatted:
		default:
		}
	}
	s.logMu.Unlock()
}

func (s *Supervisor) GetRecentLogs() []string {
	s.logMu.RLock()
	defer s.logMu.RUnlock()
	result := make([]string, len(s.logBuffer))
	copy(result, s.logBuffer)
	return result
}

func (s *Supervisor) SubscribeLogs() chan string {
	s.logMu.Lock()
	defer s.logMu.Unlock()
	ch := make(chan string, 100)
	s.subscribers[ch] = true
	return ch
}

func (s *Supervisor) UnsubscribeLogs(ch chan string) {
	s.logMu.Lock()
	defer s.logMu.Unlock()
	if _, ok := s.subscribers[ch]; ok {
		delete(s.subscribers, ch)
		close(ch)
	}
}
