#!/usr/bin/env python3
from pathlib import Path
import sys

ROOT = Path(__file__).resolve().parents[1]
WORKFLOW = ROOT / ".github/workflows/unreal-source-build.yml"


def main() -> int:
    text = WORKFLOW.read_text(encoding="utf-8")
    errors: list[str] = []

    required = {
        "manual trigger": "workflow_dispatch:",
        "self-hosted UE runner": "runs-on: [self-hosted, linux, unreal-5.8]",
        "runner preflight": "bash tools/ue-linux.sh preflight",
        "runner disk floor": 'UE_MIN_FREE_GB: "20"',
        "runner preflight evidence": "runner-preflight.txt",
        "installed-build info": "bash tools/ue-linux.sh info",
        "project generation": "bash tools/ue-linux.sh generate",
        "client build": "bash tools/ue-linux.sh build-target NeonDriveClient",
        "server build": "bash tools/ue-linux.sh build-target NeonDriveServer",
        "evidence directory": ".runtime/ue-build-evidence",
        "build-version evidence": "Engine/Build/Build.version",
        "checksums": "sha256sum",
        "evidence upload": "actions/upload-artifact@v4",
        "retention": "retention-days: 30",
        "failure evidence": "if: always()",
    }
    for label, token in required.items():
        if token not in text:
            errors.append(f"{label}: missing {token!r}")

    forbidden = [
        "GAME_SERVER_SHARED_KEY",
        "ZNEON_GAME_SERVER_KEY",
        "postgres://",
        "redis://",
    ]
    for token in forbidden:
        if token in text:
            errors.append(f"workflow must not embed runtime secret/authority material: {token}")

    if text.count("bash tools/ue-linux.sh preflight") != 1:
        errors.append("build-evidence workflow must run runner preflight exactly once")

    if "package-client" in text or "package-server" in text:
        errors.append(
            "build-evidence workflow must not silently promote package evidence; packaging remains a separate retained-evidence gate"
        )

    if errors:
        print("Unreal evidence workflow validation FAILED", file=sys.stderr)
        for error in errors:
            print(f"- {error}", file=sys.stderr)
        return 1

    print("Unreal evidence workflow validation OK")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
