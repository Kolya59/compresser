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
		trees[s[0]] = &tree.Tree{
			Value:       s[0],
			Probability: pb,
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
		leafs[fst.Value+snd.Value] = tree.Concat(fst, snd)
	}

	var l *tree.Tree
	for _, el := range leafs {
		l = el
	}

	return l, nil
}

func pack(source string, packer *tree.Tree) ([]byte, error) {
	var res []byte
	for _, c := range source {
		b, err := packer.GetCode(c)
		if err != nil {
			return nil, fmt.Errorf("failed to find code for %c: %w", c, err)
		}
		res = append(res, b...)
	}
	return res, nil
}

func unpack(b []byte, packer *tree.Tree) (string, error) {
	var res, c string
	var err error
	for len(b) != 0 {
		c, b, err = packer.GetValue(b)
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
	delete(*leafs, fst.Value)
	delete(*leafs, snd.Value)
	return
}

func main() {
	// Read freq
	leafs, err := readData("./data/freq.txt")
	if err != nil {
		log.Fatalf("Failed to read data: %v", err)
	}
	leafs[" "] = &tree.Tree{
		Value:       " ",
		Probability: 0.190767000311362,
	}
	leafs["\n"] = &tree.Tree{
		Value:       "\n",
		Probability: 0.0152592971737005,
	}

	// Create tree
	root, err := createTree(leafs)
	if err != nil {
		log.Fatalf("Failed to create tree: %v", err)
	}

	// Get absolute path
	originalPath, err := filepath.Abs("./data/original.txt")
	if err != nil {
		log.Fatalf("failed to get path: %v", err)
	}

	// Read file
	original, err := ioutil.ReadFile(originalPath)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}

	// Pack freq
	packed, err := pack(string(original), root)
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
	}
}
