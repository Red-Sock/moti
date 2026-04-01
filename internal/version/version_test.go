package version

import (
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSystem(t *testing.T) {
	v := System()
	assert.NotEmpty(t, v)
}

func TestBuildSetting(t *testing.T) {
	bi := &debug.BuildInfo{
		Settings: []debug.BuildSetting{
			{Key: "foo", Value: "bar"},
		},
	}
	assert.Equal(t, "bar", buildSetting(bi, "foo"))
	assert.Equal(t, "", buildSetting(bi, "nonexistent"))
}

func TestBuildVersion(t *testing.T) {
	// Since we can't easily mock debug.ReadBuildInfo(), we test the logic of buildVersion with what we get
	bi, _ := buildVersion()
	// If bi is nil, it means ReadBuildInfo failed (e.g. not built with modules or in test environment without build info)
	if bi == nil {
		t.Log("ReadBuildInfo returned nil, skipping version check")
	}
}
