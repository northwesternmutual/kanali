# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [1.2.3] - 2017-11-12
### Changed
- Allow for batching of InfluxDB writes.
- Fixed [#83](https://github.com/northwesternmutual/kanali/issues/83).
- Fixed bug that did not properly handle all config value types.

## [1.2.2] - 2017-11-02
### Changed
- Support regular expressions for ApiKeyBinding key subpaths.

## [1.2.1] - 2017-10-13
### Added
- Adding `ClusterRole` in `scripts/install.sh` so that install works with latest Minikube releases.
### Changed
- Upstream URLs will now be properly encoded/decoded.
- Improved test coverage.
### Removed
- Gzip support. Kanali will be a transparent proxy.

## [1.2.0] - 2017-09-24
### Added
- Gzip support.
- `Update` method to `Store` interface.
- Additional Opentracing span to capture target request.
- _Average Target Response Time_ Grafana panel.
### Changed
- Fixed [#61](https://github.com/northwesternmutual/kanali/issues/61)
- Fixed [#57](https://github.com/northwesternmutual/kanali/issues/57).
- Using versioned API key plugin in distributed `Dockefile` and examples.
- Refactored configuration items.
- Using `TOML` as configuration format in distribution and examples.
- Improved test coverage.
### Removed
- `controller.Controller` from `Plugin` interface method parameters

## [1.1.4] - 2017-08-31
### Added
- Store interface implementation for mock responses.
- Slack integration for TravisCI and Coveralls.
### Changed
- Improved test coverage.
### Removed
- `Contains` method from `Store` interface.

## [1.1.3] - 2017-08-26
### Added
- *Go Report* GitHub badge
- [License](./LICENSE) file
### Changed
- If connection to the k8s apiserver is lost, it will be reattempted after 5 seconds.
- When TPRs are created, the server will not start until it is confirmed that they are created.
- Kanali process will not terminate if a connection to Jaeger cannot be made.
- Response will not wait for metrics to be written to InfluxDB before completing.
### Removed
- Deprecated `status-port` flag. The status server has been removed.
- Deprecated `enable-tracing` flag. Like InfluxDB, a best effort at a connection will be made.

## [1.1.2] - 2017-08-18
### Changed
- If InfluxDB database doesn't exist when writing, it will be created.
- Raw URL path will be used as the primary path when proxying to upstream services.
- Cleaning up logging.

## [1.1.1] - 2017-08-15
### Changed
- `TrafficFactory` now implements the `Store` interface.

## [1.1.0] - 2017-08-15
### Changed
- Extracted request metrics out of context resulting in an updated `Plugin` interface.
- Enhanced Grafana dashboard.
- Updated docs to reflect change to `Plugin` interface.

## [1.0.3] - 2017-08-11
### Changed
- Fixed issue that used wrong map key for total time metric
- Fixed issue that reported api key tag to wrong span

## [1.0.2] - 2017-08-11
### Added
- Improved test coverage.
### Changed
- Fixed issue where Docker images were not pushed for tag builds.
- Fixed issue that always reported 0 for total request time
- Fixed issue that did not report span if plugin return an error

## [1.0.1] - 2017-08-10
### Added
- Improved test coverage.
### Changed
- Upgraded Go version to `1.8.3`
- Resolved bug that did not handle potential panics from plugins.
- Fixed typos in README.
- Updated Postman collection.
- Using read/write mutex to Metrics map. Before this change, it was not safe for concurrent use.
- Resolved bug that resulted in requested HTTP header values to not be properly omitted when being reported to Influxdb.
- Travis will now retry the goveralls step if not successful the first time.
- Docker images will now only be built and pushed on pushes to the master branch. This resolves an issue where images were being built and pushed on pull requests where the master branch was the target.

## [1.0.0] - 2017-07-31
### Added
- Initial Project Commit.