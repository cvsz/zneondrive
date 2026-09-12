# Support

PROJECT: NEON DRIVE is under active development. This repository is not yet a production live-service support channel.

## Where to ask

- **Bug in repository code or documentation:** open a GitHub issue using the bug-report template.
- **Feature/design proposal:** use the feature-request template.
- **Security vulnerability:** follow [SECURITY.md](./SECURITY.md). Do not post exploitable details publicly.
- **Commercial/client questions:** use the client package under [client/](./client/).
- **Operational incident in a future deployed environment:** follow [docs/incident-response.md](./docs/incident-response.md).

## What to include

For technical reports include:
- commit SHA or release/tag,
- environment and platform,
- exact reproduction steps,
- expected vs actual behavior,
- logs with secrets and personal data removed,
- screenshots/video where useful,
- whether the issue is deterministic.

## Current support boundary

Supported today:
- repository build/validation tooling,
- Python reference oracle,
- Go service-plane prototype,
- PostgreSQL persistence prototype,
- Unreal source baseline,
- design/content catalogs and documentation.

Not yet represented as supported production service:
- public game service,
- ranked live competition,
- production anti-cheat,
- production HA/DR,
- console certification,
- 24/7 player support.

Response times are best-effort until a commercial support SLA is explicitly agreed.
