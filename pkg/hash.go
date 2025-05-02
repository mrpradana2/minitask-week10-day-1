package pkg

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type HashConfig struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
	SaltLen uint32
}

func InitHashConfig() *HashConfig {
	return &HashConfig{}
}

func (h *HashConfig) UseConfig(time, memory, keylen, saltLen uint32, threads uint8) {
	h.Threads = threads
	h.Time = time
	h.Memory = memory
	h.KeyLen = keylen
	h.SaltLen = saltLen
}

func (h *HashConfig) UseDefaultConfig() {
	h.Threads = 2
	h.Time = 3
	h.Memory = 64 * 1024
	h.KeyLen = 32
	h.SaltLen = 16
}

func (h *HashConfig) genSalt() ([]byte, error) {
	salt := make([]byte, h.SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	return salt, nil
}

func (h *HashConfig) GenHashedPassword(password string) (string, error) {
	// hash = password + salt + config
	salt, err := h.genSalt()
	if err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(password), salt, h.Time, h.Memory, h.Threads, h.KeyLen)
	// dalam penulisan hash ada format
	// $jenisKey$varsi
	version := argon2.Version
	base64Salt := base64.RawStdEncoding.EncodeToString(salt)
	base64Hash := base64.RawStdEncoding.EncodeToString(hash)
	hashedPwd := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", version, h.Memory, h.Time, h.Threads, base64Salt, base64Hash)
	return hashedPwd, nil
}

// compare hash and password
func (h *HashConfig) CompareHashAndPassword(hashedPass string, password string) (bool, error) {
	// komparasi antara pass yang datang dari db dan body
	// yang datang dari body adalah hashedpass, bagaimana kira ekstrak confignya
	salt, hash, err := h.decodeHash(hashedPass)
	if err != nil {
		return false, err
	}

	// hati hati dengan timing attack
	// if slices.Compare(hash, newHash) != 0 (
	newHash := argon2.IDKey([]byte(password), salt, h.Time, h.Memory, h.Threads, h.KeyLen)

	if subtle.ConstantTimeCompare(hash, newHash) == 0 {
		return false, err
	}

	return true, nil
}

// decode hash
func (h *HashConfig) decodeHash(hashedPass string) ([]byte, []byte, error) {
	values := strings.Split(hashedPass, "$")

	// mengecek panjang hashedPassword panjangnya harus 6
	if len(values) != 6 {
		// log.Println("[DEBUGGGG]: length", len(values))
		return nil, nil, errors.New("invalid format")
	}

	// mengecek type hash harus menggunakan argon2
	if values[1] != "argon2id" {
		return nil, nil, errors.New("invalid hash type")
	}

	// mengecek versi hash, harus sesuai versinya
	var version int
	if _, err := fmt.Sscanf(values[2], "v=%d", &version); err != nil {
		// log.Println("[DEBUGGGG]", err)
		return nil, nil, errors.New("invalid format")
	}

	if version != argon2.Version {
		return nil, nil, errors.New("invalid argon2id version")
	}

	// mengecek config (memory, time, dan threads nya sudah selesai)
	if _, err := fmt.Sscanf(values[3], "m=%d,t=%d,p=%d", &h.Memory, &h.Time, &h.Threads); err != nil {
		// log.Println("[DEBUGGGG]", err)
		return nil, nil, errors.New("invalid fromat")
	}

	// mengecek panjang dari salt di config dengan salt hasil konversi dari base 64
	salt, err := base64.RawStdEncoding.DecodeString(values[4])
	if err != nil {
		// log.Println("[DEBUGGGG]", err)
		return nil, nil, err
	}
	// sini
	h.SaltLen = uint32(len(salt))

	hash, err := base64.RawStdEncoding.DecodeString(values[5])
	if err != nil {
		// log.Println("[DEBUGGGG]", err)
		return nil, nil, err
	}

	h.KeyLen = uint32(len(hash))
	return salt, hash, nil
}

