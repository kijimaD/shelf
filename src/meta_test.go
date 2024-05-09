package shelf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMetas(t *testing.T) {
	tomlContent := `
["20240310T224413832103518"]
title = "example1"
todo = "NONE"
tags = ["new"]
["20240310T224413832109999"]
title = "example2"
todo = "NONE"
tags = ["new"]
`

	metas, err := GetMetas(tomlContent)
	assert.NoError(t, err)

	expect := Metas{
		"20240310T224413832103518": Meta{
			Title: GetPtr("example1"),
			TODO:  GetPtr(TODOTypeNONE),
			Tags:  GetPtr([]string{"new"}),
		},
		"20240310T224413832109999": Meta{
			Title: GetPtr("example2"),
			TODO:  GetPtr(TODOTypeNONE),
			Tags:  GetPtr([]string{"new"}),
		},
	}
	assert.Equal(t, expect, metas)
}
