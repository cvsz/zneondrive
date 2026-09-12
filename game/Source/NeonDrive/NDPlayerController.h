#pragma once

#include "CoreMinimal.h"
#include "GameFramework/PlayerController.h"
#include "Interfaces/IHttpRequest.h"
#include "Interfaces/IHttpResponse.h"
#include "NDPlayerController.generated.h"

UCLASS()
class NEONDRIVE_API ANDPlayerController : public APlayerController
{
    GENERATED_BODY()

public:
    virtual void BeginPlay() override;

    UFUNCTION(BlueprintCallable, Category = "NeonDrive|Service")
    void SubmitGameTicket(const FString& Ticket);

protected:
    UFUNCTION(Server, Reliable)
    void ServerSubmitGameTicket(const FString& Ticket);

private:
    void HandleTicketRedeemed(FHttpRequestPtr Request, FHttpResponsePtr Response, bool bSucceeded);
};
