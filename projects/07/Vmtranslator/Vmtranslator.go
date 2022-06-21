package main

func main() {

}

type CommandType int

const (
	NOT_COMMAND CommandType = iota
	C_PUSH
	C_POP
	C_LABLE
	C_GOTO
	C_IF
	C_FUNCTION
	C_RETURN
	C_CALL
)

type Parser struct {
}

func (p *Parser) Parser() {
}

func (p *Parser) hasMoreCommands() bool {
	return true
}

func (p *Parser) advance() {

}

func (p *Parser) commandType() CommandType {
	return NOT_COMMAND
}

func (p *Parser) arg1() string {
	return ""
}

func (p *Parser) Arg2() int {
	return 0
}

type CodeWriter struct {
}

func (c *CodeWriter) setFileName(filename string) {

}

func (c *CodeWriter) writeArithmetic(Command string) {

}

func (c *CodeWriter) WritePushPop(CommandType, segment string, index int) {

}

func (c *CodeWriter) Close() {

}
