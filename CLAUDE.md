# CLAUDE.md

This file is read automatically at the start of every Claude Code
session in this repo. It's the source of truth for how to work here —
if anything below conflicts with a one-off instruction in chat, ask
before assuming which one wins.

## What this project is

Formula Engine: a stateless Go library that evaluates spreadsheet-style
formulas (`MAX(a, b)`, `IF(cond, x, y)`) against supplied data and
returns results. No storage, no HTTP service, no UI — it's imported
and called directly, like `lodash` but for formula evaluation instead
of data utilities.

## Docs structure

This repo tracks **one PRD + RFC per feature**, not one global set of
documents for the whole project. This is what lets new features get
added later without editing or overloading old documents.

```
docs/
├── features/
│   ├── 0001-core-engine/
│   │   ├── prd.md              # product requirements
│   │   └── rfc.md              # technical design
│   ├── 0002-<next-feature-slug>/
│   │   ├── prd.md
│   │   └── rfc.md
│   └── ...
└── adr/
    ├── 0001-iterative-dependency-resolution.md
    └── ...   (global — decisions can apply across features)
```

- Feature folders are numbered sequentially (`0001`, `0002`, ...) in
  the order they were started, regardless of whether earlier ones are
  finished.
- **`prd.md` and `rfc.md` are for human readers** (the user, other
  engineers) — no agent-execution instructions belong in them.
  Execution instructions (how to break a feature into tickets, work
  them, when to write an ADR) are handled separately, not stored as a
  file this pipeline auto-generates.
