package tunnel

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	quickTunnelRegex = regexp.MustCompile(`https://[a-zA-Z0-9-]+\.trycloudflare\.com`)
)

// Mode defines the tunneling method.
type Mode string

const (
	ModeNone        Mode = "none"
	ModeQuickTunnel Mode = "quick"
	ModeZeroTrust   Mode = "zero_trust"
)

// Status represents the live state of the tunnel.
type Status struct {
	Active              bool      `json:"active"`
	Mode                Mode      `json:"mode"`
	URL                 string    `json:"url,omitempty"`
	StartedAt           time.Time `json:"started_at,omitempty"`
	Error               string    `json:"error,omitempty"`
	CloudflaredDetected bool      `json:"cloudflared_detected"`
	CloudflaredPath     string    `json:"cloudflared_path,omitempty"`
}

// Manager supervises the cloudflared process lifecycle.
type Manager struct {
	mu        sync.RWMutex
	cmd       *exec.Cmd
	mode      Mode
	url       string
	startedAt time.Time
	lastErr   string
	cancel    context.CancelFunc
}

// NewManager constructs a new Tunnel Manager.
func NewManager() *Manager {
	return &Manager{
		mode: ModeNone,
	}
}

// DetectCloudflared checks whether cloudflared binary is installed.
func DetectCloudflared() (bool, string) {
	binName := "cloudflared"
	if runtime.GOOS == "windows" {
		binName = "cloudflared.exe"
	}

	// 1. Check PATH
	if p, err := exec.LookPath(binName); err == nil {
		return true, p
	}

	// 2. Check standard Windows probing locations
	if runtime.GOOS == "windows" {
		candidates := []string{
			`C:\Program Files\cloudflared\cloudflared.exe`,
			`C:\Program Files (x86)\cloudflared\cloudflared.exe`,
			filepathJoinSafe(os.Getenv("LOCALAPPDATA"), "Programs", "cloudflared", "cloudflared.exe"),
			filepathJoinSafe(os.Getenv("USERPROFILE"), "bin", "cloudflared.exe"),
			filepathJoinSafe(os.Getenv("USERPROFILE"), "cloudflared.exe"),
		}
		for _, c := range candidates {
			if c != "" {
				if fi, err := os.Stat(c); err == nil && !fi.IsDir() {
					return true, c
				}
			}
		}
	} else {
		candidates := []string{
			"/usr/local/bin/cloudflared",
			"/usr/bin/cloudflared",
			"/opt/cloudflared/cloudflared",
		}
		for _, c := range candidates {
			if fi, err := os.Stat(c); err == nil && !fi.IsDir() {
				return true, c
			}
		}
	}

	return false, ""
}

func filepathJoinSafe(base string, parts ...string) string {
	if base == "" {
		return ""
	}
	all := append([]string{base}, parts...)
	return strings.Join(all, string(os.PathSeparator))
}

// GetStatus returns the current status snapshot.
func (m *Manager) GetStatus() Status {
	m.mu.RLock()
	defer m.mu.RUnlock()

	detected, path := DetectCloudflared()

	active := m.cmd != nil && m.cmd.Process != nil
	return Status{
		Active:              active,
		Mode:                m.mode,
		URL:                 m.url,
		StartedAt:           m.startedAt,
		Error:               m.lastErr,
		CloudflaredDetected: detected,
		CloudflaredPath:     path,
	}
}

// StartQuickTunnel spawns cloudflared tunnel --url http://127.0.0.1:<port>
func (m *Manager) StartQuickTunnel(port int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.cmd != nil && m.cmd.Process != nil {
		return errors.New("a tunnel is already running; stop it before starting another")
	}

	detected, binPath := DetectCloudflared()
	if !detected {
		return errors.New("cloudflared binary not found; please install cloudflared first")
	}

	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, binPath, "tunnel", "--url", fmt.Sprintf("http://127.0.0.1:%d", port))

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		cancel()
		return fmt.Errorf("failed to open stderr pipe: %w", err)
	}

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		cancel()
		return fmt.Errorf("failed to open stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		cancel()
		return fmt.Errorf("failed to start cloudflared: %w", err)
	}

	m.cmd = cmd
	m.cancel = cancel
	m.mode = ModeQuickTunnel
	m.startedAt = time.Now()
	m.url = ""
	m.lastErr = ""

	// Process output asynchronously to discover the trycloudflare.com URL
	urlFound := make(chan string, 1)

	readOutput := func(r io.Reader) {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			line := scanner.Text()
			if m := quickTunnelRegex.FindString(line); m != "" {
				select {
				case urlFound <- m:
				default:
				}
			}
		}
	}

	go readOutput(stderrPipe)
	go readOutput(stdoutPipe)

	go func() {
		_ = cmd.Wait()
		m.mu.Lock()
		if m.cmd == cmd {
			m.cmd = nil
			m.mode = ModeNone
			m.url = ""
		}
		m.mu.Unlock()
	}()

	// Wait up to 10 seconds for the URL to be emitted
	go func() {
		select {
		case u := <-urlFound:
			m.mu.Lock()
			m.url = u
			m.mu.Unlock()
		case <-time.After(15 * time.Second):
			m.mu.Lock()
			if m.url == "" && m.cmd != nil {
				m.lastErr = "timeout waiting for trycloudflare.com tunnel URL"
			}
			m.mu.Unlock()
		}
	}()

	return nil
}

// StartTokenTunnel spawns cloudflared tunnel run --token <token>
func (m *Manager) StartTokenTunnel(token string) error {
	trimmed := strings.TrimSpace(token)
	if trimmed == "" {
		return errors.New("tunnel token cannot be empty")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.cmd != nil && m.cmd.Process != nil {
		return errors.New("a tunnel is already running; stop it before starting another")
	}

	detected, binPath := DetectCloudflared()
	if !detected {
		return errors.New("cloudflared binary not found; please install cloudflared first")
	}

	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, binPath, "tunnel", "run", "--token", trimmed)

	if err := cmd.Start(); err != nil {
		cancel()
		return fmt.Errorf("failed to start cloudflared Zero Trust tunnel: %w", err)
	}

	m.cmd = cmd
	m.cancel = cancel
	m.mode = ModeZeroTrust
	m.startedAt = time.Now()
	m.url = "Configured Domain (Zero Trust)"
	m.lastErr = ""

	go func() {
		_ = cmd.Wait()
		m.mu.Lock()
		if m.cmd == cmd {
			m.cmd = nil
			m.mode = ModeNone
			m.url = ""
		}
		m.mu.Unlock()
	}()

	return nil
}

// Stop terminates any active cloudflared process.
func (m *Manager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.cmd == nil || m.cmd.Process == nil {
		return nil
	}

	if m.cancel != nil {
		m.cancel()
	}

	_ = m.cmd.Process.Kill()
	m.cmd = nil
	m.mode = ModeNone
	m.url = ""
	m.lastErr = ""
	return nil
}
