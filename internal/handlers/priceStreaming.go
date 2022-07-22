package handlers

import (
	"time"

	"priceService/internal/repository"
)

type PriceStreamingServerImplement struct {
	clientManager *ClientManager
	PriceStreamingServer
}

func NewPriceStreamingServerImplement(rep *repository.Redis) *PriceStreamingServerImplement {
	clMng := NewClientManager(rep)
	return &PriceStreamingServerImplement{
		clientManager: clMng,
	}
}

func (p *PriceStreamingServerImplement) HighLoadStream(req *GetPriceStreamRequest, resp PriceStreaming_HighLoadStreamServer) error {
	p.clientManager.Add(resp)
	time.Sleep(1 * time.Minute)
	return nil
}
