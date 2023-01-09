package httpNostr

import (
	"bufio"
	"bytes"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip04"
	"io"
	"net/http"
	"net/url"
)

func ReverseProxyHandler(relay *nostr.Relay, target *url.URL) SubscribeCallback {
	return func(message nostr.Event, sub *nostr.Subscription) {
		secret, err := nip04.ComputeSharedSecret(Configuration.PrivateKey, message.PubKey)
		if err != nil {
			panic(err)
		}
		r, err := nip04.Decrypt(message.Content, secret)
		if err != nil {
			panic(err)
		}
		buf := &bytes.Buffer{}
		buf.WriteString(r)
		request, err := http.ReadRequest(bufio.NewReader(buf))
		if err != nil {
			panic(err)
		}
		request.RequestURI = ""
		request.URL = target
		c := http.Client{}
		res, err := c.Do(request)
		if err != nil {
			panic(err)
		}
		response, err := io.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}
		Publish(request.Context(), string(response), message.PubKey, relay)
	}
}
