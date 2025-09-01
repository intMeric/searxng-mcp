package config

import (
	_ "embed"
)

//go:embed settings.yml
var DefaultSettings string