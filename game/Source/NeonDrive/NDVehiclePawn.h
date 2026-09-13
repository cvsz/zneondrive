#pragma once

#include "CoreMinimal.h"
#include "GameFramework/Pawn.h"
#include "NDVehiclePawn.generated.h"

class UBoxComponent;
class UCameraComponent;
class USpringArmComponent;
class UStaticMeshComponent;

UCLASS()
class NEONDRIVE_API ANDVehiclePawn : public APawn
{
    GENERATED_BODY()

public:
    ANDVehiclePawn();
    virtual void Tick(float DeltaSeconds) override;
    virtual void SetupPlayerInputComponent(UInputComponent* PlayerInputComponent) override;
    virtual void GetLifetimeReplicatedProps(TArray<FLifetimeProperty>& OutLifetimeProps) const override;

    void ApplyDurableIdentity(
        const FString& VehicleID,
        int32 BuildRevision,
        const TArray<FString>& PartIDs,
        bool bRoadworthy
    );

    UFUNCTION(BlueprintPure, Category = "NeonDrive|Vehicle")
    bool HasDurableIdentity() const { return bDurableIdentityBound; }

    UPROPERTY(Replicated, BlueprintReadOnly, Category = "NeonDrive|Vehicle")
    FString DurableVehicleID;

    UPROPERTY(Replicated, BlueprintReadOnly, Category = "NeonDrive|Vehicle")
    int32 DurableBuildRevision = 0;

    UPROPERTY(Replicated, BlueprintReadOnly, Category = "NeonDrive|Vehicle")
    TArray<FString> DurablePartIDs;

    UPROPERTY(Replicated, BlueprintReadOnly, Category = "NeonDrive|Vehicle")
    bool bDurableRoadworthy = false;

    UPROPERTY(Replicated, BlueprintReadOnly, Category = "NeonDrive|Vehicle")
    bool bDurableIdentityBound = false;

protected:
    UPROPERTY(VisibleAnywhere, Category = "Vehicle")
    TObjectPtr<UBoxComponent> Collision;

    UPROPERTY(VisibleAnywhere, Category = "Vehicle")
    TObjectPtr<UStaticMeshComponent> Body;

    UPROPERTY(VisibleAnywhere, Category = "Camera")
    TObjectPtr<USpringArmComponent> SpringArm;

    UPROPERTY(VisibleAnywhere, Category = "Camera")
    TObjectPtr<UCameraComponent> Camera;

    UPROPERTY(Replicated)
    float AuthoritativeThrottle = 0.0f;

    UPROPERTY(Replicated)
    float AuthoritativeSteering = 0.0f;

    float LocalThrottle = 0.0f;
    float LocalSteering = 0.0f;

    UPROPERTY(EditDefaultsOnly, Category = "Prototype Driving")
    float MaxSpeedCmPerSecond = 1400.0f;

    UPROPERTY(EditDefaultsOnly, Category = "Prototype Driving")
    float TurnRateDegreesPerSecond = 75.0f;

    UPROPERTY(EditDefaultsOnly, Category = "Prototype Driving|Integrity")
    float AuthorityDisplacementSlackCm = 150.0f;

    FVector LastAuthorityLocation = FVector::ZeroVector;
    bool bAuthorityLocationBaselineValid = false;

    void InputThrottle(float Value);
    void InputSteering(float Value);
    void PushDrivingInput();

    UFUNCTION(Server, Unreliable)
    void ServerSetDrivingInput(float Throttle, float Steering);
};
