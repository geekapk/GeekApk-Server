package modelmap

type Deserializer func (out interface{}) error

type Provider interface {
	GetName() string
	Create(
		rc *RequestContext,
		createInfo Deserializer,
	) interface{}
	Read(
		rc *RequestContext,
		filterRules map[string]FilterRule,
	) interface{}
	Update(
		rc *RequestContext,
		filterRules map[string]FilterRule,
		updateInfo Deserializer,
	) interface{}
	Delete(
		rc *RequestContext,
		filterRules map[string]FilterRule,
	) interface{}
}
