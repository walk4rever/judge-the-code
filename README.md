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

## Monorepo Layout

This repository is a **skill monorepo**:

```text
skills/
  judge-the-code/   # root router skill
  code-explore/     # codebase understanding
  design-lens/      # design philosophy review
  demon-hunter/     # security and risk scanning
  token-optimize/   # LLM/token cost review
  skill-review/     # skill/prompt project review
tools/
  judge-the-code/   # shared CLI wrappers and dashboard binary
  view/             # dashboard source code
```

Each directory under `skills/` is an independent skill. Shared runtime tooling lives under `tools/`, not inside a single skill.

## Installation

Clone the repo, then run setup from the monorepo root:

```bash
git clone <repo-url>
cd judge-the-code
./setup
```

`setup` prepares shared CLI tooling under `tools/judge-the-code/` and installs any per-skill helper binaries in place.

## Usage

Use a single entrypoint by default:

```bash
/judge-the-code .
/judge-the-code /path/to/project
```

The root skill routes automatically:

- code repo: `code-explore -> design-lens -> demon-hunter -> token-optimize`
- skill/prompt repo: `skill-review`
- hybrid repo: auto-detects and chooses the appropriate path

Natural-language requests should map to the same root skill, for example:

- "Review this repo end-to-end"
- "Help me understand this codebase and find risks"
- "Audit this AI-generated project"
- "Check this skill project for prompt and orchestration issues"

### Advanced Usage

Use sub-skills directly only when you want one focused slice:

```bash
/code-explore .       # Structure and onboarding
/design-lens .        # Design philosophy
/demon-hunter .       # Security and hazard scan
/token-optimize .     # LLM/token efficiency review
/skill-review .       # Skill/prompt engineering review
```

### Deterministic CLI Entrypoint

```bash
./tools/judge-the-code/bin/judge-the-code .
```

`judge-the-code` is the non-chat unified entrypoint for mixed repositories:

- hybrid/skill: runs `skill-review`
- hybrid/code: runs the full baseline
- all outputs are centralized at `TARGET/.judge-the-code/`

`run-judge` remains as the underlying implementation script for compatibility.

### Dashboard

```bash
./tools/judge-the-code/bin/view .
```

Generates `.judge-the-code/summary.html` and opens it in your browser.

### Output files

```
.judge-the-code/
├── code-explore.md     ← code-explore report
├── design-lens.md      ← design-lens report
├── demon-hunter.md     ← demon-hunter report
├── token-optimize.md   ← token-optimize report
├── skill-review.md     ← skill-review report
├── summary.html        ← visual summary
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
