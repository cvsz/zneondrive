#include "NDServiceSubsystem.h"

#include "Dom/JsonObject.h"
#include "HAL/FileManager.h"
#include "HttpModule.h"
#include "Kismet/GameplayStatics.h"
#include "Misc/ConfigCacheIni.h"
#include "Misc/FileHelper.h"
#include "Misc/Guid.h"
#include "Misc/Paths.h"
#include "NDPlayerController.h"
#include "Serialization/JsonReader.h"
#include "Serialization/JsonSerializer.h"

namespace
{
TSharedRef<IHttpRequest, ESPMode::ThreadSafe> NewJsonRequest(
    const FString& Url,
    const FString& Verb,
    const FString& SessionToken = FString())
{
    TSharedRef<IHttpRequest, ESPMode::ThreadSafe> Request = FHttpModule::Get().CreateRequest();
    Request->SetURL(Url);
    Request->SetVerb(Verb);
    Request->SetHeader(TEXT("Accept"), TEXT("application/json"));
    Request->SetHeader(TEXT("Content-Type"), TEXT("application/json"));
    if (!SessionToken.IsEmpty())
    {
        Request->SetHeader(TEXT("Authorization"), TEXT("Bearer ") + SessionToken);
    }
    return Request;
}

bool ParseJsonObject(const FHttpResponsePtr& Response, TSharedPtr<FJsonObject>& OutObject)
{
    if (!Response.IsValid())
    {
        return false;
    }

    TSharedRef<TJsonReader<>> Reader = TJsonReaderFactory<>::Create(Response->GetContentAsString());
    return FJsonSerializer::Deserialize(Reader, OutObject) && OutObject.IsValid();
}
}

void UNDServiceSubsystem::Initialize(FSubsystemCollectionBase& Collection)
{
    Super::Initialize(Collection);

    FString ConfiguredUrl;
    if (GConfig && GConfig->GetString(TEXT("NeonDrive.Service"), TEXT("BaseUrl"), ConfiguredUrl, GGameIni))
    {
        ConfiguredUrl.TrimStartAndEndInline();
        if (!ConfiguredUrl.IsEmpty())
        {
            ServiceBaseUrl = ConfiguredUrl;
        }
    }
    ServiceBaseUrl.RemoveFromEnd(TEXT("/"));

    LoadResumeKey();
    BootstrapOrResume();
}

FString UNDServiceSubsystem::ResumeKeyPath() const
{
    return FPaths::Combine(FPaths::ProjectSavedDir(), TEXT("NeonDrive"), TEXT("resume_key.txt"));
}

void UNDServiceSubsystem::LoadResumeKey()
{
    FString Loaded;
    if (FFileHelper::LoadFileToString(Loaded, *ResumeKeyPath()))
    {
        Loaded.TrimStartAndEndInline();
        if (Loaded.Len() >= 32)
        {
            ResumeKey = Loaded;
        }
    }
}

void UNDServiceSubsystem::SaveResumeKey() const
{
    if (ResumeKey.IsEmpty())
    {
        return;
    }

    const FString Directory = FPaths::GetPath(ResumeKeyPath());
    IFileManager::Get().MakeDirectory(*Directory, true);
    FFileHelper::SaveStringToFile(
        ResumeKey,
        *ResumeKeyPath(),
        FFileHelper::EEncodingOptions::ForceUTF8WithoutBOM
    );
}

void UNDServiceSubsystem::BootstrapOrResume()
{
    if (bBootstrapInFlight)
    {
        return;
    }
    bBootstrapInFlight = true;

    TSharedRef<FJsonObject> Body = MakeShared<FJsonObject>();
    if (!ResumeKey.IsEmpty())
    {
        Body->SetStringField(TEXT("resume_key"), ResumeKey);
    }

    FString Payload;
    TSharedRef<TJsonWriter<>> Writer = TJsonWriterFactory<>::Create(&Payload);
    FJsonSerializer::Serialize(Body, Writer);

    auto Request = NewJsonRequest(ServiceBaseUrl + TEXT("/v1/sessions/bootstrap"), TEXT("POST"));
    Request->SetContentAsString(Payload);
    Request->OnProcessRequestComplete().BindUObject(this, &UNDServiceSubsystem::HandleBootstrap);
    if (!Request->ProcessRequest())
    {
        bBootstrapInFlight = false;
        BroadcastError(TEXT("bootstrap_request_start_failed"));
    }
}

