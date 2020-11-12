### etchash

This is a standalone, pure go etchash module with ECIP-1099 implemented. It can be used by your pool in place of the go/cpp github.com/ethereum/ethash module. Simply copy this directory to your project, and import as you would any other go module.

See https://github.com/etclabscore/open-etc-pool/blob/master/proxy/miner.go#L9

### ECIP-1099 Activation

The activation block `ecip1099FBlock` is currently set for the **mordor** test network.
You can change this to classic mainnet here: https://github.com/etclabscore/open-etc-pool/blob/master/etchash/etchash.go#L61

Simply uncommet the mainnet line, and comment the mordor line. Rebuild.