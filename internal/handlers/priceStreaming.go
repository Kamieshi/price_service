package handlers

import (
	"context"

	"priceService/internal/repository"
)

type PriceStreamingServerImplement struct {
	ctx           context.Context
	clientManager *ClientManager
	PriceStreamingServer
}

func NewPriceStreamingServerImplement(rep *repository.Redis) *PriceStreamingServerImplement {
	clMng := NewClientManager(rep)
	return &PriceStreamingServerImplement{
		ctx:           context.Background(),
		clientManager: clMng,
	}
}

func (p *PriceStreamingServerImplement) HighLoadStream(req *GetPriceStreamRequest, resp PriceStreaming_HighLoadStreamServer) error {
	p.clientManager.Add(resp)
	<-resp.Context().Done()
	return resp.Context().Err()
}
