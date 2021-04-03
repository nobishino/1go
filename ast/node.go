package ast

// Node represents AST node
type Node struct {
	Value int    // only used when Kind = Num
	Name  string // only used when Kind = LocalVar
	// TODO: delete Name field (全ての変数のoffsetはあらかじめ決めておくので名前は必要ないけどデバッグ用に残しておく)
	Kind
	Lhs    *Node
	Rhs    *Node
	Offset int // only used when Kind = LocalVar
}

// Kind represents kind of a node
type Kind string

const (
	Num      Kind = "Num"
	Add      Kind = "Add"
	Sub      Kind = "Sub"
	Mul      Kind = "Mul"
	Div      Kind = "Div"
	Eq       Kind = "Equality"
	Neq      Kind = "NonEquality"
	LT       Kind = "LessThan"
	GT       Kind = "GreaterThan"
	LE       Kind = "LessThanOrEqual"
	GE       Kind = "GreaterThanOrEqual"
	Assign   Kind = "Assignment"
	LocalVar Kind = "Identifier"
)

func NewNode(k Kind, lhs, rhs *Node) *Node {
	return &Node{
		Kind: k,
		Lhs:  lhs,
		Rhs:  rhs,
	}
}

func newNumber(value int) *Node {
	return &Node{
		Kind:  Num,
		Value: value,
	}
}
