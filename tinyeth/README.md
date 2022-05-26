# Tiny ETH client

## Usage examples

### Reading data from the ledger

```go
rawHash, err := hashToBytes32(hash)
if err != nil {
    return nil, err
}
result, err := q.client.QueryBlockchain(ctx,
    q.Account,
    q.ContractAddress,
    q.ContractAbi,
    "getEvidence", // 0x4a7221a0
    rawHash,
)
if err != nil {
    return nil, err
}
// if no error, we always expect 1 result
// according to Solidity method definition
return result[0], nil
```

### Sending transactions

```go
rawHash, err := hashToBytes32(data.Hash)
if err != nil {
    return nil, err
}
return q.client.SendTransaction(ctx,
    q.Account,
    q.ContractAddress,
    q.ContractAbi,
    "addEvidence", // 0x55f0d402
    rawHash,
)
```