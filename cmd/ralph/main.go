package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ANSI color codes
const (
	reset   = "\033[0m"
	bold    = "\033[1m"
	dim     = "\033[2m"
	italic  = "\033[3m"
	cyan    = "\033[36m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	red     = "\033[31m"
	magenta = "\033[35m"
	white   = "\033[37m"
	blue    = "\033[34m"
)

// version is set by goreleaser via ldflags at build time.
var version = "dev"

func ralphArtLines() []string {
	e := "\033"
	return []string{
		e + "[49m   " + e + "[38;5;94;49m▄▄" + e + "[38;5;94;48;5;94m▄" + e + "[38;5;220;48;5;94m▄" + e + "[38;5;94;48;5;94m▄" + e + "[38;5;220;48;5;94m▄" + e + "[38;5;94;48;5;94m▄▄" + e + "[38;5;94;49m▄▄" + e + "[49m   " + e + "[m",
		e + "[49m " + e + "[38;5;94;49m▄" + e + "[38;5;220;48;5;94m▄▄" + e + "[38;5;94;48;5;220m▄" + e + "[38;5;220;48;5;94m▄" + e + "[48;5;220m " + e + "[38;5;230;48;5;94m▄" + e + "[38;5;230;48;5;220m▄" + e + "[38;5;220;48;5;94m▄" + e + "[38;5;94;48;5;220m▄" + e + "[48;5;220m " + e + "[38;5;230;48;5;58m▄" + e + "[38;5;220;48;5;94m▄" + e + "[38;5;94;49m▄" + e + "[49m " + e + "[m",
		e + "[49m " + e + "[48;5;94m " + e + "[48;5;220m " + e + "[38;5;220;48;5;94m▄" + e + "[48;5;220m  " + e + "[38;5;220;48;5;15m▄" + e + "[38;5;15;48;5;0m▄" + e + "[38;5;220;48;5;15m▄" + e + "[48;5;220m " + e + "[38;5;94;48;5;184m▄" + e + "[38;5;220;48;5;15m▄" + e + "[38;5;15;48;5;232m▄" + e + "[38;5;220;48;5;15m▄" + e + "[38;5;94;48;5;94m▄" + e + "[49m " + e + "[m",
		e + "[49;38;5;94m▀▀" + e + "[48;5;220m     " + e + "[38;5;94;48;5;220m▄▄" + e + "[38;5;94;48;5;0m▄" + e + "[38;5;94;48;5;220m▄" + e + "[38;5;220;48;5;94m▄" + e + "[48;5;220m  " + e + "[38;5;220;48;5;94m▄" + e + "[38;5;94;49m▄" + e + "[m",
		e + "[49m  " + e + "[38;5;37;48;5;94m▄" + e + "[38;5;0;48;5;220m▄" + e + "[48;5;220m " + e + "[38;5;220;48;5;232m▄" + e + "[38;5;220;48;5;0m▄" + e + "[48;5;0m " + e + "[48;5;220m " + e + "[38;5;220;48;5;233m▄" + e + "[38;5;220;48;5;232m▄" + e + "[48;5;94m   " + e + "[49;38;5;94m▀" + e + "[49m " + e + "[m",
		e + "[38;5;0;49m▄" + e + "[38;5;44;48;5;0m▄" + e + "[38;5;38;48;5;232m▄" + e + "[38;5;0;48;5;38m▄" + e + "[38;5;37;48;5;37m▄" + e + "[38;5;233;48;5;0m▄" + e + "[38;5;44;48;5;38m▄" + e + "[38;5;38;48;5;0m▄" + e + "[38;5;233;48;5;220m▄" + e + "[48;5;220m " + e + "[38;5;232;48;5;94m▄" + e + "[38;5;220;48;5;0m▄" + e + "[38;5;0;48;5;44m▄" + e + "[38;5;44;48;5;232m▄" + e + "[49m  " + e + "[m",
		e + "[38;5;0;48;5;0m▄" + e + "[38;5;0;48;5;44m▄" + e + "[38;5;44;48;5;44m▄" + e + "[38;5;44;48;5;0m▄▄" + e + "[38;5;44;48;5;38m▄" + e + "[38;5;44;48;5;44m▄" + e + "[38;5;37;48;5;44m▄" + e + "[38;5;0;48;5;37m▄" + e + "[38;5;0;48;5;232m▄" + e + "[38;5;44;48;5;44m▄" + e + "[38;5;38;48;5;44m▄" + e + "[38;5;44;48;5;38m▄" + e + "[38;5;44;48;5;0m▄" + e + "[48;5;0m " + e + "[49m " + e + "[m",
		e + "[49m " + e + "[49;38;5;0m▀" + e + "[38;5;0;48;5;38m▄▄" + e + "[38;5;232;48;5;38m▄" + e + "[38;5;232;48;5;37m▄" + e + "[38;5;0;48;5;0m▄" + e + "[38;5;37;48;5;0m▄" + e + "[38;5;38;48;5;232m▄" + e + "[38;5;44;48;5;44m▄▄▄▄▄" + e + "[38;5;44;48;5;38m▄" + e + "[48;5;0m " + e + "[m",
	}
}

// visibleLen returns the printed width of a string, ignoring ANSI escape sequences.
func visibleLen(s string) int {
	inEsc := false
	n := 0
	for _, r := range s {
		if r == '\033' {
			inEsc = true
			continue
		}
		if inEsc {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEsc = false
			}
			continue
		}
		n++
	}
	return n
}

// padRight pads s with spaces so its visible width reaches width.
func padRight(s string, width int) string {
	pad := width - visibleLen(s)
	if pad <= 0 {
		return s
	}
	return s + strings.Repeat(" ", pad)
}

func printRalph() {
	for _, line := range ralphArtLines() {
		fmt.Fprintf(os.Stderr, "    %s\n", line)
	}
}

