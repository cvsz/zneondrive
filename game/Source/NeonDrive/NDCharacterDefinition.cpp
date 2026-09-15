#include "NDCharacterDefinition.h"

FPrimaryAssetId UNDCharacterDefinition::GetPrimaryAssetId() const
{
    return FPrimaryAssetId(TEXT("NeonCharacter"), GetFName());
}
