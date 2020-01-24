package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	uuid "github.com/satori/go.uuid"

	"github.com/kolya59/compresser/pkg/tree"
)

func readData(path string) (map[string]*tree.Tree, error) {
	// Get absolute path
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get path: %w", err)
	}
	// Read file
	file, err := ioutil.ReadFile(abs)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	// Parse data
	splits := strings.Split(string(file), "\n")
	trees := make(map[string]*tree.Tree, len(splits))
	for _, split := range splits {
		s := strings.Split(split, ";")
		pb, err := strconv.ParseFloat(s[1], 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", s[1], err)
		}
		id, _ := uuid.NewV4()
		trees[id.String()] = &tree.Tree{
			Value:       s[0],
			Probability: pb,
			UUID:        id.String(),
		}
	}
	return trees, nil
}

func createTree(leafs map[string]*tree.Tree) (*tree.Tree, error) {
	for len(leafs) > 1 {
		fst, snd, err := popMinCoupleFromSlice(&leafs)
		if err != nil {
			return nil, fmt.Errorf("failed to get min elements: %w", err)
		}
		id, _ := uuid.NewV4()
		leafs[id.String()] = tree.Concat(fst, snd, id.String())
	}

	var l *tree.Tree
	for _, el := range leafs {
		l = el
	}

	return l, nil
}

func bytesToString(b []byte) string {
	res := ""
	n := len(b)
	if n <= 0 {
		return ""
	}
	for _, el := range b {
		res += fmt.Sprintf("%08b", el)
	}
	return res
}

func stringToBytes(s string) []byte {
	var res []byte
	var tmp int64
	for len(s)%8 != 0 {
		s += "0"
	}
	for i := 0; len(s)-i >= 8; i += 8 {
		tmp, _ = strconv.ParseInt(s[i:i+8], 2, 16)
		res = append(res, byte(tmp))
	}
	return res
}

func pack(source string, packer map[string]string) ([]byte, error) {
	var res string
	for _, c := range source {
		b := tree.GetCode(c, packer)
		if b == "" {
			return nil, fmt.Errorf("failed to find code for %c", c)
		}
		res += b
	}
	return stringToBytes(res), nil
}

func unpack(b []byte, packer *tree.Tree) (string, error) {
	var res, c string
	var err error
	s := bytesToString(b)
	for len(s) != 0 {
		c, s, err = packer.GetValue(s)
		if err != nil {
			return "", fmt.Errorf("failed to parse freq: %w", err)
		}
		res += c
	}
	return res, nil
}

func popMinCoupleFromSlice(leafs *map[string]*tree.Tree) (fst, snd *tree.Tree, err error) {
	n := len(*leafs)
	if n < 2 {
		return nil, nil, errors.New("EOF")
	}
	min1 := math.MaxFloat64
	min2 := math.MaxFloat64
	fst = &tree.Tree{Probability: min1}
	for _, v := range *leafs {
		if v.Probability < min1 {
			snd = fst
			min2 = snd.Probability
			fst = v
			min1 = fst.Probability
			continue
		}
		if v.Probability < min2 {
			snd = v
			min2 = v.Probability
		}
	}
	delete(*leafs, fst.UUID)
	delete(*leafs, snd.UUID)
	return
}

func createMap(root *tree.Tree) map[string]string {
	res := make(map[string]string)
	dfs("", root, res)
	return res
}

func dfs(prefix string, node *tree.Tree, m map[string]string) {
	if node.Left != nil {
		dfs(fmt.Sprintf("%s0", prefix), node.Left, m)
	}
	if node.Right != nil {
		dfs(fmt.Sprintf("%s1", prefix), node.Right, m)
	}
	if len(node.Value) == 1 {
		m[node.Value] = prefix
	}
}

func main() {
	// Read freq
	leafs, err := readData("./data/freq1.txt")
	if err != nil {
		log.Fatalf("Failed to read data: %v", err)
	}
	leafs["uuid"] = &tree.Tree{
		Value:       " ",
		Probability: 0.190767000311362,
		UUID:        "uuid",
	}
	leafs["uuid2"] = &tree.Tree{
		Value:       "\n",
		Probability: 0.0152592971737005,
		UUID:        "uuid2",
	}

	// Create tree
	root, err := createTree(leafs)
	if err != nil {
		log.Fatalf("Failed to create tree: %v", err)
	}

	// Create map
	rootMap := createMap(root)

	// Get absolute path
	originalPath, err := filepath.Abs("./data/test.txt")
	if err != nil {
		log.Fatalf("Failed to get path: %v", err)
	}

	// Read file
	original, err := ioutil.ReadFile(originalPath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Pack freq
	packed, err := pack(string(original), rootMap)
	if err != nil {
		log.Fatalf("Failed to pack original: %v", err)
	}

	// Save packet to the file
	packedPath, err := filepath.Abs("./data/packed.bin")
	if err != nil {
		log.Fatalf("Failed to get abs packed path: %v", err)
	}
	if err := ioutil.WriteFile(packedPath, packed, os.ModePerm); err != nil {
		log.Fatalf("Failed to save packed to the file: %v", err)
	}

	// Unpack the packed file
	unpacked, err := unpack(packed, root)
	if err != nil {
		log.Fatalf("Failed to unpack the packed file: %v", err)
	}

	// Compare original and unpacked
	if string(original) != unpacked {
		log.Fatal("Original != unpacked")
	} else {
		log.Print("Original == unpacked")
	}
}
