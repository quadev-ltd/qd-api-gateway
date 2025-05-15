package mock

//go:generate mockgen -destination=middleware_mock.go -package=mock github.com/quadev-ltd/qd-qpi-gateway/internal/middleware AutheticationMiddlewarer
