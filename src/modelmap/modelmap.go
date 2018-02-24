package modelmap

import (
	"log"
	"sync"
	"errors"
	"unicode"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
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

func (r *Registry) BuildHandler(
	cookieSecret string,
) (http.Handler, error) {
	mux := http.NewServeMux()

	r.Lock()
	defer r.Unlock()

	sessionStore := sessions.NewCookieStore([]byte(cookieSecret))

	for k, _provider := range r.providers {
		urlName := transformModelNameToUrlRepr(k)
		provider := _provider
		mux.HandleFunc("/" + urlName + "/", func (w http.ResponseWriter, r *http.Request) {
			var filter map[string]FilterRule = make(map[string]FilterRule)

			session, _ := sessionStore.Get(r, "geekapk")

			rc := &RequestContext {
				Session: session,
			}

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
				ret = provider.Read(rc, filter)
				break
			case "PUT": // Update
				ret = provider.Update(rc, filter, defaultDeserializeFeed(string(body)))
				break
			case "POST": // Create
				ret = provider.Create(rc, defaultDeserializeFeed(string(body)))
				break
			case "DELETE": // Delete
				ret = provider.Delete(rc, filter)
				break
			default:
				panic("Unknown HTTP method")
			}

			session.Save(r, w)
			marshalAndWrite(w, ret)
		})
	}

	return context.ClearHandler(mux), nil
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
