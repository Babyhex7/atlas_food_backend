package collab

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"

	redisInfra "atlas_food/internal/infra/redis"
)

type PubSubEnvelope struct {
	NodeID  string   `json:"node_id"`
	Message *Message `json:"message"`
}

type RedisPubSubBroker struct {
	nodeID      string
	redisClient *redisInfra.Client
	hub         *Hub
	stopCh      chan struct{}
	once        sync.Once
}

func NewRedisPubSubBroker(hub *Hub) *RedisPubSubBroker {
	return &RedisPubSubBroker{
		nodeID:      uuid.New().String(),
		redisClient: redisInfra.GetInstance(),
		hub:         hub,
		stopCh:      make(chan struct{}),
	}
}

func (b *RedisPubSubBroker) Start(ctx context.Context) {
	if b.redisClient == nil || !b.redisClient.IsAvailable() {
		log.Println("ℹ️ [Redis PubSub] Redis unavailable, running in single-node mode")
		return
	}

	go func() {
		log.Printf("📡 [Redis PubSub] Starting broker for Node %s (Subscribing to pattern atlas:prod:ws:pubsub:room:*)", b.nodeID)
		pubsub := b.redisClient.RawClient().PSubscribe(ctx, "atlas:prod:ws:pubsub:room:*")
		defer pubsub.Close()

		ch := pubsub.Channel()
		for {
			select {
			case msg, ok := <-ch:
				if !ok {
					return
				}
				var envelope PubSubEnvelope
				if err := json.Unmarshal([]byte(msg.Payload), &envelope); err != nil {
					continue
				}

				// Skip messages sent by this same Go node instance to prevent duplicate broadcasts
				if envelope.NodeID == b.nodeID {
					continue
				}

				// Broadcast message to local WebSocket clients in the room
				if envelope.Message != nil {
					b.hub.mu.RLock()
					room, exists := b.hub.rooms[envelope.Message.RoomID]
					b.hub.mu.RUnlock()

					if exists {
						b.hub.broadcastToRoom(room, envelope.Message, nil)
					}
				}

			case <-b.stopCh:
				return
			}
		}
	}()
}

func (b *RedisPubSubBroker) Publish(ctx context.Context, msg *Message) {
	if b.redisClient == nil || !b.redisClient.IsAvailable() || msg == nil {
		return
	}

	go func() {
		channel := fmt.Sprintf("atlas:prod:ws:pubsub:room:%s", msg.RoomID)
		envelope := PubSubEnvelope{
			NodeID:  b.nodeID,
			Message: msg,
		}

		data, err := json.Marshal(envelope)
		if err != nil {
			return
		}

		ctxTimeout, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err := b.redisClient.RawClient().Publish(ctxTimeout, channel, data).Err(); err != nil {
			log.Printf("⚠️ [Redis PubSub] Publish error to %s: %v", channel, err)
		}
	}()
}

func (b *RedisPubSubBroker) Stop() {
	b.once.Do(func() {
		close(b.stopCh)
	})
}
