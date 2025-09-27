package garif

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithLineColumn(t *testing.T) {
	l := NewLocation()
	require.NotNil(t, l)
	assert.Nil(t, l.PhysicalLocation)

	const line = 10
	const column = 20
	l.WithLineColumn(line, column)

	require.NotNil(t, l.PhysicalLocation)
	require.NotNil(t, l.PhysicalLocation.Region)

	assert.Equal(t, line, l.PhysicalLocation.Region.StartLine)
	assert.Equal(t, column, l.PhysicalLocation.Region.StartColumn)

	l.WithLineColumn(0, 0)

	assert.Equal(t, 0, l.PhysicalLocation.Region.StartLine)
	assert.Equal(t, 0, l.PhysicalLocation.Region.StartColumn)
}

func TestWithURI(t *testing.T) {
	l := NewLocation()
	require.NotNil(t, l)
	assert.Nil(t, l.PhysicalLocation)

	const uri = "URI4TEST"
	l.WithURI(uri)

	require.NotNil(t, l.PhysicalLocation)
	require.NotNil(t, l.PhysicalLocation.ArtifactLocation)

	assert.Equal(t, uri, l.PhysicalLocation.ArtifactLocation.Uri)

	l.WithURI("")

	assert.Equal(t, "", l.PhysicalLocation.ArtifactLocation.Uri)
}

func TestWithKeyValue(t *testing.T) {
	b := *NewPropertyBag()

	const key = "a-key"
	const value = 0.0

	b.WithKeyValue(key, value)

	v, ok := b[key]
	assert.True(t, ok)
	assert.Equal(t, value, v)

	const key2 = "another-key"
	const value2 = "0.0"

	b.WithKeyValue(key2, value2)

	v, ok = b[key]
	assert.True(t, ok)
	assert.Equal(t, value, v)

	v, ok = b[key2]
	assert.True(t, ok)
	assert.Equal(t, value2, v)
}

func TestWithWithHelpUri(t *testing.T) {
	r := NewReportingDescriptor("foo")
	require.NotNil(t, r)

	const uri = "URI4TEST"
	r.WithHelpUri(uri)

	assert.Equal(t, uri, r.HelpUri)

	r.WithHelpUri("")

	assert.Equal(t, "", r.HelpUri)
}

func TestWithArtifactsURIs(t *testing.T) {
	r := NewRun(nil)
	require.NotNil(t, r)

	var uris = []string{"a", "b", "c"}

	r.WithArtifactsURIs(uris...)

	artifacts := r.Artifacts
	require.NotNil(t, artifacts)

	require.Equal(t, len(uris), len(artifacts))

	for i := 0; i < len(uris); i++ {
		assert.NotNil(t, artifacts[i].Location.Uri)
		aURI := artifacts[i].Location.Uri
		assert.Equal(t, uris[i], aURI)
	}
}

func TestWithWithResult(t *testing.T) {
	r := NewRun(nil)
	require.NotNil(t, r)

	const line = 10
	const column = 20
	const uri = "uri"
	const ruleId = "id"
	const message = "msg"

	r.WithResult(ruleId, message, uri, line, column)

	results := r.Results
	require.Equal(t, 1, len(results))

	theResult := results[0]
	assert.Equal(t, ruleId, theResult.RuleId)
	require.NotNil(t, theResult.Locations)
	require.Equal(t, 1, len(theResult.Locations))
	require.NotNil(t, theResult.Locations[0].PhysicalLocation)
	assert.Equal(t, line, theResult.Locations[0].PhysicalLocation.Region.StartLine)
	assert.Equal(t, column, theResult.Locations[0].PhysicalLocation.Region.StartColumn)
	assert.Equal(t, uri, theResult.Locations[0].PhysicalLocation.ArtifactLocation.Uri)
	require.NotNil(t, theResult.Message)
	assert.Equal(t, message, theResult.Message.Text)
}
