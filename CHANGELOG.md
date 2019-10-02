# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.3.0] - 2019-10-02
### Added
- vcl_snippets: Support for `hash` type ([#133](https://github.com/fastly/go-fastly/pull/133)).
- service: Support for `type` field ([#132](https://github.com/fastly/go-fastly/pull/132)).
- token: Support for API tokens ([#131](https://github.com/fastly/go-fastly/pull/131)).

### Changed
- client: Codebase dependency management from `dep` to Go modules ([#130](https://github.com/fastly/go-fastly/pull/130)).

## [1.2.1] - 2019-07-25
### Added
- acl: Constant to represent the maximum number of entries that can be placed within an ACL ([#129](https://github.com/fastly/go-fastly/pull/129)).
- dictionary: Constant to represent the maximum number of items that can be placed within an Edge Dictionary ([#129](https://github.com/fastly/go-fastly/pull/129)).

## [1.2.0] - 2019-07-24
### Added
- acl: Support for Create, Delete and Update BatchOperations ([#126](https://github.com/fastly/go-fastly/pull/126)).

## [1.1.0] - 2019-07-22
### Added
- dictionary: Support for Create, Delete, Update and Upsert BatchOperations ([#125](https://github.com/fastly/go-fastly/pull/125)).

## [1.0.0] - 2019-06-14
### Added
- bigquery: Support for `format_version` field ([#97](https://github.com/fastly/go-fastly/pull/97)).
- ftp: Support for `format_version` field ([#97](https://github.com/fastly/go-fastly/pull/97)).
- gcs: Support for `format_version` field ([#97](https://github.com/fastly/go-fastly/pull/97)).
- papertrail: Support for `format_version` field ([#97](https://github.com/fastly/go-fastly/pull/97)).
- backend: Support for `override_host` field ([#120](https://github.com/fastly/go-fastly/pull/120)).
- backend: Support for `created_at`, `updated_at`, and `deleted_at` fields ([#121](https://github.com/fastly/go-fastly/pull/121)).
- cache_setting: Support for `created_at`, `updated_at`, and `deleted_at` fields ([#121](https://github.com/fastly/go-fastly/pull/121)).
- condition: Support for `created_at`, `updated_at`, and `deleted_at` fields ([#121](https://github.com/fastly/go-fastly/pull/121)).
- domain: Support for `created_at`, `updated_at`, and `deleted_at` fields ([#121](https://github.com/fastly/go-fastly/pull/121)).
- gcs: Support for `created_at`, `updated_at`, and `deleted_at` fields ([#121](https://github.com/fastly/go-fastly/pull/121)).
- gzip: Support for `created_at`, `updated_at`, and `deleted_at` fields ([#121](https://github.com/fastly/go-fastly/pull/121)).
- header: Support for `created_at`, `updated_at`, and `deleted_at` fields ([#121](https://github.com/fastly/go-fastly/pull/121)).
- health_check: Support for `created_at`, `updated_at`, and `deleted_at` fields ([#121](https://github.com/fastly/go-fastly/pull/121)).
- request_setting: Support for `created_at`, `updated_at`, and `deleted_at` fields ([#121](https://github.com/fastly/go-fastly/pull/121)).
- response_object: Support for `created_at`, `updated_at`, and `deleted_at` fields ([#121](https://github.com/fastly/go-fastly/pull/121)).
- vcl: Support for `created_at`, `updated_at`, and `deleted_at` fields ([#121](https://github.com/fastly/go-fastly/pull/121)).
- version: Support for `created_at`, `updated_at`, and `deleted_at` fields ([#121](https://github.com/fastly/go-fastly/pull/121)).

### Changed
- bigquery: Function signature to list all of the BigQuery logging objects ([#97](https://github.com/fastly/go-fastly/pull/97)).
- acl: Data types of all timestamp fields ([#121](https://github.com/fastly/go-fastly/pull/121)).
- acl_entry: Data types of all timestamp fields ([#121](https://github.com/fastly/go-fastly/pull/121)).
- dictionary: Data types of all timestamp fields ([#121](https://github.com/fastly/go-fastly/pull/121)).
- dictionary_item: Data types of all timestamp fields ([#121](https://github.com/fastly/go-fastly/pull/121)).
- director: Data types of all timestamp fields ([#121](https://github.com/fastly/go-fastly/pull/121)).
- event_logs: Data types of all timestamp fields ([#121](https://github.com/fastly/go-fastly/pull/121)).
- service: Data types of all timestamp fields ([#121](https://github.com/fastly/go-fastly/pull/121)).
- vcl_snippets: Data types of all timestamp fields ([#121](https://github.com/fastly/go-fastly/pull/121)).

## [0.4.3] - 2019-05-24
### Added
- waf: Support for OWASP ([74c03ce](https://github.com/fastly/go-fastly/commit/74c03cec4549738bb1fc20fd881fbd9c750f6928), [cf4b60f](https://github.com/fastly/go-fastly/commit/cf4b60fa9e61358753b9cafa4481adc5977d8432), [cee4e1d](https://github.com/fastly/go-fastly/commit/cee4e1df35a27e507ef0809991fe8ffb462a94d4)).
- waf: Support for rules ([34f0b9f](https://github.com/fastly/go-fastly/commit/34f0b9f50526c1a4fb0a866d363dd089cd628a7c), [c1e80fc](https://github.com/fastly/go-fastly/commit/c1e80fc42a3d4c7d85ece91b682cd2db602c4df7)).
- waf: Support for rule sets ([a3cdd36](https://github.com/fastly/go-fastly/commit/a3cdd3621d72750658d3bdb988576e3401dc50f9), [bb50c5a](https://github.com/fastly/go-fastly/commit/bb50c5a9265476ce3a4161f450f6055478b84b7e)).
- loggers: Support for `placement` field ([22a781b](https://github.com/fastly/go-fastly/commit/22a781b35924ba511bda44a21d1d240026b7cf8c), [#88](https://github.com/fastly/go-fastly/pull/88), [#89](https://github.com/fastly/go-fastly/pull/89), [#92](https://github.com/fastly/go-fastly/pull/92)).
- waf: Support for rule statuses ([41be101](https://github.com/fastly/go-fastly/commit/41be101a60f72090a9e0f2da6f9916b1d3e9fc2b), [959f73b](https://github.com/fastly/go-fastly/commit/959f73bf50b63ce73e705ccee25003af983e9306), [9ac3984](https://github.com/fastly/go-fastly/commit/9ac3984c0bed8a2d459c86798caa82379a4ad64e), [f07eebf](https://github.com/fastly/go-fastly/commit/f07eebff0c9f1de723e4f09e5c6e34708ad895df), [629aea3](https://github.com/fastly/go-fastly/commit/629aea38479bd207a6e65db6eaf5b0012df665f7), [7c0a4cb](https://github.com/fastly/go-fastly/commit/7c0a4cb048d5f81fb4bfc62e116c4c58771a4dfc), [8dea9d3](https://github.com/fastly/go-fastly/commit/8dea9d3d6c37cb795edd3c90ceca2185de6f7411), [d0959bb](https://github.com/fastly/go-fastly/commit/d0959bbb9d7014852a2d156ca4f781f8b85af0f6), [787809f](https://github.com/fastly/go-fastly/commit/787809f40a2c2e1f22bda069461a570ef5c0b3c7), [313dc49](https://github.com/fastly/go-fastly/commit/313dc494e149a53cee03f858d197a8cc40bebee0), [10e52ed](https://github.com/fastly/go-fastly/commit/10e52edcd71b88cb15195ee5380f77a149089e94), [ad02a14](https://github.com/fastly/go-fastly/commit/ad02a140be06849fe79c56cb1b1da5bf9825118d)).
- waf: Support for configuration sets ([572ae53](https://github.com/fastly/go-fastly/commit/572ae535464793fddb5f393ec82f69051a747227)).
- bigquery: Support for BigQuery ([c4d7e54](https://github.com/fastly/go-fastly/commit/c4d7e54baa4a84f2e6c897ee59aeac24a4c33a9d), [ba3228c](https://github.com/fastly/go-fastly/commit/ba3228c61ddc08f8f76c348e50bafc5c389e8d30), [a5ccfe9](https://github.com/fastly/go-fastly/commit/a5ccfe98e321b197af0205f7cb45c52a913f725b), [d9957b4](https://github.com/fastly/go-fastly/commit/d9957b462a473c7a535288cbb8fccb643a34aa59), [cd9a5e6](https://github.com/fastly/go-fastly/commit/cd9a5e6c218407189f76dc13d74ab5fb59868671), [#81](https://github.com/fastly/go-fastly/pull/81)).
- event_logs: Support for event logs ([2a7fdb8](https://github.com/fastly/go-fastly/commit/2a7fdb8643082c66ee2fc248dd7aa8610b91a109), [dc37f61](https://github.com/fastly/go-fastly/commit/dc37f61f9f8392895cf4051bdeefef039f2c106f), [f833123](https://github.com/fastly/go-fastly/commit/f833123061198183771036c847988f143846f717), [6c902bc](https://github.com/fastly/go-fastly/commit/6c902bc3064338335a8eba4f301c039239309e52), [f771fa0](https://github.com/fastly/go-fastly/commit/f771fa072e71dd535853d2a52e999af90070d2af), [c7f7044](https://github.com/fastly/go-fastly/commit/c7f7044614055db8d79a1333c5d2ade78f1be0ea), [cee4e1d](https://github.com/fastly/go-fastly/commit/cee4e1df35a27e507ef0809991fe8ffb462a94d4), [10d525c](https://github.com/fastly/go-fastly/commit/10d525c164a3f6760f47ab43da8c895f7713c25c)).
- dictionary_item: Support to create multiple items ([a162398](https://github.com/fastly/go-fastly/commit/a162398362d3a967cdc5518190344b16b5421060)).
- vcl_snippets: Support for VCL snippets ([#80](https://github.com/fastly/go-fastly/pull/80), [#82](https://github.com/fastly/go-fastly/pull/82/files), [#84](https://github.com/fastly/go-fastly/pull/84), [#85](https://github.com/fastly/go-fastly/pull/85), [#96](https://github.com/fastly/go-fastly/pull/96)).
- acl: Support for `deleted_at`, `created_at`, and `updated_at` fields ([#86](https://github.com/fastly/go-fastly/pull/86)).
- acl_entry: Support for `deleted_at`, `created_at`, and `updated_at` fields ([#86](https://github.com/fastly/go-fastly/pull/86)).
- backend: Support for `comment` field ([#86](https://github.com/fastly/go-fastly/pull/86)).
- condition: Support for `name` and `comment` fields ([#86](https://github.com/fastly/go-fastly/pull/86)).
- dictionary: Support for `write_only` field ([#86](https://github.com/fastly/go-fastly/pull/86)).
- director: Support for `shield`, `name`, `deleted_at`, `created_at`, and `updated_at` fields ([#86](https://github.com/fastly/go-fastly/pull/86)).
- health_check: Support for `comment` field ([#86](https://github.com/fastly/go-fastly/pull/86)).
- service: Support to list the domains within a service ([#90](https://github.com/fastly/go-fastly/pull/90)).
- vcl: Support for `main` field ([#93](https://github.com/fastly/go-fastly/pull/93)).
- version: Support for `comment` field ([#103](https://github.com/fastly/go-fastly/pull/103)).
- splunk: Support for Splunk ([#101](https://github.com/fastly/go-fastly/pull/101)).
- blobstorage: Support for Azure Blob Storage ([#99](https://github.com/fastly/go-fastly/pull/99)).
- settings: Support for `stale_if_error` and `stale_if_error_ttl` fields ([#104](https://github.com/fastly/go-fastly/pull/104)).

### Changed
- dictionary: Response struct to align with API ([6a8a1c6](https://github.com/fastly/go-fastly/commit/6a8a1c62e61097da752bae838375cb139f4e9cc3)).
- dictionary_item: Response struct to align with API ([7d31c4a](https://github.com/fastly/go-fastly/commit/7d31c4aa34ef904f4426ef8af4be915c2c373e70)).
- user-agent: client.go and fixtures to reference fastly in the user-agent ([#109](https://github.com/fastly/go-fastly/pull/113).

### Removed
- domain: `locked` field ([#86](https://github.com/fastly/go-fastly/pull/86)).

### Fixed
- request: URL encoding for names ([8b3e2d6](https://github.com/fastly/go-fastly/commit/8b3e2d653b2d4a32ecbd050056370c95e9f6cbd8)).
- condition: Request struct for updating ([1fe3fda](https://github.com/fastly/go-fastly/commit/1fe3fda765dcbc8d75d0e2e501a2c326a0b6fafb)).

## [0.4.2] - 2017-09-05
### Added
- logentries: Support for `format_version` field ([#50](https://github.com/fastly/go-fastly/pull/50)).
- gcs: Support for `message_type` field ([#52](https://github.com/fastly/go-fastly/pull/52)).
- waf: Support for firewall ([216f9cb](https://github.com/fastly/go-fastly/commit/216f9cb6a92bc6e3c4653b7ebc9206f78a80d69b), [c6feafe](https://github.com/fastly/go-fastly/commit/c6feafe0fc5ed2b74bef9d3105f2f20c6197b19e), [50fef06](https://github.com/fastly/go-fastly/commit/50fef061051d188edcc37749a108fc6d025e495c)).

## [0.4.1] - 2017-08-07
### Added
- syslog: Support for `hostname`, `ipv4`, `tls_hostname`, and `message_type` fields ([2b863da](https://github.com/fastly/go-fastly/commit/2b863da88fc1033a68538ccdc5c9dc82fa52681f)).

## [0.4.0] - 2017-07-27
### Added
- realtime_stats: Support for real-time analytics ([#48](https://github.com/fastly/go-fastly/pull/48)).

### Changed
- acl: Names of all types, functions, and variables to follow Go standards ([#46](https://github.com/fastly/go-fastly/pull/46)).

### Fixed
- condition: URL encoding for forward slashes ([3d6dabb](https://github.com/fastly/go-fastly/commit/3d6dabb37bd2df7195d28aef08b1edd98895b960)).

## [0.3.0] - 2017-07-19

- Initial tagged release

[Unreleased]: https://github.com/fastly/go-fastly/compare/v1.3.0...HEAD
[1.3.0]: https://github.com/fastly/go-fastly/compare/v1.2.1...v1.3.0
[1.2.1]: https://github.com/fastly/go-fastly/compare/v1.2.0...v1.2.1
[1.2.0]: https://github.com/fastly/go-fastly/compare/v1.1.0...v1.2.0
[1.1.0]: https://github.com/fastly/go-fastly/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/fastly/go-fastly/compare/v0.4.3...v1.0.0
[0.4.3]: https://github.com/fastly/go-fastly/compare/v0.4.2...v0.4.3
[0.4.2]: https://github.com/fastly/go-fastly/compare/v0.4.1...v0.4.2
[0.4.1]: https://github.com/fastly/go-fastly/compare/v0.4.0...v0.4.1
[0.4.0]: https://github.com/fastly/go-fastly/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/fastly/go-fastly/releases/tag/v0.3.0
