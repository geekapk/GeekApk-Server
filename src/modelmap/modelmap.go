package modelmap

import (
	"log"
	"sync"
	"errors"
	"unicode"
	"strings"
	"net/url"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type Registry struct {
	sync.Mutex

	providers map[string]Provider
}

func NewRegistry() *Registry {
	return &Registry {
		providers: make(map[string]Provider),
	}
}

func (r *Registry) AddProvider(p Provider) error {
	name := p.GetName()
	nameChars := []rune(name)
	if len(nameChars) < 1 || !unicode.IsUpper(nameChars[0]) {
		return errors.New("Provider name must start with an upper case letter.")
	}

	r.Lock()
	defer r.Unlock()

	r.providers[name] = p

	log.Printf("Provider added: %s\n", name)
	return nil
}

func (r *Registry) RemoveProvider(name string) error {
	r.Lock()
	defer r.Unlock()

	delete(r.providers, name)
	return nil
}

func (r *Registry) BuildServeMux() (*http.ServeMux, error) {
	mux := http.NewServeMux()

	r.Lock()
	defer r.Unlock()

	for k, provider := range r.providers {
		urlName := transformModelNameToUrlRepr(k)
		mux.HandleFunc("/" + urlName + "/", func (w http.ResponseWriter, r *http.Request) {
			var filter map[string]FilterRule = make(map[string]FilterRule)

			parseImplicitFilterRules(filter, r.URL)

			err := parseFilterRules(filter, r.URL.Query().Get("filter"))
			if err != nil {
				marshalAndWrite(w, err.Error())
				return
			}

			var body []byte = nil

			if r.Method == "PUT" || r.Method == "POST" {
				body, err = ioutil.ReadAll(r.Body)
				if err != nil {
					panic(err)
				}
				defer r.Body.Close()
			}

			var ret interface{} = nil

			switch r.Method {
			case "GET": // Read
				ret = provider.Read(filter)
				break
			case "PUT": // Update
				ret = provider.Update(filter, defaultDeserializeFeed(string(body)))
				break
			case "POST": // Create
				ret = provider.Create(defaultDeserializeFeed(string(body)))
				break
			case "DELETE": // Delete
				ret = provider.Delete(filter)
				break
			default:
				panic("Unknown HTTP method")
			}

			marshalAndWrite(w, ret)
		})
	}

	return mux, nil
}

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

func marshalAndWrite(w http.ResponseWriter, data interface{}) {
	dataSer, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(dataSer)
}

func defaultDeserializeFeed(input string) func(out interface{}) error {
	inputBytes := []byte(input)
	return func(out interface{}) error {
		return json.Unmarshal(inputBytes, out)
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

func transformModelNameToUrlRepr(name string) string {
	ret := ""
	nameChars := []rune(name)

	for i, ch := range nameChars {
		if unicode.IsUpper(ch) {
			if i != 0 {
				ret += "_"
			}
			ret += string(unicode.ToLower(ch))
		} else {
			ret += string(ch)
		}
	}
	ret += "s"

	return ret
}
