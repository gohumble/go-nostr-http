package httpNostr

import (
	"context"
	"fmt"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip04"
	"log"
	"time"
)

type SubscribeCallback func(message nostr.Event, sub *nostr.Subscription)

func Subscribe(pool *nostr.Relay, filter nostr.Filters, callback SubscribeCallback) {
	sub := pool.Subscribe(context.Background(), filter)
	go func() {
		for event := range sub.Events {
			callback(event, sub)
		}
	}()
}

func Publish(ctx context.Context, content, publicKey string, pool *nostr.Relay) {
	secret, err := nip04.ComputeSharedSecret(Configuration.PrivateKey, publicKey)
	if err != nil {
		panic(err)
	}
	myPublicKey, err := nostr.GetPublicKey(Configuration.PrivateKey)
	if err != nil {
		panic(err)
	}
	tags := make(nostr.Tags, 0)
	tags = append(tags, nostr.Tag{"p", publicKey})
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
	status := pool.Publish(ctx, event)
	if err != nil {
		fmt.Printf("error calling PublishEvent(): %s\n", err.Error())
	}
	log.Printf("Status: %s\n", status.String())
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
