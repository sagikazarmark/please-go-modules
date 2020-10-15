package main

// Platform represents a single build targe platform.
type Platform struct {
	OS   string
	Arch string
}

// SupportedPlatforms lists all the supported platforms.
var SupportedPlatforms = []Platform{
	{"linux", "amd64"},
	{"darwin", "amd64"},
}
