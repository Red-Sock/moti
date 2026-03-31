package models

import (
	"strings"
)

const (
	// If version was omitted
	Omitted RequestedVersion = ""
)

type (
	// ModuleHash alias for module's hash
	// used in lock file for verification
	ModuleHash string

	// RequestedVersion for installing
	RequestedVersion string
)

// Module contain requested dependency name and its version
type Module struct {
	Name    string           // Full path on remote repository
	Version RequestedVersion // Version obtained from config (Omitted if version was omitted)
}

type GeneratedVersionParts struct {
	CommitHash string
}

func NewModule(dependency string) Module {
	parts := strings.Split(dependency, "@")
	name := parts[0]

	version := Omitted

	if len(parts) > 1 {
		version = RequestedVersion(parts[1])
	}

	return Module{
		Name:    name,
		Version: version,
	}
}

// IsGenerated check if requested version is a commit hash
func (v RequestedVersion) IsGenerated() bool {
	return v.IsHex()
}

// IsHex check if requested version is a hex string (commit hash)
func (v RequestedVersion) IsHex() bool {
	return isHex(string(v))
}

// IsCommitHash check if requested version is a full commit hash (40-char hex)
func (v RequestedVersion) IsCommitHash() bool {
	return len(v) == 40 && v.IsHex()
}

func isHex(s string) bool {
	if len(s) < 7 {
		return false
	}

	for _, c := range s {
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') && (c < 'A' || c > 'F') {
			return false
		}
	}

	return true
}

// IsOmitted check if requested version is omitted
func (v RequestedVersion) IsOmitted() bool {
	return v == Omitted
}

func (g GeneratedVersionParts) GetVersionString() string {
	return g.CommitHash
}
