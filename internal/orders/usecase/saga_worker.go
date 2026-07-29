package usecase

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/aclgo/simple-api-gateway/internal/orders"
)


type SagaWorker struct {
	taskArray []*orders.CompensationTask
	maxAttempts int
	delay time .Duration
	mu sync.Mutex
}

func NewSagaWorker(maxAttempts int, delay time.Duration)*SagaWorker{
	return &SagaWorker{
		taskArray: make([]*orders.CompensationTask, 0),
	}
}

func(s *SagaWorker)Start(ctx context.Context){

	ticker := time.NewTicker(5 *time.Second)
	defer ticker.Stop()

	go func() {
		for{
			select{
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.proccessNext(ctx)
			}
		}
	}()
}

func (s *SagaWorker)AppendTask(task *orders.CompensationTask){
	s.mu.Lock()
	defer s.mu.Unlock()
	s.taskArray = append(s.taskArray, task)
}

func(s *SagaWorker)proccessNext(ctx context.Context){
	s.mu.Lock()

	if len(s.taskArray) == 0  {
		s.mu.Unlock()
		return
	}

	task := s.taskArray[0]

	s.taskArray = s.taskArray[1:]

	s.mu.Unlock()

	s.processTask(ctx, task)
}

func(s *SagaWorker)processTask(ctx context.Context, task *orders.CompensationTask){
	for i := len(task.Compensations)-1;i>=0;i--{
		compFn := task.Compensations[i]

		retryErr := compensateWithRetry(ctx, s.maxAttempts,s.delay,compFn)
		if retryErr != nil {
			log.Printf("[CRITICAL ALARM] worker falhou permanentemente na compensação! erro original: %v | erro compensação: %v", task.OriginalErr, retryErr)
		}
	}
}