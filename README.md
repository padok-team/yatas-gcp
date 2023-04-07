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
      - project: "project-2"
```

<!-- BEGIN_YATAS -->

## GCP - 1 Checks

### GCS
- GCP_GCS_001 GCS buckets are versioned

<!-- END_YATAS -->
