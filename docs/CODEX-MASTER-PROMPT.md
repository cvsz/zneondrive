# Codex Master Prompt — Unreal Recovery and Production Readiness

> **Operator instructions:** Open this repository in Codex and ask: `Read docs/CODEX-MASTER-PROMPT.md and execute all safely actionable phases. Start with the Unreal manifest blocker.`
>
> This file is an execution brief, **not evidence** that the described tasks have run. Follow root and nested `AGENTS.md` instructions. Report in Thai; keep code, commands, configuration, and identifiers in English.

## Mission

Act as the principal Unreal Engine 5.8 C++ engineer, Go engineer, security engineer, SRE, QA lead, and release engineer for PROJECT: NEON DRIVE. Implement safe, testable improvements toward evidence-backed production readiness. Prioritize the actual Unreal Client/Server blocker, then real builds, packaging, and live gameplay E2E. Continue independent work while external engine prerequisites are unavailable. Do not stop at a planning document.

## Reported operator baseline — verify before use

- The operator reported that PR #38 and #88–#95 were merged and main CI, Web Stack, and CodeQL passed.
- `make ci` reportedly passed, including 42 Python reference tests, Go tests, vet, and build. A Python 3.14 Docker image reportedly built and ran.
- GitHub CLI previously received HTTP 401 because an invalid `GITHUB_TOKEN` overrode the stored `cvsz` credentials. Do not print credentials.
- Candidate: `UE_ROOT=/mnt/zworkforce-storage/dev/UnrealEngine-5.8`, reportedly Unreal Engine 5.8.2 with `Engine/Build/InstalledBuild.txt`.
- `make doctor` and `make client-generate` passed. Android generation was skipped because its SDK is unavailable.
- `make client-build` and `make game-server-build` failed with UnrealBuildTool `OtherCompilationError`, reportedly missing Client/Server precompiled manifests for engine modules such as `Core` and `TraceLog`.
- Logs are on the operator's machine at `artifacts/ue-build/client.log` and `artifacts/ue-build/server.log`; they are **not** committed evidence.
- Reported free disk space was approximately 45.8 GB. Project generation modified four `game/.vscode/` compile-command files.
- No successful real Unreal Client/Server build, package, or packaged gameplay E2E has been reported.

Recheck current HEAD, worktree, open PRs, runner/toolchain availability, engine version, disk space, and actual logs. Do not treat the above as current CI or build proof.

## Architecture and non-negotiable invariants

Preserve Unreal Engine 5.8 Client and dedicated gameplay server, Go 1.27 service plane, PostgreSQL durable authority, and Redis ephemeral coordination/rate limiting. Unreal dedicated servers own gameplay authority; clients are untrusted. Python remains a reference test oracle. Preserve one account → one primary character → one starter vehicle, immutable build revisions, idempotent operations, server-derived rewards, and VIP fairness. Never put server credentials in a client package.

## Phase 0 — Discovery and safety

1. Read `AGENTS.md`, applicable nested instructions, README, ROADMAP, IMPLEMENTATION-CHECKLIST, CHANGELOG, Makefile, `game/README.md`, Unreal tooling, workflows, ADRs, and existing evidence documents.
2. Inspect `git status --short`, `HEAD`, `origin/main`, open PRs, issues, required checks, and recent Actions runs.
3. Inspect the four changed `game/.vscode/` files; preserve user changes and restore generated files only after verifying their provenance.
4. Read both Unreal build logs. Extract the first meaningful failure, exact target/configuration, missing manifest paths, and invoked UnrealBuildTool commands.
5. Create a focused working branch when appropriate. Never reset, clean, or overwrite unrelated work.

## Phase 1 — Unreal engine capability diagnosis

1. Resolve `UE_ROOT`; inspect `Build.version`, `InstalledBuild.txt`, `Build.sh`, UnrealBuildTool, installed-platform metadata, and real precompiled manifest locations.
2. Determine Editor, Linux Client, Linux Server, and Linux packaging capabilities independently.
3. Distinguish Editor-only installed builds, incomplete installations, unsupported target/configuration combinations, stale intermediates, and genuine game-project defects.
4. Do not invent/copy manifests, disable UnrealBuildTool validation, or claim that `InstalledBuild.txt` proves Client/Server support.
5. Produce an evidence-backed capability matrix and exact remediation prerequisites.

## Phase 2 — Safe recovery

Prefer an authorized local UE 5.8 installation already containing valid Client/Server support. Otherwise inspect authorized source trees and existing build caches. Before any source/BuildGraph build, check licensing, toolchain, memory, swap, free disk, and output capacity. Preserve the existing installed engine and use separate outputs. Do not launch a full engine rebuild with insufficient resources or delete user data to free space. If blocked, report the exact missing prerequisite and continue Phase 6.

