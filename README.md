# Google検索サマリーアプリ

Go言語で作成されたWebアプリケーションで、Googleの検索を行い、上位5件の結果をサマリーして表示します。

## 機能

- 🔍 Google検索の実行
- 📝 検索結果の自動サマリー生成
- 🌐 日本語対応のWebインターフェース
- 📱 レスポンシブデザイン

## 使用方法

### アプリケーションの起動

```bash
go run main.go
```

または

```bash
go build -o google-search-summary
./google-search-summary
```

### アクセス

ブラウザで http://localhost:8080 にアクセスしてください。

### 検索の実行

1. トップページで検索キーワードを入力
2. 「検索」ボタンをクリック
3. 検索結果とサマリーが表示されます

## 技術仕様

- **言語**: Go 1.24+
- **依存関係**: 
  - `github.com/PuerkitoBio/goquery` - HTMLパースing
- **ポート**: 8080

## 注意事項

- Google検索の直接スクレイピングは制限される場合があります
- 制限された場合は、デモ用のモックデータが表示されます
- 本格的な運用には Google Custom Search API の使用を推奨します

## ファイル構成

```
.
├── main.go          # メインアプリケーション
├── go.mod           # Go モジュール定義
├── go.sum           # 依存関係チェックサム
└── README.md        # このファイル
```

## ライセンス

MIT License