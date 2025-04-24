# avrasm

An AVR assembler I work on in my spare time.
It currently supports only a small subset of the AVR instruction set.


## Example Program

This program uses all features the assembler currently supports:

```asm
.byte 1, 2, 3
    nop
    mov r1, r2 ; hello world!
```
## How to Build
```sh
git clone https://github.com/jochemarends/avrasm.git
cd avrasm
go build -o avrasm cmd/main.go
```

