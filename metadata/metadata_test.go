package metadata

import (
	"encoding/xml"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMetadata(t *testing.T) {
	bs, err := ioutil.ReadFile("./testdata/idp.test-metadata.xml")
	require.NoError(t, err, "error reading metadata test file")

	var es Entities

	require.NoError(t, xml.Unmarshal(bs, &es))

	require.Len(t, es, 3, "wrong number of entities returned")

	expectationsTable := []struct{ eid, t string }{
		{eid: "https://portal.astuart.co/uPortal", t: TypeSP},
		{eid: "https://saml2.test.astuart.co/sso/saml2", t: TypeSP},
		{eid: "https://idp.astuart.co/idp/shibboleth", t: TypeIDP},
	}

	for i := range expectationsTable {
		require.Equal(t, expectationsTable[i].eid, es[i].EntityID, "wrong entityID for item %d", i)
		require.Equal(t, expectationsTable[i].t, es[i].Type, "wrong type for item %d", i)
	}

	f := es[0]
	require.Len(t, f.Keys, 2, "wrong number of keys for first item")
	require.Equal(t, f.Keys[0].Usage, "signing")
	require.Equal(t, f.Keys[1].Usage, "encryption")

	require.Len(t, f.NameIDFormats, 5, "wrong number of nameidformats")
	require.Len(t, f.LogoutServices, 2)
	require.Equal(t, f.LogoutServices[0].Binding, "urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST")
	require.Equal(t, f.LogoutServices[0].Location, "https://portal.ccctcportal.org/uPortal/saml/SingleLogout")

	require.Len(t, f.Consumers, 1)
	require.Equal(t, true, f.Consumers[0].Default, "first consumer was not a default")
}

func TestKey(t *testing.T) {
	bs, err := ioutil.ReadFile("./testdata/idp.test-metadata.xml")
	require.NoError(t, err, "error reading metadata test file")

	var es Entities
	require.NoError(t, xml.Unmarshal(bs, &es))

	c, err := es[0].Keys[0].Cert()
	require.NoError(t, err)

	require.Len(t, c.Certificate, 1)
	require.Equal(t, "apollo", c.Leaf.Subject.CommonName)
}
