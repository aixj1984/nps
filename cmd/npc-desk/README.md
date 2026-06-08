# NPC Desk

基于 **[Wails v2](https://github.com/wailsapp/wails)** 的 **NPC 内网穿透客户端图形界面**，保留 Blueprint 风格主界面，内嵌 `github.com/djylb/nps` 客户端核心。

支持完整 `npc.conf` 配置编辑、连接状态展示、隧道/域名/本地服务列表与运行日志。

---

## 旧版说明（Deskreen 投屏）

本仓库 `internal/` 仍保留 Deskreen 投屏相关代码；当前默认入口为 **NPC Desk**（`npc_app.go` / `main.go`）。

原先 Deskreen 说明如下：

使用 **[Wails v2](https://github.com/wailsapp/wails)** 仿制的 **[Deskreen CE](https://github.com/pavlobu/deskreen)**：将带浏览器的设备变成电脑副屏，通过 **WebRTC** 共享桌面。

| 原版 Deskreen | 本项目 |
|---------------|--------|
| Electron + React 主机 | **Wails + Vue 3** 主机 |
| 内嵌 React 查看端 | 内嵌 **Vue 3** 查看端 (`webui/`) |
| WebRTC 视频轨 (VP8/H264) | **信令/ICE/DataChannel 走 WebRTC**；画面为 **PrintWindow/显示器采集 + JPEG**（免 libvpx，后续可换媒体轨） |
| 端口 3131/3132 | 相同 |
| 三步向导 + 允许/拒绝 | 相同 |

## 功能（Deskreen CE 对齐）

- 步骤 1 **连接**：二维码、链接、`--ip` 指定局域网 IP
- 步骤 2 **选择**：整个屏幕 / 应用窗口（缩略图网格）
- 步骤 3 **确认**：预览后开始共享
- **单查看端**（第二台收到 ROOM_LOCKED）
- 查看端：播放 / 暂停 / 全屏 / 画质档位
- WiFi/LAN 检测提示

## 环境

- Go 1.22+
- Node.js 18+
- Windows 10/11 [WebView2](https://developer.microsoft.com/microsoft-edge/webview2/)

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
wails doctor
```

## NPC Desk 构建

```bat
cd cmd\npc-desk
build-npc.bat
```

产出：`npc-desk.exe`

首次运行会自动查找同目录或 `conf\npc.conf`，也可在界面中「打开配置」选择文件。

## 构建（Deskreen 旧流程）

按 [Wails 手动构建指南](https://wails.io/docs/guides/manual-builds/) 编译（不依赖 `wails build` 打包步骤）：

```bat
cd mochi-deskreen
build.bat
```

流程概要：

1. 构建 `webui` → `web/dist`（浏览器查看端）
2. 构建 `frontend` → `frontend/dist`（Wails 主机 UI）
3. （可选）`wails generate module` 生成 `wailsjs` 绑定
4. `go run ./tools/genwindows` 生成 `icon.ico` 与 `mochi-deskreen-res.syso`
5. `go build -tags desktop,production -ldflags "-w -s -H windowsgui"`

产出：`mochi-deskreen.exe`

调试版（保留控制台、`dev` 标签）：

```bat
build-debug.bat
```

开发热重载（需 Wails CLI）：

```bat
go install github.com/wailsapp/wails/v2/cmd/wails@latest
wails dev
```

### 启动失败排查

- 安装 [WebView2 运行时](https://developer.microsoft.com/microsoft-edge/webview2/)
- 查看 exe 同目录下的 `mochi-deskreen-startup.log`
- 勿运行旧的 `deskreen-gui.exe`（已废弃的 Fyne 版本）
- 勿单独 `go build` 而跳过 `tools/genwindows` 与 `desktop,production` 标签

## 使用

1. 运行 `mochi-deskreen.exe`
2. 浏览器打开 `http://<局域网IP>:3131/<房间ID>`（界面可复制/扫码）
3. 主机 **允许** → 选择屏幕或窗口 → **确认**
4. 查看端显示桌面画面

```bat
mochi-deskreen.exe --ip 192.168.1.100
```

## 结构

```
main.go / app.go     Wails 入口 + Go API
frontend/            主机 UI（仿 Deskreen Blueprint 风格）
webui/               浏览器查看端 → go:embed 到 web/dist
internal/            房间、采集、WebRTC、HTTP :3131
```

### 投屏链路说明

- **连接层**：与 Deskreen 一样用 **WebRTC**（SDP/ICE/DataChannel），不是「没用 WebRTC」。
- **画面层（当前）**：主机用 Win32 **PrintWindow**（应用）或显示器截图（整屏）得到位图，编码为 **JPEG**，经 DataChannel（或 WS 回退）送到浏览器 **Canvas** 绘制。
- **与原版差异**：Deskreen Electron 用 Chromium **desktopCapturer** + **VP8 媒体轨**；本项目为 Wails/Go，暂未接 libvpx，因此不是 `<video>` + `ontrack`，但应用投屏语义应对齐 **只渲窗口、不裁桌面**（已去掉 BitBlt/屏幕区域回退，避免后台时闪成整屏）。
- **限制**：窗口 **最小化**、部分 GPU/Chrome 窗口可能 PrintWindow 失败；最大化窗口看起来像整屏属正常。

## 许可说明

[Deskreen](https://github.com/pavlobu/deskreen) 为 AGPL-3.0。本项目为学习用途的独立实现，请遵守相关开源许可。
