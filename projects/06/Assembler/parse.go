package parse

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func padding(s string, l int, pad string) string {
	for len([]rune(s)) < l {
		s = pad + s
	}
	return s
}

func main1() {
	fmt.Printf("%#v\n", os.Args)
	if len(os.Args) < 2 {
		fmt.Print("assembler need filename")
	}
	filename := os.Args[1]
	name, path, err := fileinfo(filename)
	if err != nil {
		log.Fatal(err)
	}
	println(name, path)

	asmFile, err := os.OpenFile(fmt.Sprintf("%s.jack", name), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer asmFile.Close()

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	parser := NewParser(f, asmFile)
	parser.Parse()

}

func fileinfo(filename string) (name string, path string, err error) {
	name = strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
	path, err = filepath.Abs(filepath.Dir(filename))
	return
}

type commandType int

// commandType
const (
	A_COMMAND commandType = iota
	C_COMMAND
	L_COMMAND
)

// Parser parse jack assembler to machinecode
type Parser struct {
	in             io.Reader
	out            io.Writer
	scanner        *bufio.Scanner
	CurrentCommand string
}

type a_command struct {
}

func NewParser(in io.Reader, out io.Writer) *Parser {
	scanner := bufio.NewScanner(in)
	return &Parser{in: in, out: out, scanner: scanner}
}

func (p *Parser) Parse() {
	for p.HasMoreCommands() {
		if ok := p.Advance(); ok {
			if p.CommandType() == A_COMMAND {
				p.Symbol()
			}
		}
	}
}
func (p *Parser) HasMoreCommands() bool {
	return p.scanner.Scan()
}
func (p *Parser) Advance() bool {
	if t := p.scanner.Text(); t != "" {
		p.CurrentCommand = t
		return true
	}
	return false
}

func (p *Parser) CommandType() commandType {
	if strings.HasPrefix(p.CurrentCommand, "@") {
		if regexp.MustCompile(`@\d+`).MatchString(p.CurrentCommand) {
			return A_COMMAND
		}
		return L_COMMAND
	}
	return C_COMMAND
}
func (p *Parser) Symbol() string {
	symbol := strings.TrimPrefix(p.CurrentCommand, "@")
	if p.CommandType() == A_COMMAND {
		i64, err := strconv.ParseInt(symbol, 10, 0)
		if err != nil {
			log.Fatal(err)
		}
		return padding(strconv.FormatInt(i64, 2), 15, "0")
	}
	return ""
}

const (
	M uint = 1 << iota
	D
	A
)

func (p *Parser) Dest() string {
	var r uint = 0
	switch (strings.Split(p.CurrentCommand, "="))[0] {
	case "M":
		r = M
		break
	case "D":
		r = D
		break
	case "MD":
		r = M | D
		break
	case "A":
		r = A
		break
	case "AM":
		r = A | M
		break
	case "AD":
		r = A | D
		break
	case "AMD":
		r = A | M | D
	}
	return padding(strconv.FormatInt(int64(r), 2), 3, "0")
}

var CompMap = map[string]string{
	"0":   "0101010",
	"1":   "0111111",
	"-1":  "0111010",
	"D":   "0001100",
	"A":   "0110000",
	"!D":  "0001101",
	"!A":  "0110001",
	"-D":  "0001111",
	"-A":  "0110011",
	"D+1": "0011111",
	"A+1": "0110111",
	"D-1": "0001110",
	"A-1": "0110010",
	"D+A": "0000010",
	"D-A": "0010011",
	"A-D": "0000111",
	"D&A": "0000000",
	"D|A": "0010101",
}

func (p *Parser) Comp() (compCode string) {
	compCode = "0"
	comp := (strings.Split(p.CurrentCommand, "="))[1]
	if strings.Contains(comp, "M") {
		compCode = "1"
	}
	comp = strings.ReplaceAll(comp, "M", "A")
	compCode += CompMap[comp]
	return
}

var JumpMap = []string{
	"", "JGT", "JEQ", "JGE", "JLT", "JNE", "JLE", "JMP",
}

func (p *Parser) Jump() string {
	var r int
	for index, jumpCommand := range JumpMap {
		if strings.Contains(p.CurrentCommand, jumpCommand) {
			r = index
			break
		}
	}
	return padding(strconv.FormatInt(int64(r), 2), 3, "0")
}

type Code struct {
}

func (c *Code) Dest() int {
	return 0
}
func (c *Code) Comp() int {
	return 0
}
func (c *Code) Jump() int {
	return 0
}