void UNDServiceSubsystem::HandleBootstrap(FHttpRequestPtr Request, FHttpResponsePtr Response, bool bSucceeded)
{
    bBootstrapInFlight = false;

    if (!bSucceeded || !Response.IsValid() || Response->GetResponseCode() != 201)
    {
        OnSessionReady.Broadcast(false);
        BroadcastError(TEXT("bootstrap_failed"));
        return;
    }

    TSharedPtr<FJsonObject> Root;
    if (!ParseJsonObject(Response, Root))
    {
        OnSessionReady.Broadcast(false);
        BroadcastError(TEXT("bootstrap_invalid_json"));
        return;
    }

    FString NewSessionToken;
    if (!Root->TryGetStringField(TEXT("session_token"), NewSessionToken) || NewSessionToken.IsEmpty())
    {
        OnSessionReady.Broadcast(false);
        BroadcastError(TEXT("bootstrap_missing_session"));
        return;
    }

    FString NewResumeKey;
    if (Root->TryGetStringField(TEXT("resume_key"), NewResumeKey) && !NewResumeKey.IsEmpty())
    {
        ResumeKey = NewResumeKey;
        SaveResumeKey();
    }

    const TSharedPtr<FJsonObject>* SnapshotObject = nullptr;
    if (!Root->TryGetObjectField(TEXT("snapshot"), SnapshotObject) ||
        SnapshotObject == nullptr ||
        !ParseSnapshotObject(*SnapshotObject, Snapshot))
    {
        OnSessionReady.Broadcast(false);
        BroadcastError(TEXT("bootstrap_invalid_snapshot"));
        return;
    }

    SessionToken = NewSessionToken;
    OnSessionReady.Broadcast(true);
    OnSnapshotUpdated.Broadcast(Snapshot);
    EnsureGameplayBinding();
}

void UNDServiceSubsystem::RefreshState()
{
    if (SessionToken.IsEmpty())
    {
        BroadcastError(TEXT("session_missing"));
        return;
    }

    auto Request = NewJsonRequest(ServiceBaseUrl + TEXT("/v1/state"), TEXT("GET"), SessionToken);
    Request->OnProcessRequestComplete().BindUObject(this, &UNDServiceSubsystem::HandleState);
    if (!Request->ProcessRequest())
    {
        BroadcastError(TEXT("state_request_start_failed"));
    }
}

void UNDServiceSubsystem::HandleState(FHttpRequestPtr Request, FHttpResponsePtr Response, bool bSucceeded)
{
    if (!bSucceeded || !Response.IsValid() || Response->GetResponseCode() != 200)
    {
        BroadcastError(TEXT("state_refresh_failed"));
        return;
    }

    TSharedPtr<FJsonObject> Root;
    if (!ParseJsonObject(Response, Root) || !ParseSnapshotObject(Root, Snapshot))
    {
        BroadcastError(TEXT("state_invalid_snapshot"));
        return;
    }
    OnSnapshotUpdated.Broadcast(Snapshot);
}

void UNDServiceSubsystem::CompleteQuest(const FString& QuestID)
{
    if (SessionToken.IsEmpty() || QuestID.IsEmpty())
    {
        BroadcastError(TEXT("quest_request_invalid"));
        return;
    }

    TSharedRef<FJsonObject> Body = MakeShared<FJsonObject>();
    Body->SetStringField(
        TEXT("operation_id"),
        FGuid::NewGuid().ToString(EGuidFormats::DigitsWithHyphensLower)
    );

    FString Payload;
    TSharedRef<TJsonWriter<>> Writer = TJsonWriterFactory<>::Create(&Payload);
    FJsonSerializer::Serialize(Body, Writer);

    auto Request = NewJsonRequest(
        ServiceBaseUrl + TEXT("/v1/quests/") + QuestID + TEXT("/complete"),
        TEXT("POST"),
        SessionToken
    );
    Request->SetContentAsString(Payload);
    Request->OnProcessRequestComplete().BindUObject(this, &UNDServiceSubsystem::HandleQuest);
    if (!Request->ProcessRequest())
    {
        BroadcastError(TEXT("quest_request_start_failed"));
    }
}

