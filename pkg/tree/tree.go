package tree

import (
	"fmt"
)

type Tree struct {
	UUID        string
	Left        *Tree
	Right       *Tree
	Value       string
	Probability float64
}

func GetCode(c int32, packer map[string]string) string {
	s := fmt.Sprintf("%c", c)
	return packer[s]
}

func (t *Tree) GetValue(code string) (string, string, error) {
	tmp := &*t
	for i := 0; i < len(code); i++ {
		switch true {
		case tmp.isLeaf():
			return tmp.Value, code[i:], nil
		case code[i] == '0':
			tmp = tmp.Left
		case code[i] == '1':
			tmp = tmp.Right
		default:
			return "", "", fmt.Errorf("invalid character: %c", code[i])
		}
	}
	return tmp.Value, "", nil
}

func (t *Tree) isLeaf() bool {
	return t.Left == nil && t.Right == nil
}

func Concat(left, right *Tree, uuid string) *Tree {
	return &Tree{
		Left:        left,
		Right:       right,
		Probability: left.Probability + right.Probability,
		Value:       "",
		UUID:        uuid,
	}
}