func printWelcome() {
	art := ralphArtLines()
	cfg := loadConfig()
	cwd, _ := os.Getwd()
	dir := filepath.Base(cwd)

	// Right-side content
	right := []string{
		"",
		bold + white + "Tips" + reset,
		dim + "  Run " + reset + cyan + "ralph" + reset + dim + " for guided mode" + reset,
		dim + "  Run " + reset + cyan + "ralph init" + reset + dim + " to scaffold a project" + reset,
		dim + "  Run " + reset + cyan + "ralph plan" + reset + dim + " to create a PRD" + reset,
		dim + "  " + strings.Repeat("─", 40) + reset,
	}

	// PRD status
	prd, err := loadPRD()
	if err != nil {
		right = append(right, bold+white+"Status"+reset)
		right = append(right, dim+"  No project found"+reset)
	} else if len(prd) == 0 {
		right = append(right, bold+white+"Status"+reset)
		right = append(right, dim+"  PRD empty — run "+reset+cyan+"ralph plan"+reset)
	} else {
		done := 0
		for _, item := range prd {
			if item.Passes {
				done++
			}
		}
		right = append(right, bold+white+"PRD Status"+reset)
		right = append(right, fmt.Sprintf("  %s%d/%d items complete%s", dim, done, len(prd), reset))
	}

	leftWidth := 24
	rightWidth := 46

	// Ensure enough rows for both columns
	rows := len(art) + 2 // art + blank top/bottom
	if len(right)+1 > rows {
		rows = len(right) + 1
	}

	// Top border
	titleStr := fmt.Sprintf("─── ralph v%s ", version)
	titleRuneLen := len([]rune(titleStr))
	fmt.Fprintf(os.Stderr, "╭%s%s┬%s╮\n",
		titleStr, strings.Repeat("─", leftWidth-titleRuneLen), strings.Repeat("─", rightWidth))

	for i := 0; i < rows; i++ {
		// Left: art centered vertically with 1-row top padding
		leftContent := ""
		artIdx := i - 1 // 1 row of padding before art
		if artIdx >= 0 && artIdx < len(art) {
			leftContent = " " + art[artIdx]
		}

		// Right
		rightContent := ""
		if i < len(right) {
			rightContent = " " + right[i]
		}

		fmt.Fprintf(os.Stderr, "│%s│%s│\n",
			padRight(leftContent, leftWidth),
			padRight(rightContent, rightWidth))
	}

	// Bottom info line
	infoStr := fmt.Sprintf(" %s · %s", cfg.Model, dir)
	fmt.Fprintf(os.Stderr, "│%s│%s│\n",
		padRight(dim+infoStr+reset, leftWidth),
		strings.Repeat(" ", rightWidth))

	// Bottom border
	fmt.Fprintf(os.Stderr, "╰%s┴%s╯\n",
		strings.Repeat("─", leftWidth), strings.Repeat("─", rightWidth))
}

const planSystemPrompt = `You are a product manager helping define requirements for a software feature.

Your job is to ask clarifying questions to understand the feature fully, then produce a structured PRD.

Rules:
1. Ask 2-4 clarifying questions per round to understand scope, users, edge cases, and acceptance criteria.
2. When you have enough information to write a complete PRD, output exactly <DONE> on its own line (with nothing else on that line) to signal you are ready.
3. Do NOT produce the PRD until explicitly asked after signaling <DONE>.
4. When asked to produce the PRD, output ONLY a valid JSON array with no markdown fencing, no explanation, just raw JSON.

PRD JSON format:
[
  {
    "category": "ui|backend|api|data|infra",
    "description": "what this requirement covers",
    "steps": ["implementation step 1", "step 2", ...],
    "passes": false
  }
]

Set "passes" to false for all items (they haven't been implemented yet).
Be thorough but concise. Group related work into logical categories.`

const tddTestWriterPrompt = `You are writing tests for a specific feature requirement. Study the codebase first to understand the language, framework, testing conventions, and project structure.

Your task: write FAILING tests for this PRD item:
%s

Implementation steps:
%s

Rules:
1. Study existing test files to match the project's testing patterns and conventions.
2. Write tests that verify the described behavior — they MUST reference code/functions/endpoints that DO NOT EXIST YET.
3. Tests should fail because the implementation is missing, NOT because of syntax errors.
4. Only create or modify test files. Do NOT implement any production code.
5. Make a git commit with your test files when done.
6. If you cannot determine the testing approach, output <BLOCKED> followed by what you need.`

const tddImplementPrompt = `You are implementing a feature to make failing tests pass.

The PRD item you are implementing:
%s

Implementation steps:
%s

Rules:
1. Read the failing test files first to understand what is expected.
2. Implement the feature so ALL tests pass.
3. You may fix test compilation errors (e.g. import paths, type mismatches from your test-writing phase) but do NOT weaken assertions or delete test cases.
4. Run the test command to verify tests pass.
5. If check commands are configured in .ralph/config.json, run them to verify your work.
6. Update .ralph/prd.json — set passes: true for the item you completed.
7. Append your progress to .ralph/progress.md.
8. Make a git commit of your work.
9. If you are blocked, output <BLOCKED> followed by what you need.`

const testCommandSuggestPrompt = `Analyze this project to determine the correct command to run its test suite.

Inspect the project files — look at build files, config files, existing test files, and project structure to determine:
1. The programming language(s) used
2. The test framework in use (or the standard one for this language)
3. The exact shell command to run all tests

Output ONLY a single JSON object with no markdown fencing:
{
  "test_command": "the exact command to run tests",
  "reasoning": "brief explanation of what you found"
}

Examples of test commands by language:
- Go: "go test ./..."
- Dart/Flutter: "dart test" or "flutter test"
- Node.js: "npm test" or "npx jest" or "npx vitest run"
- Python: "pytest" or "python -m pytest"
- Rust: "cargo test"
- Java/Gradle: "./gradlew test"
- Java/Maven: "mvn test"

If you cannot determine a test command, output:
{"test_command": "", "reasoning": "why you couldn't determine it"}`

const buildSystemPrompt = `Study .ralph/prd.json and .ralph/progress.md.

1. Find the highest-priority incomplete item to work on (ignore anything with passes: true) and work only on that item. This should be the one YOU decide has the highest priority — not necessarily the first item in the list.
2. Implement the feature. Read existing code first to understand the codebase.
3. If check commands are configured in .ralph/config.json, run them to verify your work.
4. Update .ralph/prd.json — set passes: true for the item you completed, or leave it false if not finished.
5. Append your progress to .ralph/progress.md. Use this to leave a note for the next iteration.
6. Make a git commit of your work.

ONLY WORK ON A SINGLE ITEM PER ITERATION.
If all items have passes: true, output <COMPLETE> on its own line.
If you need additional permissions or are blocked, output <BLOCKED> followed by what you need.`

// PRDItem represents a single requirement in the PRD.
type PRDItem struct {
	Category    string   `json:"category"`
	Description string   `json:"description"`
	Steps       []string `json:"steps"`
	Passes      bool     `json:"passes"`
}

// RalphConfig holds project-level configuration.
type RalphConfig struct {
	CheckCommands []string `json:"check_commands"`
	Model         string   `json:"model"`
	AllowedTools  []string `json:"allowed_tools"`
	ClaudeCommand string   `json:"claude_command"`
	TestCommand   string   `json:"test_command"`
}

type ClaudeResponse struct {
	SessionID string `json:"session_id"`
	Result    string `json:"result"`
}

// spinner runs a terminal spinner animation until stopped.
type spinner struct {
	mu      sync.Mutex
	running bool
	done    chan struct{}
}

func newSpinner(label string) *spinner {
	s := &spinner{running: true, done: make(chan struct{})}
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	go func() {
		i := 0
		for {
			s.mu.Lock()
			if !s.running {
				s.mu.Unlock()
				fmt.Printf("\r\033[K") // clear spinner line
				close(s.done)
				return
			}
			s.mu.Unlock()
			fmt.Printf("\r  %s%s %s%s", magenta, frames[i%len(frames)], label, reset)
			i++
			time.Sleep(80 * time.Millisecond)
		}
	}()
	return s
}

func (s *spinner) stop() {
	s.mu.Lock()
	s.running = false
	s.mu.Unlock()
	<-s.done
}

