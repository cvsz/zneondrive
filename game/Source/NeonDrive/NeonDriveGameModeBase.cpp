#include "NeonDriveGameModeBase.h"
#include "NDPlayerController.h"
#include "NDVehiclePawn.h"

ANeonDriveGameModeBase::ANeonDriveGameModeBase()
{
    DefaultPawnClass = ANDVehiclePawn::StaticClass();
    PlayerControllerClass = ANDPlayerController::StaticClass();
}
