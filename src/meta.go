package shelf

type TODOType string

var (
	TODOTypeNONE TODOType = "NONE"
	TODOTypeTODO TODOType = "TODO"
	TODOTypeWIP  TODOType = "WIP"
	TODOTypeDONE TODOType = "DONE"
)

type Meta struct {
	Title string   `toml:"title"`
	TODO  TODOType `toml:"todo"`
	Tags  []string `toml:"tags"`
}
