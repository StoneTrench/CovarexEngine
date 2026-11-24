# Contributing

Feel free to contribute, but since I'm only one person it could take awhile to approve or respond to a PR.

## Fork and Pull-Request

1. Fork the repository to your GitHub account.
1. Clone your fork to your local machine.
1. Optionally you can create a feature branch locally.
1. Make your changes. Make sure to follow the project's coding standards.
1. Commit and push your changes to your fork's main branch.
1. Open a pull request (PR) against the main repository.
1. If your PR contains multiple commits or big changes, you should include extra information in the PR description covering the following:
- Problem: What was problem that was solved?
- Solution: What is the solution that was implemented?
- Open Issues: Which related issues does this PR address? Link them where applicable.

A detailed description is not required for a simple PR containing only minor changes.

## Commit messages

Ignoring the mess that the old commit history is, please stick to the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/#specification) specification from this point forward.

Here's a very rough syntax description of it in EBNF:

```ebnf
COMMIT_MESSAGE = TYPE, [ '(', SCOPE, ')' ], [ '!' ], ':', MESSAGE;
```

Types can be:
- `build`: Things relating to the build system.
- `chore`: Maintenance tasks.
- `docs`: Updates to documentation.
- `feat`: New feature.
- `fix`: Bug fixes.
- `perf`: Performance related change.
- `refactor`: Change that improves code quality.
- `revert`: Reverts a previous commit.
- `style`: Formatting, whitespaces.
- `test`: Adding or fixing tests.

The `!` marker means that this is a breaking change.

Examples:

```r
refactor(systems/modding)!: remove schema version
feat(systems/workers): implement worker pool
docs(platform/vfs): add documentation to functions
revert: refactor(systems/modding): remove schema version
```

Use the imperative mood in commit descriptions (add, remove, implement).

## Versioning

```r
<MAJOR>.<MINOR>.<PATCH>[-rc.<RELEASE_CANDIDATE>]
```

Example:
- Release candidate: `9.12.951-rc.3`
- Stable release: `9.12.951`

Patch does not increment for every release candidate, or the release for a single cycle.