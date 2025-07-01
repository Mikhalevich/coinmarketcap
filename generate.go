package coinmarketcap

//go:generate go tool mockgen -source=./request_executor.go -destination=./request_executor_mock.go -package=coinmarketcap
//go:generate go tool mockgen -source=./api/cryptocurrency/cryptocurrency.go -destination=./api/cryptocurrency/cryptocurrency_mock.go -package=cryptocurrency
