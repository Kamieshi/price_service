package handler

import (
	"context"

	"github.com/Kamieshi/price_service/internal/model"
	"github.com/Kamieshi/price_service/internal/service"
	"github.com/Kamieshi/price_service/protoc"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// CommonPriceStreamServerImplement Implement gRPC interface PriceStreamingServer
type CommonPriceStreamServerImplement struct {
	poolListeners *service.PoolListeners
	protoc.CommonPriceStreamServer
}

// NewCommonPriceStreamServerImplement Constructor and runner goroutine Listeners.StartStream
func NewCommonPriceStreamServerImplement(ctx context.Context, ch chan model.Price) *CommonPriceStreamServerImplement {
	poolListeners := service.NewPoolListeners(ctx, ch, 100)
	go poolListeners.StartListenPool(ctx)
	return &CommonPriceStreamServerImplement{
		poolListeners: poolListeners,
	}
}

// GetPriceStream Handler for proto service CommonPriceStream.GetPriceStream() For this connection stream will not create
// new connection to Redis, There is only one common connection(group reader) to redis, This type price stream is slower
// than OwnPriceStream, but you don't have any count limit like count connection to one instance price service.
func (p *CommonPriceStreamServerImplement) GetPriceStream(_ *protoc.GetPriceStreamRequest, resp protoc.CommonPriceStream_GetPriceStreamServer) error {
	listener := p.poolListeners.GetListener()

	ch := service.Chanel{
		AddrListeners: listener,
		NameChanel:    uuid.New().String(),
		Chanel:        make(chan model.Price),
	}
	listener.AddChanel(&ch)
	for {
		select {
		case <-resp.Context().Done():
			log.WithError(resp.Context().Err()).Info()
			listener.DropChanel(ch.NameChanel)
			return resp.Context().Err()
		case data := <-ch.Chanel:
			err := resp.Send(&protoc.GetPriceStreamResponse{
				Company: &protoc.Company{
					ID:   data.Company.ID,
					Name: data.Company.Name,
				},
				Ask:  data.Ask,
				Bid:  data.Bid,
				Time: data.Time,
			})
			if err != nil {
				log.WithError(err).Error()
			}
		}
	}
}