- ADRs are **global**, not per-feature — a decision made while
  building one feature (e.g. "dependency resolution must be
  iterative") may constrain or inform later features too.

## When asked to build a new feature

Follow this pipeline, in order — don't skip straight to coding:

1. **Check if a feature folder already exists** for this request
   (`docs/features/*/prd.md` — search by topic). If yes, this feature
   already has documentation — don't duplicate it; ask the user how
   they'd like to proceed. If no, this is a new feature — go to
   step 2.
2. **Write a PRD** at `docs/features/000N-<slug>/prd.md` (next
   sequential number). Cover: problem statement, product concept,
   goals, non-goals, use cases, functional & non-functional
   requirements, success metrics, open questions. Keep it
   implementation-free — no code, no folder structure, no algorithm
   detail. Use existing PRDs in this repo as a style/format
   reference, but don't copy their specific content.
3. **Write an RFC** at `docs/features/000N-<slug>/rfc.md` in the same
   folder. This is a **pure technical design document for other
   software engineers** — public API, package structure, algorithms,
   alternatives considered, and an "Implementation Phases" section at
   the end (sequential phases with goal/modules/definition-of-done/
   estimate — not individual tickets). Do **not** put agent
   instructions in the RFC.
4. **Confirm the PRD and RFC with the user before writing any code**,
   unless they've explicitly said to proceed straight through.
5. Stop there. Ticket breakdown and implementation happen when the
   user separately provides execution instructions for that feature —
   don't start generating tickets or writing code on your own
   initiative just because a PRD/RFC exists.

## Git workflow

- **`develop`** is the default branch. All active work happens here
  or on short-lived ticket branches — never commit directly to
  `main`.
- **One GitHub issue = one branch = one merge.** Don't use a single
  long-lived branch per feature folder — each ticket gets its own
  short-lived branch off `develop`, its own PR, and its own merge.
  This keeps every reviewable diff scoped to one ticket's acceptance
  criteria instead of piling up into one large review at the end of
  a feature.
  ```
  git checkout develop
  git pull
  git checkout -b issue-<N>-<short-slug>
  ```
  e.g. `issue-5-ast-node-types` for issue #5
  (`[FASE-3.1] AST node types + literal/identifier/binary expression
  parsing`).
- **Per ticket:** implement, commit with a message referencing the
  issue (e.g. `closes #5`), run the full test suite (`make test &&
  make test-race`) — do not open a PR if anything fails. Push the
  branch and open a PR against `develop` (`gh pr create`). Merging a
  passing per-ticket PR into `develop` is routine ticket work, not a
  release action (see the guardrail below) — it doesn't need a
  separate confirmation each time. Delete the branch after merge.
- **`main` is release-only.** It only ever receives merges from
  `develop` (a "cut a release" action), never direct commits and
  never a merge straight from a ticket branch.
- **On merge to `main` (releasing a version):**
    1. Run the full test suite, including the race detector:
       `make test && make test-race`. Do not proceed if anything fails.
    2. Decide the version bump using semver, based on what changed
       since the last tag:
        - **MAJOR** — breaking change to the public API (`formula.go`'s
          exported types/functions)
        - **MINOR** — new backwards-compatible functionality (e.g. a new
          built-in function, a new feature folder shipped)
        - **PATCH** — bug fix, internal-only change, no public API
          change
    3. Tag the merge commit and push the tag — this is how Go modules
       are versioned; there's no separate build/publish step for a
       library beyond this:
       ```
       git tag vX.Y.Z
       git push origin vX.Y.Z
       ```
       Consumers then pin to it via `go get module/path@vX.Y.Z`.
    4. Update `CHANGELOG.md` with the version and a summary of what's
       included, in the same commit as the merge if possible.
- **Guardrail:** merging into `main` and creating a version tag are
  release actions, not routine ticket work — confirm with the user
  before doing either, even if everything above passes cleanly.
- **GitHub's default branch is `develop`** (set via `gh repo edit
  --default-branch develop`) — this is what makes `closes #N` in a
  per-ticket PR auto-close the issue on merge, since GitHub only
  honors that keyword when merging into the repository's default
  branch.

## Code conventions

- Package layout is fixed by the relevant feature's RFC — don't
  restructure folders without discussing it first.
- Public API surface stays small: only what's exported from the
  repo's top-level package(s) is public. Everything else lives under
  `internal/` and must stay there — if a helper feels like it should
  be public, raise it as a design question (possibly a new ADR),
  don't just export it.
- File naming: `snake_case.go`, per Go convention.
- Each built-in function (math/logic/text/date/comparison, or whatever
  categories a future feature adds) gets its own file and its own
  test file — don't bundle multiple functions into one file.
- Follow the `Function` interface pattern already established (see
  `internal/registry/function.go`) for any new built-in function;
  don't invent a parallel pattern without an ADR justifying why.
- Dependency resolution and cycle detection must stay **iterative**
  (explicit stack/queue), never recursive — this is a cross-feature
  correctness requirement (see the relevant ADR once written), not a
  style preference.
- Type coercion behavior must match the coercion table in the core
  engine's RFC exactly. If a new feature needs a behavior the table
  doesn't cover, update the table first, then implement.

## Environment

- Pinned Go version: `1.25`, via the `golang:1.25-bookworm` container
  image. `go.mod` declares `go 1.25`.
- All builds and tests run inside that container, not against
  whatever Go version is installed on the host — use the `Makefile`
  targets (`make test`, `make test-race`, `make build`, `make tidy`,
  `make fmt`, `make vet`, `make shell`) rather than invoking `go`
  directly. This applies whether you're a human contributor or Claude
  Code — don't shell out to a host-installed `go` binary.

## Testing

- Every acceptance criterion in a ticket must map to an actual test,
  not just manual verification.
- Run `make test` (full suite, containerized) before every commit,
  not just the package you touched.
- Use `make test-race` for anything touching shared state or
  concurrency.
- Prefer table-driven tests for functions with many input variations.

## Commits

- One ticket ≈ one commit, on that ticket's own branch (see Git
  workflow above). Don't bundle multiple tickets into one commit, and
  don't split one ticket across many commits that leave
  intermediate commits in a broken/non-compiling state.
- Commit message format: short summary line, then `closes #N` if it
  finishes a GitHub issue.
- Update `CHANGELOG.md` in the same commit as the work it describes.
- When a ticket adds or changes user-visible functionality (a new
  built-in function, new supported syntax, a new error case a caller
  can observe, etc.), also update `README.md`'s supported/not-yet-
  supported lists in that same commit. `CHANGELOG.md` is a history of
  what changed; `README.md` must stay an accurate snapshot of what
  the library can do *right now* — don't let it drift stale as phases
  complete. Routine internal-only tickets (no observable behavior
  change) don't need a README update.

## What NOT to do

- Don't add a network layer, database, or persistence of any kind
  without a new PRD/RFC first — the core design is stateless
  throughout. If a future feature genuinely needs this, that's a new
  RFC, not a quiet addition to existing code.
- Don't add a public API for consumer-registered custom functions
  without a new PRD/RFC — currently out of scope.
- Don't write ADRs for routine tickets (e.g. "add the ROUND
  function") — reserve them for decisions with real tradeoffs.
- Don't manipulate commit timestamps or batch large amounts of work
  into one sitting — commit history should reflect genuine,
  incremental progress.
- Don't start writing code for a new feature before its PRD and RFC
  exist and have been confirmed — see the pipeline above.
- Don't mix agent-execution instructions into `prd.md` or `rfc.md`.
- Don't generate GitHub issues or start implementation on your own
  initiative right after writing a PRD/RFC — wait for separate
  execution instructions from the user.