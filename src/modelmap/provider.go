package modelmap

type Deserializer func (out interface{}) error

type Provider interface {
	GetName() string
	Create(
		createInfo Deserializer,
	) interface{}
	Read(
		filterRules map[string]FilterRule,
	) interface{}
	Update(
		filterRules map[string]FilterRule,
		updateInfo Deserializer,
	) interface{}
	Delete(
		filterRules map[string]FilterRule,
	) interface{}
}
