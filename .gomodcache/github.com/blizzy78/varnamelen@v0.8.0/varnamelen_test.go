package varnamelen

import (
	"os"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestVarNameLen_Run_Warnings(t *testing.T) {
	run(t, "warnings", nil)
}

func TestVarNameLen_Run_Warnings_Return(t *testing.T) {
	run(t, "warningsreturn", map[string]string{
		"checkReturn": "true",
	})
}

func TestVarNameLen_Run_Warnings_Receiver(t *testing.T) {
	run(t, "warningsreceiver", map[string]string{
		"checkReceiver": "true",
	})
}

func TestVarNameLen_Run_NameLen(t *testing.T) {
	run(t, "namelen", nil)
}

func TestVarNameLen_Run_NameLen_Return(t *testing.T) {
	run(t, "namelenreturn", map[string]string{
		"checkReturn": "true",
	})
}

func TestVarNameLen_Run_Distance(t *testing.T) {
	run(t, "distance", nil)
}

func TestVarNameLen_Run_Distance_Return(t *testing.T) {
	run(t, "distancereturn", map[string]string{
		"checkReturn": "true",
	})
}

func TestVarNameLen_Run_Distance_Receiver(t *testing.T) {
	run(t, "distancereceiver", map[string]string{
		"checkReceiver": "true",
	})
}

func TestVarNameLen_Run_IgnoreNames(t *testing.T) {
	run(t, "ignorenames", map[string]string{
		"ignoreNames": "i, I",
	})
}

func TestVarNameLen_Run_IgnoreNames_Return(t *testing.T) {
	run(t, "ignorenamesreturn", map[string]string{
		"ignoreNames": "i",
		"checkReturn": "true",
	})
}

func TestVarNameLen_Run_IgnoreNames_Receiver(t *testing.T) {
	run(t, "ignorenamesreceiver", map[string]string{
		"ignoreNames":   "f",
		"checkReceiver": "true",
	})
}

func TestVarNameLen_Run_IgnoreDecls(t *testing.T) {
	run(t, "ignoredecls", map[string]string{
		"ignoreDecls": "i int, i map[string]string, i *bytes.Buffer, i *strs.Builder, const I",
	})
}

func TestVarNameLen_Run_IgnoreDecls_Return(t *testing.T) {
	run(t, "ignoredeclsreturn", map[string]string{
		"ignoreDecls": "i int",
		"checkReturn": "true",
	})
}

func TestVarNameLen_Run_IgnoreDecls_Receiver(t *testing.T) {
	run(t, "ignoredeclsreceiver", map[string]string{
		"ignoreDecls":   "f foo, f *foo",
		"checkReceiver": "true",
	})
}

func TestVarNameLen_Run_IgnoreTypeAssertOk(t *testing.T) {
	run(t, "typeassert", map[string]string{
		"ignoreTypeAssertOk": "true",
	})
}

func TestVarNameLen_Run_IgnoreChanRecvOk(t *testing.T) {
	run(t, "chanrecv", map[string]string{
		"ignoreChanRecvOk": "true",
	})
}

func TestVarNameLen_Run_IgnoreMapIndexOk(t *testing.T) {
	run(t, "mapindex", map[string]string{
		"ignoreMapIndexOk": "true",
	})
}

func TestVarNameLen_Run_Conventional(t *testing.T) {
	run(t, "conventional", nil)
}

func TestVarNameLen_Run_FalseNegative_TypeAssertOk(t *testing.T) {
	run(t, "falsenegativetypeassert", map[string]string{
		"ignoreTypeAssertOk": "true",
	})
}

func TestVarNameLen_Run_FalseNegative_ChanRecvOk(t *testing.T) {
	run(t, "falsenegativechanrecv", map[string]string{
		"ignoreChanRecvOk": "true",
	})
}

func TestVarNameLen_Run_FalseNegative_MapIndexOk(t *testing.T) {
	run(t, "falsenegativemapindex", map[string]string{
		"ignoreMapIndexOk": "true",
	})
}

func run(t *testing.T, pkg string, flags map[string]string) {
	t.Helper()

	analyzer := NewAnalyzer()

	for k, v := range flags {
		_ = analyzer.Flags.Set(k, v)
	}

	wd, _ := os.Getwd()
	analysistest.Run(t, wd+"/testdata", analyzer, pkg)
}
