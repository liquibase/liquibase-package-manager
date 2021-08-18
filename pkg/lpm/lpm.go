package lpm

import (
	_ "embed" // Embed Import for Package Files
)

//go:embed "embeds/VERSION"
var VersionNumber string

//PackagesJSON is embedded for first time run
//go:embed "embeds/packages.json"
var PackagesJSON []byte
