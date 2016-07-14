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

	xml.Unmarshal(bs, &es)

	require.Len(t, es, 3, "wrong number of entities returned")

	expectedIDs := []string{
		"https://portal.astuart.co/uPortal",
		"https://saml2.test.astuart.co/sso/saml2",
		"https://idp.astuart.co/idp/shibboleth",
	}

	for i := range expectedIDs {
		require.Equal(t, expectedIDs[i], es[i].EntityID, "wrong entityID for item %d", i)
	}
}
