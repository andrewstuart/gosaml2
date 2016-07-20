package metadata

import (
	"encoding/xml"
	"fmt"
	"io"

	"astuart.co/gosaml2"
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
	descSSO    = "SingleSignOnService"
)

// Entity represents a described entity within SAML2.0 Metadata
type Entity struct {
	ID string `xml:",attr"`
	// EntityID is the unique id for the specified entity. This is most often a
	// URI.
	EntityID string `xml:"entityID,attr"`
	// Type indicates whether this entity should be treated as an SP or IDP
	Type string

	Keys               []saml2.Key
	NameIDFormats      []string
	AssertionConsumers []AssertionConsumer
	SignAuthnReq       bool
	LogoutServices     []LogoutService `xml:"SingleLogoutService"`
	SSOServices        []SSOService
	Consumers          []AssertionConsumer
}

type endpoint struct {
	Binding  string `xml:",attr"`
	Location string `xml:",attr"`
}

// LogoutService is a combination of Logout URL and the SAML binding required
// for exercising the logout mechanism.
type LogoutService struct {
	endpoint
}

// SSOService is an endpoint/binding for Single Signon
type SSOService struct {
	endpoint
}

// AssertionConsumer is a type that represents a SAML2 service provider
// endpoint
type AssertionConsumer struct {
	endpoint
	// Default indicates whether or not this assertion consumer should be treated
	// as the default
	Default bool `xml:"isDefault,attr"`
}

// UnmarshalXML implements xml.Unmarshaler for Entity
func (e *Entity) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*e = Entity{
		Keys:           []saml2.Key{},
		NameIDFormats:  []string{},
		LogoutServices: []LogoutService{},
		SSOServices:    []SSOService{},
		Consumers:      []AssertionConsumer{},
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
				//SP Descriptor
				e.Type = TypeSP
				for _, att := range t.Attr {
					switch att.Name.Local {
					case "AuthnRequestsSigned":
						e.SignAuthnReq = att.Value == "true"
					}
				}
			case descIDP:
				e.Type = TypeIDP
			case descKey:
				var k saml2.Key
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
			case descLogout:
				var ls LogoutService
				err = d.DecodeElement(&ls, &t)
				if err != nil {
					return err
				}
				e.LogoutServices = append(e.LogoutServices, ls)
			case descSSO:
				var ls SSOService
				err = d.DecodeElement(&ls, &t)
				if err != nil {
					return err
				}
				e.SSOServices = append(e.SSOServices, ls)
			case descACS:
				var ac AssertionConsumer
				err = d.DecodeElement(&ac, &t)
				if err != nil {
					return err
				}
				e.Consumers = append(e.Consumers, ac)
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
