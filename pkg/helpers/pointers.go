// Package helpers provides simple helper to deal with various nil events and slice objects
package helpers

// IsStringPtrNilOrEmtpy returns true either if the *string is nil or the value is empty
func IsStringPtrNilOrEmtpy(ptr *string) bool {
	return ptr == nil || *ptr == ""
}

// IsBoolPtrTrue returns true if the *bool is not nil and the value is true
func IsBoolPtrTrue(ptr *bool) bool {
	return ptr != nil && *ptr == true
}
