module github.com/router-for-me/CLIProxyAPI/v7/examples/plugin/gateway/go

go 1.26.0

require (
	github.com/redis/go-redis/v9 v9.19.0
	github.com/router-for-me/CLIProxyAPI/v7 v7.0.0
	github.com/tidwall/gjson v1.18.0
	github.com/tidwall/sjson v1.2.5
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
)

replace github.com/router-for-me/CLIProxyAPI/v7 => ../../CLIProxyAPI
