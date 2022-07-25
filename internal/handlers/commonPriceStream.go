package handlers

import (
	"context"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"priceService/internal/models"
	"priceService/internal/service"
	"priceService/protoc"
)

// CommonPriceStreamServerImplement Implement gRPC interface PriceStreamingServer
type CommonPriceStreamServerImplement struct {
	ch        chan models.Price
	listeners *service.Listeners
	protoc.CommonPriceStreamServer
}

// NewCommonPriceStreamServerImplement Constructor and runner goroutine Listeners.StartStream
func NewCommonPriceStreamServerImplement(ctx context.Context, ch chan models.Price) *CommonPriceStreamServerImplement {
	listener := &service.Listeners{
		ChanelPrices: ch,
		Channels:     make(map[string]*service.Chanel),
	}
	go listener.StartStream(ctx)
	return &CommonPriceStreamServerImplement{
		ch:        ch,
		listeners: listener,
	}
}

// GetPriceStream Handler for proto service CommonPriceStream.GetPriceStream() For this connection stream will not create
// new connection to Redis, There is only one common connection(group reader) to redis, This type price stream is slower
// than OwnPriceStream, but you don't have any count limit like count connection to one instance price service.
func (p *CommonPriceStreamServerImplement) GetPriceStream(_ *protoc.GetPriceStreamRequest, resp protoc.CommonPriceStream_GetPriceStreamServer) error {
	ch := service.Chanel{
		AddrListeners: p.listeners,
		NameChanel:    uuid.New().String(),
		Chanel:        make(chan models.Price),
	}
	p.listeners.AddChanel(&ch)

	for {
		select {
		case <-resp.Context().Done():
			log.WithError(resp.Context().Err()).Info()
			p.listeners.DropChanel(ch.NameChanel)
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
