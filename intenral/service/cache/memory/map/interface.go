package memory_map

type Provider interface {
	Upsert(token, currency string, price float32)
	Get(token, currency string) (float32, bool)
}
