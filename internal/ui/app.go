package ui

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/exec"
	"portmanager/internal/core"
	"portmanager/internal/util"
	"runtime"
	"sync"
)

type App struct {
	monitor *core.PortMonitor
	mu      sync.Mutex
	running bool
	server  *http.Server
}

func NewApp() *App {
	return &App{
		monitor: core.NewPortMonitor(),
	}
}

type ScanResult struct {
	Success bool              `json:"success"`
	Ports   []core.PortInfo   `json:"ports"`
	Message string            `json:"message"`
	Count   int               `json:"count"`
}

func (a *App) Run() error {
	// 设置 HTTP 路由
	http.HandleFunc("/", a.handleIndex)
	http.HandleFunc("/api/scan", a.handleScan)
	http.HandleFunc("/api/close-port", a.handleClosePort)
	http.HandleFunc("/api/startup", a.handleAutoStartup)
	http.HandleFunc("/api/startup-status", a.handleStartupStatus)

	// 找一个可用的端口
	port, err := findAvailablePort()
	if err != nil {
		return fmt.Errorf("无法找到可用端口: %v", err)
	}

	addr := fmt.Sprintf("127.0.0.1:%d", port)
	a.server = &http.Server{Addr: addr}

	log.Printf("✓ PortManager 启动在 http://%s\n", addr)

	// 在浏览器中打开
	go openBrowser(fmt.Sprintf("http://%s", addr))

	// 启动服务器
	return a.server.ListenAndServe()
}

func (a *App) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(getHTMLContent()))
}

func (a *App) handleScan(w http.ResponseWriter, r *http.Request) {
	a.mu.Lock()
	if a.running {
		a.mu.Unlock()
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "扫描正在进行中"})
		return
	}
	a.running = true
	a.mu.Unlock()

	defer func() {
		a.mu.Lock()
		a.running = false
		a.mu.Unlock()
	}()

	ports, err := a.monitor.ScanPorts()
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	result := ScanResult{
		Success: true,
		Ports:   ports,
		Count:   len(ports),
		Message: fmt.Sprintf("发现 %d 个开放端口", len(ports)),
	}

	respondJSON(w, http.StatusOK, result)
}

func (a *App) handleClosePort(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	var req struct {
		Port int `json:"port"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}

	if err := core.ClosePort(req.Port); err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("端口 %d 已关闭", req.Port)})
}

func (a *App) handleAutoStartup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	var req struct {
		Enabled bool `json:"enabled"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}

	var err error
	if req.Enabled {
		err = util.EnableAutoStartup()
	} else {
		err = util.DisableAutoStartup()
	}

	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "已更新"})
}

func (a *App) handleStartupStatus(w http.ResponseWriter, r *http.Request) {
	enabled := util.IsAutoStartupEnabled()
	respondJSON(w, http.StatusOK, map[string]bool{"enabled": enabled})
}

func respondJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

func findAvailablePort() (int, error) {
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 0,
	})
	if err != nil {
		return 0, err
	}
	defer listener.Close()

	addr := listener.Addr().(*net.TCPAddr)
	return addr.Port, nil
}

func openBrowser(url string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}

	cmd.Run()
}

