package delivery

import (
	"context"
	"fmt"
	"github.com/Nyxox-debug/Cronyx/pkg/cronyx"
)

type ConsoleDelivery struct{}

func (ConsoleDelivery) Deliver(ctx context.Context, target cronyx.DeliveryConfig, files []cronyx.OutputFile) error {
	for _, f := range files {
		fmt.Printf("[Cronyx] Delivered file: %s (%s)\n", f.Name, f.Path)
	}
	return nil
}
