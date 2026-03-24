package generate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuery_Build(t *testing.T) {
	t.Run("grouping_files_by_directory", func(t *testing.T) {
		q := Query{
			Compiler: "protoc",
			Files: []string{
				"proto/a.proto",
				"proto/b.proto",
				"proto/sub/c.proto",
				"external/d.proto",
			},
			Imports: []string{"."},
			Plugins: []Plugin{
				{Name: "go", Out: "gen/go"},
			},
		}

		_, args := q.Build()

		// It should contain unique imports including parents of the files
		assert.Contains(t, args, "-I .")
		assert.Contains(t, args, "-I proto")
		assert.Contains(t, args, "-I proto/sub")
		assert.Contains(t, args, "-I external")

		// It should contain all files
		assert.Contains(t, args, "proto/a.proto")
		assert.Contains(t, args, "proto/b.proto")
		assert.Contains(t, args, "proto/sub/c.proto")
		assert.Contains(t, args, "external/d.proto")
	})
}
