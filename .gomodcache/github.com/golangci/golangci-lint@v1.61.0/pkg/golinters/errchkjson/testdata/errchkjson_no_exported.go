//golangcitest:args -Eerrchkjson
//golangcitest:config_path testdata/errchkjson_no_exported.yml
package testdata

import (
	"encoding/json"
)

// JSONMarshalStructWithoutExportedFields contains a struct without exported fields.
func JSONMarshalStructWithoutExportedFields() {
	var withoutExportedFields struct {
		privateField            bool
		ExportedButOmittedField bool `json:"-"`
	}
	_, err := json.Marshal(withoutExportedFields) // want "Error argument passed to `encoding/json.Marshal` does not contain any exported field"
	_ = err
}

// JSONMarshalStructWithNestedStructWithoutExportedFields contains a struct without exported fields.
func JSONMarshalStructWithNestedStructWithoutExportedFields() {
	var withNestedStructWithoutExportedFields struct {
		ExportedStruct struct {
			privatField bool
		}
	}
	_, err := json.Marshal(withNestedStructWithoutExportedFields)
	_ = err
}
