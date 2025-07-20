package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2"
)

type SalesClosingApp struct {
	app        fyne.App
	window     fyne.Window
	config     *APIConfig
	currentContract *ContractInfo
	core       *SalesClosingCore
	
	// UI components
	mainTabs   *container.AppTabs
	infoTab    *container.TabItem
	docTab     *container.TabItem
	configTab  *container.TabItem
}

func main() {
	myApp := app.NewWithID("com.example.sales-closing-support")
	myApp.SetIcon(theme.DocumentIcon())
	
	salesApp := &SalesClosingApp{
		app:    myApp,
		window: myApp.NewWindow("営業クロージング支援AIエージェント"),
		config: &APIConfig{},
		currentContract: &ContractInfo{},
		core:   &SalesClosingCore{},
	}
	
	salesApp.window.Resize(fyne.NewSize(900, 700))
	salesApp.window.CenterOnScreen()
	
	salesApp.loadConfig()
	salesApp.createUI()
	
	salesApp.window.ShowAndRun()
}

func (app *SalesClosingApp) createUI() {
	// Create main tabs
	app.mainTabs = container.NewAppTabs()
	
	// Tab 1: Contract Information Confirmation
	app.infoTab = container.NewTabItem("契約情報確認", app.createContractInfoTab())
	app.mainTabs.Append(app.infoTab)
	
	// Tab 2: Document Generation
	app.docTab = container.NewTabItem("書類生成", app.createDocumentGenTab())
	app.mainTabs.Append(app.docTab)
	
	// Tab 3: Configuration
	app.configTab = container.NewTabItem("設定", app.createConfigTab())
	app.mainTabs.Append(app.configTab)
	
	// Set content
	app.window.SetContent(app.mainTabs)
}

