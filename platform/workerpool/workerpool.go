package workerpool

import (
	"context"
	"folderstructure/internal/errors"
	"log"
	"sync"
)



type Task func(ctx context.Context) error

type WorkerPool struct {
	tasks  chan Task
	wg     sync.WaitGroup
	once   sync.Once
	ctx    context.Context
	cancel context.CancelFunc

	maxWorkers int
	semaphores sync.Map

	started bool
	mu      sync.Mutex
}

func New(maxWorkers int, buffer int) *WorkerPool {
	if maxWorkers <= 0 {
		maxWorkers = 50
	}
	if buffer <= 0 {
		buffer = 1000
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &WorkerPool{
		tasks:      make(chan Task, buffer),
		maxWorkers: maxWorkers,
		ctx:        ctx,
		cancel:     cancel,
	}
}

func (p *WorkerPool) Start() {
	p.once.Do(func() {
		p.mu.Lock()
		p.started = true
		p.mu.Unlock()

		for i := 0; i < p.maxWorkers; i++ {
			p.wg.Add(1)
			go p.worker(i)
		}
	})
}

func (p *WorkerPool) worker(id int) {
	defer p.wg.Done()

	for {
		select {
		case <-p.ctx.Done():
			return

		case task, ok := <-p.tasks:
			if !ok {
				return
			}

			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("[worker-%d] panic recovered: %v", id, r)
					}
				}()

				_ = task(p.ctx)
			}()
		}
	}
}

func (p *WorkerPool) Submit(task Task) error {
	p.mu.Lock()
	if !p.started {
		p.mu.Unlock()
		return errors.ErrPoolNotStarted.New("pool not started")
	}
	p.mu.Unlock()

	select {
	case <-p.ctx.Done():
		return p.ctx.Err()

	case p.tasks <- task:
		return nil
	}
}

func (p *WorkerPool) SubmitWithKey(key string, maxConcurrent int, task Task) error {
	if maxConcurrent <= 0 {
		maxConcurrent = 1
	}

	val, _ := p.semaphores.LoadOrStore(key, make(chan struct{}, maxConcurrent))
	sem := val.(chan struct{})

	return p.Submit(func(ctx context.Context) error {
		select {
		case sem <- struct{}{}:
			defer func() { <-sem }()
		case <-ctx.Done():
			return ctx.Err()
		}

		return task(ctx)
	})
}

func (p *WorkerPool) Stop(ctx context.Context) error {
	p.cancel()

	done := make(chan struct{})

	go func() {
		p.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		close(p.tasks)
		return nil

	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *WorkerPool) Drain() {
	p.wg.Wait()
}
