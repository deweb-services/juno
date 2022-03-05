# To test bridge preparing ledger with transactions

Create account for bridge:
```
./dewebd keys add bridge
```


Save accounts addresses to variables:
```
ALICE_ADDRESS='deweb1hhgzez20yv6vvf5a7w6sr0s5p4vwkn7cyxh8w5'
BOB_ADDRESS='deweb1xegmesuhx9jman8psxgynaupmchtqjf9us9f45'
BRIDGE_ADDRESS='deweb18dxv2ke6efedju6xwkmsfgyzls87zmxm96ufud'
```

In these requests created records for external wallets, removed one record. Created chain mappings, removed one mapping
```
./dewebd tx deweb save-wallet cosmosaddress private_cosmos_key cosmos --from alice --chain-id deweb-testnet-0 -y
./dewebd tx deweb save-wallet siaaddress private_cosmos_key sia --from alice --chain-id deweb-testnet-0 -y
./dewebd tx deweb delete-wallet siaaddress --from alice --chain-id deweb-testnet-0 -y

./dewebd tx deweb connect-chain cosmos cosmosaddress --from alice --chain-id deweb-testnet-0 -y
./dewebd tx deweb connect-chain sia siaaddress --from alice --chain-id deweb-testnet-0 -y
./dewebd tx deweb delete-chain-connect cosmos cosmosaddress --from alice --chain-id deweb-testnet-0 -y
```

Checking states. Expected that exist one wallet record and one mapping record
```
./dewebd q deweb filter-user-wallet-records $ALICE_ADDRESS

./dewebd q deweb filter-chain-mappings-records $ALICE_ADDRESS "" "" false 100 0
```

Building CW20 contract and store code in chain:
```
./dewebd tx wasm store cw20_base.wasm --from alice --chain-id deweb-testnet-0 --gas 2000000 --output json -b block 
```

From result set code:
```
CODE_ID=1
```

Deploy contract for wrapped SIA:
```
INIT="{\"name\":\"wrapped SIA\",\"symbol\":\"wSIA\",\"decimals\":18,\"initial_balances\":[{\"address\":\"$ALICE_ADDRESS\",\"amount\":\"20000\"},{\"address\":\"$BOB_ADDRESS\",\"amount\":\"30000\"},{\"address\":\"$BRIDGE_ADDRESS\",\"amount\":\"100000\"}],\"mint\":{\"minter\":\"$ALICE_ADDRESS\"}}"
./dewebd tx wasm instantiate $CODE_ID "$INIT" --from alice --label "wSIA contract" --chain-id deweb-testnet-0 --no-admin
```

Set contract address from deployment result. Preparing requests.
```
CONTRACT_ADDRESS='deweb14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s03q9ga'
BALANCE_QUERY_ALICE="{\"balance\": {\"address\": \"$ALICE_ADDRESS\"}}"
BALANCE_QUERY_BOB="{\"balance\": {\"address\": \"$BOB_ADDRESS\"}}"
BALANCE_QUERY_BRIDGE="{\"balance\": {\"address\": \"$BRIDGE_ADDRESS\"}}"
TRANSFER_REQUIEST="{\"transfer\": {\"recipient\": \"$BOB_ADDRESS\", \"amount\": \"100\"}}"
TRANSFER_REQUIEST_BRIDGE="{\"transfer\": {\"recipient\": \"$BRIDGE_ADDRESS\", \"amount\": \"200\"}}"
```

Executing transfer requests:
```
./dewebd tx wasm execute $CONTRACT_ADDRESS "$TRANSFER_REQUIEST" --from alice --chain-id deweb-testnet-0
./dewebd tx wasm execute $CONTRACT_ADDRESS "$TRANSFER_REQUIEST_BRIDGE" --from alice --chain-id deweb-testnet-0
```

Requesting for balances:
```
./dewebd query wasm contract-state smart $CONTRACT_ADDRESS "$BALANCE_QUERY_ALICE" --output json
./dewebd query wasm contract-state smart $CONTRACT_ADDRESS "$BALANCE_QUERY_BOB" --output json
./dewebd query wasm contract-state smart $CONTRACT_ADDRESS "$BALANCE_QUERY_BRIDGE" --output json
```

Expected balances:
- alice: 19700
- bob: 30100
- bridge: 100200

And created one bridge record from alice.