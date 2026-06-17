package console

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// reader est le lecteur partagé de l'entrée standard.
var reader = bufio.NewReader(os.Stdin)

// ReadLine lit une ligne complète, nettoyée des espaces et du retour à la ligne.
func ReadLine() string {
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

// ReadInt lit l'entrée jusqu'à obtenir un entier valide.
func ReadInt() int {
	for {
		if n, err := strconv.Atoi(ReadLine()); err == nil {
			return n
		}
		fmt.Print("Entrée invalide, entrez un nombre : ")
	}
}
