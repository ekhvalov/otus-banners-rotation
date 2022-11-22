//go:generate mockgen -destination=./mock/id_generator.gen.go -package mock . IDGenerator

package redis

import "github.com/google/uuid"

type IDGenerator interface {
	GenerateID() string
}

func NewUUIDGenerator() IDGenerator {
	return generator{}
}

type generator struct{}

func (g generator) GenerateID() string {
	return uuid.New().String()
}
