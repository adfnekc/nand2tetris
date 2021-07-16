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

type commandType int

const debug = false

const (
	NOT_COMMAND commandType = iota
	L_COMMAND
	A_COMMAND
	C_COMMAND
)

func IntTosBin(i int) string {
	s := strconv.Itoa(i)
	i64, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		log.Fatal(err)
	}
	return strconv.FormatInt(i64, 2)
}

var COMMAND_FUNC = map[commandType]func(p *Parser) string{
	NOT_COMMAND: func(p *Parser) string {
		return ""
	},
	A_COMMAND: func(p *Parser) string {
		symbol := p.Symbol()
		if _, err := strconv.Atoi(symbol); err != nil {
			// symbol is predefined variable or lable
			if p.ContainsSymbol(symbol) {
				address := p.st.GetAddress(symbol)
				symbol = IntTosBin(address)
			} else {
				// symbol is variable
				p.variableAddress = p.variableAddress + 1
				p.st.addEntry(symbol, p.variableAddress)
				symbol = IntTosBin(p.variableAddress)
			}
		} else {
			// symbol is value
			i64, err := strconv.ParseInt(symbol, 10, 0)
			if err != nil {
				log.Fatal(err)
			}
			symbol = strconv.FormatInt(i64, 2)
		}
		return fmt.Sprintf("0%s", padding(symbol, 15, "0"))
	},
	C_COMMAND: func(p *Parser) string {
		return fmt.Sprintf("111%s%s%s", p.Comp(), p.Dest(), p.Jump())
	},
	L_COMMAND: func(p *Parser) string {
		return ""
	},
}

func getType(command string) commandType {
	if strings.HasPrefix(command, "@") {
		return A_COMMAND
	}
	if strings.HasPrefix(command, "(") {
		return L_COMMAND
	}
	return C_COMMAND
}

// must call when command is A_COMMAND
func isSymbol(command string) bool {
	return !regexp.MustCompile(`@\d+`).MatchString(command)
}

type SymbolTable struct {
	in      io.Reader
	scanner *bufio.Scanner
	table   map[string]int
	address int
}

func NewSymbolTable(filename string) *SymbolTable {
	table := map[string]int{
		"R0":     0x0000,
		"R1":     0x0001,
		"R2":     0x0002,
		"R3":     0x0003,
		"R4":     0x0004,
		"R5":     0x0005,
		"R6":     0x0006,
		"R7":     0x0007,
		"R8":     0x0008,
		"R9":     0x0009,
		"R10":    0x000A,
		"R11":    0x000B,
		"R12":    0x000B,
		"R13":    0x000D,
		"R14":    0x000E,
		"R15":    0x000F,
		"SP":     0x0000,
		"LCL":    0x0001,
		"ARG":    0x0002,
		"THIS":   0x0003,
		"THAT":   0x0004,
		"SCREEN": 0x4000,
		"KBD":    0x6000,
	}
	inputFile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(inputFile)
	scanner.Split(bufio.ScanLines)
	return &SymbolTable{
		in:      inputFile,
		scanner: scanner,
		table:   table,
		address: 0,
	}
}

func (s *SymbolTable) buildTable() {
	for s.scanner.Scan() {
		if command := normalize(s.scanner.Text()); command != "" {
			switch getType(command) {
			case NOT_COMMAND:
				break
			case L_COMMAND:
				s.addEntry(getSymbol(command), s.address)
				break
			case A_COMMAND:
				s.address++
				break
			case C_COMMAND:
				s.address++
				break
			}
		}
	}
}

func (s *SymbolTable) addEntry(symbol string, address int) {
	s.table[symbol] = address
}

func (s *SymbolTable) contains(symbol string) bool {
	_, ok := s.table[symbol]
	return ok
}

func (s *SymbolTable) GetAddress(symbol string) int {
	address, ok := s.table[symbol]
	if !ok {
		log.Fatal("get uncantains symbol" + symbol)
	}
	return address
}

// Parser parse jack assembler to machinecode
type Parser struct {
	in              io.Reader
	out             io.Writer
	scanner         *bufio.Scanner
	st              *SymbolTable
	variableAddress int
	CurrentCommand  string
}

func NewParser(filename string) *Parser {
	name, path, err := fileinfo(filename)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(name, path)

	outputFile, err := os.OpenFile(fmt.Sprintf("%s.hack", name), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}

	inputFile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(inputFile)
	scanner.Split(bufio.ScanLines)
	return &Parser{in: inputFile, out: outputFile, scanner: scanner, st: NewSymbolTable(filename), variableAddress: 15}
}

func fileinfo(filename string) (name string, path string, err error) {
	name = strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
	path, err = filepath.Abs(filepath.Dir(filename))
	return
}

func (p *Parser) WriteLines(s string) {
	p.out.Write([]byte(s + "\n"))
}

func (p *Parser) Main() {
	p.buildTable()
	p.Parse()
}

// Scan will build the SymbolTable
func (p *Parser) buildTable() {
	p.st.buildTable()
}

func (p *Parser) ContainsSymbol(symbol string) bool {
	return p.st.contains(symbol)
}

func (p *Parser) Parse() {
	i := 0
	for p.HasMoreCommands() {
		i++
		if ok := p.Advance(); ok {
			//fmt.Println(p.CurrentCommand)
			if code := COMMAND_FUNC[p.CommandType()](p); code != "" {
				if debug {
					p.WriteLines(fmt.Sprintf("%d %s %s", i, code, p.CurrentCommand))
				} else {
					p.WriteLines(code)
				}

			}

		}
	}
}

func (p *Parser) HasMoreCommands() bool {
	return p.scanner.Scan()
}

func normalize(command string) string {
	command = strings.TrimSpace(command)
	if index := strings.Index(command, "//"); index > -1 {
		command = strings.TrimSpace(string([]byte(command)[:index]))
		if command != "" {
			return command
		}
		return ""
	}
	return command
}

func (p *Parser) Advance() bool {
	if command := normalize(p.scanner.Text()); command != "" {
		p.CurrentCommand = command
		return true
	}
	return false
}

func (p *Parser) CommandType() commandType {
	return getType(p.CurrentCommand)
}

func getSymbol(command string) string {
	symbol := strings.TrimPrefix(command, "@")
	symbol = strings.TrimPrefix(symbol, "(")
	symbol = strings.TrimSuffix(symbol, ")")
	return symbol
}

func (p *Parser) Symbol() string {
	return getSymbol(p.CurrentCommand)
}

const (
	M uint = 1 << iota
	D
	A
)

func (p *Parser) isJump() bool {
	return p.CommandType() == C_COMMAND && strings.Contains(p.CurrentCommand, ";")
}

func (p *Parser) dest() string {
	return (strings.Split(p.CurrentCommand, "="))[0]
}

func (p *Parser) Dest() string {
	var r uint = 0
	if !p.isJump() {
		switch p.dest() {
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
	if p.isJump() {
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
