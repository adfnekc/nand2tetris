// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Fill.asm

// Runs an infinite loop that listens to the keyboard input.
// When a key is pressed (any key), the program blackens the screen,
// i.e. writes "black" in every pixel;
// the screen should remain fully black as long as the key is pressed. 
// When no key is pressed, the program clears the screen, i.e. writes
// "white" in every pixel;
// the screen should remain fully clear as long as no key is pressed.

// Put your code here.

(LOOP)
    @KBD
    D=M
    @white
    D;JEQ
    @black
    0;JMP
    @LOOP
    0;JMP


(white)
    @colur
    M=0
    @draw
    0;JMP

(black)
    @colur
    M=-1
    @draw
    0;JMP

(draw)
    @SCREEN
    D=A
    @pos
    M=D
    @KBD
    D=A
    @last_pos
    M=D

    (drwa_loop)
        @last_pos
        D=M
        @pos
        D=M-D
        @LOOP
        D;JEQ

        @colur
        D=M
        @pos
        A=M
        M=D
        @pos
        M=M+1
        @drwa_loop
        0;JMP
