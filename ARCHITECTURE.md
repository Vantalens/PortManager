# PortManager 架构说明

适用版本：v1.1.0

## 架构总览

PortManager 是一个纯本地 Windows 桌面应用，不依赖 Web 服务或 REST API。

- 语言：Go 1.25.5
- UI：Walk（原生窗口）
- 运行方式：直接调用系统命令 netstat、tasklist、taskkill、reg

## 模块分层

### main.go

- 程序入口
- 初始化配置
- 创建并运行 UI 应用

### internal/ui

- 窗口布局、事件绑定、用户交互
- 扫描任务并发控制
- 结果列表刷新、端口洞察弹窗、状态提示

### internal/core

- 端口扫描与解析
- 进程名批量映射
- 端口关闭
- 端口洞察与建议逻辑
- 中文系统输出兼容与编码处理

### internal/util

- 配置读写
- 开机自启动（注册表 Run 项）

## 核心数据流

### 扫描流程

1. UI 发起 快速扫描 或 全面扫描
2. core 执行 netstat -ano，解析监听端口
3. core 批量调用 tasklist /FO CSV /NH 映射 PID 到进程名
4. UI 刷新列表与数量统计

### 关闭流程

1. UI 读取当前选中端口
2. core 再次遍历 netstat 定位 PID
3. core 调用 taskkill /PID /F
4. UI 反馈结果并触发刷新

### 自启动流程

1. UI 切换 开机自启动
2. util 调用 reg add 或 reg delete
3. 同步更新配置状态

## 关键设计点

- 批量进程查询：降低扫描阶段命令调用次数
- 输出编码兼容：优先 UTF-8，不可用时回退 GB18030
- 中文状态兼容：识别 LISTENING、LISTEN、侦听
- UI 稳定性：最小化恢复后自动修正尺寸与右下角位置

## 目录结构

```text
PortManager/
├─ main.go
├─ internal/
│  ├─ ui/app.go
│  ├─ core/port.go
│  └─ util/
│     ├─ config.go
│     └─ startup.go
├─ run.bat
├─ run.vbs
├─ build_gui.bat
└─ vendor/
```

## 非目标范围

- 不提供 HTTP 接口
- 不做远程控制
- 不做驱动级流量拦截