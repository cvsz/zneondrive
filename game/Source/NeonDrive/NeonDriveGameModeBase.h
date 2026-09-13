#pragma once

#include "CoreMinimal.h"
#include "GameFramework/GameModeBase.h"
#include "NeonDriveGameModeBase.generated.h"

UCLASS()
class NEONDRIVE_API ANeonDriveGameModeBase : public AGameModeBase
{
    GENERATED_BODY()

public:
    ANeonDriveGameModeBase();

    virtual void Tick(float DeltaSeconds) override;
    virtual void PostLogin(APlayerController* NewPlayer) override;
    virtual void Logout(AController* Exiting) override;
};
