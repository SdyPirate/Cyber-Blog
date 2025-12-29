# Big B Bro - 个人博客系统 / Personal Blog System

<div align="center">

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)
![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)
![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux%20%7C%20macOS-lightgrey?style=flat-square)

**一个轻量级、现代化的个人博客系统 / A lightweight, modern personal blog system**

[中文](#-中文文档) | [English](#-english-documentation)

</div>

---

## 🌟 功能特点 / Features

| 中文 | English |
|------|---------|
| 🎨 赛博朋克风格 UI | 🎨 Cyberpunk-style UI |
| 📝 Markdown 博客支持 | 📝 Markdown blog support |
| 🔐 安全的后台管理 | 🔐 Secure admin panel |
| 🏷️ 标签分类系统 | 🏷️ Tag classification system |
| 📱 响应式设计 | 📱 Responsive design |
| ⚡ 单文件部署 | ⚡ Single binary deployment |
| 🗂️ JSON 数据存储 | 🗂️ JSON data storage |

---

# 📖 中文文档

## 快速开始

### 1. 下载

从 [Releases](https://github.com/your-username/big-b-bro/releases) 下载适合你操作系统的版本。

### 2. 运行

```bash
# Windows
.\blog.exe

# Linux / macOS
chmod +x blog
./blog
```

### 3. 访问

- **博客首页**: http://localhost:12302/blog
- **后台管理**: http://localhost:12302/admin
- **默认账号**: `admin`
- **默认密码**: `admin123`

## 项目结构

```
big-b-bro/
├── main.go              # 程序入口
├── handlers/            # HTTP 处理器
│   ├── admin.go         # 管理员认证
│   ├── blog.go          # 博客前台
│   ├── blog_admin.go    # 博客后台管理
│   └── content.go       # 内容管理
├── models/              # 数据模型
│   ├── content.go       # 内容模型
│   └── post.go          # 博客文章模型
├── middleware/          # 中间件
│   └── auth.go          # 认证中间件
├── templates/           # HTML 模板
├── css/                 # 样式文件
├── js/                  # JavaScript 文件
└── data/                # 数据存储 (运行时生成)
    ├── admin.json       # 管理员账号
    ├── content.json     # 网站内容
    └── posts.json       # 博客文章
```

## 从源码构建

```bash
# 克隆项目
git clone https://github.com/your-username/big-b-bro.git
cd big-b-bro

# 安装依赖
go mod tidy

# 运行开发服务器
go run main.go

# 构建可执行文件
go build -o blog.exe .

# 跨平台构建
GOOS=linux GOARCH=amd64 go build -o blog-linux .
GOOS=darwin GOARCH=amd64 go build -o blog-macos .
```

## 配置说明

应用默认运行在端口 `12302`，如需修改请编辑 `main.go` 中的 `port` 变量。

### 修改管理员密码

1. 删除 `data/admin.json` 文件
2. 重启应用
3. 使用默认密码 `admin123` 登录
4. 在后台修改密码（功能开发中）

或手动编辑 `data/admin.json`，使用 bcrypt 加密后的密码。

## 部署指南

### 本地运行

```bash
./blog.exe
```

### Linux 服务器部署 (systemd)

```bash
# 创建服务文件
sudo nano /etc/systemd/system/blog.service
```

```ini
[Unit]
Description=Big B Bro Blog
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/blog
ExecStart=/opt/blog/blog
Restart=always

[Install]
WantedBy=multi-user.target
```

```bash
# 启动服务
sudo systemctl enable blog
sudo systemctl start blog
```

### Nginx 反向代理

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:12302;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

---

# 📖 English Documentation

## Quick Start

### 1. Download

Download the appropriate version for your OS from [Releases](https://github.com/your-username/big-b-bro/releases).

### 2. Run

```bash
# Windows
.\blog.exe

# Linux / macOS
chmod +x blog
./blog
```

### 3. Access

- **Blog Homepage**: http://localhost:12302/blog
- **Admin Panel**: http://localhost:12302/admin
- **Default Username**: `admin`
- **Default Password**: `admin123`

## Project Structure

```
big-b-bro/
├── main.go              # Entry point
├── handlers/            # HTTP handlers
│   ├── admin.go         # Admin authentication
│   ├── blog.go          # Blog frontend
│   ├── blog_admin.go    # Blog admin management
│   └── content.go       # Content management
├── models/              # Data models
│   ├── content.go       # Content model
│   └── post.go          # Blog post model
├── middleware/          # Middleware
│   └── auth.go          # Auth middleware
├── templates/           # HTML templates
├── css/                 # Stylesheets
├── js/                  # JavaScript files
└── data/                # Data storage (generated at runtime)
    ├── admin.json       # Admin credentials
    ├── content.json     # Website content
    └── posts.json       # Blog posts
```

## Build from Source

```bash
# Clone the repository
git clone https://github.com/your-username/big-b-bro.git
cd big-b-bro

# Install dependencies
go mod tidy

# Run development server
go run main.go

# Build executable
go build -o blog.exe .

# Cross-platform builds
GOOS=linux GOARCH=amd64 go build -o blog-linux .
GOOS=darwin GOARCH=amd64 go build -o blog-macos .
```

## Configuration

The app runs on port `12302` by default. To change it, edit the `port` variable in `main.go`.

### Changing Admin Password

1. Delete `data/admin.json`
2. Restart the application
3. Login with default password `admin123`
4. Change password in the admin panel (coming soon)

Or manually edit `data/admin.json` with a bcrypt-hashed password.

## Deployment Guide

### Local Run

```bash
./blog.exe
```

### Linux Server Deployment (systemd)

```bash
# Create service file
sudo nano /etc/systemd/system/blog.service
```

```ini
[Unit]
Description=Big B Bro Blog
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/blog
ExecStart=/opt/blog/blog
Restart=always

[Install]
WantedBy=multi-user.target
```

```bash
# Start the service
sudo systemctl enable blog
sudo systemctl start blog
```

### Nginx Reverse Proxy

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:12302;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

---

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

---

<div align="center">
Made with ❤️ by Big B Bro
</div>