// renderMarkdown converts common markdown formatting to ANSI terminal styles.
func renderMarkdown(text string) string {
	var out strings.Builder

	boldItalic := regexp.MustCompile(`\*\*\*(.+?)\*\*\*`)
	boldText := regexp.MustCompile(`\*\*(.+?)\*\*`)
	italicText := regexp.MustCompile(`\*(.+?)\*`)
	codeText := regexp.MustCompile("`([^`]+)`")

	lines := strings.Split(text, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "---" || trimmed == "***" || trimmed == "___" {
			out.WriteString(dim + strings.Repeat("─", 50) + reset)
			if i < len(lines)-1 {
				out.WriteString("\n")
			}
			continue
		}

		if strings.HasPrefix(trimmed, "### ") {
			out.WriteString(bold + cyan + trimmed[4:] + reset)
			if i < len(lines)-1 {
				out.WriteString("\n")
			}
			continue
		}
		if strings.HasPrefix(trimmed, "## ") {
			out.WriteString(bold + cyan + trimmed[3:] + reset)
			if i < len(lines)-1 {
				out.WriteString("\n")
			}
			continue
		}
		if strings.HasPrefix(trimmed, "# ") {
			out.WriteString(bold + cyan + trimmed[2:] + reset)
			if i < len(lines)-1 {
				out.WriteString("\n")
			}
			continue
		}

		line = boldItalic.ReplaceAllString(line, bold+italic+"$1"+reset)
		line = boldText.ReplaceAllString(line, bold+white+"$1"+reset)
		line = italicText.ReplaceAllString(line, italic+"$1"+reset)
		line = codeText.ReplaceAllString(line, yellow+"`$1`"+reset)

		if strings.HasPrefix(trimmed, "- ") {
			indent := line[:len(line)-len(trimmed)]
			line = indent + cyan + "•" + reset + " " + trimmed[2:]
		}

		numberedItem := regexp.MustCompile(`^(\s*)(\d+)\.\s+`)
		if numberedItem.MatchString(line) {
			line = numberedItem.ReplaceAllString(line, "${1}"+cyan+"${2}."+reset+" ")
		}

		out.WriteString(line)
		if i < len(lines)-1 {
			out.WriteString("\n")
		}
	}

	return out.String()
}

func runClaude(claudeCmd string, sessionID string, prompt string, model string, allowedTools []string) (ClaudeResponse, error) {
	args := []string{"-p", "--output-format", "json", "--model", model}
	if sessionID != "" {
		args = append(args, "--resume", sessionID)
	}
	for _, tool := range allowedTools {
		args = append(args, "--allowedTools", tool)
	}

	cmd := exec.Command(claudeCmd, args...)
	cmd.Stdin = strings.NewReader(prompt)
	cmd.Stderr = os.Stderr

	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return ClaudeResponse{}, fmt.Errorf("%s exited with code %d — this may indicate a quota limit, authentication issue, or interrupted session", claudeCmd, exitErr.ExitCode())
		}
		if _, lookErr := exec.LookPath(claudeCmd); lookErr != nil {
			return ClaudeResponse{}, fmt.Errorf("%s not found in PATH", claudeCmd)
		}
		return ClaudeResponse{}, fmt.Errorf("%s command failed: %w", claudeCmd, err)
	}

	var resp ClaudeResponse
	if err := json.Unmarshal(output, &resp); err != nil {
		return ClaudeResponse{}, fmt.Errorf("failed to parse response: %w\nraw output: %s", err, string(output))
	}

	return resp, nil
}

func printHeader(description string, model string, maxRounds int, output string) {
	fmt.Println()
	fmt.Printf("  %s%s╭─ ralph plan%s\n", bold, cyan, reset)
	fmt.Printf("  %s│%s %s\n", cyan, reset, description)
	fmt.Printf("  %s│%s\n", cyan, reset)
	fmt.Printf("  %s│%s  %smodel:%s %s  %srounds:%s %d  %soutput:%s %s\n",
		cyan, reset, dim, reset, model, dim, reset, maxRounds, dim, reset, output)
	fmt.Printf("  %s╰─%s\n", cyan, reset)
	fmt.Println()
}

func printQuestions(result string, round int, maxRounds int) {
	fmt.Printf("  %s%sClaude%s %s(round %d/%d)%s\n", bold, cyan, reset, dim, round, maxRounds, reset)
	fmt.Printf("  %s%s%s\n", dim, strings.Repeat("─", 50), reset)

	rendered := renderMarkdown(strings.TrimSpace(result))
	lines := strings.Split(rendered, "\n")
	for _, line := range lines {
		fmt.Printf("  %s\n", line)
	}
	fmt.Println()
}

func printDone() {
	fmt.Println()
	fmt.Printf("  %s%s✓ Claude is satisfied with the requirements%s\n", bold, green, reset)
	fmt.Println()
}

func printPRDSummary(prd []PRDItem, outputPath string) {
	fmt.Println()
	fmt.Printf("  %s%s╭─ PRD Generated%s\n", bold, green, reset)
	fmt.Printf("  %s│%s\n", green, reset)

	for _, item := range prd {
		fmt.Printf("  %s│%s  %s%-8s%s %s %s(%d steps)%s\n",
			green, reset, yellow, item.Category, reset, item.Description, dim, len(item.Steps), reset)
	}

	fmt.Printf("  %s│%s\n", green, reset)
	fmt.Printf("  %s╰─%s %sSaved to %s%s\n", green, reset, dim, outputPath, reset)
	fmt.Println()
}

func readUserInput(scanner *bufio.Scanner) string {
	fmt.Printf("  %s%sYour answers%s %s(blank line to submit)%s\n", bold, white, reset, dim, reset)
	fmt.Printf("  %s%s%s\n", dim, strings.Repeat("─", 50), reset)
	fmt.Printf("  %s❯%s ", green, reset)

	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		lines = append(lines, line)
		fmt.Printf("  %s❯%s ", green, reset)
	}
	fmt.Println()
	return strings.Join(lines, "\n")
}

// --- Init helpers ---

func defaultConfig() RalphConfig {
	return RalphConfig{
		CheckCommands: []string{},
		Model:         "sonnet",
		AllowedTools:  []string{"Read", "Edit", "Write", "Bash"},
		ClaudeCommand: "claude",
	}
}

// resolveClaudeCmd returns the claude command, preferring: flag > config > default.
func resolveClaudeCmd(flagVal string, cfg RalphConfig) string {
	if flagVal != "" {
		return flagVal
	}
	if cfg.ClaudeCommand != "" {
		return cfg.ClaudeCommand
	}
	return "claude"
}

// testCommandSuggestion is the JSON response from Claude when suggesting a test command.
type testCommandSuggestion struct {
	TestCommand string `json:"test_command"`
	Reasoning   string `json:"reasoning"`
}

// suggestTestCommand asks Claude to inspect the project and suggest the right test command.
func suggestTestCommand(claudeCmd string, model string) (string, string, error) {
	resp, err := runClaude(claudeCmd, "", testCommandSuggestPrompt, model, []string{"Read", "Bash"})
	if err != nil {
		return "", "", fmt.Errorf("failed to detect test command: %w", err)
	}

	var suggestion testCommandSuggestion
	if err := json.Unmarshal([]byte(strings.TrimSpace(resp.Result)), &suggestion); err != nil {
		return "", "", fmt.Errorf("could not parse suggestion: %w\nraw: %s", err, resp.Result)
	}

	if suggestion.TestCommand == "" {
		return "", suggestion.Reasoning, fmt.Errorf("could not determine test command: %s", suggestion.Reasoning)
	}

	return suggestion.TestCommand, suggestion.Reasoning, nil
}

