package decisions

import (
	"bytes"
	"math/rand"

	u "../utils"
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

// isCondition returns true if the object passed in is a Condition
func isCondition(v interface{}) bool {
	switch v.(type) {
	case Condition:
		return true
	}
	return false
}

// InitializeMetricsMap returns an initialize map of each Metric type to 0
func InitializeMetricsMap() map[Metric]float64 {
	return map[Metric]float64{
		MetricHealth: 0.0,
	}
}

// CopyTreeByValue recursively copies an existing tree by value given an
// existing one, initializing metrics to 0.
func CopyTreeByValue(source *Node) *Node {
	if source == nil {
		return nil
	}
	destination := Node{
		ID:          source.ID,
		Color:       u.MutateColor(source.Color),
		NodeType:    source.NodeType,
		Metrics:     InitializeMetricsMap(),
		MetricsAvgs: InitializeMetricsMap(),
		Uses:        source.Uses,
		YesNode:     CopyTreeByValue(source.YesNode),
		NoNode:      CopyTreeByValue(source.NoNode),
		UsedYes:     source.UsedYes,
		UsedNo:      source.UsedNo,
	}
	return &destination
}

// MutateTree copies a root Node, makes changes to the full tree, and returns
func MutateTree(original *Node) *Node {
	mutated := CopyTreeByValue(original)
	MutateNode(mutated)
	mutated.UpdateNodeIDs()
	return mutated
}

// MutateNode randomly mutates nodes of a tree
func MutateNode(node *Node) {
	node.Uses = 0
	if isCondition(node.NodeType) {
		// If node is a condition and one of its paths has 0 uses, try
		// switching the condition type
		if !node.UsedYes || !node.UsedNo {
			node.NodeType = GetRandomCondition()
		} else {
			if rand.Float64() < 0.5 {
				MutateNode(node.YesNode)
			} else {
				MutateNode(node.NoNode)
			}
		}
	} else {
		if rand.Float64() < ChanceOfAddingNewSubTree {
			originalAction := node.NodeType.(Action)
			node.NodeType = GetRandomCondition()
			node.UsedYes = false
			node.UsedNo = false
			yesNode := Node{}
			noNode := Node{}
			if rand.Float64() < 0.5 {
				yesNode = TreeFromAction(GetRandomAction())
				noNode = TreeFromAction(originalAction)
			} else {
				yesNode = TreeFromAction(originalAction)
				noNode = TreeFromAction(GetRandomAction())
			}
			node.YesNode = &yesNode
			node.NoNode = &noNode
		} else {
			node.NodeType = GetRandomAction()
		}
	}
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
		buffer.WriteString("├─Then... ")
		buffer.WriteString(PrintNode(*node.YesNode, spaces+1))
		for i := 0; i < spaces; i++ {
			buffer.WriteString("  ")
		}
		buffer.WriteString("└─Else: ")
		buffer.WriteString(PrintNode(*node.NoNode, spaces+1))
	}
	return buffer.String()
}
