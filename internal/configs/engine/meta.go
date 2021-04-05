package engine

type MetaData struct {
	Type    string `json:"type" enums:"OK,CREATED,ERROR"`
	Message string `json:"message"`
}

func NewMetaData(Type string, message string) MetaData {
	return MetaData{Type: Type, Message: message}
}
