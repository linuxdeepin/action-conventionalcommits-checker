## Conventional Commits Checker

A GitHub Action that checks commits inside a Pull Request against the Conventional Commits guideline.

### Usage

#### Default mode (no configuration)

```yaml
- name: Call Conventional Commits Checker
  uses: linuxdeepin/action-conventionalcommits-checker@master
```

This uses the built-in 6 predefined rules. Behavior is identical for all existing callers — no changes required.

#### Full custom mode (`rules`)

Provide a JSON array of rules to completely replace the default rule set:

```yaml
- name: Call Conventional Commits Checker
  uses: linuxdeepin/action-conventionalcommits-checker@master
  with:
    rules: '[{"rule_type":"always_check","check_type":"conventional_commits_base"},{"rule_type":"always_check","check_type":"require_body"}]'
```

You can also place the JSON in a file in your repository and reference it with `rules-file`:

```yaml
- name: Call Conventional Commits Checker
  uses: linuxdeepin/action-conventionalcommits-checker@master
  with:
    rules-file: .github/cc-rules.json
```

#### Hybrid mode (`extra-rules`)

Append custom rules on top of the default (or custom) rule set:

```yaml
- name: Call Conventional Commits Checker
  uses: linuxdeepin/action-conventionalcommits-checker@master
  with:
    extra-rules: '[{"rule_type":"fallback_check","check_type":"allowed_subject_pattern","rule_data":"Allow chore commits","check_data":"chore: .*"}]'
```

`extra-rules` can be combined with `rules`, `rules-file`, or used alone (appends to defaults).

### Inputs

| Input | Description | Default |
|-------|-------------|---------|
| `rules` | JSON string defining all rules (replaces default rules) | `""` |
| `extra-rules` | JSON string defining extra rules to append to the current rule set | `""` |
| `rules-file` | Path to a JSON file in the repo defining all rules (takes priority over `rules`) | `""` |

Priority: `rules-file` > `rules` > default rules. `extra-rules` can be combined with any of the above.

### Rule JSON Format

Each rule is a JSON object with the following fields:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `rule_type` | string | Yes | The type of rule (see table below) |
| `check_type` | string | Yes | The check to perform (see table below) |
| `condition_type` | string | Only for `conditional_check` | The condition matching strategy |
| `rule_data` | string | No | Data for the rule (e.g. email domain, type string, or fallback description) |
| `check_data` | string | No | Data for the check (e.g. trailer keys, regex pattern, or type list) |

#### `rule_type` values

| Value | Description |
|-------|-------------|
| `always_check` | Always required to pass (AND logic) |
| `conditional_check` | Applied only when the condition matches (e.g. author email or commit type) |
| `fallback_check` | If any fallback check passes, the commit is allowed (OR logic) |

#### `check_type` values

| Value | Description |
|-------|-------------|
| `conventional_commits_base` | Require subject line to follow `type: description` or `type(scope): description` |
| `type_in_range` | Check that the commit type is in the given list (`check_data`, pipe-separated) |
| `require_scope` | Require the commit to have a scope |
| `require_body` | Require a non-empty commit body |
| `require_trailers` | Require the given trailer key(s) to be present (`check_data`, pipe-separated) |
| `allowed_subject_pattern` | Fallback: allow the commit if the subject matches the regex in `check_data` |
| `strict_require_trailers` | Like `require_trailers` but the trailer must also have a non-empty value |

#### `condition_type` values

| Value | Description |
|-------|-------------|
| `by_email` | Apply the rule when the author email contains `rule_data` |
| `by_type` | Apply the rule when the commit type contains `rule_data` |

### Default Rules

When no rules are specified, the following 6 rules are used:

```json
[
  {
    "rule_type": "always_check",
    "check_type": "conventional_commits_base"
  },
  {
    "rule_type": "always_check",
    "check_type": "require_body"
  },
  {
    "rule_type": "conditional_check",
    "check_type": "require_trailers",
    "condition_type": "by_email",
    "rule_data": "uniontech.com",
    "check_data": "Log"
  },
  {
    "rule_type": "conditional_check",
    "check_type": "require_trailers",
    "condition_type": "by_email",
    "rule_data": "uniontech.com",
    "check_data": "Issue|PMS"
  },
  {
    "rule_type": "fallback_check",
    "check_type": "allowed_subject_pattern",
    "rule_data": "Default merge commit caused by workflow sync",
    "check_data": "File Sync from [\\w\\-\\.]+\\/[\\w\\-\\.]+ \\(#[\\d]+\\)"
  },
  {
    "rule_type": "fallback_check",
    "check_type": "allowed_subject_pattern",
    "rule_data": "Default GitHub auto merge commit message",
    "check_data": "Merge (?:pull request #[\\d]+ from [\\w\\-\\.\\/]+|branch '[\\w\\-\\.:\\/]+' into [\\w\\-\\/])"
  }
]
```

### Example: Full custom rule set

```json
[
  {
    "rule_type": "always_check",
    "check_type": "conventional_commits_base"
  },
  {
    "rule_type": "conditional_check",
    "check_type": "require_trailers",
    "condition_type": "by_email",
    "rule_data": "uniontech.com",
    "check_data": "Log"
  },
  {
    "rule_type": "fallback_check",
    "check_type": "allowed_subject_pattern",
    "rule_data": "Allow merge commits",
    "check_data": "Merge (?:pull request #[\\d]+ from [\\w\\-\\.\\/]+)"
  }
]
```

This configuration checks conventional commits format, requires the `Log` trailer for `uniontech.com` authors, and allows merge commits as a fallback — but drops the body requirement and the `Issue|PMS` trailer requirement.