func saveConfig(cfg RalphConfig) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(".ralph/config.json", append(data, '\n'), 0644)
}

func loadConfig() RalphConfig {
	cfg := defaultConfig()
	data, err := os.ReadFile(".ralph/config.json")
	if err != nil {
		return cfg
	}
	json.Unmarshal(data, &cfg)
	return cfg
}

func loadPRD() ([]PRDItem, error) {
	data, err := os.ReadFile(".ralph/prd.json")
	if err != nil {
		return nil, fmt.Errorf("cannot read .ralph/prd.json: %w", err)
	}
	var prd []PRDItem
	if err := json.Unmarshal(data, &prd); err != nil {
		return nil, fmt.Errorf("invalid prd.json: %w", err)
	}
	return prd, nil
}

// prdActuallyComplete re-reads the PRD from disk and checks whether all items
// truly have passes: true. Returns (allDone, done, total).
func prdActuallyComplete() (bool, int, int) {
	prd, err := loadPRD()
	if err != nil || len(prd) == 0 {
		return false, 0, 0
	}
	done := 0
	for _, item := range prd {
		if item.Passes {
			done++
		}
	}
	return done == len(prd), done, len(prd)
}

func savePRD(prd []PRDItem) error {
	data, err := json.MarshalIndent(prd, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(".ralph/prd.json", data, 0644)
}

func ensureInit() error {
	// Ensure git repo exists (no-op if already initialized)
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		cmd := exec.Command("git", "init", "-b", "main")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to initialize git repo: %w", err)
		}
	}

	if err := os.MkdirAll(".ralph", 0755); err != nil {
		return fmt.Errorf("failed to create .ralph/: %w", err)
	}

	// prd.json — empty array
	if _, err := os.Stat(".ralph/prd.json"); os.IsNotExist(err) {
		if err := os.WriteFile(".ralph/prd.json", []byte("[]\n"), 0644); err != nil {
			return err
		}
	}

	// progress.md — header
	if _, err := os.Stat(".ralph/progress.md"); os.IsNotExist(err) {
		if err := os.WriteFile(".ralph/progress.md", []byte("# Ralph Progress Log\n\n"), 0644); err != nil {
			return err
		}
	}

	// config.json — defaults
	if _, err := os.Stat(".ralph/config.json"); os.IsNotExist(err) {
		cfg := defaultConfig()
		data, err := json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			return err
		}
		if err := os.WriteFile(".ralph/config.json", append(data, '\n'), 0644); err != nil {
			return err
		}
	}

	return nil
}

// --- Subcommand: init ---

func cmdInit() {
	if err := ensureInit(); err != nil {
		fmt.Fprintf(os.Stderr, "  %s%s✗ %v%s\n", bold, red, err, reset)
		os.Exit(1)
	}
	fmt.Printf("  %s%s✓ .ralph/ initialized%s\n", bold, green, reset)
	fmt.Printf("    %s•%s config.json\n", cyan, reset)
	fmt.Printf("    %s•%s prd.json\n", cyan, reset)
	fmt.Printf("    %s•%s progress.md\n", cyan, reset)
	fmt.Println()
}

// --- Subcommand: status ---

func cmdStatus() {
	prd, err := loadPRD()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  %s%s✗ %v%s\n", bold, red, err, reset)
		os.Exit(1)
	}

	if len(prd) == 0 {
		fmt.Printf("  %s%s⚠ PRD is empty. Run `ralph plan` first.%s\n\n", bold, yellow, reset)
		return
	}

	fmt.Println()
	fmt.Printf("  %s%s╭─ ralph status%s\n", bold, cyan, reset)
	fmt.Printf("  %s│%s\n", cyan, reset)

	done := 0
	for i, item := range prd {
		indicator := fmt.Sprintf("%s✗%s", red, reset)
		if item.Passes {
			indicator = fmt.Sprintf("%s✓%s", green, reset)
			done++
		}
		fmt.Printf("  %s│%s  %s %s%-8s%s %s\n",
			cyan, reset, indicator, yellow, item.Category, reset, item.Description)
		_ = i
	}

	fmt.Printf("  %s│%s\n", cyan, reset)
	fmt.Printf("  %s╰─%s %s%d/%d complete%s\n", cyan, reset, dim, done, len(prd), reset)
	fmt.Println()
}

// --- Subcommand: plan ---

func cmdPlan(args []string) {
	fs := flag.NewFlagSet("plan", flag.ExitOnError)
	model := fs.String("model", "", "Claude model to use")
	claudeCmd := fs.String("claude-cmd", "", "Claude CLI command (default: claude)")
	maxRounds := fs.Int("max-rounds", 10, "Maximum Q&A rounds")
	output := fs.String("output", ".ralph/prd.json", "PRD output path")
	fs.Parse(args)

	if fs.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "  %s%sUsage:%s ralph plan [flags] \"feature description\"\n\n", bold, red, reset)
		os.Exit(1)
	}

	// Ensure .ralph/ exists
	if err := ensureInit(); err != nil {
		fmt.Fprintf(os.Stderr, "  %s%s✗ %v%s\n", bold, red, err, reset)
		os.Exit(1)
	}

	// Resolve model and command: flag > config > default
	cfg := loadConfig()
	effectiveModel := cfg.Model
	if *model != "" {
		effectiveModel = *model
	}
	effectiveCmd := resolveClaudeCmd(*claudeCmd, cfg)

	description := fs.Arg(0)

	// Load existing PRD if present, so Claude can build on prior work
	existingPRD := ""
	if data, err := os.ReadFile(*output); err == nil {
		content := strings.TrimSpace(string(data))
		if content != "" && content != "[]" {
			existingPRD = fmt.Sprintf("\n\nExisting PRD (from previous planning sessions):\n%s\n\nBuild on this existing PRD. You may add new items, refine existing ones, or leave completed items as-is.", content)
		}
	}

	initialPrompt := fmt.Sprintf("%s\n\nFeature to plan: %s%s\n\nPlease ask your clarifying questions.", planSystemPrompt, description, existingPRD)

	printHeader(description, effectiveModel, *maxRounds, *output)

	spin := newSpinner("Thinking...")
	resp, err := runClaude(effectiveCmd, "", initialPrompt, effectiveModel, nil)
	spin.stop()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  %s%s✗ Error: %v%s\n", bold, red, err, reset)
		os.Exit(1)
	}

	sessionID := resp.SessionID
	scanner := bufio.NewScanner(os.Stdin)

	for round := 1; round <= *maxRounds; round++ {
		if strings.Contains(resp.Result, "<DONE>") {
			printDone()
			break
		}

		printQuestions(resp.Result, round, *maxRounds)

		userInput := readUserInput(scanner)
		if userInput == "" {
			fmt.Printf("  %s%s⚠ No input provided, ending session.%s\n\n", bold, yellow, reset)
			break
		}

		spin := newSpinner("Thinking...")
		resp, err = runClaude(effectiveCmd, sessionID, userInput, effectiveModel, nil)
		spin.stop()
		if err != nil {
			fmt.Fprintf(os.Stderr, "  %s%s✗ Error in round %d: %v%s\n", bold, red, round, err, reset)
			os.Exit(1)
		}

		if round == *maxRounds && !strings.Contains(resp.Result, "<DONE>") {
			fmt.Printf("\n  %s%s⚠ Max rounds reached. Generating PRD with current information...%s\n", bold, yellow, reset)
		}
	}

	// Final prompt to generate the PRD
	spin = newSpinner("Generating PRD...")
	finalPrompt := "Now produce the PRD as a JSON array. Output ONLY valid JSON, no markdown fencing, no explanation."
	resp, err = runClaude(effectiveCmd, sessionID, finalPrompt, effectiveModel, nil)
	spin.stop()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  %s%s✗ Error generating PRD: %v%s\n", bold, red, err, reset)
		os.Exit(1)
	}

	var prd []PRDItem
	prdText := strings.TrimSpace(resp.Result)
	if err := json.Unmarshal([]byte(prdText), &prd); err != nil {
		fmt.Fprintf(os.Stderr, "  %s%s⚠ Warning: PRD output is not valid JSON: %v%s\n", bold, yellow, err, reset)
		fmt.Fprintf(os.Stderr, "  %sRaw output:%s\n%s\n", dim, reset, prdText)
		// Write raw text as fallback
		if dir := filepath.Dir(*output); dir != "." {
			os.MkdirAll(dir, 0755)
		}
		os.WriteFile(*output, []byte(prdText), 0644)
		return
	}

	if dir := filepath.Dir(*output); dir != "." {
		os.MkdirAll(dir, 0755)
	}

	formatted, err := json.MarshalIndent(prd, "", "  ")
	if err != nil {
		formatted = []byte(prdText)
	}

	if err := os.WriteFile(*output, formatted, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "  %s%s✗ Error writing PRD: %v%s\n", bold, red, err, reset)
		os.Exit(1)
	}

	printPRDSummary(prd, *output)
}

