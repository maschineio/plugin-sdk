package sdk

import (
	"encoding/json"
)

// RequestBuilder helps build ExecuteRequest objects
type RequestBuilder struct {
	resource    string
	input       []byte
	parameters  map[string][]byte
	credentials map[string]string
	context     map[string]string
}

// NewRequestBuilder creates a new request builder
func NewRequestBuilder(resource string) *RequestBuilder {
	return &RequestBuilder{
		resource:    resource,
		parameters:  make(map[string][]byte),
		credentials: make(map[string]string),
		context:     make(map[string]string),
	}
}

// SetInput sets the input data
func (b *RequestBuilder) SetInput(input []byte) *RequestBuilder {
	b.input = input
	return b
}

// SetInputJSON sets the input data from a JSON-serializable object
func (b *RequestBuilder) SetInputJSON(input any) *RequestBuilder {
	data, _ := json.Marshal(input)
	b.input = data
	return b
}

// AddParameter adds a parameter with JSON encoding
func (b *RequestBuilder) AddParameter(key string, value any) *RequestBuilder {
	data, _ := json.Marshal(value)
	b.parameters[key] = data
	return b
}

// AddStringParameter adds a string parameter
func (b *RequestBuilder) AddStringParameter(key, value string) *RequestBuilder {
	return b.AddParameter(key, value)
}

// AddIntParameter adds an int parameter
func (b *RequestBuilder) AddIntParameter(key string, value int) *RequestBuilder {
	return b.AddParameter(key, value)
}

// AddBoolParameter adds a bool parameter
func (b *RequestBuilder) AddBoolParameter(key string, value bool) *RequestBuilder {
	return b.AddParameter(key, value)
}

// AddCredential adds a credential
func (b *RequestBuilder) AddCredential(key, value string) *RequestBuilder {
	b.credentials[key] = value
	return b
}

// AddContext adds a context value
func (b *RequestBuilder) AddContext(key, value string) *RequestBuilder {
	b.context[key] = value
	return b
}

// Build creates the ExecuteRequest
func (b *RequestBuilder) Build() *ExecuteRequest {
	return &ExecuteRequest{
		Resource:    b.resource,
		Input:       b.input,
		Parameters:  b.parameters,
		Credentials: b.credentials,
		Context:     b.context,
	}
}

// BuildTyped creates a TypedExecuteRequest
func (b *RequestBuilder) BuildTyped() (*TypedExecuteRequest, error) {
	req := b.Build()
	return NewTypedExecuteRequest(req)
}