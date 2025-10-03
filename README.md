# Achievement Management System

Go言語で開発された達成目録管理アプリケーション

## プロジェクト構造

```
achievement-management/
├── cmd/
│   ├── api/           # REST APIサーバー
│   └── cli/           # コマンドラインツール
├── internal/
│   ├── models/        # データモデル
│   ├── services/      # ビジネスロジック層
│   ├── repository/    # データアクセス層
│   ├── handlers/      # HTTPハンドラー
│   ├── config/        # 設定管理
│   └── errors/        # エラーハンドリング
├── go.mod
└── README.md
```

## データモデル

- **Achievement**: 達成目録
- **Reward**: 報酬
- **CurrentPoints**: 現在のポイント
- **RewardHistory**: 報酬獲得履歴

## 開発環境

- Go 1.24+
- AWS DynamoDB (Local/Cloud)
- Docker (オプション)

## ビルドとデプロイメント

### 前提条件

- Go 1.24以上がインストールされていること
- Git がインストールされていること（バージョン情報取得のため）
- Docker がインストールされていること（Dockerビルドを使用する場合）

### ビルド方法

#### 1. Makefileを使用したビルド

```bash
# ヘルプを表示
make help

# 現在のプラットフォーム用にビルド
make build

# 全プラットフォーム用にクロスコンパイル
make build-all

# テスト実行
make test

# 配布パッケージ作成
make package

# Dockerイメージ作成
make docker-build
```

#### 2. ビルドスクリプトを使用したビルド

```bash
# ヘルプを表示
./scripts/build.sh help

# 完全なビルドプロセス（テスト + ビルド + パッケージ）
./scripts/build.sh all

# テストなしでビルドのみ
./scripts/build.sh build-only

# 配布パッケージ作成
./scripts/build.sh package
```

### 対応プラットフォーム

- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

### バイナリファイル

ビルド後、以下のバイナリファイルが生成されます：

- `achievement-api`: REST APIサーバー
- `achievement-app`: コマンドラインツール

### 配布パッケージ

`make package` または `./scripts/build.sh package` を実行すると、`dist/` ディレクトリに以下のファイルが作成されます：

- `achievement-app-{version}-linux-amd64.tar.gz`
- `achievement-app-{version}-linux-arm64.tar.gz`
- `achievement-app-{version}-darwin-amd64.tar.gz`
- `achievement-app-{version}-darwin-arm64.tar.gz`
- `achievement-app-{version}-windows-amd64.zip`

各パッケージには以下が含まれます：
- バイナリファイル（API サーバーとCLI）
- 設定ファイル（config/）
- README.md
- .env.example

### Dockerを使用したデプロイメント

```bash
# Dockerイメージをビルド
make docker-build

# Dockerコンテナを実行
make docker-run

# または直接実行
docker run -p 8080:8080 achievement-app:latest
```

### 環境変数

アプリケーションは以下の環境変数で設定できます：

```bash
# AWS設定
AWS_REGION=ap-northeast-1
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
DYNAMODB_ENDPOINT=http://localhost:8000  # ローカル開発用

# サーバー設定
SERVER_PORT=8080
LOG_LEVEL=info
ENVIRONMENT=development
```

## 使用方法

### APIサーバー起動

```bash
# 開発環境で起動
go run cmd/api/main.go

# ビルド済みバイナリで起動
./build/achievement-api

# Dockerで起動
docker run -p 8080:8080 achievement-app:latest
```

### CLIツール使用

```bash
# 開発環境で実行
go run cmd/cli/main.go [command]

# ビルド済みバイナリで実行
./build/achievement-app [command]

# バージョン確認
./build/achievement-app --version

# ヘルプ表示
./build/achievement-app --help
```

### 開発環境セットアップ

```bash
# 依存関係のダウンロード
make deps

# 開発環境のセットアップ（linter等のインストール）
make setup

# ローカルDynamoDB起動（Docker必要）
make setup-dynamodb

# フォーマット
make fmt

# リント
make lint

# テスト実行
make test

# カバレッジ付きテスト
make test-coverage
```

## API エンドポイント

### ヘルスチェック

```bash
# ヘルスチェック
curl -X GET http://localhost:8080/health
```

### 達成目録管理

```bash
# 達成目録作成
curl -X POST http://localhost:8080/api/achievements \
  -H "Content-Type: application/json" \
  -d '{
    "title": "初回ログイン",
    "description": "アプリに初回ログインした",
    "point": 10
  }'

# 達成目録一覧取得
curl -X GET http://localhost:8080/api/achievements

# 達成目録詳細取得
curl -X GET http://localhost:8080/api/achievements/{achievement_id}

# 達成目録更新
curl -X PUT http://localhost:8080/api/achievements/{achievement_id} \
  -H "Content-Type: application/json" \
  -d '{
    "title": "初回ログイン（更新）",
    "description": "アプリに初回ログインした（説明更新）",
    "point": 15
  }'

# 達成目録削除
curl -X DELETE http://localhost:8080/api/achievements/{achievement_id}
```

### 報酬管理

```bash
# 報酬作成
curl -X POST http://localhost:8080/api/rewards \
  -H "Content-Type: application/json" \
  -d '{
    "title": "コーヒー券",
    "description": "スターバックスのコーヒー券",
    "point": 100
  }'

# 報酬一覧取得
curl -X GET http://localhost:8080/api/rewards

# 報酬詳細取得
curl -X GET http://localhost:8080/api/rewards/{reward_id}

# 報酬更新
curl -X PUT http://localhost:8080/api/rewards/{reward_id} \
  -H "Content-Type: application/json" \
  -d '{
    "title": "コーヒー券（更新）",
    "description": "スターバックスのコーヒー券（説明更新）",
    "point": 120
  }'

# 報酬削除
curl -X DELETE http://localhost:8080/api/rewards/{reward_id}

# 報酬獲得
curl -X POST http://localhost:8080/api/rewards/{reward_id}/redeem
```

### ポイント管理

```bash
# 現在のポイント取得
curl -X GET http://localhost:8080/api/points/current

# ポイント集計取得
curl -X GET http://localhost:8080/api/points/aggregate

# 報酬獲得履歴取得
curl -X GET http://localhost:8080/api/points/history
```

## 要件

このプロジェクトは以下の要件を満たします：
- 要件 1.1: 達成目録の作成・編集機能
- 要件 2.1: 報酬の作成・編集機能  
- 要件 3.1: ポイント自動更新機能
- 要件 4.1: 報酬獲得機能
