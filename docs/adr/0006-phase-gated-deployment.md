# ADR-0006: Phase-Gated Deployment — Simple First, Kubernetes After Evidence

- Status: Accepted
- Date: 2026-09-12

## Context

The project already targets persistent multiplayer and eventual large-scale operations, but adopting a full orchestration platform before gameplay/network load is measured would add complexity without proving player value.

## Decision

Use phase-gated deployment.

### Vertical slice / internal demo
- local Unreal client + dedicated server
- service plane via Docker Compose
- PostgreSQL + Redis
- deterministic seed/demo data

### Online alpha
- containerized Go service plane
- containerized Linux dedicated servers where supported
- ingress/TLS, managed or self-hosted PostgreSQL/Redis depending environment
- OpenTelemetry-compatible logs/metrics/traces
- repeatable infrastructure manifests

### Scale/production candidate
Move to k3s/Kubernetes only when concurrency, deployment frequency, failover, autoscaling, multi-instance scheduling, or operational evidence justifies it.

## Required observability

Before production-readiness claims:
- request/error/latency metrics,
- gameplay server health/session counts,
- race validation failures,
- economy/idempotency rejection metrics,
- database saturation,
- reconnect failure metrics,
- structured audit logs,
- distributed trace correlation for durable mutations.

## Consequences

This avoids building an expensive platform before the vertical slice proves the game loop, while retaining a clean path to the user's existing k3s-oriented operating model.
