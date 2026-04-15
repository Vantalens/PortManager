package core

import (
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type PortInfo struct {
	Port        int
	Process     string
	PID         int32
	Path        string
	Protocol    string
	State       string // OPEN, CLOSED
	TrafficRate float64
	RiskLevel   int // 0: normal, 1: medium, 2: critical
}

type PortMonitor struct {
	ports        map[int]*PortInfo
	lastCheck    time.Time
	trafficHist  map[int][]float64
}

func NewPortMonitor() *PortMonitor {
	return &PortMonitor{
		ports:       make(map[int]*PortInfo),
		trafficHist: make(map[int][]float64),
	}
}

// ScanPorts 扫描常用端口
func (pm *PortMonitor) ScanPorts() ([]PortInfo, error) {
	commonPorts := []int{
		22, 80, 443, 445, 1433, 3306, 3389, 5432, 6379,
		8000, 8080, 8443, 8888, 9000, 9200, 27017, 50070,
	}

	var results []PortInfo

	for _, port := range commonPorts {
		if info := pm.checkPort(port); info != nil {
			results = append(results, *info)
		}
	}

	pm.lastCheck = time.Now()
	return results, nil
}

func (pm *PortMonitor) checkPort(port int) *PortInfo {
	// 检查 TCP
	if pm.isPortOpen(port, "tcp") {
		info := &PortInfo{
			Port:     port,
			Protocol: "TCP",
			State:    "OPEN",
		}

		// 获取进程信息
		if proc, err := pm.getProcessByPort(port); err == nil {
			info.Process = proc.Name
			info.PID = proc.PID
			info.Path = proc.Path
		}

		return info
	}

	return nil
}

func (pm *PortMonitor) isPortOpen(port int, protocol string) bool {
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	conn, err := net.DialTimeout(protocol, addr, 100*time.Millisecond)
	if err == nil {
		conn.Close()
		return true
	}
	return false
}

type ProcessInfo struct {
	Name string
	PID  int32
	Path string
}

func (pm *PortMonitor) getProcessByPort(port int) (*ProcessInfo, error) {
	cmd := exec.Command("netstat", "-ano")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	portStr := fmt.Sprintf(":%d", port)
	for _, line := range strings.Split(string(output), "\n") {
		if strings.Contains(line, portStr) && strings.Contains(line, "LISTENING") {
			fields := strings.Fields(line)
			if len(fields) >= 5 {
				pidStr := fields[len(fields)-1]
				pid, _ := strconv.ParseInt(pidStr, 10, 32)

				name, path := pm.getProcessInfo(int32(pid))
				return &ProcessInfo{
					Name: name,
					PID:  int32(pid),
					Path: path,
				}, nil
			}
		}
	}

	return nil, fmt.Errorf("process not found")
}

func (pm *PortMonitor) getProcessInfo(pid int32) (string, string) {
	cmd := exec.Command("wmic", "process", "where", fmt.Sprintf("ProcessId=%d", pid), "get", "Name,ExecutablePath")
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		if len(lines) >= 2 {
			fields := strings.Fields(lines[1])
			if len(fields) >= 1 {
				name := fields[0]
				path := ""
				if len(fields) > 1 {
					path = strings.Join(fields[1:], " ")
				}
				return name, path
			}
		}
	}
	return fmt.Sprintf("PID_%d", pid), ""
}

// ClosePort 关闭端口
func ClosePort(port int) error {
	cmd := exec.Command("netstat", "-ano")
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	portStr := fmt.Sprintf(":%d", port)
	for _, line := range strings.Split(string(output), "\n") {
		if strings.Contains(line, portStr) && strings.Contains(line, "LISTENING") {
			fields := strings.Fields(line)
			if len(fields) >= 5 {
				pidStr := fields[len(fields)-1]
				killCmd := exec.Command("taskkill", "/PID", pidStr, "/F")
				if err := killCmd.Run(); err == nil {
					return nil
				}
			}
		}
	}

	return fmt.Errorf("failed to close port %d", port)
}
