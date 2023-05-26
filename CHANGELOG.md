# Changelog

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.

## [1.3.0](https://github.com/padok-team/yatas-gcp/compare/v1.2.1...v1.3.0) (2023-05-26)


### Features

* **checks/lb:** check if Forwarding Rules have SSL certificated auto-renewed ([29b0f05](https://github.com/padok-team/yatas-gcp/commit/29b0f05b61dfe4aa1af04ee269736746cef0d3c8))
* **checks/lb:** WIP: getter for forwarding rules ([1d5d9bb](https://github.com/padok-team/yatas-gcp/commit/1d5d9bba8f97c1b99e357003c533f5c0507a78b6))
* **checks/sql:** add two first checks for SQL ([7a2d96d](https://github.com/padok-team/yatas-gcp/commit/7a2d96d4477f9cd09e135093a9591515f13b55b2))
* **checks/sql:** GCP_SQL_003: check RequireSsl on instances ([decaf52](https://github.com/padok-team/yatas-gcp/commit/decaf529d0f03af2fda44d620a260dd5501dd500))
* **checks/sql:** GCP_SQL_004: check no public IP on instances ([55bfeea](https://github.com/padok-team/yatas-gcp/commit/55bfeea9ca11cb3a0906d000f750a395c6133916))
* **checks/sql:** GCP_SQL_005: check instance encrypted with KMS key ([e1db4da](https://github.com/padok-team/yatas-gcp/commit/e1db4da22a347e0057b835c811a908c55c77c272))
* **README.md:** update documentation ([17ad2e1](https://github.com/padok-team/yatas-gcp/commit/17ad2e18bc5cafc6f6052956bf133bf822afb7cf))


### Bug Fixes

* **checks/lb:** add TODO comment for SSLProxies ([db894a5](https://github.com/padok-team/yatas-gcp/commit/db894a545b7ee24dbab5a210e67adf6e55305519))
* **checks/lb:** rename forwarding rules type ([feb3311](https://github.com/padok-team/yatas-gcp/commit/feb331124a20afe1b2d7fcf6242e7944299c7e8d))

### [1.2.1](https://github.com/padok-team/yatas-gcp/compare/v1.2.0...v1.2.1) (2023-05-19)


### Bug Fixes

* **checks/gcs:** handle error when getting bucket policy ([bec174c](https://github.com/padok-team/yatas-gcp/commit/bec174ca168e8f5e5686b66b0334b413a8a275be))

## [1.2.0](https://github.com/padok-team/yatas-gcp/compare/v1.1.0...v1.2.0) (2023-04-24)


### Features

* **compute:** Add instance checks GCS_VM_ ([#11](https://github.com/padok-team/yatas-gcp/issues/11)) ([4da74a7](https://github.com/padok-team/yatas-gcp/commit/4da74a7d527976ab46bdfc76d3112fc96e2745e2))
* **gcs:** add checks ([#10](https://github.com/padok-team/yatas-gcp/issues/10)) ([ff84bca](https://github.com/padok-team/yatas-gcp/commit/ff84bca6c2e497d8dca50982b6a5a48118f317cf))
* **logging:** improve logs during account unmarshal ([500bf66](https://github.com/padok-team/yatas-gcp/commit/500bf663a0564b79810c4fc9cf5cc80ecb019b9a))
* **new-yatas:** update imports and function calls from YATAS ([abee8e1](https://github.com/padok-team/yatas-gcp/commit/abee8e16614043989b33fcae85076e3540d80b60))
* **refacto:** refacto with new YATAS interfaces ([#8](https://github.com/padok-team/yatas-gcp/issues/8)) ([7e0cb81](https://github.com/padok-team/yatas-gcp/commit/7e0cb816b8d7f3cf81132e07f27a1b0200b09e2a))

## 1.1.0 (2023-04-07)


### Features

* **bootstrap:** Add GCP Authentication + First GCS check ([#3](https://github.com/padok-team/yatas-gcp/issues/3)) ([5f2dedd](https://github.com/padok-team/yatas-gcp/commit/5f2dedd58ca55dd0e9a2f634399c0dfc2174c33a))
