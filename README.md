## Conventional Commits Checker

A GitHub Action that checks commits inside a Pull Request against the [Conventional Commits](https://www.conventionalcommits.org/) guideline.

### Usage

```yaml
- name: Call Conventional Commits Checker
  uses: linuxdeepin/action-conventionalcommits-checker@master
```

### Inputs

All inputs are optional and default to `"true"`. When all inputs are left at their defaults, the behavior is identical to the previous hardcoded version (backward compatible).

| Input | Description | Default |
|-------|-------------|---------|
| `github_token` | Action GitHub token | — |
| `check-conventional-commits` | Check that the commit subject follows `type: description` format | `"true"` |
| `check-require-body` | Check that commits include a body | `"true"` |
| `check-require-log-trailer` | Check that `uniontech.com` commits include a `Log` trailer | `"true"` |
| `check-require-issue-trailer` | Check that `uniontech.com` commits include an `Issue` or `PMS` trailer | `"true"` |
| `allow-file-sync-commits` | Allow file sync merge commits | `"true"` |
| `allow-merge-commits` | Allow GitHub auto merge commits | `"true"` |

### Examples

Use all default checks (backward compatible):

```yaml
- name: Call Conventional Commits Checker
  uses: linuxdeepin/action-conventionalcommits-checker@master
```

Disable specific checks (e.g. skip body requirement and merge commit allowance):

```yaml
- name: Call Conventional Commits Checker
  uses: linuxdeepin/action-conventionalcommits-checker@master
  with:
    check-require-body: "false"
    allow-merge-commits: "false"
```
