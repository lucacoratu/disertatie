package detection

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/lucacoratu/disertatie/agent/logging"
)

// Marker printed by the python classifier scripts once the model is loaded
const classifierReadyMarker = "READY"

// How long to wait for the python process to load the model
const classifierStartupTimeout = 60 * time.Second

// How long to wait for a single classification response
const classifierRequestTimeout = 5 * time.Second

/*
 * ClassifierProcess keeps a long-running python classifier script alive and
 * communicates with it through pipes (stdin / stdout).
 * Protocol: one line of comma separated features is written to stdin,
 * one line containing the classification is read back from stdout.
 * Requests are serialized with a mutex because the pipe is a single stream.
 */
type ClassifierProcess struct {
	logger      logging.ILogger
	interpreter string
	scriptPath  string

	mu     sync.Mutex
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	lines  chan string // lines read from the python stdout
	exited chan struct{}
}

func NewClassifierProcess(logger logging.ILogger, interpreter string, scriptPath string) *ClassifierProcess {
	return &ClassifierProcess{logger: logger, interpreter: interpreter, scriptPath: scriptPath}
}

// Start launches the python process and waits for it to signal that the model is loaded
func (cp *ClassifierProcess) Start() error {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	return cp.startLocked()
}

func (cp *ClassifierProcess) startLocked() error {
	// -u => unbuffered stdout so responses are delivered immediately
	cmd := exec.Command(cp.interpreter, "-u", cp.scriptPath)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("could not create stdin pipe for classifier: %w", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("could not create stdout pipe for classifier: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("could not create stderr pipe for classifier: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("could not start classifier process: %w", err)
	}

	lines := make(chan string)
	exited := make(chan struct{})

	// Forward stdout lines to the channel
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			lines <- scanner.Text()
		}
		close(lines)
	}()

	// Log anything the script writes on stderr
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			cp.logger.Error("Classifier stderr:", scanner.Text())
		}
	}()

	// Reap the process when it exits
	go func() {
		err := cmd.Wait()
		if err != nil {
			cp.logger.Error("Classifier process exited", err.Error())
		} else {
			cp.logger.Warning("Classifier process exited")
		}
		close(exited)
	}()

	cp.cmd = cmd
	cp.stdin = stdin
	cp.lines = lines
	cp.exited = exited

	// Wait for the model to be loaded
	select {
	case line, ok := <-lines:
		if !ok {
			cp.killLocked()
			return errors.New("classifier process exited before becoming ready")
		}
		if strings.TrimSpace(line) != classifierReadyMarker {
			cp.killLocked()
			return fmt.Errorf("unexpected startup message from classifier: %q", line)
		}
	case <-time.After(classifierStartupTimeout):
		cp.killLocked()
		return errors.New("timed out waiting for classifier process to start")
	}

	cp.logger.Info("Classifier process started", cp.scriptPath, "pid", cmd.Process.Pid)
	return nil
}

// Kills the current process (if any). Must be called with the mutex held
func (cp *ClassifierProcess) killLocked() {
	if cp.cmd == nil {
		return
	}
	//Closing stdin makes the python loop end gracefully
	if cp.stdin != nil {
		cp.stdin.Close()
	}
	if cp.exited != nil {
		select {
		case <-cp.exited:
		case <-time.After(time.Second):
			//The process did not exit by itself, kill it
			if cp.cmd.Process != nil {
				cp.cmd.Process.Kill()
			}
			<-cp.exited
		}
	}
	// Drain remaining stdout so the reader goroutine can finish
	if cp.lines != nil {
		go func(lines chan string) {
			for range lines {
			}
		}(cp.lines)
	}
	cp.cmd = nil
	cp.stdin = nil
	cp.lines = nil
	cp.exited = nil
}

// Returns true if the process is running. Must be called with the mutex held
func (cp *ClassifierProcess) aliveLocked() bool {
	if cp.cmd == nil {
		return false
	}
	select {
	case <-cp.exited:
		return false
	default:
		return true
	}
}

// Classify sends the features to the python process and returns the classification
func (cp *ClassifierProcess) Classify(features string) (string, error) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	// Restart the process if it died
	if !cp.aliveLocked() {
		cp.killLocked()
		cp.logger.Warning("Classifier process is not running, restarting it")
		if err := cp.startLocked(); err != nil {
			return "", err
		}
	}

	if _, err := io.WriteString(cp.stdin, features+"\n"); err != nil {
		cp.killLocked()
		return "", fmt.Errorf("could not write features to classifier: %w", err)
	}

	select {
	case line, ok := <-cp.lines:
		if !ok {
			cp.killLocked()
			return "", errors.New("classifier process closed its output")
		}
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "ERROR") {
			return "", fmt.Errorf("classifier returned an error: %s", line)
		}
		return line, nil
	case <-time.After(classifierRequestTimeout):
		// The stream is now out of sync, restart on next request
		cp.killLocked()
		return "", errors.New("timed out waiting for classifier response")
	}
}

// Stop terminates the python process
func (cp *ClassifierProcess) Stop() {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	cp.killLocked()
}
