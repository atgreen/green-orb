package rules

import (
	"errors"
	"go/types"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseRuleFile_IOError(t *testing.T) {
	r := iotest.ErrReader(errors.New("broken IO"))
	_, err := ParseRuleFile(r)
	assert.EqualError(t, err, "broken IO")
}

func TestParseFuncRule(t *testing.T) {
	testCases := []struct {
		name              string
		rule              string
		wantError         error
		wantPackageImport string
		wantRule          FuncRule
	}{
		{
			name:      "invalid-rule-missing-paren",
			rule:      "(*go.uber.org/zap/SugaredLogger.Debugw",
			wantError: ErrInvalidRule,
		},
		{
			name:      "invalid-rule-receiver-no-type",
			rule:      "(*go.uber.org/zap/SugaredLogger).Debugw",
			wantError: ErrInvalidRule,
		},
		{
			name:      "invalid-rule-just-import",
			rule:      "go.uber.org/zap",
			wantError: ErrInvalidRule,
		},
		{
			name:              "zap",
			rule:              "(*go.uber.org/zap.SugaredLogger).Debugw",
			wantPackageImport: "go.uber.org/zap",
			wantRule: FuncRule{
				IsReceiver:   true,
				ReceiverType: "*SugaredLogger",
				FuncName:     "Debugw",
			},
		},
		{
			name:              "klog-no-receiver",
			rule:              "k8s.io/klog/v2.InfoS",
			wantPackageImport: "k8s.io/klog/v2",
			wantRule: FuncRule{
				FuncName: "InfoS",
			},
		},
		{
			name:              "logr",
			rule:              "(github.com/go-logr/logr.Logger).Error",
			wantPackageImport: "github.com/go-logr/logr",
			wantRule: FuncRule{
				IsReceiver:   true,
				ReceiverType: "Logger",
				FuncName:     "Error",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotPackageImport, gotRule, err := ParseFuncRule(tc.rule)
			if tc.wantError != nil {
				assert.EqualError(t, err, tc.wantError.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.wantPackageImport, gotPackageImport)
				assert.Equal(t, tc.wantRule, gotRule)
			}
		})
	}
}

func TestReceiverTypeOf_InvalidType(t *testing.T) {
	t.Parallel()

	basicType := types.Universe.Lookup("byte").Type()
	assert.Equal(t, "", receiverTypeOf(basicType))
}
