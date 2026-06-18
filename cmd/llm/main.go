package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/SupRemaZie/loupGOrou/internal/orchestrator"
)

func main() {
	host := os.Getenv("OLLAMA_HOST")
	if host == "" {
		host = "http://localhost:11434"
	}
	model := os.Getenv("OLLAMA_MODEL")
	if model == "" {
		model = "llama3.1"
	}
	fmt.Printf("Ollama: %s — modèle: %s\n", host, model)

	players := []string{"Alice", "Bob", "Charlie", "Diana", "Eve", "Frank"}
	nbWolves := 2

	state := orchestrator.Init("partie-llm-1", players, nbWolves)

	fmt.Println("=== Loup-Garou LLM ===")
	fmt.Printf("Joueurs (%d), loups (%d):\n", len(players), nbWolves)
	for _, p := range state.Players {
		fmt.Printf("  %-10s %s\n", p.Name, p.Role)
	}

	result, err := orchestrator.Run(context.Background(), state)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n=== Fin de partie ===\n")
	fmt.Printf("Vainqueur : %s\n\n", result.Result)
	fmt.Println("Survivants:")
	for _, p := range result.Players {
		if p.IsAlive {
			fmt.Printf("  %-10s %s\n", p.Name, p.Role)
		}
	}
	fmt.Println("Elimines:")
	for _, p := range result.Players {
		if !p.IsAlive {
			fmt.Printf("  %-10s %s\n", p.Name, p.Role)
		}
	}
}
