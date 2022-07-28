// Package service work with CommonPriceStream
package service

import (
	"context"
	"sync"

	"github.com/Kamieshi/price_service/internal/model"
	log "github.com/sirupsen/logrus"
)

// PoolListeners Pool from common handler
type PoolListeners struct {
	ChanelPrices  chan model.Price
	PullListeners []*Listeners
	CountChInPool int
}

// NewPoolListeners constructor
func NewPoolListeners(ctx context.Context, ch chan model.Price, countChInPool int) *PoolListeners {
	pools := make([]Listeners, 100, 100)
	poolsL := make([]*Listeners, 100, 100)
	for i := range pools {
		chListener := make(chan model.Price)
		pools[i].ChanelPrices = chListener
		pools[i].Channels = make(map[string]*Chanel)
		poolsL[i] = &pools[i]
		go poolsL[i].StartStream(ctx)
	}

	return &PoolListeners{
		PullListeners: poolsL,
		CountChInPool: countChInPool,
		ChanelPrices:  ch,
	}
}

// StartListenPool constructor
func (p *PoolListeners) StartListenPool(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case data := <-p.ChanelPrices:
			for _, listener := range p.PullListeners {
				listener.ChanelPrices <- data
			}
		}
	}
}

func (p *PoolListeners) GetListener() *Listeners {
	for {
		for _, list := range p.PullListeners {
			if len(list.Channels) <= p.CountChInPool {
				return list
			}
		}
	}
}

// Chanel struct
type Chanel struct {
	AddrListeners *Listeners
	NameChanel    string
	Chanel        chan model.Price
}

// Listeners Main listener for management map Chanel
type Listeners struct {
	ChanelPrices chan model.Price
	Channels     map[string]*Chanel
	sMutex       sync.RWMutex
}

// AddChanel Add Chanel in Listeners.Channels
func (l *Listeners) AddChanel(chanel *Chanel) {
	l.sMutex.Lock()
	l.Channels[chanel.NameChanel] = chanel
	l.sMutex.Unlock()
}

// DropChanel Drop Chanel from Listeners.Channels and close chan *model.Price
func (l *Listeners) DropChanel(chName string) {
	l.sMutex.Lock()
	delete(l.Channels, chName)
	l.sMutex.Unlock()
}

// StartStream Goroutine for started check Listeners.Channels
func (l *Listeners) StartStream(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.WithError(ctx.Err()).Info()
		case data := <-l.ChanelPrices:
			l.sMutex.RLock()
			for _, ch := range l.Channels {
				ch.Chanel <- data
			}
			l.sMutex.RUnlock()
		}
	}
}
