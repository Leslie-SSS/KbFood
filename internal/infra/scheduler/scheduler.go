package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/robfig/cron/v3"
)

// Job represents a scheduled job
type Job interface {
	Name() string
	Run(ctx context.Context) error
}

// Scheduler manages scheduled jobs
type Scheduler struct {
	cron       *cron.Cron
	jobs       map[string]Job
	timeWindow *TimeWindow
	running    bool
	mu         sync.RWMutex
}

// NewScheduler creates a new scheduler
func NewScheduler(timeWindow *TimeWindow) *Scheduler {
	if timeWindow == nil {
		timeWindow = DefaultTimeWindow()
	}

	return &Scheduler{
		cron:       cron.New(cron.WithSeconds()),
		jobs:       make(map[string]Job),
		timeWindow: timeWindow,
		running:    false,
	}
}

// RegisterJob registers a job with a cron expression
func (s *Scheduler) RegisterJob(job Job, cronExpr string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	wrapped := s.wrapJob(job)

	_, err := s.cron.AddFunc(cronExpr, wrapped)
	if err != nil {
		return fmt.Errorf("add job %s: %w", job.Name(), err)
	}

	s.jobs[job.Name()] = job
	log.Info().
		Str("name", job.Name()).
		Str("cron", cronExpr).
		Msg("Job registered")

	return nil
}

// RegisterJobWithTimeWindow registers a job that only runs during time window
func (s *Scheduler) RegisterJobWithTimeWindow(job Job, cronExpr string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	wrapped := s.wrapJobWithTimeWindow(job)

	_, err := s.cron.AddFunc(cronExpr, wrapped)
	if err != nil {
		return fmt.Errorf("add job %s: %w", job.Name(), err)
	}

	s.jobs[job.Name()] = job
	log.Info().
		Str("name", job.Name()).
		Str("cron", cronExpr).
		Str("timeWindow", s.timeWindow.String()).
		Msg("Job registered with time window")

	return nil
}

// wrapJob wraps a job with logging
func (s *Scheduler) wrapJob(job Job) func() {
	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		log.Info().
			Str("job", job.Name()).
			Msg("Job started")

		start := time.Now()

		if err := job.Run(ctx); err != nil {
			log.Error().Err(err).
				Str("job", job.Name()).
				Dur("duration", time.Since(start)).
				Msg("Job failed")
		} else {
			log.Info().
				Str("job", job.Name()).
				Dur("duration", time.Since(start)).
				Msg("Job completed")
		}
	}
}

// wrapJobWithTimeWindow wraps a job with time window check
func (s *Scheduler) wrapJobWithTimeWindow(job Job) func() {
	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		if !s.timeWindow.IsActiveNow() {
			log.Debug().
				Str("job", job.Name()).
				Str("timeWindow", s.timeWindow.String()).
				Msg("Job skipped due to time window")
			return
		}

		log.Info().
			Str("job", job.Name()).
			Msg("Job started")

		start := time.Now()

		if err := job.Run(ctx); err != nil {
			log.Error().Err(err).
				Str("job", job.Name()).
				Dur("duration", time.Since(start)).
				Msg("Job failed")
		} else {
			log.Info().
				Str("job", job.Name()).
				Dur("duration", time.Since(start)).
				Msg("Job completed")
		}
	}
}

// Start starts the scheduler
func (s *Scheduler) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return // Already started
	}

	s.cron.Start()
	s.running = true
	log.Info().Msg("Scheduler started")
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return // Already stopped
	}

	stopCtx := s.cron.Stop()
	<-stopCtx.Done()
	s.running = false
	log.Info().Msg("Scheduler stopped")
}

// Shutdown gracefully shuts down the scheduler
func (s *Scheduler) Shutdown(ctx context.Context) error {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return nil
	}
	s.running = false
	stopCtx := s.cron.Stop()
	s.mu.Unlock()

	select {
	case <-stopCtx.Done():
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// IsRunning checks if the scheduler is running
func (s *Scheduler) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

// GetJobCount returns the number of registered jobs
func (s *Scheduler) GetJobCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.jobs)
}
