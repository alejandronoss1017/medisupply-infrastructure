# Novelties (logistics-distributions/novelties)

A small Go service and library that demonstrates how to interact with an Ethereum smart contract from within the MediSupply infrastructure. It includes:

- An HTTP server with a health endpoint (/ping)
- A lightweight Ethereum client wrapper (SmartContractClient) to:
  - Call read-only (view/pure) contract methods
  - Send state-changing transactions (writes)
  - Wait for transaction mining and check success
  - Query ETH balances

This module is meant for experimentation and as a starting point to integrate on-chain actions within the logistics/distributions domain.

---

## Prerequisites

- Go 1.24 (toolchain go1.24.5), as declared in `go.mod`
- Access to an Ethereum-compatible node (HTTP or WS RPC)
- A deployed contract and its ABI JSON
- A funded account’s private key for sending transactions

> Important: The private key must be a hex string WITHOUT the `0x` prefix (format expected by go-ethereum's `crypto.HexToECDSA`).

---

## Configuration

The web example (`cmd/web/main.go`) reads the following environment variables:

- `RCP_URL` — RPC endpoint (HTTP/HTTPS or WS), e.g. `http://localhost:8545`
- `SMART_CONTRACT_ADDRESS` — Contract address (0x-prefixed), e.g. `0xabc...`
- `PRIVATE_KEY` — Hex-encoded private key WITHOUT `0x`
- `ABI_PATH` — Filesystem path to the ABI JSON file

Example (PowerShell):

```powershell
$env:RCP_URL = "http://localhost:8545"
$env:SMART_CONTRACT_ADDRESS = "0xYourContractAddressHere"
$env:PRIVATE_KEY = "0123abcd...deadbeef"
$env:ABI_PATH = "C:\\Users\\you\\project\\logistics-distributions\\novelties\\abi.json"
```

> ABI format: The client parses the file with `abi.JSON(...)`, which expects the ABI JSON itself (array/object). If your build tool produces an artifact with many fields, extract and save only the `abi` section into a standalone file for this demo.

---

## Running the example web server

From the repository root:

```powershell
cd logistics-distributions\novelties
go run ./cmd/web
```

If the environment is set correctly, the program will:

1. Create an Ethereum client with your RPC, contract, key, and ABI
2. Call the contract method `retrieveNumber` (as written in `cmd/web/main.go`)
3. Log the returned value to the console
4. Start an HTTP server (default Gin settings) and expose:
   - `GET /ping` → `{ "message": "pong" }`

Open: http://localhost:8080/ping

---

## Using the Ethereum client in code

The core wrapper lives in `internal/adapter/ethereum/client.go`.

### Create a client

```
c, err := ethereum.NewSmartContractClient(rcpURL, contractAddr, privateKeyHex, abiPath)
if err != nil { /* handle */ }
```

- `rcpURL`: RPC endpoint URL
- `contractAddr`: 0x-prefixed contract address
- `privateKeyHex`: hex private key without `0x`
- `abiPath`: path to ABI JSON (the ABI itself, not the full artifact)

### Read-only call (view/pure)

The current signature expects a pointer to an `any` variable as the destination. For example, calling `balanceOf(address)` or a simple getter:

```
ctx := context.Background()
var out any
if err := c.InvoqueContract(ctx, "retrieveNumber", &out /* args... */); err != nil {
    // handle error
}

// If the method returns a single value, `out` will hold that value.
// For multiple return values, `out` will typically be a []interface{}.
log.Printf("result: %#v\n", out)
```

Type assert as needed:

```
if n, ok := out.(*big.Int); ok {
    fmt.Println("number:", n)
} else if items, ok := out.([]interface{}); ok {
    // handle multiple return values in order
}
```

> Note: For this demo, the decode target is `*any`, which is flexible but requires type assertions. If you want stronger typing in your app, you can wrap this helper and map results into your own types.

### Write transaction (state-changing)

Send a transaction to a method like `setNumber(uint256)` and wait until it’s mined:

```
ctx := context.Background()

// If your method doesn’t accept ETH, pass nil as the value
// and the client will default it to 0 wei.
tx, err := c.InvoqueContractWrite(ctx, "setNumber", nil, big.NewInt(42))
if err != nil {
    // handle error
}
fmt.Println("broadcasted tx:", tx.Hash().Hex())

// Wait for mining and verify success
rcpt, err := c.WaitTransaction(ctx, tx)
if err != nil {
    // handle error (includes failed receipt)
}
fmt.Printf("mined in block %v, status %d\n", rcpt.BlockNumber, rcpt.Status)
```

The client:
- Packs inputs with the loaded ABI
- Estimates gas and adds a small buffer
- Uses EIP-1559 fees if available (falls back to legacy gas price otherwise)
- Signs with your private key and broadcasts the transaction

### Check balances

```
bal, err := c.GetBalance(ctx, someAddress)
if err != nil { /* handle */ }
fmt.Println("balance (wei):", bal)
```

---

## Troubleshooting

- Private key format: must be hex without `0x`. If you see "error converting private key", remove the prefix or verify the value.
- ABI parse errors: ensure your `ABI_PATH` file contains the ABI JSON itself (array/object), not an entire artifact. Many build tools output large JSONs — copy only the `abi` section into a new file.
- Method not found: verify the method name exists in the ABI and that argument types and order match.
- Gas estimation failures: they usually mean invalid arguments, revert conditions in the contract, or that the call would fail. Double-check inputs and contract state.

---

## Project structure (subset)

```
logistics-distributions/novelties
├── abi.json                     # Example ABI file (you can replace with your own)
├── cmd/web/main.go              # Minimal web server + example read call
└── internal/adapter/ethereum/
    └── client.go                # SmartContractClient implementation
```

---

## Security notes

- Never commit real private keys. Use development accounts or environment variables managed by secure tooling.
- For production, consider using an external signer (e.g., HashiCorp Vault, AWS KMS) and stricter nonce/fee policies.

---

## License

Check the [LICENSE](./LICENSE) file.