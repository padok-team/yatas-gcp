<p align="center">
<img src="docs/auditory.png" alt="yatas-logo" width="30%">
<p align="center">

# YATAS

[![codecov](https://codecov.io/gh/StanGirard/yatas-template/branch/main/graph/badge.svg?token=OFGny8Za4x)](https://codecov.io/gh/StanGirard/YATAS) [![goreport](https://goreportcard.com/badge/github.com/stangirard/yatas-template)](https://goreportcard.com/badge/github.com/stangirard/yatas)

Yet Another Testing &amp; Auditing Solution

The goal of YATAS is to help you create a secure Cloud environment without too much hassle. It won't check for all best practices but only for the ones that are important for you based on my experience. Please feel free to tell me if you find something that is not covered.

## GCP Plugin setup

In your `.yatas.yml` file. More information on [YATAS homepage](https://github.com/padok-team/yatas).

```yaml
plugins:
  - name: "gcp"
    enabled: true
    source: "github.com/padok-team/yatas-gcp"
    version: "latest"
    description: "Check for GCP good practices"

pluginsConfiguration:
  - pluginName: "gcp"
    accounts:
      - project: "project-1"
        computeRegions:
          - europe-west1
          - europe-west2
          - europe-west3
      - project: "project-2"
        computeRegions:
          - us-east1
```

<!-- BEGIN_YATAS -->

## GCP - 25 Checks

### FUN
- GCP_FUN_001 CloudFunctions are not directly exposed on the internet
- GCP_FUN_002 CloudFunctions do not use the default Compute Engine service account
- GCP_FUN_003 CloudFunctions do not have plain text secrets in environment variables
- GCP_FUN_004 CloudFunctions require IAM authentication

### GCS
- GCP_GCS_001 Check if GCS buckets are using object versioning
- GCP_GCS_002 Check if GCS buckets are encrypted with a custom KMS key
- GCP_GCS_003 Check if GCS buckets are not public

### GKE
- GCP_GKE_001 GKE Control Plane is Regional (HA)
- GCP_GKE_002 Workload Identity is enabled
- GCP_GKE_003 GKE Control Plane does not have a public endpoint

### IAM
- GCP_IAM_001 Service accounts cannot escalate privileges
- GCP_IAM_002 Service accounts keys are not older than 90 days

### LB
- GCP_LB_001 Check if SSL certificates attached to HTTPS forwarding rules are in auto-renewed (managed mode)

### RUN
- GCP_RUN_001 CloudRun services are not directly exposed on the internet
- GCP_RUN_002 CloudRun services do not use the default Compute Engine service account
- GCP_RUN_003 CloudRun services do not have plain text secrets in environment variables

### SQL
- GCP_SQL_001 Check if SQL Instances are Regional (HA)
- GCP_SQL_002 Check if SQL Instances have backups enabled with Point in Time Recovery
- GCP_SQL_003 Check if SQL Instances have encrypted traffic enforced
- GCP_SQL_004 Check if SQL Instances are not exposed with a public IP
- GCP_SQL_005 Check if SQL Instances are encrypted at rest with a customer-managed key
- GCP_SQL_006 Check if SQL backups are stored in a multi-regional location

### VM
- GCP_VM_001 Check if VM instance is not using a public IP address
- GCP_VM_002 Check if VM Disk is encrypted with a customer-managed key
- GCP_VM_003 Check VM instance group is multi-zonal

<!-- END_YATAS -->
