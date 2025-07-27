# ブランチ保護ルール

このドキュメントでは、このプロジェクトで設定されているブランチ保護ルールについて説明します。

## 概要

重要なブランチが誤って削除されたり、不適切な変更が加えられることを防ぐために、GitHubのブランチ保護ルールを設定しています。

## 保護されているブランチ

### main/masterブランチ
- **削除禁止**: `allow_deletions: false`
- **強制プッシュ禁止**: `allow_force_pushes: false`
- **直線的な履歴要求**: `required_linear_history: true`
- **最新状態要求**: `require_up_to_date: true`
- **プルリクエストレビュー必須**: 最低1人の承認が必要
- **コードオーナーレビュー必須**: `require_code_owner_reviews: true`
- **署名されたコミット必須**: `require_signed_commits: true`
- **管理者も制限に従う**: `enforce_admins: true`

### staginブランチ
- **削除禁止**: `allow_deletions: false`
- **強制プッシュ禁止**: `allow_force_pushes: false`
- **プルリクエストレビュー必須**: 最低1人の承認が必要
- **署名されたコミット**: 任意

### developブランチ
- **削除禁止**: `allow_deletions: false`
- **強制プッシュ禁止**: `allow_force_pushes: false`
- **プルリクエストレビュー必須**: 最低1人の承認が必要
- **署名されたコミット**: 任意

## 設定ファイル

ブランチ保護ルールは以下のファイルで管理されています：

- `.github/branch-protection.yml`: ブランチ保護ルールの設定
- `.github/workflows/branch-protection.yml`: 自動適用ワークフロー

## 自動適用

ブランチ保護ルールは以下のタイミングで自動的に適用されます：

1. **設定ファイルの変更時**: `.github/branch-protection.yml`が変更された場合
2. **手動実行**: GitHub Actionsのワークフローを手動で実行
3. **定期実行**: 毎日午前9時に自動実行

## 手動での設定

GitHubのWebインターフェースから手動でブランチ保護ルールを設定する場合：

1. リポジトリの **Settings** タブに移動
2. 左サイドバーの **Branches** をクリック
3. **Add rule** または既存のルールを編集
4. 以下の設定を有効にする：
   - ✅ **Restrict deletions**
   - ✅ **Restrict pushes that create files that use the Git LFS (Large File Storage)**
   - ✅ **Require a pull request before merging**
   - ✅ **Require approvals** (1人以上)
   - ✅ **Dismiss stale PR approvals when new commits are pushed**
   - ✅ **Require review from code owners**
   - ✅ **Require status checks to pass before merging**
   - ✅ **Require branches to be up to date before merging**
   - ✅ **Require signed commits**
   - ✅ **Require linear history**
   - ✅ **Restrict force pushes**
   - ✅ **Do not allow bypassing the above settings**

## トラブルシューティング

### ブランチが削除できない場合
保護されたブランチを削除しようとすると、GitHubはエラーを表示します。これは意図された動作です。

### 強制プッシュができない場合
保護されたブランチへの強制プッシュは禁止されています。代わりに：
1. 新しいブランチを作成
2. 変更をコミット
3. プルリクエストを作成
4. レビュー後にマージ

### ワークフローが失敗する場合
1. リポジトリの設定で **Actions** の権限が有効になっているか確認
2. ワークフローファイルの構文エラーがないか確認
3. GitHub Actionsのログを確認して詳細なエラーメッセージを確認

## セキュリティ上の注意

- 管理者も制限に従う設定により、リポジトリオーナーでも保護されたブランチを削除できません
- 緊急時は、GitHubのWebインターフェースから一時的に保護を無効にできます
- 保護ルールの変更は、必ずプルリクエストを通して行ってください

## 関連リンク

- [GitHub Branch Protection Rules](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/defining-the-mergeability-of-pull-requests/managing-a-branch-protection-rule)
- [GitHub Actions Documentation](https://docs.github.com/en/actions) 