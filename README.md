# openshift-to-cloud-operator

An operator that scans an OpenShift namespace and generates **cloud-portable Kubernetes manifests**: Currently, it converts to EKS supported manifests.. 
- **DeploymentConfig → Deployment**
- **Route → Ingress** (with a chosen ingress class, e.g. `alb`)
- **Service → Service**
- **PVC → PVC** (with optional storageClass mapping, e.g. `gp3` on EKS)

The operator writes the result into a ConfigMap (`converted.yaml`) you can apply to EKS/AKS/GKE.

---

## Prerequisites

- Go 1.22+
- `operator-sdk` and `kubectl/oc`
- An OpenShift cluster (OpenShift Local works great!!!!)
- Docker Hub account (for publishing)

---

## Quick Start (local dev loop)

```bash
# 1) Generate CRDs + manifests and install CRDs
make generate
make install

# 2) Run the controller locally (uses your kubeconfig context)
make run
