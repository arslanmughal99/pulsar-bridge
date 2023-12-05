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
	_             = godotenv.Load()
	tnxTimeout, _ = strconv.ParseInt(os.Getenv("TNX_TIMEOUT_SEC"), 10, 64)
)

// Service instance
type Service struct {
	client *Connection
}

// NewService create new instance of Service
func NewService(client *Connection) *Service {
	s := new(Service)
	s.client = client
	return s
}

// ProduceMessagesTnx publish messages to their respective topics in a single transaction
// if any message failed then no message will be published
func (s *Service) ProduceMessagesTnx(dto *ProducerRequestDto) (*uint, error) {
	if dto == nil {
		log.
			Panic().
			Msg("dto cannot be nil")
	}
	total := new(uint)
	errch := make(chan error)

	var tnx pulsar.Transaction

	tnx, err := s.client.con.NewTransaction(time.Second * time.Duration(tnxTimeout))

	if err != nil {
		log.
			Error().
			Err(err).
			Msg("unable to initiate transaction")
		return nil, err
	}

	go func() {
		wg := sync.WaitGroup{}
		for _, tpm := range dto.Data {
			wg.Add(1)
			*total += uint(len(tpm.Messages))

			p, err := NewTopicProducer(
				&TopicProducerOptions{
					Trx:    tnx,
					ErrCh:  errch,
					Client: s.client.con,
					Topic:  tpm.Topic,
				},
			)
			if err != nil {
				errch <- err
				close(errch)
				wg.Done()
				continue
			}

			go func(tpm TopicMessagesDto) {
				defer wg.Done()
				p.SendMessages(tpm.Messages)
			}(tpm)
		}
		wg.Wait()
		close(errch)
	}()

	err = <-errch

	if err != nil {
		log.Error().Err(err).Msg("operation failed")
		if err := tnx.Abort(context.Background()); err != nil {
			log.Error().Err(err).Msg("failed to rollback transaction")
			return nil, err
		}

		return nil, errors.New("something went wrong")
	}

	log.
		Trace().
		Interface("tnx", tnx.GetTxnID()).
		Msg("committing transaction")
	if err := tnx.Commit(context.Background()); err != nil {
		return nil, err
	}

	log.
		Debug().
		Interface("tnx", tnx.GetTxnID()).
		Msg("transaction committed")

	return total, nil
}
