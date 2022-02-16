package internal

import "runtime"

var Version = "Unknown"
var GoVersion = runtime.Version()
var Buildtime = "Unknown"
var Commit = "Unknown"
var Os = runtime.GOOS
var Arch = runtime.GOARCH
