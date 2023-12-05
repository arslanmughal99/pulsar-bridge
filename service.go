package main

import (
	"context"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

var (
	_                = godotenv.Load()
	tnxTimeout int64 = 10
)

// Service instance
type Service struct {
	client *Connection
}

// NewService create new instance of Service
func NewService(client *Connection) *Service {
	tt, err := strconv.ParseInt(os.Getenv("TNX_TIMEOUT_SEC"), 10, 64)

	if err != nil {
		log.Panic().Err(err).Msg("failed to parse transaction timeout.")
	}

	tnxTimeout = tt

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

createTnx:
	tnx, err := s.client.con.NewTransaction(time.Second * time.Duration(tnxTimeout))
	if err != nil {
		if err.Error() == "connection closed" {
			log.Warn().Err(err).Msg("connection dropped re-attempting to connect")
			s.client.Connect()
			goto createTnx
		}

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
