## Phase1. 基盤構築
### 簡単な実装内容
Go + Gin でのREST API基盤
PostgreSQL でのデバイス管理
基本的なCRUD操作

#### データモデル
Device: IoTデバイスの基本情報（ID、名前、タイプ、場所、ステータスなど）
DeviceData: センサーデータ（デバイスID、タイムスタンプ、データ）
CreateDeviceRequest: デバイス作成時のリクエスト
UpdateDeviceRequest: デバイス更新時のリクエスト

#### 設定管理
環境変数から設定を読み込み
サーバー設定（ポート、ホスト）
データベース設定（PostgreSQL接続情報）
MQTT設定（ブローカーURL、認証情報）
JWT設定（認証用）

#### データベース接続
PostgreSQLへの接続
テーブル自動作成
devicesテーブル：デバイス情報
device_dataテーブル：センサーデータ
インデックス作成（パフォーマンス向上

#### デバイス管理リポジトリ
CRUD操作：
Create: デバイス作成
Read: デバイス取得（単体・全件）
Update: デバイス更新
Delete: デバイス削除
ステータス管理: デバイスのオンライン/オフライン状態更新

#### APIハンドラー
RESTful API：
POST /api/devices - デバイス作成
GET /api/devices - 全デバイス取得
GET /api/devices/:id - 特定デバイス取得
PUT /api/devices/:id - デバイス更新
DELETE /api/devices/:id - デバイス削除
GET /api/devices/:id/status - デバイスステータス取得
エラーハンドリング: 適切なHTTPステータスコードとエラーメッセージ

#### メインサーバー
GinフレームワークでWebサーバー起動
ルーティング設定
ミドルウェア（ログ、CORS、リカバリー）
ヘルスチェックエンドポイント（/health）

#### 開発環境
PostgreSQL: データベース
Mosquitto: MQTTブローカー（Phase 2で使用）
Grafana: 可視化ツール（オプション）

#### 留意点
クリーンアーキテクチャ: レイヤー分離（API → Repository → Database）
エラーハンドリング: 適切なHTTPステータスコード
開発効率: Makefile、Docker Compose