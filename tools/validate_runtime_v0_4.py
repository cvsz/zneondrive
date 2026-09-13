#!/usr/bin/env python3
"""Static checks for the Phase 4 Unreal/Go runtime baseline."""

from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]

REQUIRED = [
    "game/NeonDrive.uproject",
    "game/Source/NeonDrive/NeonDrive.Build.cs",
    "game/Source/NeonDrive/NDVehiclePawn.h",
    "game/Source/NeonDrive/NDVehiclePawn.cpp",
    "game/Source/NeonDriveServer.Target.cs",
    "services/game-api/go.mod",
    "services/game-api/cmd/server/main.go",
    "services/game-api/internal/core/core.go",
    "services/game-api/internal/store/postgres.go",
    "services/game-api/internal/store/migrations/001_init.sql",
    "services/game-api/internal/httpapi/api.go",
    "compose.yaml",
]

for relative in REQUIRED:
    path = ROOT / relative
    assert path.is_file(), f"missing runtime file: {relative}"

uproject = (ROOT / "game/NeonDrive.uproject").read_text(encoding="utf-8")
assert '"EngineAssociation": "5.8"' in uproject
assert '"Name": "NeonDrive"' in uproject

server_target = (ROOT / "game/Source/NeonDriveServer.Target.cs").read_text(encoding="utf-8")
assert "TargetType.Server" in server_target

pawn_h = (ROOT / "game/Source/NeonDrive/NDVehiclePawn.h").read_text(encoding="utf-8")
pawn_cpp = (ROOT / "game/Source/NeonDrive/NDVehiclePawn.cpp").read_text(encoding="utf-8")
assert "UFUNCTION(Server" in pawn_h
assert "SetReplicateMovement(true)" in pawn_cpp
assert "if (!HasAuthority())" in pawn_cpp
assert "FMath::Clamp" in pawn_cpp

store_go = (ROOT / "services/game-api/internal/store/postgres.go").read_text(encoding="utf-8")
for required in (
    "resume_key_hash",
    "CreateSession",
    "CompleteQuest",
    "ReviseBuild",
    "FOR UPDATE",
):
    assert required in store_go, f"missing durable-store invariant: {required}"

assert (
    'questID == "MQ012"' in store_go or 'case "MQ012":' in store_go
), "missing durable-store invariant: MQ012 Roadworthy transition"

migration = (ROOT / "services/game-api/internal/store/migrations/001_init.sql").read_text(encoding="utf-8")
for table in ("accounts", "characters", "vehicles", "vehicle_builds", "sessions", "quest_completions"):
    assert f"CREATE TABLE IF NOT EXISTS {table}" in migration

compose = (ROOT / "compose.yaml").read_text(encoding="utf-8")
assert "55432:5432" not in compose, "ports should remain environment-configurable"
assert "POSTGRES_PORT:-55432" in compose
assert "REDIS_PORT:-56379" in compose
assert "GAME_API_PORT:-18080" in compose

print("runtime v0.4 static validation OK")
