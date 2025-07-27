# セキュリティ設定ガイド

## GitHub Dependency Graph の有効化

`actions/dependency-review-action@v4`を使用するには、GitHubのDependency graphを有効にする必要があります。

### 手順

1. **リポジトリの設定ページに移動**
   - GitHubでリポジトリを開く
   - `Settings` タブをクリック
   - 左サイドバーから `Security & analysis` を選択

2. **Dependency graph を有効化**
   - `Dependency graph` セクションを見つける
   - `Enable` ボタンをクリック
   - 確認ダイアログで `Enable dependency graph` をクリック

3. **GitHub Advanced Security の有効化（プライベートリポジトリの場合）**
   - プライベートリポジトリの場合は、GitHub Advanced Securityも必要です
   - `Code scanning` セクションで `Enable` をクリック
   - 必要に応じてGitHub Advanced Securityライセンスを購入

### 設定後の確認

設定が完了すると、以下の機能が利用可能になります：

- **Dependency graph**: リポジトリの依存関係を視覚化
- **Dependency review**: プルリクエストでの依存関係の変更を自動チェック
- **Security advisories**: 既知の脆弱性の自動検出

## 代替手段

Dependency graphが利用できない場合、以下の代替手段が実装されています：

1. **手動依存関係チェック** (`security.yml`)
   - Go modulesの検証
   - 基本的な脆弱性チェック

2. **包括的依存関係チェック** (`dependency-check.yml`)
   - `govulncheck`を使用した脆弱性スキャン
   - 依存関係の更新チェック
   - 詳細なレポート生成

## 推奨設定

セキュリティを最大限に高めるために、以下を推奨します：

1. **Dependency graph の有効化**
2. **GitHub Advanced Security の有効化**
3. **定期的なセキュリティスキャンの実行**
4. **依存関係の自動更新設定**

## トラブルシューティング

### よくあるエラー

**Error: Dependency review is not supported on this repository**

このエラーは、Dependency graphが有効になっていない場合に発生します。

**解決方法:**
1. リポジトリの設定でDependency graphを有効化
2. プライベートリポジトリの場合はGitHub Advanced Securityを有効化

### サポート

問題が解決しない場合は、以下を確認してください：

1. リポジトリの権限設定
2. GitHub Advanced Securityライセンスの状況
3. 組織のセキュリティポリシー 