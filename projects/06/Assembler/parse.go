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
	// fmt.Printf("%#v\n", os.Args)
	// if len(os.Args) < 2 {
	// 	fmt.Print("assembler need filename")
	// }
	// filename := os.Args[1]
	filename := "../rect/RectL.asm"
	name, path, err := fileinfo(filename)
	if err != nil {
		log.Fatal(err)
	}
	println(name, path)

	asmFile, err := os.OpenFile(fmt.Sprintf("%s.hack", name), os.O_RDWR|os.O_CREATE, 0644)
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
	NOT_COMMAND commandType = iota
	A_COMMAND
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

func NewParser(in io.Reader, out io.Writer) *Parser {
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)
	return &Parser{in: in, out: out, scanner: scanner}
}

func (p *Parser) WriteLines(s string) {
	p.out.Write([]byte(s + "\n"))
}

func (p *Parser) Parse() {
	for p.HasMoreCommands() {
		if ok := p.Advance(); ok {
			fmt.Println(p.CurrentCommand)
			if p.CommandType() == A_COMMAND {
				p.WriteLines(fmt.Sprintf("0%s", p.Symbol()))
			} else if p.CommandType() == C_COMMAND {
				fmt.Println(p.Comp(), p.Dest(), p.Jump())
				p.WriteLines(fmt.Sprintf("111%s%s%s", p.Comp(), p.Dest(), p.Jump()))
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
	if strings.HasPrefix(p.CurrentCommand, "//") {
		return NOT_COMMAND
	}

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
	if strings.Contains(p.CurrentCommand, ";") {
		return "000"
	}
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
	"0":   "101010",
	"1":   "111111",
	"-1":  "111010",
	"D":   "001100",
	"A":   "110000",
	"!D":  "001101",
	"!A":  "110001",
	"-D":  "001111",
	"-A":  "110011",
	"D+1": "011111",
	"A+1": "110111",
	"D-1": "001110",
	"A-1": "110010",
	"D+A": "000010",
	"D-A": "010011",
	"A-D": "000111",
	"D&A": "000000",
	"D|A": "010101",
}

func (p *Parser) Comp() (compCode string) {
	comp := ""
	if strings.Contains(p.CurrentCommand, ";") {
		comp = (strings.Split(p.CurrentCommand, ";"))[0]
	} else {
		comp = (strings.Split(p.CurrentCommand, "="))[1]
	}
	compCode = "0"
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

	for index := len(JumpMap) - 1; index > 0; index-- {
		jumpCommand := JumpMap[index]
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
