package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Endpoints struct {
	API           string `json:"api"`
	OriginTracker string `json:"origin_tracker"`
	Proxy         string `json:"proxy"`
	Telemetry     string `json:"telemetry"`
}

type TokenResponse struct {
	AnnotationsEnabled    bool      `json:"annotations_enabled"`
	ChatEnabled           bool      `json:"chat_enabled"`
	ChatJetbrainsEnabled  bool      `json:"chat_jetbrains_enabled"`
	CodeQuoteEnabled      bool      `json:"code_quote_enabled"`
	CodeReviewEnabled     bool      `json:"code_review_enabled"`
	Codesearch            bool      `json:"codesearch"`
	CopilotignoreEnabled  bool      `json:"copilotignore_enabled"`
	Endpoints             Endpoints `json:"endpoints"`
	ExpiresAt             int64     `json:"expires_at"`
	Individual            bool      `json:"individual"`
	LimitedUserQuotas     *bool     `json:"limited_user_quotas"`
	LimitedUserResetDate  *string   `json:"limited_user_reset_date"`
	NesEnabled            bool      `json:"nes_enabled"`
	Prompt8k              bool      `json:"prompt_8k"`
	PublicSuggestions     string    `json:"public_suggestions"`
	RefreshIn             int       `json:"refresh_in"`
	Sku                   string    `json:"sku"`
	SnippyLoadTestEnabled bool      `json:"snippy_load_test_enabled"`
	Telemetry             string    `json:"telemetry"`
	Token                 string    `json:"token"`
	TrackingID            string    `json:"tracking_id"`
	VscElectronFetcherV2  bool      `json:"vsc_electron_fetcher_v2"`
	Xcode                 bool      `json:"xcode"`
	XcodeChat             bool      `json:"xcode_chat"`
}

type Extra struct {
	Language          string `json:"language"`
	NextIndent        int    `json:"next_indent"`
	TrimByIndentation bool   `json:"trim_by_indentation"`
	PromptTokens      int    `json:"prompt_tokens"`
	SuffixTokens      int    `json:"suffix_tokens"`
}

type CompletionRequest struct {
	Prompt      string   `json:"prompt"`
	Suffix      string   `json:"suffix"`
	MaxTokens   int      `json:"max_tokens"`
	Temperature float64  `json:"temperature"`
	TopP        float64  `json:"top_p"`
	N           int      `json:"n"`
	Stop        []string `json:"stop"`
	Nwo         *string  `json:"nwo"`
	Stream      bool     `json:"stream"`
	Extra       Extra    `json:"extra"`
}

var apiKey = os.Getenv("MISTRAL_API_KEY")

func getToken(w http.ResponseWriter, r *http.Request, port int) {
	server := fmt.Sprintf("http://localhost:%d", port)
	response := TokenResponse{
		AnnotationsEnabled:   true,
		ChatEnabled:          true,
		ChatJetbrainsEnabled: true,
		CodeQuoteEnabled:     true,
		CodeReviewEnabled:    true,
		Codesearch:           true,
		CopilotignoreEnabled: false,
		Endpoints: Endpoints{
			API:           server,
			OriginTracker: server,
			Proxy:         server,
			Telemetry:     server,
		},
		ExpiresAt:             time.Now().Add(24 * time.Hour).Unix(),
		Individual:            true,
		LimitedUserQuotas:     nil,
		LimitedUserResetDate:  nil,
		NesEnabled:            true,
		Prompt8k:              true,
		PublicSuggestions:     "disabled",
		RefreshIn:             15000,
		Sku:                   "sku",
		SnippyLoadTestEnabled: false,
		Telemetry:             "disabled",
		Token:                 "my-token",
		TrackingID:            "my-tracking-id",
		VscElectronFetcherV2:  false,
		Xcode:                 true,
		XcodeChat:             false,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	log.Printf("Handled request: %s %s - Success", r.Method, r.URL.Path)
}

func createCompletion(w http.ResponseWriter, r *http.Request) {
	var request CompletionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Handled request: %s %s - Error: %v", r.Method, r.URL.Path, err)
		return
	}

	reqBody, err := json.Marshal(map[string]interface{}{
		"model":       "codestral-latest",
		"prompt":      request.Prompt,
		"suffix":      request.Suffix,
		"temperature": request.Temperature,
		"top_p":       request.TopP,
		"max_tokens":  request.MaxTokens,
		"stop":        request.Stop,
		"stream":      true,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Handled request: %s %s - Error: %v", r.Method, r.URL.Path, err)
		return
	}

	req, err := http.NewRequest("POST", "https://api.mistral.ai/v1/fim/completions", io.NopCloser(bytes.NewReader(reqBody)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Handled request: %s %s - Error: %v", r.Method, r.URL.Path, err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Handled request: %s %s - Error: %v", r.Method, r.URL.Path, err)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		log.Printf("Handled request: %s %s - Error: Streaming unsupported", r.Method, r.URL.Path)
		return
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		fmt.Fprintf(w, "%s\n\n", scanner.Text())
		flusher.Flush()
	}
	if err := scanner.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Handled request: %s %s - Error: %v", r.Method, r.URL.Path, err)
		return
	}
	log.Printf("Handled request: %s %s - Success", r.Method, r.URL.Path)
}

func patchCopilotLSP(port int) {
	mainJSPath := filepath.Join(
		os.Getenv("HOME"),
		".local",
		"share",
		"nvim",
		"lazy",
		"copilot.lua",
		"copilot",
		"js",
		"main.js",
	)

	server := fmt.Sprintf("http://localhost:%d", port)
	replacements := []struct {
		old string
		new string
	}{
		{"this.apiUrl=i.href", "i.href=\"" + server + "\";this.apiUrl=i.href"},
		{"\"https://copilot-telemetry.githubusercontent.com\"", "\"" + server + "\""},
		{"\"https://origin-tracker.githubusercontent.com\"", "\"" + server + "\""},
		{"\"https://api.githubcopilot.com\"", "\"http://localhost:" + server + "\""},
		{"\"https://copilot-proxy.githubusercontent.com\"", "\"" + server + "\""},
	}

	backupPath := mainJSPath + ".bak"
	_, err := os.Stat(backupPath)
	backupExists := err == nil
	if backupExists {
		backupContent, err := ioutil.ReadFile(backupPath)
		if err != nil {
			log.Fatalf("Error reading backup: %v", err)
		}
		if err := ioutil.WriteFile(mainJSPath, backupContent, 0644); err != nil {
			log.Fatalf("Error restoring backup: %v", err)
		}
	}

	content, err := ioutil.ReadFile(mainJSPath)
	if err != nil {
		log.Fatalf("Error reading LSP server: %v", err)
	}
	if !backupExists {
		if err := ioutil.WriteFile(backupPath, content, 0644); err != nil {
			log.Fatalf("Error creating backup: %v", err)
		}
	}

	for _, replacement := range replacements {
		content = []byte(strings.ReplaceAll(string(content), replacement.old, replacement.new))
	}

	if err := ioutil.WriteFile(mainJSPath, content, 0644); err != nil {
		log.Fatalf("Error writing file: %v", err)
	}

	log.Println("Patched Copilot LSP.")
}

func main() {
	port := 12000 + rand.Intn(4000) // Choose a random port between 12000 and 16000

	patchCopilotLSP(port)
	http.HandleFunc("/copilot_internal/v2/token", func(w http.ResponseWriter, r *http.Request) {
		getToken(w, r, port)
	})
	http.HandleFunc("/v1/engines/{engine_id}/completions", createCompletion)

	log.Printf("Server is running on :%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
