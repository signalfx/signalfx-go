# 1.18.0, 5 Apr 2022

* Add `syncCustomNamespacesOnly` fields to the `AwsCloudWatchIntegration` struct [#160](https://github.com/signalfx/signalfx-go/pull/160)

# 1.17.0, 1 Apr 2022

No functional changes, this is an equivalent of `1.12.0`.

We are releasing it to succeed mistakenly published `1.16.13` instead of `1.6.13`, which since Jan 2020 shows as the latest.

# 1.12.0, 17 Mar 2022

Add support for ServiceNow integration [#159]

# 1.11.0, 16 Mar 2022

* Add `IncludeList` field to `GCPIntegration` struct [#157](https://github.com/signalfx/signalfx-go/pull/157) and stop using deprecated `Whitelist` field

# 1.10.0, 14 Mar 2022

* Add `AdditionalServices` and `ResourceFilterRules` fields to the `AzureIntegration` struct [#156](https://github.com/signalfx/signalfx-go/pull/156)

# 1.9.0, 11 Mar 2022

* Add support for read permissions for dashboard and dashboard groups

* Add `AWS/S3/Storage-Lens` to AWS services list. Sorted AWS/Azure/GCP services by name. [#154](https://github.com/signalfx/signalfx-go/pull/154)

* Add `PollRateMs` (type `int64`, value in milliseconds) to define Azure and GCP integration poll rate [#153](https://github.com/signalfx/signalfx-go/pull/153) [#155](https://github.com/signalfx/signalfx-go/pull/155)

`PollRate` (type `PollRate`) is deprecated. Please use `PollRateMs` instead.<BR>
`PollRateMs` accepts any value between 1 and 10 minutes (`60000` - `600000`).

# 1.8.9, 3 Mar 2022

* signalflow: Fix race in computation buffering when closing (previously the messages were not getting fully flushed to the downstream channels) [#151](https://github.com/signalfx/signalfx-go/pull/151)

# 1.8.8, 28 Feb 2022

* Add MetricStreamsSyncState to AWS Cloudwatch integration [#147](https://github.com/signalfx/signalfx-go/pull/147)

# 1.8.6, 8 Feb 2022

* Adding AuthScopes attribute to the org token struct [#146](https://github.com/signalfx/signalfx-go/pull/146)

# 1.8.5, 27 Jan 2022

* signalflow fakebackend support for timestamp control and END_OF_CHANNEL [#145](https://github.com/signalfx/signalfx-go/pull/145)

# 1.8.4, 24 Jan 2022

* Adding ValidateDetector and its unit tests [#142](https://github.com/signalfx/signalfx-go/pull/142)

# 1.8.3, 17 Nov 2021

* Add newly supported AWS services to the list. [#140](https://github.com/signalfx/signalfx-go/pull/140)

# 1.8.2, 30 Sep 2021

*  Update dependencies [#135](https://github.com/signalfx/signalfx-go/pull/135)

# 1.8.1, 24 Aug 2021

* Context cancellation handling [#131](https://github.com/signalfx/signalfx-go/pull/131)
* Allow to run PostDisconnectCallback

# 1.8.0, 2021-05-11

* Don't send empty orderBy query parameter for metric metadata calls [#127](https://github.com/signalfx/signalfx-go/pull/127).
* Replace nonfunctioning `RetrieveMetricMetadataResponseModel.Result` field with `Results` [#127](https://github.com/signalfx/signalfx-go/pull/127).

# 1.7.18, 2021-03-29

* Fix handling of client base URLs that have paths in them

# 1.7.17, 2021-03-23

* Update list of Azure services [#123](https://github.com/signalfx/signalfx-go/pull/123)

# 1.7.15, 2021-02-10

* Added `TimeZone` property to Detector. [#119](https://github.com/signalfx/signalfx-go/pull/119)

# 1.7.14, 2021-01-28

* Return events/alerts from signalflow events/alerts queries

# 1.7.13, 2020-12-07

*  signalflow: Fix channel end/abort control message handling

# 1.7.12, 2021-02-10

* Added `hideMissingValues` to chart.Options to show or hide missing values in the chart. [#111](https://github.com/signalfx/signalfx-go/pull/111)
* Added `minDelay` field to detector.
* Added `PagerDutyIntegrationGetByName` method.

# 1.7.11, 2020-10-19

* Fix processing ControlMessage in signalflow.Computation

# 1.7.10, 2020-08-06
* Don't omit chart.description from JSON, as API doesn't clear it if absent. [#107](https://github.com/signalfx/signalfx-go/pull/107)

# 1.7.9, 2020-08-06
* Allow AWS poll rate to be an int64. [#106](https://github.com/signalfx/signalfx-go/pull/106)

# 1.7.8, 2020-08-06
* Added `detectorName` to incident model. Thanks [choo-stripe](https://github.com/choo-stripe)! [#103](https://github.com/signalfx/signalfx-go/pull/103)
* Clean up some resources when a context is canceled during Signalflow computations. Thanks [kerbyhughes](https://github.com/kerbyhughes)! [#104](https://github.com/signalfx/signalfx-go/pull/104)
* Add new team linking methods for detectors and dashboard groups. [#102](https://github.com/signalfx/signalfx-go/pull/102)
* Add new team fields. [#105](https://github.com/signalfx/signalfx-go/pull/105)

# 1.7.7, 2020-07-22
* Added new fields to Azure integration.
* Fix race condition in SignalFlow client Close operation that could cause
  panics.

# 1.7.5, 2020-06-11
* Fix busted accessor.

# 1.7.4, 2020-06-11
* Finish support for additional computation fields for job metadata. Thanks [shrivu-stripe](https://github.com/shrivu-stripe)! [#93](https://github.com/signalfx/signalfx-go/pull/93)

# 1.7.3, 2020-06-11
* Added additional computation fields for job metadata. Thanks [shrivu-stripe](https://github.com/shrivu-stripe)! [#92](https://github.com/signalfx/signalfx-go/pull/92)

# 1.7.2, 2020-06-10
* Added `packageSpecifications` field to Detector`

# 1.7.0, 2020-06-02

* Added contexts to everything, switch to `NewRequestWithContext` for HTTP.

# 1.6.39, 2020-06-02

## Improvements

* Removed the go:generate comment out of the writer codegen template to make it
  easier to consume in external projects
* Added `NamedToken` to integrations.

# 1.6.37, 2020-05-18

## Improvements

Added GCP service type to match the other integrations. [#86](https://github.com/signalfx/signalfx-go/pull/86)

# 1.6.36, 2020-05-14

## Improvements

Various `Search*` methods now verify the response code. [#83](https://github.com/signalfx/signalfx-go/issues/83)

## Improvements

* Add better resolution management to the fake SignalFlow backend for external
  testing uses.

# 1.6.35, 2020-05-12

Fix some typos in Azure services

# 1.6.34, 2020-05-11

* Add new Azure services. [#82](https://github.com/signalfx/signalfx-go/pull/82)

# 1.6.33, 2020-04-29

* Adjust behavior to more reliably close/fully-read HTTP bodies with sketchy replies. [#81](https://github.com/signalfx/signalfx-go/pull/81)

# 1.6.32, 2020-04-13

## Improvements

* Added User-Agent client param and header to client.

# 1.6.31, 2020-04-10

## Bugfixes

* Don't set orderBy on dimension search if an empty string.

# 1.6.30, 2020-03-27

## Bugfixes

* Don't set tags on chart search if an empty string. Thanks [ChimeraCoder](https://github.com/ChimeraCoder)! [#79](https://github.com/signalfx/signalfx-go/pull/79)

# 1.6.29, 2020-03-20

## Bugfixes

* Adjust Data Link targets to include isDefault even if "empty".
* Don't set tags on detectors if an empty string. Thanks [rma-stripe](https://github.com/rma-stripe)! [#76](https://github.com/signalfx/signalfx-go/pull/77)

# 1.6.28, 2020-03-19

## Improvements

* Added `MatchedSize` and `LimitSize` to `Computation`. Thanks [rma-stripe](https://github.com/rma-stripe)! [#76](https://github.com/signalfx/signalfx-go/pull/76)

# 1.6.27, 2020-03-16

## Bugfixes

* Chart Axes' `HighWatermark`, `LowWatermark`, `Max`, and `Min` are now correctly typed as `float64`. Same for `ColorScale.Thresholds` and `SecondaryVisualization`'s fields (`Gt`, `Gte`, `Lt`, and `Lte`). [#75](https://github.com/signalfx/signalfx-go/pull/75).

## Improvements

* `Computation`s internal errors are now richer, allowing users to get to the code, message and type. Thanks [rma-stripe](https://github.com/rma-stripe)! [#74](https://github.com/signalfx/signalfx-go/pull/74)


# 1.6.25, 2020-03-13

## Bugfixes

* Fix typos in AWS service name for `AWS/VPN`

# 1.6.24, 2020-03-11

## Bugfixes

* Protect `Client` with a mutex so that multiple calls don't races. Thanks [rma-stripe](https://github.com/rma-stripe)! [#73](https://github.com/signalfx/signalfx-go/pull/73)

# 1.6.23, 2020-03-10

## Bugfixes

* Make the writer package work properly on 32-bit systems by aligning struct
  fields on 64-bit boundaries.

# 1.6.22, 2020-03-09

## Added
* New package `realm` to help with constructing SignalFx ingest and API urls
  from the realm name

# 1.6.21, 2020-03-04

## Bugfixes

* Fixed some errors in new AWS services

# 1.6.20, 2020-03-03

## Added

* Many new AWS services

# 1.6.19, 2020-02-26

## Added

* Webhook integration client functions

# 1.6.18, 2020-02-25

## Added

* New `datalink.EpochSeconds` for it's `TimeFormat`

# 1.6.17, 2020-02-18

## Added

* New methods `GetDetectorEvents` and `GetDetectorIncidents`

# 1.6.16, 2020-02-13

## Added

* Added `UseGetMetricDataMethod` to AWS integration

# 1.6.15, 2020-01-27

## Added

* Add `PublishLabelOptions` to Detector correctly this time

# 1.6.13, 2020-01-27

## Added

* Added `PublishLabelOptions` to Detector

# 1.6.12, 2020-01-21

## Added

* Field `sfxAwsAccountArn` added to AWS response

# 1.6.11, 2019-12-18

## Added

* Support for creating and deleting tokens using the Session API

# 1.6.10, 2019-12-16

## Added

* Methods for Data Links

# 1.6.9, 2019-12-09

## Added

* New datapoint and span writer for high volume output

## Bugfixes

* Token operations now URL encode the name.

# 1.6.8, Pending

## Added

* Methods for Alert Muting Rules

# 1.6.7, 2019-11-05

## Bugfixes

* Added `AuthorizedWriters` to Detector model

# 1.6.6, 2019-11-05

## Bugfixes

* Detector and DashboardGroup structs modified to use a pointer for `AuthorizedWriters`.

# 1.6.5, 2019-10-30

## Added

* Additional reconnect delays upon SignalFlow socket errors to reduce load on
backend.
* Added `*JiraIntegration` methods
* Added `notification.JiraNotification`

# 1.6.4, 2019-09-27

## Added

Event Overlays now support a detector id.

# 1.6.3, 2019-09-19

## Bugfixes

* Changed detector's time fields to be `*int64`

# 1.6.2, 2019-09-16

## Added

* VictorOps integration functions

## Updated

* Adjusted `EventPublishLabelOptions.PalleteIndex` to an `*int32` to match other uses.
* SignalFlow computation Handle() method wait for handle to come in until
  returning (with timeout).
* Renamed `BinaryPayload` to `DataPayload` in the `messages` package.
* Exported `BinaryMessageHeader` and `DataMessageHeader` from `messages`
  package to facilitate low-level SignalFlow parsing.

## Bugfixes

* SignalFlow client connection handling was refactored to prevent deadlocks
  that could occur on reconnects and bad authentication.

## Removed

# 1.6.1, 2019-08-16

## Updated

* Adjusted detector.CreateUpdateDetectorRequest to use pointer for Rules

# 1.6.0, 2019-08-16

## Added

* Added `*GCPIntegration` methods
* Added `*Opsgenie` methods
* Added `*PagerDutyIntegration` methods
* Added `*SlackIntegration` methods

## Updated

* `Detector.Rules` now uses `Notification` as it's type instead of an untyped `[]map[string]interface{}`.

## Removed
* Renamed `integration.GcpIntegration` and it's sub-types to `GCP`, fixing case.

# 1.5.0, 2019-08-05

## Added

* Add OrgToken methods

## Bugfixes

* Properly recognize the SignalFlow keep alive event message and ignore it.

## Updated

* Moved various notification bits into a `notification` package

# 1.4.0, 2019-07-29

## Added

* Add `*AzureIntregration` functions to client.

## Updated

## Bugfixes

## Removed

# 1.3.0, 2019-07-24

## Added

* Added OpenAPI code for integrations, experimental for now.
* Add `*AwsCloudWatchIntegration` functions to client.

## Removed

* Removed `credentialName` from Opsgenie notifications, not a real field in the API.

# 1.2.0, 2019-07-16

## Updated

* Many numeric properties have been adjusted to pointers to play better with Go's JSON (un)marshaling.

# 1.1.0, 2019-07-15

## Added
* Added `DashboardConfigs` to `CreateUpdateDashboardRequest`
* `DashboardGroupCreate` now has an option to create an empty group.

## Updated
* Many types have been changed to pointers to add (de)serialization
* Moved `StringOrSlice` into a `util` package, cuz all projects must have one

## Bugfixes
* Switched to `StringOrSlice` for some fields that needed it.
* Added `StringOrInteger` to handle failures in some Chart filter responses, thanks to (doctornkz)[https://github.com/doctornkz] for flagging!

## Removed

# 1.0.0, Forgot the date

Tagged!

## Added

## Updated

## Bugfixes

## Removed
