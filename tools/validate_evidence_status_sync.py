#!/usr/bin/env python3
"""Fail CI when top-level production-readiness evidence status drifts.

This validator checks promoted evidence status rather than treating every test-harness
increment as a runtime milestone. It prevents documentation from silently closing
real UE build/package/live-integration gates before retained evidence exists, while
also requiring post-baseline reference-parity evidence to remain discoverable.
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

    # Post-baseline reference-parity evidence is test-harness evidence, not a promoted
    # runtime milestone, but it must remain discoverable from the canonical docs index.
    require(docs_index, "runtime-reference-parity-v2.4.md", "docs index parity v2.4")
    require(docs_index, "runtime-reference-parity-integrity-v2.5.md", "docs index parity v2.5")
    require(docs_index, "runtime-reference-operation-parity-v2.6.md", "docs index parity v2.6")
    require(docs_index, "runtime-reference-rebuild-parity-v2.7.md", "docs index parity v2.7")
    require(docs_index, "runtime-reference-parity-index-v2.8.md", "docs index parity evidence index v2.8")
    require(checklist, "[ ] full reference-oracle parity in Go", "full parity remains open")

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
