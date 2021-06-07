package service

import "github.com/lucsky/cuid"

func generateId() string {
	return cuid.New()
}
