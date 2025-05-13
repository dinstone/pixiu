package repository

type PreferenceRepository interface {
	Load() ([]byte, error)
	Store(data []byte) error
}
