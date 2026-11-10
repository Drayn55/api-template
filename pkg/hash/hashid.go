package hash

import (
	"os"

	"github.com/speps/go-hashids"
)

func getSalt() string {
	salt := os.Getenv("HASHID_SALT")
	if salt == "" {
		salt = "hansway-secret-key"
	}
	return salt
}

func EncodeID(id int64) (string, error) {
	hd := hashids.NewData()
	hd.Salt = getSalt()
	hd.MinLength = 8

	h, err := hashids.NewWithData(hd)
	if err != nil {
		return "", err
	}

	return h.EncodeInt64([]int64{id})
}

func DecodeID(hash string) (int64, error) {
	hd := hashids.NewData()
	hd.Salt = getSalt()
	hd.MinLength = 8

	h, err := hashids.NewWithData(hd)
	if err != nil {
		return 0, err
	}

	ids, err := h.DecodeInt64WithError(hash)
	if err != nil {
		return 0, err
	}

	return ids[0], nil
}

