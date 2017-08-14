# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Changed
- Extracted request metrics out of context

## [1.0.3] - 2017-08-11
### Changed
- Fixed issue that used wrong map key for total time metric
- Fixed issue that reported api key tag to wrong span

## [1.0.2] - 2017-08-11
## Added
- Improved test coverage.
### Changed
- Fixed issue where Docker images were not pushed for tag builds.
- Fixed issue that always reported 0 for total request time
- Fixed issue that did not report span if plugin return an error

## [1.0.1] - 2017-08-10
### Added
- Improved test coverage.
### Changed
- Upgraded to Go version 1.8.3
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