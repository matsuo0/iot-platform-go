# API仕様書

## 概要

IoT Platform GoのRESTful API仕様です。デバイス管理、データ収集、認証などの機能を提供します。

## 基本情報

- **Base URL**: `http://localhost:8080`
- **Content-Type**: `application/json`
- **認証**: JWT Bearer Token（一部エンドポイントを除く）

## 共通レスポンス形式

### 成功レスポンス

```json
{
  "success": true,
  "data": {
    // レスポンスデータ
  },
  "message": "操作が成功しました"
}
```

### エラーレスポンス

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "バリデーションエラーが発生しました",
    "details": {
      "field": "name",
      "reason": "名前は必須です"
    }
  }
}
```

### HTTPステータスコード

| コード | 説明 |
|--------|------|
| 200 | OK - 成功 |
| 201 | Created - 作成成功 |
| 400 | Bad Request - リクエストエラー |
| 401 | Unauthorized - 認証エラー |
| 403 | Forbidden - 認可エラー |
| 404 | Not Found - リソースが見つかりません |
| 422 | Unprocessable Entity - バリデーションエラー |
| 500 | Internal Server Error - サーバーエラー |

## エンドポイント一覧

### ヘルスチェック

#### GET /health

システムの健全性を確認します。

**リクエスト**
```http
GET /health
```

**レスポンス**
```json
{
  "status": "ok",
  "timestamp": "2024-01-01T00:00:00Z",
  "version": "1.0.0",
  "uptime": "2h30m15s"
}
```

### デバイス管理

#### GET /api/devices

デバイス一覧を取得します。

**リクエスト**
```http
GET /api/devices?page=1&limit=10&type=temperature&status=online
```

**クエリパラメータ**
| パラメータ | 型 | 必須 | 説明 |
|-----------|----|------|------|
| page | integer | 否 | ページ番号（デフォルト: 1） |
| limit | integer | 否 | 1ページあたりの件数（デフォルト: 10, 最大: 100） |
| type | string | 否 | デバイスタイプでフィルタ |
| status | string | 否 | ステータスでフィルタ（online/offline） |
| location | string | 否 | 場所でフィルタ |

**レスポンス**
```json
{
  "success": true,
  "data": {
    "devices": [
      {
        "id": "dev-001",
        "name": "Temperature Sensor 1",
        "type": "temperature",
        "location": "Living Room",
        "status": "online",
        "last_seen": "2024-01-01T12:00:00Z",
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T12:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 25,
      "total_pages": 3
    }
  }
}
```

#### POST /api/devices

新しいデバイスを作成します。

**リクエスト**
```http
POST /api/devices
Content-Type: application/json

{
  "name": "Temperature Sensor 1",
  "type": "temperature",
  "location": "Living Room",
  "description": "リビングルームの温度センサー",
  "metadata": {
    "manufacturer": "SensorCorp",
    "model": "TC-100",
    "firmware_version": "1.2.3"
  }
}
```

**リクエストボディ**
| フィールド | 型 | 必須 | 説明 |
|-----------|----|------|------|
| name | string | 是 | デバイス名（1-100文字） |
| type | string | 是 | デバイスタイプ（temperature, humidity, pressure等） |
| location | string | 是 | 設置場所（1-200文字） |
| description | string | 否 | 説明（最大500文字） |
| metadata | object | 否 | 追加メタデータ |

**レスポンス**
```json
{
  "success": true,
  "data": {
    "id": "dev-001",
    "name": "Temperature Sensor 1",
    "type": "temperature",
    "location": "Living Room",
    "status": "offline",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "message": "デバイスが正常に作成されました"
}
```

#### GET /api/devices/{id}

特定のデバイス情報を取得します。

**リクエスト**
```http
GET /api/devices/dev-001
```

**パスパラメータ**
| パラメータ | 型 | 必須 | 説明 |
|-----------|----|------|------|
| id | string | 是 | デバイスID |

**レスポンス**
```json
{
  "success": true,
  "data": {
    "id": "dev-001",
    "name": "Temperature Sensor 1",
    "type": "temperature",
    "location": "Living Room",
    "status": "online",
    "description": "リビングルームの温度センサー",
    "metadata": {
      "manufacturer": "SensorCorp",
      "model": "TC-100",
      "firmware_version": "1.2.3"
    },
    "last_seen": "2024-01-01T12:00:00Z",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

#### PUT /api/devices/{id}

デバイス情報を更新します。

**リクエスト**
```http
PUT /api/devices/dev-001
Content-Type: application/json

{
  "name": "Updated Temperature Sensor",
  "location": "Kitchen",
  "description": "キッチンの温度センサーに更新"
}
```

**リクエストボディ**
| フィールド | 型 | 必須 | 説明 |
|-----------|----|------|------|
| name | string | 否 | デバイス名（1-100文字） |
| location | string | 否 | 設置場所（1-200文字） |
| description | string | 否 | 説明（最大500文字） |
| metadata | object | 否 | 追加メタデータ |

**レスポンス**
```json
{
  "success": true,
  "data": {
    "id": "dev-001",
    "name": "Updated Temperature Sensor",
    "type": "temperature",
    "location": "Kitchen",
    "status": "online",
    "updated_at": "2024-01-01T13:00:00Z"
  },
  "message": "デバイスが正常に更新されました"
}
```

#### DELETE /api/devices/{id}

デバイスを削除します。

**リクエスト**
```http
DELETE /api/devices/dev-001
```

**レスポンス**
```json
{
  "success": true,
  "message": "デバイスが正常に削除されました"
}
```

#### GET /api/devices/{id}/status

デバイスのステータス情報を取得します。

**リクエスト**
```http
GET /api/devices/dev-001/status
```

**レスポンス**
```json
{
  "success": true,
  "data": {
    "device_id": "dev-001",
    "status": "online",
    "last_seen": "2024-01-01T12:00:00Z",
    "uptime": "2h30m15s",
    "connection_quality": "excellent",
    "battery_level": 85,
    "signal_strength": -45
  }
}
```

### デバイスデータ

#### GET /api/devices/{id}/data

デバイスのセンサーデータを取得します。

**リクエスト**
```http
GET /api/devices/dev-001/data?start=2024-01-01T00:00:00Z&end=2024-01-01T23:59:59Z&limit=100
```

**クエリパラメータ**
| パラメータ | 型 | 必須 | 説明 |
|-----------|----|------|------|
| start | string | 否 | 開始時刻（ISO 8601形式） |
| end | string | 否 | 終了時刻（ISO 8601形式） |
| limit | integer | 否 | 取得件数（デフォルト: 100, 最大: 1000） |
| aggregation | string | 否 | 集約方法（raw, avg, min, max, sum） |
| interval | string | 否 | 集約間隔（1m, 5m, 1h, 1d） |

**レスポンス**
```json
{
  "success": true,
  "data": {
    "device_id": "dev-001",
    "data_points": [
      {
        "timestamp": "2024-01-01T12:00:00Z",
        "temperature": 23.5,
        "humidity": 45.2,
        "pressure": 1013.25
      }
    ],
    "summary": {
      "count": 100,
      "avg_temperature": 23.2,
      "min_temperature": 22.1,
      "max_temperature": 24.8
    }
  }
}
```

#### POST /api/devices/{id}/data

デバイスからセンサーデータを受信します。

**リクエスト**
```http
POST /api/devices/dev-001/data
Content-Type: application/json

{
  "timestamp": "2024-01-01T12:00:00Z",
  "temperature": 23.5,
  "humidity": 45.2,
  "pressure": 1013.25,
  "battery_level": 85
}
```

**リクエストボディ**
| フィールド | 型 | 必須 | 説明 |
|-----------|----|------|------|
| timestamp | string | 否 | データ時刻（ISO 8601形式、デフォルト: 現在時刻） |
| temperature | number | 否 | 温度（摂氏） |
| humidity | number | 否 | 湿度（%） |
| pressure | number | 否 | 気圧（hPa） |
| battery_level | number | 否 | バッテリーレベル（%） |

**レスポンス**
```json
{
  "success": true,
  "data": {
    "id": "data-001",
    "device_id": "dev-001",
    "timestamp": "2024-01-01T12:00:00Z",
    "created_at": "2024-01-01T12:00:01Z"
  },
  "message": "データが正常に保存されました"
}
```

### 認証

#### POST /api/auth/login

ユーザーログインを行います。

**リクエスト**
```http
POST /api/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "password123"
}
```

**リクエストボディ**
| フィールド | 型 | 必須 | 説明 |
|-----------|----|------|------|
| username | string | 是 | ユーザー名 |
| password | string | 是 | パスワード |

**レスポンス**
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "refresh_token_here",
    "expires_in": 3600,
    "user": {
      "id": "user-001",
      "username": "admin",
      "email": "admin@example.com",
      "role": "admin"
    }
  }
}
```

#### POST /api/auth/refresh

アクセストークンを更新します。

**リクエスト**
```http
POST /api/auth/refresh
Content-Type: application/json

{
  "refresh_token": "refresh_token_here"
}
```

**レスポンス**
```json
{
  "success": true,
  "data": {
    "token": "new_access_token_here",
    "expires_in": 3600
  }
}
```

#### POST /api/auth/logout

ログアウトを行います。

**リクエスト**
```http
POST /api/auth/logout
Authorization: Bearer <access_token>
```

**レスポンス**
```json
{
  "success": true,
  "message": "正常にログアウトしました"
}
```

## エラーコード

### バリデーションエラー

| コード | 説明 |
|--------|------|
| VALIDATION_ERROR | バリデーションエラー |
| REQUIRED_FIELD | 必須フィールドが不足 |
| INVALID_FORMAT | 形式が不正 |
| FIELD_TOO_LONG | フィールドが長すぎる |
| FIELD_TOO_SHORT | フィールドが短すぎる |

### ビジネスロジックエラー

| コード | 説明 |
|--------|------|
| DEVICE_NOT_FOUND | デバイスが見つかりません |
| DEVICE_ALREADY_EXISTS | デバイスが既に存在します |
| DEVICE_OFFLINE | デバイスがオフラインです |
| INSUFFICIENT_PERMISSIONS | 権限が不足しています |

### システムエラー

| コード | 説明 |
|--------|------|
| INTERNAL_ERROR | 内部サーバーエラー |
| DATABASE_ERROR | データベースエラー |
| EXTERNAL_SERVICE_ERROR | 外部サービスエラー |

## レート制限

- **認証エンドポイント**: 5回/分
- **データ送信エンドポイント**: 1000回/分
- **その他のエンドポイント**: 100回/分

## WebSocket API

### 接続

```javascript
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onopen = function() {
  console.log('WebSocket接続が確立されました');
  
  // デバイスデータの購読
  ws.send(JSON.stringify({
    type: 'subscribe',
    device_id: 'dev-001'
  }));
};

ws.onmessage = function(event) {
  const data = JSON.parse(event.data);
  console.log('受信したデータ:', data);
};
```

### メッセージ形式

#### 購読メッセージ
```json
{
  "type": "subscribe",
  "device_id": "dev-001"
}
```

#### リアルタイムデータ
```json
{
  "type": "device_data",
  "device_id": "dev-001",
  "data": {
    "timestamp": "2024-01-01T12:00:00Z",
    "temperature": 23.5,
    "humidity": 45.2
  }
}
```

#### デバイスステータス更新
```json
{
  "type": "device_status",
  "device_id": "dev-001",
  "status": "online",
  "last_seen": "2024-01-01T12:00:00Z"
}
```

## 使用例

### cURLでの使用例

```bash
# デバイス一覧の取得
curl -X GET "http://localhost:8080/api/devices" \
  -H "Authorization: Bearer <token>"

# デバイスの作成
curl -X POST "http://localhost:8080/api/devices" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "name": "Temperature Sensor",
    "type": "temperature",
    "location": "Living Room"
  }'

# データの送信
curl -X POST "http://localhost:8080/api/devices/dev-001/data" \
  -H "Content-Type: application/json" \
  -d '{
    "temperature": 23.5,
    "humidity": 45.2
  }'
```

### JavaScriptでの使用例

```javascript
// デバイス一覧の取得
async function getDevices() {
  const response = await fetch('/api/devices', {
    headers: {
      'Authorization': 'Bearer ' + token
    }
  });
  const data = await response.json();
  return data.data.devices;
}

// デバイスの作成
async function createDevice(deviceData) {
  const response = await fetch('/api/devices', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' + token
    },
    body: JSON.stringify(deviceData)
  });
  return await response.json();
}
``` 