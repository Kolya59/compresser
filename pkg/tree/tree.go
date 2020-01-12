package tree

import (
	"fmt"
	"strings"
)

type Tree struct {
	Left        *Tree
	Right       *Tree
	Value       string
	Probability float64
}

func (t *Tree) GetCode(c int32) ([]byte, error) {
	s := fmt.Sprintf("%c", c)
	var res []byte
	tmp := &*t
loop:
	for {
		switch true {
		case tmp != nil && tmp.Value == s:
			break loop
		case tmp.Left != nil && strings.Contains(tmp.Left.Value, s):
			res = append(res, 0)
			tmp = tmp.Left
		case tmp.Right != nil && strings.Contains(tmp.Right.Value, s):
			res = append(res, 1)
			tmp = tmp.Right
		default:
			return nil, fmt.Errorf("the tree doesn't contain %s", s)
		}
	}
	return res, nil
}

func (t *Tree) GetValue(code []byte) (string, []byte, error) {
	tmp := &*t
	for i := 0; len(code) != 0; i++ {
		switch true {
		case tmp.isLeaf():
			return tmp.Value, code[i+1:], nil
		case code[i] == '0':
			tmp = tmp.Left
		case code[i] == '1':
			tmp = tmp.Right
		default:
			return "", nil, fmt.Errorf("invalid character: %c", code[i])
		}
	}
	return tmp.Value, nil, nil
}

func (t *Tree) isLeaf() bool {
	return t.Left == nil && t.Right == nil
}

func Concat(left, right *Tree) *Tree {
	return &Tree{
		Left:        left,
		Right:       right,
		Probability: left.Probability + right.Probability,
		Value:       left.Value + right.Value,
	}
}
