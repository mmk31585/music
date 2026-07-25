package sentry

import (
	"time"

	sentrygo "github.com/getsentry/sentry-go"
)

func Init(dsn, environment string) error {
	return sentrygo.Init(sentrygo.ClientOptions{
		Dsn:              dsn,
		Environment:      environment,
		TracesSampleRate: 0.1,
	})
}

func Recover() {
	sentrygo.Recover()
}

func Flush() {
	sentrygo.Flush(2 * time.Second)
}