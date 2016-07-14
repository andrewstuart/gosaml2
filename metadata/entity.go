package metadata

import (
	"encoding/xml"
	"fmt"
	"io"

	"astuart.co/astuart.co/gosaml2"
)

const (
	// TypeIDP is the well-known type for IDP entities
	TypeIDP = "IDP"
	// TypeSP is the well-known type for SP entities
	TypeSP = "SP"

	descIDP    = "IDPSSODescriptor"
	descSP     = "SPSSODescriptor"
	descKey    = "KeyDescriptor"
	descEnts   = "EntitiesDescriptor"
	descEnt    = "EntityDescriptor"
	descLogout = "SingleLogoutService"
	descNID    = "NameIDFormat"
	descACS    = "AssertionConsumerService"
)

type Entity struct {
	ID       string `xml:",attr"`
	EntityID string `xml:"entityID,attr"`
	Type     string

	Keys []saml.EncryptedKey
	// LogoutServices []LogoutService
	NameIDFormats []string
	// Consumers     []AssertionConsumer
}

func (e *Entity) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	*e = Entity{
		Keys:          []saml.EncryptedKey{},
		NameIDFormats: []string{},
	}

	// StartElement should always be an EntityDescriptor, thus having ID and
	// entityID attrs
	err := e.parseDescriptor(start)
	if err != nil {
		return err
	}

	// Start going through the tokens and parse the data
	for {
		t, err := d.Token()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("error decoding entity xml: %s", err)
		}

		switch t := t.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case descSP:
				e.Type = TypeSP
			case descIDP:
				e.Type = TypeIDP
			case descKey:
				var k saml.EncryptedKey
				err = d.DecodeElement(&k, &t)
				if err != nil {
					return err
				}
				e.Keys = append(e.Keys, k)
			case descNID:
				var nid string
				err = d.DecodeElement(&nid, &t)
				if err != nil {
					return err
				}
				e.NameIDFormats = append(e.NameIDFormats, nid)
			}

		}

		if err != nil {
			return fmt.Errorf("error decoding entity xml: %s", err)
		}
	}
}

func (e *Entity) parseDescriptor(desc xml.StartElement) error {
	for _, attr := range desc.Attr {
		switch attr.Name.Local {
		case "ID":
			e.ID = attr.Value
		case "entityID":
			e.EntityID = attr.Value
		}
	}
	return nil
}
