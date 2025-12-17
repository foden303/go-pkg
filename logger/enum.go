package logger

type ENV string

const (
	EnvDev  ENV = "dev"
	EnvProd ENV = "prod"
	EnvTest ENV = "test"
)

func (e ENV) String() string {
	return string(e)
}
