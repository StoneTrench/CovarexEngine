package resources

import (
	"errors"
	"fmt"
	"regexp"
)

type ResourceError struct {
	Ident string
	Op    string
	Err   error
}

func (e *ResourceError) Error() string {
	return fmt.Sprintf("resource %s: %s: %s", e.Ident, e.Op, e.Err.Error())
}

var regex_invalid = regexp.MustCompile("[^a-zA-Z0-9._]")

// IsNameInvalid is used to check the incoming names for validity.
//
// A name can only contain [a-zA-Z0-9._] and cant have a length of 0.
func IsNameInvalid(name string) bool {
	return regex_invalid.MatchString(name) || len(name) == 0
}

// IsNameInvalid is used to check the incoming names for validity.
//
// A name can only contain [a-zA-Z0-9._] and cant have a length of 0.
func IsNameInvalidError(name string) error {
	errs := []error{}
	if regex_invalid.MatchString(name) {
		errs = append(errs, fmt.Errorf("name %q contains invalid characters", name))
	}
	if len(name) == 0 {
		errs = append(errs, fmt.Errorf("name %q has a length of zero", name))
	}
	return errors.Join(errs...)
}

// NameToIdentifier concatenates the two strings with a colon in between, returning a ResourceError if the names are not valid. See IsNameInvalidError.
func NameToIdentifier(mod, name string) (string, error) {
	ident := fmt.Sprintf("%s:%s", mod, name)
	errlist := []error{}
	if err := IsNameInvalidError(mod); err != nil {
		errlist = append(errlist, &ResourceError{Ident: ident, Op: "name_to_identifier", Err: err})
	}
	if err := IsNameInvalidError(name); err != nil {
		errlist = append(errlist, &ResourceError{Ident: ident, Op: "name_to_identifier", Err: err})
	}
	return ident, errors.Join(errlist...)
}
