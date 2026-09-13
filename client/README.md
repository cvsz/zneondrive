# Client Meeting Package — 13 September 2026

Start here before the meeting.

## Presentation assets

1. **Interactive demo:** `demo/index.html`
2. **Executive brief:** [EXECUTIVE-BRIEF-TH.md](./EXECUTIVE-BRIEF-TH.md)
3. **Talk track:** [PRESENTATION-TALK-TRACK-TH.md](./PRESENTATION-TALK-TRACK-TH.md)
4. **Demo runbook:** [DEMO-RUNBOOK-TH.md](./DEMO-RUNBOOK-TH.md)
5. **Commercial scope:** [COMMERCIAL-SCOPE-TH.md](./COMMERCIAL-SCOPE-TH.md)
6. **Prepared Q&A:** [CLIENT-QA-TH.md](./CLIENT-QA-TH.md)

## Fast start

```bash
git pull
python3 -m http.server 8080 --directory demo
```

Open `http://localhost:8080`.

## Evidence to keep open

- `docs/game-design-bible-v0.2.md`
- `docs/vertical-slice-v0.2.md`
- `docs/architecture.md`
- `docs/adr/0004-unreal-engine-client-and-gameplay-server.md`
- `docs/adr/0005-durable-service-plane.md`
- `src/zneondrive/domain.py`
- `tests/test_reference_runtime.py`
- GitHub Actions CI / CodeQL

## Client-safe status language

Use:

> Pre-production design, vertical-slice contract, authority proof, production architecture direction, and interactive presentation demo are ready. The next commercial milestone is the playable Unreal vertical slice.

Do not use:

> The MMORPG is finished / production-ready.

## GitHub Pages option

A manual deployment workflow is included at `.github/workflows/pages-client-demo.yml`.

Before first use, repository Pages must be configured to use **GitHub Actions** as its publishing source. The workflow is manual so an unconfigured Pages setting cannot create a red deployment on normal pushes.
