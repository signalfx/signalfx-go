package signalfx

//go:generate gojson -input chart.json -o chart.go -fmt json -pkg signalfx -name Chart
//go:generate gojson -input dashboard.json -o dashboard.go -fmt json -pkg signalfx -name Dashboard
//go:generate gojson -input dashboard_group.json -o dashboard_group.go -fmt json -pkg signalfx -name DashboardGroup
//go:generate gojson -input detector.json -o detector.go -fmt json -pkg signalfx -name Detector
//go:generate gojson -input team.json -o team.go -fmt json -pkg signalfx -name Team
