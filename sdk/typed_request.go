package sdk

import (
	"encoding/json"
	"fmt"
	"maschine.io/core/params"
)

// TypedExecuteRequest wraps ExecuteRequest with typed parameter access
type TypedExecuteRequest struct {
	*ExecuteRequest
	params *params.Parameter
}

// NewTypedExecuteRequest creates a typed request from a raw request
func NewTypedExecuteRequest(req *ExecuteRequest) (*TypedExecuteRequest, error) {
	// Convert Parameters map[string][]byte to map[string]any
	paramsMap := make(map[string]any)
	for key, data := range req.Parameters {
		var value any
		if err := json.Unmarshal(data, &value); err != nil {
			paramsMap[key] = string(data)
		} else {
			paramsMap[key] = value
		}
	}
	
	return &TypedExecuteRequest{
		ExecuteRequest: req,
		params:         params.NewParameter(&paramsMap),
	}, nil
}

// GetParams returns the typed parameter accessor
func (r *TypedExecuteRequest) GetParams() *params.Parameter {
	return r.params
}

// GetParam retrieves a typed parameter
func (r *TypedExecuteRequest) GetParam(key string) (any, error) {
	return params.GetParam[any](r.params, key)
}

// GetStringParam retrieves a string parameter
func (r *TypedExecuteRequest) GetStringParam(key string) (string, error) {
	return params.GetParam[string](r.params, key)
}

// GetIntParam retrieves an int parameter (handles float64 from JSON)
func (r *TypedExecuteRequest) GetIntParam(key string) (int, error) {
	val := r.params.Get(key)
	if val == nil {
		return 0, fmt.Errorf("'%v' parameter expected", key)
	}
	
	switch v := val.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		// JSON numbers are unmarshaled as float64
		return int(v), nil
	case float32:
		return int(v), nil
	default:
		return 0, fmt.Errorf("'%v' parameter must be a number, got %T", key, val)
	}
}

// GetBoolParam retrieves a bool parameter
func (r *TypedExecuteRequest) GetBoolParam(key string) (bool, error) {
	return params.GetParam[bool](r.params, key)
}

// GetSliceParam retrieves a slice parameter
func (r *TypedExecuteRequest) GetSliceParam(key string) ([]any, error) {
	return params.GetParam[[]any](r.params, key)
}

// GetParamWithDefault retrieves a typed parameter with a default value
func (r *TypedExecuteRequest) GetParamWithDefault(key string, defaultValue any) (any, error) {
	return params.GetParamDefault[any](r.params, key, defaultValue)
}

// GetStringParamWithDefault retrieves a string parameter with a default value
func (r *TypedExecuteRequest) GetStringParamWithDefault(key string, defaultValue string) (string, error) {
	return params.GetParamDefault[string](r.params, key, defaultValue)
}

// GetIntParamWithDefault retrieves an int parameter with a default value (handles float64 from JSON)
func (r *TypedExecuteRequest) GetIntParamWithDefault(key string, defaultValue int) (int, error) {
	val := r.params.Get(key)
	if val == nil {
		return defaultValue, nil
	}
	
	switch v := val.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		// JSON numbers are unmarshaled as float64
		return int(v), nil
	case float32:
		return int(v), nil
	default:
		return 0, fmt.Errorf("'%v' parameter must be a number, got %T", key, val)
	}
}

// GetBoolParamWithDefault retrieves a bool parameter with a default value
func (r *TypedExecuteRequest) GetBoolParamWithDefault(key string, defaultValue bool) (bool, error) {
	return params.GetParamDefault[bool](r.params, key, defaultValue)
}

// GetSliceParamWithDefault retrieves a slice parameter with a default value
func (r *TypedExecuteRequest) GetSliceParamWithDefault(key string, defaultValue []any) ([]any, error) {
	return params.GetParamDefault[[]any](r.params, key, defaultValue)
}

// GetInt64Param retrieves an int64 parameter (handles float64 from JSON)
func (r *TypedExecuteRequest) GetInt64Param(key string) (int64, error) {
	val := r.params.Get(key)
	if val == nil {
		return 0, fmt.Errorf("'%v' parameter expected", key)
	}
	
	switch v := val.(type) {
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	case float64:
		// JSON numbers are unmarshaled as float64
		return int64(v), nil
	case float32:
		return int64(v), nil
	default:
		return 0, fmt.Errorf("'%v' parameter must be a number, got %T", key, val)
	}
}

// GetInt64ParamWithDefault retrieves an int64 parameter with a default value
func (r *TypedExecuteRequest) GetInt64ParamWithDefault(key string, defaultValue int64) (int64, error) {
	val := r.params.Get(key)
	if val == nil {
		return defaultValue, nil
	}
	
	switch v := val.(type) {
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	case float64:
		// JSON numbers are unmarshaled as float64
		return int64(v), nil
	case float32:
		return int64(v), nil
	default:
		return 0, fmt.Errorf("'%v' parameter must be a number, got %T", key, val)
	}
}

// GetFloat64Param retrieves a float64 parameter
func (r *TypedExecuteRequest) GetFloat64Param(key string) (float64, error) {
	val := r.params.Get(key)
	if val == nil {
		return 0, fmt.Errorf("'%v' parameter expected", key)
	}
	
	switch v := val.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	default:
		return 0, fmt.Errorf("'%v' parameter must be a number, got %T", key, val)
	}
}

// GetFloat64ParamWithDefault retrieves a float64 parameter with a default value
func (r *TypedExecuteRequest) GetFloat64ParamWithDefault(key string, defaultValue float64) (float64, error) {
	val := r.params.Get(key)
	if val == nil {
		return defaultValue, nil
	}
	
	switch v := val.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	default:
		return 0, fmt.Errorf("'%v' parameter must be a number, got %T", key, val)
	}
}