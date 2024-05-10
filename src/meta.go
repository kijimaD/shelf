package shelf

import "github.com/BurntSushi/toml"

type TODOType string

var (
	TODOTypeNONE TODOType = "NONE"
	TODOTypeTODO TODOType = "TODO"
	TODOTypeWIP  TODOType = "WIP"
	TODOTypeDONE TODOType = "DONE"
)

// 1ドキュメントがもつメタ情報
type Meta struct {
	Title *string   `toml:"title"`
	TODO  *TODOType `toml:"todo"`
	Tags  *[]string `toml:"tags"`
}

type Metas map[string]Meta

func GetMetas(tomlContent string) (Metas, error) {
	var metas Metas
	if _, err := toml.Decode(tomlContent, &metas); err != nil {
		return nil, err
	}

	return metas, nil
}
