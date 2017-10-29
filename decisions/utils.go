package decisions

import (
	"bytes"
	"math/rand"
)

// GetRandomCondition returns a random Condition from the Conditions array
func GetRandomCondition() Condition {
	return Conditions[rand.Intn(len(Conditions))]
}

// GetRandomAction returns a random Action from the Actions array
func GetRandomAction() Action {
	return Actions[rand.Intn(len(Actions))]
}

// isAction returns true if the object passed in is an Action
func isAction(v interface{}) bool {
	switch v.(type) {
	case Action:
		return true
	}
	return false
}

// InitializeMetricsMap returns an initialize map of each Metric type to 0
func InitializeMetricsMap() map[Metric]float32 {
	return map[Metric]float32{
		MetricHealth: 0.0,
	}
}

// CopyTreeByValue recursively copies an existing tree by value given an
// existing one, initializing uses and metrics to 0.
func CopyTreeByValue(source *Node) *Node {
	destination := Node{
		ID:       source.ID,
		NodeType: source.NodeType,
		Metrics:  InitializeMetricsMap(),
		Uses:     0,
		YesNode:  CopyTreeByValue(source.YesNode),
		NoNode:   CopyTreeByValue(source.NoNode),
	}
	return &destination
}

// PrintNode prints node and all children showing hierarchy
func PrintNode(node Node, spaces int) string {
	var buffer bytes.Buffer
	buffer.WriteString(Map[node.NodeType])
	buffer.WriteString("\n")
	if !isAction(node.NodeType) {
		for i := 0; i < spaces; i++ {
			buffer.WriteString("  ")
		}
		buffer.WriteString("Then: ")
		buffer.WriteString(PrintNode(*node.YesNode, spaces+1))
		for i := 0; i < spaces; i++ {
			buffer.WriteString("  ")
		}
		buffer.WriteString("Otherwise: ")
		buffer.WriteString(PrintNode(*node.NoNode, spaces+1))
	}
	return buffer.String()
}
