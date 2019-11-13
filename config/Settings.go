package config

import "encoding/xml"

type Settings struct {
	XMLName        xml.Name         `xml:"settings"`
	LocalRepository string `xml:"localRepository"`
}
