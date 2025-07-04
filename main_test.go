package main

import (
	"strings"
	"testing"
)

func TestGenerateSummary(t *testing.T) {
	// Test data
	results := []SearchResult{
		{
			Title:       "Go言語の基礎",
			URL:         "https://example.com/go-basics",
			Description: "Go言語は、Googleが開発したプログラミング言語です。シンプルで効率的な開発が可能です。",
		},
		{
			Title:       "Go言語チュートリアル",
			URL:         "https://example.com/go-tutorial",
			Description: "初心者向けのGo言語チュートリアル。基本的な文法から応用まで学べます。",
		},
	}
	
	query := "Go言語"
	summary := generateSummary(results, query)
	
	// Check if summary contains expected elements
	if !strings.Contains(summary, query) {
		t.Errorf("Summary should contain query '%s'", query)
	}
	
	if !strings.Contains(summary, "検索結果") {
		t.Error("Summary should contain '検索結果'")
	}
	
	if !strings.Contains(summary, "Go言語の基礎") {
		t.Error("Summary should contain the first result title")
	}
	
	if len(summary) < 50 {
		t.Error("Summary should be reasonably long")
	}
}

func TestGetMockResults(t *testing.T) {
	query := "テストクエリ"
	results := getMockResults(query)
	
	if len(results) != 5 {
		t.Errorf("Expected 5 mock results, got %d", len(results))
	}
	
	for i, result := range results {
		if result.Title == "" {
			t.Errorf("Result %d should have a title", i)
		}
		if result.URL == "" {
			t.Errorf("Result %d should have a URL", i)
		}
		if result.Description == "" {
			t.Errorf("Result %d should have a description", i)
		}
		if !strings.Contains(result.Title, query) && !strings.Contains(result.Description, query) {
			t.Errorf("Result %d should contain the query", i)
		}
	}
}

func TestIsStopWord(t *testing.T) {
	// Test Japanese stop words
	japaneseStopWords := []string{"の", "に", "は", "を", "が"}
	for _, word := range japaneseStopWords {
		if !isStopWord(word) {
			t.Errorf("'%s' should be identified as a stop word", word)
		}
	}
	
	// Test English stop words
	englishStopWords := []string{"a", "an", "the", "and", "or"}
	for _, word := range englishStopWords {
		if !isStopWord(word) {
			t.Errorf("'%s' should be identified as a stop word", word)
		}
	}
	
	// Test non-stop words
	normalWords := []string{"Go", "言語", "プログラミング", "development"}
	for _, word := range normalWords {
		if isStopWord(word) {
			t.Errorf("'%s' should not be identified as a stop word", word)
		}
	}
}