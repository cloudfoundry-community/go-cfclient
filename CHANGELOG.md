# Changelog

## [v3.0.0-beta.2]

### Added

- `AppPushOperation.WithDockerPassword` to set the docker registry password per push operation. It takes precedence over the `CF_DOCKER_PASSWORD` environment variable, so concurrent pushes in one process no longer have to share a process-wide variable.

### Changed

- Minimum supported Go version increased from 1.23.0 to 1.26.0.
- Updated `golang.org/x/oauth2` from v0.30.0 to v0.38.0.

### Fixed

- Timestamp filters now serialize in UTC, preserving the supplied instant. Filters such as `CreatedAts` and `UpdatedAts` now accept timestamps with non-UTC offsets without triggering an invalid timestamp query error from the Cloud Foundry API ([#472](https://github.com/cloudfoundry/go-cfclient/issues/472)). Caller-provided timestamps are unchanged.

## [v3.0.0-beta.1]

### Added

- Support for setting and reading application lifecycle (`buildpack`, `docker`, and `cnb`/Cloud Native Buildpacks) in V3 API.
- Custom marshaling/unmarshaling for `Lifecycle` struct to support multiple lifecycle types.
- Unit tests for all lifecycle types, including CNB.
- Support for Service Broker-provided metadata (labels and attributes) on Service Instances. This includes the `BrokerProvidedMetadata` field on `ServiceInstance`, `ServiceInstanceManagedCreate`, and `ServiceInstanceManagedUpdate` structs, along with fluent builder methods for managing labels and attributes.

### Changed

- All lifecycle-related test expectations updated to match new marshaling output (both `type` and `data` fields).

### Notes

- For CNB, only `buildpacks` and `stack` fields are currently supported. Additional fields may be added in the future as the API evolves.
- No breaking changes for existing users of buildpack or docker lifecycles.
