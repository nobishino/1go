package c

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/nobishino/1go/ast"
	"golang.org/x/xerrors"
)

// NaiveAddSub compiles source code and returns assembly code
func NaiveAddSub(src string) (string, error) {
	buf := new(bytes.Buffer)
	fmt.Fprintln(buf, ".intel_syntax noprefix")
	fmt.Fprintln(buf, ".globl main")
	fmt.Fprintln(buf, "")
	fmt.Fprintln(buf, "main:")
	op := "mov"
	tokens, err := Tokenize(src)
	if err != nil {
		return "", err
	}
	for _, token := range tokens {
		switch token {
		case "+":
			op = "add"
		case "-":
			op = "sub"
		default: // digit token
			fmt.Fprintf(buf, "    %s rax, %s\n", op, token)
		}
	}
	fmt.Fprintln(buf, "    ret")
	return buf.String(), nil
}

// Compile compiles source code and returns asssembly w/ intel syntax
func Compile(src string) (string, error) {
	p, err := ast.NewTParser(src)
	if err != nil {
		return "", err
	}
	parsed, err := p.Program()
	if err != nil {
		return "", err
	}
	result := Gen(parsed, p.GetOffset())
	return strings.Join(result, "\n"), nil
}

func Tokenize(src string) ([]string, error) {
	if err := Validate(src); err != nil {
		return nil, err
	}
	isOp := func(r rune) bool {
		operators := []rune{'+', '-', '*', '/', '(', ')'}
		for _, op := range operators {
			if r == op {
				return true
			}
		}
		return false
	}
	var i, j int
	var result []string
	s := []rune(src)
	for i < len(src) {
		for j < len(src) && !isOp(s[j]) {
			j++
		}
		var token string
		if i == j {
			token = string(s[i : i+1])
			j++
		} else {
			token = string(s[i:j])
		}
		i = j
		token = strings.Trim(token, " ")
		result = append(result, token)
	}
	return result, nil
}

