#!/usr/bin/env python3
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]

telemetry_h = (ROOT / "game/Source/NeonDrive/NDServerTelemetry.h").read_text(encoding="utf-8")
telemetry_cpp = (ROOT / "game/Source/NeonDrive/NDServerTelemetry.cpp").read_text(encoding="utf-8")
game_mode = (ROOT / "game/Source/NeonDrive/NeonDriveGameModeBase.cpp").read_text(encoding="utf-8")
controller = (ROOT / "game/Source/NeonDrive/NDPlayerController.cpp").read_text(encoding="utf-8")
pawn = (ROOT / "game/Source/NeonDrive/NDVehiclePawn.cpp").read_text(encoding="utf-8")

required_telemetry = {
    "bounded metric prefix": "metric=zneondrive_unreal_server",
    "active sessions": "active_sessions=%d",
    "tick latency": "tick_ms=%.3f",
    "peak tick latency": "peak_tick_ms=%.3f",
    "join counter": "joins_total=%llu",
    "leave counter": "leaves_total=%llu",
    "input clamp counter": "input_clamps_total=%llu",
    "authority movement ticks": "authority_movement_ticks_total=%llu",
    "identity gate blocks": "identity_gate_blocks_total=%llu",
    "collision blocks": "collision_blocks_total=%llu",
    "net update requests": "net_update_requests_total=%llu",
    "ticket attempts": "ticket_redeem_attempts_total=%llu",
    "ticket successes": "ticket_redeem_successes_total=%llu",
    "ticket failures": "ticket_redeem_failures_total=%llu",
    "fixed log cadence": "LogIntervalSeconds = 10.0",
}
for label, token in required_telemetry.items():
    if token not in telemetry_cpp:
        raise SystemExit(f"missing Unreal server telemetry contract: {label}: {token}")

for label, token in {
    "server tick hook": "FNDServerTelemetry::RecordServerTick(DeltaSeconds)",
    "join hook": "FNDServerTelemetry::RecordPlayerJoin()",
    "leave hook": "FNDServerTelemetry::RecordPlayerLeave()",
}.items():
    if token not in game_mode:
        raise SystemExit(f"missing GameMode telemetry hook: {label}")

for label, token in {
    "ticket attempt hook": "FNDServerTelemetry::RecordTicketRedeemAttempt()",
    "ticket success hook": "FNDServerTelemetry::RecordTicketRedeemSuccess()",
    "ticket failure hook": "FNDServerTelemetry::RecordTicketRedeemFailure()",
}.items():
    if token not in controller:
        raise SystemExit(f"missing gameplay-ticket telemetry hook: {label}")

for label, token in {
    "authoritative input clamp": "FNDServerTelemetry::RecordInputClamp()",
    "authority movement tick": "FNDServerTelemetry::RecordAuthorityMovementTick()",
    "durable identity gate": "FNDServerTelemetry::RecordIdentityGateBlock()",
    "collision block": "FNDServerTelemetry::RecordCollisionBlock()",
    "explicit replication update request": "FNDServerTelemetry::RecordNetUpdateRequest()",
}.items():
    if token not in pawn:
        raise SystemExit(f"missing authoritative vehicle telemetry hook: {label}")

if "Hit.bBlockingHit" not in pawn:
    raise SystemExit("collision telemetry must derive from authoritative swept movement result")
if "ForceNetUpdate();" not in pawn:
    raise SystemExit("net update telemetry must remain paired with an actual ForceNetUpdate request")

forbidden_metric_tokens = (
    "%s",
    "vehicle_id",
    "race_id",
    "session_token",
    "gameplay_ticket",
    "game_server_key",
    "peer_address",
    "client_ip",
)
metric_line = next(
    (line.lower() for line in telemetry_cpp.splitlines() if "metric=zneondrive_unreal_server" in line),
    "",
)
for forbidden in forbidden_metric_tokens:
    if forbidden in metric_line:
        raise SystemExit(f"Unreal server metric line exposes forbidden dynamic/string field: {forbidden}")

if "FString" in telemetry_h:
    raise SystemExit("bounded telemetry snapshot must not retain FString identifiers")
if "GetEnvironmentVariable" in telemetry_cpp or "GetEnvironmentVariable" in telemetry_h:
    raise SystemExit("telemetry accumulator must not read runtime credentials or endpoints")

print("Unreal dedicated-server observability source contract OK")
