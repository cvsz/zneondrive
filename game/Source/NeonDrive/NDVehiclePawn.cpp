#include "NDVehiclePawn.h"

#include "Camera/CameraComponent.h"
#include "Components/BoxComponent.h"
#include "Components/InputComponent.h"
#include "Components/StaticMeshComponent.h"
#include "GameFramework/SpringArmComponent.h"
#include "Net/UnrealNetwork.h"
#include "UObject/ConstructorHelpers.h"

ANDVehiclePawn::ANDVehiclePawn()
{
    PrimaryActorTick.bCanEverTick = true;
    bReplicates = true;
    SetReplicateMovement(true);

    Collision = CreateDefaultSubobject<UBoxComponent>(TEXT("Collision"));
    Collision->SetBoxExtent(FVector(115.0f, 55.0f, 35.0f));
    Collision->SetCollisionProfileName(TEXT("Pawn"));
    SetRootComponent(Collision);

    Body = CreateDefaultSubobject<UStaticMeshComponent>(TEXT("Body"));
    Body->SetupAttachment(Collision);
    Body->SetCollisionEnabled(ECollisionEnabled::NoCollision);

    static ConstructorHelpers::FObjectFinder<UStaticMesh> CubeMesh(TEXT("/Engine/BasicShapes/Cube.Cube"));
    if (CubeMesh.Succeeded())
    {
        Body->SetStaticMesh(CubeMesh.Object);
        Body->SetRelativeScale3D(FVector(2.3f, 1.1f, 0.45f));
        Body->SetRelativeLocation(FVector(0.0f, 0.0f, 20.0f));
    }

    SpringArm = CreateDefaultSubobject<USpringArmComponent>(TEXT("SpringArm"));
    SpringArm->SetupAttachment(Collision);
    SpringArm->TargetArmLength = 650.0f;
    SpringArm->SetRelativeLocation(FVector(0.0f, 0.0f, 120.0f));
    SpringArm->bEnableCameraLag = true;
    SpringArm->CameraLagSpeed = 6.0f;

    Camera = CreateDefaultSubobject<UCameraComponent>(TEXT("Camera"));
    Camera->SetupAttachment(SpringArm, USpringArmComponent::SocketName);

    AutoPossessPlayer = EAutoReceiveInput::Player0;
}

void ANDVehiclePawn::Tick(float DeltaSeconds)
{
    Super::Tick(DeltaSeconds);

    if (!HasAuthority())
    {
        return;
    }

    const float ClampedThrottle = FMath::Clamp(AuthoritativeThrottle, -1.0f, 1.0f);
    const float ClampedSteering = FMath::Clamp(AuthoritativeSteering, -1.0f, 1.0f);

    AddActorWorldRotation(FRotator(
        0.0f,
        ClampedSteering * TurnRateDegreesPerSecond * DeltaSeconds,
        0.0f
    ));

    const FVector Delta = GetActorForwardVector() * ClampedThrottle * MaxSpeedCmPerSecond * DeltaSeconds;
    FHitResult Hit;
    AddActorWorldOffset(Delta, true, &Hit);
}

void ANDVehiclePawn::SetupPlayerInputComponent(UInputComponent* PlayerInputComponent)
{
    Super::SetupPlayerInputComponent(PlayerInputComponent);

    PlayerInputComponent->BindAxis(TEXT("Throttle"), this, &ANDVehiclePawn::InputThrottle);
    PlayerInputComponent->BindAxis(TEXT("Steering"), this, &ANDVehiclePawn::InputSteering);
}

void ANDVehiclePawn::InputThrottle(float Value)
{
    const float Clamped = FMath::Clamp(Value, -1.0f, 1.0f);
    if (HasAuthority())
    {
        AuthoritativeThrottle = Clamped;
    }
    else
    {
        ServerSetDrivingInput(Clamped, AuthoritativeSteering);
    }
}

void ANDVehiclePawn::InputSteering(float Value)
{
    const float Clamped = FMath::Clamp(Value, -1.0f, 1.0f);
    if (HasAuthority())
    {
        AuthoritativeSteering = Clamped;
    }
    else
    {
        ServerSetDrivingInput(AuthoritativeThrottle, Clamped);
    }
}

void ANDVehiclePawn::ServerSetDrivingInput_Implementation(float Throttle, float Steering)
{
    AuthoritativeThrottle = FMath::Clamp(Throttle, -1.0f, 1.0f);
    AuthoritativeSteering = FMath::Clamp(Steering, -1.0f, 1.0f);
}

void ANDVehiclePawn::GetLifetimeReplicatedProps(TArray<FLifetimeProperty>& OutLifetimeProps) const
{
    Super::GetLifetimeReplicatedProps(OutLifetimeProps);

    DOREPLIFETIME(ANDVehiclePawn, AuthoritativeThrottle);
    DOREPLIFETIME(ANDVehiclePawn, AuthoritativeSteering);
}
