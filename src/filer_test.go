package shelf

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateMeta(t *testing.T) {
	tempdir := t.TempDir()
	{
		fname := filepath.Join(tempdir, "20010101T010101_abc.pdf")
		assert.NoError(t, os.WriteFile(fname, []byte(""), 0666))

	}
	{
		fname := filepath.Join(tempdir, "20010101T010101.toml")
		assert.NoError(t, os.WriteFile(fname, []byte(""), 0666))
	}

	GenerateFiler(tempdir)

}
