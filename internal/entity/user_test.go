package entity_test

import (
	"testing"

	"github.com/ilovepitsa/jwt_tokens/internal/entity"
)

func TestUser(t *testing.T) {
	user := entity.User{
		Id: 1,
	}
	t.Run("TestId", func(t *testing.T) {
		want := uint32(1)
		got := user.GetID()

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

}
