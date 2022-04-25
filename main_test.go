package main

import (
	"reflect"
	"testing"
)

func Test_parseDiceTrows(t *testing.T) {
	type args struct {
		dice string
	}
	tests := []struct {
		name string
		args args
		want []DiceThrow
	}{
		{
			name: "one of 20",
			args: args{dice: "1d20"},
			want: []DiceThrow{
				{Times: 1, Max: 20, Ok: true},
			},
		},
		{
			name: "two of 20",
			args: args{dice: "2d20"},
			want: []DiceThrow{
				{Times: 2, Max: 20, Ok: true},
			},
		},
		{
			name: "one of 20 and one of 8",
			args: args{dice: "1d20 1d8"},
			want: []DiceThrow{
				{Times: 1, Max: 20, Ok: true},
				{Times: 1, Max: 8, Ok: true},
			},
		},
		{
			name: "invalid aaa",
			args: args{dice: "aaa"},
			want: []DiceThrow{
				{Msg: HelpStr, Ok: false},
			},
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseDiceTrows(tt.args.dice); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseDiceTrows() = %v, want %v", got, tt.want)
			}
		})
	}
}
