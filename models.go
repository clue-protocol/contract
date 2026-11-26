package contract

type StaticUpdateMeta struct {
	Version      uint64
	SMTRoot      string
	ProviderCode uint32
	ProviderType string
	LeafAccounts [][]byte
}

type RouteUpdateMeta struct {
	Version      uint64
	SMTRoot      string
	ProviderCode uint32
	ProviderType string
	LeafAccounts [][]byte
	Linked       map[uint32]StaticUpdateMeta
}
