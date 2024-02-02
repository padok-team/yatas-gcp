# Contributing

## Add a check to an existing category

Categories are represented by a prefix in the check code. For example, `GCP_FUN_001` is a check in the `FUN` category (CloudFunctions). If you want to add a check to an existing category, you must work in the corresponding Go package in the `gcp` directory.

Example: to add a new check for Cloud Functions, you must work in the `functions` package.

**Process:**

- (Optional) modify or add a new getter function in the `getter.go` file
- Add a new condition in the `conditions.go` file
- Add the checks `<package_name>.go` file, in the checks array

## Add a new category

If you want to add a new category, you must create a new Go package in the `gcp` directory.

For example, if you want to add a new category for BiqQuery, call it `BQ`, and create a new Go package called `bigquery`.

**Requirements:**

- A `getter.go` file
  - Write here the functions to get the resources you need to be checked (for example BigQuery datasets, tables, etc.).
- A `<package_name>.go` file
  - Define here the structures that encapsulate the resources you want to check (respect the commons.Check interface from YATAS)
  - Define the RunChecks functions that will run checks and be called by the main function of the YATAS GCP plugin
- A `conditions.go` file
  - Define here the conditions that will be used to check the resources; the functions must return a boolean value