void UNDServiceSubsystem::HandleQuest(FHttpRequestPtr Request, FHttpResponsePtr Response, bool bSucceeded)
{
    if (!bSucceeded || !Response.IsValid() || Response->GetResponseCode() != 200)
    {
        BroadcastError(TEXT("quest_mutation_failed"));
        return;
    }

    TSharedPtr<FJsonObject> Root;
    if (!ParseJsonObject(Response, Root))
    {
        BroadcastError(TEXT("quest_invalid_json"));
        return;
    }

    const TSharedPtr<FJsonObject>* SnapshotObject = nullptr;
    if (!Root->TryGetObjectField(TEXT("snapshot"), SnapshotObject) ||
        SnapshotObject == nullptr ||
        !ParseSnapshotObject(*SnapshotObject, Snapshot))
    {
        BroadcastError(TEXT("quest_invalid_snapshot"));
        return;
    }

    OnSnapshotUpdated.Broadcast(Snapshot);
}

void UNDServiceSubsystem::ReviseBuild(const TArray<FString>& PartIDs)
{
    if (SessionToken.IsEmpty() || Snapshot.VehicleID.IsEmpty() || Snapshot.ActiveBuildRevision < 1)
    {
        BroadcastError(TEXT("build_request_invalid"));
        return;
    }

    TSharedRef<FJsonObject> Body = MakeShared<FJsonObject>();
    Body->SetNumberField(TEXT("expected_revision"), Snapshot.ActiveBuildRevision);
    Body->SetStringField(
        TEXT("operation_id"),
        FGuid::NewGuid().ToString(EGuidFormats::DigitsWithHyphensLower)
    );

    TArray<TSharedPtr<FJsonValue>> JsonParts;
    for (const FString& PartID : PartIDs)
    {
        JsonParts.Add(MakeShared<FJsonValueString>(PartID));
    }
    Body->SetArrayField(TEXT("part_ids"), JsonParts);

    FString Payload;
    TSharedRef<TJsonWriter<>> Writer = TJsonWriterFactory<>::Create(&Payload);
    FJsonSerializer::Serialize(Body, Writer);

    auto Request = NewJsonRequest(
        ServiceBaseUrl + TEXT("/v1/vehicles/") + Snapshot.VehicleID + TEXT("/builds"),
        TEXT("POST"),
        SessionToken
    );
    Request->SetContentAsString(Payload);
    Request->OnProcessRequestComplete().BindUObject(this, &UNDServiceSubsystem::HandleBuild);
    if (!Request->ProcessRequest())
    {
        BroadcastError(TEXT("build_request_start_failed"));
    }
}

void UNDServiceSubsystem::HandleBuild(FHttpRequestPtr Request, FHttpResponsePtr Response, bool bSucceeded)
{
    if (!bSucceeded || !Response.IsValid() || Response->GetResponseCode() != 200)
    {
        BroadcastError(TEXT("build_mutation_failed"));
        return;
    }

    TSharedPtr<FJsonObject> Root;
    if (!ParseJsonObject(Response, Root) || !ParseSnapshotObject(Root, Snapshot))
    {
        BroadcastError(TEXT("build_invalid_snapshot"));
        return;
    }
    OnSnapshotUpdated.Broadcast(Snapshot);
}

void UNDServiceSubsystem::EnsureGameplayBinding()
{
    if (SessionToken.IsEmpty() || bTicketInFlight)
    {
        return;
    }

    UWorld* World = GetWorld();
    if (!World)
    {
        return;
    }

    APlayerController* Controller = UGameplayStatics::GetPlayerController(World, 0);
    if (!Cast<ANDPlayerController>(Controller))
    {
        return;
    }

    RequestGameTicket();
}

void UNDServiceSubsystem::RequestGameTicket()
{
    bTicketInFlight = true;
    auto Request = NewJsonRequest(ServiceBaseUrl + TEXT("/v1/game-tickets"), TEXT("POST"), SessionToken);
    Request->SetContentAsString(TEXT("{}"));
    Request->OnProcessRequestComplete().BindUObject(this, &UNDServiceSubsystem::HandleGameTicket);
    if (!Request->ProcessRequest())
    {
        bTicketInFlight = false;
        BroadcastError(TEXT("game_ticket_request_start_failed"));
    }
}

