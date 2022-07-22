package handlers

import (
	"context"
	"sync"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"priceService/internal/models"
	"priceService/internal/repository"
)

type StreamingMap struct {
	owner      *ClientManager
	name       string
	ch         chan *models.Price
	StreamsMap map[*PriceStreaming_HighLoadStreamServer]PriceStreaming_HighLoadStreamServer
	sync.RWMutex
}

func (s *StreamingMap) Stream(cancelFunc context.CancelFunc) {
	for {
		data, ok := <-s.ch
		if !ok {
			log.Fatalf("sad")
		}
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
				log.WithError(err).Info(" incorrect stream connection ")
				connectionFromDelete = append(connectionFromDelete, addr)
			}
		}
		for _, ad := range connectionFromDelete {
			s.RWMutex.Lock()
			delete(s.StreamsMap, ad)
			s.RWMutex.Unlock()
		}
		if len(s.StreamsMap) == 0 {
			cancelFunc()
			defer func() {
				delete(s.owner.StreamingMaps, s.name)
			}()
			break
		}
	}

}

type ClientManager struct {
	rep           *repository.Redis
	StreamingMaps map[string]*StreamingMap
	sync.RWMutex
}

func NewClientManager(rep *repository.Redis) *ClientManager {
	return &ClientManager{
		rep:           rep,
		StreamingMaps: make(map[string]*StreamingMap, 0),
	}
}

func (c *ClientManager) Add(resp PriceStreaming_HighLoadStreamServer) {
	for _, sm := range c.StreamingMaps {
		if len(sm.StreamsMap) < 500 {
			sm.Lock()
			sm.StreamsMap[&resp] = resp
			sm.Unlock()
			return
		}
	}
	ctx, cFunc := context.WithCancel(context.Background())
	ch, err := c.rep.ListenChanel(ctx)
	if err != nil {
		log.WithError(err).Fatal()
	}
	c.Lock()
	newName := uuid.New().String()
	c.StreamingMaps[newName] = &StreamingMap{
		owner:      c,
		name:       newName,
		StreamsMap: make(map[*PriceStreaming_HighLoadStreamServer]PriceStreaming_HighLoadStreamServer),
		ch:         ch,
	}
	c.StreamingMaps[newName].StreamsMap[&resp] = resp
	c.Unlock()
	go c.StreamingMaps[newName].Stream(cFunc)
}