// --- Subcommand: build ---

func cmdBuild(args []string) {
	fs := flag.NewFlagSet("build", flag.ExitOnError)
	model := fs.String("model", "", "Claude model to use")
	claudeCmd := fs.String("claude-cmd", "", "Claude CLI command (default: claude)")
	tdd := fs.Bool("tdd", false, "Enable TDD gated workflow")
	maxRetries := fs.Int("retries", 3, "Max retries per TDD phase")
	fs.Parse(args)

	iterations := 1
	if fs.NArg() > 0 {
		n, err := strconv.Atoi(fs.Arg(0))
		if err != nil || n < 1 {
			fmt.Fprintf(os.Stderr, "  %s%sUsage:%s ralph build [flags] [N]\n\n", bold, red, reset)
			os.Exit(1)
		}
		iterations = n
	}

	cfg := loadConfig()
	effectiveModel := cfg.Model
	if *model != "" {
		effectiveModel = *model
	}
	effectiveCmd := resolveClaudeCmd(*claudeCmd, cfg)

	if *tdd {
		cmdBuildTDD(iterations, effectiveModel, effectiveCmd, cfg, *maxRetries, false)
		return
	}

	allowedTools := cfg.AllowedTools
	if len(allowedTools) == 0 {
		allowedTools = []string{"Read", "Edit", "Write", "Bash"}
	}

	fmt.Println()
	fmt.Printf("  %s%s╭─ ralph build%s\n", bold, cyan, reset)
	fmt.Printf("  %s│%s  %smodel:%s %s  %siterations:%s %d\n",
		cyan, reset, dim, reset, effectiveModel, dim, reset, iterations)
	fmt.Printf("  %s╰─%s\n", cyan, reset)
	fmt.Println()

	for iter := 1; iter <= iterations; iter++ {
		fmt.Printf("  %s%s▸ Iteration %d/%d%s\n", bold, cyan, iter, iterations, reset)

		spin := newSpinner("Working...")
		resp, err := runClaude(effectiveCmd, "", buildSystemPrompt, effectiveModel, allowedTools)
		spin.stop()
		if err != nil {
			fmt.Fprintf(os.Stderr, "  %s%s✗ %v%s\n", bold, red, err, reset)
			fmt.Fprintf(os.Stderr, "  %sStopping build due to error.%s\n\n", dim, reset)
			break
		}

		if strings.Contains(resp.Result, "<COMPLETE>") {
			allDone, done, total := prdActuallyComplete()
			if allDone {
				fmt.Printf("  %s%s✓ PRD complete!%s\n\n", bold, green, reset)
				break
			}
			fmt.Printf("  %s%s⚠ Claude reported complete, but %d/%d PRD items remain — continuing%s\n\n", bold, yellow, total-done, total, reset)
		}

		if strings.Contains(resp.Result, "<BLOCKED>") {
			fmt.Printf("  %s%s⚠ Blocked: %s%s\n\n", bold, yellow, strings.TrimSpace(resp.Result), reset)
			break
		}

		fmt.Printf("  %s%s✓ Iteration %d complete%s\n\n", bold, green, iter, reset)

		// Show current PRD status
		_, done, total := prdActuallyComplete()
		fmt.Printf("  %s%d/%d PRD items complete%s\n\n", dim, done, total, reset)
	}
}

// --- TDD gated build ---

// runGate runs a shell command and checks the exit code.
// If expectFailure is true, the gate passes when the command exits non-zero.
// Returns (passed, combinedOutput).
func runGate(command string, expectFailure bool) (bool, string) {
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	if expectFailure {
		// Gate passes if the command fails (non-zero exit)
		return err != nil, outputStr
	}
	// Gate passes if the command succeeds (zero exit)
	return err == nil, outputStr
}

