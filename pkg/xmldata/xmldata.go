package xmldata

// Store holds XML data in a map
var Data = make(map[string]string)

// StoreXML stores XML data with the given name in the global Store
func StoreXML(name string, data string) {
	Data[name] = data
}
