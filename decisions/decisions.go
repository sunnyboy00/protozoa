package decisions

import (
	"bytes"
	"fmt"
	"math/rand"
)

// Action is the custom type for all Organism actions
type Action int

// Condition is the custom type for all Organism conditions
type Condition int

// Sequence is a slice representing a serialized tree of NodeType values
type Sequence []interface{}

// Define all possible actions for Organism
const (
	ActEat Action = iota
	ActIdle
	ActMove
	ActTurnLeft
	ActTurnRight
	CanMove Condition = iota
	IsFoodAhead
	IsFoodLeft
	IsFoodRight
)

// Define slices
var (
	Actions    = [...]Action{ActEat, ActIdle, ActMove, ActTurnLeft, ActTurnRight}
	Conditions = [...]Condition{CanMove, IsFoodAhead, IsFoodLeft, IsFoodRight}
	Map        = map[interface{}]string{
		ActEat:       "A_Eat",
		ActIdle:      "A_Idle",
		ActMove:      "A_Move",
		ActTurnLeft:  "A_Left",
		ActTurnRight: "A_Right",
		CanMove:      "C_Move",
		IsFoodAhead:  "C_FoodAhead",
		IsFoodLeft:   "C_FoodLeft",
		IsFoodRight:  "C_FoodRight",
	}
)

// Node includes an Action or Condition value
type Node struct {
	NodeType interface{}
	UseCount int
	YesNode  *Node
	NoNode   *Node
}

// IsAction returns true if Node's type is Action (false if Condition)
func (n *Node) IsAction() bool {
	return isAction(n.NodeType)
}

// NewRandomSequence generates a new Sequence of random length
func NewRandomSequence() Sequence {
	sequence := NewRandomSubSequence()
	return sequence
}

// MutateSequence mutates a given sequence by replacing a random number of
// actions with condition - action - action blocks
func MutateSequence(sequence Sequence) Sequence {
	mutatedSequence := make(Sequence, len(sequence))
	copy(mutatedSequence, sequence)
	index := rand.Intn(len(mutatedSequence))
	if isAction(mutatedSequence[index]) {
		// in 25 % of cases where the node is an Action, replace with C-A-A
		if rand.Float32() < 0.25 {
			return MutateByAddingSubSequence(mutatedSequence, index)
		}
		// fmt.Println("\nMutating (changing action)")
		// fmt.Printf("\nBefore: %s", PrintSequence(mutatedSequence))
		mutatedSequence[index] = GetRandomAction()
		// fmt.Printf("\nAfter: %s\n", PrintSequence(mutatedSequence))
	} else {
		// fmt.Println("\nMutating (changing condition)")
		// fmt.Printf("\nBefore: %s", PrintSequence(mutatedSequence))
		mutatedSequence[index] = GetRandomCondition()
		// fmt.Printf("\nAfter: %s\n", PrintSequence(mutatedSequence))
	}
	return mutatedSequence
}

// MutateByAddingSubSequence takes a sequence and index of an action in that
// sequence. Replaces the action with a randomly generated sub-sequence and
// returns the result
func MutateByAddingSubSequence(sequence Sequence, index int) Sequence {
	mutatedSequence := sequence
	// insert random subsquence in place of action index to be replaced
	// fmt.Println("\nMutating (adding 2 nodes)")
	// fmt.Printf("\nBefore: %s", PrintSequence(sequence))
	subSequence := NewRandomSubSequence()
	subSequence = append(subSequence, mutatedSequence[index+1:]...)
	mutatedSequence = append(mutatedSequence[:index], subSequence...)
	// fmt.Printf("\nAfter: %s\n", PrintSequence(mutatedSequence))
	return mutatedSequence
}

func MutateByChangingAction(sequence Sequence, index int) Sequence {
	mutatedSequence := sequence
	// insert random subsquence in place of action index to be replaced
	subSequence := NewRandomSubSequence()
	subSequence = append(subSequence, mutatedSequence[index+1:]...)
	mutatedSequence = append(mutatedSequence[:index], subSequence...)
	return mutatedSequence
}

// TreeFromSequence recursively calls itself to create a Node and its
// children from a sequence slice.
func TreeFromSequence(sequence, fullSequence Sequence) Node {
	nodeType := sequence[0]
	if isAction(nodeType) {
		return Node{NodeType: nodeType, UseCount: 0}
	}
	index := 1
	numActionsMinusConditions := 0
	for numActionsMinusConditions < 1 {
		if index >= len(sequence) {
			fmt.Printf("\nThis is gonna fail: %s\n", PrintSequence(fullSequence))
		}
		sequenceItem := sequence[index]
		if isAction(sequenceItem) {
			numActionsMinusConditions++
		} else {
			numActionsMinusConditions--
		}
		index++
	}
	yesNode := TreeFromSequence(sequence[1:index], fullSequence)
	noNode := TreeFromSequence(sequence[index:], fullSequence)
	node := Node{
		NodeType: nodeType,
		UseCount: 0,
		YesNode:  &yesNode,
		NoNode:   &noNode,
	}
	// if len(sequence) == len(fullSequence) {
	// 	fmt.Printf("Tree from sequence: %s", PrintSequence(fullSequence))
	// 	fmt.Print("\n\n")
	// 	fmt.Print(PrintNode(node, 1))
	// 	fmt.Print("\n\n")
	// }
	return node
}

// PrintSequence prints sequence chronologically
func PrintSequence(sequence Sequence) string {
	var buffer bytes.Buffer
	for i, s := range sequence {
		if i > 0 {
			buffer.WriteString("-")
		}
		buffer.WriteString(Map[s])
	}
	return buffer.String()
}

// PrintSequence prints node and all children showing hierarchy
func PrintNode(node Node, spaces int) string {
	var buffer bytes.Buffer
	buffer.WriteString(Map[node.NodeType])
	buffer.WriteString("\n")
	if !isAction(node.NodeType) {
		for i := 0; i < spaces; i++ {
			buffer.WriteString("\t")
		}
		buffer.WriteString("Y: ")
		buffer.WriteString(PrintNode(*node.YesNode, spaces+1))
		// buffer.WriteString("\n")
		for i := 0; i < spaces; i++ {
			buffer.WriteString("\t")
		}
		buffer.WriteString("N: ")
		buffer.WriteString(PrintNode(*node.NoNode, spaces+1))
	}
	return buffer.String()
}
