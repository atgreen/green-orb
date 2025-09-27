package global

func dontImportPackage() {
	header.Get("Test-HEader") // want `non-canonical header "Test-HEader", instead use: "Test-Header"`
}
