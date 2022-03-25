package vueglue

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"path/filepath"
	"reflect"
	"strconv"
)

const (
	AssetsDir = "dist/assets"
	RootItem  = "src/main.ts" // set in vite.config.ts
)

type VueGlue struct {

	// Entry point for JS
	MainModule string

	// JS Dependencies / Vendor libs
	Imports []string

	// Bundled CSS
	CSSModule []string

	// I use a 'data-entryp' attribute to find what
	// components to load. Lookup is in src/main.ts.
	MountPoint string

	// An embed that points to the Vue/Vite dist
	// directory.
	DistFS fs.ReadFileFS
}

type ManifestNode struct {
	Key      string
	Type     reflect.Kind
	Value    reflect.Value
	Children []*ManifestNode
}

type ManifestTarget struct {
	File    string   `json:"file"`
	Source  string   `json:"src"`
	IsEntry bool     `json:"isEntry"`
	Imports []string `json:"imports"`
	CSS     []string `json:"css"`
	Nodes   []*ManifestNode
}

func (n *ManifestNode) subKey(key string) *ManifestNode {
	if len(n.Children) == 0 {
		return nil
	}
	for _, leaf := range n.Children {
		if leaf.Key == key {
			return leaf
		}
	}
	return nil
}

// @see https://yourbasic.org/golang/json-example/
func (m *ManifestTarget) parseWithoutReflection(jsonData []byte) (*VueGlue, error) {
	var v interface{}
	json.Unmarshal(jsonData, &v)
	topNode := ManifestNode{
		Key: "top",
	}
	m.Nodes = append(m.Nodes, &topNode)
	m.siftCollections(&topNode, "", "", v)

	// Get entry point
	entry := (*ManifestNode)(nil)
	glue := &VueGlue{}

	for _, leaf := range topNode.Children {
		if leaf.subKey("isEntry") != nil {
			entry = leaf
			glue.MainModule = leaf.subKey("file").Value.String()
			break
		}
	}
	if entry == nil {
		return nil, ErrNoEntryPoint
	}

	imports := entry.subKey("imports")
	if imports == nil || len(imports.Children) == 0 {
		// return nil, errors.New("expected code to have js dependencies")
		// turns out this will become optional as of Vite 2.9, so:
	} else {
		for _, child := range imports.Children {
			// these have a level of indirection for some reason
			deref := topNode.subKey(child.Value.String())
			if deref == nil {
				return nil, ErrNoInputFile
			}
			item := deref.subKey("file")
			if item == nil {
				return nil, ErrManifestBadlyFormed
			}
			glue.Imports = append(glue.Imports, item.Value.String())
		}
	}

	css := entry.subKey("css")
	if css == nil || len(css.Children) == 0 {
		// not an error, since CSS is optional
		return glue, nil
	}

	for _, child := range css.Children {
		glue.CSSModule = append(glue.CSSModule, child.Value.String())
	}

	return glue, nil
}

func (m *ManifestTarget) siftCollections(leaf *ManifestNode, indent, key string, v interface{}) {
	data, ok := v.(map[string]interface{})
	if ok {
		leaf.Type = reflect.Map
		for k, v := range data {
			child := &ManifestNode{
				Key: k,
			}
			leaf.Children = append(leaf.Children, child)
			m.processInterface(child, indent, k, v)
		}
	} else if arrayData, ok := v.([]interface{}); ok {
		leaf.Type = reflect.Slice
		for i, v := range arrayData {
			child := &ManifestNode{}
			leaf.Children = append(leaf.Children, child)
			m.processInterface(child, indent, strconv.Itoa(i), v)
		}
	} else {
		m.processInterface(leaf, indent, key, v)
	}
}

// call this for recurisve structures.
func (m *ManifestTarget) processInterface(leaf *ManifestNode, indent, k string, v interface{}) {
	// Cover types we know we get in JSON
	switch v := v.(type) {
	case string:
		leaf.Type = reflect.String
		leaf.Value = reflect.ValueOf(v)
	case float64:
		leaf.Type = reflect.Float64
		leaf.Value = reflect.ValueOf(v)
	case bool:
		leaf.Type = reflect.Bool
		leaf.Value = reflect.ValueOf(v)
	case []interface{}:
		m.siftCollections(leaf, indent+"    ", k, v)
	case map[string]interface{}:
		m.siftCollections(leaf, indent+"    ", k, v)
	default:
		fmt.Printf("%s %s ?? %T (unknown)", indent, k, v)
	}
}

func ParseManifest(contents []byte) (*VueGlue, error) {
	var testRslt ManifestTarget
	glue, err := testRslt.parseWithoutReflection(contents)
	if err != nil {
		return nil, err
	}
	return glue, nil
}

func NewVueGlue(dist fs.ReadFileFS, pathToDist string) (*VueGlue, error) {

	if !fs.ValidPath(pathToDist) {
		return nil, ErrManifestDNF
	}

	// Get the manifest file
	manifestFile := filepath.Join(pathToDist, "manifest.json")
	contents, err := dist.ReadFile(manifestFile)
	if err != nil {
		return nil, err
	}
	glue, err := ParseManifest(contents)
	if err != nil {
		return nil, err
	}
	glue.DistFS = dist

	output, _ := json.MarshalIndent(glue, "", "  ")
	fmt.Println(string(output))

	tags, err := glue.RenderTags()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(tags)

	return glue, nil

}
