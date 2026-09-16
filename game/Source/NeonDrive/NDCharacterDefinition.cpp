#include "NDCharacterDefinition.h"

bool UNDCharacterDefinition::HasValidIdentity() const
{
    return !CharacterId.IsNone() && !FactionId.IsNone() && !ArchetypeId.IsNone();
}

bool UNDCharacterDefinition::UsesHighFidelityPresentation() const
{
    return PresentationTier == ENDCharacterPresentationTier::Hero ||
           PresentationTier == ENDCharacterPresentationTier::Player ||
           bUseHighFidelityFacialRig ||
           bUseGroom;
}

FPrimaryAssetId UNDCharacterDefinition::GetPrimaryAssetId() const
{
    return FPrimaryAssetId(TEXT("NeonCharacter"), GetFName());
}
