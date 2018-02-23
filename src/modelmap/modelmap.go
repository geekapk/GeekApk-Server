package modelmap

import (
	"log"
	"sync"
	"errors"
	"unicode"
	"strings"
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
		urlPrefix := "/" + urlName + "/"
		mux.HandleFunc(urlPrefix + "new", func (w http.ResponseWriter, r *http.Request) {
			rawCreateInfoBytes, err := ioutil.ReadAll(r.Body)
			if err != nil {
				panic(err)
			}
			defer r.Body.Close()

			ret := provider.Create(defaultDeserializeFeed(string(rawCreateInfoBytes)))
			marshalAndWrite(w, ret)
		})
		mux.HandleFunc(urlPrefix + "get", func (w http.ResponseWriter, r *http.Request) {
			filter, err := parseFilterRules(r.URL.Query().Get("filter"))
			if err != nil {
				marshalAndWrite(w, err.Error())
				return
			}

			ret := provider.Read(filter)
			marshalAndWrite(w, ret)
		})
		mux.HandleFunc(urlPrefix + "update", func (w http.ResponseWriter, r *http.Request) {
			rawUpdateInfoBytes, err := ioutil.ReadAll(r.Body)
			if err != nil {
				panic(err)
			}
			defer r.Body.Close()

			filter, err := parseFilterRules(r.URL.Query().Get("filter"))
			if err != nil {
				marshalAndWrite(w, err.Error())
				return
			}

			ret := provider.Update(filter, defaultDeserializeFeed(string(rawUpdateInfoBytes)))
			marshalAndWrite(w, ret)
		})
		mux.HandleFunc(urlPrefix + "remove", func (w http.ResponseWriter, r *http.Request) {
			filter, err := parseFilterRules(r.URL.Query().Get("filter"))
			if err != nil {
				marshalAndWrite(w, err.Error())
				return
			}

			ret := provider.Delete(filter)
			marshalAndWrite(w, ret)
		})
	}

	return mux, nil
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

func parseFilterRules(input string) (map[string]FilterRule, error) {
	ret := make(map[string]FilterRule)

	if len(input) == 0 {
		return ret, nil
	}

	parts := strings.Split(input, ";")

	for _, p := range parts {
		operands := strings.Split(p, ",")
		if len(operands) != 3 {
			return nil, errors.New("Expecting exactly 3 operands for filter rule")
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
			return nil, errors.New("Expecting one of eq, ne, gt, ge, lt, le")
		}
		ret[operands[0]] = FilterRule {
			Key: operands[0],
			CompareType: cmpType,
			Value: operands[2],
		}
	}

	return ret, nil
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
