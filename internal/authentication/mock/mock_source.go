package mock

//go:generate mockgen -destination=authentication_mock.go -package=mock github.com/quadev-ltd/qd-common/pb/gen/go/pb_authentication AuthenticationServiceClient