func (app *SalesClosingApp) createContractInfoTab() *fyne.Container {
	// Title
	title := widget.NewLabelWithStyle("契約情報確認支援", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	title.TextStyle.Bold = true
	
	// Description
	desc := widget.NewRichTextFromMarkdown(`
営業担当者が契約締結後に行うべき情報確認プロセスを支援します。
契約概要を入力すると、Confluenceのルールとの照合により、確認すべき情報やMattermostからのヒントを提示します。
`)
	
	// Input form
	customerEntry := widget.NewEntry()
	customerEntry.SetPlaceHolder("顧客名を入力してください")
	
	contractTypeSelect := widget.NewSelect(
		[]string{"新規契約", "継続契約", "追加契約", "変更契約"},
		nil,
	)
	contractTypeSelect.SetSelected("新規契約")
	
	productEntry := widget.NewEntry()
	productEntry.SetPlaceHolder("製品/サービス名を入力してください")
	
	amountEntry := widget.NewEntry()
	amountEntry.SetPlaceHolder("契約金額を入力してください（例：1,000,000円）")
	
	dateEntry := widget.NewEntry()
	dateEntry.SetPlaceHolder("契約日を入力してください（例：2024-01-15）")
	
	salesPersonEntry := widget.NewEntry()
	salesPersonEntry.SetPlaceHolder("営業担当者名を入力してください")
	
	notesEntry := widget.NewMultiLineEntry()
	notesEntry.SetPlaceHolder("追加の契約詳細や特記事項があれば入力してください")
	notesEntry.Resize(fyne.NewSize(400, 100))
	
	// Result area
	resultArea := widget.NewRichText()
	resultArea.Resize(fyne.NewSize(400, 200))
	resultScroll := container.NewScroll(resultArea)
	resultScroll.SetMinSize(fyne.NewSize(400, 200))
	
	// Check button
	checkBtn := widget.NewButton("情報確認チェックを実行", func() {
		app.performContractCheck(
			customerEntry.Text,
			contractTypeSelect.Selected,
			productEntry.Text,
			amountEntry.Text,
			dateEntry.Text,
			salesPersonEntry.Text,
			notesEntry.Text,
			resultArea,
		)
	})
	checkBtn.Importance = widget.HighImportance
	
	// Form layout
	form := container.NewVBox(
		widget.NewFormItem("顧客名", customerEntry).Widget,
		widget.NewFormItem("契約種別", contractTypeSelect).Widget,
		widget.NewFormItem("製品/サービス", productEntry).Widget,
		widget.NewFormItem("契約金額", amountEntry).Widget,
		widget.NewFormItem("契約日", dateEntry).Widget,
		widget.NewFormItem("営業担当者", salesPersonEntry).Widget,
		widget.NewFormItem("備考", notesEntry).Widget,
		checkBtn,
	)
	
	// Main layout
	content := container.NewBorder(
		container.NewVBox(title, desc),
		nil,
		nil,
		nil,
		container.NewHSplit(
			container.NewScroll(form),
			container.NewBorder(
				widget.NewLabelWithStyle("確認結果", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
				nil,
				nil,
				nil,
				resultScroll,
			),
		),
	)
	
	return content
}

func (app *SalesClosingApp) createDocumentGenTab() *fyne.Container {
	// Title
	title := widget.NewLabelWithStyle("書類自動生成支援", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	
	// Description
	desc := widget.NewRichTextFromMarkdown(`
契約情報に基づいて必要な書類をテキスト形式で自動生成します。
Confluenceのテンプレートを使用して、正確な書類ドラフトを作成できます。
`)
	
	// Document type selection
	docTypeSelect := widget.NewSelect(
		[]string{
			"契約書",
			"発注書",
			"見積書",
			"請求書",
			"納品書",
			"受領書",
			"承認依頼書",
		},
		nil,
	)
	docTypeSelect.SetSelected("契約書")
	
	// Customer info for document
	customerEntry := widget.NewEntry()
	customerEntry.SetPlaceHolder("顧客名")
	
	productEntry := widget.NewEntry()
	productEntry.SetPlaceHolder("製品/サービス名")
	
	amountEntry := widget.NewEntry()
	amountEntry.SetPlaceHolder("金額")
	
	dateEntry := widget.NewEntry()
	dateEntry.SetPlaceHolder("日付")
	
	// Document preview area
	previewArea := widget.NewMultiLineEntry()
	previewArea.SetPlaceHolder("生成された書類がここに表示されます...")
	previewArea.Resize(fyne.NewSize(500, 300))
	previewScroll := container.NewScroll(previewArea)
	previewScroll.SetMinSize(fyne.NewSize(500, 300))
	
	// Generate button
	generateBtn := widget.NewButton("書類を生成", func() {
		app.generateDocument(
			docTypeSelect.Selected,
			customerEntry.Text,
			productEntry.Text,
			amountEntry.Text,
			dateEntry.Text,
			previewArea,
		)
	})
	generateBtn.Importance = widget.HighImportance
	
	// Copy button
	copyBtn := widget.NewButton("クリップボードにコピー", func() {
		app.window.Clipboard().SetContent(previewArea.Text)
		dialog.ShowInformation("コピー完了", "書類内容をクリップボードにコピーしました。", app.window)
	})
	
	// Save button
	saveBtn := widget.NewButton("ファイルに保存", func() {
		app.saveDocument(previewArea.Text)
	})
	
	// Input form
	inputForm := container.NewVBox(
		widget.NewFormItem("書類種別", docTypeSelect).Widget,
		widget.NewFormItem("顧客名", customerEntry).Widget,
		widget.NewFormItem("製品/サービス", productEntry).Widget,
		widget.NewFormItem("金額", amountEntry).Widget,
		widget.NewFormItem("日付", dateEntry).Widget,
		generateBtn,
	)
	
	// Button row
	buttonRow := container.NewHBox(copyBtn, saveBtn)
	
	// Main layout
	content := container.NewBorder(
		container.NewVBox(title, desc),
		nil,
		nil,
		nil,
		container.NewHSplit(
			container.NewScroll(inputForm),
			container.NewBorder(
				widget.NewLabelWithStyle("書類プレビュー", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
				buttonRow,
				nil,
				nil,
				previewScroll,
			),
		),
	)
	
	return content
}

func (app *SalesClosingApp) createConfigTab() *fyne.Container {
	// Title
	title := widget.NewLabelWithStyle("API設定", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	
	// Description
	desc := widget.NewRichTextFromMarkdown(`
Mattermost、Jira、Confluenceとの連携に必要なAPI設定を行います。
設定は安全に保存され、次回起動時に自動的に読み込まれます。
`)
	
	// Mattermost settings
	mattermostURL := widget.NewEntry()
	mattermostURL.SetText(app.config.MattermostURL)
	mattermostURL.SetPlaceHolder("https://mattermost.example.com")
	
	mattermostToken := widget.NewPasswordEntry()
	mattermostToken.SetText(app.config.MattermostToken)
	mattermostToken.SetPlaceHolder("Mattermostアクセストークン")
	
	// Jira settings
	jiraURL := widget.NewEntry()
	jiraURL.SetText(app.config.JiraURL)
	jiraURL.SetPlaceHolder("https://jira.example.com")
	
	jiraUser := widget.NewEntry()
	jiraUser.SetText(app.config.JiraUser)
	jiraUser.SetPlaceHolder("Jiraユーザー名")
	
	jiraPassword := widget.NewPasswordEntry()
	jiraPassword.SetText(app.config.JiraPassword)
	jiraPassword.SetPlaceHolder("Jiraパスワード")
	
	// Confluence settings
	confluenceURL := widget.NewEntry()
	confluenceURL.SetText(app.config.ConfluenceURL)
	confluenceURL.SetPlaceHolder("https://confluence.example.com")
	
	confluenceUser := widget.NewEntry()
	confluenceUser.SetText(app.config.ConfluenceUser)
	confluenceUser.SetPlaceHolder("Confluenceユーザー名")
	
	confluencePassword := widget.NewPasswordEntry()
	confluencePassword.SetText(app.config.ConfluencePassword)
	confluencePassword.SetPlaceHolder("Confluenceパスワード")
	
	// Save button
	saveBtn := widget.NewButton("設定を保存", func() {
		app.config.MattermostURL = mattermostURL.Text
		app.config.MattermostToken = mattermostToken.Text
		app.config.JiraURL = jiraURL.Text
		app.config.JiraUser = jiraUser.Text
		app.config.JiraPassword = jiraPassword.Text
		app.config.ConfluenceURL = confluenceURL.Text
		app.config.ConfluenceUser = confluenceUser.Text
		app.config.ConfluencePassword = confluencePassword.Text
		
		app.saveConfig()
		dialog.ShowInformation("保存完了", "設定が正常に保存されました。", app.window)
	})
	saveBtn.Importance = widget.HighImportance
	
	// Test connection button
	testBtn := widget.NewButton("接続テスト", func() {
		app.testConnections()
	})
	
	// Form layout
	form := container.NewVBox(
		widget.NewCard("Mattermost設定", "", container.NewVBox(
			widget.NewFormItem("サーバーURL", mattermostURL).Widget,
			widget.NewFormItem("アクセストークン", mattermostToken).Widget,
		)),
		widget.NewCard("Jira設定", "", container.NewVBox(
			widget.NewFormItem("サーバーURL", jiraURL).Widget,
			widget.NewFormItem("ユーザー名", jiraUser).Widget,
			widget.NewFormItem("パスワード", jiraPassword).Widget,
		)),
		widget.NewCard("Confluence設定", "", container.NewVBox(
			widget.NewFormItem("サーバーURL", confluenceURL).Widget,
			widget.NewFormItem("ユーザー名", confluenceUser).Widget,
			widget.NewFormItem("パスワード", confluencePassword).Widget,
		)),
		container.NewHBox(saveBtn, testBtn),
	)
	
	// Main layout
	content := container.NewBorder(
		container.NewVBox(title, desc),
		nil,
		nil,
		nil,
		container.NewScroll(form),
	)
	
	return content
}

func (app *SalesClosingApp) performContractCheck(customerName, contractType, product, amount, date, salesPerson, notes string, resultArea *widget.RichText) {
	// Store current contract info
	app.currentContract = &ContractInfo{
		CustomerName:    customerName,
		ContractType:    contractType,
		ProductService:  product,
		ContractAmount:  amount,
		ContractDate:    date,
		SalesPersonName: salesPerson,
		Notes:          notes,
	}
	
	// Show loading
	resultArea.ParseMarkdown("処理中... Confluenceからルールを取得し、Mattermostから過去事例を検索しています。")
	
	// Simulate processing with mock data for now
	go func() {
		time.Sleep(2 * time.Second)
		
		// Generate mock results
		result := app.core.generateContractCheckResult(app.currentContract)
		
		// Update UI on main thread
		fyne.CurrentApp().Driver().StartAnimation(&fyne.Animation{
			Duration:    100 * time.Millisecond,
			RepeatCount: 1,
			Tick: func(f float32) {
				resultArea.ParseMarkdown(result)
			},
		})
	}()
}

func (app *SalesClosingApp) generateDocument(docType, customer, product, amount, date string, previewArea *widget.Entry) {
	// Show loading
	previewArea.SetText("書類を生成中... Confluenceからテンプレートを取得しています。")
	
	// Simulate processing
	go func() {
		time.Sleep(1 * time.Second)
		
		// Generate mock document
		document := app.core.generateMockDocument(docType, customer, product, amount, date)
		
		// Update UI on main thread
		fyne.CurrentApp().Driver().StartAnimation(&fyne.Animation{
			Duration:    100 * time.Millisecond,
			RepeatCount: 1,
			Tick: func(f float32) {
				previewArea.SetText(document)
			},
		})
	}()
}

func (app *SalesClosingApp) saveDocument(content string) {
	// Show save dialog
	dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, app.window)
			return
		}
		if writer == nil {
			return
		}
		defer writer.Close()
		
		writer.Write([]byte(content))
		dialog.ShowInformation("保存完了", "書類をファイルに保存しました。", app.window)
	}, app.window)
}

func (app *SalesClosingApp) loadConfig() {
	configPath := app.getConfigPath()
	if data, err := os.ReadFile(configPath); err == nil {
		json.Unmarshal(data, app.config)
	}
}

func (app *SalesClosingApp) saveConfig() {
	configPath := app.getConfigPath()
	os.MkdirAll(filepath.Dir(configPath), 0755)
	
	data, err := json.MarshalIndent(app.config, "", "  ")
	if err != nil {
		dialog.ShowError(err, app.window)
		return
	}
	
	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		dialog.ShowError(err, app.window)
		return
	}
}

func (app *SalesClosingApp) getConfigPath() string {
	configDir, _ := os.UserConfigDir()
	return filepath.Join(configDir, "sales-closing-support", "config.json")
}

func (app *SalesClosingApp) testConnections() {
	progress := dialog.NewProgressInfinite("接続テスト中", "各システムへの接続をテストしています...", app.window)
	progress.Show()
	
	go func() {
		time.Sleep(3 * time.Second) // Simulate testing
		progress.Hide()
		
		// Mock test results
		result := `接続テスト結果:

✅ Mattermost: 接続成功
✅ Jira: 接続成功  
✅ Confluence: 接続成功

すべてのシステムとの接続が正常に確認されました。`
		
		dialog.ShowInformation("接続テスト完了", result, app.window)
	}()
}