# Deployment

## Deployment philosophy

Deployment complexity is phase-gated. The project does not adopt Kubernetes merely to satisfy an “enterprise” label.

## Local / developer stack

Implemented:
- Go service container,
- PostgreSQL,
- Redis reservation for later ephemeral workloads,
- Docker Compose orchestration.

Default host ports intentionally avoid common conflicts:
- API: 18080
- PostgreSQL: 55432
- Redis: 56379

## Unreal development

The Unreal project requires UE 5.8.x. Dedicated source-build CI is manual and runs on a self-hosted Linux runner labeled `unreal-5.8`.

A workflow definition is not build evidence. Retain successful logs/artifacts before marking the source-build gate complete.

## Alpha deployment target

Minimum deployable topology:
- TLS ingress/load balancer,
- one or more Go service instances,
- PostgreSQL with tested backup,
- Redis only for workloads that tolerate/recover from ephemeral loss,
- one or more Unreal dedicated gameplay-server instances,
- centralized logs/metrics/traces,
- external secret management.

## Kubernetes threshold

Move to k3s/Kubernetes when measured needs justify:
- many gameplay server instances,
- autoscaling/scheduling,
- rolling deployments,
- failover requirements,
- operational multi-service complexity.

## Environment separation

Use separate local/dev/staging/production credentials and data. Never use production secrets in CI fixtures or committed files.

## Release deployment requirements

Before production:
- immutable versioned artifacts,
- schema migration plan,
- health/readiness checks,
- rollback plan,
- secret rotation path,
- observability dashboards/alerts,
- backup age/restore verification,
- deployment evidence recorded by commit/release.
