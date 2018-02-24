package modelmap

import (
	"errors"
	"strings"
	"net/url"
)

func parseImplicitFilterRules(rules map[string]FilterRule, urlInfo *url.URL) {
	path := strings.Split(urlInfo.EscapedPath(), "/")

	if len(path) >= 3 {
		v, err := url.PathUnescape(path[2])
		if err == nil {
			rules["id"] = FilterRule {
				Key: "id",
				CompareType: CmpEq,
				Value: v,
			}
		}
	}
	if len(path) >= 4 {
		v, err := url.PathUnescape(path[3])
		if err == nil {
			rules["property"] = FilterRule {
				Key: "property",
				CompareType: CmpEq,
				Value: v,
			}
		}
	}
}

func parseFilterRules(rules map[string]FilterRule, input string) error {
	if len(input) == 0 {
		return nil
	}

	parts := strings.Split(input, ";")

	for _, p := range parts {
		operands := strings.Split(p, ",")
		if len(operands) != 3 {
			return errors.New("Expecting exactly 3 operands for filter rule")
		}
		cmpType := CmpUnknown
		switch operands[1] {
		case "eq":
			cmpType = CmpEq
			break
		case "ne":
			cmpType = CmpNe
			break
		case "gt":
			cmpType = CmpGt
			break
		case "ge":
			cmpType = CmpGe
			break
		case "lt":
			cmpType = CmpLt
			break
		case "le":
			cmpType = CmpLe
			break
		default:
			return errors.New("Expecting one of eq, ne, gt, ge, lt, le")
		}
		rules[operands[0]] = FilterRule {
			Key: operands[0],
			CompareType: cmpType,
			Value: operands[2],
		}
	}

	return nil
}
