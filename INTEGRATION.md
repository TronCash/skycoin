# TRONCash Exchange Integration

A TRONCash node offers multiple interfaces:

* REST API on port 6930 (when running from source; if you are using the releases downloaded from the website, the port is randomized)
* JSON-RPC 2.0 API accessible on `/api/v1/webrpc` endpoint **[deprecated]**

A CLI tool is provided in `cmd/cli/cli.go`. This tool communicates over the JSON-RPC 2.0 API. In the future it will communicate over the REST API.

*Note*: Do not interface with the JSON-RPC 2.0 API directly, it is deprecated and will be removed in a future version.

The API interfaces do not support authentication or encryption so they should only be used over localhost.

If your application is written in Go, you can use these client libraries to interface with the node:

* [TRONCash REST API Client Godoc](https://godoc.org/github.com/troncash/troncash/src/api#Client)
* [TRONCash CLI Godoc](https://godoc.org/github.com/troncash/troncash/src/cli)

*Note*: The CLI interface will be deprecated and replaced with a better one in the future.

The wallet APIs in the REST API operate on wallets loaded from and saved to `~/.troncash/wallets`.
Use the CLI tool to perform seed generation and transaction signing outside of the TRONCash node.

The TRONCash node's wallet APIs can be enabled with `-enable-wallet-api`.

For a node used to support another application,
it is recommended to use the REST API for blockchain queries and disable the wallet APIs,
and to use the CLI tool for wallet operations (seed and address generation, transaction signing).

<!-- MarkdownTOC autolink="true" bracket="round" levels="1,2,3,4,5,6" -->

- [API Documentation](#api-documentation)
    - [Wallet REST API](#wallet-rest-api)
    - [TRONCash command line interface](#troncash-command-line-interface)
    - [TRONCash REST API Client Documentation](#troncash-rest-api-client-documentation)
    - [TRONCash Go Library Documentation](#troncash-go-library-documentation)
    - [libtroncash Documentation](#libtroncash-documentation)
- [Implementation guidelines](#implementation-guidelines)
    - [Scanning deposits](#scanning-deposits)
        - [Using the CLI](#using-the-cli)
        - [Using the REST API](#using-the-rest-api)
        - [Using troncash as a library in a Go application](#using-troncash-as-a-library-in-a-go-application)
    - [Sending coins](#sending-coins)
        - [General principles](#general-principles)
        - [Using the CLI](#using-the-cli-1)
        - [Using the REST API](#using-the-rest-api-1)
        - [Using troncash as a library in a Go application](#using-troncash-as-a-library-in-a-go-application-1)
        - [Coinhours](#coinhours)
            - [REST API](#rest-api)
            - [CLI](#cli)
    - [Verifying addresses](#verifying-addresses)
        - [Using the CLI](#using-the-cli-2)
        - [Using the REST API](#using-the-rest-api-2)
        - [Using troncash as a library in a Go application](#using-troncash-as-a-library-in-a-go-application-2)
        - [Using troncash as a library in other applications](#using-troncash-as-a-library-in-other-applications)
    - [Checking TRONCash node connections](#checking-troncash-node-connections)
        - [Using the CLI](#using-the-cli-3)
        - [Using the REST API](#using-the-rest-api-3)
        - [Using troncash as a library in a Go application](#using-troncash-as-a-library-in-a-go-application-3)
    - [Checking TRONCash node status](#checking-troncash-node-status)
        - [Using the CLI](#using-the-cli-4)
        - [Using the REST API](#using-the-rest-api-4)
        - [Using troncash as a library in a Go application](#using-troncash-as-a-library-in-a-go-application-4)

<!-- /MarkdownTOC -->


## API Documentation

### Wallet REST API

[Wallet REST API](src/api/README.md).

### TRONCash command line interface

[CLI command API](cmd/cli/README.md).

### TRONCash REST API Client Documentation

[TRONCash REST API Client](https://godoc.org/github.com/troncash/troncash/src/api#Client)

### TRONCash Go Library Documentation

[TRONCash Godoc](https://godoc.org/github.com/troncash/troncash)

### libtroncash Documentation

[libtroncash documentation](/lib/cgo/README.md)

## Implementation guidelines

### Scanning deposits

There are multiple approaches to scanning for deposits, depending on your implementation.

One option is to watch for incoming blocks and check them for deposits made to a list of known deposit addresses.
Another option is to check the unspent outputs for a list of known deposit addresses.

#### Using the CLI

To scan the blockchain, use `troncash-cli lastBlocks` or `troncash-cli blocks`. These will return block data as JSON
and new unspent outputs sent to an address can be detected.

To check address outputs, use `troncash-cli addressOutputs`. If you only want the balance, you can use `troncash-cli addressBalance`.

#### Using the REST API

To scan the blockchain, call `GET /api/v1/last_blocks?num=` or `GET /api/v1/blocks?start=&end=`. There will return block data as JSON
and new unspent outputs sent to an address can be detected.

To check address outputs, call `GET /api/v1/outputs?addrs=`. If you only want the balance, you can call `GET /api/v1/balance?addrs=`.

* [`GET /api/v1/last_blocks` docs](src/api/README.md#get-last-n-blocks)
* [`GET /api/v1/blocks` docs](src/api/README.md#get-blocks-in-specific-range)
* [`GET /api/v1/outputs` docs](src/api/README.md#get-unspent-output-set-of-address-or-hash)
* [`GET /api/v1/balance` docs](src/api/README.md#get-balance-of-addresses)

#### Using troncash as a library in a Go application

We recommend using the [TRONCash REST API Client](https://godoc.org/github.com/troncash/troncash/src/api#Client).

### Sending coins

#### General principles

After each spend, wait for the transaction to confirm before trying to spend again.

For higher throughput, combine multiple spends into one transaction.

TRONCash uses "coin hours" to ratelimit transactions.
The total number of coinhours in a transaction's outputs must be 50% or less than the number of coinhours in a transaction's inputs,
or else the transaction is invalid and will not be accepted. A transaction must have at least 1 input with at least 1 coin hour.
Sending too many transactions in quick succession will use up all available coinhours.
Coinhours are earned at a rate of 1 coinhour per coin per hour, calculated per second.
This means that 3600 coins will earn 1 coinhour per second.
However, coinhours are only updated when a new block is published to the blockchain.
New blocks are published every 10 seconds, but only if there are pending transactions in the network.

To avoid running out of coinhours in situations where the application may frequently send,
the sender should batch sends into a single transaction and send them on a
30 second to 1 minute interval.

There are other strategies to minimize the likelihood of running out of coinhours, such
as splitting up balances into many unspent outputs and having a large balance which generates
coinhours quickly.

#### Using the CLI

When sending coins from the CLI tool, a wallet file local to the caller is used.
The CLI tool allows you to specify the wallet file on disk to use for operations.

See [CLI command API](cmd/cli/README.md) for documentation of the CLI interface.

To perform a send, the preferred method follows these steps in a loop:

* `troncash-cli createRawTransaction -m '[{"addr:"$addr1,"coins:"$coins1"}, ...]` - `-m` flag is send-to-many
* `troncash-cli broadcastTransaction` - returns `txid`
* `troncash-cli transaction $txid` - repeat this command until `"status"` is `"confirmed"`

That is, create a raw transaction, broadcast it, and wait for it to confirm.

#### Using the REST API

Create a transaction with [POST /wallet/transaction](https://github.com/troncash/troncash/blob/develop/src/api/README.md#create-transaction),
then inject it to the network with [POST /injectTransaction](https://github.com/troncash/troncash/blob/develop/src/api/README.md#inject-raw-transaction).

When using `POST /wallet/transaction`, a wallet file local to the troncash node is used.
The wallet file is specified by wallet ID, and all wallet files are in the
configured data directory (which is `$HOME/.troncash/wallets` by default).

#### Using troncash as a library in a Go application

If your application is written in Go, you can interface with the CLI library
directly, see [TRONCash CLI Godoc](https://godoc.org/github.com/troncash/troncash/src/cli).

A REST API client is also available: [TRONCash REST API Client Godoc](https://godoc.org/github.com/troncash/troncash/src/api#Client).

#### Coinhours

Transaction fees in troncash is paid in coinhours and is currently set to `50%`,
every transaction created burns `50%` of the total coinhours in all the input
unspents.

You need a minimum of `1` of coinhour to create a transaction.

Coinhours are generated at a rate of `1 coinsecond` per `second`
which are then converted to `coinhours`, `1` coinhour = `3600` coinseconds.

> Note: Coinhours don't have decimals and only show up in whole numbers.

##### REST API

When using the REST API, the coin hours sent to the destination and change can be controlled.
The 50% burn fee is still required.

See the [POST /wallet/transaction](https://github.com/troncash/troncash/blob/develop/src/api/README.md#create-transaction)
documentation for more information on how to control the coin hours.

We recommend sending at least 1 coin hour to each destination, otherwise the receiver will have to
wait for another coin hour to accumulate before they can make another transaction.

##### CLI

When using the CLI the amount of coinhours sent to the receiver is capped to
the number of coins they receive with a minimum of `1` coinhour for transactions
with `<1` troncash being sent.

The coinhours left after burning `50%` and sending to receivers are sent to the change address.

For eg. If an address has `10` troncashs and `50` coinhours and only `1` unspent.
If we send `5` troncashs to another address then that address will receive
`5` troncashs and `5` coinhours, `26` coinhours will be burned.
The sending address will be left with `5` troncashs and `19` coinhours which
will then be sent to the change address.


### Verifying addresses

#### Using the CLI

```sh
troncash-cli verifyAddress $addr
```

#### Using the REST API

Not directly supported, but API calls that have an address argument will return `400 Bad Request` if they receive an invalid address.

#### Using troncash as a library in a Go application

https://godoc.org/github.com/troncash/troncash/src/cipher#DecodeBase58Address

```go
if _, err := cipher.DecodeBase58Address(address); err != nil {
    fmt.Println("Invalid address:", err)
    return
}
```

#### Using troncash as a library in other applications

Address validation is available through a C wrapper, `libtroncash`.

See the [libtroncash documentation](/lib/cgo/README.md) for usage instructions.

### Checking TRONCash node connections

#### Using the CLI

Not implemented

#### Using the REST API

* `GET /api/v1/network/connections`

#### Using troncash as a library in a Go application

Use the [TRONCash REST API Client](https://godoc.org/github.com/troncash/troncash/src/api#Client)

### Checking TRONCash node status

#### Using the CLI

```sh
troncash-cli status
```

#### Using the REST API

A method similar to `troncash-cli status` is not implemented, but these endpoints can be used:

* `GET /api/v1/health`
* `GET /api/v1/version`
* `GET /api/v1/blockchain/metadata`
* `GET /api/v1/blockchain/progress`

#### Using troncash as a library in a Go application

Use the [TRONCash CLI package](https://godoc.org/github.com/troncash/troncash/src/cli)
