#!/usr/bin/env python3
"""Static contract checks for Phase 4.1 runtime integration v0.5."""

from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]

required = [
    "docs/runtime-integration-v0.5.md",
    "game/Source/NeonDrive/NDServiceSubsystem.h",
    "game/Source/NeonDrive/NDServiceSubsystem.cpp",
    "game/Source/NeonDrive/NDPlayerController.h",
    "game/Source/NeonDrive/NDPlayerController.cpp",
    "services/game-api/internal/store/migrations/002_game_tickets.sql",
    "services/game-api/internal/store/tickets.go",
    "services/game-api/internal/httpapi/api_integration_test.go",
]

for relative in required:
    assert (ROOT / relative).is_file(), f"missing v0.5 integration file: {relative}"

api = (ROOT / "services/game-api/internal/httpapi/api.go").read_text(encoding="utf-8")
for token in (
    'POST /v1/game-tickets',
    'POST /v1/internal/game-tickets/redeem',
    'X-Game-Server-Key',
    "constantTimeSecretEqual",
    "gameplayTicketTTL",
):
    assert token in api, f"missing gameplay-ticket API invariant: {token}"

tickets = (ROOT / "services/game-api/internal/store/tickets.go").read_text(encoding="utf-8")
for token in (
    "IssueGameTicket",
    "RedeemGameTicket",
    "consumed_at=now()",
    "expires_at > now()",
):
    assert token in tickets, f"missing one-time ticket store invariant: {token}"

postgres = (ROOT / "services/game-api/internal/store/postgres.go").read_text(encoding="utf-8")
assert postgres.index('"001_init"') < postgres.index('"002_game_tickets"'), "migrations must be ordered"

controller = (ROOT / "game/Source/NeonDrive/NDPlayerController.cpp").read_text(encoding="utf-8")
for token in (
    "ServerSubmitGameTicket",
    "ZNEON_GAME_SERVER_KEY",
    "/v1/internal/game-tickets/redeem",
    "ApplyDurableIdentity",
):
    assert token in controller, f"missing dedicated-server binding invariant: {token}"

subsystem = (ROOT / "game/Source/NeonDrive/NDServiceSubsystem.cpp").read_text(encoding="utf-8")
for token in (
    "/v1/sessions/bootstrap",
    "/v1/state",
    "/v1/game-tickets",
    "resume_key.txt",
    "EnsureGameplayBinding",
):
    assert token in subsystem, f"missing client session invariant: {token}"
assert "ZNEON_GAME_SERVER_KEY" not in subsystem, "client subsystem must never read the game-server shared key"
assert "X-Game-Server-Key" not in subsystem, "client subsystem must never send the internal server key"

pawn = (ROOT / "game/Source/NeonDrive/NDVehiclePawn.cpp").read_text(encoding="utf-8")
for token in (
    "bDurableIdentityBound",
    "DurableVehicleID",
    "DurableBuildRevision",
    "DurablePartIDs",
    "ForceNetUpdate",
):
    assert token in pawn, f"missing durable pawn invariant: {token}"

migration = (ROOT / "services/game-api/internal/store/migrations/002_game_tickets.sql").read_text(encoding="utf-8")
assert "game_tickets" in migration
assert "ticket_hash TEXT PRIMARY KEY" in migration
assert "consumed_at TIMESTAMPTZ" in migration

print("runtime integration v0.5 static validation OK")
