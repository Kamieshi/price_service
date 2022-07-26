// Package handler  handler for work with gRPC server
package handler

import (
	"context"
	"errors"

	"github.com/Kamieshi/price_service/internal/service"
	"github.com/Kamieshi/price_service/protoc"
	log "github.com/sirupsen/logrus"
)

// PriceServerImplement Implement Price Server
type PriceServerImplement struct {
	protoc.OwnPriceStreamServer
	RedisListener *service.RedisListener
}

// GetPriceStream Create new connection stream with client
func (p *PriceServerImplement) GetPriceStream(_ *protoc.GetPriceStreamRequest, resp protoc.OwnPriceStream_GetPriceStreamServer) error {
	chResp, err := p.RedisListener.ListenChanel(context.Background())
	if err != nil {
		log.WithError(err)
	}
	for m := range chResp {
		err := resp.Send(&protoc.GetPriceStreamResponse{
			Company: &protoc.Company{
				ID:   m.Company.ID,
				Name: m.Company.Name,
			},
			Ask:  m.Ask,
			Bid:  m.Bid,
			Time: m.Time,
		})

		if err != nil {
			log.WithError(err).Info()
			break
		}
	}
	return errors.New("listening was stopped")
}
