package domain

type MetaDataKeys int

const (
	UserIDKey MetaDataKeys = iota
	DocIDKey
)

func (m MetaDataKeys) String() string {
	switch m {
	case UserIDKey:
		return "user_id"
	case DocIDKey:
		return "doc_id"
	default:
		return ""
	}
}

type Chunks struct {
	ChunkID   string
	MetaData  map[MetaDataKeys]interface{}
	ChunkText string
}
