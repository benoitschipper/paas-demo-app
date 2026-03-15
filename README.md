# 🎖️ Defensie Mission Control — OpenShift Demo App

A military-themed Mission Status Dashboard built for a live booth demonstration at the **Ministry of Defence (Ministerie van Defensie)**. It showcases how fast a code change goes from commit to production on OpenShift.

```
┌─────────────────────────────────────────────────────────────────┐
│  MINISTERIE VAN DEFENSIE // OPERATIONEEL CENTRUM                │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │         THREAT LEVEL: GREEN                             │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────┐  │
│  │ MISSION      │  │ UNIT STATUS  │  │ ACTIVE OPERATIONS    │  │
│  │ READINESS    │  │ ALPHA-1 ✓    │  │ OP IRON SENTINEL     │  │
│  │    94%       │  │ BRAVO-2 ✓    │  │ OP STEEL HORIZON     │  │
│  │ ████████░░   │  │ CHARLIE-3 ⚠  │  │ OP CYBER SHIELD      │  │
│  └──────────────┘  └──────────────┘  └──────────────────────┘  │
│                                                                 │
│  PLATFORM: OpenShift  │  UPTIME: 2h 14m  │  BUILD: abc1234    │
└─────────────────────────────────────────────────────────────────┘
```

## 🎯 The Demo

Change **one line** of code → push → watch the pipeline deploy → see the dashboard update live.

```go
// internal/config/config.go — change this line:
const ThreatLevel = "GREEN"   // → "AMBER", "RED", or "BLACK"
```

## 🏗️ Architecture

### Demo Flow (actual)

For the demo we pre-build the container image via **GitHub Actions** (see [`.github/workflows/`](.github/workflows/)) and push it to the registry. This keeps the demo snappy — no waiting for an in-cluster build during the limited booth time.

```mermaid
flowchart TD
    A([👨‍💻 Code change<br/>push to GitHub]) --> B[GitHub Actions<br/>CI/CD workflow]
    B --> C[(Container Registry<br/>e.g. ghcr.io)]
    C --> D

    subgraph ArgoCD Bootstrap ["ArgoCD Bootstrap (argocd/bootstrap/)"]
        D[ArgoCD Application<br/>paas-demo-build<br/>→ argocd/build/]
        E[ArgoCD Application<br/>paas-demo-deploy-dev<br/>→ argocd/deploy/dev/]
        F[ArgoCD Application<br/>paas-demo-deploy-tst<br/>→ argocd/deploy/tst/]
        G[ArgoCD Application<br/>paas-demo-deploy-acc<br/>→ argocd/deploy/acc/]
        H[ArgoCD Application<br/>paas-demo-deploy-prd<br/>→ argocd/deploy/prd/]
        I[ArgoCD Application<br/>paas-demo-observe<br/>→ argocd/deploy/observe/]
    end

    D -->|Syncs Tekton resources| NS_TEKTON[Namespace<br/>example-paas-tekton<br/>Pipeline + RBAC + SA]
    E -->|Syncs manifests| NS_DEV[Namespace<br/>example-paas-dev<br/>Deployment · Service · Route]
    F -->|Syncs manifests| NS_TST[Namespace<br/>example-paas-tst<br/>Deployment · Service · Route]
    G -->|Syncs manifests| NS_ACC[Namespace<br/>example-paas-acc<br/>Deployment · Service · Route]
    H -->|Syncs manifests| NS_PRD[Namespace<br/>example-paas-prd<br/>Deployment · Service · Route]
    I -->|Syncs GrafanaDashboard CR| NS_GRAFANA[Namespace<br/>example-paas-grafana<br/>GrafanaDashboard]

    NS_DEV --> APP([🖥️ Mission Status Dashboard<br/>Go HTTP Server :8080])
    NS_TST --> APP
    NS_ACC --> APP
    NS_PRD --> APP
    NS_GRAFANA --> GRAFANA([📊 Grafana Dashboard<br/>Defensie Mission Control])
```

### Full GitOps Flow (with in-cluster Tekton build)

The [`argocd/build/`](argocd/build/) directory contains a **reference Tekton setup** showing how an in-cluster build pipeline would look. It is **not used during the demo** — it is included as an example of how you would wire up Tekton on OpenShift for a real production pipeline.

```
Git Push
   │
   ▼
Tekton Pipeline (on OpenShift)          ← example only (argocd/build/)
   ├── git-clone
   ├── buildah (build + push to registry)
   └── oc rollout (trigger deployment)
         │
         ▼
   ArgoCD (GitOps sync)
         │
         ▼
   OpenShift Deployment (per environment)
         │
         ▼
   Go HTTP Server (port 8080)
   ├── GET /          → Mission Status Dashboard (HTMX + Tailwind CDN)
   ├── GET /metrics   → Prometheus text format
   ├── GET /health    → Liveness probe {"status":"ok"}
   └── GET /ready     → Readiness probe {"status":"ready"}
```

**Stack:** Go · HTMX · Tailwind CSS (CDN) · Prometheus · Tekton (example) · ArgoCD · OpenShift 4.x

## 🚀 Quick Start (local)

```bash
# Run locally
make run

# Run with a different threat level
THREAT_LEVEL=RED make run

# Run tests
make test

# Build binary
make build
```

Open http://localhost:8080

## 🐳 Container Build

```bash
# Build image
make docker-build REGISTRY=quay.io/your-org IMAGE_TAG=dev

# Run with read-only filesystem (restricted-v2 SCC compatible)
docker run --read-only -p 8080:8080 quay.io/your-org/paas-demo-app:dev
```

