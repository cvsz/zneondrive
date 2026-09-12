#pragma once

#include "CoreMinimal.h"
#include "Interfaces/IHttpRequest.h"
#include "Interfaces/IHttpResponse.h"
#include "Subsystems/GameInstanceSubsystem.h"
#include "NDServiceSubsystem.generated.h"

class FJsonObject;

USTRUCT(BlueprintType)
struct FNDInventoryItem
{
    GENERATED_BODY()

    UPROPERTY(BlueprintReadOnly)
    FString ItemID;

    UPROPERTY(BlueprintReadOnly)
    int64 Quantity = 0;
};

USTRUCT(BlueprintType)
struct FNDPlayerSnapshot
{
    GENERATED_BODY()

    UPROPERTY(BlueprintReadOnly)
    FString AccountID;

    UPROPERTY(BlueprintReadOnly)
    FString CharacterID;

    UPROPERTY(BlueprintReadOnly)
    FString VehicleID;

    UPROPERTY(BlueprintReadOnly)
    int64 Money = 0;

    UPROPERTY(BlueprintReadOnly)
    int64 XP = 0;

    UPROPERTY(BlueprintReadOnly)
    int64 Reputation = 0;

    UPROPERTY(BlueprintReadOnly)
    int32 ActiveBuildRevision = 0;

    UPROPERTY(BlueprintReadOnly)
    TArray<FString> ActivePartIDs;

    UPROPERTY(BlueprintReadOnly)
    TArray<FNDInventoryItem> Inventory;

    UPROPERTY(BlueprintReadOnly)
    TArray<FString> Blueprints;

    UPROPERTY(BlueprintReadOnly)
    bool bStarterLineage = false;

    UPROPERTY(BlueprintReadOnly)
    bool bRoadworthy = false;

    UPROPERTY(BlueprintReadOnly)
    TArray<FString> CompletedQuests;
};

DECLARE_DYNAMIC_MULTICAST_DELEGATE_OneParam(FNDSessionReadySignature, bool, bSuccess);
DECLARE_DYNAMIC_MULTICAST_DELEGATE_OneParam(FNDSnapshotUpdatedSignature, FNDPlayerSnapshot, Snapshot);
DECLARE_DYNAMIC_MULTICAST_DELEGATE_OneParam(FNDServiceErrorSignature, FString, ErrorCode);

UCLASS()
class NEONDRIVE_API UNDServiceSubsystem : public UGameInstanceSubsystem
{
    GENERATED_BODY()

public:
    virtual void Initialize(FSubsystemCollectionBase& Collection) override;

    UFUNCTION(BlueprintCallable, Category = "NeonDrive|Service")
    void BootstrapOrResume();

    UFUNCTION(BlueprintCallable, Category = "NeonDrive|Service")
    void RefreshState();

    UFUNCTION(BlueprintCallable, Category = "NeonDrive|Service")
    void CompleteQuest(const FString& QuestID);

    UFUNCTION(BlueprintCallable, Category = "NeonDrive|Service")
    void ReviseBuild(const TArray<FString>& PartIDs);

    UFUNCTION(BlueprintCallable, Category = "NeonDrive|Service")
    void EnsureGameplayBinding();

    UFUNCTION(BlueprintPure, Category = "NeonDrive|Service")
    bool HasSession() const { return !SessionToken.IsEmpty(); }

    UFUNCTION(BlueprintPure, Category = "NeonDrive|Service")
    FNDPlayerSnapshot GetSnapshot() const { return Snapshot; }

    UPROPERTY(BlueprintAssignable, Category = "NeonDrive|Service")
    FNDSessionReadySignature OnSessionReady;

    UPROPERTY(BlueprintAssignable, Category = "NeonDrive|Service")
    FNDSnapshotUpdatedSignature OnSnapshotUpdated;

    UPROPERTY(BlueprintAssignable, Category = "NeonDrive|Service")
    FNDServiceErrorSignature OnServiceError;

private:
    FString ServiceBaseUrl = TEXT("http://127.0.0.1:18080");
    FString SessionToken;
    FString ResumeKey;
    FNDPlayerSnapshot Snapshot;
    bool bBootstrapInFlight = false;
    bool bTicketInFlight = false;

    FString ResumeKeyPath() const;
    void LoadResumeKey();
    void SaveResumeKey() const;
    void RequestGameTicket();

    void HandleBootstrap(FHttpRequestPtr Request, FHttpResponsePtr Response, bool bSucceeded);
    void HandleState(FHttpRequestPtr Request, FHttpResponsePtr Response, bool bSucceeded);
    void HandleQuest(FHttpRequestPtr Request, FHttpResponsePtr Response, bool bSucceeded);
    void HandleBuild(FHttpRequestPtr Request, FHttpResponsePtr Response, bool bSucceeded);
    void HandleGameTicket(FHttpRequestPtr Request, FHttpResponsePtr Response, bool bSucceeded);

    bool ParseSnapshotObject(const TSharedPtr<FJsonObject>& Object, FNDPlayerSnapshot& OutSnapshot) const;
    void BroadcastError(const FString& Code);
};
