# PortManager 图标说明

## 文件说明

- **PortManager.ico** - 应用程序图标（包含16x16到256x256多个分辨率）
- **PortManager.png** - 图标的高分辨率PNG版本（512x512）
- **generate_icon.py** - Python脚本，用于生成ICO文件
- **rsrc_windows_amd64.syso** - 编译后的资源文件（自动生成）

## 图标设计

图标设计理念：
- **深蓝色背景**：代表网络、稳定性和专业性
- **中心白点**：代表主端口/监听点
- **周围金黄色点和连线**：代表多个端口连接和实时监控
- **P字母**：代表"Port"（端口）

## 重新生成图标

如果需要修改或重新生成图标，请运行：

```bash
python generate_icon.py
```

然后使用build_gui.bat重新编译应用程序。

## 集成到应用

rsrc工具会自动将以下资源嵌入到PortManager.exe中：
- Manifest（PortManager.exe.manifest）- 用于Modern Windows控制
- 图标（PortManager.ico）- 在资源管理器和任务栏中显示

## 在不同场景下的显示

1. **文件浏览器** - 显示应用图标
2. **任务栏** - 显示应用图标
3. **开始菜单** - 显示应用图标（如果安装了快捷方式）
4. **应用窗口标题栏** - 显示图标
5. **系统托盘** - 应用启用了托盘功能时显示

## 自定义图标

如果要创建自己的图标，可以：
1. 修改generate_icon.py中的颜色值和设计元素
2. 或使用第三方工具（如GIMP、Photoshop）创建256x256的PNG
3. 使用在线转换工具或命令行工具将PNG转换为ICO

## 技术细节

- 使用Python PIL库生成多分辨率ICO
- 使用rsrc工具将资源嵌入到Go编译的exe中
- .syso文件由Go链接器自动识别和链接
