# go-nostr-http 
This library will enable sending and receiving http requests and responses in golang applications trough nostr relay servers.

We make use of `EncryptedDirectMessage` event types to send encrypted requests to any nostr public key.

Listening services are able to subscribe to requests directed to their public key. This will make services publicly available even if the 
host server is behind a NAT. It should also increase privacy, because we evade direct TCP/IP connections between clients and servers.

## Installation
```
go get github.com/gohumble/go-nostr-http
```
## Example
### Client
Clients are http.Client's with a special nostr transport layer. Therefor the integration to any client application should be flawless and minimally invasive.
It's also possible to integrate go-nostr-http into third party http client libraries like https://github.com/imroc/req
```go
import (
    "github.com/gohumble/go-nostr-http"
    "github.com/nbd-wtf/go-nostr"
)
// ...
relay, err := nostr.RelayConnect(context.Background(), "wss://nostr.relay.com")
if err != nil {
    panic(err)
}
client := httpNostr.NewClient(relay,"myNostrPublicKey")
req, err := http.NewRequest("GET", "", nil)
if err != nil {
    panic(err)
}
req.Header.Set("TO-NOSTR-PUBLIC-KEY", "recipient-nostr-public-key")
response,err := WalletClient.client.Do(req)
```

### Server 
This server will transform incoming httpNostr requests to valid http requests to his localhost (nostr reverse proxy to local api).

```go 
import (
    "github.com/gohumble/go-nostr-http"
    "github.com/nbd-wtf/go-nostr"
)
relay, err := nostr.RelayConnect(context.Background(), "wss://nostr.relay.com")
if err != nil {
    panic(err)
}
localhost, err := url.Parse("http://localhost:3338")
if err != nil {
    panic(err)
}
httpNostr.Subscribe(relay, httpNostr.GetSubscriptionFilter(mintNostrPublicKey), httpNostr.ReverseProxyHandler(m.Nostr, localhost))
```

