# 🚀 HTTPS Proxy with TLS Fingerprinting

一个基于 Cloudflare Workers 和 Go 的高级 HTTPS 代理工具，支持 TLS 指纹伪装和代码混淆。

## ✨ 特性

- 🔒 **TLS 指纹识别** - 使用 uTLS 库模拟浏览器指纹，绕过反爬虫检测 [[1](https://www.bing.com/ck/a?!&&p=2717aa131993f50bfe222f197f359521cc4dd254821a9bed4697f24522d664f8JmltdHM9MTc2NDYzMzYwMA&ptn=3&ver=2&hsh=4&fclid=285629af-c457-69b6-0834-3ab3c5c76865&u=a1aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L2dpdGJsb2dfMDA3MDYvYXJ0aWNsZS9kZXRhaWxzLzE0MjA4MTE4MA&ntb=1)][[5](https://www.bing.com/ck/a?!&&p=8b38ada0799a325c5726651ebedb45c316994f2463ad8b86b66d164fdfe0716dJmltdHM9MTc2NDYzMzYwMA&ptn=3&ver=2&hsh=4&fclid=285629af-c457-69b6-0834-3ab3c5c76865&u=a1aHR0cHM6Ly9yZXEuY29vbC96aC9kb2NzL3R1dG9yaWFsL3Rscy1maW5nZXJwcmludC8&ntb=1)]
- 🛡️ **代码混淆** - 保护 Worker 代码安全
- 🔑 **密码保护** - 环境变量配置密码
- ⚡ **高性能** - 基于 Go 编译的本地可执行文件

## 📦 快速开始

### 1️⃣ Cloudflare Workers 部署

1. 访问 [代码混淆工具](https://obfuscator.io/)
2. 混淆 `_worker.js` 代码
3. 在 Cloudflare Workers 中设置环境变量 `PASSWORD`

### 2️⃣ 本地编译运行
```bash
# 初始化 Go 模块
go mod init {自定义名字}

# 设置代理（可选）
set https_proxy=http://{自定义IP}:{自定义端口}

# 下载依赖
go mod tidy

# 编译
go build
```

编译完成后会生成可执行文件，直接运行即可启动代理服务 [[2](https://www.bing.com/ck/a?!&&p=380eebf58e070c3a571020d5b690a02c30e46f1b29a8e0e3836020d13a89195aJmltdHM9MTc2NDYzMzYwMA&ptn=3&ver=2&hsh=4&fclid=285629af-c457-69b6-0834-3ab3c5c76865&u=a1aHR0cHM6Ly9nby5kZXYvZG9jL3R1dG9yaWFsL2NyZWF0ZS1tb2R1bGU&ntb=1)][[4](https://www.bing.com/ck/a?!&&p=ce5cb1f1b34c0e2dc9f2a60527a76e0d3cb880d96470ab05c78f3d96e4f39ad0JmltdHM9MTc2NDYzMzYwMA&ptn=3&ver=2&hsh=4&fclid=285629af-c457-69b6-0834-3ab3c5c76865&u=a1aHR0cHM6Ly93d3cuY25ibG9ncy5jb20vYWJjbGlmZS9wLzE4MDk2MTgy&ntb=1)]。

## 🔧 配置说明

| 配置项 | 说明 | 示例 |
|--------|------|------|
| `PASSWORD` | Worker 密码 | 在 CF 环境变量设置 |
| `https_proxy` | 编译时代理 | `http://127.0.0.1:7890` |
| 模块名称 | Go 模块名 | `myproxy` |

## 📚 技术栈

- **Go** - 主程序语言
- **uTLS** - TLS 指纹伪装 [[3](https://www.bing.com/ck/a?!&&p=728d4522c4f5f9b65a3641868c9927dc946d3a1714e8ae37f63eb007a7604582JmltdHM9MTc2NDYzMzYwMA&ptn=3&ver=2&hsh=4&fclid=285629af-c457-69b6-0834-3ab3c5c76865&u=a1aHR0cHM6Ly9naXRodWIuY29tL3JlZnJhY3Rpb24tbmV0d29ya2luZy91dGxz&ntb=1)]
- **Cloudflare Workers** - 边缘计算平台
- **Gorilla WebSocket** - WebSocket 支持

## 📝 许可证

MIT License
