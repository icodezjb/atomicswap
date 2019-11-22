#### AtomicSwap example between two blockchain which are based on EVM
the Hashed Timelock Contract refers to [hashed-timelock-contract-ethereum](https://github.com/chatch/hashed-timelock-contract-ethereum)

##### solc, ganache-cli, go1.11
- `solc: Version: 0.5.10+commit.5a6ea5b1.Linux.g++`
- `ganache-cli: Ganache CLI v6.4.1 (ganache-core: 2.5.3)`
- `golang: go version go1.11.10 linux/amd64`

##### test
```bash
go test -v -run Test ./contract
```

##### install & lanuch  ganache client 

```bash
 sudo npm install -g ganache-cli
```

```bash
ganache-cli --account "0xa5a1aca01671e2660f1ee47abfd7065d5d38f99fa4a53495f02df939cd5b86f6,111111111111111111111" -p 7545
```
##### deploy contract
```bash
go run main.go
```
