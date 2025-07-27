package manifest

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationErrors is a collection of validation errors
type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	var msgs []string
	for _, err := range e {
		msgs = append(msgs, err.Error())
	}
	return "validation failed:\n" + strings.Join(msgs, "\n")
}

// Validate validates a plugin manifest
func (m *PluginManifest) Validate() error {
	var errors ValidationErrors

	// Validate manifest version
	if m.ManifestVersion == "" {
		errors = append(errors, ValidationError{
			Field:   "manifestVersion",
			Message: "is required",
		})
	} else if m.ManifestVersion != "2.0" {
		errors = append(errors, ValidationError{
			Field:   "manifestVersion",
			Message: fmt.Sprintf("unsupported version '%s', expected '2.0'", m.ManifestVersion),
		})
	}

	// Validate plugin info
	if err := m.Plugin.Validate(); err != nil {
		if ve, ok := err.(ValidationErrors); ok {
			errors = append(errors, ve...)
		} else {
			errors = append(errors, ValidationError{Field: "plugin", Message: err.Error()})
		}
	}

	// Validate runtime
	if err := m.Runtime.Validate(); err != nil {
		if ve, ok := err.(ValidationErrors); ok {
			errors = append(errors, ve...)
		} else {
			errors = append(errors, ValidationError{Field: "runtime", Message: err.Error()})
		}
	}

	// Validate requirements
	if err := m.Requirements.Validate(); err != nil {
		if ve, ok := err.(ValidationErrors); ok {
			errors = append(errors, ve...)
		} else {
			errors = append(errors, ValidationError{Field: "requirements", Message: err.Error()})
		}
	}

	// Validate resources
	if len(m.Resources) == 0 {
		errors = append(errors, ValidationError{
			Field:   "resources",
			Message: "at least one resource must be defined",
		})
	}
	for i, r := range m.Resources {
		if err := r.Validate(); err != nil {
			if ve, ok := err.(ValidationErrors); ok {
				for _, e := range ve {
					e.Field = fmt.Sprintf("resources[%d].%s", i, e.Field)
					errors = append(errors, e)
				}
			} else {
				errors = append(errors, ValidationError{
					Field:   fmt.Sprintf("resources[%d]", i),
					Message: err.Error(),
				})
			}
		}
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}

// Validate validates plugin info
func (p *PluginInfo) Validate() error {
	var errors ValidationErrors

	// Required fields
	if p.ID == "" {
		errors = append(errors, ValidationError{Field: "plugin.id", Message: "is required"})
	} else if !isValidPluginID(p.ID) {
		errors = append(errors, ValidationError{
			Field:   "plugin.id",
			Message: "must follow reverse domain notation (e.g., io.maschine.plugins.mail)",
		})
	}

	if p.Name == "" {
		errors = append(errors, ValidationError{Field: "plugin.name", Message: "is required"})
	}
	if p.DisplayName == "" {
		errors = append(errors, ValidationError{Field: "plugin.displayName", Message: "is required"})
	}
	if p.Version == "" {
		errors = append(errors, ValidationError{Field: "plugin.version", Message: "is required"})
	}
	if p.Description == "" {
		errors = append(errors, ValidationError{Field: "plugin.description", Message: "is required"})
	}
	if p.Category == "" {
		errors = append(errors, ValidationError{Field: "plugin.category", Message: "is required"})
	}
	if p.License == "" {
		errors = append(errors, ValidationError{Field: "plugin.license", Message: "is required"})
	}

	// Validate URLs
	if p.Homepage != "" && !isValidURL(p.Homepage) {
		errors = append(errors, ValidationError{Field: "plugin.homepage", Message: "must be a valid URL"})
	}
	if p.Repository.URL != "" && !isValidURL(p.Repository.URL) {
		errors = append(errors, ValidationError{Field: "plugin.repository.url", Message: "must be a valid URL"})
	}
	if p.Bugs.URL != "" && !isValidURL(p.Bugs.URL) {
		errors = append(errors, ValidationError{Field: "plugin.bugs.url", Message: "must be a valid URL"})
	}

	// Validate author
	if p.Author.Name == "" {
		errors = append(errors, ValidationError{Field: "plugin.author.name", Message: "is required"})
	}
	if p.Author.Email == "" {
		errors = append(errors, ValidationError{Field: "plugin.author.email", Message: "is required"})
	} else if !isValidEmail(p.Author.Email) {
		errors = append(errors, ValidationError{Field: "plugin.author.email", Message: "must be a valid email"})
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}

// Validate validates runtime info
func (r *RuntimeInfo) Validate() error {
	var errors ValidationErrors

	if r.Type == "" {
		errors = append(errors, ValidationError{Field: "runtime.type", Message: "is required"})
	}
	if r.Executable.Path == "" {
		errors = append(errors, ValidationError{Field: "runtime.executable.path", Message: "is required"})
	}
	if r.Protocol == "" {
		errors = append(errors, ValidationError{Field: "runtime.protocol", Message: "is required"})
	}
	if r.ProtocolVersion == "" {
		errors = append(errors, ValidationError{Field: "runtime.protocolVersion", Message: "is required"})
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}

// Validate validates requirements
func (r *Requirements) Validate() error {
	var errors ValidationErrors

	if r.MaschineVersion == "" {
		errors = append(errors, ValidationError{Field: "requirements.maschineVersion", Message: "is required"})
	}
	if len(r.OS) == 0 {
		errors = append(errors, ValidationError{Field: "requirements.os", Message: "at least one OS must be specified"})
	}
	if len(r.Arch) == 0 {
		errors = append(errors, ValidationError{Field: "requirements.arch", Message: "at least one architecture must be specified"})
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}

// Validate validates a resource definition
func (r *ResourceDef) Validate() error {
	var errors ValidationErrors

	if r.Type == "" {
		errors = append(errors, ValidationError{Field: "type", Message: "is required"})
	} else if !isValidMRN(r.Type) {
		errors = append(errors, ValidationError{Field: "type", Message: "must be a valid MRN (e.g., mrn:mail:smtp:send)"})
	}

	if r.Name == "" {
		errors = append(errors, ValidationError{Field: "name", Message: "is required"})
	}
	if r.Description == "" {
		errors = append(errors, ValidationError{Field: "description", Message: "is required"})
	}
	if r.Category == "" {
		errors = append(errors, ValidationError{Field: "category", Message: "is required"})
	}

	// Validate parameters
	paramNames := make(map[string]bool)
	for i, p := range r.Parameters {
		if p.Name == "" {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("parameters[%d].name", i),
				Message: "is required",
			})
		} else if paramNames[p.Name] {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("parameters[%d].name", i),
				Message: fmt.Sprintf("duplicate parameter name '%s'", p.Name),
			})
		}
		paramNames[p.Name] = true

		if p.Type == "" {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("parameters[%d].type", i),
				Message: "is required",
			})
		}
		if p.Description == "" {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("parameters[%d].description", i),
				Message: "is required",
			})
		}
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}

// Helper functions

func isValidPluginID(id string) bool {
	// Should follow reverse domain notation
	match, _ := regexp.MatchString(`^[a-z][a-z0-9]*(\.[a-z][a-z0-9]*)+$`, id)
	return match
}

func isValidMRN(mrn string) bool {
	// MRN format: mrn:service:resource or mrn:service:resource:action etc.
	// Allow alphanumeric characters in each part, minimum 2 parts after "mrn:"
	match, _ := regexp.MatchString(`^mrn:[a-z][a-z0-9]*(:[a-z][a-z0-9]*)+$`, mrn)
	return match
}

func isValidURL(s string) bool {
	u, err := url.Parse(s)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func isValidEmail(email string) bool {
	match, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, email)
	return match
}