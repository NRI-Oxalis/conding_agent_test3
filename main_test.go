package main

import (
	"strings"
	"testing"
)

func TestGenerateContractCheckResult(t *testing.T) {
	app := &SalesClosingApp{}
	
	// Test data
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
	
	// Check if result contains expected elements
	if !strings.Contains(result, contract.CustomerName) {
		t.Errorf("Result should contain customer name '%s'", contract.CustomerName)
	}
	
	if !strings.Contains(result, contract.ContractType) {
		t.Errorf("Result should contain contract type '%s'", contract.ContractType)
	}
	
	if !strings.Contains(result, "確認すべき情報リスト") {
		t.Error("Result should contain confirmation checklist")
	}
	
	if !strings.Contains(result, "必須確認項目") {
		t.Error("Result should contain required confirmation items")
	}
	
	if !strings.Contains(result, "Confluenceリンク") {
		t.Error("Result should contain Confluence links")
	}
	
	if !strings.Contains(result, "Mattermostからのヒント") {
		t.Error("Result should contain Mattermost hints")
	}
	
	if len(result) < 500 {
		t.Error("Result should be reasonably comprehensive")
	}
}

func TestGenerateMockDocument(t *testing.T) {
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
		{"発注書", "例示会社", "Webシステム", "500,000円", "2024-03-01"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.docType, func(t *testing.T) {
			document := app.generateMockDocument(tc.docType, tc.customer, tc.product, tc.amount, tc.date)
			
			// Check if document contains basic required elements
			if !strings.Contains(document, tc.customer) {
				t.Errorf("Document should contain customer name '%s'", tc.customer)
			}
			
			if !strings.Contains(document, tc.product) {
				t.Errorf("Document should contain product '%s'", tc.product)
			}
			
			if !strings.Contains(document, tc.amount) {
				t.Errorf("Document should contain amount '%s'", tc.amount)
			}
			
			// Check document type specific content
			switch tc.docType {
			case "契約書":
				if !strings.Contains(document, "業務委託契約書") {
					t.Error("Contract document should contain contract title")
				}
				if !strings.Contains(document, "甲") || !strings.Contains(document, "乙") {
					t.Error("Contract document should contain parties")
				}
			case "見積書":
				if !strings.Contains(document, "見　積　書") {
					t.Error("Quote document should contain quote title")
				}
				if !strings.Contains(document, "小計") {
					t.Error("Quote document should contain subtotal")
				}
			case "発注書":
				if !strings.Contains(document, "発　注　書") {
					t.Error("Purchase order should contain PO title")
				}
				if !strings.Contains(document, "納期") {
					t.Error("Purchase order should contain delivery date")
				}
			}
			
			if len(document) < 100 {
				t.Errorf("Document should be reasonably long, got %d characters", len(document))
			}
		})
	}
}

func TestContractInfoValidation(t *testing.T) {
	contract := &ContractInfo{
		CustomerName:    "テスト株式会社",
		ContractType:    "新規契約",
		ProductService:  "CRMシステム",
		ContractAmount:  "1,000,000円",
		ContractDate:    "2024-01-15",
		SalesPersonName: "田中太郎",
		Notes:          "テスト用契約",
	}
	
	// Test that all fields are set
	if contract.CustomerName == "" {
		t.Error("Customer name should not be empty")
	}
	
	if contract.ContractType == "" {
		t.Error("Contract type should not be empty")
	}
	
	if contract.ProductService == "" {
		t.Error("Product/service should not be empty")
	}
	
	if contract.ContractAmount == "" {
		t.Error("Contract amount should not be empty")
	}
	
	if contract.ContractDate == "" {
		t.Error("Contract date should not be empty")
	}
	
	if contract.SalesPersonName == "" {
		t.Error("Sales person name should not be empty")
	}
}

func TestAPIConfigStructure(t *testing.T) {
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
	
	// Test that all API configuration fields are accessible
	if config.MattermostURL == "" {
		t.Error("Mattermost URL should be accessible")
	}
	
	if config.JiraURL == "" {
		t.Error("Jira URL should be accessible")
	}
	
	if config.ConfluenceURL == "" {
		t.Error("Confluence URL should be accessible")
	}
}

func TestDocumentTemplateStructure(t *testing.T) {
	template := &DocumentTemplate{
		Type:     "契約書",
		Title:    "業務委託契約書",
		Template: "契約書のテンプレート内容",
		Fields:   []string{"顧客名", "製品名", "金額", "日付"},
		Rules:    []string{"承認必須", "法務チェック必要"},
	}
	
	// Test that document template structure is valid
	if template.Type == "" {
		t.Error("Document type should not be empty")
	}
	
	if len(template.Fields) == 0 {
		t.Error("Document should have fields")
	}
	
	if len(template.Rules) == 0 {
		t.Error("Document should have rules")
	}
	
	// Test that expected fields are present
	expectedFields := []string{"顧客名", "製品名", "金額", "日付"}
	for _, expected := range expectedFields {
		found := false
		for _, field := range template.Fields {
			if field == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Document template should contain field '%s'", expected)
		}
	}
}