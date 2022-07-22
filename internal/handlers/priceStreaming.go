package handlers

import (
	"sync"

	log "github.com/sirupsen/logrus"

	"priceService/internal/models"
)

type StreamingMap struct {
	StreamsMap map[*PriceStreaming_HighLoadStreamServer]PriceStreaming_HighLoadStreamServer
	sync.RWMutex
}

type PriceStreamingServerImplement struct {
	dataChanel chan *models.Price
	PriceStreamingServer
	Streams    []PriceStreaming_HighLoadStreamServer
	StreamsMap *StreamingMap
}

func NewPriceStreamingServerImplement(dtCh chan *models.Price) *PriceStreamingServerImplement {
	streams := make([]PriceStreaming_HighLoadStreamServer, 0)
	strMap := &StreamingMap{
		StreamsMap: make(map[*PriceStreaming_HighLoadStreamServer]PriceStreaming_HighLoadStreamServer),
	}
	return &PriceStreamingServerImplement{
		dataChanel: dtCh,
		Streams:    streams,
		StreamsMap: strMap,
	}
}

func (p *PriceStreamingServerImplement) HighLoadStream(req *GetPriceStreamRequest, resp PriceStreaming_HighLoadStreamServer) error {
	p.StreamsMap.RWMutex.Lock()
	p.StreamsMap.StreamsMap[&resp] = resp
	p.StreamsMap.RWMutex.Unlock()
	for {
		m := <-p.dataChanel
		connectionFromDelete := make([]*PriceStreaming_HighLoadStreamServer, 0)
		p.StreamsMap.RWMutex.Lock()
		for addr, stream := range p.StreamsMap.StreamsMap {
			err := stream.Send(&GetPriceStreamResponse{
				Company: &Company{
					ID:   m.Company.ID,
					Name: m.Company.Name,
				},
				Ask:  m.Ask,
				Bid:  m.Bid,
				Time: m.Time,
			})
			if err != nil {
				log.WithError(err).Error(" incorrect stream connection ")
				connectionFromDelete = append(connectionFromDelete, addr)
			}
		}
		p.StreamsMap.RWMutex.Unlock()
		for _, ad := range connectionFromDelete {
			p.StreamsMap.RWMutex.Lock()
			delete(p.StreamsMap.StreamsMap, ad)
			p.StreamsMap.RWMutex.Unlock()
		}
	}
	return nil
}
