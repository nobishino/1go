package ast

// Node represents AST node
type Node struct {
	Value int // only used when Kind = Num
	Kind
	Lhs *Node
	Rhs *Node
}

// Kind represents kind of a token
type Kind string

const (
	Num Kind = "Num"
	Add Kind = "Add"
	Sub Kind = "Sub"
)

func NewNode(k Kind, lhs, rhs *Node) *Node {
	return &Node{
		Kind: k,
		Lhs:  lhs,
		Rhs:  rhs,
	}
}
