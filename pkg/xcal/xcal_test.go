package xcal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertToCSV(t *testing.T) {
	fixtures, err := os.ReadDir("./fixtures")
	require.NoError(t, err)
	for _, f := range fixtures {
		f := f
		t.Run(f.Name(), func(t *testing.T) {
			f := "./fixtures/" + f.Name()
			res, err := ConvertToCSV(f)
			assert.NoError(t, err)
			t.Log(string(res))
			assert.NoError(t, os.WriteFile(fmt.Sprintf("%s.csv", strings.TrimSuffix(f, filepath.Ext(f))), res, 0644))
		})
	}
}
