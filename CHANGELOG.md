# go-fastly CHANGELOG

## v0.4.3 (Unreleased)

- Add WAF methods for fetching status of rules, both one at a time and in filtered lists
- Add WAF methods for modifying the status of rules, both one at a time and based on tags
- Rename `UpdateWafRuleSets` function to `UpdateWAFRuleSets` to match other names

## v0.4.2 (September 5, 2017)

- Add support for specifying `message_type` in GCS [GH-52]
- Add support for WAF [GH-35, GH-55]

## v0.4.1 (August 7, 2017)

- Add `hostname`, `ipv4`, `tls_hostname`, and `message_type` to syslog [GH-49]

## v0.4.0 (July 27, 2017)

FEATURES:

- Add support for real-time stats [GH-48]

## v0.3.0 (July 19, 2017)

- Initial tagged release
