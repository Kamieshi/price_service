package handlers

import (
	"context"
	"sync"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"priceService/internal/models"
)

// Chanel struct
type Chanel struct {
	addrListeners *Listeners
	nameChanel    string
	chanel        chan *models.Price
}

// Listeners Main listener for management map Chanel
type Listeners struct {
	ChanelPrices chan *models.Price
	Channels     map[string]*Chanel
	sync.RWMutex
}

// Listen Goroutine for each grpc stream connection
func Listen(resp PriceStreaming_HighLoadStreamServer, chanel Chanel) {
	for {
		select {
		case <-resp.Context().(context.Context).Done():
			chanel.addrListeners.DropChanel(chanel.nameChanel)
			close(chanel.chanel)
			return
		case data := <-chanel.chanel:
			err := resp.Send(&GetPriceStreamResponse{
				Company: &Company{
					ID:   data.Company.ID,
					Name: data.Company.Name,
				},
				Ask:  data.Ask,
				Bid:  data.Bid,
				Time: data.Time,
			})
			if err != nil {
				log.WithError(err).Error()
			}
		}
	}
}

// AddChanel Add chanel in Listeners.Channels
func (l *Listeners) AddChanel(chanel *Chanel) {
	l.Lock()
	l.Channels[chanel.nameChanel] = chanel
	l.Unlock()
}

// DropChanel Drop chanel from Listeners.Channels and close chan *models.Price
func (l *Listeners) DropChanel(chName string) {
	l.Lock()
	delete(l.Channels, chName)
	l.Unlock()
}

// StartStream Goroutine for started check Listeners.Channels
func (l *Listeners) StartStream(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.WithError(ctx.Err()).Info()
		case data := <-l.ChanelPrices:
			l.RLock()
			for _, ch := range l.Channels {
				ch.chanel <- data
			}
			l.RUnlock()
		}
	}
}

// PriceStreamingServerImplement Implement gRPC interface PriceStreamingServer
type PriceStreamingServerImplement struct {
	ch        chan *models.Price
	listeners *Listeners
	PriceStreamingServer
}

// NewPriceStreamingServerImplement Constructor and runner goroutine Listeners.StartStream
func NewPriceStreamingServerImplement(ctx context.Context, ch chan *models.Price) *PriceStreamingServerImplement {
	listener := &Listeners{
		ChanelPrices: ch,
		Channels:     make(map[string]*Chanel),
	}
	go listener.StartStream(ctx)
	return &PriceStreamingServerImplement{
		ch:        ch,
		listeners: listener,
	}
}

// HighLoadStream Handler for proto service PriceStreaming.HighLoadStream()
func (p *PriceStreamingServerImplement) HighLoadStream(_ *GetPriceStreamRequest, resp PriceStreaming_HighLoadStreamServer) error {
	ch := Chanel{
		addrListeners: p.listeners,
		nameChanel:    uuid.New().String(),
		chanel:        make(chan *models.Price),
	}
	p.listeners.AddChanel(&ch)
	go Listen(resp, ch)
	for d := range resp.Context().(context.Context).Done() {
		log.Info(d)
		return resp.Context().(context.Context).(context.Context).Err()
	}
	return resp.Context().(context.Context).(context.Context).Err()
}
