package types

type Expander interface {
	Expand(text string) (string, string, bool)
	TypeExpansion(text string)
}
