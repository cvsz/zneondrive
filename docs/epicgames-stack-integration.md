# Epic Games Stack Integration

## Scope

This document records the Epic Games repositories and capabilities that are directly relevant to zNeonDrive's current Unreal Engine 5.8 architecture. It is an integration baseline, not a claim that every Epic subsystem is already production-ready in zNeonDrive.

## Current zNeonDrive baseline

- Unreal Engine 5.8 client and dedicated gameplay server.
- Go service plane remains authoritative for durable account, character, vehicle, build, quest, inventory, and race state.
- PostgreSQL remains the durable source of truth.
- Redis remains ephemeral coordination/abuse-control state.
- `game/NeonDrive.uproject` currently targets UE 5.8 and contains only the NeonDrive runtime module.
- `NeonDrive.Build.cs` currently depends on core UE runtime, HTTP, and JSON modules; no Epic Online Services or Pixel Streaming dependency is hard-coded.

## Epic repositories reviewed

### 1. EpicGames/EOS-Getting-Started — adopt as the online-services reference

The repository targets Unreal Engine 5.8 and demonstrates the Online Subsystem EOS integration path. Its current sample scope includes authentication, stats/achievements/leaderboards, Ecom/IAP, P2P/lobbies/voice, sessions, friends/presence/social overlay, player/title data storage, player reports, sanctions, Easy Anti-Cheat, and trusted-server voice.

**zNeonDrive action:** use EOS as an optional platform/integration layer around the existing server-authoritative Go architecture. EOS identity, social, sessions, commerce, reports, sanctions, and anti-cheat signals must not become the durable gameplay authority.

### 2. EpicGames/PixelStreamingInfrastructure — adopt for remote-client delivery

The repository is the current home for Epic's Pixel Streaming infrastructure and maintains a UE 5.8 branch. It contains signalling, SFU, common frontend libraries, frontend implementations, and protocol test tooling. Epic documents breaking changes between UE-version branches, so the zNeonDrive integration must pin the UE 5.8-compatible infrastructure rather than tracking an arbitrary branch.

Epic also documents Pixel Streaming 2 for UE 5.5+ and provides a migration path from the original Pixel Streaming plugin.

**zNeonDrive action:** treat Pixel Streaming as a delivery/operations plane, not a gameplay authority. Add it as an optional deployment profile after packaged UE 5.8 client/server evidence exists.

### 3. EpicGames/UnrealEngine — engine source baseline

Epic identifies Unreal Engine source as a licensed repository requiring Epic account/GitHub linkage and acceptance of the applicable Unreal Engine license. The repository is therefore not treated as an ordinary open-source dependency.

**zNeonDrive action:** retain UE 5.8 as the explicit engine contract. Do not vendor Epic's licensed engine source into this repository. Build evidence must continue to come from an authorized UE 5.8 environment.

### 4. EpicGames/BlenderTools — content-pipeline reference

Epic's BlenderTools repository provides Blender tooling intended to improve workflows between Blender and Unreal Engine.

**zNeonDrive action:** evaluate this only for the asset-authoring pipeline. Runtime code should not depend on Blender tooling.

### 5. EpicGames/include-what-you-use / EpicGames/Linter — code-quality references

These repositories are useful references for C++ include hygiene and Unreal-oriented linting/tooling. They should be evaluated as CI/tooling improvements rather than runtime dependencies.

## Target integration architecture

```text
                         +-----------------------------+
                         |       Epic Services         |
                         | EOS Auth / Social / EAC      |
                         | Sessions / Reports / etc.    |
                         +--------------+--------------+
                                        |
                                        | optional platform adapter
                                        v
+----------------+       +-------------+-------------+       +------------------+
| Web / Remote   |<----->| UE 5.8 Client / Dedicated |<----->| Go Game API      |
| Pixel Streaming| WebRTC| Server                     | HTTPS | authoritative API |
+----------------+       +-------------+-------------+       +--------+---------+
                                                                        |
                                                               +--------+--------+
                                                               | PostgreSQL      |
                                                               | durable state  |
                                                               +-----------------+
                                                                        |
                                                               +--------+--------+
                                                               | Redis           |
                                                               | coordination   |
                                                               +-----------------+
```

## Non-negotiable boundaries

1. EOS credentials/tokens must not be persisted as durable gameplay credentials in PostgreSQL.
2. The dedicated server must continue to receive short-lived gameplay tickets rather than long-lived player credentials.
3. Pixel Streaming must never be used as an authorization boundary.
4. EOS availability must not make core PvE/gameplay persistence unavailable; define explicit degraded-mode behavior.
5. Ecom/VIP entitlements must map into the existing fairness policy and must not directly mutate competitive vehicle performance.
6. Anti-cheat signals are evidence/inputs to enforcement workflows; authoritative race acceptance remains server-side.
7. Epic licensed source must not be copied into this public repository.

## Implementation order

### Track A — EOS foundation

- Add an optional UE plugin configuration profile for Online Subsystem EOS.
- Add a platform adapter boundary so EOS-specific code is isolated from `NeonDrive` gameplay/domain code.
- Add login/session bootstrap mapping from EOS identity to the existing Go account/session model.
- Add EOS session/presence integration without replacing the Go authoritative session/ticket model.
- Add EAC/report/sanction integration only after the basic identity/session path is evidenced.

### Track B — Pixel Streaming

- Pin the Epic Pixel Streaming Infrastructure UE 5.8 branch/release used by deployment.
- Add a dedicated streaming deployment profile and health checks.
- Add signalling/SFU observability to the existing private metrics model.
- Add packaged client streaming E2E evidence.
- Keep direct native client and Pixel Streaming delivery paths functionally equivalent.

### Track C — Content tooling

- Evaluate Epic BlenderTools against the current character/vehicle import pipeline.
- Add automated asset validation before Unreal packaging.
- Evaluate Epic C++ lint/include tooling for CI without making it a runtime dependency.

## Evidence gates

The following must remain open until executable evidence exists:

- UE 5.8 Client/Server source build.
- UE 5.8 package/cook.
- EOS authenticated client + dedicated-server gameplay-ticket E2E.
- EOS reconnect/session recovery.
- Pixel Streaming packaged-client E2E.
- Pixel Streaming multi-instance lifecycle and cleanup.
- EOS anti-cheat/report/sanction event flow.
- Deployment-scale streaming/load/soak evidence.

## Source references

- EpicGames organization: https://github.com/EpicGames
- EOS Getting Started: https://github.com/EpicGames/EOS-Getting-Started
- Pixel Streaming Infrastructure: https://github.com/EpicGames/PixelStreamingInfrastructure
- Unreal Engine access information: https://github.com/EpicGames/Signup
- BlenderTools: https://github.com/EpicGames/BlenderTools
- Include-What-You-Use: https://github.com/EpicGames/include-what-you-use
