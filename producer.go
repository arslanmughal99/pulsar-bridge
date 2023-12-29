package main

import (
	"context"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

var (
	_                 = godotenv.Load()
	maxPendingMsgsStr = os.Getenv("MAX_PENDING_MSGS")
	batchMaxSizeStr   = os.Getenv("BATCHING_MAX_SIZE")
)

// TopicProducer producer with a specific topic
type TopicProducer struct {
	topic    string
	Producer pulsar.Producer
	trx      pulsar.Transaction
	errch    chan error
}

// TopicProducerOptions options for new TopicProducer
type TopicProducerOptions struct {
	// ErrCh error channel
	ErrCh chan error
	// Message topic
	Topic string
	// Pulsar client instance
	Client pulsar.Client
	// Pulsar transaction instance
	Trx pulsar.Transaction
}

// Close flush producer queue and close it
func (tp *TopicProducer) Close() {
	if err := tp.Producer.Flush(); err != nil {
		log.
			Error().
			Err(err).
			Msg("cannot flush producer")
	}
	tp.Producer.Close()
}

// SendMessages publish messages
func (tp *TopicProducer) SendMessages(msgs []string) {
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancelCause(context.Background())
	defer tp.Close()

	for _, m := range msgs {
		wg.Add(1)
		tp.Producer.SendAsync(
			ctx,
			&pulsar.ProducerMessage{
				Payload:     []byte(m),
				Transaction: tp.trx,
			},
			func(id pulsar.MessageID, message *pulsar.ProducerMessage, err error) {
				defer wg.Done()
				log.
					Trace().
					Err(err).
					Stringer("mid", id).
					Str("topic", tp.topic).
					Str("payload", string(message.Payload)).
					Msg("msg report")
				if err != nil {
					tp.errch <- err
					cancel(err)
					close(tp.errch)
				}

			},
		)
	}

	wg.Wait()
}

// NewTopicProducer create a new instance of TopicProducer
func NewTopicProducer(opt *TopicProducerOptions) (*TopicProducer, error) {
	if opt == nil {
		log.
			Panic().
			Msg("topic producer options cannot be nil")
	}

	maxPendingMsgs, err := strconv.ParseInt(maxPendingMsgsStr, 10, 64)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse `MAX_PENDING_MSGS`")
		return nil, err
	}

	batchMaxSize, err := strconv.ParseInt(batchMaxSizeStr, 10, 64)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse `BATCHING_MAX_SIZE`")
		return nil, err
	}

	p, err := opt.Client.CreateProducer(pulsar.ProducerOptions{
		Topic:                   opt.Topic,
		BatchingMaxSize:         uint(batchMaxSize),
		MaxPendingMessages:      int(maxPendingMsgs),
		BatchingMaxPublishDelay: time.Duration(0),
	})
	if err != nil {
		log.
			Error().
			Err(err).
			Str("topic", opt.Topic).
			Msg("failed to initialise producer")
		return nil, errors.Wrap(err, "failed to initialise producer")
	}

	tp := new(TopicProducer)
	tp.Producer = p
	tp.trx = opt.Trx
	tp.errch = opt.ErrCh
	tp.topic = opt.Topic

	return tp, nil
}
