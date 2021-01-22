module github.com/etclabscore/core-pool

go 1.15

require (
	github.com/btcsuite/btcd v0.21.0-beta // indirect
	github.com/edsrzf/mmap-go v1.0.0
	github.com/etclabscore/go-etchash v0.0.0-20201116172336-ee6bd2e29b29
	github.com/ethereum/go-ethereum v1.9.24
	github.com/garyburd/redigo v1.6.2 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/golang-lru v0.5.4
	github.com/yvasiyarov/go-metrics v0.0.0-20150112132944-c25f46c4b940 // indirect
	github.com/yvasiyarov/gorelic v0.0.7
	github.com/yvasiyarov/newrelic_platform_go v0.0.0-20160601141957-9c099fbc30e9 // indirect
	golang.org/x/crypto v0.0.0-20201116153603-4be66e5b6582
	golang.org/x/sys v0.0.0-20201109165425-215b40eba54c // indirect
	gopkg.in/bsm/ratelimit.v1 v1.0.0-20160220154919-db14e161995a // indirect
	gopkg.in/redis.v3 v3.6.4
)

// develop
//replace github.com/etclabscore/go-etchash => /home/xocel/go/src/github.com/etclabscore/go-etchash
//replace github.com/etclabscore/core-pool => /home/xocel/go/src/github.com/etclabscore/core-pool
