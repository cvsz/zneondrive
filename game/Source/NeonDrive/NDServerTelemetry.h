#pragma once

#include "CoreMinimal.h"

struct FNDServerTelemetrySnapshot
{
    int32 ActiveSessions = 0;
    uint64 PlayerJoinsTotal = 0;
    uint64 PlayerLeavesTotal = 0;
    uint64 InputClampsTotal = 0;
    uint64 TicketRedeemAttemptsTotal = 0;
    uint64 TicketRedeemSuccessesTotal = 0;
    uint64 TicketRedeemFailuresTotal = 0;
    double LastTickMilliseconds = 0.0;
    double PeakTickMilliseconds = 0.0;
};

// Game-thread-only aggregate telemetry for the authoritative Unreal gameplay server.
// No player, vehicle, race, ticket, credential, address, or request identifier is retained.
class NEONDRIVE_API FNDServerTelemetry final
{
public:
    static void Reset();
    static void RecordServerTick(float DeltaSeconds);
    static void RecordPlayerJoin();
    static void RecordPlayerLeave();
    static void RecordInputClamp();
    static void RecordTicketRedeemAttempt();
    static void RecordTicketRedeemSuccess();
    static void RecordTicketRedeemFailure();
    static FNDServerTelemetrySnapshot Snapshot();
    static void MaybeLog(double WorldSeconds);

private:
    static int32 ActiveSessions;
    static uint64 PlayerJoinsTotal;
    static uint64 PlayerLeavesTotal;
    static uint64 InputClampsTotal;
    static uint64 TicketRedeemAttemptsTotal;
    static uint64 TicketRedeemSuccessesTotal;
    static uint64 TicketRedeemFailuresTotal;
    static double LastTickMilliseconds;
    static double PeakTickMilliseconds;
    static double LastLogWorldSeconds;
};
