module github.com/ssotops/gitspace-plugin/examples/hello-world

go 1.23.1

require (
	github.com/charmbracelet/lipgloss v0.13.0
	github.com/charmbracelet/log v0.4.0
	github.com/ssotops/gitspace-plugin v0.0.0-20230701000000-abcdefghijkl
)

// This replace directive ensures we use the local version of gitspace-plugin
// Remove this line when using the published version
replace github.com/ssotops/gitspace-plugin => ../..
