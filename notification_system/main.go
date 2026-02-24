package main

import (
	"context"
	"fmt"
	"math/rand"
	"notification_system/src/channels"
	"notification_system/src/database"
	pubsub "notification_system/src/infrastructure/pub_sub"
	"notification_system/src/models"
	"notification_system/src/repository"
	"notification_system/src/services"
	"sync"
	"sync/atomic"
	"time"
)

// ─── counters shared across all goroutines ────────────────────────────────────

var (
	delivered  int64
	failed     int64
	dlqCount   int64
	queuedMsgs int64
)

func makeHandler(
	svc channels.ExternalService,
	daRepo repository.DeliveryAttemptRepository,
	notifRepo repository.NotificationRepository,
) pubsub.EventHandler {
	return func(event pubsub.Event) error {
		ne, ok := event.Payload.(models.NotificationEvent)
		if !ok {
			return fmt.Errorf("invalid payload type")
		}
		n := ne.Notification
		ch := ne.Channels[0]

		// Record initial attempt
		attempt := &models.DeliveryAttempt{
			NotificationID: n.ID,
			Channel:        ch,
			Status:         models.DeliveryPending,
		}
		attempt.IncrementAttempt()
		_ = daRepo.Save(attempt)

		// Try to send
		if err := svc.Send(n); err != nil {
			attempt.Status = models.DeliveryFailed
			_ = daRepo.Update(attempt)
			atomic.AddInt64(&failed, 1)
			return err // triggers retry
		}

		// Success
		attempt.Status = models.DeliverySuccess
		_ = daRepo.Update(attempt)
		atomic.AddInt64(&delivered, 1)
		return nil
	}
}

func main() {
	fmt.Println("╔══════════════════════════════════════════════╗")
	fmt.Println("║     1M+ Notification System Simulation       ║")
	fmt.Println("╚══════════════════════════════════════════════╝")

	// 1. Database (optimized string for in-memory SQLite)
	db, err := database.InitDB("file::memory:?cache=shared&mode=memory&_journal_mode=OFF&_synchronous=OFF")
	if err != nil {
		panic(fmt.Sprintf("failed to init DB: %v", err))
	}
	fmt.Println("[init] SQLite DB ready (in-memory, optimized)")

	notifRepo := repository.NewNotificationRepository(db)
	daRepo := repository.NewDeliveryAttemptRepository(db)

	// 2. Channel services
	emailSvc := channels.NewEmailService()
	smsSvc := channels.NewSMSService()
	pushSvc := channels.NewPushService()
	inappSvc := channels.NewInAppService()

	// 3. Pub/Sub Broker + Consumers
	broker := pubsub.NewBroker()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	topics := []struct {
		name string
		svc  channels.ExternalService
	}{
		{"email-notifications", emailSvc},
		{"sms-notifications", smsSvc},
		{"push-notifications", pushSvc},
		{"inapp-notifications", inappSvc},
	}

	for _, t := range topics {
		handler := makeHandler(t.svc, daRepo, notifRepo)
		// 50k buffer, 500 workers per channel to handle massive concurrency
		consumer := pubsub.NewConsumer(t.name, 50000, 500, 3, handler)
		consumer.Start(ctx)
		consumer.StartDLQLogger(ctx, &dlqCount)
		broker.Subscribe(t.name, consumer)
		fmt.Printf("[init] consumer started: %s (workers=50)\n", t.name)
	}

	// 4. Application services
	ruleEngine := services.NewRuleEngine()
	rateLimiter := services.NewRateLimiter(1000, time.Minute) // 5k per user/min
	idempotency := services.NewIdempotencyService()
	notifService := services.NewNotificationService(notifRepo, ruleEngine, rateLimiter, idempotency, broker)

	// 5. Generate 10,000 Users
	fmt.Println("[init] Generating 10,000 users...")
	users := make([]string, 10000)
	prefs := make(map[string]*models.UserPreference)
	allChans := []models.ChannelType{models.ChannelEmail, models.ChannelSMS, models.ChannelPush, models.ChannelInApp}
	for i := 0; i < 10000; i++ {
		uid := fmt.Sprintf("user-%d", i)
		users[i] = uid
		// Randomly enable 2-4 channels per user
		numChans := 2 + rand.Intn(3)
		var userChans []models.ChannelType
		for j := 0; j < numChans; j++ {
			userChans = append(userChans, allChans[rand.Intn(len(allChans))])
		}
		prefs[uid] = &models.UserPreference{UserID: uid, EnabledChannels: userChans}
	}

	// 6. Publish 10,000 notifications via 100 goroutines
	totalEvents := 10000
	eventsPerWorker := totalEvents / 100

	fmt.Printf("\n[sim] Publishing %d notifications...\n", totalEvents)
	startTime := time.Now()

	// Progress ticker
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				q := atomic.LoadInt64(&queuedMsgs)
				d := atomic.LoadInt64(&delivered)
				f := atomic.LoadInt64(&failed)
				dlq := atomic.LoadInt64(&dlqCount)
				fmt.Printf("[progress] Queued: %7d | Delivered: %7d | Failed: %7d | DLQ: %7d\n", q, d, f, dlq)
			case <-ctx.Done():
				return
			}
		}
	}()

	categories := []models.NotificationCategory{
		models.CategoryTransaction,
		models.CategoryMarketing,
		models.CategorySecurity,
		models.CategorySystem,
	}

	var wg sync.WaitGroup
	for w := 0; w < 100; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for i := 0; i < eventsPerWorker; i++ {
				uid := users[rand.Intn(len(users))]
				n := &models.Notification{
					UserID:   uid,
					Category: categories[rand.Intn(len(categories))],
					Title:    fmt.Sprintf("Bulk Event %d-%d", workerID, i),
					Content:  "Load test message content",
					Status:   models.StatusCreated,
					Priority: 1,
				}
				_ = notifService.CreateAndQueue(n, prefs[uid])
				atomic.AddInt64(&queuedMsgs, 1)
			}
		}(w)
	}

	// Wait for all publishers
	wg.Wait()
	publishTime := time.Since(startTime)
	fmt.Printf("\n[sim] ✅ All %d notifications published to broker in %v\n", totalEvents, publishTime)
	fmt.Println("[sim] Waiting for consumers to drain queues... (can take a minute or two)")

	// Wait for queues to drain
	for {
		time.Sleep(1 * time.Second)
		d := atomic.LoadInt64(&delivered)
		dlq := atomic.LoadInt64(&dlqCount)
		if atomic.LoadInt64(&queuedMsgs) == int64(totalEvents) && d+dlq > int64(totalEvents)*2 {
			// Check if progress stalls
			time.Sleep(2 * time.Second)
			d2 := atomic.LoadInt64(&delivered)
			dlq2 := atomic.LoadInt64(&dlqCount)
			if d2 == d && dlq2 == dlq {
				break
			}
		}
	}

	cancel() // stop consumers and ticker

	// Summary
	fmt.Println("\n╔══════════════════════════════════════════════╗")
	fmt.Println("║              1M+ Simulation Summary          ║")
	fmt.Printf("║  Total Queued: %-30d║\n", atomic.LoadInt64(&queuedMsgs))
	fmt.Printf("║  Delivered   : %-30d║\n", atomic.LoadInt64(&delivered))
	fmt.Printf("║  Failed      : %-30d║\n", atomic.LoadInt64(&failed))
	fmt.Printf("║  DLQ         : %-30d║\n", atomic.LoadInt64(&dlqCount))
	fmt.Printf("║  Total Time  : %-30v║\n", time.Since(startTime))
	fmt.Println("╚══════════════════════════════════════════════╝")
}
