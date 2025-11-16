package shopping

import (
	"context"
	"log"
	"time"

	"ecom/service/transaction"
)

func StartTransactionExpireJob(svc transaction.Service) {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			modified, err := svc.RunExpireJob(ctx)
			cancel()

			if err != nil {
				log.Printf("expire job error: %v", err)
				continue
			}
			if modified > 0 {
				log.Printf("expire job: %d transactions expired", modified)
			}
		}
	}()
}
