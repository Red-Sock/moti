package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewModule(t *testing.T) {
	tests := map[string]struct {
		dependency     string
		expectedResult Module
	}{
		"with version": {
			dependency: "github.com/company/repository@v1.2.3",
			expectedResult: Module{
				Name:    "github.com/company/repository",
				Version: "v1.2.3",
			},
		},
		"without version": {
			dependency: "github.com/company/repository",
			expectedResult: Module{
				Name:    "github.com/company/repository",
				Version: Omitted,
			},
		},
	}

	for name, tc := range tests {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			result := NewModule(tc.dependency)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func Test_RequestedVersion_IsGenerated(t *testing.T) {
	tests := map[string]struct {
		requestedVersion RequestedVersion
		expectedResult   bool
	}{
		"not generated, simple tag": {
			requestedVersion: RequestedVersion("v1.2.3"),
			expectedResult:   false,
		},
		"not generated, tag with no `v` prefix": {
			requestedVersion: RequestedVersion("some_tag"),
			expectedResult:   false,
		},
		"not generated, with `-`": {
			requestedVersion: RequestedVersion("v1.2.3-rc"),
			expectedResult:   false,
		},
		"not generated, with several `-`": {
			requestedVersion: RequestedVersion("v1.2.3-rc-111222"),
			expectedResult:   false,
		},
		"Use Omitted": {
			requestedVersion: Omitted,
			expectedResult:   false,
		},
		"not generated, pseudo-version": {
			requestedVersion: "v0.0.0-20240222234643-814bf88cf225",
			expectedResult:   false,
		},
		"commit hash": {
			requestedVersion: "220e0db758f9ce96d9b1f457234616284530622b",
			expectedResult:   true,
		},
	}

	for name, tc := range tests {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			result := tc.requestedVersion.IsGenerated()
			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func Test_RequestedVersion_IsCommitHash(t *testing.T) {
	tests := map[string]struct {
		requestedVersion RequestedVersion
		expectedResult   bool
	}{
		"not commit hash, simple tag": {
			requestedVersion: RequestedVersion("v1.2.3"),
			expectedResult:   false,
		},
		"not commit hash, pseudo-version": {
			requestedVersion: "v0.0.0-20240222234643-814bf88cf225",
			expectedResult:   false,
		},
		"is commit hash": {
			requestedVersion: RequestedVersion("220e0db758f9ce96d9b1f457234616284530622b"),
			expectedResult:   true,
		},
		"too short hash": {
			requestedVersion: RequestedVersion("220e0db"),
			expectedResult:   false,
		},
		"invalid characters": {
			requestedVersion: RequestedVersion("220e0db758f9ce96d9b1f457234616284530622g"),
			expectedResult:   false,
		},
	}

	for name, tc := range tests {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			result := tc.requestedVersion.IsCommitHash()
			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func Test_RequestedVersion_IsOmitted(t *testing.T) {
	tests := map[string]struct {
		requestedVersion RequestedVersion
		expectedResult   bool
	}{
		"not omitted": {
			requestedVersion: RequestedVersion("v1.2.3"),
			expectedResult:   false,
		},
		"omitted": {
			requestedVersion: Omitted,
			expectedResult:   true,
		},
	}

	for name, tc := range tests {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			result := tc.requestedVersion.IsOmitted()
			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func Test_GeneratedVersionParts_GetVersionString(t *testing.T) {
	tests := map[string]struct {
		parts          GeneratedVersionParts
		expectedResult string
	}{
		"case 1": {
			parts:          GeneratedVersionParts{CommitHash: "814bf88cf225"},
			expectedResult: "814bf88cf225",
		},
		"case 2": {
			parts:          GeneratedVersionParts{CommitHash: "914af88cf235"},
			expectedResult: "914af88cf235",
		},
	}

	for name, tc := range tests {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			result := tc.parts.GetVersionString()
			require.Equal(t, tc.expectedResult, result)
		})
	}
}
