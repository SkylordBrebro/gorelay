package services

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
)

// XMLToJSON converts XML data to JSON format
type XMLToJSON struct{}

// NewXMLToJSON creates a new XML to JSON converter
func NewXMLToJSON() *XMLToJSON {
	return &XMLToJSON{}
}

// Convert converts XML data to JSON
func (x *XMLToJSON) Convert(xmlData []byte) ([]byte, error) {
	var xmlMap map[string]interface{}
	if err := xml.Unmarshal(xmlData, &xmlMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal XML: %v", err)
	}

	// Convert XML map to JSON
	jsonData, err := json.Marshal(xmlMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %v", err)
	}

	return jsonData, nil
}

// ConvertStream converts XML data from a reader to JSON
func (x *XMLToJSON) ConvertStream(reader io.Reader) ([]byte, error) {
	xmlData, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read XML data: %v", err)
	}

	return x.Convert(xmlData)
}

// ConvertFile converts an XML file to JSON
func (x *XMLToJSON) ConvertFile(xmlPath, jsonPath string) error {
	xmlData, err := ioutil.ReadFile(xmlPath)
	if err != nil {
		return fmt.Errorf("failed to read XML file: %v", err)
	}

	jsonData, err := x.Convert(xmlData)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(jsonPath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %v", err)
	}

	return nil
}
