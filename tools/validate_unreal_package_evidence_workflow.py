#!/usr/bin/env python3
from pathlib import Path
import sys

ROOT = Path(__file__).resolve().parents[1]
WORKFLOW = ROOT / ".github/workflows/unreal-package-evidence.yml"

REQUIRED = [
    "name: Unreal Package Evidence",
    "workflow_dispatch:",
    "runs-on: [self-hosted, linux, unreal-5.8]",
    "bash tools/ue-linux.sh info",
    "bash tools/ue-linux.sh generate",
    "bash tools/ue-linux.sh package-client",
    "bash tools/ue-linux.sh package-server",
    "for kind in client server",
    'dist/packages/${kind}-linux',
    "client-manifest.txt",
    "server-manifest.txt",
    "client-sha256sums.txt",
    "server-sha256sums.txt",
    "Assert package evidence exists",
    "if: success()",
    "if: always()",
    "actions/upload-artifact@v4",
    "retention-days: 30",
]
FORBIDDEN = [
    "GAME_SERVER_SHARED_KEY",
    "ZNEON_GAME_SERVER_KEY",
    "DATABASE_URL",
    "REDIS_URL",
]


def main() -> int:
    errors: list[str] = []
    if not WORKFLOW.is_file():
        errors.append("missing .github/workflows/unreal-package-evidence.yml")
    else:
        text = WORKFLOW.read_text(encoding="utf-8")
        for token in REQUIRED:
            if token not in text:
                errors.append(f"package evidence workflow missing {token!r}")
        for token in FORBIDDEN:
            if token in text:
                errors.append(f"package evidence workflow must not contain runtime secret/data-plane token {token!r}")
        if text.count('dist/packages/${kind}-linux') != 1:
            errors.append("package inventory must derive exactly one root per explicit client/server kind")
        if text.count("bash tools/ue-linux.sh package-client") != 1:
            errors.append("workflow must invoke the Client package command exactly once")
        if text.count("bash tools/ue-linux.sh package-server") != 1:
            errors.append("workflow must invoke the dedicated Server package command exactly once")

    if errors:
        print("Unreal package evidence workflow validation FAILED", file=sys.stderr)
        for error in errors:
            print(f"- {error}", file=sys.stderr)
        return 1
    print("Unreal package evidence workflow validation OK")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