func cmdBuildTDD(iterations int, model string, claudeCmd string, cfg RalphConfig, maxRetries int, guided bool) {
	if cfg.TestCommand == "" {
		fmt.Printf("  %s%s⚡ No test_command configured — analyzing project...%s\n\n", bold, yellow, reset)
		spin := newSpinner("Detecting test command...")
		suggested, reasoning, err := suggestTestCommand(claudeCmd, model)
		spin.stop()

		if err != nil || suggested == "" {
			fmt.Fprintf(os.Stderr, "  %s%s✗ Could not detect test command%s\n", bold, red, reset)
			if reasoning != "" {
				fmt.Fprintf(os.Stderr, "  %s%s%s\n", dim, reasoning, reset)
			}
			fmt.Fprintf(os.Stderr, "  %sAdd test_command to .ralph/config.json manually, e.g.: \"test_command\": \"go test ./...\"%s\n\n", dim, reset)
			os.Exit(1)
		}

		fmt.Printf("  %s%s▸ Suggested:%s %s%s%s\n", bold, cyan, reset, yellow, suggested, reset)
		if reasoning != "" {
			fmt.Printf("  %s%s%s\n", dim, reasoning, reset)
		}
		fmt.Println()

		if guided {
			scanner := bufio.NewScanner(os.Stdin)
			choice := readLine(scanner, "Accept (y), edit (e), or reject (n)?")
			switch strings.ToLower(choice) {
			case "y", "yes", "":
				cfg.TestCommand = suggested
			case "e", "edit":
				custom := readLine(scanner, "Enter test command:")
				if custom == "" {
					fmt.Fprintf(os.Stderr, "  %s%s✗ No command entered, aborting.%s\n\n", bold, red, reset)
					os.Exit(1)
				}
				cfg.TestCommand = custom
			default:
				fmt.Fprintf(os.Stderr, "  %s%s✗ Aborted.%s\n\n", bold, red, reset)
				os.Exit(1)
			}
		} else {
			cfg.TestCommand = suggested
		}

		if err := saveConfig(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "  %s%s✗ Failed to save config: %v%s\n\n", bold, red, err, reset)
			os.Exit(1)
		}
		fmt.Printf("  %s%s✓ Saved test_command to .ralph/config.json%s\n\n", bold, green, reset)
	}

	allowedTools := cfg.AllowedTools
	if len(allowedTools) == 0 {
		allowedTools = []string{"Read", "Edit", "Write", "Bash"}
	}

	fmt.Println()
	fmt.Printf("  %s%s╭─ ralph build (TDD)%s\n", bold, cyan, reset)
	fmt.Printf("  %s│%s  %smodel:%s %s  %siterations:%s %d  %sretries:%s %d\n",
		cyan, reset, dim, reset, model, dim, reset, iterations, dim, reset, maxRetries)
	fmt.Printf("  %s│%s  %stest:%s %s\n", cyan, reset, dim, reset, cfg.TestCommand)
	if len(cfg.CheckCommands) > 0 {
		fmt.Printf("  %s│%s  %schecks:%s %s\n", cyan, reset, dim, reset, strings.Join(cfg.CheckCommands, ", "))
	}
	fmt.Printf("  %s╰─%s\n", cyan, reset)
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	for iter := 1; iter <= iterations; iter++ {
		prd, err := loadPRD()
		if err != nil {
			fmt.Fprintf(os.Stderr, "  %s%s✗ %v%s\n", bold, red, err, reset)
			os.Exit(1)
		}

		// Find next incomplete item
		itemIdx := -1
		for i, item := range prd {
			if !item.Passes {
				itemIdx = i
				break
			}
		}
		if itemIdx == -1 {
			fmt.Printf("  %s%s✓ All PRD items complete!%s\n\n", bold, green, reset)
			break
		}

		item := prd[itemIdx]
		stepsStr := strings.Join(item.Steps, "\n- ")
		if stepsStr != "" {
			stepsStr = "- " + stepsStr
		}

		fmt.Printf("  %s%s▸ Item %d/%d: %s%s\n", bold, cyan, itemIdx+1, len(prd), item.Description, reset)
		fmt.Println()

		// Phase 1: Write Tests
		phase1Success := false
		var lastGateOutput string
		for retry := 0; retry <= maxRetries; retry++ {
			prompt := fmt.Sprintf(tddTestWriterPrompt, item.Description, stepsStr)
			if retry > 0 {
				prompt += fmt.Sprintf("\n\nPREVIOUS ATTEMPT FAILED — the tests were expected to fail but they passed.\nGate output:\n%s\n\nPlease write tests that reference unimplemented code so they FAIL.", lastGateOutput)
			}

			retryLabel := ""
			if retry > 0 {
				retryLabel = fmt.Sprintf(" (retry %d/%d)", retry, maxRetries)
			}
			spin := newSpinner(fmt.Sprintf("Writing tests%s...", retryLabel))
			resp, err := runClaude(claudeCmd, "", prompt, model, allowedTools)
			spin.stop()
			if err != nil {
				fmt.Fprintf(os.Stderr, "  %s%s✗ %v%s\n", bold, red, err, reset)
				os.Exit(1)
			}

			if strings.Contains(resp.Result, "<BLOCKED>") {
				fmt.Printf("  %s%s⚠ Blocked: %s%s\n\n", bold, yellow, strings.TrimSpace(resp.Result), reset)
				break
			}

			// Gate: TDD Red — tests must fail
			passed, output := runGate(cfg.TestCommand, true)
			lastGateOutput = output
			if passed {
				fmt.Printf("  %s%s✓ TDD Red — tests fail as expected%s\n", bold, green, reset)
				phase1Success = true
				break
			}

			fmt.Printf("  %s%s✗ TDD Red — tests should fail but passed%s\n", bold, red, reset)
			if retry == maxRetries {
				fmt.Printf("  %s%s⚠ Retries exhausted for test-writing phase%s\n", bold, yellow, reset)
			}
		}

		if !phase1Success {
			if guided {
				choice := handleRetryExhausted(scanner, "test-writing", item.Description)
				if choice == "skip" {
					continue
				} else if choice == "quit" {
					return
				}
				// "retry" — but we already exhausted retries, so skip for now
				continue
			}
			fmt.Fprintf(os.Stderr, "  %s%s✗ TDD Red gate failed after %d retries — halting build%s\n\n", bold, red, maxRetries, reset)
			os.Exit(1)
		}

		// Phase 2: Implement
		phase2Success := false
		lastGateOutput = ""
		for retry := 0; retry <= maxRetries; retry++ {
			prompt := fmt.Sprintf(tddImplementPrompt, item.Description, stepsStr)
			if retry > 0 {
				prompt += fmt.Sprintf("\n\nPREVIOUS ATTEMPT FAILED — gates did not pass.\nGate output:\n%s\n\nFix the implementation so all tests and checks pass.", lastGateOutput)
			}

			retryLabel := ""
			if retry > 0 {
				retryLabel = fmt.Sprintf(" (retry %d/%d)", retry, maxRetries)
			}
			spin := newSpinner(fmt.Sprintf("Implementing%s...", retryLabel))
			resp, err := runClaude(claudeCmd, "", prompt, model, allowedTools)
			spin.stop()
			if err != nil {
				fmt.Fprintf(os.Stderr, "  %s%s✗ %v%s\n", bold, red, err, reset)
				os.Exit(1)
			}

			if strings.Contains(resp.Result, "<BLOCKED>") {
				fmt.Printf("  %s%s⚠ Blocked: %s%s\n\n", bold, yellow, strings.TrimSpace(resp.Result), reset)
				break
			}

			// Gate: TDD Green — tests must pass
			passed, output := runGate(cfg.TestCommand, false)
			if !passed {
				lastGateOutput = output
				fmt.Printf("  %s%s✗ TDD Green — tests still failing%s\n", bold, red, reset)
				if retry == maxRetries {
					fmt.Printf("  %s%s⚠ Retries exhausted for implementation phase%s\n", bold, yellow, reset)
				}
				continue
			}
			fmt.Printf("  %s%s✓ TDD Green — all tests pass%s\n", bold, green, reset)

			// Gate: Check commands
			allChecksPass := true
			for _, check := range cfg.CheckCommands {
				checkPassed, checkOutput := runGate(check, false)
				if !checkPassed {
					lastGateOutput = checkOutput
					fmt.Printf("  %s%s✗ Check failed: %s%s\n", bold, red, check, reset)
					allChecksPass = false
					break
				}
			}
			if !allChecksPass {
				if retry == maxRetries {
					fmt.Printf("  %s%s⚠ Retries exhausted for implementation phase%s\n", bold, yellow, reset)
				}
				continue
			}

			if len(cfg.CheckCommands) > 0 {
				fmt.Printf("  %s%s✓ All checks pass%s\n", bold, green, reset)
			}
			phase2Success = true
			break
		}

		if !phase2Success {
			if guided {
				choice := handleRetryExhausted(scanner, "implementation", item.Description)
				if choice == "skip" {
					continue
				} else if choice == "quit" {
					return
				}
				continue
			}
			fmt.Fprintf(os.Stderr, "  %s%s✗ Implementation gates failed after %d retries — halting build%s\n\n", bold, red, maxRetries, reset)
			os.Exit(1)
		}

		fmt.Printf("  %s%s✓ Item complete (TDD verified)%s\n", bold, green, reset)

		// Show PRD status
		prd, _ = loadPRD()
		done := 0
		for _, p := range prd {
			if p.Passes {
				done++
			}
		}
		fmt.Printf("  %s%d/%d PRD items complete%s\n\n", dim, done, len(prd), reset)
	}
}

