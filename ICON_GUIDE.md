# PortManager v1.1.0 图标集成指南

## 🎨 新增功能

PortManager v1.1.0 现在包含专业的应用程序图标设计，具有现代Windows风格。

### 图标设计特点

**视觉元素**：
- 🔵 **深蓝色背景** - 代表网络、稳定和专业性
- ⚪ **中心白点** - 代表主监听端口
- 🟡 **5个金黄色端口点** - 代表多个网络端口
- 🔗 **连接线** - 代表端口间的网络连接
- **P 字母** - "Port"（端口）的首字母

### 显示位置

图标将在以下位置显示：

1. **文件浏览器** 📁 - 文件详情中显示
2. **任务栏** 📌 - 应用运行时显示
3. **系统托盘** 🔔 - 应用最小化到托盘时显示
4. **窗口标题栏** 🪟 - 应用窗口的左上角
5. **快捷方式** 🎯 - 如果创建了桌面快捷方式

## 📦 文件包含

新版分发包包含以下图标相关文件：

```
PortManager-v1.1.0-windows-amd64.zip
├── PortManager.exe           ✓ 已嵌入图标和manifest
├── PortManager.ico           图标文件（多分辨率）
├── PortManager.exe.manifest  Windows manifest（Modern controls）
├── ICON.md                  图标详细说明
└── 其他文件...
```

## 🔧 技术实现

### 资源集成流程

```
1. Python生成图标
   generate_icon.py → PortManager.png + PortManager.ico
   
2. rsrc工具编译资源
   rsrc + manifest + ico → rsrc_windows_amd64.syso
   
3. Go编译器链接
   go build + .syso → PortManager.exe (含资源)
```

### 关键工具

- **PIL/Pillow** - 用于生成多分辨率ICO文件
- **rsrc** - 将manifest和图标嵌入到exe的Windows资源
- **Go编译器** - 自动识别和链接.syso资源文件

## 🛠️ 自定义图标

### 修改图标

如需修改图标设计（颜色、尺寸、形状等）：

1. **编辑Python脚本**
   ```bash
   # 修改generate_icon.py中的颜色值和绘图代码
   nano generate_icon.py
   ```

2. **重新生成图标**
   ```bash
   python generate_icon.py
   ```

3. **重新编译应用**
   ```bash
   build_gui.bat
   # 或
   go build -ldflags "-H=windowsgui" -o PortManager.exe .
   ```

### 使用现有图标

如果要替换为其他ICO文件：

1. 将新的 `icon.ico` 放在项目根目录
2. 重命名为 `PortManager.ico`
3. 运行 `build_gui.bat` 重新编译

## ✅ 验证图标

### 在Windows中查看

1. **文件属性**
   - 右击 PortManager.exe
   - 选择"属性"
   - 查看"详情"标签的图标缩略图

2. **运行应用**
   - 双击 run.bat
   - 查看任务栏图标
   - 关闭到托盘时查看系统托盘

### 技术验证

```powershell
# 检查资源嵌入
D:\CodeFromGoland\bin\rsrc.exe 可用时自动嵌入

# 查看exe大小变化
# 带资源的exe会比基础版本大约2-3MB
```

## 📋 版本变更

### v1.1.0 新增

- ✨ 专业应用图标（5个分辨率级别）
- 🎨 现代Windows控制样式（via manifest）
- 🔧 自动资源嵌入构建流程
- 📖 完整的图标自定义文档

### 与v1.0.0的区别

| 特性 | v1.0.0 | v1.1.0 |
|-----|--------|--------|
| 应用图标 | ❌ 无 | ✅ 蓝色网络主题 |
| 控制样式 | ⚠️ 经典 | ✅ 现代Modern |
| Manifest | ❌ 无 | ✅ 已嵌入 |
| 文件体积 | ~8.3MB | ~8.4MB |

## 🚀 快速开始

### 使用预编译版本

```bash
# 直接运行（包含图标）
run.bat

# 或者在PowerShell中
.\run.bat
```

### 从源码构建

```bash
# 自动生成并嵌入图标
build_gui.bat

# 或手动
python generate_icon.py
D:\CodeFromGoland\bin\rsrc.exe -manifest PortManager.exe.manifest -ico PortManager.ico
go build -ldflags "-H=windowsgui" -o PortManager.exe .
```

## 💡 常见问题

**Q: 为什么我的图标不显示？**
> A: 确保 `PortManager.exe` 是由rsrc工具生成的资源版本。检查exe大小应该在8.3MB以上。

**Q: 如何更改图标颜色？**
> A: 编辑 `generate_icon.py`，修改RGB颜色值，然后运行脚本和重新编译。

**Q: 支持透明背景吗？**
> A: 是的，ICO格式支持透明度。在 `generate_icon.py` 中修改 Image.new() 的background color参数。

**Q: 可以使用自己的PNG图标吗？**
> A: 可以。创建256x256的PNG，然后使用PIL或在线工具转换为ICO，替换 `PortManager.ico`。

## 📚 参考资源

- [Pillow文档](https://pillow.readthedocs.io/)
- [rsrc工具](https://github.com/akavel/rsrc)
- [Windows ICO格式规范](https://en.wikipedia.org/wiki/ICO_(file_format))
- [Go Build Tags](https://pkg.go.dev/cmd/cgo)

---

**最后更新**: 2026-04-16  
**版本**: v1.1.0
