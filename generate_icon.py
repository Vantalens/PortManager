#!/usr/bin/env python3
"""
生成PortManager应用图标
设计概念：网络端口、连接、监控
"""

from PIL import Image, ImageDraw, ImageFont
import os

def create_icon():
    # 创建512x512基础图像（高分辨率）
    size = 512
    image = Image.new('RGBA', (size, size), (0, 0, 0, 0))
    draw = ImageDraw.Draw(image)
    
    # 背景圆形（渐变效果模拟）- 使用深蓝色到蓝色
    # 绘制外圆（深蓝色背景）
    bg_color = (25, 118, 210)  # 深蓝色
    draw.ellipse([10, 10, size-10, size-10], fill=bg_color, outline=bg_color)
    
    # 内圆（略浅的蓝色）
    inner_margin = 50
    inner_color = (33, 150, 243)  # 更亮的蓝色
    draw.ellipse([inner_margin, inner_margin, size-inner_margin, size-inner_margin], 
                 fill=inner_color, outline=inner_color)
    
    # 绘制网络连接点（代表端口）
    # 中心点
    center = size // 2
    center_radius = 40
    draw.ellipse([center-center_radius, center-center_radius, 
                  center+center_radius, center+center_radius], 
                 fill=(255, 255, 255), outline=(255, 255, 255))
    
    # 周围的5个端口点
    num_ports = 5
    circle_radius = size // 3
    port_radius = 28
    port_color = (255, 193, 7)  # 金黄色
    
    import math
    for i in range(num_ports):
        angle = (2 * math.pi * i) / num_ports - math.pi / 2
        x = center + circle_radius * math.cos(angle)
        y = center + circle_radius * math.sin(angle)
        
        # 绘制端口连接线
        draw.line([(center, center), (x, y)], fill=(200, 200, 200), width=3)
        
        # 绘制端口点
        draw.ellipse([x-port_radius, y-port_radius, x+port_radius, y+port_radius],
                     fill=port_color, outline=(255, 255, 255), width=2)
    
    # 添加"P"字母在中心（代表Port）
    try:
        # 尝试使用系统字体
        font = ImageFont.truetype("C:\\Windows\\Fonts\\segoeui.ttf", 80)
    except:
        # 如果没有指定字体，使用默认字体
        font = ImageFont.load_default()
    
    text = "P"
    bbox = draw.textbbox((0, 0), text, font=font)
    text_width = bbox[2] - bbox[0]
    text_height = bbox[3] - bbox[1]
    text_x = center - text_width // 2
    text_y = center - text_height // 2
    
    draw.text((text_x, text_y), text, fill=(33, 150, 243), font=font)
    
    # 调整大小并保存为不同规格的ICO
    output_path = os.path.join(os.path.dirname(__file__), 'PortManager.ico')
    
    # 保存为ICO（包含多个分辨率）
    sizes = [(16, 16), (32, 32), (48, 48), (64, 64), (128, 128), (256, 256)]
    images = []
    
    for size_tuple in sizes:
        resized = image.resize(size_tuple, Image.Resampling.LANCZOS)
        images.append(resized)
    
    # 保存ICO文件
    images[0].save(output_path, format='ICO', sizes=sizes)
    print(f"✓ 图标已生成: {output_path}")
    
    # 也保存一份PNG作为备用
    png_path = os.path.join(os.path.dirname(__file__), 'PortManager.png')
    image.save(png_path, 'PNG')
    print(f"✓ PNG图标已生成: {png_path}")
    
    return output_path

if __name__ == '__main__':
    create_icon()