// handleRetryExhausted prompts the user in guided mode when retries are exhausted.
// Returns "skip", "retry", or "quit".
func handleRetryExhausted(scanner *bufio.Scanner, phase string, itemDesc string) string {
	fmt.Printf("\n  %s%s⚠ %s phase failed for: %s%s\n", bold, yellow, phase, itemDesc, reset)
	fmt.Printf("  %s(s)kip this item, (r)etry, or (q)uit?%s ", dim, reset)
	if scanner.Scan() {
		switch strings.ToLower(strings.TrimSpace(scanner.Text())) {
		case "r", "retry":
			return "retry"
		case "q", "quit":
			return "quit"
		}
	}
	return "skip"
}

// --- Guided flow (no-arg mode) ---

func readLine(scanner *bufio.Scanner, prompt string) string {
	fmt.Printf("  %s%s%s %s❯%s ", bold, white, prompt, green, reset)
	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text())
	}
	return ""
}

func cmdGuided(claudeCmdFlag string) {
	fmt.Println()
	printWelcome()
	fmt.Println()

	// Step 1: Init
	if err := ensureInit(); err != nil {
		fmt.Fprintf(os.Stderr, "  %s%s✗ %v%s\n", bold, red, err, reset)
		os.Exit(1)
	}

	cfg := loadConfig()
	effectiveModel := cfg.Model
	effectiveCmd := resolveClaudeCmd(claudeCmdFlag, cfg)
	scanner := bufio.NewScanner(os.Stdin)

	// Step 2: Check for existing PRD with incomplete items
	prd, _ := loadPRD()
	incomplete := 0
	for _, item := range prd {
		if !item.Passes {
			incomplete++
		}
	}

	skipToBuild := false
	if incomplete > 0 {
		fmt.Printf("  %s%s⚡ Existing PRD found with %d incomplete item(s)%s\n\n", bold, yellow, incomplete, reset)
		for _, item := range prd {
			if item.Passes {
				fmt.Printf("    %s✓%s %s%-8s%s %s\n", green, reset, yellow, item.Category, reset, item.Description)
			} else {
				fmt.Printf("    %s✗%s %s%-8s%s %s\n", red, reset, yellow, item.Category, reset, item.Description)
			}
		}
		fmt.Println()

		choice := readLine(scanner, "Continue building (b), TDD build (t), add to plan (p), start fresh (f), or quit (q)?")
		switch strings.ToLower(choice) {
		case "t", "tdd":
			cmdBuildTDD(incomplete, effectiveModel, effectiveCmd, cfg, 3, true)
			return
		case "b", "build":
			skipToBuild = true
		case "f", "fresh":
			// Wipe PRD and start over
			if err := os.WriteFile(".ralph/prd.json", []byte("[]\n"), 0644); err != nil {
				fmt.Fprintf(os.Stderr, "  %s%s✗ Failed to reset PRD: %v%s\n", bold, red, err, reset)
				os.Exit(1)
			}
			prd = nil
		case "q", "quit":
			fmt.Printf("  %s%s✓ Goodbye%s\n\n", bold, green, reset)
			return
		case "p", "plan":
			// Fall through to planning
		default:
			// Default to planning
		}
	}

	// Step 3: Plan (skip if user chose to build directly)
	if !skipToBuild {
		fmt.Printf("  %s%sWhat feature would you like to plan?%s\n", bold, white, reset)
		description := readLine(scanner, "")
		if description == "" {
			fmt.Printf("  %s%s⚠ No description provided, exiting.%s\n\n", bold, yellow, reset)
			return
		}

		// Load existing PRD context
		existingPRD := ""
		if data, err := os.ReadFile(".ralph/prd.json"); err == nil {
			content := strings.TrimSpace(string(data))
			if content != "" && content != "[]" {
				existingPRD = fmt.Sprintf("\n\nExisting PRD (from previous planning sessions):\n%s\n\nBuild on this existing PRD. You may add new items, refine existing ones, or leave completed items as-is.", content)
			}
		}

		maxRounds := 10
		initialPrompt := fmt.Sprintf("%s\n\nFeature to plan: %s%s\n\nPlease ask your clarifying questions.", planSystemPrompt, description, existingPRD)

		printHeader(description, effectiveModel, maxRounds, ".ralph/prd.json")

		spin := newSpinner("Thinking...")
		resp, err := runClaude(effectiveCmd, "", initialPrompt, effectiveModel, nil)
		spin.stop()
		if err != nil {
			fmt.Fprintf(os.Stderr, "  %s%s✗ %v%s\n", bold, red, err, reset)
			os.Exit(1)
		}

		sessionID := resp.SessionID

		for round := 1; round <= maxRounds; round++ {
			if strings.Contains(resp.Result, "<DONE>") {
				printDone()
				break
			}

			printQuestions(resp.Result, round, maxRounds)

			userInput := readUserInput(scanner)
			if userInput == "" {
				fmt.Printf("  %s%s⚠ No input provided, ending Q&A.%s\n\n", bold, yellow, reset)
				break
			}

			spin := newSpinner("Thinking...")
			resp, err = runClaude(effectiveCmd, sessionID, userInput, effectiveModel, nil)
			spin.stop()
			if err != nil {
				fmt.Fprintf(os.Stderr, "  %s%s✗ %v%s\n", bold, red, err, reset)
				os.Exit(1)
			}

			if round == maxRounds && !strings.Contains(resp.Result, "<DONE>") {
				fmt.Printf("\n  %s%s⚠ Max rounds reached. Generating PRD with current information...%s\n", bold, yellow, reset)
			}
		}

		// Generate PRD
		spin = newSpinner("Generating PRD...")
		finalPrompt := "Now produce the PRD as a JSON array. Output ONLY valid JSON, no markdown fencing, no explanation."
		resp, err = runClaude(effectiveCmd, sessionID, finalPrompt, effectiveModel, nil)
		spin.stop()
		if err != nil {
			fmt.Fprintf(os.Stderr, "  %s%s✗ %v%s\n", bold, red, err, reset)
			os.Exit(1)
		}

		prdText := strings.TrimSpace(resp.Result)
		if err := json.Unmarshal([]byte(prdText), &prd); err != nil {
			fmt.Fprintf(os.Stderr, "  %s%s⚠ PRD output is not valid JSON: %v%s\n", bold, yellow, err, reset)
			fmt.Fprintf(os.Stderr, "  %sRaw output:%s\n%s\n", dim, reset, prdText)
			os.WriteFile(".ralph/prd.json", []byte(prdText), 0644)
			return
		}

		formatted, _ := json.MarshalIndent(prd, "", "  ")
		if err := os.WriteFile(".ralph/prd.json", formatted, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "  %s%s✗ Error writing PRD: %v%s\n", bold, red, err, reset)
			os.Exit(1)
		}

		printPRDSummary(prd, ".ralph/prd.json")
	}

	// Step 4: Ask to build
	var buildErr error
	prd, buildErr = loadPRD()
	if buildErr != nil {
		fmt.Fprintf(os.Stderr, "  %s%s✗ %v%s\n", bold, red, buildErr, reset)
		os.Exit(1)
	}

	incomplete = 0
	for _, item := range prd {
		if !item.Passes {
			incomplete++
		}
	}

	if incomplete == 0 {
		fmt.Printf("  %s%s✓ All items already complete!%s\n\n", bold, green, reset)
		return
	}

	iterInput := readLine(scanner, fmt.Sprintf("How many iterations? (%d remaining, enter=all, 0=skip, tdd=TDD mode)", incomplete))
	if iterInput == "tdd" {
		cmdBuildTDD(incomplete, effectiveModel, effectiveCmd, cfg, 3, true)
		return
	}
	iterations := incomplete
	if iterInput == "0" {
		fmt.Printf("\n  %s%s✓ Skipping build. Run `ralph build` when ready.%s\n\n", bold, green, reset)
		return
	} else if iterInput != "" {
		n, err := strconv.Atoi(iterInput)
		if err != nil || n < 1 {
			fmt.Fprintf(os.Stderr, "  %s%s⚠ Invalid number, building all %d items%s\n", bold, yellow, incomplete, reset)
		} else {
			iterations = n
		}
	}
	fmt.Println()

	// Step 5: Build
	allowedTools := cfg.AllowedTools
	if len(allowedTools) == 0 {
		allowedTools = []string{"Read", "Edit", "Write", "Bash"}
	}

	fmt.Printf("  %s%s╭─ ralph build%s\n", bold, cyan, reset)
	fmt.Printf("  %s│%s  %smodel:%s %s  %siterations:%s %d\n",
		cyan, reset, dim, reset, effectiveModel, dim, reset, iterations)
	fmt.Printf("  %s╰─%s\n", cyan, reset)
	fmt.Println()

	for iter := 1; iter <= iterations; iter++ {
		fmt.Printf("  %s%s▸ Iteration %d/%d%s\n", bold, cyan, iter, iterations, reset)

		spin := newSpinner("Working...")
		resp, err := runClaude(effectiveCmd, "", buildSystemPrompt, effectiveModel, allowedTools)
		spin.stop()
		if err != nil {
			fmt.Fprintf(os.Stderr, "  %s%s✗ %v%s\n", bold, red, err, reset)
			fmt.Fprintf(os.Stderr, "  %sStopping build due to error.%s\n\n", dim, reset)
			break
		}

		if strings.Contains(resp.Result, "<COMPLETE>") {
			allDone, done, total := prdActuallyComplete()
			if allDone {
				fmt.Printf("  %s%s✓ PRD complete!%s\n\n", bold, green, reset)
				break
			}
			fmt.Printf("  %s%s⚠ Claude reported complete, but %d/%d PRD items remain — continuing%s\n\n", bold, yellow, total-done, total, reset)
		}

		if strings.Contains(resp.Result, "<BLOCKED>") {
			fmt.Printf("  %s%s⚠ Blocked: %s%s\n\n", bold, yellow, strings.TrimSpace(resp.Result), reset)
			break
		}

		fmt.Printf("  %s%s✓ Iteration %d complete%s\n\n", bold, green, iter, reset)

		// Show current PRD status
		_, done, total := prdActuallyComplete()
		fmt.Printf("  %s%d/%d PRD items complete%s\n\n", dim, done, total, reset)
	}

	// Final status
	_, done, total := prdActuallyComplete()
	remaining := total - done
	if remaining > 0 {
		fmt.Printf("  %s%d/%d complete — %d remaining. Run `ralph` or `ralph build %d` to continue%s\n", dim, done, total, remaining, remaining, reset)
	} else if total > 0 {
		fmt.Printf("  %s%s✓ All %d PRD items complete!%s\n", bold, green, total, reset)
	}
	fmt.Println()
}

