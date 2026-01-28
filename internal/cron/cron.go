package cron

import (
	"context"
	"log"
	"time"

	"github.com/go-co-op/gocron"

	"github.com/azzamdhx/moneybro/backend/internal/services"
)

type Scheduler struct {
	scheduler           *gocron.Scheduler
	notificationService *services.NotificationService
}

func NewScheduler(notificationService *services.NotificationService) *Scheduler {
	s := gocron.NewScheduler(time.UTC)
	return &Scheduler{
		scheduler:           s,
		notificationService: notificationService,
	}
}

func (s *Scheduler) Start() {
	s.scheduler.Every(1).Day().At("08:00").Do(func() {
		log.Println("Running daily notification job...")
		ctx := context.Background()
		if err := s.notificationService.SendDueReminders(ctx); err != nil {
			log.Printf("Error running notification job: %v", err)
		}
		log.Println("Daily notification job completed")
	})

	s.scheduler.StartAsync()
	log.Println("Cron scheduler started")
}

func (s *Scheduler) Stop() {
	s.scheduler.Stop()
	log.Println("Cron scheduler stopped")
}
