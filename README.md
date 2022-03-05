# DWS bridge component

This product based on juno (https://github.com/forbole/juno) implementation of blockchain data aggregator and exporter.
It is part of DWS bridge to connected non-Cosmos networks. 

## Config

Create config file in bridge home directory. Create empty sceleton with **init** command

Config consists of parts:
- chain - chain specific settings. Important parameter _modules_ - list of modules to run
- bridge - settings for tokens transfers
- node - setting for blockchain node
- parsing - settings for parser
- database - settings for DB for persistence
- logging - logging settings

### Chain config

Value for DWS bridge:
```
chain:
    bech32_prefix: deweb
    modules: [bridge_transactions]
```

### Bridge config
Settings:
- address - wallet address of bridge, watching for transfers to this address
- networks - list of bridged networks and wrapped tokens addresses
- consensus_host - address for notifications of requested transfers

Example value:
```
bridge:
    address: deweb1rl3wl5v7cln6m3hekp39lfe6244t420mnc540m
    networks:
        sia:
            token: deweb14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s03q9ga
    consensus_host: http://localhost:8083
```

Networks config part used for mapping between wrapped token and target network. In example added mapping between token
and network SIA. If we find transaction to _address_, checking transferred token. By token determined network name
to which transfer requested. Then performing request ot node (grpc) for stored mapping between transaction creator and
address in target network. When mapping found creating JSON and sending POST request to _consensus_host_. Amount of
tokens determined from ERC20 wrapped token transfer. Request JSON
example:
```
{
    "chain":"sia",
    "address":"siaaddress",
    "amount":"200"
}
```