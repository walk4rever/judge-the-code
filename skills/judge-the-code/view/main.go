package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

//go:embed template.html
var htmlTemplate string

type skillInfo struct {
	ID      string
	Label   string
	Icon    string
	File    string
	Content string
	Empty   bool
}

func main() {
	// Determine project root
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}
	dir := filepath.Join(root, ".judge-the-code")

	skills := []skillInfo{
		{ID: "explore", Label: "code-explore", Icon: "🗺", File: "code-explore.md"},
		{ID: "lens", Label: "design-lens", Icon: "🔍", File: "design-lens.md"},
		{ID: "hunter", Label: "demon-hunter", Icon: "👹", File: "demon-hunter.md"},
		{ID: "optimize", Label: "token-optimize", Icon: "🪙", File: "token-optimize.md"},
		{ID: "review", Label: "skill-review", Icon: "🧭", File: "skill-review.md"},
	}

	for i, s := range skills {
		content, err := os.ReadFile(filepath.Join(dir, s.File))
		if err == nil {
			skills[i].Content = string(content)
		} else {
			skills[i].Empty = true
		}
	}

	// Build JSON payload for each skill
	html := htmlTemplate
	for _, s := range skills {
		placeholder := "<!--CONTENT_" + strings.ToUpper(s.ID) + "-->"
		encoded, _ := json.Marshal(s.Content)
		html = strings.ReplaceAll(html, placeholder, string(encoded))
	}

	// Inject skill-review history JSON (optional)
	historyRaw, err := os.ReadFile(filepath.Join(dir, "state", "skill-review-history.json"))
	historyJSON := "null"
	if err == nil {
		var anyJSON interface{}
		if json.Unmarshal(historyRaw, &anyJSON) == nil {
			if b, e := json.Marshal(anyJSON); e == nil {
				historyJSON = string(b)
			}
		}
	}
	html = strings.ReplaceAll(html, "<!--CONTENT_REVIEW_HISTORY-->", historyJSON)

	// Inject generated time
	html = strings.ReplaceAll(html, "<!--GENERATED_AT-->", time.Now().Format("2006-01-02 15:04"))

	// Write summary HTML
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "error: cannot create %s: %v\n", dir, err)
		os.Exit(1)
	}
	outPath := filepath.Join(dir, "summary.html")
	legacyPath := filepath.Join(dir, "dashboard.html")
	if err := os.WriteFile(outPath, []byte(html), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error: cannot write summary html: %v\n", err)
		os.Exit(1)
	}
	_ = os.Remove(legacyPath)

	absPath, _ := filepath.Abs(outPath)
	fmt.Printf("📊 Summary HTML → %s\n", absPath)
	openBrowser(absPath)
}

func openBrowser(path string) {
	url := "file://" + path
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		fmt.Printf("Open manually: %s\n", url)
		return
	}
	_ = cmd.Start()
}
