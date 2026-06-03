# My RALPH Workflow for Claude Code

By [Adam Tuttle](https://adamtuttle.codes/) on Jan 27, 2026

Source: https://adamtuttle.codes/blog/2026/my-ralph-workflow-for-claude-code/

## What is RALPH?

The super, super short version: it's a technique for running Claude Code, headlessly, in a for loop, which has some nice advantages.

**It's worth noting:** RALPH is probably not worth using on a Claude Pro subscription. It's going to chew through your credits quickly. Even on the $100/month plan, if I'm being particularly productive, I have run out of quota with about an hour remaining in my (5-hour) session window.

## Setup

First, add `~/.bin` to your `$PATH` if you haven't already. Then create three executable scripts in that directory: `plan`, `ralph`, and `ralph-install`. You'll also need a `~/.ralph/` directory to store your template files.

Here's the basic structure:

```
~/.bin/
├── plan               # Starts planning sessions
├── ralph              # Runs iterative development cycles
└── ralph-install      # Sets up .ralph directory in your project

~/.ralph/
├── prd.example.json   # PRD template
├── progress.md        # Example progress document
└── ralph.sh           # Ralph script template
```

You don't technically need the `ralph` script, it just executes the `./.ralph/ralph.sh` script, but who wants to type `./.ralph/ralph.sh` when you can just type `ralph`?

### The Files

**~/.bin/plan**

```bash
#!/bin/bash

if [ $# -eq 0 ]; then
	echo "Usage: plan <feature description>"
	exit 1
fi

claude "Study .ralph/prd.example.json and, if it exists, .ralph/prd.json.
Help me plan changes to my project. First give me an opportunity to brain-dump my thoughts, and then when I ask if you
have any questions, ask relevant/clarifying questions to better understand the problem/solution and then add the result to .ralph/prd.json.
The change we're planning is: $*"
```

**~/.bin/ralph**

```bash
sh ./.ralph/ralph.sh "$@"
```

**~/.bin/ralph-install**

```bash
#!/usr/bin/env bash

echo "Adding ralph-scripts to git ignore"
echo "" >> .gitignore # a blank line
echo ".ralph" >> .gitignore

echo "Setting up .ralph directory"
mkdir -p ./.ralph
cp ~/.ralph/ralph.sh ./.ralph/ralph.sh
cp ~/.ralph/prd.example.json ./.ralph/prd.example.json
cp ~/.ralph/progress.md ./.ralph/progress.md

echo "Install complete"
echo "--------------------------------"
echo "Customize this project's ralph script at .ralph/ralph.sh"
echo "--------------------------------"
echo 'Usage: plan "name of feature"'
echo "Usage: ralph <iterations>"
```

**~/.ralph/prd.example.json**

```json
[
	{
		"category": "ui",
		"description": "some ui feature",
		"steps": ["step 1", "step 2"],
		"passes": false
	}
]
```

**~/.ralph/progress.md**

```markdown
# Progress Log

## (current date in yyyy-mm-dd)

Summary of lessons learned and work done.
```

**~/.ralph/ralph.sh**

```bash
# set error mode to ensure any errors are thrown
set -e

if [ -z "$1" ]; then
	echo "Usage: $0 <iterations>"
	exit 1
fi

for ((i=1; i<=$1; i++)); do
	echo "--------------------------------"
	echo "Iteration $i"
	echo "--------------------------------"
	result=$(claude --permission-mode acceptEdits -p "study ./.ralph/prd.json

1. Find the highest-priority item to work on (ignore anything with passes: true) and work only on that item. This should be the one YOU decide has the highest priority - not necessarily the first item in the list.
2. Check that the types check via 'npm run check' (if available) and that the tests pass via 'npm run test' (if available).
3. Update the PRD with the work that was done.
4. Append your progress to the progress.md file. Use this to leave a note for the next person workin in the codebase.
5. Make a git commit of that feature.

ONLY WORK ON A SINGLE FEATURE.
If, while implementing the feature, you notice the PRD is complete, output <PROMISE>COMPLETE</PROMISE>.

If you need additional permissions to complete a task, first double-check that you don't already have the permission. If you have the necessary permission, then proceed to use it. Else print the permissions you need along with <PROMISE>NEED_PERMISSIONS</PROMISE> and exit. I will add them to the .claude/settings.local.json file and re-run.
")

	echo "$result"
	echo ""

	if [[ "$result" == *"<PROMISE>NEED_PERMISSIONS</PROMISE>"* ]]; then
		exit 1
	fi
	if [[ "$result" == *"<PROMISE>COMPLETE</PROMISE>"* ]]; then
		echo "PRD complete, exiting."
		afplay /System/Library/Sounds/Hero.aiff
		exit 0
	fi
done
```

## The Workflow

### Step 1: Make (or add to) your plan

When you're ready to start a new feature, run:

```bash
plan "Add user authentication with OAuth2 support"
```

This kicks off a planning session where you brain-dump your requirements and Claude asks clarifying questions. The goal is to produce a solid PRD that captures everything Claude needs to know to implement the feature autonomously.

You do not have to get it all written down in one prompt. Your thoughts do not need to be well organized. Claude will organize them for you. The prompt built in to the plan command tells claude to assume you're going to send in another prompt with more details, and another, and another, until you tell it you're done. Take as many as you need. Take breaks to think about it.

I've taken to using dictation software to avoid typing out really long prompts. However, it is really helpful to mention specific paths/files that are relevant, so sometimes you have to adjust the output from your dictation. Still worth it.

### Step 2: Get to work

Once your PRD is ready, run:

```bash
ralph 42
```

The number specifies the maximum iterations Claude should run. Each iteration involves Claude reviewing the current state, making progress on the implementation, and documenting what was done. With the prompt I've provided, it will commit at the end of each iteration.

### Step 3: Interrupts

I am specifically AVOIDING using the `--dangerously-skip-permissions` flag (aka "YOLO mode") in this workflow because it's not safe. Instead, my prompts indicate that if extra permissions are needed to accomplish the task, claude should print what that new permission is, and exit.

Since we're not running in YOLO mode, by default claude has very little permission to do anything other than read and write files in the project directory. You provide it additional permissions by specifying them in the file `<project-root>/.claude/settings.local.json`.

If Claude gets stuck or needs input, it'll tell you what it needs and exit. Then you decide whether or not to allow it. If you do want to allow it, add the requested permission to your settings.local.json file. Then run `claude --continue` to resume the session. Tell claude, "I granted that permission. keep going." It will pick right back up where it left off (it just won't continue the for-loop).

### Progress Tracking

The progress.md file serves as a running log of work completed and lessons learned during each iteration. Despite ralph-install blocking the whole .ralph directory in .gitignore, I am currently allowing RALPH to commit both the PRD and the progress.md file in my feature branch. I like the idea of being able to roll back, and when I'm done I can just delete those files and squash the branch and it's like they were never there.

If you're not sure why/how ralph accomplished something, it's probably in progress.md. You could read it (like a chump) or you could tell claude to read it and then ask claude to explain itself to you.

### MODIFY THE SCRIPT

These things are super fluid. Modify them to fit your project, your tools, and your workflow. That's why I copy the script into the project directory: so it can be unique not just to you, but to your project, too.
