# judge-the-code

> Help humans maintain **Judgment** and **Taste** over code in the age of AI-generated software.

English | [中文](README.zh.md)

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

---

## Why this exists

Before AI, writing code and understanding code were the same act. You wrote it, you understood it.

Now AI writes the code. Writing and understanding have decoupled.

**AI makes code run. But running isn't the same as good.**

AI-generated code can:
- Introduce security vulnerabilities you never noticed
- Break the design philosophy your project was built on
- Plant performance time bombs that explode at 100k users
- Burn tokens invisibly — turning your LLM budget into a black hole
- Take "working" shortcuts that become the next developer's nightmare

Spotting these requires truly understanding a codebase's DNA — its design intent, historical decisions, what it cares about.

That understanding can't come from lint. It can't come from tests. **It can only come from human judgment.**

`judge-the-code` is the tool that helps you keep that judgment sharp.

---

## Two things

```
Taste                               Judgment
────────────────────────────────────────────────────
This design is clever — why?        There's a trap here, watch out
This abstraction level is just right  This pattern looks clean but will break
This is a decision worth learning   There's a hidden security hole here
This API makes misuse hard          This assumption fails under high concurrency
                                    This prompt is burning 10x the tokens it needs

Seeing what's good in code          Seeing what's dangerous in code
```

---

## Architecture

```
Tools find problems (deterministic)  +  Claude explains them (semantic)

Skill layer (Claude)            Tool layer (Go binaries)
────────────────────────────────────────────────────────
code-explore                   ← reads directory structure, scc, syft
design-lens                    ← samples source files
demon-hunter  ←─────────────── bearer / trivy / gitleaks
token-optimize                 ← static analysis of LLM call sites
• interprets scan results           deterministic scanning, CVE databases
• judges with project context        single binaries, one-line setup
• explains why it's dangerous
• gives fix recommendations
```

## Skills

| Component | Type | Purpose | Status |
|-----------|------|---------|--------|
| `code-explore` | Skill + Tools | Build global understanding of a codebase (structure, stack, entry points, dependencies) | ✅ Ready |
| `design-lens` | Skill | Extract design philosophy and key decisions — find what's brilliant, reasonable, or questionable | ✅ Ready |
| `demon-hunter` | Skill + Tools | Find security vulnerabilities, dependency CVEs, leaked secrets, performance traps, design hazards | ✅ Ready |
| `token-optimize` | Skill | Discover token waste in LLM integrations — wallet black holes, attention pollution, unnecessary context | ✅ Ready |
| `skill-review` | Skill | Review quality of Skill/Prompt engineering projects — prompt clarity, agent orchestration, injection risks | 🚧 MVP |

Together they form the full `judge-the-code` workflow:

```
code-explore  →  design-lens  →  demon-hunter  →  token-optimize
"What does this   "What's good or     "Where are        "Where is money
 project look      bad about the       the demons?"      being burned?"
 like?"            design?"
  Structure layer    Taste layer       Judgment layer     Economy layer
```

### Skill/Prompt path: `skill-review`

As the AI Agent ecosystem grows, more and more projects are not traditional source code — they are **Skill projects**: natural language prompts, agent definitions, execution flow orchestrations. Current tools (linters, SAST scanners, dependency auditors) are useless on these.

`skill-review` brings judge-the-code's philosophy to this new frontier:

- **Prompt clarity** — Are instructions ambiguous? Would a weaker model misinterpret them?
- **Execution flow design** — Are phases well-structured? Any dead ends or information gaps?
- **Agent orchestration** — Are parallel agents truly independent? Hidden serial dependencies?
- **Fault tolerance** — What happens when an agent fails? Is there a fallback?
- **Security boundaries** — Prompt injection risks? Overly broad file system access?
- **Model portability** — Does it over-rely on one model's quirks (e.g., Claude's XML tags)?

---

## Usage

```bash
/code-explore .       # Step 1: understand the codebase structure
/design-lens .        # Step 2: extract design philosophy
/demon-hunter .       # Step 3: hunt for demons
/token-optimize .     # Step 4: find token waste
/skill-review .       # Skill/Prompt path: review prompt-engineering projects

view .                # Open dashboard in browser

# One-command deterministic skill-review (report + history + dashboard)
~/.agents/skills/judge-the-code/skill-review/bin/run-skill-review .

# Unified mixed-repo entrypoint (auto-detect code/skill/hybrid)
~/.agents/skills/judge-the-code/bin/run-judge .
# - hybrid/skill: runs skill-review automatically
# - hybrid/code:  default full baseline (code-explore/design-lens/demon-hunter/token-optimize)
# - all outputs are centralized at TARGET/.judge-the-code/
```

---

## Installation

```bash
# 1. Copy the skill
cp -r skills/judge-the-code ~/.agents/skills/

# 2. One-time setup (builds dashboard binary + downloads scan tools)
~/.agents/skills/judge-the-code/setup
```

> ⚠️ **Upgrading**: re-run `cp` after each update to overwrite the installed version.

## Dashboard

```bash
~/.agents/skills/judge-the-code/bin/view .
```

Generates `.judge-the-code/dashboard.html` and opens it in your browser. Renders Mermaid architecture diagrams, color-coded design decisions, and severity-graded security findings.

### Output files

```
.judge-the-code/
├── code-explore.md     ← code-explore report
├── design-lens.md      ← design-lens report
├── demon-hunter.md     ← demon-hunter report
├── token-optimize.md   ← token-optimize report
├── skill-review.md     ← skill-review report
├── dashboard.html      ← visual dashboard
└── state/              ← internal skill state (ignore this)
```

---

## Use cases

- **Evaluating whether to adopt a library** — not just what it does, but what traps it hides
- **Learning from well-designed projects** — with a critical eye, finding what's genuinely worth stealing
- **Reviewing AI-generated code** — verifying it didn't break the design philosophy or bury a mine
- **Onboarding to an unfamiliar codebase** — building real judgment, not just a surface tour
- **Auditing LLM integration costs** — finding where tokens are wasted, contexts bloated, money burned

---

## Roadmap

| Milestone | Description | Status |
|-----------|-------------|--------|
| code-explore | Codebase structure analysis with Mermaid visualization | ✅ Shipped |
| design-lens | Design philosophy extraction and decision archaeology | ✅ Shipped |
| demon-hunter | Security scanning (bearer + trivy + gitleaks) + semantic analysis | ✅ Shipped |
| token-optimize | LLM token waste detection and optimization recommendations | ✅ Shipped |
| code-explore hybrid | Deterministic tools (scc + syft) for architecture & dependency analysis | 🚧 In Progress |
| skill-review | Quality review for Skill/Prompt engineering projects | 🚧 MVP |

---

## License

[MIT](LICENSE) — free to use, modify, and distribute.
