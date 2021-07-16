package parse

import (
	"bufio"
	"fmt"
	"io"
	"testing"
)

func TestParser_Dest(t *testing.T) {
	type fields struct {
		in             io.Reader
		out            io.Writer
		scanner        *bufio.Scanner
		CurrentCommand string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "M",
			fields: fields{
				CurrentCommand: "M=D+M",
			},
			want: "001"},
		{name: "D",
			fields: fields{
				CurrentCommand: "D=D-A",
			},
			want: "010"},
		{name: "A",
			fields: fields{
				CurrentCommand: "A=D-A",
			},
			want: "100"},
		{name: "AM",
			fields: fields{
				CurrentCommand: "AM=D-A",
			},
			want: "101"},
		{name: "AD",
			fields: fields{
				CurrentCommand: "AD=D-A",
			},
			want: "110"},
		{name: "MD",
			fields: fields{
				CurrentCommand: "MD=D-A",
			},
			want: "011"},
		{name: "AMD",
			fields: fields{
				CurrentCommand: "AMD=D-A",
			},
			want: "111"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				// in:             tt.fields.in,
				// out:            tt.fields.out,
				scanner:        tt.fields.scanner,
				CurrentCommand: tt.fields.CurrentCommand,
			}
			if got := p.Dest(); got != tt.want {
				t.Errorf("Parser.Dest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_Jump(t *testing.T) {
	type fields struct {
		in             io.Reader
		out            io.Writer
		scanner        *bufio.Scanner
		CurrentCommand string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "000",
			fields: fields{
				CurrentCommand: "M=D+M",
			},
			want: "000"},
		{name: "JMP",
			fields: fields{
				CurrentCommand: "D:JMP",
			},
			want: "111"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				// in:             tt.fields.in,
				// out:            tt.fields.out,
				scanner:        tt.fields.scanner,
				CurrentCommand: tt.fields.CurrentCommand,
			}
			if got := p.Jump(); got != tt.want {
				t.Errorf("Parser.Jump() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_Main(t *testing.T) {
	type fields struct {
		filename string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// {
		// 	"Pong", fields{"../pong/Pong.asm"},
		// },
		// {
		// 	"max", fields{"../max/Max.asm"},
		// },
		{
			"rect", fields{"../rect/Rect.asm"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.fields.filename)
			p.Main()
		})
	}
}

func TestSymbolTable_buildTable(t *testing.T) {
	type fields struct {
		filename string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			"Max", fields{"../max/Max.asm"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSymbolTable(tt.fields.filename)
			s.buildTable()
			fmt.Printf("%#v", s.table)
		})
	}
}
