package parse

import (
	"bufio"
	"io"
	"testing"
)

func TestMain1(t *testing.T) {
	type fields struct {
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "main",
			fields: fields{},
			want:   "1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main1()
			if got := "1"; got != tt.want {
				t.Errorf("main() = %v, want %v", got, tt.want)
			}
		})
	}
}
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
