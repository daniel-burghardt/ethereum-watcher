# Ethereum Watcher

Ethereum Watcher is a service built in Go to track transactions from subscribed addresses on the Ethereum network.

It provides an API for easy subscribing and transactions fetching.

Additionally, a webhook URL can be configured to receive live notifications about new transactions pertaining to the subscribed addresses.

## API

### Configuration

The following environment variables can be configured:

- `ETH_SERVER_URL` (optional): specified which RPC node to connect to. Defaults to `https://cloudflare-eth.com`;
- `WEBHOOK_URL` (optional): when provided the service will make a post request to this URL informing about new transactions for one or more subscribed addresses;

### Endpoints

#### `GET /current-block`

Returns the latest block the service has ingested.

#### `POST /subscribe/{address}`

Subscribes `{address}` to the watcher service. Once subscribed, its transactions will be tracked and persistent in the service's storage.

#### `GET /transactions/{address}`

Returns an array of all transactions sent to or received by `{address}` after the time of subscription.

## Running

### Requirements

- `Go 1.22` (or higher) installed _(this is essential for the API to function properly)_;

### Command Line

The service can be started with `go run .` on the root folder.

If you want to connect to an alternative RCP node (such as a testnet one) the environment variable can be set with `export ETH_SERVER_URL=https://rpc.sepolia.org` before running `go run .`.

In the same fashion, the webhook URL can be set with `export WEBHOOK_URL=https://myserver/webhook`.

## Considerations

### Performance

The current implementation uses the `eth_blockNumber` method to fetch the latest block from the ledger. However, this method presented some lag in the tests conducted, taking a few seconds after the closing time of the latest block to return the updated value.

An alternative version calling `eth_getBlockByNumber` directly until the response is not empty is implemented on separate branch `performance-improvement`. This approach presented a significant performance improvement (delays below 1s) over the current one, but it relies on RPC method calls to error, which resulted in a different response objects in the tests conducted (with `https://rpc.sepolia.org` and `https://cloudflare-eth.com`) depending on the RPC node. Although promising this approach requires some further experimentations.

A third alternative was also tested using `eth_newBlockFilter` and `eth_getFilterChanges` to get the latest block information. Although it presented great performance improvement (similar to the one described above), Cloudflare's Public RPC node doesn't support either of the methods as of now: see [docs](https://developers.cloudflare.com/web3/ethereum-gateway/reference/supported-api-methods/).

### Webhook

The webhook calls don't currently implement any retry logic, which should be done in future work.

An even more reliable approach would be to implement an event-based system such as Kafka as to guarantee delivery and minimize latency.
