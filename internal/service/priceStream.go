// Package service work with CommonPriceStream
package service

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"

	"priceService/internal/models"
)

// Chanel struct
type Chanel struct {
	AddrListeners *Listeners
	NameChanel    string
	Chanel        chan models.Price
}

// Listeners Main listener for management map Chanel
type Listeners struct {
	ChanelPrices chan models.Price
	Channels     map[string]*Chanel
	sMutex       sync.RWMutex
}

// AddChanel Add Chanel in Listeners.Channels
func (l *Listeners) AddChanel(chanel *Chanel) {
	l.sMutex.Lock()
	l.Channels[chanel.NameChanel] = chanel
	l.sMutex.Unlock()
}

// DropChanel Drop Chanel from Listeners.Channels and close chan *models.Price
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
