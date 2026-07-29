package usecase

import (
	"context"
	"fmt"
	"log"
	"time"
)

func compensateWithRetry(ctx context.Context, maxAttempts int, delay time.Duration, fn func(context.Context)error)error{
	var err error
	for attempt := 1; attempt <= maxAttempts; attempt++{
		err = fn(ctx)
		if err == nil{
			return nil
		}

		log.Printf("compensação falhou na tentativa %d/%d. erro: %v. retentando em %v...", attempt, maxAttempts, err, delay)

		if attempt == maxAttempts{
			break
		}

		timer := time.NewTicker(delay)
		defer timer.Stop()
		select {
		case  <-ctx.Done():
			return fmt.Errorf("contexto cancelado durante o retry: %w", ctx.Err())
		case <-timer.C:
		}	
	}

	return fmt.Errorf("esgotadas as %d tentativas. Último erro: %w", maxAttempts, err)
}