## Phase 3 — Implement preflight diagnostics

Improve the canonical `tools/ue-linux.sh`, `make doctor`, or other existing preflight only where appropriate. Detect actual target support and missing manifests early without rejecting valid alternative engine layouts. Add deterministic tests/fixtures for source trees, Editor-only installed builds, Client/Server-capable builds, and incomplete installations. Avoid hardcoded host paths. Validate that the reported failure produces actionable diagnostics.

## Phase 4 — Real builds and packaging

Once a compatible engine is available, generate project files, build the actual Linux Client and dedicated Server, and fix genuine C++/target/configuration failures. Then separately cook, stage, and archive Client and Server using canonical Makefile tooling. Retain sanitized full logs, exit codes, engine/toolchain metadata, artifact inventories, sizes, and SHA-256 checksums. Check that packages contain expected executable artifacts and no secrets. Mark each gate complete only after real successful execution.

## Phase 5 — Packaged gameplay E2E

Using isolated PostgreSQL, Redis, Go API, a packaged Unreal dedicated server, and a packaged client where graphics are available, verify:

- Fresh account, primary character, and starter vehicle bootstrap.
- Short-lived single-use gameplay ticket issuance/redemption and authority-only durable binding.
- Client input → server RPC → authoritative movement → replication; disconnect/reconnect.
- Garage 17, inventory/blueprints, rebuild, MQ001–MQ012, First Ignition, and Roadworthy where implemented.
- Race registration bound to active immutable build revision/hash, ordered checkpoints, finish, idempotent replay, and durable results.
- Negative paths: forged rewards/ownership, duplicate tickets, stale revisions, invalid checkpoints, and impossible movement.

Retain sanitized reproducible evidence. Do not substitute mocks or static validators for live packaged proof, or claim unimplemented gameplay was exercised.

## Phase 6 — Independent work while engine is blocked

- Inspect all `github/codeql-action@v3` references and the reported December 2026 deprecation. Verify current supported version, runner compatibility, permissions, SARIF behavior, and required checks; migrate in a focused PR when safe.
- Check that Redis Sentinel majority v5.3 documentation is indexed consistently. Do not close broad HA, external fencing, asymmetric-partition, zero-RPO, regional DR, or production-scale gates without proof.
- Improve isolated E2E orchestration using existing Makefile/Compose tools: bounded waits, health checks, sanitized logs, deterministic cleanup, and retained failure artifacts.
- Run relevant Go, Python, PostgreSQL, Redis, security, and documentation tests. Fix concrete failures without unnecessary refactoring.

## Phase 7 — Operational evidence

When infrastructure permits, exercise bounded load/soak, private metrics, credential-safe logs, isolated PostgreSQL backup/restore, Redis Sentinel recovery, rollback, and deployment failure handling. Measure actual RPO/RTO; retain evidence and keep unverified production gates open.

## Phase 8 — Validation, PRs, and release governance

For each focused increment, implement code/tests/docs, run relevant validation (including `make ci` when appropriate), inspect the complete diff for secrets/generated artifacts, commit, and open a focused PR when authorized. Use existing `cvsz` authentication without printing tokens; diagnose invalid environment overrides without changing unrelated credentials. Verify exact PR head SHA, required checks, and review requirements. Merge only when all requirements pass; never force-merge or bypass gates. Verify post-merge main and CI. If authentication or infrastructure is unavailable, preserve a clean local branch and report precise publication blockers.

Synchronize README, ROADMAP, IMPLEMENTATION-CHECKLIST, CHANGELOG, and documentation indexes only to the level justified by real evidence. Never rewrite a long checklist from truncated excerpts.

## Execution rules and completion criteria

Execute all safely actionable phases in dependency order. Prefer source changes, tests, and retained evidence over prose-only changes. Do not redownload or redistribute licensed engine assets without authorization; do not perform destructive cleanup or production database operations. If a phase is blocked, record the exact prerequisite and continue an independent task. Do not declare production-ready until real packaged gameplay, security, operational, and deployment gates are evidenced.

**Final report in Thai:** initial/final SHA and worktree; engine capability matrix; recovery path; exact files changed; Client/Server build/package/E2E commands, exit codes, artifacts and hashes; tests and CI run URLs; PRs and merge SHAs; remaining P0/P1/P2 blockers; and one next highest-priority executable action.

**Start now:** inspect the actual Unreal logs and engine manifests before retrying compilation.
