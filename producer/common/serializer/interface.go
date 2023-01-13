package serializer

type Serializer interface {
	Serialize(object interface{}) ([]byte, error)
	Deserialize(data []byte, dest interface{}) error
}
