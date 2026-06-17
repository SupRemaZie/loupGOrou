package console

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)

func ReadLine() string {
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func ReadInt() int {
	for {
		if n, err := strconv.Atoi(ReadLine()); err == nil {
			return n
		}
		fmt.Print("Entrée invalide, entrez un nombre : ")
	}
}
