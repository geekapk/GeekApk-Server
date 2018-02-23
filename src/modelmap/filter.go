package modelmap

type CmpBase int

type FilterRule struct {
	Key string
	CompareType CmpBase
	Value string
}

const (
	CmpUnknown CmpBase = iota
	CmpEq
	CmpNe
	CmpGt
	CmpGe
	CmpLt
	CmpLe
	CmpLike
)
