package dbfs

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"golang.org/x/crypto/scrypt"
	"gorm.io/gorm"
	"io"
)

type PasswordProtect struct {
	FileId         string `gorm:"primaryKey"`
	FileKey        string
	FileIv         string
	PasswordActive bool
	PasswordSalt   string
	PasswordNonce  string
	PasswordHint   string
}

func EncryptWithKeyIV(plaintext, key, iv string) (string, error) {

	keyBytes, err := hex.DecodeString(key)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ivByte, err := hex.DecodeString(iv)
	plaintextByte, err := hex.DecodeString(plaintext)

	ciphertext := aesgcm.Seal(nil, ivByte, plaintextByte, nil)

	return hex.EncodeToString(ciphertext), nil
}

func DecryptWithKeyIV(ciphertext, key, iv string) (string, error) {

	keyBytes, err := hex.DecodeString(key)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)

	if err != nil {
		return "", err
	}

	ivByte, err := hex.DecodeString(iv)
	ciphertextByte, err := hex.DecodeString(ciphertext)

	plaintext, err := aesgcm.Open(nil, ivByte, ciphertextByte, nil)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(plaintext), nil

}

func GenerateKeyIV() (string, string, error) {
	key := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, key)
	if err != nil {
		return "", "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", "", err
	}

	iv := make([]byte, aesgcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return "", "", err
	}

	return hex.EncodeToString(key), hex.EncodeToString(iv), nil
}

func (p *PasswordProtect) encryptWithPassword(tx *gorm.DB, oldPassword, newPassword, hint string) error {

	var ogKey, ogIv string
	var err error

	ogKey = p.FileKey
	ogIv = p.FileIv

	if p.PasswordActive {
		// get original key and iv
		ogKey, ogIv, err = p.DecryptWithPassword(oldPassword)
		if err != nil {
			return err
		}
	}

	// Generate a secure salt
	salt := make([]byte, 32)
	_, err = rand.Read(salt)
	if err != nil {
		return err
	}

	key, err := scrypt.Key([]byte(newPassword), salt, 16384, 8, 1, 32)
	if err != nil {
		return err
	}

	// ENCRYPT
	fileKeyBytes, err := hex.DecodeString(ogKey)
	if err != nil {
		return err
	}

	fileIvBytes, err := hex.DecodeString(ogIv)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}
	p.PasswordNonce = hex.EncodeToString(nonce)

	newFileKeyBytes := aesgcm.Seal(nil, nonce, fileKeyBytes, nil)
	newFileIvBytes := aesgcm.Seal(nil, nonce, fileIvBytes, nil)

	p.FileKey = hex.EncodeToString(newFileKeyBytes)
	p.FileIv = hex.EncodeToString(newFileIvBytes)
	p.PasswordActive = true
	p.PasswordNonce = hex.EncodeToString(nonce)
	p.PasswordSalt = hex.EncodeToString(salt)
	p.PasswordHint = hint

	return tx.Save(p).Error

}

func (p *PasswordProtect) DecryptWithPassword(password string) (string, string, error) {

	if password != "" && p.PasswordActive {
		// Getting everything in bytes first
		passwordSaltBytes, err := hex.DecodeString(p.PasswordSalt)
		passwordNonceBytes, err := hex.DecodeString(p.PasswordNonce)
		fileKeyBytes, err := hex.DecodeString(p.FileKey)
		fileIvBytes, err := hex.DecodeString(p.FileIv)

		// get key from password
		key, err := scrypt.Key([]byte(password), passwordSaltBytes, 16384, 8, 1, 32)
		if err != nil {
			return "", "", err
		}

		// try to decrypt file key with PasswordProtect

		block, err := aes.NewCipher(key)
		if err != nil {
			return "", "", err
		}

		aesgcm, err := cipher.NewGCM(block)
		if err != nil {
			return "", "", err
		}

		actualKey, err := aesgcm.Open(nil, passwordNonceBytes, fileKeyBytes, nil)
		if err != nil {
			return "", "", ErrPasswordIncorrect
		}
		actualIv, err := aesgcm.Open(nil, passwordNonceBytes, fileIvBytes, nil)
		if err != nil {
			return "", "", ErrPasswordIncorrect
		}

		return hex.EncodeToString(actualKey), hex.EncodeToString(actualIv), nil

	} else if !p.PasswordActive {
		return p.FileKey, p.FileIv, nil
	} else {
		return "", "", ErrPasswordRequired
	}
}

func (p *PasswordProtect) removePassword(tx *gorm.DB, password string) error {
	if !p.PasswordActive {
		return ErrNoPassword
	}

	ogKey, ogIv, err := p.DecryptWithPassword(password)
	if err != nil {
		return err
	}
	p.FileKey = ogKey
	p.FileIv = ogIv
	p.PasswordActive = false

	return tx.Save(p).Error

}
