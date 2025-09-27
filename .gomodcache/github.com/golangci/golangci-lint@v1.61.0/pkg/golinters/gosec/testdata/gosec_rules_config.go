//golangcitest:args -Egosec
//golangcitest:config_path testdata/gosec.yml
package testdata

import "io/ioutil"

const gosecToken = "62ebc7a03d6ca24dca1258fd4b48462f6fed1545"
const gosecSimple = "62ebc7a03d6ca24dca1258fd4b48462f6fed1545" // want "G101: Potential hardcoded credentials"

func gosecCustom() {
	ioutil.WriteFile("filename", []byte("test"), 0755) // want "G306: Expect WriteFile permissions to be 0666 or less"
}
