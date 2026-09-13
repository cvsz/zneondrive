#!/usr/bin/env python3
"""Fail CI when top-level production-readiness evidence status drifts.

This validator intentionally checks evidence language, not implementation maturity.
It prevents documentation from silently promoting UE build/package/live integration
claims before retained evidence exists and ensures newer reference-parity increments
are recorded without changing the evidence-gated runtime phase.
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

    # Reference-parity increments are evidence improvements, not runtime-phase promotion.
    for version in ("v2.4", "v2.5", "v2.6"):
        require(changelog, version, f"CHANGELOG {version} evidence entry")

    print("evidence status synchronization: ok")


if __name__ == "__main__":
    main()
