#include "NDServerTelemetry.h"

DEFINE_LOG_CATEGORY_STATIC(LogNDServerTelemetry, Log, All);

int32 FNDServerTelemetry::ActiveSessions = 0;
uint64 FNDServerTelemetry::PlayerJoinsTotal = 0;
uint64 FNDServerTelemetry::PlayerLeavesTotal = 0;
uint64 FNDServerTelemetry::InputClampsTotal = 0;
uint64 FNDServerTelemetry::AuthorityMovementTicksTotal = 0;
uint64 FNDServerTelemetry::IdentityGateBlocksTotal = 0;
uint64 FNDServerTelemetry::CollisionBlocksTotal = 0;
uint64 FNDServerTelemetry::NetUpdateRequestsTotal = 0;
uint64 FNDServerTelemetry::ImpossibleDisplacementsTotal = 0;
uint64 FNDServerTelemetry::TicketRedeemAttemptsTotal = 0;
uint64 FNDServerTelemetry::TicketRedeemSuccessesTotal = 0;
uint64 FNDServerTelemetry::TicketRedeemFailuresTotal = 0;
double FNDServerTelemetry::LastTickMilliseconds = 0.0;
double FNDServerTelemetry::PeakTickMilliseconds = 0.0;
double FNDServerTelemetry::LastLogWorldSeconds = -10.0;

void FNDServerTelemetry::Reset()
{
    ActiveSessions = 0;
    PlayerJoinsTotal = 0;
    PlayerLeavesTotal = 0;
    InputClampsTotal = 0;
    AuthorityMovementTicksTotal = 0;
    IdentityGateBlocksTotal = 0;
    CollisionBlocksTotal = 0;
    NetUpdateRequestsTotal = 0;
    ImpossibleDisplacementsTotal = 0;
    TicketRedeemAttemptsTotal = 0;
    TicketRedeemSuccessesTotal = 0;
    TicketRedeemFailuresTotal = 0;
    LastTickMilliseconds = 0.0;
    PeakTickMilliseconds = 0.0;
    LastLogWorldSeconds = -10.0;
}

void FNDServerTelemetry::RecordServerTick(float DeltaSeconds)
{
    LastTickMilliseconds = FMath::Max(0.0, static_cast<double>(DeltaSeconds) * 1000.0);
    PeakTickMilliseconds = FMath::Max(PeakTickMilliseconds, LastTickMilliseconds);
}

void FNDServerTelemetry::RecordPlayerJoin()
{
    ++ActiveSessions;
    ++PlayerJoinsTotal;
}

void FNDServerTelemetry::RecordPlayerLeave()
{
    ActiveSessions = FMath::Max(0, ActiveSessions - 1);
    ++PlayerLeavesTotal;
}

void FNDServerTelemetry::RecordInputClamp()
{
    ++InputClampsTotal;
}

void FNDServerTelemetry::RecordAuthorityMovementTick()
{
    ++AuthorityMovementTicksTotal;
}

void FNDServerTelemetry::RecordIdentityGateBlock()
{
    ++IdentityGateBlocksTotal;
}

void FNDServerTelemetry::RecordCollisionBlock()
{
    ++CollisionBlocksTotal;
}

void FNDServerTelemetry::RecordNetUpdateRequest()
{
    ++NetUpdateRequestsTotal;
}

void FNDServerTelemetry::RecordImpossibleDisplacement()
{
    ++ImpossibleDisplacementsTotal;
}

void FNDServerTelemetry::RecordTicketRedeemAttempt()
{
    ++TicketRedeemAttemptsTotal;
}

void FNDServerTelemetry::RecordTicketRedeemSuccess()
{
    ++TicketRedeemSuccessesTotal;
}

void FNDServerTelemetry::RecordTicketRedeemFailure()
{
    ++TicketRedeemFailuresTotal;
}

FNDServerTelemetrySnapshot FNDServerTelemetry::Snapshot()
{
    FNDServerTelemetrySnapshot Result;
    Result.ActiveSessions = ActiveSessions;
    Result.PlayerJoinsTotal = PlayerJoinsTotal;
    Result.PlayerLeavesTotal = PlayerLeavesTotal;
    Result.InputClampsTotal = InputClampsTotal;
    Result.AuthorityMovementTicksTotal = AuthorityMovementTicksTotal;
    Result.IdentityGateBlocksTotal = IdentityGateBlocksTotal;
    Result.CollisionBlocksTotal = CollisionBlocksTotal;
    Result.NetUpdateRequestsTotal = NetUpdateRequestsTotal;
    Result.ImpossibleDisplacementsTotal = ImpossibleDisplacementsTotal;
    Result.TicketRedeemAttemptsTotal = TicketRedeemAttemptsTotal;
    Result.TicketRedeemSuccessesTotal = TicketRedeemSuccessesTotal;
    Result.TicketRedeemFailuresTotal = TicketRedeemFailuresTotal;
    Result.LastTickMilliseconds = LastTickMilliseconds;
    Result.PeakTickMilliseconds = PeakTickMilliseconds;
    return Result;
}

void FNDServerTelemetry::MaybeLog(double WorldSeconds)
{
    constexpr double LogIntervalSeconds = 10.0;
    if (WorldSeconds < 0.0 || WorldSeconds - LastLogWorldSeconds < LogIntervalSeconds)
    {
        return;
    }

    LastLogWorldSeconds = WorldSeconds;
    const FNDServerTelemetrySnapshot Current = Snapshot();
    UE_LOG(
        LogNDServerTelemetry,
        Log,
        TEXT("metric=zneondrive_unreal_server active_sessions=%d tick_ms=%.3f peak_tick_ms=%.3f joins_total=%llu leaves_total=%llu input_clamps_total=%llu authority_movement_ticks_total=%llu identity_gate_blocks_total=%llu collision_blocks_total=%llu net_update_requests_total=%llu impossible_displacements_total=%llu ticket_redeem_attempts_total=%llu ticket_redeem_successes_total=%llu ticket_redeem_failures_total=%llu"),
        Current.ActiveSessions,
        Current.LastTickMilliseconds,
        Current.PeakTickMilliseconds,
        static_cast<unsigned long long>(Current.PlayerJoinsTotal),
        static_cast<unsigned long long>(Current.PlayerLeavesTotal),
        static_cast<unsigned long long>(Current.InputClampsTotal),
        static_cast<unsigned long long>(Current.AuthorityMovementTicksTotal),
        static_cast<unsigned long long>(Current.IdentityGateBlocksTotal),
        static_cast<unsigned long long>(Current.CollisionBlocksTotal),
        static_cast<unsigned long long>(Current.NetUpdateRequestsTotal),
        static_cast<unsigned long long>(Current.ImpossibleDisplacementsTotal),
        static_cast<unsigned long long>(Current.TicketRedeemAttemptsTotal),
        static_cast<unsigned long long>(Current.TicketRedeemSuccessesTotal),
        static_cast<unsigned long long>(Current.TicketRedeemFailuresTotal)
    );
}
