# Changelog

## [UNRELEASED]

[Full Changelog](https://github.com/fastly/go-fastly/compare/v9.14.0...)

### Breaking:

- fix(vcl_snippets): Correct type of 'Priority' field from integer to string. [#644](https://github.com/fastly/go-fastly/pull/644)
- fix(datacenters): fix spelling of `Longitude` field [#646](https://github.com/fastly/go-fastly/pull/646)

### Enhancements:
- fix(client): add ability to pass context into the client [#647](https://github.com/fastly/go-fastly/pull/647)

### Bug fixes:

### Dependencies:
- build(deps): `golang.org/x/crypto` from 0.35.0 to 0.36.0 ([#639](https://github.com/fastly/go-fastly/pull/639))
- build(deps): `golang.org/x/mod` from 0.23.0 to 0.24.0 ([#642](https://github.com/fastly/go-fastly/pull/642))
- build(deps): `golang.org/x/tools` from 0.30.0 to 0.31.0 ([#640](https://github.com/fastly/go-fastly/pull/640))
- build(deps): `honnef.co/go/tools` from 0.6.0 to 0.6.1 ([#641](https://github.com/fastly/go-fastly/pull/641))
- build(deps): `github.com/BurntSushi/toml` from 1.4.1-0.20240526193622-a339e1f7089c to 1.5.0 ([#645](https://github.com/fastly/go-fastly/pull/645))

## [v9.14.0](https://github.com/fastly/go-fastly/releases/tag/v9.14.0) (2025-03-05)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v9.13.1...v9.14.0)

### Enhancements:

 - feat(fastly/objectstorage): adds crud operations for access keys in object storage [#612](https://github.com/fastly/go-fastly/pull/612)
 - feat(kv_store): Adds support for all documented KV Store API features. [#630](https://github.com/fastly/go-fastly/pull/630)

### Bug fixes:

 - fix(automation_tokens): Fix decodeBodyMap for string to time.Time [#619](https://github.com/fastly/go-fastly/pull/619)

### Dependencies:

- build(deps): `github.com/google/go-cmp` from 0.6.0 to 0.7.0 ([#617](https://github.com/fastly/go-fastly/pull/617))
- build(deps): upgrade Go from 1.22 to 1.23 ([#624](https://github.com/fastly/go-fastly/pull/624/files))
- build(deps): `honnef.co/go/tools` from 0.5.1 to 0.6.0 ([#610](https://github.com/fastly/go-fastly/pull/610))
- build(deps): `golang.org/x/crypto` from 0.33.0 to 0.35.0 ([#618](https://github.com/fastly/go-fastly/pull/618))

## [v9.13.1](https://github.com/fastly/go-fastly/releases/tag/v9.13.1) (2025-02-14)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v9.13.0...v9.13.1)

**Bug fixes:**

- fix(computeacls): lookup call to treat status 204 properly [#605](https://github.com/fastly/go-fastly/pull/605)
- fix(fastly): ensure that HTTP response body objects are always closed [#592](https://github.com/fastly/go-fastly/pull/592)

**Dependencies:**

- build(deps): upgrade Go from 1.20 to 1.22 [#606](https://github.com/fastly/go-fastly/pull/606)
- build(deps): bump golang.org/x/crypto from 0.32.0 to 0.33.0 [#601](https://github.com/fastly/go-fastly/pull/601)
- build(deps): bump golang.org/x/tools from 0.29.0 to 0.30.0 [#598](https://github.com/fastly/go-fastly/pull/598)
- build(deps): bump golang.org/x/sys from 0.29.0 to 0.30.0 [#597](https://github.com/fastly/go-fastly/pull/597)

## [v9.13.0](https://github.com/fastly/go-fastly/releases/tag/v9.13.0) (2025-01-27)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v9.12.0...v9.13.0)

**Enhancements:**

- feat(origin_inspector): Add `limit` query parameter [#568](https://github.com/fastly/go-fastly/pull/568)
- feat(products): Product-specific enablement and configuration [#570](https://github.com/fastly/go-fastly/pull/570)
- feat(domains_v1): Add support for new UDM endpoints [#577](https://github.com/fastly/go-fastly/pull/577)
- feat(computeacls): Add support for compute platform ACLs [#574](https://github.com/fastly/go-fastly/pull/574)

**Bug fixes:**

- fix(domains): Parse error response correctly [#579](https://github.com/fastly/go-fastly/pull/579)
- fix(products): Improve API usability [#572](https://github.com/fastly/go-fastly/pull/572)

**Dependencies:**

- build(deps): bump github.com/google/go-cmp from 0.5.8 to 0.6.0 [#580](https://github.com/fastly/go-fastly/pull/580)
- build(deps): bump github.com/mitchellh/mapstructure from 1.4.3 to 1.5.0 [#580](https://github.com/fastly/go-fastly/pull/580)
- build(deps): bump github.com/peterhellberg/link from 1.1.0 to 1.2.0 [#580](https://github.com/fastly/go-fastly/pull/580)
- build(deps): bump golang.org/x/crypto from 0.31.0 to 0.32.0 [#580](https://github.com/fastly/go-fastly/pull/580)
- build(deps): bump honnef.co/go/tools from 0.4.7 to 0.5.1 [#586](https://github.com/fastly/go-fastly/pull/586)
- build(deps): bump golang.org/x/tools from 0.15.0 to 0.29.0 [#587](https://github.com/fastly/go-fastly/pull/587)

## [v9.12.0](https://github.com/fastly/go-fastly/releases/tag/v9.12.0) (2024-11-07)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v9.11.0...v9.12.0)

**Breaking:**

Note that in spite of this *breaking* change, the major version number
of the package was not incremented, as the feature which was removed
did not work and no users of the package could have been relying on
it.

- breaking(product_enablement): Remove support for NGWAF product. [#560](https://github.com/fastly/go-fastly/pull/560)

**Enhancements:**

- feat(stats_historical): add origin statistics [#552](https://github.com/fastly/go-fastly/pull/552)
- feat(stats_historical): add fields with all prefix [#553](https://github.com/fastly/go-fastly/pull/553)
- Add GrafanaCloudLogs as an logging enpoint [#556](https://github.com/fastly/go-fastly/pull/556)
- feat(product_enablement): Add support for Log Explorer & Insights product. [#558](https://github.com/fastly/go-fastly/pull/558)

**Bug fixes:**

- fix(logging_grafanacloudlogs): Fix Grafana Cloud Logs errors [#559](https://github.com/fastly/go-fastly/pull/559)
- fix(debug_mode): Fix FASTLY_DEBUG_MODE when used in combination with go-vcr. [#561](https://github.com/fastly/go-fastly/pull/561)
- test(infrastructure): Add support for testing against both Delivery and Compute services. [#562](https://github.com/fastly/go-fastly/pull/562)
- test(product_enablement): Add test suites for all supported products. [#563](https://github.com/fastly/go-fastly/pull/563)

**Dependencies:**

- build(deps): Unpin staticcheck, 'latest' version is acceptable. [#549](https://github.com/fastly/go-fastly/pull/549)

## [v9.11.0](https://github.com/fastly/go-fastly/releases/tag/v9.11.0) (2024-10-01)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v9.10.0...v9.11.0)

**Enhancements:**

- feat(custom_dashboards): Add support for Custom Dashboards APIs. [#546](https://github.com/fastly/go-fastly/pull/546)
- feat(automation_tokens): Add support for automation-tokens [#547](https://github.com/fastly/go-fastly/pull/547)
- feat(product_enablement): Add support for NGWAF. [#550](https://github.com/fastly/go-fastly/pull/550)

## [v9.10.0](https://github.com/fastly/go-fastly/releases/tag/v9.10.0) (2024-09-09)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v9.9.0...v9.10.0)

**Enhancements:**

- feat(service, version): Add support for 'environments'. [#542](https://github.com/fastly/go-fastly/pull/542)

**Bug fixes:**

- fix: Ensure that all API endpoint URLs are constructed safely. [#544](https://github.com/fastly/go-fastly/pull/544)

## [v9.9.0](https://github.com/fastly/go-fastly/releases/tag/v9.9.0) (2024-08-20)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v9.8.0...v9.9.0)

**Enhancements:**

- feat(product_enablement): Add support for the 'bot_management' product. [#539](https://github.com/fastly/go-fastly/pull/539)

## [v9.8.0](https://github.com/fastly/go-fastly/releases/tag/v9.8.0) (2024-08-07)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v9.7.0...v9.8.0)

**Enhancements:**

- Add TCP Keep-Alive parameters to backend structs [#537](https://github.com/fastly/go-fastly/pull/537)

## [v9.7.0](https://github.com/fastly/go-fastly/releases/tag/v9.7.0) (2024-06-21)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v9.6.0...v9.7.0)

**Enhancements:**

- feat(tls_custom): add `mutual_authentication` field to response [#535](https://github.com/fastly/go-fastly/pull/535)

- feat(tls/custom_certs): add `in_use` filter [#534](https://github.com/fastly/go-fastly/pull/534)

- feat(client): prepend custom user-agent via `FASTLY_USER_AGENT` env variable to user [#531](https://github.com/fastly/go-fastly/pull/531)

## [v9.6.0](https://github.com/fastly/go-fastly/releases/tag/v9.6.0) (2024-06-04)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v9.5.0...v9.6.0)

**Enhancements:**

- add `limit_workspaces` to `User` model [#529](https://github.com/fastly/go-fastly/pull/529)

## [v9.5.0](https://github.com/fastly/go-fastly/releases/tag/v9.5.0) (2024-05-29)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v9.4.0...v9.5.0)

**Enhancements:**

- feat(package): add ClonedFrom field to PackageMetadata [#527](https://github.com/fastly/go-fastly/pull/527)

## [v9.4.0](https://github.com/fastly/go-fastly/releases/tag/v9.4.0) (2024-05-15)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v9.3.2...v9.4.0)

**Enhancements:**

- feat(helpers): Extend `ToPointer` to cover `[]string` [#525](https://github.com/fastly/go-fastly/pull/525)
- feat(image_opto): Add Image Optimizer default settings API [#516](https://github.com/fastly/go-fastly/pull/516)

## [v9.3.2](https://github.com/fastly/go-fastly/releases/tag/v9.3.2) (2024-05-07)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v9.3.1...v9.3.2)

**Enhancements:**

- feat(tls_subscription): add `tls_certificate` attributes [#523](https://github.com/fastly/go-fastly/pull/523)
- tests(alerts): add percent alerts tests [#522](https://github.com/fastly/go-fastly/pull/522)

## [v9.3.1](https://github.com/fastly/go-fastly/releases/tag/v9.3.1) (2024-04-17)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v9.3.0...v9.3.1)

**Bug fixes:**

- fix(client): support tcp keepalive [#519](https://github.com/fastly/go-fastly/pull/519)

## [v9.3.0](https://github.com/fastly/go-fastly/releases/tag/v9.3.0) (2024-04-16)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v9.2.2...v9.3.0)

**Enhancements:**

- Add support for specifying location when creating KV stores [#517](https://github.com/fastly/go-fastly/pull/517)
- feat: notifications [#513](https://github.com/fastly/go-fastly/pull/513)

## [v9.2.2](https://github.com/fastly/go-fastly/releases/tag/v9.2.2) (2024-04-10)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v9.2.1...v9.2.2)

**Bug fixes:**

- fix(tls_mutual_authentication): support null for ID value [#514](https://github.com/fastly/go-fastly/pull/514)

## [v9.2.1](https://github.com/fastly/go-fastly/releases/tag/v9.2.1) (2024-04-09)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v9.2.0...v9.2.1)

**Bug fixes:**

- fix(tls_custom_activation): only error if no fields provided [#511](https://github.com/fastly/go-fastly/pull/511)

## [v9.2.0](https://github.com/fastly/go-fastly/releases/tag/v9.2.0) (2024-03-15)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v9.1.0...v9.2.0)

**Enhancements:**

- feat: expose project_id field of scalyr logging config [#508](https://github.com/fastly/go-fastly/pull/508)

## [v9.1.0](https://github.com/fastly/go-fastly/releases/tag/v9.1.0) (2024-03-13)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v9.0.1...v9.1.0)

**Enhancements:**

- feat(kv-store-entry): allow parallelization of batch put item [#502](https://github.com/fastly/go-fastly/pull/502)
- feat(historical-stats): add aggregate endpoint [#506](https://github.com/fastly/go-fastly/pull/506)

## [v9.0.1](https://github.com/fastly/go-fastly/releases/tag/v9.0.1) (2024-03-12)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v9.0.0...v9.0.1)

**Bug fixes:**

- fix(config_store): decoding [#503](https://github.com/fastly/go-fastly/pull/503)

**Enhancements:**

- style: resolve gofumpt issues [#500](https://github.com/fastly/go-fastly/pull/500)

**Dependencies:**

- build: bump minimum go version [#504](https://github.com/fastly/go-fastly/pull/504)

## [v9.0.0](https://github.com/fastly/go-fastly/releases/tag/v9.0.0) (2024-02-05)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v8.6.4...v9.0.0)

**Breaking**:

To summarize the key items in the v9 release:

- The paginators for Services, ACL Entries and Dictionary Items have been replaced with a generic paginator.
- The KV Stores Entries paginator was NOT replaced as it has a incompatible implementation.
- The `List*` methods no longer sort results, and will now return _all_ results.
- All endpoints (except those using JSONAPI) use pointers consistently for Create/Update methods.

The complete list of breaking changes are:

- fix(token): add ListTokensInput to ListTokens [#487](https://github.com/fastly/go-fastly/pull/487)
- feat(config store): adds `ListConfigStoresInput` when returning `ListConfigStores` [#481](https://github.com/fastly/go-fastly/pull/481)
- Use integer instead of string [#486](https://github.com/fastly/go-fastly/pull/486)
- fix(request_settings): make action a pointer in update [#488](https://github.com/fastly/go-fastly/pull/488)
- feat: generic paginator [#491](https://github.com/fastly/go-fastly/pull/491)
- fix: all relevant fields to be pointers [#493](https://github.com/fastly/go-fastly/pull/493)

**Enhancements:**

- feat: domain inspector [#483](https://github.com/fastly/go-fastly/pull/483)
- Move CBool helper with the other helpers [#484](https://github.com/fastly/go-fastly/pull/484)
- Support retrieving a secret store by name [#485](https://github.com/fastly/go-fastly/pull/486)
- replace: pointer helpers to generic ToPointer function [#489](https://github.com/fastly/go-fastly/pull/489)
- refactor(helpers): avoid explicit types in favour of tilda type constraint [#490](https://github.com/fastly/go-fastly/pull/490)
- refactor(secret_store): replace json.NewEncoder with c.PostJSON [#492](https://github.com/fastly/go-fastly/pull/492)
- refactor: rename files [#494](https://github.com/fastly/go-fastly/pull/494)
- refactor: resolves linter issues [#495](https://github.com/fastly/go-fastly/pull/495)
- refactor: use consistent naming conventions [#496](https://github.com/fastly/go-fastly/pull/496)
- feat: fastly alerts [#499](https://github.com/fastly/go-fastly/pull/499)
- refactor: use explicit naming for IDs [#497](https://github.com/fastly/go-fastly/pull/497)

**Dependencies:**

- build(deps): bump golang.org/x/crypto from 0.1.0 to 0.17.0 [#498](https://github.com/fastly/go-fastly/pull/498)

## [v8.6.4](https://github.com/fastly/go-fastly/releases/tag/v8.6.4) (2023-10-31)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v8.6.3...v8.6.4)

**Enhancements:**

- feat: add kv-store list keys `consistency` param [#479](https://github.com/fastly/go-fastly/pull/479)

## [v8.6.3](https://github.com/fastly/go-fastly/releases/tag/v8.6.3) (2023-10-31)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v8.6.2...v8.6.3)

**Enhancements:**

- feat(client): expose DebugMode field [#477](https://github.com/fastly/go-fastly/pull/477)

**Documentation:**

- docs: support for Compute [#476](https://github.com/fastly/go-fastly/pull/476)

## [v8.6.2](https://github.com/fastly/go-fastly/releases/tag/v8.6.2) (2023-09-28)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v8.6.1...v8.6.2)

**Bug fixes:**

- fix(client): simplify key-redaction logic [#472](https://github.com/fastly/go-fastly/pull/472)

## [v8.6.1](https://github.com/fastly/go-fastly/releases/tag/v8.6.1) (2023-08-31)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v8.6.0...v8.6.1)

**Bug fixes:**

- fastly: fix return value for errors during kvstore pagination [#468](https://github.com/fastly/go-fastly/pull/468)

**Dependencies:**

- build: bump minimum go version to 1.19 [#467](https://github.com/fastly/go-fastly/pull/467)

## [v8.6.0](https://github.com/fastly/go-fastly/releases/tag/v8.6.0) (2023-08-30)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v8.5.9...v8.6.0)

**Enhancements:**

- feat(backend): add share_key field [#463](https://github.com/fastly/go-fastly/pull/463)

## [v8.5.9](https://github.com/fastly/go-fastly/releases/tag/v8.5.9) (2023-08-09)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v8.5.8...v8.5.9)

**Bug fixes:**

- revert(backend): revert removal of redundant ErrorThreshold field [#464](https://github.com/fastly/go-fastly/pull/464)

## [v8.5.8](https://github.com/fastly/go-fastly/releases/tag/v8.5.8) (2023-08-04)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v8.5.7...v8.5.8)

**Enhancements:**

- Add support for the New Relic OTLP logging endpoint [#460](https://github.com/fastly/go-fastly/pull/460)

**Bug fixes:**

- remove(backend): remove redundant ErrorThreshold field [#461](https://github.com/fastly/go-fastly/pull/461)

## [v8.5.7](https://github.com/fastly/go-fastly/releases/tag/v8.5.7) (2023-07-24)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v8.5.6...v8.5.7)

**Enhancements:**

- feat(kv_store): add interface for easier testing of paginator [#458](https://github.com/fastly/go-fastly/pull/458)

## [v8.5.6](https://github.com/fastly/go-fastly/releases/tag/v8.5.6) (2023-07-24)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v8.5.5...v8.5.6)

**Enhancements:**

- fix(kv_store): support parallel calls to delete keys [#456](https://github.com/fastly/go-fastly/pull/456)

## [v8.5.5](https://github.com/fastly/go-fastly/releases/tag/v8.5.5) (2023-07-21)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v8.5.4...v8.5.5)

**Enhancements:**

- fix(errors): support KV Store batch endpoint error response [#454](https://github.com/fastly/go-fastly/pull/454)

## [v8.5.4](https://github.com/fastly/go-fastly/releases/tag/v8.5.4) (2023-06-29)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v8.5.3...v8.5.4)

**Bug fixes:**

- fix: ensure Fastly-Key is stripped from the dump [#453](https://github.com/fastly/go-fastly/pull/453)

## [v8.5.3](https://github.com/fastly/go-fastly/releases/tag/v8.5.3) (2023-06-29)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v8.5.2...v8.5.3)

**Enhancements:**

- feat: add debug mode [#451](https://github.com/fastly/go-fastly/pull/451)

## [v8.5.2](https://github.com/fastly/go-fastly/releases/tag/v8.5.2) (2023-06-23)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v8.5.1...v8.5.2)

**Bug fixes:**

- fix(errors): use title if present [#449](https://github.com/fastly/go-fastly/pull/449)

## [v8.5.1](https://github.com/fastly/go-fastly/releases/tag/v8.5.1) (2023-06-14)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v8.5.0...v8.5.1)

**Enhancements:**

- Add support for file_max_bytes field to s3 logging endpoint [#446](https://github.com/fastly/go-fastly/pull/446)
- feat(errors): expose rate limit headers [#447](https://github.com/fastly/go-fastly/pull/447)

## [v8.5.0](https://github.com/fastly/go-fastly/releases/tag/v8.5.0) (2023-06-07)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v8.4.1...v8.5.0)

**Enhancements:**

- feat(config_store): implement batch endpoint [#444](https://github.com/fastly/go-fastly/pull/444)

## [v8.4.1](https://github.com/fastly/go-fastly/releases/tag/v8.4.1) (2023-06-01)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v8.4.0...v8.4.1)

**Enhancements:**

- feat(mutual_authentication): add Enforced to CREATE [#442](https://github.com/fastly/go-fastly/pull/442)

## [v8.4.0](https://github.com/fastly/go-fastly/releases/tag/v8.4.0) (2023-05-31)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v8.3.0...v8.4.0)

**Enhancements:**

- Support Mutual Authentication (mTLS) endpoints [#440](https://github.com/fastly/go-fastly/pull/440)

## [v8.3.0](https://github.com/fastly/go-fastly/releases/tag/v8.3.0) (2023-05-12)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v8.2.1...v8.3.0)

**Enhancements:**

- fix(kv_store): allow buffering io for key insert [#437](https://github.com/fastly/go-fastly/pull/437)

## [v8.2.1](https://github.com/fastly/go-fastly/releases/tag/v8.2.1) (2023-05-11)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v8.2.0...v8.2.1)

**Enhancements:**

- feat(package): FilesHash [#435](https://github.com/fastly/go-fastly/pull/435)

## [v8.2.0](https://github.com/fastly/go-fastly/releases/tag/v8.2.0) (2023-05-11)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v8.1.0...v8.2.0)

**Enhancements:**

- Secret Store: Support PUT & PATCH methods when creating secret [#433](https://github.com/fastly/go-fastly/pull/433)

## [v8.1.0](https://github.com/fastly/go-fastly/releases/tag/v8.1.0) (2023-05-10)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v8.0.3...v8.1.0)

**Enhancements:**

- feat: kv_store batch API endpoint [#431](https://github.com/fastly/go-fastly/pull/431)

## [v8.0.3](https://github.com/fastly/go-fastly/releases/tag/v8.0.3) (2023-05-09)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v8.0.2...v8.0.3)

**Enhancements:**

- fix(kv_store): support parallel requests [#429](https://github.com/fastly/go-fastly/pull/429)

## [v8.0.2](https://github.com/fastly/go-fastly/releases/tag/v8.0.2) (2023-05-09)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v8.0.1...v8.0.2)

**Bug fixes:**

- fix(kv_store): allow file read support [#427](https://github.com/fastly/go-fastly/pull/427)

## [v8.0.1](https://github.com/fastly/go-fastly/releases/tag/v8.0.1) (2023-04-26)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v8.0.0...v8.0.1)

**Enhancements:**

- Add project id for GCP in logging [#425](https://github.com/fastly/go-fastly/pull/425)

## [v8.0.0](https://github.com/fastly/go-fastly/releases/tag/v8.0.0) (2023-04-12)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v7.5.5...v8.0.0)

**Breaking**:

- breaking(secret-store): rename `ValidateSignature` to `VerifySignature` [#397](https://github.com/fastly/go-fastly/pull/397)
- breaking(object-store): rename `ObjectStore` to `KVStore` [#422](https://github.com/fastly/go-fastly/pull/422)

## [v7.5.5](https://github.com/fastly/go-fastly/releases/tag/v7.5.5) (2023-03-31)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v7.5.4...v7.5.5)

**Enhancements:**

- feat(ratelimiter): add uri_dictionary_name property [#420](https://github.com/fastly/go-fastly/pull/420)

## [v7.5.4](https://github.com/fastly/go-fastly/releases/tag/v7.5.4) (2023-03-31)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v7.5.3...v7.5.4)

**Bug fixes:**

- fix(ratelimiter): use separate structs for input/response serialization [#418](https://github.com/fastly/go-fastly/pull/418)

## [v7.5.3](https://github.com/fastly/go-fastly/releases/tag/v7.5.3) (2023-03-31)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v7.5.2...v7.5.3)

**Bug fixes:**

- fix(ratelimiter): serialize with mapstructure tag [#416](https://github.com/fastly/go-fastly/pull/416)

## [v7.5.2](https://github.com/fastly/go-fastly/releases/tag/v7.5.2) (2023-03-30)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v7.5.1...v7.5.2)

**Enhancements:**

- feat(erl): add missing fields [#414](https://github.com/fastly/go-fastly/pull/414)

## [v7.5.1](https://github.com/fastly/go-fastly/releases/tag/v7.5.1) (2023-03-28)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v7.5.0...v7.5.1)

**Bug fixes:**

- fix(lint): avoid range variable being captured by func literal [#411](https://github.com/fastly/go-fastly/pull/411)
- fix(ratelimit): add missing LoggerType input for CREATE/UPDATE [#412](https://github.com/fastly/go-fastly/pull/412)

## [v7.5.0](https://github.com/fastly/go-fastly/releases/tag/v7.5.0) (2023-03-15)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v7.4.0...v7.5.0)

**Enhancements:**

- feat: Add Config Store support [#398](https://github.com/fastly/go-fastly/pull/398)

**Bug fixes:**

- fix(tls-subscription): serialise warnings field [#409](https://github.com/fastly/go-fastly/pull/409)

## [v7.4.0](https://github.com/fastly/go-fastly/releases/tag/v7.4.0) (2023-03-07)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v7.3.0...v7.4.0)

**Enhancements:**

- fix(lint): semgrep [#406](https://github.com/fastly/go-fastly/pull/406)
- Expose `KeepAliveTime` field from Backend API [#405](https://github.com/fastly/go-fastly/pull/405)

**Dependencies:**

- Bump golang.org/x/crypto from 0.0.0-20210921155107-089bfa567519 to 0.1.0 [#404](https://github.com/fastly/go-fastly/pull/404)

## [v7.3.0](https://github.com/fastly/go-fastly/releases/tag/v7.3.0) (2023-02-17)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v7.2.0...v7.3.0)

**Enhancements:**

- feat(product_enablement): implement API [#399](https://github.com/fastly/go-fastly/pull/399)
- build(make): setup semgrep configuration [#400](https://github.com/fastly/go-fastly/pull/400)

## [v7.2.0](https://github.com/fastly/go-fastly/releases/tag/v7.2.0) (2023-02-08)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v7.1.0...v7.2.0)

**Enhancements:**

- allow Go binary to be specified when running make [#393](https://github.com/fastly/go-fastly/pull/393)
- add support for Secret Store client keys [#392](https://github.com/fastly/go-fastly/pull/392)
- resource: improve parameter naming; support JSON encode [#394](https://github.com/fastly/go-fastly/pull/394)
- feat(http3): implement HTTP3 API endpoints [#395](https://github.com/fastly/go-fastly/pull/395)

## [v7.1.0](https://github.com/fastly/go-fastly/releases/tag/v7.1.0) (2023-01-19)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v7.0.0...v7.1.0)

**Enhancements:**

- fastly: allow uploading compute package from memory [#384](https://github.com/fastly/go-fastly/pull/384)
- Improve error response handling [#386](https://github.com/fastly/go-fastly/pull/386)
- secret-store: Add created_at fields to store and secret [#387](https://github.com/fastly/go-fastly/pull/387)
- Add compute_requests to usage response struct [#388](https://github.com/fastly/go-fastly/pull/388)
- feat(resource): implement Resource API [#390](https://github.com/fastly/go-fastly/pull/390)

**Bug fixes:**

- fastly: fix invalid package uploading test [#385](https://github.com/fastly/go-fastly/pull/385)
- Fix typo in AttackRequestHeaderBytes field name [#389](https://github.com/fastly/go-fastly/pull/389)

## [v7.0.0](https://github.com/fastly/go-fastly/releases/tag/v7.0.0) (2022-11-07)

The biggest change in this release is the majority of input fields are now pointers rather than primitives. This is to enable users to set an empty value for fields when they are being sent to the Fastly API, in cases where the user wishes to override an otherwise default value provided by the API.

There has been a bunch of interface fixes (e.g. consistent use of `int` over `uint` equivalents, typos in field names fixed, response types renamed to be less ambiguous etc) and also many other changes to better support consistency and documentation across the code base.

[Full Changelog](https://github.com/fastly/go-fastly/compare/v6.8.0...v7.0.0)

**Breaking**:

- Backend resource tweaks to allow setting zero value [#375](https://github.com/fastly/go-fastly/pull/375)
- Sort all struct fields [#376](https://github.com/fastly/go-fastly/pull/376)
- Address FIXME notes [#379](https://github.com/fastly/go-fastly/pull/379)
- Revive fixes [#378](https://github.com/fastly/go-fastly/pull/378)
- Ensure all relevant 'create' input fields are using pointers [#382](https://github.com/fastly/go-fastly/pull/382)

**Bug fixes:**

- Ensure parameters are sent to API [#380](https://github.com/fastly/go-fastly/pull/380)

**Enhancements:**

- Add google account name to all google logging endpoints [#377](https://github.com/fastly/go-fastly/pull/377)

## [v6.8.0](https://github.com/fastly/go-fastly/releases/tag/v6.8.0) (2022-10-10)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v6.7.0...v6.8.0)

**Enhancements:**

- Support Health Check Headers [#373](https://github.com/fastly/go-fastly/pull/373)

## [v6.7.0](https://github.com/fastly/go-fastly/releases/tag/v6.7.0) (2022-10-06)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v6.6.0...v6.7.0)

**Enhancements:**

- Add Secret Store API support [#367](https://github.com/fastly/go-fastly/pull/367)

## [v6.6.0](https://github.com/fastly/go-fastly/releases/tag/v6.6.0) (2022-10-05)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v6.5.2...v6.6.0)

**Enhancements:**

- Make gzip level consistent everywhere [#368](https://github.com/fastly/go-fastly/pull/368)
- Add Object Store API support [#359](https://github.com/fastly/go-fastly/pull/359)

## [v6.5.2](https://github.com/fastly/go-fastly/releases/tag/v6.5.2) (2022-09-08)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v6.5.1...v6.5.2)

**Bug fixes:**

- Avoid duplicate `Close()` call on `http.Response.Body` [#365](https://github.com/fastly/go-fastly/pull/365)

## [v6.5.1](https://github.com/fastly/go-fastly/releases/tag/v6.5.1) (2022-09-07)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v6.5.0...v6.5.1)

**Bug fixes:**

- Ensure `http.Response.Body` is closed [#363](https://github.com/fastly/go-fastly/pull/363)

## [v6.5.0](https://github.com/fastly/go-fastly/releases/tag/v6.5.0) (2022-08-22)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v6.4.3...v6.5.0)

**Enhancements:**

- TLS Subscriptions: add filter for active orders [#357](https://github.com/fastly/go-fastly/pull/357)

## [v6.4.3](https://github.com/fastly/go-fastly/releases/tag/v6.4.3) (2022-07-20)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v6.4.2...v6.4.3)

**Bug fixes:**

- `Token` field is required for some log endpoints [#355](https://github.com/fastly/go-fastly/pull/355)

## [v6.4.2](https://github.com/fastly/go-fastly/releases/tag/v6.4.2) (2022-07-11)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v6.4.1...v6.4.2)

**Bug fixes:**

- purge: fail to purge with query string [#353](https://github.com/fastly/go-fastly/pull/353)

## [v6.4.1](https://github.com/fastly/go-fastly/releases/tag/v6.4.1) (2022-06-23)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v6.4.0...v6.4.1)

**Enhancements:**

- Implement missing `ListServiceAuthorizations` [#351](https://github.com/fastly/go-fastly/pull/351)

## [v6.4.0](https://github.com/fastly/go-fastly/releases/tag/v6.4.0) (2022-06-20)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v6.3.2...v6.4.0)

**Enhancements:**

- Support for Service Authorizations [#349](https://github.com/fastly/go-fastly/pull/349)
- logging/gcp: add AccountName [#346](https://github.com/fastly/go-fastly/pull/346)

## [v6.3.2](https://github.com/fastly/go-fastly/releases/tag/v6.3.2) (2022-04-27)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v6.3.1...v6.3.2)

**Enhancements:**

- Add `Backends` field to `Director` [#347](https://github.com/fastly/go-fastly/pull/347)

## [v6.3.1](https://github.com/fastly/go-fastly/releases/tag/v6.3.1) (2022-04-04)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v6.3.0...v6.3.1)

**Bug fixes:**

- Fix typo in `filter[tls_certificate.id]` [#344](https://github.com/fastly/go-fastly/pull/344)

## [v6.3.0](https://github.com/fastly/go-fastly/releases/tag/v6.3.0) (2022-03-25)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v6.2.0...v6.3.0)

**Enhancements:**

- Support for Edge Rate Limiting [#341](https://github.com/fastly/go-fastly/pull/341)

**Bug fixes:**

- go.mod: downgrade to 1.16 [#340](https://github.com/fastly/go-fastly/pull/340)

## [v6.2.0](https://github.com/fastly/go-fastly/releases/tag/v6.2.0) (2022-02-23)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v6.1.0...v6.2.0)

**Enhancements:**

- Origin Inspector: Get Historical Origin Metrics for a Service [#335](https://github.com/fastly/go-fastly/pull/335)

## [v6.1.0](https://github.com/fastly/go-fastly/releases/tag/v6.1.0) (2022-02-22)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v6.0.1...v6.1.0)

**Enhancements:**

- Support `Fastly-RateLimit-Remaining` and `Fastly-RateLimit-Reset` [#337](https://github.com/fastly/go-fastly/pull/337)
- Stats: Add test for fetching stats by service and field [#333](https://github.com/fastly/go-fastly/pull/333)

**Documentation:**

- Fix README typo [#334](https://github.com/fastly/go-fastly/pull/334)

## [v6.0.1](https://github.com/fastly/go-fastly/releases/tag/v6.0.1) (2022-01-17)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v6.0.0...v6.0.1)

**Bug fixes:**

- Correctly set and advance current page for pagination types [#330](https://github.com/fastly/go-fastly/pull/330)

## [v6.0.0](https://github.com/fastly/go-fastly/releases/tag/v6.0.0) (2022-01-13)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v5.3.0...v6.0.0)

**Enhancements:**

- Pagination to return interface not concrete type [#328](https://github.com/fastly/go-fastly/pull/328)
- Update go modules and fixtures [#325](https://github.com/fastly/go-fastly/pull/325)

## [v5.3.0](https://github.com/fastly/go-fastly/releases/tag/v5.3.0) (2022-01-10)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v5.2.0...v5.3.0)

**Enhancements:**

- Add pagination support for `GET /service` endpoint [#326](https://github.com/fastly/go-fastly/pull/326)

## [v5.2.0](https://github.com/fastly/go-fastly/releases/tag/v5.2.0) (2022-01-07)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v5.1.3...v5.2.0)

**Enhancements:**

- Bump deps to side-step `gopkg.in/yaml.v2` security vulnerability [#323](https://github.com/fastly/go-fastly/pull/323)
- Add support for `modsec_rule_id` filter parameter [#322](https://github.com/fastly/go-fastly/pull/322)
- Add support for Link header pagination [#321](https://github.com/fastly/go-fastly/pull/321)

## [v5.1.3](https://github.com/fastly/go-fastly/releases/tag/v5.1.3) (2021-12-03)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v5.1.2...v5.1.3)

**Bug fixes:**

- Do not omit zero int values on create operations [#318](https://github.com/fastly/go-fastly/pull/318)
- Fix `len < 0` in test: len cannot return negative values [#319](https://github.com/fastly/go-fastly/pull/319)

## [v5.1.2](https://github.com/fastly/go-fastly/releases/tag/v5.1.2) (2021-11-04)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v5.1.1...v5.1.2)

**Bug fixes:**

- Use pointer for `subnet` field [#316](https://github.com/fastly/go-fastly/pull/316)
- Replace form dependency with google's go-querystring [#315](https://github.com/fastly/go-fastly/pull/315)

## [v5.1.1](https://github.com/fastly/go-fastly/releases/tag/v5.1.1) (2021-10-11)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v5.1.0...v5.1.1)

**Bug fixes:**

- Manually fill `ActiveVersion` value in `GetService` [#312](https://github.com/fastly/go-fastly/pull/312)

## [v5.1.0](https://github.com/fastly/go-fastly/releases/tag/v5.1.0) (2021-10-05)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v5.0.0...v5.1.0)

**Enhancements:**

- Add `UseTLS` field for Splunk logging [#309](https://github.com/fastly/go-fastly/pull/309)

## [v5.0.0](https://github.com/fastly/go-fastly/releases/tag/v5.0.0) (2021-09-23)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v4.0.0...v5.0.0)

**Breaking**:

- Switch from `bool` to `Compatibool` for ACL entries [#307](https://github.com/fastly/go-fastly/pull/307)

## [v4.0.0](https://github.com/fastly/go-fastly/releases/tag/v4.0.0) (2021-09-21)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v3.12.0...v4.0.0)

**Breaking**:

- Change Backend `SSLCiphers` field from `[]string` to `string` [#304](https://github.com/fastly/go-fastly/pull/304)

**Bug fixes:**

- Fix issues with serialising x-www-form-urlencoded data [#304](https://github.com/fastly/go-fastly/pull/304)

## [v3.12.0](https://github.com/fastly/go-fastly/releases/tag/v3.12.0) (2021-09-15)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v3.11.0...v3.12.0)

**Enhancements:**

- Implement `DELETE /tokens` (bulk token deletion) API [#302](https://github.com/fastly/go-fastly/pull/302)

## [v3.11.0](https://github.com/fastly/go-fastly/releases/tag/v3.11.0) (2021-09-15)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v3.10.0...v3.11.0)

**Enhancements:**

- Implement user password reset [#300](https://github.com/fastly/go-fastly/pull/300)

## [v3.10.0](https://github.com/fastly/go-fastly/releases/tag/v3.10.0) (2021-09-13)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v3.9.5...v3.10.0)

**Enhancements:**

- Add domain validation API [#298](https://github.com/fastly/go-fastly/pull/298)

## [v3.9.5](https://github.com/fastly/go-fastly/releases/tag/v3.9.5) (2021-09-08)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v3.9.4...v3.9.5)

**Bug fixes:**

- Allow creating dynamic VCL snippets with empty content [#296](https://github.com/fastly/go-fastly/pull/296)

## [v3.9.4](https://github.com/fastly/go-fastly/releases/tag/v3.9.4) (2021-09-07)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v3.9.3...v3.9.4)

**Enhancements:**

- Add region support for NewRelic logging endpoint [#292](https://github.com/fastly/go-fastly/pull/292)

**Bug fixes:**

- Remove `omitempty` from `SSLCheckCert` field [#294](https://github.com/fastly/go-fastly/pull/294)

**Documentation:**

- Correct indentation in README code example [#293](https://github.com/fastly/go-fastly/pull/293)

## [v3.9.3](https://github.com/fastly/go-fastly/releases/tag/v3.9.3) (2021-07-23)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v3.9.2...v3.9.3)

**Bug fixes:**

- ACL entry `Subnet` field should be integer type (missed references) [#290](https://github.com/fastly/go-fastly/pull/290)

## [v3.9.2](https://github.com/fastly/go-fastly/releases/tag/v3.9.2) (2021-07-23)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v3.9.1...v3.9.2)

**Bug fixes:**

- ACL entry `Subnet` field should be integer type [#288](https://github.com/fastly/go-fastly/pull/288)

## [v3.9.1](https://github.com/fastly/go-fastly/releases/tag/v3.9.1) (2021-07-21)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v3.9.0...v3.9.1)

**Bug fixes:**

- `omitempty` should be set for ContentTypes and Extensions [#282](https://github.com/fastly/go-fastly/pull/282)
- Check pointer struct type isn't nil before referencing its fields [#286](https://github.com/fastly/go-fastly/pull/286)

## [v3.9.0](https://github.com/fastly/go-fastly/releases/tag/v3.9.0) (2021-06-29)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v3.8.0...v3.9.0)

**Enhancements:**

- Batch Surrogate-Key purge [#284](https://github.com/fastly/go-fastly/pull/284)

## [v3.8.0](https://github.com/fastly/go-fastly/releases/tag/v3.8.0) (2021-06-28)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v3.7.0...v3.8.0)

**Enhancements:**

- Update VCL snippets to more closely align with API [#281](https://github.com/fastly/go-fastly/pull/281)

## [v3.7.0](https://github.com/fastly/go-fastly/releases/tag/v3.7.0) (2021-06-10)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v3.6.0...v3.7.0)

**Enhancements:**

- Add S3 Logging ACL parameter and add Redundancy options [#279](https://github.com/fastly/go-fastly/pull/279)
- Add HTTP `206 Partial Content` to stats [#277](https://github.com/fastly/go-fastly/pull/277)

## [v3.6.0](https://github.com/fastly/go-fastly/releases/tag/v3.6.0) (2021-04-15)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v3.5.0...v3.6.0)

**Enhancements:**

- Fix testing behavior for logging endpoints the support compression [#271](https://github.com/fastly/go-fastly/pull/271)

## [v3.5.0](https://github.com/fastly/go-fastly/releases/tag/v3.5.0) (2021-04-06)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v3.4.1...v3.5.0)

**Enhancements:**

- Support usage of IAM role in S3 and Kinesis logging endpoints [#269](https://github.com/fastly/go-fastly/pull/269)

## [v3.4.1](https://github.com/fastly/go-fastly/releases/tag/v3.4.1) (2021-03-25)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v3.4.0...v3.4.1)

**Bug fixes:**

- Purge with Soft=true bug [#266](https://github.com/fastly/go-fastly/issues/266)
- Initialise Headers map to avoid runtime panic when purging. [#267](https://github.com/fastly/go-fastly/pull/267)

**Closed issues:**

- Potentially misleading comment in README.md [#260](https://github.com/fastly/go-fastly/issues/260)

## [v3.4.0](https://github.com/fastly/go-fastly/releases/tag/v3.4.0) (2021-02-18)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v3.3.0...v3.4.0)

**Enhancements:**

- Add PATCH endpoint for TLS Subscriptions [#262](https://github.com/fastly/go-fastly/pull/262)

## [v3.3.0](https://github.com/fastly/go-fastly/releases/tag/v3.3.0) (2021-02-15)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v3.2.0...v3.3.0)

**Enhancements:**

- Updates needed to support Terraform TLS resources [#259](https://github.com/fastly/go-fastly/pull/259)

## [v3.2.0](https://github.com/fastly/go-fastly/releases/tag/v3.2.0) (2021-02-04)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v3.1.0...v3.2.0)

**Enhancements:**

- Add more support for compression_codec to logging endpoints [#257](https://github.com/fastly/go-fastly/pull/257)

## [v3.1.0](https://github.com/fastly/go-fastly/releases/tag/v3.1.0) (2021-01-28)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v3.0.0...v3.1.0)

**Enhancements:**

- Add support for file_max_bytes configuration for azure logging endpoint [#255](https://github.com/fastly/go-fastly/pull/255)

## [v3.0.0](https://github.com/fastly/go-fastly/releases/tag/v3.0.0) (2021-01-19)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v2.1.0...v3.0.0)

There were a few breaking changes introduced in v3:

1. A new `FieldError` abstraction for validating API struct fields.
1. Changing some mandatory fields to Optional (and vice-versa) to better support more _practical_ API usage.
1. Avoid generic ID field when more explicit naming would be clearer.

**Enhancements:**

- Strip TLS prefix from TLS struct fields [#250](https://github.com/fastly/go-fastly/pull/250)
- Avoid generic ID field when more explicit naming would be clearer [#247](https://github.com/fastly/go-fastly/pull/247)
- Update Test Fixtures [#245](https://github.com/fastly/go-fastly/pull/245)
- Add region support to logentries logging endpoint [#243](https://github.com/fastly/go-fastly/pull/243)
- Add basic managed logging endpoint support to go-fastly [#241](https://github.com/fastly/go-fastly/pull/241)
- Create new error abstraction for field validation [#239](https://github.com/fastly/go-fastly/pull/239)

**Bug fixes:**

- NewName should be optional [#252](https://github.com/fastly/go-fastly/pull/252)
- Dictionary ItemValue isn't optional [#251](https://github.com/fastly/go-fastly/pull/251)
- Ensure consistent naming for ServiceID (fixes missed references) [#249](https://github.com/fastly/go-fastly/pull/249)
- Update to RequestMaxBytes to align with updated API and regenerate fixtures [#248](https://github.com/fastly/go-fastly/pull/248)
- Cleanup naming of Kinesis to be more consistent. [#246](https://github.com/fastly/go-fastly/pull/246)
- Reword expected error message based on API changes [#244](https://github.com/fastly/go-fastly/pull/244)

**Closed issues:**

- Remove uninitialized ActiveVersion field from Service struct? [#242](https://github.com/fastly/go-fastly/issues/242)
- Cache setting is missing the 'deliver' action [#136](https://github.com/fastly/go-fastly/issues/136)

## [v2.1.0](https://github.com/fastly/go-fastly/releases/tag/v2.1.0) (2020-12-11)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v2.0.0...v2.1.0)

**Enhancements:**

- Add support for TLS client and batch size options for splunk [#236](https://github.com/fastly/go-fastly/pull/236)
- Add support for compression_codec to logging file sink endpoints [#235](https://github.com/fastly/go-fastly/pull/235)
- Add support for Kinesis logging endpoint [#234](https://github.com/fastly/go-fastly/pull/234)
- Add SASL fields support for Kafka Logging Endpoint [#226](https://github.com/fastly/go-fastly/pull/226)
- Custom TLS API  [#225](https://github.com/fastly/go-fastly/pull/225)

**Closed issues:**

- Any plan to add custom TLS certificates? [#224](https://github.com/fastly/go-fastly/issues/224)

## [v2.0.0](https://github.com/fastly/go-fastly/releases/tag/v2.0.0) (2020-11-17)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v1.18.0...v2.0.0)

The move from major version v1 to v2 has resulted in a couple of fundamental changes to the library:

- Consistent field name format for IDs and Versions (e.g. `DictionaryID`, `PoolID`, `ServiceID`, `ServiceVersion` etc).
- Input struct fields (for write/update operations) that are optional (i.e. `omitempty`) and use basic types, are now defined as pointers.

The move to more consistent field names in some cases will have resulted in the corresponding sentinel error name to be updated also. For example, `ServiceID` has resulted in a change from `ErrMissingService` to `ErrMissingServiceID`.

The change in type for [basic types](https://tour.golang.org/basics/11) that are optional on input structs related to write/update operations is designed to avoid unexpected behaviours when dealing with their zero value (see [this reference](https://willnorris.com/2014/05/go-rest-apis-and-pointers/) for more details). As part of this change we now provide [helper functions](./fastly/basictypes_helper.go) to assist with generating the new pointer types required.

> Note: some read/list operations require fields to be provided but if omitted a zero value will be used when marshaling the data structure into JSON. This too can cause confusion, which is why some input structs define their mandatory fields as pointers (to ensure that the backend can distinguish between a zero value and an omitted field).

**Enhancements:**

- v2 [#230](https://github.com/fastly/go-fastly/pull/230)

**Closed issues:**

- Fails to Parse Historic Stats when no Service Provided [#214](https://github.com/fastly/go-fastly/issues/214)

## [v1.18.0](https://github.com/fastly/go-fastly/releases/tag/v1.18.0) (2020-10-28)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v1.17.0...v1.18.0)

## [v1.17.0](https://github.com/fastly/go-fastly/releases/tag/v1.17.0) (2020-07-20)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v1.16.2...v1.17.0)

**Enhancements:**

- Added support to list all datacenters [#210](https://github.com/fastly/go-fastly/pull/210)

## [v1.16.2](https://github.com/fastly/go-fastly/releases/tag/v1.16.2) (2020-07-13)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v1.16.1...v1.16.2)

**Bug fixes:**

- Allow message_type support for FTP endpoint [#212](https://github.com/fastly/go-fastly/pull/212)

## [v1.16.1](https://github.com/fastly/go-fastly/releases/tag/v1.16.1) (2020-07-07)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v1.16.0...v1.16.1)

**Bug fixes:**

- ci: add GH Action for fmt, vet, staticcheck, test [#184](https://github.com/fastly/go-fastly/pull/184)

## [v1.16.0](https://github.com/fastly/go-fastly/releases/tag/v1.16.0) (2020-06-25)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v1.15.0...v1.16.0)

**Enhancements:**

- Add support for Compute@Edge Packages [#203](https://github.com/fastly/go-fastly/pull/203)

## [v1.15.0](https://github.com/fastly/go-fastly/releases/tag/v1.15.0) (2020-06-04)

[Full Changelog](https://github.com/fastly/go-fastly/compare/v1.14.0...v1.15.0)

**Enhancements:**

- Add PublicKey field to S3 all CRUD actions [#198](https://github.com/fastly/go-fastly/pull/198)
- Add User field to Cloudfiles Updates [#197](https://github.com/fastly/go-fastly/pull/197)
- Remove extraneous Token field from all Kafka CRUD [#196](https://github.com/fastly/go-fastly/pull/196)
- Add Region field to all Scalyr CRUD actions [#195](https://github.com/fastly/go-fastly/pull/195)
- Add MessageType field to all SFTP CRUD actions [#194](https://github.com/fastly/go-fastly/pull/194)
- Add MessageType field to GCS Updates [#193](https://github.com/fastly/go-fastly/pull/193)

# Historical Manual Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v1.14.0] (2020-05-26)

### Added

- ip: Support for accessing Fastly's IPv6 ranges ([#189](https://github.com/fastly/go-fastly/pull/189)).

## [v1.13.0] (2020-05-19)

### Added

- helpers: Add `NullString` helper ([#187](https://github.com/fastly/go-fastly/pull/187)).

## [v1.12.0] (2020-05-15)

### Added

- waf: Support for `link` field ([#179](https://github.com/fastly/go-fastly/pull/179)).

## [v1.11.0] (2020-05-13)

### Added

- test: Updates testing helper for fix fixtures to support macOS ([#177](https://github.com/fastly/go-fastly/pull/177)).
- helpers: Add raw access to stats JSON responses ([#176](https://github.com/fastly/go-fastly/pull/176)).
- datadog: Add Datadog logging endpoint support ([#182](https://github.com/fastly/go-fastly/pull/182)).
- sftp: Add SFTP logging endpoint support ([#175](https://github.com/fastly/go-fastly/pull/175)).
- scalyr: Add Scalyr logging endpoint support ([#174](https://github.com/fastly/go-fastly/pull/174)).
- pubsub: Add Google Cloud Pub/Sub logging endpoint support ([#173](https://github.com/fastly/go-fastly/pull/173)).
- openstack: Add OpenStack logging endpoint support ([#172](https://github.com/fastly/go-fastly/pull/172)).
- newrelic: Add New Relic logging endpoint support ([#171](https://github.com/fastly/go-fastly/pull/171)).
- logshuttle: Add Log Shuttle logging endpoint support ([#170](https://github.com/fastly/go-fastly/pull/170)).
- loggly: Add Loggly logging endpoint support ([#169](https://github.com/fastly/go-fastly/pull/169)).
- kafka: Add Kafka logging endpoint support ([#168](https://github.com/fastly/go-fastly/pull/168)).
- honeycomb: Add Honeycomb logging endpoint support ([#167](https://github.com/fastly/go-fastly/pull/167)).
- heroku: Add Heroku Logplex logging endpoint support ([#166](https://github.com/fastly/go-fastly/pull/166)).
- ftp: Update FTP logging endpoint support to include `PublicKey` ([#165](https://github.com/fastly/go-fastly/pull/165)).
- elasticsearch: Add Elasticsearch logging endpoint support ([#164](https://github.com/fastly/go-fastly/pull/164)).
- digitalocean: Add DigitalOcean Spaces logging endpoint support ([#163](https://github.com/fastly/go-fastly/pull/163)).
- rackspace: Add Rackspace Cloud Files logging endpoint support ([#162](https://github.com/fastly/go-fastly/pull/162)).
- test: Improve testing experience ([#161](https://github.com/fastly/go-fastly/pull/161)).
- doc: Fix typos in `GetRealtimeStats` documentation ([#160](https://github.com/fastly/go-fastly/pull/160)).

## [v1.10.0] (2020-04-24)

### Added

- tls: Add support for Platform TLS API endpoints ([#154](https://github.com/fastly/go-fastly/pull/154)).

## [v1.9.0] (2020-04-23)

### Changed

- splunk: Add missing TLS fields to the Splunk logging endpoint ([#156](https://github.com/fastly/go-fastly/pull/156)).
- https: Add support for HTTPS logging endpoints ([#155](https://github.com/fastly/go-fastly/pull/155)).

## [v1.8.0] (2020-04-21)

### Changed

- client: Add NewRealtimeStatsClientForEndpoint API ([#152](https://github.com/fastly/go-fastly/pull/152)).

## [v1.7.2] (2020-03-30)

### Changed

- client: Allow purge requests to run in parallel ([#147](https://github.com/fastly/go-fastly/pull/147)).

## [v1.7.1] (2020-03-24)

### Changed

- client: Serialize all non readable requests ([#146](https://github.com/fastly/go-fastly/pull/146)).

## [v1.7.0] (2020-02-26)

### Added

- user: Support for Fastly's User Management ([#145](https://github.com/fastly/go-fastly/pull/145)).

### Changed

- purge: Request method for purging an individual URL ([#116](https://github.com/fastly/go-fastly/pull/116)).

## [v1.6.0] (2020-02-18)

### Added

- s3: Support for `server_side_encryption_kms_key_id` and `server_side_encryption` fields ([#144](https://github.com/fastly/go-fastly/pull/144)).

## [v1.5.0] (2020-01-29)

### Added

- pool/server: Support for Fastly's Load Balancer ([#142](https://github.com/fastly/go-fastly/pull/142)).

## [v1.4.0] (2020-01-06)

### Added

- dictionary_info: Support to retrieve metadata for a single dictionary ([#122](https://github.com/fastly/go-fastly/pull/122)).
- syslog: Support for `tls_client_cert` and `tls_client_key` fields ([#139](https://github.com/fastly/go-fastly/pull/139)).

## [v1.3.0] (2019-10-02)

### Added

- vcl_snippets: Support for `hash` type ([#133](https://github.com/fastly/go-fastly/pull/133)).
- service: Support for `type` field ([#132](https://github.com/fastly/go-fastly/pull/132)).
- token: Support for API tokens ([#131](https://github.com/fastly/go-fastly/pull/131)).

### Changed

- client: Codebase dependency management from `dep` to Go modules ([#130](https://github.com/fastly/go-fastly/pull/130)).

## [v1.2.1] (2019-07-25)

### Added

- acl: Constant to represent the maximum number of entries that can be placed within an ACL ([#129](https://github.com/fastly/go-fastly/pull/129)).
- dictionary: Constant to represent the maximum number of items that can be placed within an Edge Dictionary ([#129](https://github.com/fastly/go-fastly/pull/129)).

## [v1.2.0] (2019-07-24)

### Added

- acl: Support for Create, Delete and Update BatchOperations ([#126](https://github.com/fastly/go-fastly/pull/126)).

## [v1.1.0] (2019-07-22)

### Added

- dictionary: Support for Create, Delete, Update and Upsert BatchOperations ([#125](https://github.com/fastly/go-fastly/pull/125)).

## [v1.0.0] (2019-06-14)

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

## [v0.4.3] (2019-05-24)

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

## [v0.4.2] (2017-09-05)

### Added

- logentries: Support for `format_version` field ([#50](https://github.com/fastly/go-fastly/pull/50)).
- gcs: Support for `message_type` field ([#52](https://github.com/fastly/go-fastly/pull/52)).
- waf: Support for firewall ([216f9cb](https://github.com/fastly/go-fastly/commit/216f9cb6a92bc6e3c4653b7ebc9206f78a80d69b), [c6feafe](https://github.com/fastly/go-fastly/commit/c6feafe0fc5ed2b74bef9d3105f2f20c6197b19e), [50fef06](https://github.com/fastly/go-fastly/commit/50fef061051d188edcc37749a108fc6d025e495c)).

## [v0.4.1] (2017-08-07)

### Added

- syslog: Support for `hostname`, `ipv4`, `tls_hostname`, and `message_type` fields ([2b863da](https://github.com/fastly/go-fastly/commit/2b863da88fc1033a68538ccdc5c9dc82fa52681f)).

## [v0.4.0] (2017-07-27)

### Added

- realtime_stats: Support for real-time analytics ([#48](https://github.com/fastly/go-fastly/pull/48)).

### Changed

- acl: Names of all types, functions, and variables to follow Go standards ([#46](https://github.com/fastly/go-fastly/pull/46)).

### Fixed

- condition: URL encoding for forward slashes ([3d6dabb](https://github.com/fastly/go-fastly/commit/3d6dabb37bd2df7195d28aef08b1edd98895b960)).

## [v0.3.0] (2017-07-19)

- Initial tagged release

[v0.3.0]: https://github.com/fastly/go-fastly/releases/tag/v0.3.0
[v0.4.0]: https://github.com/fastly/go-fastly/compare/v0.3.0...v0.4.0
[v0.4.1]: https://github.com/fastly/go-fastly/compare/v0.4.0...v0.4.1
[v0.4.2]: https://github.com/fastly/go-fastly/compare/v0.4.1...v0.4.2
[v0.4.3]: https://github.com/fastly/go-fastly/compare/v0.4.2...v0.4.3
[v1.0.0]: https://github.com/fastly/go-fastly/compare/v0.4.3...v1.0.0
[v1.1.0]: https://github.com/fastly/go-fastly/compare/v1.0.0...v1.1.0
[v1.10.0]: https://github.com/fastly/go-fastly/compare/v1.9.0...v1.10.0
[v1.11.0]: https://github.com/fastly/go-fastly/compare/v1.10.0...v1.11.0
[v1.12.0]: https://github.com/fastly/go-fastly/compare/v1.11.0...v1.12.0
[v1.13.0]: https://github.com/fastly/go-fastly/compare/v1.12.0...v1.13.0
[v1.14.0]: https://github.com/fastly/go-fastly/compare/v1.13.0...v1.14.0
[v1.2.0]: https://github.com/fastly/go-fastly/compare/v1.1.0...v1.2.0
[v1.2.1]: https://github.com/fastly/go-fastly/compare/v1.2.0...v1.2.1
[v1.3.0]: https://github.com/fastly/go-fastly/compare/v1.2.1...v1.3.0
[v1.4.0]: https://github.com/fastly/go-fastly/compare/v1.3.0...v1.4.0
[v1.5.0]: https://github.com/fastly/go-fastly/compare/v1.4.0...v1.5.0
[v1.6.0]: https://github.com/fastly/go-fastly/compare/v1.5.0...v1.6.0
[v1.7.0]: https://github.com/fastly/go-fastly/compare/v1.6.0...v1.7.0
[v1.7.1]: https://github.com/fastly/go-fastly/compare/v1.7.0...v1.7.1
[v1.7.2]: https://github.com/fastly/go-fastly/compare/v1.7.1...v1.7.2
[v1.8.0]: https://github.com/fastly/go-fastly/compare/v1.7.2...v1.8.0
[v1.9.0]: https://github.com/fastly/go-fastly/compare/v1.8.0...v1.9.0
