# 設計書

## 概要

「達成目録」管理アプリケーションは、Go言語で開発されるRESTful APIサーバーとコマンドラインツールを提供するシステムです。AWS DynamoDBをデータストアとして使用し、達成目録、報酬、ポイント管理機能を提供します。

## アーキテクチャ

### システム構成

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   CLI Tool      │    │   REST API      │    │   AWS DynamoDB  │
│                 │    │                 │    │                 │
│ - 達成目録管理   │    │ - HTTP Server   │    │ - achievements  │
│ - 報酬管理      │    │ - JSON API      │    │ - rewards       │
│ - ポイント集計   │    │ - Validation    │    │ - current_points│
│                 │    │                 │    │ - reward_history│
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │  Service Layer  │
                    │                 │
                    │ - Business Logic│
                    │ - Data Access   │
                    │ - Validation    │
                    └─────────────────┘
```

### レイヤー構成

1. **プレゼンテーション層**
   - REST API ハンドラー
   - CLI コマンドハンドラー
   - リクエスト/レスポンス変換

2. **ビジネスロジック層**
   - 達成目録管理サービス
   - 報酬管理サービス
   - ポイント管理サービス
   - データ検証

3. **データアクセス層**
   - DynamoDB リポジトリ
   - データモデル変換
   - エラーハンドリング

## コンポーネントとインターフェース

### 主要コンポーネント

#### 1. データモデル

```go
// Achievement 達成目録
type Achievement struct {
    ID          string    `json:"id" dynamodbav:"id"`
    Title       string    `json:"title" dynamodbav:"title"`
    Description string    `json:"description" dynamodbav:"description"`
    Point       int       `json:"point" dynamodbav:"point"`
    CreatedAt   time.Time `json:"created_at" dynamodbav:"created_at"`
}

// Reward 報酬
type Reward struct {
    ID          string    `json:"id" dynamodbav:"id"`
    Title       string    `json:"title" dynamodbav:"title"`
    Description string    `json:"description" dynamodbav:"description"`
    Point       int       `json:"point" dynamodbav:"point"`
    CreatedAt   time.Time `json:"created_at" dynamodbav:"created_at"`
}

// CurrentPoints 現在のポイント
type CurrentPoints struct {
    ID        string    `json:"id" dynamodbav:"id"` // 固定値 "current"
    Point     int       `json:"point" dynamodbav:"point"`
    UpdatedAt time.Time `json:"updated_at" dynamodbav:"updated_at"`
}

// RewardHistory 報酬獲得履歴
type RewardHistory struct {
    ID         string    `json:"id" dynamodbav:"id"`
    RewardID   string    `json:"reward_id" dynamodbav:"reward_id"`
    RewardTitle string   `json:"reward_title" dynamodbav:"reward_title"`
    PointCost  int       `json:"point_cost" dynamodbav:"point_cost"`
    RedeemedAt time.Time `json:"redeemed_at" dynamodbav:"redeemed_at"`
}
```

#### 2. サービスインターフェース

```go
// AchievementService 達成目録サービス
type AchievementService interface {
    Create(achievement *Achievement) error
    Update(id string, achievement *Achievement) error
    GetByID(id string) (*Achievement, error)
    List() ([]*Achievement, error)
    Delete(id string) error
}

// RewardService 報酬サービス
type RewardService interface {
    Create(reward *Reward) error
    Update(id string, reward *Reward) error
    GetByID(id string) (*Reward, error)
    List() ([]*Reward, error)
    Delete(id string) error
    Redeem(rewardID string) error
}

