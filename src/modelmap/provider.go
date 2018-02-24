package modelmap

type Deserializer func (out interface{}) error

// A Provider is a Model used to serve requests matching some rules.
//
// Note that the Model here does not correspond perfectly to the "Model"
// in MVC. The Model here directly serves HTTP requests.
//
// For complex logic that should not be processed in a Model directly,
// it should be abstracted out as a Service.
type Provider interface {
	// Gets the name of the Model.
	//
	// The name returned here should be a constant string and will be
	// registered and cached while being added to a Registry instance.
	//
	// The name MUST start with a upper-case letter.
	//
	// A model's name has the following uses:
	//
	// a) Used in the first level routing path with its plural form.
	//    For example, the `Account` model will be accessible from `/accounts/`.
	// b) Debug info
	GetName() string

	// Creates a resource related to the model.
	//
	// This corresponds to the `POST` HTTP method.`
	//
	// Request body can be deserialized by calling `loadCreateInfo`, which is
	// passed in as an argument.
	//
	// The return value is the response and will be serialized to JSON.
	Create(
		rc *RequestContext,
		loadCreateInfo Deserializer,
	) interface{}

	// Gets a resource related to the model.
	//
	// This corresponds to the `GET` HTTP method.
	//
	// See `filter.go` for `filterRules`.
	Read(
		rc *RequestContext,
		filterRules map[string]FilterRule,
	) interface{}

	// Updates a resource related to the model.
	//
	// This corresponds to the `PUT` HTTP method.
	Update(
		rc *RequestContext,
		filterRules map[string]FilterRule,
		loadUpdateInfo Deserializer,
	) interface{}

	// Deletes a resource related to the model.
	//
	// This corresponds to the `DELETE` HTTP method.
	Delete(
		rc *RequestContext,
		filterRules map[string]FilterRule,
	) interface{}
}
