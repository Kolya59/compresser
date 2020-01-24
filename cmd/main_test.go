package main

import (
	"reflect"
	"testing"

	"github.com/kolya59/compresser/pkg/tree"
)

func Test_bytesToString(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "10001000",
			args: args{
				b: []byte{17},
			},
			want: "10001000",
		},
		{
			name: "Empty",
			args: args{
				b: []byte{},
			},
			want: "",
		},
		{
			name: "1000100010000000",
			args: args{
				b: []byte{136, 128},
			},
			want: "1000100010000000",
		},
		{
			name: "0100100010",
			args: args{
				b: []byte{72, 2},
			},
			want: "0100100010000000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bytesToString(tt.args.b); got != tt.want {
				t.Errorf("bytesToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createMap(t *testing.T) {
	type args struct {
		root *tree.Tree
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createMap(tt.args.root); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createTree(t *testing.T) {
	type args struct {
		leafs map[string]*tree.Tree
	}
	tests := []struct {
		name    string
		args    args
		want    *tree.Tree
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createTree(tt.args.leafs)
			if (err != nil) != tt.wantErr {
				t.Errorf("createTree() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createTree() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pack(t *testing.T) {
	type args struct {
		source string
		packer map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "AbC",
			args: args{
				source: "AbC",
				packer: map[string]string{
					"A": "1",
					"b": "01",
					"C": "001",
					"d": "0001",
					"a": "00001",
				},
			},
			want:    []byte{164},
			wantErr: false,
		},
		{
			name: "adC",
			args: args{
				source: "adC",
				packer: map[string]string{
					"A": "1",
					"b": "01",
					"C": "001",
					"d": "0001",
					"a": "0000",
				},
			},
			want:    []byte{1, 32},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pack(tt.args.source, tt.args.packer)
			if (err != nil) != tt.wantErr {
				t.Errorf("pack() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pack() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_popMinCoupleFromSlice(t *testing.T) {
	type args struct {
		leafs *map[string]*tree.Tree
	}
	tests := []struct {
		name    string
		args    args
		wantFst *tree.Tree
		wantSnd *tree.Tree
		wantErr bool
	}{
		{
			name: "Empty",
			args: args{
				leafs: &map[string]*tree.Tree{},
			},
			wantFst: nil,
			wantSnd: nil,
			wantErr: true,
		},
		{
			name: "One element",
			args: args{
				leafs: &map[string]*tree.Tree{
					"O": {
						Left:        nil,
						Right:       nil,
						Value:       "O",
						Probability: 0,
					},
				},
			},
			wantFst: nil,
			wantSnd: nil,
			wantErr: true,
		},
		{
			name: "Two elements",
			args: args{
				leafs: &map[string]*tree.Tree{
					"O": {
						Left:        nil,
						Right:       nil,
						Value:       "O",
						Probability: 1,
					},
					"1": {
						Left:        nil,
						Right:       nil,
						Value:       "1",
						Probability: 2,
					},
				},
			},
			wantFst: &tree.Tree{
				Left:        nil,
				Right:       nil,
				Value:       "O",
				Probability: 1,
			},
			wantSnd: &tree.Tree{
				Left:        nil,
				Right:       nil,
				Value:       "1",
				Probability: 2,
			},
			wantErr: false,
		},
		{
			name: "Three elements",
			args: args{
				leafs: &map[string]*tree.Tree{
					"O": {
						Left:        nil,
						Right:       nil,
						Value:       "O",
						Probability: 4,
					},
					"1": {
						Left:        nil,
						Right:       nil,
						Value:       "1",
						Probability: 2,
					},
					"2": {
						Left:        nil,
						Right:       nil,
						Value:       "2",
						Probability: 1,
					},
				},
			},
			wantFst: &tree.Tree{
				Left:        nil,
				Right:       nil,
				Value:       "2",
				Probability: 1,
			},
			wantSnd: &tree.Tree{
				Left:        nil,
				Right:       nil,
				Value:       "1",
				Probability: 2,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFst, gotSnd, err := popMinCoupleFromSlice(tt.args.leafs)
			if (err != nil) != tt.wantErr {
				t.Errorf("popMinCoupleFromSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotFst, tt.wantFst) {
				t.Errorf("popMinCoupleFromSlice() gotFst = %v, want %v", gotFst, tt.wantFst)
			}
			if !reflect.DeepEqual(gotSnd, tt.wantSnd) {
				t.Errorf("popMinCoupleFromSlice() gotSnd = %v, want %v", gotSnd, tt.wantSnd)
			}
		})
	}
}

func Test_stringToBytes(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "10001000",
			args: args{
				s: "10001000",
			},
			want: []byte{136},
		},
		{
			name: "Empty",
			args: args{
				s: "",
			},
			want: nil,
		},
		{
			name: "10001000100000000",
			args: args{
				s: "1000100010000000",
			},
			want: []byte{136, 128},
		},
		{
			name: "0100100010000000",
			args: args{
				s: "0100100010000000",
			},
			want: []byte{72, 128},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stringToBytes(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("stringToBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_unpack(t *testing.T) {
	type args struct {
		b      []byte
		packer *tree.Tree
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "",
			args: args{
				b: nil,
				packer: &tree.Tree{
					Left:        &tree.Tree{},
					Right:       &tree.Tree{},
					Value:       "",
					Probability: 0,
				},
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := unpack(tt.args.b, tt.args.packer)
			if (err != nil) != tt.wantErr {
				t.Errorf("unpack() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("unpack() got = %v, want %v", got, tt.want)
			}
		})
	}
}
