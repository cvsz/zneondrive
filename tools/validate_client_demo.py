#!/usr/bin/env python3
"""Smoke-check the client presentation package with Python stdlib only."""

from __future__ import annotations

import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
DEMO = ROOT / "demo"

REQUIRED_FILES = [
    DEMO / "index.html",
    DEMO / "styles.css",
    DEMO / "app.js",
    DEMO / "README.md",
    ROOT / "client" / "EXECUTIVE-BRIEF-TH.md",
    ROOT / "client" / "PRESENTATION-TALK-TRACK-TH.md",
    ROOT / "client" / "DEMO-RUNBOOK-TH.md",
    ROOT / "client" / "COMMERCIAL-SCOPE-TH.md",
    ROOT / "client" / "CLIENT-QA-TH.md",
    ROOT / "docs" / "adr" / "0004-unreal-engine-client-and-gameplay-server.md",
    ROOT / "docs" / "adr" / "0005-durable-service-plane.md",
    ROOT / "docs" / "adr" / "0006-phase-gated-deployment.md",
]

REQUIRED_SECTIONS = [
    "hero",
    "world",
    "story",
    "vehicle",
    "slice",
    "architecture",
    "roadmap",
]

FORBIDDEN_PLACEHOLDERS = ("TODO", "FIXME", "TBD")


def fail(message: str) -> None:
    raise AssertionError(message)


def main() -> int:
    for path in REQUIRED_FILES:
        if not path.is_file():
            fail(f"missing client-readiness file: {path.relative_to(ROOT)}")

    html = (DEMO / "index.html").read_text(encoding="utf-8")
    css = (DEMO / "styles.css").read_text(encoding="utf-8")
    js = (DEMO / "app.js").read_text(encoding="utf-8")

    for section in REQUIRED_SECTIONS:
        if f'id="{section}"' not in html:
            fail(f"missing demo section #{section}")

    if 'href="./styles.css"' not in html:
        fail("demo stylesheet is not linked")
    if 'src="./app.js"' not in html:
        fail("demo JavaScript is not linked")

    combined = "\n".join((html, css, js))
    for placeholder in FORBIDDEN_PLACEHOLDERS:
        if re.search(rf"\b{placeholder}\b", combined, flags=re.IGNORECASE):
            fail(f"presentation source contains placeholder marker: {placeholder}")

    if re.search(r'https?://', combined):
        fail("demo source must remain offline-first and contain no external HTTP dependencies")

    required_copy = [
        "100",
        "PROJECT DRIVE ZERO",
        "Unreal Engine 5.8",
        "PostgreSQL + Redis",
        "Presentation prototype",
    ]
    for text in required_copy:
        if text not in html:
            fail(f"demo is missing required client-facing copy: {text}")

    if "MQ001" not in js or "MQ012" not in js:
        fail("vertical-slice demo must cover MQ001 through MQ012")

    print("client demo validation OK: offline assets, required sections, stack and slice present")
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except AssertionError as exc:
        print(f"client demo validation FAILED: {exc}", file=sys.stderr)
        raise SystemExit(1)
