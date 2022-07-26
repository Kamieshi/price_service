// Package service work with CommonPriceStream
package service

import (
	"container/list"
	"context"
	"sync"

	log "github.com/sirupsen/logrus"

	"priceService/internal/models"
)

// Chanel struct
type Chanel struct {
	Disabled bool
	Chanel   chan models.Price
	sync.RWMutex
}

// Listeners Main listener for management map Chanel
type Listeners struct {
	ChanelPrices chan models.Price
	Channels     *list.List
	sMutex       sync.RWMutex
}

// AddChanel Add Chanel in Listeners.Channels
func (l *Listeners) AddChanel(chanel *Chanel) {
	l.sMutex.Lock()
	l.Channels.PushBack(chanel)
	l.sMutex.Unlock()
}

// StartStream Goroutine for started check Listeners.Channels
func (l *Listeners) StartStream(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.WithError(ctx.Err()).Info()
		case data := <-l.ChanelPrices:
			for ch := l.Channels.Front(); ch != nil; ch = ch.Next() {
				ch.Value.(*Chanel).RLock()
				if ch.Value.(*Chanel).Disabled {
					ch.Value.(*Chanel).RUnlock()
					ch.Value.(*Chanel).Lock()
					l.Channels.Remove(ch)
					ch.Value.(*Chanel).Unlock()
					continue
				}
				ch.Value.(*Chanel).Chanel <- data
				ch.Value.(*Chanel).RUnlock()
			}
		}
	}
}
