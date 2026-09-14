#!/usr/bin/env python3
"""Fail CI when top-level production-readiness evidence status drifts.

This validator checks promoted evidence status rather than treating every test-harness
increment as a runtime milestone. It prevents documentation from silently closing
real UE build/package/live-integration gates before retained evidence exists, while
also requiring post-baseline evidence increments to remain discoverable.
"""
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]


def read(path: str) -> str:
    return (ROOT / path).read_text(encoding="utf-8")


def require(text: str, needle: str, label: str) -> None:
    if needle not in text:
        raise SystemExit(f"status-sync failure: {label}: missing {needle!r}")


def main() -> None:
    readme = read("README.md")
    roadmap = read("ROADMAP.md")
    checklist = read("IMPLEMENTATION-CHECKLIST.md")
    changelog = read("CHANGELOG.md")
    docs_index = read("docs/README.md")

    # The promoted repository status is still Phase 4.19 / Runtime v2.3.
    require(readme, "Phase-4.19%20Unreal%20Rotation%20Envelope", "README phase badge")
    require(readme, "Runtime-v2.3", "README runtime badge")
    require(readme, "Phase 4.19 Unreal Rotation Envelope v2.3", "README current phase")
    require(roadmap, "Phase 4 — Runtime prototype / integration v2.3", "ROADMAP runtime phase")
    require(checklist, "Unreal rotation-envelope telemetry v2.3", "checklist promoted evidence")
    require(changelog, "Runtime Unreal Authority Rotation Envelope v2.3", "CHANGELOG promoted evidence")

    # Post-baseline evidence is not a promoted gameplay-runtime milestone, but it must
    # remain discoverable and historically recorded.
    require(docs_index, "runtime-reference-parity-v2.4.md", "docs index parity v2.4")
    require(docs_index, "runtime-reference-parity-integrity-v2.5.md", "docs index parity v2.5")
    require(docs_index, "runtime-reference-operation-parity-v2.6.md", "docs index parity v2.6")
    require(docs_index, "runtime-reference-rebuild-parity-v2.7.md", "docs index parity v2.7")
    require(docs_index, "runtime-reference-parity-index-v2.8.md", "docs index parity evidence index v2.8")
    require(docs_index, "runtime-reference-race-parity-v2.9.md", "docs index race parity v2.9")
    require(changelog, "Runtime Reference Race Input Parity v2.9", "CHANGELOG race parity v2.9")
    require(docs_index, "runtime-reference-race-result-parity-v3.5.md", "docs index race result parity v3.5")
    require(changelog, "Runtime Reference Race Result Parity v3.5", "CHANGELOG race result parity v3.5")
    require(checklist, "race-result hash reference parity v3.5", "checklist race result parity v3.5")
    require(docs_index, "runtime-reference-entitlement-parity-v3.6.md", "docs index entitlement parity v3.6")
    require(changelog, "Runtime Reference Entitlement Parity v3.6", "CHANGELOG entitlement parity v3.6")
    require(checklist, "entitlement fairness reference parity v3.6", "checklist entitlement parity v3.6")
    require(docs_index, "runtime-reference-quest-sequence-parity-v3.7.md", "docs index quest sequence parity v3.7")
    require(changelog, "Runtime Reference Quest Sequence Parity v3.7", "CHANGELOG quest sequence parity v3.7")
    require(checklist, "quest-sequence reference parity v3.7", "checklist quest sequence parity v3.7")
    require(docs_index, "runtime-reference-quest-reward-parity-v3.8.md", "docs index quest reward parity v3.8")
    require(changelog, "Runtime Reference Quest Reward Parity v3.8", "CHANGELOG quest reward parity v3.8")
    require(checklist, "quest-reward reference parity v3.8", "checklist quest reward parity v3.8")
    require(docs_index, "runtime-reference-quest-reward-authority-v3.9.md", "docs index quest reward authority v3.9")
    require(changelog, "Runtime Reference Quest Reward Authority v3.9", "CHANGELOG quest reward authority v3.9")
    require(checklist, "quest-reward authority reference v3.9", "checklist quest reward authority v3.9")
    require(docs_index, "runtime-reference-race-lifecycle-parity-v4.0.md", "docs index race lifecycle parity v4.0")
    require(changelog, "Runtime Reference Race Lifecycle Parity v4.0", "CHANGELOG race lifecycle parity v4.0")
    require(checklist, "race-lifecycle reference parity v4.0", "checklist race lifecycle parity v4.0")
    require(docs_index, "runtime-reference-race-start-parity-v4.1.md", "docs index race start parity v4.1")
    require(changelog, "Runtime Reference Race Start Parity v4.1", "CHANGELOG race start parity v4.1")
    require(checklist, "race-start reference parity v4.1", "checklist race start parity v4.1")
    require(docs_index, "runtime-reference-race-operation-parity-v4.2.md", "docs index race operation parity v4.2")
    require(changelog, "Runtime Reference Race Operation Parity v4.2", "CHANGELOG race operation parity v4.2")
    require(checklist, "race-operation replay reference parity v4.2", "checklist race operation parity v4.2")
    require(docs_index, "runtime-reference-quest-side-effects-parity-v4.3.md", "docs index quest side effects parity v4.3")
    require(changelog, "Runtime Reference Quest Side Effects Parity v4.3", "CHANGELOG quest side effects parity v4.3")
    require(checklist, "quest side-effects reference parity v4.3", "checklist quest side effects parity v4.3")
    require(docs_index, "runtime-reference-build-operation-replay-v4.4.md", "docs index build operation replay v4.4")
    require(changelog, "Runtime Reference Build Operation Replay v4.4", "CHANGELOG build operation replay v4.4")
    require(checklist, "rebuild-operation replay hardening v4.4", "checklist build operation replay v4.4")
    require(docs_index, "runtime-reference-quest-operation-scope-v4.5.md", "docs index quest operation scope v4.5")
    require(changelog, "Runtime Reference Quest Operation Scope v4.5", "CHANGELOG quest operation scope v4.5")
    require(checklist, "quest-operation scope hardening v4.5", "checklist quest operation scope v4.5")
    require(docs_index, "runtime-service-restart-recovery-v4.6.md", "docs index service restart recovery v4.6")
    require(changelog, "Runtime Service Restart Recovery Evidence v4.6", "CHANGELOG service restart recovery v4.6")
    require(checklist, "service-process restart/reconnect recovery v4.6", "checklist service restart recovery v4.6")
    require(docs_index, "runtime-postgres-promotion-recovery-v4.7.md", "docs index PostgreSQL promotion recovery v4.7")
    require(changelog, "Runtime PostgreSQL Promotion Recovery Evidence v4.7", "CHANGELOG PostgreSQL promotion recovery v4.7")
    require(checklist, "PostgreSQL physical-standby promotion recovery v4.7", "checklist PostgreSQL promotion recovery v4.7")
    require(checklist, "[ ] process/node reconnect recovery", "node reconnect recovery remains open")
    require(checklist, "[ ] duplicate-reward prevention under failover", "failover reward safety remains open")
    require(checklist, "[ ] full reference-oracle parity in Go", "full parity remains open")

    require(docs_index, "runtime-postgres-statements-v3.4.md", "docs index PostgreSQL statement metrics v3.4")
    require(changelog, "Runtime PostgreSQL Statement Metrics Evidence v3.4", "CHANGELOG PostgreSQL statement metrics v3.4")
    require(
        checklist,
        "[x] PostgreSQL query-level/external exporter metrics source + PostgreSQL 17 CI evidence (deployment verification still open)",
        "checklist PostgreSQL statement/exporter CI evidence",
    )
    require(checklist, "[ ] Deployed metrics scrape/dashboard evidence", "deployed metrics gate remains open")

    # Public status must remain evidence-gated until retained real-engine evidence exists.
    require(readme, "Production%20Readiness-Evidence%20Gated", "README production badge")
    require(readme, "UE%20Source%20Build-Evidence%20Pending", "README UE evidence badge")
    require(roadmap, "[ ] Produce successful retained UE 5.8 Client/Server build evidence", "ROADMAP UE build gate")
    require(roadmap, "[ ] Produce successful retained UE 5.8 Client/Server package/cook evidence", "ROADMAP UE package gate")
    require(roadmap, "[ ] Produce live Unreal ↔ Go packaged/session E2E evidence", "ROADMAP live integration gate")
    require(checklist, "[ ] successful real UE 5.8 Client build evidence artifact", "checklist Client build gate")
    require(checklist, "[ ] successful real UE 5.8 dedicated Server build evidence artifact", "checklist Server build gate")
    require(checklist, "[ ] successful retained UE 5.8 Client/Server package/cook evidence", "checklist package gate")
    require(checklist, "[ ] live packaged Unreal ↔ Go integration evidence", "checklist live integration gate")

    print("evidence status synchronization: ok")


if __name__ == "__main__":
    main()