// --- Main ---

func printUsage() {
	fmt.Fprintf(os.Stderr, "\n")
	printWelcome()
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "  %sUsage:%s\n", dim, reset)
	fmt.Fprintf(os.Stderr, "    ralph                              Guided mode — plan, build, track\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "  %sCommands:%s\n", dim, reset)
	fmt.Fprintf(os.Stderr, "    plan  \"description\"                 Generate a PRD through Q&A\n")
	fmt.Fprintf(os.Stderr, "    build [N]                           Implement N iterations (default: 1)\n")
	fmt.Fprintf(os.Stderr, "    build --tdd [N]                     TDD gated build (test → red → implement → green)\n")
	fmt.Fprintf(os.Stderr, "    build --tdd --retries 5 [N]         TDD with custom retry count (default: 3)\n")
	fmt.Fprintf(os.Stderr, "    init                                Scaffold .ralph/ directory\n")
	fmt.Fprintf(os.Stderr, "    status                              Show PRD progress\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "  %sGlobal flags:%s\n", dim, reset)
	fmt.Fprintf(os.Stderr, "    --claude-cmd <command>              Use a different CLI (default: claude)\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "  %sExamples:%s\n", dim, reset)
	fmt.Fprintf(os.Stderr, "    ralph                               # guided flow\n")
	fmt.Fprintf(os.Stderr, "    ralph plan \"add user authentication\"\n")
	fmt.Fprintf(os.Stderr, "    ralph build 3\n")
	fmt.Fprintf(os.Stderr, "    ralph --claude-cmd claude-local      # use a different CLI\n")
	fmt.Fprintf(os.Stderr, "    ralph status\n")
	fmt.Fprintf(os.Stderr, "\n")
}

func main() {
	// Check for help before flag parsing (flag pkg intercepts -h)
	for _, arg := range os.Args[1:] {
		if arg == "-h" || arg == "--help" || arg == "help" {
			printUsage()
			return
		}
	}

	// Parse global flags — stops at first non-flag arg (the subcommand)
	globalFlags := flag.NewFlagSet("ralph", flag.ContinueOnError)
	claudeCmd := globalFlags.String("claude-cmd", "", "Claude CLI command (default: claude)")
	globalFlags.Usage = func() {} // suppress default usage on parse error
	globalFlags.Parse(os.Args[1:])

	remaining := globalFlags.Args() // everything after flags

	if len(remaining) == 0 {
		cmdGuided(*claudeCmd)
		return
	}

	switch remaining[0] {
	case "init":
		cmdInit()
	case "plan":
		cmdPlan(remaining[1:])
	case "build":
		cmdBuild(remaining[1:])
	case "status":
		cmdStatus()
	default:
		fmt.Fprintf(os.Stderr, "  %s%sUnknown command:%s %s\n\n", bold, red, reset, remaining[0])
		printUsage()
		os.Exit(1)
	}
}
