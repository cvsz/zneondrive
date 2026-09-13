#!/usr/bin/env python3
"""Validate the PROJECT: NEON DRIVE documentation baseline and relative links."""

from __future__ import annotations

import re
import sys
from pathlib import Path
from urllib.parse import unquote

from validate_evidence_status_sync import main as validate_evidence_status_sync

ROOT = Path(__file__).resolve().parents[1]

REQUIRED_DOCS = [
    "README.md",
    "ABOUT.md",
    "SUPPORT.md",
    "GOVERNANCE.md",
    "MAINTAINERS.md",
    "SECURITY.md",
    "CONTRIBUTING.md",
    "CODE_OF_CONDUCT.md",
    "CHANGELOG.md",
    "ROADMAP.md",
    "IMPLEMENTATION-CHECKLIST.md",
    "docs/README.md",
    "docs/product-requirements.md",
    "docs/architecture.md",
    "docs/api-contract.md",
    "docs/data-model.md",
    "docs/multiplayer-networking.md",
    "docs/threat-model.md",
    "docs/testing-strategy.md",
    "docs/performance-budget.md",
    "docs/observability-slo.md",
    "docs/deployment.md",
    "docs/backup-restore-dr.md",
    "docs/incident-response.md",
    "docs/accessibility.md",
    "docs/localization.md",
    "docs/content-pipeline.md",
    "docs/economy-monetization.md",
    "docs/moderation-safety.md",
    "docs/privacy-data-retention.md",
    "docs/live-ops.md",
    "docs/release.md",
    "docs/release-readiness-checklist.md",
    "docs/asset-ip-policy.md",
    "docs/glossary.md",
    "docs/brand-guide.md",
    "docs/quest-authoring.md",
    "docs/vehicle-physics-and-race-integrity.md",
    "docs/world-streaming-and-districts.md",
    "docs/factions-relationships.md",
    "docs/crew-social.md",
    "docs/pets-companions.md",
    "docs/adr/README.md",
    "assets/zneondrive-banner.png",
]

MARKDOWN_LINK = re.compile(r"!?(?:\[[^\]]*\])\(([^)]+)\)")
HTML_LINK = re.compile(r"""(?:href|src)=["']([^"']+)["']""", re.IGNORECASE)

IGNORE_PREFIXES = (
    "http://",
    "https://",
    "mailto:",
    "data:",
    "javascript:",
    "#",
)


def normalized_target(raw: str) -> str:
    raw = raw.strip()
    if raw.startswith("<") and raw.endswith(">"):
        raw = raw[1:-1].strip()

    if ' "' in raw:
        raw = raw.split(' "', 1)[0].strip()
    elif " '" in raw:
        raw = raw.split(" '", 1)[0].strip()

    raw = raw.split("#", 1)[0].split("?", 1)[0]
    return unquote(raw.strip())


def validate_link(source: Path, raw: str) -> str | None:
    target = normalized_target(raw)
    if not target or target.startswith(IGNORE_PREFIXES):
        return None

    if target.startswith("//") or "://" in target:
        return None

    if target.startswith("/"):
        resolved = ROOT / target.lstrip("/")
    else:
        resolved = source.parent / target

    try:
        resolved.relative_to(ROOT)
    except ValueError:
        return f"{source.relative_to(ROOT)}: link escapes repository: {raw}"

    if not resolved.exists():
        return f"{source.relative_to(ROOT)}: missing relative link target: {raw}"
    return None


def main() -> int:
    errors: list[str] = []

    for rel in REQUIRED_DOCS:
        if not (ROOT / rel).exists():
            errors.append(f"missing required documentation asset: {rel}")

    markdown_files = [path for path in ROOT.rglob("*.md") if ".git" not in path.parts]

    checked = 0
    for path in markdown_files:
        text = path.read_text(encoding="utf-8")
        for raw in MARKDOWN_LINK.findall(text) + HTML_LINK.findall(text):
            checked += 1
            error = validate_link(path, raw)
            if error:
                errors.append(error)

    if errors:
        print("documentation validation FAILED", file=sys.stderr)
        for error in errors:
            print(f"- {error}", file=sys.stderr)
        return 1

    validate_evidence_status_sync()

    print(
        f"documentation validation OK: "
        f"{len(REQUIRED_DOCS)} required assets, "
        f"{len(markdown_files)} markdown files, "
        f"{checked} links checked"
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
