# IoT Platform Go - ドキュメント

## 📚 ドキュメント一覧

### 1. [プロジェクト概要](./01-プロジェクト概要.md)
- プロジェクトの目的と背景
- 技術スタックの詳細
- アーキテクチャの概要
- 開発フェーズの計画

### 2. [開発環境セットアップ](./02-開発環境セットアップ.md)
- 必要なツールとバージョン
- 環境変数の設定
- データベースの初期化手順
- Docker環境の構築
- トラブルシューティング

### 3. [アーキテクチャ設計書](./03-アーキテクチャ設計書.md)
- クリーンアーキテクチャの詳細
- ディレクトリ構造の説明
- レイヤー分離の原則
- 依存関係の管理
- データフローの説明

### 4. [API仕様書](./04-API仕様書.md)
- RESTful APIの詳細仕様
- エンドポイント一覧
- リクエスト/レスポンス形式
- エラーハンドリングの仕様
- WebSocket API仕様

### 5. [データベース設計書](./05-データベース設計書.md)
- PostgreSQLテーブル定義
- InfluxDB設計
- ER図
- インデックス戦略
- マイグレーション管理

### 6. [構築STEPS](./構築STEPS.md)
- Phase 1の実装内容
- データモデル
- 設定管理
- 開発環境

## 🚀 クイックスタート

### 1. 環境セットアップ
```bash
# リポジトリのクローン
git clone <repository-url>
cd iot-platform-go

# 環境変数の設定
cp configs/env.example .env

# Dockerサービスの起動
make docker-up

# 依存関係のインストール
make deps

# アプリケーションの起動
make run
```

### 2. 動作確認
```bash
# ヘルスチェック
curl http://localhost:8080/health

# デバイスの作成
curl -X POST http://localhost:8080/api/devices \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Temperature Sensor 1",
    "type": "temperature",
    "location": "Living Room"
  }'
```

## 📋 開発ガイドライン

### コード規約
- Go言語の標準的なコーディング規約に従う
- クリーンアーキテクチャの原則を守る
- 適切なエラーハンドリングを実装する
- テストコードを必ず書く

### コミットメッセージ
```
feat: 新機能の追加
fix: バグ修正
docs: ドキュメントの更新
style: コードスタイルの修正
refactor: リファクタリング
test: テストの追加・修正
chore: その他の変更
```

### ブランチ戦略
- `main`: 本番環境用
- `develop`: 開発環境用
- `feature/*`: 機能開発用
- `hotfix/*`: 緊急修正用

## 🔧 開発ツール

### 推奨ツール
- **エディタ**: VS Code, GoLand
- **APIテスト**: Postman, curl
- **データベース管理**: pgAdmin, DBeaver
- **Docker管理**: Docker Desktop

### 便利なコマンド
```bash
# 開発用コマンド
make help          # 利用可能なコマンド一覧
make build         # アプリケーションのビルド
make run           # アプリケーションの起動
make test          # テストの実行
make fmt           # コードフォーマット
make lint          # コードリント

# Docker関連
make docker-up     # Dockerサービスの起動
make docker-down   # Dockerサービスの停止
make logs          # Dockerログの表示
```

## 📊 監視・ログ

### アプリケーション監視
- ヘルスチェックエンドポイント: `/health`
- メトリクスエンドポイント: `/metrics`（予定）
- ログレベル: DEBUG, INFO, WARN, ERROR

### データベース監視
- PostgreSQL: 接続数、クエリパフォーマンス
- InfluxDB: 書き込み速度、クエリ速度

## 🔒 セキュリティ

### 認証・認可
- JWT認証
- APIキー認証
- ロールベースアクセス制御（RBAC）

### データ保護
- HTTPS通信
- データ暗号化
- 入力検証
- SQLインジェクション対策

## 📈 パフォーマンス

### 目標値
- APIレスポンス時間: < 100ms
- 同時接続数: 1000+ デバイス
- データ処理: 10000+ メッセージ/秒
- 可用性: 99.9%

### 最適化戦略
- データベースインデックスの最適化
- キャッシュの活用
- 非同期処理の実装
- ロードバランサーの設定

## 🚀 デプロイメント

### 環境別設定
- **開発環境**: ローカルDocker
- **ステージング環境**: クラウド（AWS/GCP）
- **本番環境**: クラウド（AWS/GCP）

### CI/CD
- GitHub Actions（予定）
- 自動テスト
- 自動デプロイ
- ロールバック機能

## 📞 サポート

### 問題報告
- GitHub Issuesを使用
- 詳細な再現手順を記載
- ログファイルを添付

### 質問・相談
- GitHub Discussionsを使用
- ドキュメントの改善提案も歓迎

## 📝 更新履歴

| 日付 | バージョン | 変更内容 |
|------|------------|----------|
| 2024-01-01 | 1.0.0 | 初回リリース |
| 2024-01-15 | 1.1.0 | MQTT統合追加 |
| 2024-02-01 | 1.2.0 | WebSocket機能追加 |

---

**注意**: このドキュメントは継続的に更新されます。最新の情報はGitHubリポジトリを確認してください。 