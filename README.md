# copilot-usage

A CLI tool to analyze GitHub Copilot usage across your organization. See who's using Copilot, which languages and editors get the most completions, and whether adoption is trending up.

Pulls data from the [GitHub Copilot Metrics API](https://docs.github.com/en/rest/copilot/copilot-metrics) and presents it as colored tables, bar charts, or JSON.

## Install

```sh
pip install git+https://github.com/olruss/copilot-usage.git
```

Or install locally in development mode:

```sh
git clone https://github.com/olruss/copilot-usage.git
cd copilot-usage
pip install -e .
```

You can also use [pipx](https://pipx.pypa.io/) for an isolated install:

```sh
pipx install git+https://github.com/olruss/copilot-usage.git
```

## Authentication

The tool automatically picks up your existing GitHub credentials. In order of priority:

1. `--token` flag (passed directly)
2. **`gh` CLI** — runs `gh auth token` under the hood, so if you're logged in with `gh auth login`, it just works
3. `GITHUB_TOKEN` environment variable

Your token needs the `read:org` scope (the default `gh` login includes this).

## Quick Start

```sh
# Summary of the last 28 days
copilot-usage summary -o my-org

# What happened this week
copilot-usage summary -o my-org -p week

# Day-by-day breakdown
copilot-usage daily -o my-org -p week

# Which languages are people using Copilot with
copilot-usage languages -o my-org

# Which editors are most popular
copilot-usage editors -o my-org

# Is usage going up or down
copilot-usage trends -o my-org
```

## Commands

### `summary`

Aggregated overview for a period: active/engaged users, code suggestions, acceptance rate, lines of code, chat usage, and PR summaries.

```sh
copilot-usage summary -o my-org                           # last 28 days (default)
copilot-usage summary -o my-org -p today
copilot-usage summary -o my-org -p yesterday
copilot-usage summary -o my-org -p week                   # last 7 days
copilot-usage summary -o my-org -p month                  # this calendar month
copilot-usage summary -o my-org -p last-month
copilot-usage summary -o my-org -p year                   # this year (up to 100 days)
copilot-usage summary -o my-org --since 2026-01-01
copilot-usage summary -o my-org --since 2026-01-01 --until 2026-02-01
```

### `daily`

A row for each day showing active users, suggestions, acceptances, acceptance rate, lines accepted, and chats.

```sh
copilot-usage daily -o my-org -p week
```

### `languages`

Top programming languages sorted by acceptances, with a table and a bar chart.

```sh
copilot-usage languages -o my-org
copilot-usage langs -o my-org          # alias
```

### `editors`

Editor/IDE breakdown showing engaged users, suggestions, acceptances, and chats per editor, with a bar chart.

```sh
copilot-usage editors -o my-org
```

### `trends`

Week-over-week comparison. Shows this week vs last week with color-coded deltas (green = up, red = down).

```sh
copilot-usage trends -o my-org
```

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--org` | `-o` | GitHub organization name **(required)** |
| `--period` | `-p` | Named period: `today`, `yesterday`, `week`, `month`, `last-month`, `year` |
| `--since` | | Start date in `YYYY-MM-DD` format |
| `--until` | | End date in `YYYY-MM-DD` format |
| `--json` | | Output JSON instead of tables (for scripting) |
| `--token` | | Pass a GitHub token explicitly |

When neither `--period` nor `--since`/`--until` are set, commands default to the last 28 days.

## JSON Output

Every command supports `--json` for machine-readable output:

```sh
# Pipe to jq
copilot-usage summary -o my-org --json | jq .

# Extract acceptance rate
copilot-usage summary -o my-org --json | jq '.acceptance_rate'

# Get top 5 languages as TSV
copilot-usage languages -o my-org --json | jq -r '.[:5][] | [.name, .acceptances] | @tsv'
```

## API Limitations

- GitHub returns a maximum of **100 days** of history
- Days with fewer than **5 licensed users** are excluded from the response
- Data may be delayed by up to 24 hours
