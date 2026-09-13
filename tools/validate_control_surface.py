#!/usr/bin/env python3
from pathlib import Path
import re
import sys

ROOT = Path(__file__).resolve().parents[1]

REQUIRED = [
    "Makefile",
    "tools/zneondrive-control.sh",
    "tools/install-client.ps1",
    "game/Source/NeonDriveClient.Target.cs",
    "docs/control-panel.md",
]

def require(text: str, token: str, label: str, errors: list[str]) -> None:
    if token not in text:
        errors.append(f"{label}: missing {token!r}")

def main() -> int:
    errors: list[str] = []
    for rel in REQUIRED:
        if not (ROOT / rel).is_file():
            errors.append(f"missing control-surface asset: {rel}")
    if errors:
        for error in errors:
            print(error, file=sys.stderr)
        return 1

    make = (ROOT / "Makefile").read_text(encoding="utf-8")
    shell = (ROOT / "tools/zneondrive-control.sh").read_text(encoding="utf-8")
    ps1 = (ROOT / "tools/install-client.ps1").read_text(encoding="utf-8")
    target = (ROOT / "game/Source/NeonDriveClient.Target.cs").read_text(encoding="utf-8")
    cpp = (ROOT / "game/Source/NeonDrive/NDServiceSubsystem.cpp").read_text(encoding="utf-8")

    for token in [
        "control-panel:", "ue-detect:", "server-install:", "server-up:", "server-down:",
        "client-install:", "client-build:", "client-package-linux:", "client-play:",
        "game-server-package-linux:", "package-all-linux:",
        "game-server-install:", "game-server-start:", "game-server-stop:",
        "full-install:", "full-up:", "full-down:", "status:",
    ]:
        require(make, token, "Makefile", errors)

    for token in [
        "ue_detect()", "server_install()", "client_install()", "client_package_linux()", "client_play()",
        "game_server_package_linux()", "package_all_linux()",
        "game_server_install()", "game_server_start()", "control_panel()",
        "CLIENT_PACKAGE is required for full-install", "SERVER_PACKAGE is required for full-install",
        "CONFIRM_RESET", "ZNEON_GAME_API_URL", "ZNEON_GAME_SERVER_KEY",
    ]:
        require(shell, token, "Linux control panel", errors)

    require(ps1, "Server shared keys are intentionally not accepted", "Windows installer", errors)
    if "GAME_SERVER_SHARED_KEY" in ps1 or "ZNEON_GAME_SERVER_KEY" in ps1:
        errors.append("Windows player installer must not accept or embed game-server shared keys.")

    require(target, "Type = TargetType.Client;", "Unreal client target", errors)
    require(cpp, "ZNEON_GAME_API_URL", "Unreal runtime endpoint override", errors)
    require(cpp, "ZNeonApi=", "Unreal command-line endpoint override", errors)

    client_section = re.search(r"client_play\(\).*?\n\}", shell, flags=re.S)
    if client_section and "ZNEON_GAME_SERVER_KEY" in client_section.group(0):
        errors.append("client_play must never receive the dedicated-server shared key.")

    if errors:
        print("control-surface validation FAILED", file=sys.stderr)
        for error in errors:
            print(f"- {error}", file=sys.stderr)
        return 1

    print("control-surface validation OK")
    return 0

if __name__ == "__main__":
    raise SystemExit(main())
