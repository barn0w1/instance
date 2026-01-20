package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/barn0w1/instance/pkg/vast"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Note: .env file not found, using OS environment variables")
	}

	apiKey := os.Getenv("VAST_API_KEY")
	if apiKey == "" {
		log.Fatal("Error: VAST_API_KEY is not set")
	}

	client := vast.NewClient(apiKey)

	// ビルダーパターンで条件を記述
	// 「RTX 4090を搭載し、VRAM24GB以上、信頼性が高く、0.8ドル以下」のものを探す
	query := vast.NewSearch().
		Type(vast.OnDemand).
		GpuName("RTX 4090").
		MinGpus(1).
		MinVRAM(24000). // 24GB
		MinReliability(0.98).
		MaxPrice(0.80).
		Order("dph_total", vast.Asc).
		Limit(10)

	// 実行
	offers, err := client.SearchOffers(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d offers:\n", len(offers))
	for _, o := range offers {
		fmt.Printf("- [ID:%d] %s x%d | RAM: %dMB | Price: $%.3f/hr | Score: %.1f\n",
			o.ID, o.GpuName, o.NumGpus, o.GpuRam, o.DphTotal, o.DlPerf)
	}
}
