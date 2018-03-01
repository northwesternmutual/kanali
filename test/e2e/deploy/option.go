package deploy

type TLSType int
type Option func(config)

const (
	TLSTypeNone TLSType = iota
	TLSTypePresent
	TLSTypeMutual
)

func WithServer(serverType TLSType) Option {
	return Option(func(cfg config) {
		cfg.SetServerType(serverType)
	})
}
