package xmllib

import "fmt"

type NodeType int

const (
	RootNode NodeType = 0
	AttrNode NodeType = 1
	ValNode  NodeType = 2
)

// Node is a data element on a tree
type Node struct {
	Data     string
	NType    NodeType
	Children map[string]NodeList
}

// NodeList is a list of nodes
type NodeList []*Node

// AddChild appends a node to the list of children
func (curNode *Node) AddChild(s string, c *Node) {
	if curNode.Children == nil {
		curNode.Children = map[string]NodeList{}
	}
	curNode.Children[s] = append(curNode.Children[s], c)
}

// HasChildren returns whether it is a complex type (has children)
func (curNode *Node) HasChildren() bool {
	return len(curNode.Children) > 0
}

func (nt NodeType) String() string {
	switch nt {
	case RootNode:
		return "RootNode"
	case AttrNode:
		return "AttrNode"
	case ValNode:
		return "ValNode"
	}
	return fmt.Sprintf("Unknown type of NodeType = %d\n", int(nt))
}
