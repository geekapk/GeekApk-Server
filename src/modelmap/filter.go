package modelmap

type CmpBase int

// A set of filter rules parsed from the request specifies the action
// that the client want to do or the data that the client want to fetch.
//
// There are two kinds of filter rules, explicit and implicit.
//
// Suppose that we have a URL like this:
//
// ````
// /Model/[1]/[2]
// ````
//
// By default, the [1] part corresponds to a rule with key `id`, and the
// [2] part corresponds to a rule with key `property`, and the `CompareType`s
// are both `CmpEq`. This is what called "implicit filter rules".
//
// For explicit rules, we've introduced a special URL syntax. A query param with
// name `filter` is passed in, and we will get a URL like this:
//
// ````
// /Model/[1]/[2]?filter=[3]
// ````
//
// Where [3] is a percent-encoded "filter string", with filter rules separated
// with `;`, like this:
//
// ````
// author,like,some_user;upload_time,gt,1519454311320
// ````
//
// Within the three comma-separated operands in each rule, the first is the Key,
// the second is the compare operator (one of `eq, ne, gt, lt, ge, le, like`),
// and the third is the Value.
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
