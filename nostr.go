package httpNostr

import (
	"context"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip04"
	"time"
)

type SubscribeCallback func(message nostr.Event, sub *nostr.Subscription)

func Subscribe(ctx context.Context, relay *nostr.Relay, filter nostr.Filters, callback SubscribeCallback) {
	sub := relay.Subscribe(ctx, filter)
	go func() {
		for {
			select {
			case event := <-sub.Events:
				go callback(event, sub)
			case <-ctx.Done():
				return
			default:
				time.Sleep(time.Millisecond * 25)
			}
		}
	}()
}

func Publish(ctx context.Context, content, toPublicKey string, relay *nostr.Relay) {
	secret, err := nip04.ComputeSharedSecret(Configuration.PrivateKey, toPublicKey)
	if err != nil {
		panic(err)
	}
	myPublicKey, err := nostr.GetPublicKey(Configuration.PrivateKey)
	if err != nil {
		panic(err)
	}
	tags := make(nostr.Tags, 0)
	tags = append(tags, nostr.Tag{"p", toPublicKey})
	msg, err := nip04.Encrypt(content, secret)
	event := nostr.Event{
		CreatedAt: time.Now(),
		Kind:      nostr.KindEncryptedDirectMessage,
		Tags:      tags,
		PubKey:    myPublicKey,
		Content:   msg,
	}
	err = event.Sign(Configuration.PrivateKey)
	if err != nil {
		panic(err)
	}
	relay.Publish(ctx, event)
}

func GetSubscriptionFilter(toPublicKey string) nostr.Filters {
	t := time.Now()
	subscriptionsTags := make(nostr.TagMap, 0)
	if toPublicKey != "" {
		subscriptionsTags["p"] = nostr.Tag{toPublicKey}
	}
	return nostr.Filters{
		{
			Tags:  subscriptionsTags,
			Kinds: []int{nostr.KindEncryptedDirectMessage},
			Since: &t,
		}}
}
