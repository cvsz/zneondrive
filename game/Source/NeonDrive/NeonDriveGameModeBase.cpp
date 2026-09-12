#include "NeonDriveGameModeBase.h"
#include "NDVehiclePawn.h"

ANeonDriveGameModeBase::ANeonDriveGameModeBase()
{
    DefaultPawnClass = ANDVehiclePawn::StaticClass();
}
