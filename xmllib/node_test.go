package xmllib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddChild(t *testing.T) {
	assert := assert.New(t)

	n := Node{}
	assert.Len(n.Children, 0)

	n.AddChild("a", &Node{})
	assert.Len(n.Children, 1)

	n.AddChild("b", &Node{})
	assert.Len(n.Children, 2)
}

func TestHasChildren(t *testing.T) {
	assert := assert.New(t)

	n := Node{}
	assert.False(n.HasChildren(), "nodes with no children are not complex")

	n.AddChild("b", &Node{})
	assert.True(n.HasChildren(), "nodes with children are complex")

	n.Data = "foo"
	assert.True(n.HasChildren(), "data does not impact HasChildren")
}
