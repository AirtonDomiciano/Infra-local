package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/AirtonDomiciano/envscout/internal/envfile"
	"github.com/AirtonDomiciano/envscout/internal/scanner"
)

func main() {
	base := flag.String("base", "192.168.1.", "Base da rede (ex: 192.168.1.)")
	port := flag.Int("port", 1433, "Porta para testar (ex: 1433)")
	timeoutMs := flag.Int("timeout", 300, "Timeout por host em ms")
	concurrency := flag.Int("c", 120, "Concorrência (goroutines simultâneas)")
	envPath := flag.String("env", "", "Caminho do .env para atualizar (opcional)")
	key := flag.String("key", "SERVER_DB", "Chave do .env para atualizar")
	prefer := flag.String("prefer", "", "Preferir um IP específico (ex: 192.168.1.7) se estiver disponível")

	flag.Parse()

	results := scanner.ScanOpenPort(*base, *port, *timeoutMs, *concurrency)

	if len(results) == 0 {
		fmt.Println("❌ Nenhum host com a porta aberta encontrado.")
		os.Exit(1)
	}

	chosen := results[0]
	if *prefer != "" {
		for _, ip := range results {
			if ip == *prefer {
				chosen = ip
				break
			}
		}
	}

	fmt.Println("✅ Hosts encontrados:", results)
	fmt.Println("🎯 Escolhido:", chosen)

	if *envPath != "" {
		if err := envfile.UpsertKey(*envPath, *key, chosen); err != nil {
			fmt.Println("❌ Falha ao atualizar .env:", err)
			os.Exit(1)
		}
		fmt.Printf("📝 Atualizado %s em %s => %s=%s\n", *key, *envPath, *key, chosen)
	}
}
