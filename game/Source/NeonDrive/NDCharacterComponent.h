#pragma once

#include "CoreMinimal.h"
#include "Components/ActorComponent.h"
#include "NDCharacterComponent.generated.h"

class UNDCharacterDefinition;
class USkeletalMeshComponent;

UCLASS(ClassGroup = (NeonDrive), meta = (BlueprintSpawnableComponent))
class NEONDRIVE_API UNDCharacterComponent : public UActorComponent
{
    GENERATED_BODY()

public:
    UNDCharacterComponent();

    UFUNCTION(BlueprintCallable, Category = "Character")
    bool ApplyDefinition(UNDCharacterDefinition* Definition);

    UFUNCTION(BlueprintPure, Category = "Character")
    UNDCharacterDefinition* GetDefinition() const { return CharacterDefinition; }

protected:
    UPROPERTY(VisibleInstanceOnly, BlueprintReadOnly, Category = "Character")
    TObjectPtr<UNDCharacterDefinition> CharacterDefinition;

    UPROPERTY(VisibleInstanceOnly, BlueprintReadOnly, Category = "Character")
    TObjectPtr<USkeletalMeshComponent> CharacterMesh;

    virtual void BeginPlay() override;

private:
    bool ResolveCharacterMesh();
};
