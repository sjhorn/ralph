# ralph

AI-powered feature planner and builder. Ralph uses Claude to turn a feature description into a structured PRD through interactive Q&A, then implements each requirement one at a time — with automatic retries and git commits.

Inspired by [Adam Tuttle's](https://adamtuttle.codes) workflow.

## Prerequisites

- Go 1.22+
- [Claude CLI](https://docs.anthropic.com/en/docs/claude-cli) installed and authenticated

## Installation

```bash
go install github.com/shorn/ralph/cmd/ralph@latest
```

Or build from source:

```bash
git clone https://github.com/shorn/ralph.git
cd ralph
make build
```

## Usage

### 1. Initialize a project

```bash
ralph init
```

Creates a `.ralph/` directory with:
- `config.json` — model, allowed tools, check commands
- `prd.json` — the requirements (starts empty)
- `progress.md` — build log

### 2. Plan a feature

```bash
ralph plan "add user authentication with OAuth"
```

Claude asks clarifying questions in rounds. Answer them, and when it has enough context it generates a structured PRD saved to `.ralph/prd.json`.

Options:
- `--model <model>` — override the configured Claude model
- `--max-rounds <n>` — max Q&A rounds (default: 10)
- `--output <path>` — PRD output path (default: `.ralph/prd.json`)

### 3. Build requirements

```bash
ralph build      # implement the next incomplete item
ralph build 3    # implement the next 3 items
```

For each item, Ralph:
1. Sends the requirement to Claude with implementation steps
2. Runs any check commands from `config.json`
3. Retries up to 3 times if checks fail
4. Marks the item as complete and commits

Options:
- `--model <model>` — override the configured Claude model

### 4. Check progress

```bash
ralph status
```

Shows which PRD items are complete and which remain.

## `.ralph/` directory

| File | Purpose |
|------|---------|
| `config.json` | Model, allowed tools, check commands |
| `prd.json` | The PRD — array of requirements with pass/fail status |
| `progress.md` | Append-only log of completed items |

## License

[MIT](LICENSE)
