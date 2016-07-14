package metadata

import (
	"encoding/xml"
	"fmt"
	"io"
	"time"
)

const (
	eds = "EntitiesDescriptor"
	ed  = "EntityDescriptor"
)

// Entities is an abstraction over a list of entities (including single-length
// lists)
type Entities []Entity

// UnmarshalXML implements xml.Unmarshaler
func (e *Entities) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	now := time.Now()

	*e = []Entity{}

	for _, attr := range start.Attr {
		if attr.Name.Local == "validUntil" {
			t, err := time.Parse(time.RFC3339, attr.Value)
			if err != nil {
				return fmt.Errorf("error parsing metadata expiration")
			}
			if now.After(t) {
				return &ErrExpired{EvaluatedAt: now, ValidUntil: t}
			}
		}
	}

	for {
		t, err := d.Token()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("error reading xml token: %s", err)
		}

		switch t := t.(type) {
		case xml.StartElement:
			if t.Name.Local == ed {
				var en Entity
				d.DecodeElement(&en, &t)
				*e = append(*e, en)
			}
		}
	}

	return nil
}

type Entity struct {
	ID       string `xml:",attr"`
	EntityID string `xml:"entityID,attr"`
}
