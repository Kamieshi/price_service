package service

import (
	"context"
	"encoding/json"
	"fmt"

	rds "github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"priceService/internal/model"
)

// RedisListener WorkWithRedis
type RedisListener struct {
	Client *rds.Client
}

// ListenChanel Return Chanel, for this Chanel will create new group listener to RedisListener stream
func (r *RedisListener) ListenChanel(ctx context.Context) (chan model.Price, error) {
	nameGroup := uuid.New().String()
	status := r.Client.XGroupCreate(ctx, "prices", nameGroup, "$")
	if status.Err() != nil {
		return nil, status.Err()
	}
	ch := make(chan model.Price)
	args := rds.XReadGroupArgs{
		Group:    nameGroup,
		Consumer: "Consumer",
		Streams:  []string{"prices", ">"},
		Count:    0,
		Block:    0,
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				r.Client.XGroupDestroy(context.Background(), "prices", nameGroup)
				close(ch)
				return
			default:
				resCmd := r.Client.XReadGroup(ctx, &args)
				if resCmd.Err() != nil {
					logrus.WithError(resCmd.Err()).Error()
					continue
				}
				resVal := resCmd.Val()
				for _, comm := range resVal {
					for _, mess := range comm.Messages {
						var pr model.Price
						payLoad := fmt.Sprint(mess.Values["price"])
						err := json.Unmarshal([]byte(payLoad), &pr)
						if err != nil {
							logrus.WithError(err).Error("Error parsing")
							continue
						}

						ch <- pr
					}
				}
			}
		}
	}()
	return ch, nil
}
