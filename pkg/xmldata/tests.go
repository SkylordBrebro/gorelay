package xmldata

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)

// Function to extract the structure of XML data
func analyzeXMLStructure(xmlData []byte) ([]string, error) {
	decoder := xml.NewDecoder(strings.NewReader(string(xmlData)))
	
	paths := make(map[string]bool)
	currentPath := []string{}
	
	for {
		token, err := decoder.Token()
		if err != nil || token == nil {
			break
		}
		
		switch se := token.(type) {
		case xml.StartElement:
			// Add element name to current path
			currentPath = append(currentPath, se.Name.Local)
			
			// Record the current path
			path := strings.Join(currentPath, ".")
			paths[path] = true
			
			// Record attributes if they exist
			for _, attr := range se.Attr {
				attrPath := path + "[@" + attr.Name.Local + "]"
				paths[attrPath] = true
			}
			
		case xml.EndElement:
			// Remove the last element from path when we hit a closing tag
			if len(currentPath) > 0 {
				currentPath = currentPath[:len(currentPath)-1]
			}
		}
	}
	
	// Convert map keys to sorted slice
	result := make([]string, 0, len(paths))
	for path := range paths {
		result = append(result, path)
	}
	
	sort.Strings(result)
	return result, nil
}

// Function to print the structure as a tree
func printXMLTreeStructure(paths []string) {
	// Create a map to track what we've printed to avoid duplication
	printed := make(map[string]bool)
	
	for _, path := range paths {
		parts := strings.Split(path, ".")
		
		// Print each level of the path
		for i := 1; i <= len(parts); i++ {
			subpath := strings.Join(parts[:i], ".")
			if !printed[subpath] {
				indent := strings.Repeat("  ", i-1)
				
				// Check if it's an attribute
				if strings.Contains(parts[i-1], "[@") {
					// Extract attribute name
					attrName := strings.TrimPrefix(strings.TrimSuffix(parts[i-1], "]"), "[@")
					fmt.Printf("%s- @%s (attribute)\n", indent, attrName)
				} else {
					fmt.Printf("%s- %s\n", indent, parts[i-1])
				}
				
				printed[subpath] = true
			}
		}
	}
}

// Function that combines parsing and analyzing XML structure from a file
func AnalyzeXMLFile(filename string) error {
	// Read the file
	xmlData, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}
	
	// Analyze the structure
	paths, err := analyzeXMLStructure(xmlData)
	if err != nil {
		return fmt.Errorf("error analyzing XML: %w", err)
	}
	
	// Print the tree structure
	fmt.Printf("XML Structure of %s:\n", filename)
	printXMLTreeStructure(paths)
	
	// Print the flat paths
	fmt.Println("\nFlat Path Listing:")
	for _, path := range paths {
		fmt.Println(path)
	}
	
	return nil
}

// Function that analyzes XML data directly as a string or bytes
func AnalyzeXMLString(xmlData string) error {
	// Analyze the structure
	paths, err := analyzeXMLStructure([]byte(xmlData))
	if err != nil {
		return fmt.Errorf("error analyzing XML: %w", err)
	}
	
	// Print the tree structure
	fmt.Println("XML Structure:")
	printXMLTreeStructure(paths)
	
	// Print the flat paths
	fmt.Println("\nFlat Path Listing:")
	for _, path := range paths {
		fmt.Println(path)
	}
	
	return nil
}