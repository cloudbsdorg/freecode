# Freecode

**统一的跨平台AI编程助手。**

Freecode是一个基于Go的CLI，它将opencode的最佳功能与增强的代理、钩子、车队管理和工作流功能结合在一起——作为单一统一产品。

## 功能特点

- **11个内置代理** — Sisyphus（编排器）、Hephaestus（代码生成）、Oracle、Librarian、Explore、Prometheus、Metis、Momus、Atlas、Multimodal-Looker、Sisyphus-Junior
- **60+生命周期钩子** — 会话、工具、转换、延续、Ralph、技能钩子
- **8个任务类别** — visual-engineering、ultrabrain、deep、artistry、quick、writing等
- **会话标签页** — 具有分屏视图的多个并发会话
- **车队管理** — Head/agent/client模式，用于多实例协调
- **内置MCP** — Exa网络搜索、Context7文档、Grep.app、GitHub/GitLab CLI
- **支持鼠标的TUI** — 完整的交互式界面，可点击元素
- **安全优先** — 所有服务仅绑定到本地主机，无遥测

## 平台支持

| 平台 | 状态 | 备注 |
|------|------|------|
| FreeBSD | 主要平台 | Ports中的Go 1.25 |
| Linux | 支持 | Flatpak打包 |
| macOS | 支持 | Homebrew |
| IllumOS | 支持 | tarball |
| Windows | ❌ 不支持 | |

## 快速开始

```bash
# 构建
go build -o freecode ./cmd/freecode

# 运行
./freecode

# 或通过Homebrew安装（macOS）
brew install freecode
```

## 架构

- **基于Go的CLI** — 单一静态二进制文件，无运行时依赖
- **Cobra CLI框架** — 标准Go CLI模式
- **Bubble Tea TUI** — 可组合的终端UI
- **SQLite** — 嵌入式持久存储
- **chi路由器** — 轻量级HTTP API

## 关键目录

```
cmd/freecode/          # CLI入口点
cmd/freecode-server/  # 服务器模式入口点
internal/cli/          # Cobra命令
internal/agent/        # 11个内置代理
internal/hook/         # 60+生命周期钩子
internal/session/      # 会话管理、标签页
internal/ui/           # Bubble Tea TUI
internal/fleet/        # Fleet head/agent/client
internal/platform/     # 平台特定代码
```

## 端口（仅本地主机）

| 服务 | 端口 |
|------|------|
| API服务器 | 18792 |
| MCP服务器 | 18793 |
| Web UI | 18791 |
| Fleet Head | 7842 |

## 与opencode的比较

Freecode是将opencode转换为Go版本，并将oh-my-openagent的所有功能原生集成：

| 功能 | opencode | freecode |
|------|----------|----------|
| 语言 | TypeScript | Go |
| 分发 | NPM | 静态二进制 |
| 代理 | 7 | 11 |
| 钩子 | ~20 | 60+ |
| Fleet模式 | ❌ | ✅ |
| 内置MCP | ❌ | ✅ |

## 安全性

- **仅本地主机** — 所有服务绑定到127.0.0.1和::1
- **无遥测** — 零分析或跟踪
- **权限系统** — 每个代理可配置的工具权限
- **YOLO模式** — 可选择跳过确认（默认关闭）

## 文档

有关自主代理指导，请参阅[AGENTS_START_HERE.md](AGENTS_START_HERE.md)，或浏览[.plan/](.plan/)目录获取详细规划文档。

## 作者

Mark LaPointe <mark@cloudbsd.org>

所有提交均由Mark LaPointe完成。无共同作者，无赞助。

## 许可证

授予项目无限许可证。
