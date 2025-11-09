package main

import (
"fmt"
"os"
"path/filepath"

"gopkg.in/yaml.v3"
)

func validatePort(value interface{}) bool {
switch v := value.(type) {
case int, int64:
return v > 0 && v < 65536
case float64:
return int(v) > 0 && int(v) < 65536
default:
return false
}
}

func validateYAML(filename string) {
data, err := os.ReadFile(filename)
if err != nil {
fmt.Printf("%s: unable to read file: %v\n", filename, err)
return
}

var raw map[string]interface{}
if err := yaml.Unmarshal(data, &raw); err != nil {
fmt.Printf("YAML decode error: %v\n", err)
return
}

base := filepath.Base(filename)
if metadata, ok := raw["metadata"].(map[string]interface{}); ok {
if name, ok := metadata["name"].(string); !ok || name == "" {
fmt.Printf("%s:4 name is required\n", base)
}
}

if spec, ok := raw["spec"].(map[string]interface{}); ok {
if osField, ok := spec["os"]; ok {
if osName, ok := osField.(string); ok && osName != "linux" && osName != "windows" {
fmt.Printf("%s:10 os has unsupported value '%s'\n", base, osName)
}
}
if containers, ok := spec["containers"].([]interface{}); ok {
for _, c := range containers {
container, ok := c.(map[string]interface{})
if !ok {
continue
}
if name, ok := container["name"].(string); !ok || name == "" {
fmt.Printf("%s:12 name is required\n", base)
}
}
}
}
}

func main() {
if len(os.Args) < 2 {
fmt.Println("Usage: yamlvalid <filename>")
return
}
validateYAML(os.Args[1])
}
