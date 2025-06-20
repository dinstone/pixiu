package system

type PreferenceRepository interface {
	Load() ([]byte, error)
	Store(data []byte) error
}
