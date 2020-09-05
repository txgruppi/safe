package database

type keys struct {
	mimeType []byte
	size     []byte
	data     []byte
}

func keysForID(id string) *keys {
	return &keys{
		mimeType: []byte(id + ":mimeType"),
		size:     []byte(id + ":size"),
		data:     []byte(id + ":data"),
	}
}
