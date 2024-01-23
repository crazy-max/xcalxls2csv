package xcal

import (
	"os"
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
			res, err := ConvertToCSV("./fixtures/" + f.Name())
			assert.NoError(t, err)
			t.Log(string(res))
		})
	}
}
