# gh-proxy-go

简单的 GitHub 文件代理工具。修改自 [hunshcn/gh-proxy](https://github.com/hunshcn/gh-proxy)，但是更轻量。

A simple Go proxy for GitHub API, but MUCH lighter. Ported & modified from [hunshcn/gh-proxy](https://github.com/hunshcn/gh-proxy).

已知问题：不支持 Git clone。制作本项目主要是为了节约硬盘空间，本人未用到其 git clone 代理功能，因此没有实现（或者说修复），欢迎 PR。

## 用法

### Docker

```bash
# 直接运行，只允许代理 GitHub 链接
docker run -d -p 80:80 --name gh-proxy-go anotia/gh-proxy-go
# 允许代理任意链接
docker run -d -p 80:80 --name gh-proxy-go anotia/gh-proxy-go --allow-any-url
```

### 命令行运行

在 Releases 下载对应平台的二进制文件运行，或者使用 Go 编译后运行。

```bash
# 直接运行在 0.0.0.0:80
./gh-proxy-go

# 修改端口 / 地址
./gh-proxy-go --port 8080 --host 127.0.0.1

# 允许代理任意链接
./gh-proxy-go --allow-any-url
```
