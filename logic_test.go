package main

import (
	"strings"
	"testing"
)

// Test core business logic without GUI components
func TestContractInfoLogic(t *testing.T) {
	contract := &ContractInfo{
		CustomerName:    "テスト株式会社",
		ContractType:    "新規契約",
		ProductService:  "CRMシステム",
		ContractAmount:  "1,000,000円",
		ContractDate:    "2024-01-15",
		SalesPersonName: "田中太郎",
		Notes:          "テスト用契約",
	}
	
	// Test that all fields are set correctly
	if contract.CustomerName != "テスト株式会社" {
		t.Errorf("Expected customer name to be 'テスト株式会社', got '%s'", contract.CustomerName)
	}
	
	if contract.ContractType != "新規契約" {
		t.Errorf("Expected contract type to be '新規契約', got '%s'", contract.ContractType)
	}
}

func TestAPIConfigLogic(t *testing.T) {
	config := &APIConfig{
		MattermostURL:      "https://mattermost.example.com",
		MattermostToken:    "test-token",
		JiraURL:           "https://jira.example.com",
		JiraUser:          "test-user",
		JiraPassword:      "test-password",
		ConfluenceURL:     "https://confluence.example.com",
		ConfluenceUser:    "test-user",
		ConfluencePassword: "test-password",
	}
	
	// Test configuration validation
	if !strings.HasPrefix(config.MattermostURL, "https://") {
		t.Error("Mattermost URL should use HTTPS")
	}
	
	if !strings.HasPrefix(config.JiraURL, "https://") {
		t.Error("Jira URL should use HTTPS")
	}
	
	if !strings.HasPrefix(config.ConfluenceURL, "https://") {
		t.Error("Confluence URL should use HTTPS")
	}
}

func TestDocumentTemplateLogic(t *testing.T) {
	template := &DocumentTemplate{
		Type:     "契約書",
		Title:    "業務委託契約書",
		Template: "契約書のテンプレート内容",
		Fields:   []string{"顧客名", "製品名", "金額", "日付"},
		Rules:    []string{"承認必須", "法務チェック必要"},
	}
	
	// Test template structure
	if template.Type == "" {
		t.Error("Document type should not be empty")
	}
	
	if len(template.Fields) == 0 {
		t.Error("Document should have fields")
	}
	
	// Check for required fields
	requiredFields := map[string]bool{
		"顧客名": false,
		"製品名": false,
		"金額":  false,
		"日付":  false,
	}
	
	for _, field := range template.Fields {
		if _, exists := requiredFields[field]; exists {
			requiredFields[field] = true
		}
	}
	
	for field, found := range requiredFields {
		if !found {
			t.Errorf("Required field '%s' not found in template", field)
		}
	}
}

// Test logic without GUI dependencies
func TestContractCheckResultGeneration(t *testing.T) {
	app := &SalesClosingApp{
		currentContract: &ContractInfo{
			SalesPersonName: "田中太郎",
		},
	}
	
	contract := &ContractInfo{
		CustomerName:    "テスト株式会社",
		ContractType:    "新規契約",
		ProductService:  "CRMシステム",
		ContractAmount:  "1,000,000円",
		ContractDate:    "2024-01-15",
		SalesPersonName: "田中太郎",
		Notes:          "初回導入案件",
	}
	
	result := app.generateContractCheckResult(contract)
	
	// Test result content
	if !strings.Contains(result, contract.CustomerName) {
		t.Errorf("Result should contain customer name '%s'", contract.CustomerName)
	}
	
	if !strings.Contains(result, "確認すべき情報リスト") {
		t.Error("Result should contain confirmation checklist")
	}
	
	if !strings.Contains(result, "必須確認項目") {
		t.Error("Result should contain required items")
	}
}

// Test document generation logic
func TestDocumentGeneration(t *testing.T) {
	app := &SalesClosingApp{
		currentContract: &ContractInfo{
			SalesPersonName: "田中太郎",
		},
	}
	
	testCases := []struct {
		docType  string
		customer string
		product  string
		amount   string
		date     string
	}{
		{"契約書", "テスト株式会社", "CRMシステム", "1,000,000円", "2024-01-15"},
		{"見積書", "サンプル企業", "ERPソフトウェア", "2,000,000円", "2024-02-01"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.docType, func(t *testing.T) {
			document := app.generateMockDocument(tc.docType, tc.customer, tc.product, tc.amount, tc.date)
			
			// Test basic content
			if !strings.Contains(document, tc.customer) {
				t.Errorf("Document should contain customer name '%s'", tc.customer)
			}
			
			if !strings.Contains(document, tc.product) {
				t.Errorf("Document should contain product '%s'", tc.product)
			}
			
			if len(document) < 100 {
				t.Errorf("Document should be substantial, got %d characters", len(document))
			}
		})
	}
}