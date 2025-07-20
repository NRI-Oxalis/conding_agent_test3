package main

import (
	"fmt"
	"strings"
	"time"
)

// Core data structures for sales closing support
type ContractInfo struct {
	CustomerName    string `json:"customer_name"`
	ContractType    string `json:"contract_type"`
	ProductService  string `json:"product_service"`
	ContractAmount  string `json:"contract_amount"`
	ContractDate    string `json:"contract_date"`
	SalesPersonName string `json:"sales_person_name"`
	Notes          string `json:"notes"`
}

type DocumentTemplate struct {
	Type        string `json:"type"`
	Title       string `json:"title"`
	Template    string `json:"template"`
	Fields      []string `json:"fields"`
	Rules       []string `json:"rules"`
}

type APIConfig struct {
	MattermostURL    string `json:"mattermost_url"`
	MattermostToken  string `json:"mattermost_token"`
	JiraURL          string `json:"jira_url"`
	JiraUser         string `json:"jira_user"`
	JiraPassword     string `json:"jira_password"`
	ConfluenceURL    string `json:"confluence_url"`
	ConfluenceUser   string `json:"confluence_user"`
	ConfluencePassword string `json:"confluence_password"`
}

type SalesClosingCore struct {
	config          *APIConfig
	currentContract *ContractInfo
}

// Core business logic functions
func (app *SalesClosingCore) generateContractCheckResult(contract *ContractInfo) string {
	result := fmt.Sprintf(`# 契約情報確認結果

## 📋 入力された契約情報
- **顧客名**: %s
- **契約種別**: %s  
- **製品/サービス**: %s
- **契約金額**: %s
- **契約日**: %s
- **営業担当者**: %s

## ✅ 確認すべき情報リスト

### 必須確認項目
- [ ] 顧客の信用情報確認（与信チェック）
- [ ] 契約金額の承認フロー完了確認
- [ ] 製品/サービスの在庫・提供体制確認
- [ ] 契約書の法務チェック完了確認
- [ ] 顧客の決済方法・支払い条件確認

### %s 特有の確認項目
- [ ] 継続契約の場合：前回契約の履行状況確認
- [ ] 新規契約の場合：初回導入支援体制の確認
- [ ] 追加契約の場合：既存システムとの連携確認

## 🔗 関連Confluenceリンク
- [契約手続きマニュアル](https://confluence.example.com/contracts)
- [%s向けガイドライン](https://confluence.example.com/contract-types)
- [承認フロー一覧](https://confluence.example.com/approval-flow)

## 💡 Mattermostからのヒント

### 過去の類似案件からの学び
- **顧客名**: 類似規模の案件では、事前の技術要件確認が重要でした
- **契約金額**: この金額帯では部長承認が必要です
- **製品/サービス**: 導入時のサポート体制についてお客様から質問が多い傾向があります

### 注意点
- ⚠️ 契約日が月末近くの場合、請求処理のタイミングにご注意ください
- ⚠️ この製品は導入に2週間程度かかるため、納期の調整が必要な場合があります

## 📞 次のアクション
1. 上記チェックリストの各項目を確認
2. 不明点がある場合は関連部署に問い合わせ
3. 必要に応じて書類生成タブで契約関連書類を作成
`,
		contract.CustomerName,
		contract.ContractType,
		contract.ProductService,
		contract.ContractAmount,
		contract.ContractDate,
		contract.SalesPersonName,
		contract.ContractType,
		contract.ContractType,
	)
	
	return result
}

func (app *SalesClosingCore) generateMockDocument(docType, customer, product, amount, date string) string {
	currentDate := time.Now().Format("2006年01月02日")
	
	switch docType {
	case "契約書":
		return fmt.Sprintf(`
業務委託契約書

契約日: %s
契約番号: CONTRACT-%s

甲：株式会社ABC（以下「甲」という）
乙：%s（以下「乙」という）

甲と乙は、以下の条項により業務委託契約を締結する。

第1条（業務内容）
甲は乙に対し、以下の業務を委託する。
業務名: %s
業務期間: %s より 1年間
委託料: %s（税込）

第2条（支払い条件）
甲は乙に対し、業務完了後30日以内に委託料を支払うものとする。

第3条（契約の変更）
本契約の変更は、甲乙協議の上、書面により行うものとする。

以上、本契約の成立を証するため、甲乙各1通を作成し、記名押印の上、各自保管するものとする。

%s

甲：株式会社ABC
代表取締役　田中太郎　　印

乙：%s
代表者　　　　　　　　　印
`, date, strings.Replace(strings.Replace(strings.Replace(currentDate, "年", "", -1), "月", "", -1), "日", "", -1), customer, product, date, amount, currentDate, customer)

	case "見積書":
		return fmt.Sprintf(`
見　積　書

見積日: %s
見積番号: EST-%s
有効期限: %s

%s 御中

下記の通りお見積もりいたします。

件名: %s

━━━━━━━━━━━━━━━━━━━━━━━━━━
項目                    数量    単価        金額
━━━━━━━━━━━━━━━━━━━━━━━━━━
%s                      1式     %s         %s
━━━━━━━━━━━━━━━━━━━━━━━━━━
                                小計: %s
                                税額: 未定
                               合計: %s
━━━━━━━━━━━━━━━━━━━━━━━━━━

支払条件: 請求書発行後30日以内
納期: 契約締結後2週間

株式会社ABC
営業部　営業担当者
TEL: 03-1234-5678
Email: sales@abc-corp.com
`, currentDate, strings.Replace(strings.Replace(strings.Replace(currentDate, "年", "", -1), "月", "", -1), "日", "", -1), 
		time.Now().AddDate(0, 0, 30).Format("2006年01月02日"), customer, product, product, amount, amount, amount, amount)

	case "発注書":
		return fmt.Sprintf(`
発　注　書

発注日: %s
発注番号: PO-%s

%s 御中

下記の通り発注いたします。

━━━━━━━━━━━━━━━━━━━━━━━━━━
品名・仕様                数量    単価      金額
━━━━━━━━━━━━━━━━━━━━━━━━━━
%s                        1式     %s       %s
━━━━━━━━━━━━━━━━━━━━━━━━━━
                                  合計: %s

納期: %s
納入場所: 弊社指定場所
支払条件: 納品完了後30日以内

株式会社ABC
調達部
`, currentDate, strings.Replace(strings.Replace(strings.Replace(currentDate, "年", "", -1), "月", "", -1), "日", "", -1), customer, product, amount, amount, amount, date)

	default:
		return fmt.Sprintf(`
%s

作成日: %s

件名: %sに関する%s

詳細:
顧客名: %s
製品/サービス: %s
金額: %s
日付: %s

本書類は Confluence テンプレートを基に自動生成されました。
内容をご確認の上、必要に応じて修正してください。

作成者: 営業クロージング支援AIエージェント
`, docType, currentDate, product, docType, customer, product, amount, date)
	}
}