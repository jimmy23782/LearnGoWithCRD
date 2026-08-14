# Gateway Operator

A Kubernetes Operator built using **Kubebuilder** and **Go** to provide a self-service onboarding platform for application teams.

The long-term vision of this project is to enable engineering teams to onboard themselves onto a Kubernetes platform by submitting simple Custom Resources (CRs), while the platform automatically provisions and manages the underlying Kubernetes and Gloo Gateway resources.

---

# Problem Statement

In many organisations, every application team implements Kubernetes onboarding and API gateway configuration differently.

This project demonstrates how a Platform Engineering team can expose a simple self-service API while hiding Kubernetes and Gloo Gateway complexity from feature teams.

Instead of asking teams to manually create Kubernetes resources, they simply submit a Tenant resource.

The platform provisions everything required.

---

# Current Version

**Version 1.0**

Implemented features:

- Custom Kubernetes API using Kubebuilder
- Cluster-scoped Tenant CRD
- Validation using OpenAPI schema
- Namespace reconciliation
- Idempotent reconciliation
- Unit tests

Current workflow:

```
Feature Team

        │

Submit Tenant YAML

        │

Kubernetes API Server

        │

Tenant Controller

        │

Creates Namespace
```

Example:

```yaml
apiVersion: platform.mac.com/v1
kind: Tenant

metadata:
  name: payments

spec:
  teamName: Payments
  cmdbTeamId: PAY001
  owners:
    - jimmy.joy@company.com
  environment: dev
```

Creates

```
payments-dev
```

---

# Roadmap

## Version 1.0

- [x] Tenant CRD
- [x] Namespace reconciliation

## Version 1.1

- [ ] Default RBAC
- [ ] ResourceQuota
- [ ] NetworkPolicy

## Version 1.2

- [ ] Helm packaging
- [ ] Argo CD deployment

## Version 2.0

- [ ] PlatformAPI CRD
- [ ] Gloo Gateway integration
- [ ] Authentication Policies
- [ ] HTTPRoute reconciliation

---

# Repository Structure

```
api/
    Kubernetes API definitions

internal/controller/
    Reconciliation logic

config/
    Generated Kubernetes manifests

app-teams/
    Sample Tenant resources used by feature teams

docs/
    Architecture and design documents
```

---

# Running Locally

## Prerequisites

- Go 1.24+
- Docker Desktop
- Kind
- kubectl
- Kubebuilder

---

## Install the CRD

```bash
make install
```

---

## Run the controller

```bash
make run
```

---

## Create a Tenant

```bash
kubectl apply -f app-teams/cr/payments.yaml
```

---

## Verify

```bash
kubectl get tenants

kubectl get namespaces
```

Expected namespace:

```
payments-dev
```

---

# Learning Goals

This repository demonstrates the complete lifecycle of building a Kubernetes Operator:

1. Define Kubernetes APIs using Go structs.
2. Generate CRDs using Kubebuilder.
3. Install CRDs into a Kubernetes cluster.
4. Implement reconciliation logic.
5. Validate Custom Resources.
6. Build a self-service Platform Engineering workflow.

---

# Future Architecture

```
Tenant

        │

        ▼

Namespace

        │

        ▼

RBAC

        │

        ▼

ResourceQuota

        │

        ▼

NetworkPolicy

        │

        ▼

PlatformAPI

        │

        ▼

Gloo Gateway

        │

        ▼

HTTPRoute
AuthenticationPolicy
RateLimitPolicy
```

---

# License

Licensed under the Apache License 2.0.