## ⚙️ Environment Variables

| Variable | Default | Description |
|---|---|---|
| `THREAT_LEVEL` | `"GREEN"` | Threat level: `GREEN`, `AMBER`, `RED`, `BLACK` |
| `PORT` | `"8080"` | HTTP listen port |

## 📁 Project Structure

```
.
├── cmd/server/main.go          # Entry point — HTTP server + graceful shutdown
├── internal/
│   ├── config/config.go        # ⭐ ThreatLevel constant lives here
│   ├── handlers/               # HTTP handlers (dashboard, health, ready)
│   └── metrics/metrics.go      # Prometheus metrics registration
├── templates/dashboard.html    # Military-themed HTML template (HTMX + Tailwind)
├── argocd/
│   ├── bootstrap/              # ArgoCD Applications — apply once to bootstrap the platform
│   │   ├── a_build-application.yaml   # → deploys Tekton resources into example-paas-tekton
│   │   ├── a_deploy-dev.yaml          # → deploys app into example-paas-dev
│   │   ├── a_deploy-tst.yaml          # → deploys app into example-paas-tst
│   │   ├── a_deploy-acc.yaml          # → deploys app into example-paas-acc
│   │   ├── a_deploy-prd.yaml          # → deploys app into example-paas-prd
│   │   ├── a_observe-applicaton.yaml  # → deploys Grafana dashboard into example-paas-grafana
│   │   └── kustomization.yaml
│   ├── build/                  # ⚠️ EXAMPLE ONLY — reference Tekton setup (not used in demo)
│   │   ├── pipeline.yaml              # Tekton Pipeline: git-clone → buildah → push
│   │   ├── deploy-pipeline.yaml       # Alternative pipeline with oc rollout step
│   │   ├── serviceaccount.yaml        # SA for the pipeline
│   │   ├── rbac.yaml                  # RoleBindings for build + image-push
│   │   └── kustomization.yaml
│   └── deploy/                 # Kustomize overlays per environment
│       ├── generic/            # Base manifests (Deployment, Service, Route, ServiceMonitor, …)
│       │   └── rolebindings.yaml      # Grants prometheus-k8s SA view access to app namespaces
│       ├── dev/                # Dev overlay
│       ├── tst/                # Test overlay
│       ├── acc/                # Acceptance overlay
│       ├── prd/                # Production overlay
│       └── observe/            # Grafana dashboard (GrafanaDashboard CR)
│           └── grafana-dashboard.yaml # "Defensie Mission Control — Demo App" dashboard
├── .github/workflows/          # GitHub Actions — pre-builds the image for the demo
├── Containerfile               # Multi-stage build → distroless nonroot (~20 MB)
└── DEMO_RUNBOOK.md             # Step-by-step demo script
```

## 🔒 Security (restricted-v2 SCC)

The container runs with:
- `runAsNonRoot: true` (UID 65532 via distroless nonroot)
- `allowPrivilegeEscalation: false`
- `capabilities.drop: ["ALL"]`
- `readOnlyRootFilesystem: true`
- `seccompProfile.type: RuntimeDefault`

## 📊 Prometheus Metrics

| Metric | Type | Description |
|---|---|---|
| `demo_page_views_total` | Counter | Dashboard page loads |
| `demo_threat_level` | Gauge | 0=GREEN, 1=AMBER, 2=RED, 3=BLACK |
| `demo_active_sessions` | Gauge | In-flight HTTP requests |
| `demo_request_duration_seconds` | Histogram | Request latency by handler |

## 🔭 Observability — Grafana Dashboard

A **`GrafanaDashboard` custom resource** is deployed via ArgoCD into the `example-paas-grafana` namespace. It is managed by the [Grafana Operator](https://github.com/grafana/grafana-operator) and automatically picked up by any Grafana instance labelled `dashboards: grafana`.

The dashboard is titled **"Defensie Mission Control — Demo App"** and visualises all four Prometheus metrics exposed by the app, with a namespace selector so you can switch between `dev`, `tst`, `acc`, and `prd` at runtime.

![Grafana dashboard showing Threat Level, Page Views, Active Sessions and request latency panels](argocd/deploy/observe/dashboard.png)

### Dashboard panels

| Panel | Type | Metric |
|---|---|---|
| Threat Level | Stat (colour-coded) | `demo_threat_level` |
| Total Page Views | Stat | `demo_page_views_total` |
| Active Sessions | Stat | `demo_active_sessions` |
| Page View Rate | Stat | `rate(demo_page_views_total[1m])` |
| Page Views Rate (over time) | Time series | `rate(demo_page_views_total[1m])` |
| Active Sessions (over time) | Time series | `demo_active_sessions` |
| Request Duration — p50/p90/p99 | Time series | `histogram_quantile` on `demo_request_duration_seconds_bucket` |
| Request Rate by Handler | Time series | `demo_request_duration_seconds_count` by handler |
| Threat Level (over time) | Time series | `demo_threat_level` |

### Prometheus RBAC

[`argocd/deploy/generic/rolebindings.yaml`](argocd/deploy/generic/rolebindings.yaml) grants the `prometheus-k8s` service account (from `openshift-monitoring`) the built-in `view` ClusterRole in each app namespace, so the OpenShift monitoring stack can scrape the `/metrics` endpoint.

## 📖 Demo Script

See [DEMO_RUNBOOK.md](DEMO_RUNBOOK.md) for the full step-by-step demo script.
