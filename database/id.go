package database

import (
	gonanoid "github.com/matoous/go-nanoid"
)

func generateID() (string, error) {
	return gonanoid.Generate("123456789abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ", 13)
}
