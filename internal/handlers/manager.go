package handlers

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"

	"priceService/internal/models"
	"priceService/internal/repository"
)

type StreamingMap struct {
	ch         chan *models.Price
	StreamsMap map[*PriceStreaming_HighLoadStreamServer]PriceStreaming_HighLoadStreamServer
	sync.RWMutex
}

func (s *StreamingMap) Stream() {
	for {
		data := <-s.ch
		connectionFromDelete := make([]*PriceStreaming_HighLoadStreamServer, 0, 100)
		for addr, m := range s.StreamsMap {
			err := m.Send(&GetPriceStreamResponse{
				Company: &Company{
					ID:   data.Company.ID,
					Name: data.Company.Name,
				},
				Ask:  data.Ask,
				Bid:  data.Bid,
				Time: data.Time,
			})
			if err != nil {
				log.WithError(err).Error(" incorrect stream connection ")
				connectionFromDelete = append(connectionFromDelete, addr)
			}
		}
		for _, ad := range connectionFromDelete {
			s.RWMutex.Lock()
			delete(s.StreamsMap, ad)
			s.RWMutex.Unlock()
		}
		if len(s.StreamsMap) == 0 {
			close(s.ch)
			break
		}
	}

}

type ClientManager struct {
	rep           *repository.Redis
	StreamingMaps []*StreamingMap
	sync.RWMutex
}

func NewClientManager(rep *repository.Redis) *ClientManager {
	return &ClientManager{
		rep:           rep,
		StreamingMaps: make([]*StreamingMap, 0),
	}
}

func (c *ClientManager) Add(resp PriceStreaming_HighLoadStreamServer) {
	for _, sm := range c.StreamingMaps {
		if len(sm.StreamsMap) < 100 {
			sm.Lock()
			sm.StreamsMap[&resp] = resp
			sm.Unlock()
			return
		}
	}
	ch, err := c.rep.ListenChanel(context.Background())
	if err != nil {
		log.WithError(err).Fatal()
	}
	c.Lock()
	c.StreamingMaps = append(c.StreamingMaps, &StreamingMap{
		StreamsMap: make(map[*PriceStreaming_HighLoadStreamServer]PriceStreaming_HighLoadStreamServer),
		ch:         ch,
	})
	c.StreamingMaps[len(c.StreamingMaps)-1].StreamsMap[&resp] = resp
	c.Unlock()
	go c.StreamingMaps[len(c.StreamingMaps)-1].Stream()
}
