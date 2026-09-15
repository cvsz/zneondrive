#include "NDCharacterComponent.h"

#include "NDCharacterDefinition.h"
#include "Animation/AnimInstance.h"
#include "Components/SkeletalMeshComponent.h"

UNDCharacterComponent::UNDCharacterComponent()
{
    PrimaryComponentTick.bCanEverTick = false;
    SetIsReplicatedByDefault(false);
}

void UNDCharacterComponent::BeginPlay()
{
    Super::BeginPlay();

    if (CharacterDefinition)
    {
        ApplyDefinition(CharacterDefinition);
    }
}

bool UNDCharacterComponent::ResolveCharacterMesh()
{
    if (!GetOwner())
    {
        return false;
    }

    CharacterMesh = GetOwner()->FindComponentByClass<USkeletalMeshComponent>();
    return CharacterMesh != nullptr;
}

bool UNDCharacterComponent::ApplyDefinition(UNDCharacterDefinition* Definition)
{
    if (!Definition || !ResolveCharacterMesh())
    {
        return false;
    }

    CharacterDefinition = Definition;

    if (USkeletalMesh* Mesh = Definition->BodyMesh.LoadSynchronous())
    {
        CharacterMesh->SetSkeletalMesh(Mesh);
    }

    if (UClass* AnimationClass = Definition->AnimationClass.LoadSynchronous())
    {
        CharacterMesh->SetAnimationMode(EAnimationMode::AnimationBlueprint);
        CharacterMesh->SetAnimInstanceClass(AnimationClass);
    }

    CharacterMesh->SetRelativeScale3D(FVector(Definition->CharacterScale));
    return CharacterMesh->GetSkeletalMeshAsset() != nullptr;
}