func getHTMLContent() string {
	return `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>PortManager</title>
    <style>
        *{margin:0;padding:0;box-sizing:border-box}
        body{font-family:system-ui,"Segoe UI",sans-serif;background:linear-gradient(135deg,#667eea,#764ba2);min-height:100vh;display:flex;align-items:center;justify-content:center;padding:20px}
        .container{width:100%;max-width:900px;background:#fff;border-radius:16px;box-shadow:0 20px 60px rgba(0,0,0,.3)}
        .header{background:linear-gradient(135deg,#667eea,#764ba2);color:#fff;padding:30px;text-align:center}
        .header h1{font-size:28px;margin-bottom:8px}
        .content{padding:30px}
        .grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(300px,1fr));gap:20px;margin-bottom:30px}
        .card{background:#f8f9fa;border-radius:12px;padding:20px;border:1px solid #e9ecef;transition:all .3s}
        .card:hover{box-shadow:0 8px 24px rgba(0,0,0,.08);border-color:#667eea}
        .card-title{font-size:18px;font-weight:600;margin-bottom:16px;color:#333}
        .btn{padding:10px 20px;border:none;border-radius:8px;font-size:14px;cursor:pointer;transition:all .3s;font-weight:500}
        .btn-primary{background:linear-gradient(135deg,#667eea,#764ba2);color:#fff;width:100%}
        .btn-primary:hover{transform:translateY(-2px);box-shadow:0 8px 16px rgba(102,126,234,.4)}
        .btn-primary:disabled{opacity:.6;cursor:not-allowed}
        .btn-danger{background:#ff6b6b;color:#fff;padding:6px 12px;font-size:12px}
        .btn-danger:hover{background:#ff5252}
        .status{padding:12px;background:#e3f2fd;border-left:4px solid #667eea;border-radius:4px;margin:10px 0;font-size:14px;color:#1976d2;display:none}
        .port-list{max-height:300px;overflow-y:auto}
        .port-item{padding:10px;background:#fff;border:1px solid #e9ecef;border-radius:6px;margin:8px 0;display:flex;justify-content:space-between;align-items:center;font-size:13px}
        .port-number{font-weight:600;color:#667eea}
        .port-process{color:#666;font-size:12px;margin-top:4px}
        .checkbox-group{margin:12px 0}
        .checkbox-group label{display:flex;align-items:center;gap:8px;cursor:pointer;font-size:14px}
        .checkbox-group input{width:18px;height:18px}
        .spinner{display:inline-block;width:16px;height:16px;border:2px solid #f3f3f3;border-top:2px solid #667eea;border-radius:50%;animation:spin 1s linear infinite}
        @keyframes spin{0%{transform:rotate(0deg)}100%{transform:rotate(360deg)}}
        .footer{background:#f8f9fa;padding:20px;display:flex;justify-content:space-between;border-top:1px solid #e9ecef;font-size:12px;color:#999}
        #resultCount{font-size:24px;font-weight:bold;color:#667eea;margin:10px 0}
    </style>
</head>
<body>
    <div class="container">
        <div class="header"><h1>⚙️ PortManager</h1><p>端口检测管理工具</p></div>
        <div class="content">
            <div class="grid">
                <div class="card">
                    <div class="card-title">🔍 端口扫描</div>
                    <button class="btn btn-primary" id="scanBtn" onclick="performScan()"><span id="scanText">扫描端口</span></button>
                    <div class="status" id="scanStatus"></div>
                    <div id="scanResult" style="margin-top:15px"></div>
                </div>
                <div class="card">
                    <div class="card-title">📊 扫描结果</div>
                    <div id="resultCount">0</div>
                    <div class="port-list" id="portList"><div style="color:#999;text-align:center;padding:20px">点击"扫描端口"开始</div></div>
                </div>
                <div class="card">
                    <div class="card-title">⚙️ 设置</div>
                    <div class="checkbox-group">
                        <label><input type="checkbox" id="autoStartup" onchange="toggleAutoStartup()">开机自启动</label>
                    </div>
                    <button class="btn btn-primary" onclick="showAbout()" style="width:100%;margin-top:10px">关于</button>
                </div>
            </div>
        </div>
        <div class="footer"><span>PortManager v1.0.0</span><span>💡 某些操作需要管理员权限</span></div>
    </div>
    <script>
        window.addEventListener('load',()=>loadAutoStartupStatus());
        function performScan(){
            const btn=document.getElementById('scanBtn'),status=document.getElementById('scanStatus'),result=document.getElementById('resultCount'),list=document.getElementById('portList');
            btn.disabled=true;status.style.display='block';status.innerHTML='<span class="spinner"></span> 扫描中...';
            fetch('/api/scan').then(r=>r.json()).then(data=>{
                if(data.success){
                    status.innerHTML='✓ '+data.message;status.style.background='#c8e6c9';status.style.color='#2e7d32';
                    result.textContent=data.count;
                    if(data.ports&&data.ports.length>0){
                        list.innerHTML=data.ports.map(p=>'<div class="port-item"><div><div class="port-number">端口 '+p.Port+'</div><div class="port-process">'+(p.Process||'未知')+' (PID: '+p.PID+')</div></div><button class="btn btn-danger" onclick="closePort('+p.Port+')">关闭</button></div>').join('')
                    }else{list.innerHTML='<div style="color:#999;text-align:center;padding:20px">未发现开放端口</div>'}
                }else{
                    status.innerHTML='✗ 错误: '+data.error;status.style.background='#ffcdd2';status.style.color='#c62828'
                }
                btn.disabled=false
            }).catch(e=>{status.innerHTML='✗ 错误: '+e.message;status.style.background='#ffcdd2';btn.disabled=false})
        }
        function closePort(port){
            if(!confirm('确认关闭端口 '+port+'?'))return;
            fetch('/api/close-port',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({port:port})}).then(r=>r.json()).then(data=>{alert(data.message||data.error);if(data.message)performScan()}).catch(e=>alert('错误: '+e.message))
        }
        function toggleAutoStartup(){
            const cb=document.getElementById('autoStartup');
            fetch('/api/startup',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({enabled:cb.checked})}).then(r=>r.json()).catch(e=>alert('错误: '+e.message))
        }
        function loadAutoStartupStatus(){
            fetch('/api/startup-status').then(r=>r.json()).then(data=>{document.getElementById('autoStartup').checked=data.enabled}).catch(e=>console.error(e))
        }
        function showAbout(){alert('PortManager v1.0.0\n\n端口检测与管理工具\n\n支持端口扫描、流量监控、自动关闭等功能')}
    </script>
</body>
</html>`
}
