# 🎖️ Demo Runbook — Code to Production on OpenShift

**Audience:** Booth presenter at Ministerie van Defensie  
**Goal:** Show a single-line code change going from commit to live production in under 2 minutes  
**Total demo time:** ~5 minutes (including explanation)

---

## ⏱️ Pre-Demo Checklist (do this 10 minutes before)

- [ ] OpenShift cluster is accessible and healthy
- [ ] ArgoCD is synced — `paas-demo-app` Deployment is running
- [ ] Dashboard is visible in browser at the Route URL (bookmark it)
- [ ] Terminal is open in the repo directory, logged in to OpenShift (`oc whoami`)
- [ ] Git remote is configured and you can push
- [ ] **Pre-warm the Tekton pipeline** — run it once so builder images are cached:
  ```bash
  oc create -f tekton/pipelinerun.yaml -n paas-demo
  # Wait for it to complete (~90s first run, ~45s warm)
  tkn pipelinerun logs --last -f -n paas-demo
  ```
- [ ] Dashboard shows **THREAT LEVEL: GREEN** ✅

---

## 🎬 Live Demo Script

### Step 1 — Show the dashboard (30 seconds)

> *"This is our Mission Status Dashboard, running live on OpenShift at the Ministerie van Defensie. It shows mission readiness, unit status, and active operations. Notice the Threat Level indicator — currently GREEN."*

Point to the browser showing the dashboard.

---

### Step 2 — Show the code (30 seconds)

Open `internal/config/config.go` in the editor and point to line:

```go
const ThreatLevel = "GREEN"
```

> *"This is the entire configuration for the Threat Level. One constant. One line. That's it."*

---

### Step 3 — Make the change (15 seconds)

Change the line to:

```go
const ThreatLevel = "RED"
```

> *"I'm changing the Threat Level to RED. One line of code."*

---

### Step 4 — Commit and push (15 seconds)

```bash
git add internal/config/config.go
git commit -m "chore: escalate threat level to RED"
git push origin main
```

> *"Committed. Pushed. The pipeline starts automatically."*

---

### Step 5 — Watch the pipeline (60–90 seconds)

Open the OpenShift Console → Pipelines → `paas-demo-app-pipeline`

Or watch in terminal:
```bash
tkn pipelinerun logs --last -f -n paas-demo
```

> *"OpenShift is now building the new container image, pushing it to the registry, and deploying it — all automatically. No manual steps."*

Point out the three stages: **git-clone → buildah → deploy**

---

### Step 6 — Show the result (15 seconds)

Refresh the browser.

> *"Done. The dashboard now shows THREAT LEVEL: RED — with a pulsing red alert. From a single line of code to live production in under 2 minutes."*

---

## 🎯 Threat Level Quick Reference

| Change this line | Dashboard shows | Visual effect |
|---|---|---|
| `const ThreatLevel = "GREEN"` | 🟢 THREAT LEVEL: GREEN | Green border, no animation |
| `const ThreatLevel = "AMBER"` | 🟡 THREAT LEVEL: AMBER | Amber border, no animation |
| `const ThreatLevel = "RED"` | 🔴 THREAT LEVEL: RED | Red border, pulsing glow |
| `const ThreatLevel = "BLACK"` | ⬛ THREAT LEVEL: BLACK | White-on-black, pulsing glow |

**File to edit:** [`internal/config/config.go`](internal/config/config.go), line with `const ThreatLevel`

---

## 🔧 Alternative: Config-Only Change (no rebuild)

To change the threat level without rebuilding the image, patch the Deployment env var:

```bash
oc set env deployment/paas-demo-app THREAT_LEVEL=AMBER -n paas-demo
oc rollout status deployment/paas-demo-app -n paas-demo
```

This triggers a pod restart only (no pipeline needed). Useful for a second demo variant.

---

## 🏗️ ArgoCD Application CR

To configure ArgoCD to manage this application, apply the following CR to your ArgoCD namespace:

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: paas-demo-app
  namespace: openshift-gitops   # or your ArgoCD namespace
spec:
  project: default
  source:
    repoURL: https://github.com/your-org/paas-demo-app.git
    targetRevision: main
    path: manifests
  destination:
    server: https://kubernetes.default.svc
    namespace: paas-demo
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
      - CreateNamespace=true
```

Apply with:
```bash
oc apply -f argocd-application.yaml -n openshift-gitops
```

---

## 🔑 GitHub Secrets (for GitHub Actions fallback)

Configure these secrets in your GitHub repository (Settings → Secrets → Actions):

| Secret | Description | Example |
|---|---|---|
| `REGISTRY_USERNAME` | Container registry username | `your-quay-username` |
| `REGISTRY_PASSWORD` | Container registry password or token | `your-quay-token` |
| `OC_SERVER` | OpenShift API server URL | `https://api.cluster.example.com:6443` |
| `OC_TOKEN` | OpenShift service account token | `sha256~...` |
| `OC_NAMESPACE` | Target namespace | `paas-demo` |

To create a service account token for CI:
```bash
oc create sa github-actions -n paas-demo
oc adm policy add-role-to-user edit -z github-actions -n paas-demo
oc create token github-actions -n paas-demo --duration=8760h
```

---

## 🛠️ Troubleshooting

| Problem | Solution |
|---|---|
| Pipeline fails at `buildah` | Check registry credentials: `oc get secret quay-push-secret -n paas-demo` |
| Dashboard not updating | Check rollout: `oc rollout status deployment/paas-demo-app -n paas-demo` |
| Route not accessible | Check Route: `oc get route paas-demo-app -n paas-demo` |
| Metrics not scraped | Verify user workload monitoring is enabled: `oc get configmap cluster-monitoring-config -n openshift-monitoring` |
| SCC admission failure | Verify SecurityContext in deployment.yaml matches restricted-v2 requirements |
| Pipeline takes >2 minutes | Pre-warm by running the pipeline once before the demo |
