package src

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	tempfile, err := os.CreateTemp(os.TempDir(), "test-*.pdf")
	assert.NoError(t, err)
	assert.NoError(t, Register(tempfile.Name()))
}
