package httpNostr

import (
	"bytes"
	"fmt"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip04"
	"io"
	"net/http"
	"net/http/httputil"
	"sync"
)

type Client struct {
	relay           *nostr.Relay
	clientPublicKey string
}

func NewClient(relay *nostr.Relay, publicKey string) *Client {
	return &Client{
		relay:           relay,
		clientPublicKey: publicKey,
	}
}
func (nc *Client) RoundTrip(r *http.Request) (*http.Response, error) {
	toPublicKey := r.Header.Get("NOSTR-TO-PUBLIC-KEY")
	if toPublicKey == "" {
		return nil, fmt.Errorf("please set NOSTR-TO-PUBLIC-KEY header")
	}
	request, err := httputil.DumpRequestOut(r, true)
	if err != nil {
		return nil, err
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	var response http.Response
	Subscribe(nc.relay, GetSubscriptionFilter(nc.clientPublicKey), func(message nostr.Event, sub *nostr.Subscription) {
		if message.Tags.ContainsAny("p", []string{nc.clientPublicKey}) {
			rs, err := nip04.ComputeSharedSecret(Configuration.PrivateKey, message.PubKey)
			if err != nil {
				return
			}
			resp, err := nip04.Decrypt(message.Content, rs)
			if err != nil {
				return
			}
			response = http.Response{
				Body:          io.NopCloser(bytes.NewBufferString(resp)),
				Status:        "200 OK",
				StatusCode:    200,
				Proto:         "HTTP/1.1",
				ProtoMajor:    1,
				ProtoMinor:    1,
				ContentLength: int64(len(message.Content)),
				Request:       r,
				Header:        make(http.Header, 0),
			}
			sub.Unsub()
			wg.Done()
		}
	})
	Publish(r.Context(), string(request), r.Host, nc.relay)
	wg.Wait()
	return &response, nil
}
