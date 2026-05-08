# Freecode

**統合されたクロスプラットフォームAIコーディングアシスタント。**

Freecodeは、opencodeの最佳の機能と強化されたエージェント、フック、フリート管理、ワークフロー機能を1つのまとまった製品として組み合わせた、GoベースのCLIです。

## 機能

- **11の組み込みエージェント** — Sisyphus（オーケストレイター）、Hephaestus（コード生成）、Oracle、Librarian、Explore、Prometheus、Metis、Momus、Atlas、Multimodal-Looker、Sisyphus-Junior
- **60+ライフサイクルフック** — セッション、ツール、変換、Continuation、Ralph、スキルフック
- **8つのタスクカテゴリ** — visual-engineering、ultrabrain、deep、artistry、quick、writingなど
- **セッションタビング** — 分割ビューを持つ複数の並行セッション
- **フリート管理** — マルチインスタンス調整のためのhead/agent/clientモード
- **組み込みMCP** — Exa websearch、Context7 docs、Grep.app、GitHub/GitLab CLI
- **マウスサポート付きTUI** — クリック可能な要素を備えた完全な対話式インターフェース
- **セキュリティファースト** — すべてのサービスがlocalhostにバインド、テレメトリなし

## プラットフォームサポート

| プラットフォーム | ステータス | メモ |
|-----------------|-----------|------|
| FreeBSD | プライマリ | PortsのGo 1.25 |
| Linux | サポート済み | Flatpakパッケージング |
| macOS | サポート済み | Homebrew |
| IllumOS | サポート済み | tarball |
| Windows | ❌ サポート外 | |

## クイックスタート

```bash
# ビルド
go build -o freecode ./cmd/freecode

# 実行
./freecode

# またはHomebrewでインストール（macOS）
brew install freecode
```

## アーキテクチャ

- **GoベースのCLI** — 単一の静的バイナリ、ランタイム依存なし
- **Cobra CLIフレームワーク** — 標準的なGo CLIパターン
- **Bubble Tea TUI** — 合成可能なターミナルUI
- **SQLite** — 組み込み永続ストレージ
- **chi_router** — 軽量HTTP API

## 主要ディレクトリ

```
cmd/freecode/          # CLIエントリポイント
cmd/freecode-server/  # サーバーモードエントリポイント
internal/cli/          # Cobraコマンド
internal/agent/        # 11の組み込みエージェント
internal/hook/         # 60+ライフサイクルフック
internal/session/      # セッション管理、タブ
internal/ui/           # Bubble Tea TUI
internal/fleet/        # Fleet head/agent/client
internal/platform/     # プラットフォーム固有コード
```

## ポート（localhostのみ）

| サービス | ポート |
|---------|--------|
| APIサーバ | 18792 |
| MCPサーバ | 18793 |
| Web UI | 18791 |
| Fleet Head | 7842 |

## opencodeとの比較

Freecodeは、oh-my-openagentのすべての機能をネイティブに統合したopencodeのGo変換です：

| 機能 | opencode | freecode |
|------|----------|----------|
| 言語 | TypeScript | Go |
| 配布 | NPM | 静的バイナリ |
| エージェント | 7 | 11 |
| フック | ~20 | 60+ |
| Fleetモード | ❌ | ✅ |
| 組み込みMCP | ❌ | ✅ |

## セキュリティ

- **localhostのみ** — すべてのサービスが127.0.0.1と::1にバインド
- **テレメトリなし** — 分析またはトラッキングなし
- **権限システム** — エージェントごとの構成可能なツール権限
- **YOLOモード** — 確認をスキップするオプション（デフォルトでオフ）

## ドキュメント

自律エージェントのガイダンスについては[AGENTS_START_HERE.md](AGENTS_START_HERE.md)を、詳細な計画ドキュメントについては[.plan/](.plan/)ディレクトリを参照してください。

## 著者

Mark LaPointe <mark@cloudbsd.org>

すべてのコミットはMark LaPointeによって行われます。共同作成者なし、スポンサーなし。

## ライセンス

プロジェクトに無制限のライセンスが付与されます。