func Validate(src string) error {
	isValid := func(r rune) bool {
		validChars := []rune{'+', '-', '*', '/', '(', ')', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', ' '}
		for _, c := range validChars {
			if c == r {
				return true
			}
		}
		return false
	}
	for i, r := range src {
		if !isValid(r) {
			markerLine := strings.Repeat(" ", i) + "^"
			return errors.New(fmt.Sprintf("invalid character at %d\n%s\n%s", i+1, src, markerLine))
		}
	}
	return nil
}

func Gen(nodes []*ast.Node, offset int) []string {
	if nodes == nil {
		return nil
	}
	result := []string{
		".intel_syntax noprefix",
		".globl main",
		"",
		"main:",
	}
	result = append(result, genPrologue(offset)...)
	for _, node := range nodes {
		result = append(result, genAST(node)...)
		result = append(result,
			"    pop rax",
		)
	}
	result = append(result, epilogue...)
	result = append(result, "")
	return result
}

// 指定したローカル変数オフセットから関数プロローグを生成する
func genPrologue(offset int) []string {
	return []string{
		"    push rbp",                         // 関数呼び出し前(callerの関数の実行時の)RBPレジスタの値をスタックに保存する
		"    mov rbp, rsp",                     // この関数の実行中に基準点とするメモリアドレスをRBPレジスタにセットする
		fmt.Sprintf("    sub rsp, %d", offset), // 8 x 26 bit をこの関数呼び出しインスタンスのローカル変数領域としてスタック領域に確保する
	}
}

func genAST(node *ast.Node) []string {
	if node == nil {
		return nil
	}
	var result []string
	result = append(result, genAST(node.Lhs)...)
	result = append(result, genAST(node.Rhs)...)

	switch node.Kind {
	case ast.Add:
		result = append(result, add...)
	case ast.Sub:
		result = append(result, sub...)
	case ast.Mul:
		result = append(result, mul...)
	case ast.Div:
		result = append(result, div...)
	case ast.Eq:
		result = append(result, eq...)
	case ast.Neq:
		result = append(result, neq...)
	case ast.LT:
		result = append(result, lt...)
	case ast.LE:
		result = append(result, le...)
	case ast.Num:
		result = append(result, fmt.Sprintf("    push %d", node.Value))
	case ast.Assign:
		pushMemAddr, err := genLeftValue(node.Lhs)
		if err != nil {
			panic(err) // TODO: 適切なエラー処理を行う
		}
		result = append(result, pushMemAddr...)
		result = append(result, genAST(node.Rhs)...)  // 右辺のノードを評価する
		result = append(result, assignRightToLeft...) // 代入命令を生成する
	case ast.Return:
		result = append(result, genAST(node.Lhs)...)
		result = append(result, ret...)
	}
	return result
}

// ローカル変数値(左辺値)のメモリアドレスをスタックにプッシュする命令を生成する
func genLeftValue(node *ast.Node) ([]string, error) {
	if node.Kind != ast.LocalVar {
		return nil, xerrors.Errorf("expect left value but got node of kind %q", node.Kind)
	}
	return []string{
		"    mov rax, rbp",                          // ベースポインタの値をraxにコピーする
		fmt.Sprintf("    sub rax, %d", node.Offset), // ベースポインタの値から変数名で決まるオフセットを引く
		"    push rax",
	}, nil
}

var headers = []string{
	".intel_syntax noprefix",
	".globl main",
	"",
	"main:",
}

var add = []string{
	"    pop rdi",
	"    pop rax",
	"    add rax, rdi",
	"    push rax",
}
var sub = []string{
	"    pop rdi",
	"    pop rax",
	"    sub rax, rdi",
	"    push rax",
}

var mul = []string{
	"    pop rdi",
	"    pop rax",
	"    imul rax, rdi",
	"    push rax",
}

var div = []string{
	"    pop rdi",
	"    pop rax",
	"    cqo",
	"    idiv rdi",
	"    push rax",
}

var eq = []string{
	"    pop rdi",
	"    pop rax",
	"    cmp rax, rdi",
	"    sete al",
	"    movzb rax, al",
	"    push rax",
}

var neq = []string{
	"    pop rdi",
	"    pop rax",
	"    cmp rax, rdi",
	"    setne al",
	"    movzb rax, al",
	"    push rax",
}

var lt = []string{
	"    pop rdi",
	"    pop rax",
	"    cmp rax, rdi",
	"    setl al",
	"    movzb rax, al",
	"    push rax",
}

var le = []string{
	"    pop rdi",
	"    pop rax",
	"    cmp rax, rdi",
	"    setle al",
	"    movzb rax, al",
	"    push rax",
}

var assignRightToLeft = []string{
	"    pop rdi",        // 右辺値(評価結果)
	"    pop rax",        // 左辺値のメモリアドレス
	"    mov [rax], rdi", // ローカル変数のメモリ位置に右辺値をコピーする
	"    push rdi",       // 代入された値は代入式自体の値になるのでスタックにpushする
}

var prologue = []string{
	"    push rbp",     // 関数呼び出し前(callerの関数の実行時の)RBPレジスタの値をスタックに保存する
	"    mov rbp, rsp", // この関数の実行中に基準点とするメモリアドレスをRBPレジスタにセットする
	"    sub rsp, 208", // 8 x 26 bit をこの関数呼び出しインスタンスのローカル変数領域としてスタック領域に確保する
}

var epilogue = []string{
	"    mov rsp, rbp", // ベースポインタの位置までRSPを戻してくる。これによりローカル変数領域が「捨てられる」
	"    pop rbp",      // 1つ上の関数に対するベースの値をRBPに書き戻す。このpop命令の後、RSPはこの関数のリターンアドレスが書き込まれたメモリアドレスを指している
	"    ret",          // Stackからpopし、そのpopした値のメモリアドレスに移動する。
}

// return文のNodeのgenerate結果
var ret = []string{
	"    pop rax",      // 直前に積まれた値=return x;のxの評価結果をraxにpopする
	"    mov rsp, rbp", // ベースポインタの位置までRSPを戻してくる。これによりローカル変数領域が「捨てられる」
	"    pop rbp",      // 1つ上の関数に対するベースの値をRBPに書き戻す。このpop命令の後、RSPはこの関数のリターンアドレスが書き込まれたメモリアドレスを指している
	"    ret",          // Stackからpopし、そのpopした値のメモリアドレスに移動する。
}
