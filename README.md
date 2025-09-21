# openshift-to-cloud-operator

An operator that scans an OpenShift namespace and generates **cloud-portable Kubernetes manifests**.  
Currently focused on **EKS-compatible output**:

- **DeploymentConfig → Deployment**
- **Route → Ingress** (with a chosen ingress class, e.g. `alb`)
- **Service → Service**
- **PVC → PVC** (with optional storageClass mapping, e.g. `gp3` on EKS)

The operator writes the result into a ConfigMap (`converted.yaml`) that you can extract and apply to EKS/AKS/GKE.

---

## End-user Usage

> **Note:** This operator is currently experimental and **not production ready**.  
> Instructions below uses the image is published to Docker Hub. For now, see the [Local Development](#-quick-start-local-development-loop) section.

1. **Install CRDs**

```bash
kubectl apply -f https://raw.githubusercontent.com/mostlycloudysky/openshift-to-cloud-operator/main/config/crd/bases/migrate.migrate.dev_migrationplans.yaml
```

2. **Deploy the operator from Docker Hub**

```bash
kubectl apply -f https://raw.githubusercontent.com/mostlycloudysky/openshift-to-cloud-operator/main/deploy/install.yaml
```

3. **Create a MigrationPlan**
   
```yml
apiVersion: migrate.migrate.dev/v1
kind: MigrationPlan
metadata:
  name: oc-to-cloud-sample
  namespace: oc-hosted-app
spec:
  namespaces: ["oc-hosted-app"]
  include: ["deploymentconfigs","routes","services","pvcs"]
  targetCloud: "eks"
  ingressClass: "alb"
  outputConfigMap: "oc-to-cloud-output"
```

4. **Fetch the converted YAML**

```bash
kubectl get cm oc-to-cloud-output -n oc-hosted-app -o jsonpath='{.data.converted\.yaml}' > converted.yaml
kubectl --context my-eks apply -f converted.yaml
```

## Features

- Automatically discovers OpenShift resources (DeploymentConfig, Route, Service, PVC) in one or more namespaces.
- Converts them into standard Kubernetes YAML.
- Adds cloud portability hints:
- IngressClass for ALB, NGINX, etc.
- storageClassName mapping for EKS (gp3).
- Outputs a multi-doc YAML bundle inside a ConfigMap for easy export.

## Prerequisites (for local dev/testing)

1. Go 1.22+
2. operator-sdk
3. kubectl and/or oc
4. Access to an OpenShift cluster (OpenShift Local / CRC works great!)
5. (Optional) Docker Hub account if you want to publish your own image

## Quick Start (Local Development Loop)

1. Generate CRDs and install them into the cluster

```bash
make generate
make install
```

2. Run the operator locally

```bash
make run
```

3. Apply a sample MigrationPlan

```yml
# config/samples/migrate_v1_migrationplan.yaml
apiVersion: migrate.migrate.dev/v1
kind: MigrationPlan
metadata:
  name: oc-to-cloud-sample
  namespace: oc-hosted-app
spec:
  namespaces: ["oc-hosted-app"]
  include: ["deploymentconfigs","routes","services","pvcs"]
  targetCloud: "eks"
  ingressClass: "alb"
  outputConfigMap: "oc-to-cloud-output"
```

4. Apply it. 

```bash
oc apply -f config/samples/migrate_v1_migrationplan.yaml
```

5. Inspect the results

```bash
oc get migrationplan oc-to-cloud-sample -n oc-hosted-app -o yaml
```

6. Extract the converted YAML bundle:

```bash
oc get cm oc-to-cloud-output -n oc-hosted-app -o jsonpath='{.data.converted\.yaml}' > converted.yaml
```

## Building & Deploying (Optional)

1. Build and push the image
```bash
make docker-build IMG=docker.io/<your-user>/openshift-to-cloud-operator:latest
make docker-push  IMG=docker.io/<your-user>/openshift-to-cloud-operator:latest
```

2. Deploy the operator
```bash
make deploy IMG=docker.io/<your-user>/openshift-to-cloud-operator:latest
```

3. Verify the operator pod
```bash
kubectl get pods -n openshift-to-cloud-operator-system
```