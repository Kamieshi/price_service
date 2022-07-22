// Package handlers'  handler for work with gRPC server
package handlers

import (
	"context"
	"errors"

	log "github.com/sirupsen/logrus"

	"priceService/internal/repository"
)

// PriceServerImplement Implement Price Server
type PriceServerImplement struct {
	PriceServer
	Rep *repository.Redis
}

// GetPriceStream Create new connection stream with client
func (p *PriceServerImplement) GetPriceStream(req *GetPriceStreamRequest, resp Price_GetPriceStreamServer) error {
	chResp, err := p.Rep.ListenChanel(context.Background())
	if err != nil {
		log.WithError(err)
	}
	for m := range chResp {
		err := resp.Send(&GetPriceStreamResponse{
			Company: &Company{
				ID:   m.Company.ID,
				Name: m.Company.Name,
			},
			Ask:  m.Ask,
			Bid:  m.Bid,
			Time: m.Time,
		})

		if err != nil {
			log.WithError(err).Info(req.sizeCache)
			break
		}
	}
	return errors.New("listening was stopped")
}
