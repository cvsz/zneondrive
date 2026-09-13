#include "NDPlayerController.h"

#include "Dom/JsonObject.h"
#include "HAL/PlatformMisc.h"
#include "HttpModule.h"
#include "NDServerTelemetry.h"
#include "NDServiceSubsystem.h"
#include "NDVehiclePawn.h"
#include "Serialization/JsonReader.h"
#include "Serialization/JsonSerializer.h"

namespace
{
bool IsHexTicket(const FString& Ticket)
{
    if (Ticket.Len() != 64)
    {
        return false;
    }
    for (TCHAR Character : Ticket)
    {
        if (!FChar::IsHexDigit(Character))
        {
            return false;
        }
    }
    return true;
}
}

void ANDPlayerController::BeginPlay()
{
    Super::BeginPlay();

    if (IsLocalController())
    {
        if (UGameInstance* GameInstance = GetGameInstance())
        {
            if (UNDServiceSubsystem* Service = GameInstance->GetSubsystem<UNDServiceSubsystem>())
            {
                Service->EnsureGameplayBinding();
            }
        }
    }
}

void ANDPlayerController::SubmitGameTicket(const FString& Ticket)
{
    if (!IsLocalController() || !IsHexTicket(Ticket))
    {
        return;
    }
    ServerSubmitGameTicket(Ticket);
}

void ANDPlayerController::ServerSubmitGameTicket_Implementation(const FString& Ticket)
{
    if (!HasAuthority() || !IsHexTicket(Ticket))
    {
        return;
    }

    FNDServerTelemetry::RecordTicketRedeemAttempt();

    FString GameServerKey = FPlatformMisc::GetEnvironmentVariable(TEXT("ZNEON_GAME_SERVER_KEY"));
    GameServerKey.TrimStartAndEndInline();
    if (GameServerKey.Len() < 32)
    {
        FNDServerTelemetry::RecordTicketRedeemFailure();
        UE_LOG(LogTemp, Error, TEXT("ZNEON_GAME_SERVER_KEY is missing or too short; gameplay binding denied"));
        return;
    }

    FString ServiceUrl = FPlatformMisc::GetEnvironmentVariable(TEXT("ZNEON_GAME_API_INTERNAL_URL"));
    ServiceUrl.TrimStartAndEndInline();
    if (ServiceUrl.IsEmpty())
    {
        ServiceUrl = TEXT("http://127.0.0.1:18080");
    }
    ServiceUrl.RemoveFromEnd(TEXT("/"));

    TSharedRef<FJsonObject> Body = MakeShared<FJsonObject>();
    Body->SetStringField(TEXT("ticket"), Ticket);

    FString Payload;
    TSharedRef<TJsonWriter<>> Writer = TJsonWriterFactory<>::Create(&Payload);
    FJsonSerializer::Serialize(Body, Writer);

    TSharedRef<IHttpRequest, ESPMode::ThreadSafe> Request = FHttpModule::Get().CreateRequest();
    Request->SetURL(ServiceUrl + TEXT("/v1/internal/game-tickets/redeem"));
    Request->SetVerb(TEXT("POST"));
    Request->SetHeader(TEXT("Accept"), TEXT("application/json"));
    Request->SetHeader(TEXT("Content-Type"), TEXT("application/json"));
    Request->SetHeader(TEXT("X-Game-Server-Key"), GameServerKey);
    Request->SetContentAsString(Payload);
    Request->OnProcessRequestComplete().BindUObject(this, &ANDPlayerController::HandleTicketRedeemed);

    if (!Request->ProcessRequest())
    {
        FNDServerTelemetry::RecordTicketRedeemFailure();
        UE_LOG(LogTemp, Error, TEXT("Failed to start gameplay ticket redemption request"));
    }
}

void ANDPlayerController::HandleTicketRedeemed(
    FHttpRequestPtr Request,
    FHttpResponsePtr Response,
    bool bSucceeded)
{
    if (!HasAuthority() || !bSucceeded || !Response.IsValid() || Response->GetResponseCode() != 200)
    {
        FNDServerTelemetry::RecordTicketRedeemFailure();
        UE_LOG(LogTemp, Warning, TEXT("Gameplay ticket redemption rejected"));
        return;
    }

    TSharedPtr<FJsonObject> Root;
    TSharedRef<TJsonReader<>> Reader = TJsonReaderFactory<>::Create(Response->GetContentAsString());
    if (!FJsonSerializer::Deserialize(Reader, Root) || !Root.IsValid())
    {
        FNDServerTelemetry::RecordTicketRedeemFailure();
        UE_LOG(LogTemp, Warning, TEXT("Gameplay ticket redemption returned invalid JSON"));
        return;
    }

    FString VehicleID;
    if (!Root->TryGetStringField(TEXT("vehicle_id"), VehicleID) || VehicleID.IsEmpty())
    {
        FNDServerTelemetry::RecordTicketRedeemFailure();
        UE_LOG(LogTemp, Warning, TEXT("Gameplay ticket snapshot is missing vehicle_id"));
        return;
    }

    const int32 BuildRevision = static_cast<int32>(
        Root->GetIntegerField(TEXT("active_build_revision"))
    );

    bool bRoadworthy = false;
    Root->TryGetBoolField(TEXT("roadworthy"), bRoadworthy);

    TArray<FString> PartIDs;
    const TArray<TSharedPtr<FJsonValue>>* PartValues = nullptr;
    if (Root->TryGetArrayField(TEXT("active_part_ids"), PartValues) && PartValues)
    {
        for (const TSharedPtr<FJsonValue>& Value : *PartValues)
        {
            FString PartID;
            if (Value.IsValid() && Value->TryGetString(PartID))
            {
                PartIDs.Add(PartID);
            }
        }
    }

    ANDVehiclePawn* VehiclePawn = Cast<ANDVehiclePawn>(GetPawn());
    if (!VehiclePawn)
    {
        FNDServerTelemetry::RecordTicketRedeemFailure();
        UE_LOG(LogTemp, Warning, TEXT("Gameplay ticket redeemed before vehicle pawn was available"));
        return;
    }

    VehiclePawn->ApplyDurableIdentity(VehicleID, BuildRevision, PartIDs, bRoadworthy);
    FNDServerTelemetry::RecordTicketRedeemSuccess();
}
