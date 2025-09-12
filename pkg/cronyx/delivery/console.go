package delivery

import (
	"context"
	"fmt"

	"github.com/Nyxox-debug/Cronyx/pkg/cronyx"
)

type ConsoleDelivery struct{}

func (c ConsoleDelivery) Deliver(ctx context.Context, target cronyx.DeliveryConfig, files []cronyx.OutputFile) error {
	fmt.Println("=== DELIVERY: Console Output ===")

	for _, file := range files {
		fmt.Printf("Generated file: %s\n", file.Name)
		fmt.Printf("Path: %s\n", file.Path)
		fmt.Printf("Size: %d bytes\n", len(file.Data))

		// Optionally print first few lines of content
		if len(file.Data) > 0 {
			content := string(file.Data)
			if len(content) > 500 {
				content = content[:500] + "..."
			}
			fmt.Println("Content preview:")
			fmt.Println("---")
			fmt.Println(content)
			fmt.Println("---")
		}
		fmt.Println()
	}

	return nil
}
