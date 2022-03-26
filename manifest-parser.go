package vueglue

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

type manifestNode struct {
	key      string
	nodeType reflect.Kind
	value    reflect.Value
	children []*manifestNode
}

type manifestTarget struct {
	File    string   `json:"file"`
	Source  string   `json:"src"`
	IsEntry bool     `json:"isEntry"`
	Imports []string `json:"imports"`
	CSS     []string `json:"css"`
	Nodes   []*manifestNode
}

func (n *manifestNode) subKey(key string) *manifestNode {
	if len(n.children) == 0 {
		return nil
	}
	for _, leaf := range n.children {
		if leaf.key == key {
			return leaf
		}
	}
	return nil
}

// @see https://yourbasic.org/golang/json-example/
func (m *manifestTarget) parseWithoutReflection(jsonData []byte) (*VueGlue, error) {
	var v interface{}
	json.Unmarshal(jsonData, &v)
	topNode := manifestNode{
		key: "top",
	}
	m.Nodes = append(m.Nodes, &topNode)
	m.siftCollections(&topNode, "", "", v)

	// Get entry point
	entry := (*manifestNode)(nil)
	glue := &VueGlue{}

	for _, leaf := range topNode.children {
		if leaf.subKey("isEntry") != nil {
			entry = leaf
			glue.MainModule = leaf.subKey("file").value.String()
			break
		}
	}
	if entry == nil {
		return nil, ErrNoEntryPoint
	}

	imports := entry.subKey("imports")
	if imports == nil || len(imports.children) == 0 {
		// return nil, errors.New("expected code to have js dependencies")
		// turns out this will become optional as of Vite 2.9, so:
	} else {
		for _, child := range imports.children {
			// these have a level of indirection for some reason
			deref := topNode.subKey(child.value.String())
			if deref == nil {
				return nil, ErrNoInputFile
			}
			item := deref.subKey("file")
			if item == nil {
				return nil, ErrManifestBadlyFormed
			}
			glue.Imports = append(glue.Imports, item.value.String())
		}
	}

	css := entry.subKey("css")
	if css == nil || len(css.children) == 0 {
		// not an error, since CSS is optional
		return glue, nil
	}

	for _, child := range css.children {
		glue.CSSModule = append(glue.CSSModule, child.value.String())
	}

	return glue, nil
}

func (m *manifestTarget) siftCollections(leaf *manifestNode, indent, key string, v interface{}) {
	data, ok := v.(map[string]interface{})
	if ok {
		leaf.nodeType = reflect.Map
		for k, v := range data {
			child := &manifestNode{
				key: k,
			}
			leaf.children = append(leaf.children, child)
			m.processInterface(child, indent, k, v)
		}
	} else if arrayData, ok := v.([]interface{}); ok {
		leaf.nodeType = reflect.Slice
		for i, v := range arrayData {
			child := &manifestNode{}
			leaf.children = append(leaf.children, child)
			m.processInterface(child, indent, strconv.Itoa(i), v)
		}
	} else {
		m.processInterface(leaf, indent, key, v)
	}
}

// call this for recurisve structures.
func (m *manifestTarget) processInterface(leaf *manifestNode, indent, k string, v interface{}) {
	// Cover types we know we get in JSON
	switch v := v.(type) {
	case string:
		leaf.nodeType = reflect.String
		leaf.value = reflect.ValueOf(v)
	case float64:
		leaf.nodeType = reflect.Float64
		leaf.value = reflect.ValueOf(v)
	case bool:
		leaf.nodeType = reflect.Bool
		leaf.value = reflect.ValueOf(v)
	case []interface{}:
		m.siftCollections(leaf, indent+"    ", k, v)
	case map[string]interface{}:
		m.siftCollections(leaf, indent+"    ", k, v)
	default:
		fmt.Printf("%s %s ?? %T (unknown)", indent, k, v)
	}
}
