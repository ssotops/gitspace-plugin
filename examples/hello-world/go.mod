module github.com/ssotops/gitspace-plugin/examples/hello-world

go 1.23.1

require (
	github.com/charmbracelet/lipgloss v0.13.0
	github.com/charmbracelet/log v0.4.0
	github.com/ssotops/gitspace-plugin v1.0.11
)

require (
	github.com/Masterminds/semver/v3 v3.3.0 // indirect
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/charmbracelet/x/ansi v0.3.2 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/muesli/termenv v0.15.2 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	golang.org/x/exp v0.0.0-20240909161429-701f63a606c0 // indirect
	golang.org/x/sys v0.25.0 // indirect
)

// This replace directive ensures we use the local version of gitspace-plugin
// Remove this line when using the published version
replace github.com/ssotops/gitspace-plugin => ../..
