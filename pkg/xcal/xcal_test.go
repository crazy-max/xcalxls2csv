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

func TestConvertToCSVAndWrite(t *testing.T) {
	f := "./fixtures/MCF001120_20230831_9_Short.xls"
	res, err := ConvertToCSV(f)
	assert.NoError(t, err)
	t.Log(string(res))

	out := fmt.Sprintf("%s.csv", strings.TrimSuffix(f, filepath.Ext(f)))
	t.Cleanup(func() {
		assert.NoError(t, os.Remove(out))
	})
	assert.NoError(t, os.WriteFile(out, res, 0644))
}

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
		})
	}
}
