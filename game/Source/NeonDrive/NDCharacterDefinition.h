#pragma once

#include "CoreMinimal.h"
#include "Engine/DataAsset.h"
#include "NDCharacterDefinition.generated.h"

class UAnimInstance;
class USkeletalMesh;

UENUM(BlueprintType)
enum class ENDCharacterPresentationTier : uint8
{
    Crowd,
    NPC,
    Player,
    Hero
};

UCLASS(BlueprintType)
class NEONDRIVE_API UNDCharacterDefinition : public UPrimaryDataAsset
{
    GENERATED_BODY()

public:
    UPROPERTY(EditDefaultsOnly, BlueprintReadOnly, Category = "Identity")
    FName CharacterId;

    UPROPERTY(EditDefaultsOnly, BlueprintReadOnly, Category = "Identity")
    FText DisplayName;

    UPROPERTY(EditDefaultsOnly, BlueprintReadOnly, Category = "Identity")
    FName FactionId;

    UPROPERTY(EditDefaultsOnly, BlueprintReadOnly, Category = "Identity")
    FName ArchetypeId;

    UPROPERTY(EditDefaultsOnly, BlueprintReadOnly, Category = "Visual")
    TSoftObjectPtr<USkeletalMesh> BodyMesh;

    UPROPERTY(EditDefaultsOnly, BlueprintReadOnly, Category = "Animation")
    TSoftClassPtr<UAnimInstance> AnimationClass;

    UPROPERTY(EditDefaultsOnly, BlueprintReadOnly, Category = "Animation")
    FName LocomotionProfileId = TEXT("default");

    UPROPERTY(EditDefaultsOnly, BlueprintReadOnly, Category = "Animation")
    FName MotionMatchingProfileId = TEXT("humanoid_default");

    UPROPERTY(EditDefaultsOnly, BlueprintReadOnly, Category = "Rig")
    FName SkeletonProfileId = TEXT("neondrive_humanoid_v1");

    UPROPERTY(EditDefaultsOnly, BlueprintReadOnly, Category = "Facial")
    FName FacialProfileId = TEXT("none");

    UPROPERTY(EditDefaultsOnly, BlueprintReadOnly, Category = "Visual")
    ENDCharacterPresentationTier PresentationTier = ENDCharacterPresentationTier::NPC;

    UPROPERTY(EditDefaultsOnly, BlueprintReadOnly, Category = "Visual", meta = (ClampMin = "0.5", ClampMax = "2.0"))
    float CharacterScale = 1.0f;

    UPROPERTY(EditDefaultsOnly, BlueprintReadOnly, Category = "Visual")
    bool bUseHighFidelityFacialRig = false;

    UPROPERTY(EditDefaultsOnly, BlueprintReadOnly, Category = "Visual")
    bool bUseGroom = false;

    UPROPERTY(EditDefaultsOnly, BlueprintReadOnly, Category = "Network")
    bool bReplicateCosmeticState = true;

    UFUNCTION(BlueprintPure, Category = "Character")
    bool HasValidIdentity() const;

    UFUNCTION(BlueprintPure, Category = "Character")
    bool UsesHighFidelityPresentation() const;

    virtual FPrimaryAssetId GetPrimaryAssetId() const override;
};