// PointService ポイントサービス
type PointService interface {
    GetCurrentPoints() (*CurrentPoints, error)
    AddPoints(points int) error
    SubtractPoints(points int) error
    AggregatePoints() (*PointSummary, error)
}
```

#### 3. リポジトリインターフェース

```go
// Repository DynamoDB操作の抽象化
type Repository interface {
    PutItem(tableName string, item interface{}) error
    GetItem(tableName string, key map[string]interface{}, result interface{}) error
    UpdateItem(tableName string, key map[string]interface{}, updateExpression string, expressionAttributeValues map[string]interface{}) error
    Scan(tableName string, result interface{}) error
    DeleteItem(tableName string, key map[string]interface{}) error
    TransactWrite(items []TransactWriteItem) error
}
```

## データモデル

### DynamoDBテーブル設計

#### 1. achievements テーブル
- **パーティションキー**: `id` (String)
- **属性**:
  - `title` (String) - 達成目録のタイトル
  - `description` (String) - 達成目録の説明
  - `point` (Number) - 獲得ポイント
  - `created_at` (String) - 作成日時 (ISO 8601形式)

#### 2. rewards テーブル
- **パーティションキー**: `id` (String)
- **属性**:
  - `title` (String) - 報酬のタイトル
  - `description` (String) - 報酬の説明
  - `point` (Number) - 必要ポイント
  - `created_at` (String) - 作成日時 (ISO 8601形式)

#### 3. current_points テーブル
- **パーティションキー**: `id` (String) - 固定値 "current"
- **属性**:
  - `point` (Number) - 現在のポイント
  - `updated_at` (String) - 更新日時 (ISO 8601形式)

#### 4. reward_history テーブル
- **パーティションキー**: `id` (String) - ULID形式のユニークID
- **属性**:
  - `reward_id` (String) - 獲得した報酬のID
  - `reward_title` (String) - 獲得した報酬のタイトル
  - `point_cost` (Number) - 消費したポイント
  - `redeemed_at` (String) - 獲得日時 (ISO 8601形式)
- **GSI**: `redeemed_at-index` (時系列順での取得用)

### データ整合性

- **ポイント更新**: DynamoDB Transactionsを使用してアトミックな操作を保証
- **報酬獲得**: current_pointsの減算とreward_historyの追加を単一トランザクションで実行
- **達成目録追加**: current_pointsの加算を確実に実行

## エラーハンドリング

### エラー分類

1. **バリデーションエラー**
   - 必須フィールドの欠如
   - データ型の不一致
   - 値の範囲外

2. **ビジネスロジックエラー**
   - ポイント不足
   - 存在しないリソースへのアクセス
   - 重複データの作成

3. **システムエラー**
   - DynamoDB接続エラー
   - AWS認証エラー
   - 内部サーバーエラー

### エラーレスポンス形式

```go
type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message"`
    Code    int    `json:"code"`
}
```

## テスト戦略

### テストレベル

1. **ユニットテスト**
   - サービス層のビジネスロジック
   - データ変換ロジック
   - バリデーション機能

2. **統合テスト**
   - DynamoDB操作
   - API エンドポイント
   - CLI コマンド

3. **E2Eテスト**
   - 完全なワークフロー
   - エラーシナリオ
   - パフォーマンステスト

### テストデータ管理

- DynamoDB Local を使用した開発環境
- テストデータのセットアップ/クリーンアップ
- モックを使用した外部依存関係の分離

## API設計

### RESTエンドポイント

#### 達成目録管理
- `POST /api/achievements` - 達成目録作成
- `GET /api/achievements` - 達成目録一覧取得
- `GET /api/achievements/{id}` - 達成目録詳細取得
- `PUT /api/achievements/{id}` - 達成目録更新
- `DELETE /api/achievements/{id}` - 達成目録削除

#### 報酬管理
- `POST /api/rewards` - 報酬作成
- `GET /api/rewards` - 報酬一覧取得
- `GET /api/rewards/{id}` - 報酬詳細取得
- `PUT /api/rewards/{id}` - 報酬更新
- `DELETE /api/rewards/{id}` - 報酬削除
- `POST /api/rewards/{id}/redeem` - 報酬獲得

#### ポイント管理
- `GET /api/points/current` - 現在のポイント取得
- `GET /api/points/aggregate` - ポイント集計
- `GET /api/points/history` - 報酬獲得履歴取得

### CLI コマンド設計

```bash
# 達成目録管理
achievement-app achievement create --title "タイトル" --description "説明" --point 100
achievement-app achievement list
achievement-app achievement update --id "ID" --title "新タイトル"
achievement-app achievement delete --id "ID"

# 報酬管理
achievement-app reward create --title "タイトル" --description "説明" --point 50
achievement-app reward list
achievement-app reward update --id "ID" --title "新タイトル"
achievement-app reward redeem --id "ID"
achievement-app reward delete --id "ID"

# ポイント管理
achievement-app points current
achievement-app points aggregate
achievement-app points history
```

## セキュリティ考慮事項

1. **認証・認可**
   - AWS IAM ロールベースのアクセス制御
   - DynamoDB リソースレベルの権限設定

2. **データ保護**
   - DynamoDB暗号化の有効化
   - 機密データのマスキング

3. **入力検証**
   - SQLインジェクション対策（NoSQLインジェクション）
   - XSS対策
   - データサイズ制限

## パフォーマンス考慮事項

1. **DynamoDB最適化**
   - 適切なパーティションキー設計
   - 読み取り/書き込みキャパシティの設定
   - インデックス戦略

2. **アプリケーション最適化**
   - 接続プールの使用
   - キャッシュ戦略
   - バッチ処理の活用

## 運用考慮事項

1. **ログ管理**
   - 構造化ログの出力
   - エラーレベルの適切な設定
   - 監査ログの記録

2. **監視**
   - DynamoDB メトリクスの監視
   - アプリケーションメトリクスの収集
   - アラート設定

3. **バックアップ**
   - DynamoDB Point-in-Time Recovery
   - 定期的なバックアップ戦略