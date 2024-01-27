package redirect

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetStorage tests flag wiring to configuration, as well as configuration wiring to functions that
// use it.
func TestGetRepository(t *testing.T) {
	tmpDir, err := os.MkdirTemp("/tmp", "TestGetStorage")
	assert.Nil(t, err)

	for _, tc := range []struct {
		name string

		fName, fVal string
		err         error

		setup, teardown func()
	}{
		{
			name:  "in memory storage",
			fName: flagHashMapStorage,
			fVal:  "true",
		},
		{
			name:  "disabled in memory storage",
			fName: flagHashMapStorage,
			fVal:  "false",
			err:   failedStorageSetup,
		},
		{
			name:  "yaml storage",
			fName: flagYamlStorage,
			fVal:  tmpDir + "/urls.yaml",
			setup: func() {
				file, err := os.Create(tmpDir + "/urls.yaml")
				if err != nil {
					panic(err)
				}

				file.Write([]byte(`
---
- from: //www/google
  to: //www/moutaz
`))
				file.Close()
			},
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// Allow setting up whatever resources need to exist for the storage (e.g. temporary directories to
			// store content)
			if tc.setup != nil {
				tc.setup()
			}
			if tc.teardown != nil {
				defer tc.teardown()
			}

			// Replicate the supplied user option.
			serverFlagSet.Set(tc.fName, tc.fVal)

			_, err := getRepository(serverFlagSet)
			assert.ErrorIs(t, tc.err, err)
		})
	}
}
