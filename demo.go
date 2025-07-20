package main

import (
	"fmt"
)

// Demonstrate the core functionality of the sales closing support system
func main() {
	fmt.Println("=== 営業クロージング支援AIエージェント デモ ===")
	fmt.Println()

	// Initialize the core system
	core := &SalesClosingCore{
		config: &APIConfig{
			MattermostURL:   "https://mattermost.example.com",
			MattermostToken: "demo-token",
			JiraURL:        "https://jira.example.com",
			JiraUser:       "demo-user",
			ConfluenceURL:  "https://confluence.example.com",
			ConfluenceUser: "demo-user",
		},
	}

	// Demo contract information
	contract := &ContractInfo{
		CustomerName:    "株式会社サンプル",
		ContractType:    "新規契約",
		ProductService:  "営業支援システム",
		ContractAmount:  "2,500,000円",
		ContractDate:    "2024-01-15",
		SalesPersonName: "田中太郎",
		Notes:          "IT部門との連携が重要な案件",
	}

	fmt.Println("🔍 契約情報確認支援デモ")
	fmt.Println("=" + string(make([]rune, 40)))
	
	// Generate contract check result
	checkResult := core.generateContractCheckResult(contract)
	fmt.Println(checkResult)
	fmt.Println()

	fmt.Println("📝 書類自動生成支援デモ")
	fmt.Println("=" + string(make([]rune, 40)))

	// Demo different document types
	documentTypes := []string{"契約書", "見積書", "発注書"}
	
	for _, docType := range documentTypes {
		fmt.Printf("\n--- %s ---\n", docType)
		document := core.generateMockDocument(
			docType,
			contract.CustomerName,
			contract.ProductService,
			contract.ContractAmount,
			contract.ContractDate,
		)
		
		// Display first few lines of the document
		lines := splitLines(document, 15)
		for _, line := range lines {
			fmt.Println(line)
		}
		if len(splitLines(document, -1)) > 15 {
			fmt.Println("... (省略)")
		}
		fmt.Println()
	}

	fmt.Println("✅ デモ完了")
	fmt.Println()
	fmt.Println("このデモでは以下の機能を確認できました：")
	fmt.Println("• 契約情報の確認サポート")
	fmt.Println("• Confluenceルールとの照合")
	fmt.Println("• Mattermostからのヒント提示")
	fmt.Println("• 複数種類の書類自動生成")
	fmt.Println("• テンプレートベースの書類作成")
	fmt.Println()
	fmt.Println("実際のGUIアプリケーションでは、これらの機能を")
	fmt.Println("直感的なインターフェースで利用できます。")
}

func splitLines(text string, maxLines int) []string {
	lines := []string{}
	current := ""
	
	for _, char := range text {
		if char == '\n' {
			lines = append(lines, current)
			current = ""
			if maxLines > 0 && len(lines) >= maxLines {
				break
			}
		} else {
			current += string(char)
		}
	}
	
	if current != "" && (maxLines < 0 || len(lines) < maxLines) {
		lines = append(lines, current)
	}
	
	return lines
}