void UNDServiceSubsystem::HandleGameTicket(FHttpRequestPtr Request, FHttpResponsePtr Response, bool bSucceeded)
{
    bTicketInFlight = false;

    if (!bSucceeded || !Response.IsValid() || Response->GetResponseCode() != 201)
    {
        BroadcastError(TEXT("game_ticket_issue_failed"));
        return;
    }

    TSharedPtr<FJsonObject> Root;
    FString Ticket;
    if (!ParseJsonObject(Response, Root) ||
        !Root->TryGetStringField(TEXT("ticket"), Ticket) ||
        Ticket.Len() != 64)
    {
        BroadcastError(TEXT("game_ticket_invalid"));
        return;
    }

    UWorld* World = GetWorld();
    ANDPlayerController* Controller = World
        ? Cast<ANDPlayerController>(UGameplayStatics::GetPlayerController(World, 0))
        : nullptr;
    if (!Controller)
    {
        BroadcastError(TEXT("gameplay_controller_missing"));
        return;
    }

    Controller->SubmitGameTicket(Ticket);
}

bool UNDServiceSubsystem::ParseSnapshotObject(
    const TSharedPtr<FJsonObject>& Object,
    FNDPlayerSnapshot& OutSnapshot) const
{
    if (!Object.IsValid())
    {
        return false;
    }

    FNDPlayerSnapshot Parsed;
    if (!Object->TryGetStringField(TEXT("account_id"), Parsed.AccountID) ||
        !Object->TryGetStringField(TEXT("character_id"), Parsed.CharacterID) ||
        !Object->TryGetStringField(TEXT("vehicle_id"), Parsed.VehicleID))
    {
        return false;
    }

    Parsed.Money = Object->GetIntegerField(TEXT("money"));
    Parsed.XP = Object->GetIntegerField(TEXT("xp"));
    Parsed.Reputation = Object->GetIntegerField(TEXT("reputation"));
    Parsed.ActiveBuildRevision = static_cast<int32>(
        Object->GetIntegerField(TEXT("active_build_revision"))
    );
    Object->TryGetBoolField(TEXT("starter_lineage"), Parsed.bStarterLineage);
    Object->TryGetBoolField(TEXT("roadworthy"), Parsed.bRoadworthy);

    const TArray<TSharedPtr<FJsonValue>>* PartValues = nullptr;
    if (Object->TryGetArrayField(TEXT("active_part_ids"), PartValues) && PartValues)
    {
        for (const TSharedPtr<FJsonValue>& Value : *PartValues)
        {
            FString PartID;
            if (Value.IsValid() && Value->TryGetString(PartID))
            {
                Parsed.ActivePartIDs.Add(PartID);
            }
        }
    }

    const TArray<TSharedPtr<FJsonValue>>* InventoryValues = nullptr;
    if (Object->TryGetArrayField(TEXT("inventory"), InventoryValues) && InventoryValues)
    {
        for (const TSharedPtr<FJsonValue>& Value : *InventoryValues)
        {
            if (!Value.IsValid() || Value->Type != EJson::Object)
            {
                continue;
            }
            const TSharedPtr<FJsonObject> ItemObject = Value->AsObject();
            if (!ItemObject.IsValid())
            {
                continue;
            }
            FNDInventoryItem Item;
            if (ItemObject->TryGetStringField(TEXT("item_id"), Item.ItemID))
            {
                Item.Quantity = ItemObject->GetIntegerField(TEXT("quantity"));
                Parsed.Inventory.Add(MoveTemp(Item));
            }
        }
    }

    const TArray<TSharedPtr<FJsonValue>>* BlueprintValues = nullptr;
    if (Object->TryGetArrayField(TEXT("blueprints"), BlueprintValues) && BlueprintValues)
    {
        for (const TSharedPtr<FJsonValue>& Value : *BlueprintValues)
        {
            FString BlueprintID;
            if (Value.IsValid() && Value->TryGetString(BlueprintID))
            {
                Parsed.Blueprints.Add(BlueprintID);
            }
        }
    }

    const TArray<TSharedPtr<FJsonValue>>* QuestValues = nullptr;
    if (Object->TryGetArrayField(TEXT("completed_quests"), QuestValues) && QuestValues)
    {
        for (const TSharedPtr<FJsonValue>& Value : *QuestValues)
        {
            FString QuestID;
            if (Value.IsValid() && Value->TryGetString(QuestID))
            {
                Parsed.CompletedQuests.Add(QuestID);
            }
        }
    }

    OutSnapshot = MoveTemp(Parsed);
    return true;
}

void UNDServiceSubsystem::BroadcastError(const FString& Code)
{
    OnServiceError.Broadcast(Code);
}
