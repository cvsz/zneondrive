#include "NeonDriveGameModeBase.h"
#include "NDPlayerController.h"
#include "NDServerTelemetry.h"
#include "NDVehiclePawn.h"

ANeonDriveGameModeBase::ANeonDriveGameModeBase()
{
    PrimaryActorTick.bCanEverTick = true;
    DefaultPawnClass = ANDVehiclePawn::StaticClass();
    PlayerControllerClass = ANDPlayerController::StaticClass();
}

void ANeonDriveGameModeBase::Tick(float DeltaSeconds)
{
    Super::Tick(DeltaSeconds);

    FNDServerTelemetry::RecordServerTick(DeltaSeconds);
    if (const UWorld* World = GetWorld())
    {
        FNDServerTelemetry::MaybeLog(World->GetTimeSeconds());
    }
}

void ANeonDriveGameModeBase::PostLogin(APlayerController* NewPlayer)
{
    Super::PostLogin(NewPlayer);
    FNDServerTelemetry::RecordPlayerJoin();
}

void ANeonDriveGameModeBase::Logout(AController* Exiting)
{
    FNDServerTelemetry::RecordPlayerLeave();
    Super::Logout(Exiting);
}
