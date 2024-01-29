package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/GDGVIT/configgy-backend/mailer/pkg"
	"github.com/GDGVIT/configgy-backend/pkg/logger"
	messagebroker "github.com/GDGVIT/configgy-backend/pkg/message_broker"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var numWorkers int // Number of worker goroutines for parallel processing

var RootCmd = &cobra.Command{
	Use:   "mailer",
	Short: "Start the mailer service",
	Run: func(cmd *cobra.Command, args []string) {
		logger := logger.GetInstance()
		godotenv.Load()
		ctx, cancel := context.WithCancel(context.Background())

		// Setup a signal channel to handle graceful shutdown
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			select {
			case sig := <-sigCh:
				logger.Infof("Received signal %v. Shutting down gracefully...", sig)
				cancel()
			case <-ctx.Done():
			}
		}()

		var mailer pkg.MailInstance

		mailer.Logger = logger
		smtpUsername := os.Getenv("SMTP_USERNAME")
		smtpPassword := os.Getenv("SMTP_PASSWORD")
		smtpHost := os.Getenv("SMTP_HOST")
		var smtpPort int
		fmt.Sscanf(os.Getenv("SMTP_PORT"), "%d", &smtpPort)

		mailer.SetCredentials(smtpUsername, smtpPassword)
		mailer.SetTransportDetails(smtpHost, smtpPort)

		var worker pkg.WorkerService
		worker.QueueName = "mail"
		worker.NumOfWorkers = numWorkers

		rabbitMq, err := messagebroker.NewRabbitMQHelper(os.Getenv("RABBITMQ_HOST_PORT"), worker.NumOfWorkers, logger)
		if err != nil {
			logger.Error(errors.Wrap(err, "failed to initialize RabbitMQ"))
			cancel()
		}
		defer rabbitMq.Close()

		service := pkg.NewService(logger, &mailer, &worker, rabbitMq)
		service.StartConsumer(ctx)

		// Wait for the context to be canceled (e.g., on receiving SIGINT or SIGTERM)
		<-ctx.Done()
	},
}

func init() {
	RootCmd.PersistentFlags().IntVarP(&numWorkers, "workers", "w", 1, "Number of worker goroutines")
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
