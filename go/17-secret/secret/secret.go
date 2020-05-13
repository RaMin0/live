package secret

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ramin0/live/go/secret/encrypt"
)

var (
	ErrKeyNotFound = errors.New("key not found")
)

type Vault struct {
	encodingKey string
	keysPath    string
}

func FileVault(encodingKey, keysPath string) *Vault {
	return &Vault{encodingKey, keysPath}
}

func (v *Vault) Set(keyName, keyValue string) error {
	// 1. read keys file
	encryptedData, err := v.readFile()
	if err != nil {
		return err
	}
	kv := map[string]string{}
	// if the file is not empty, i.e. keys already exist
	if encryptedData != "" {
		// 2. decrypt
		data, err := v.decryptFile(encryptedData)
		if err != nil {
			return err
		}
		kv = v.decodeFile(data)
	}
	// 3. set key
	if keyValue == "" { // delete if empty value
		delete(kv, keyName)
	} else {
		kv[keyName] = keyValue
	}
	data := v.encodeFile(kv)
	// 4. encrypt
	encryptedData, err = v.encryptFile(data)
	if err != nil {
		return err
	}
	// 5. write keys file
	return v.writeFile(encryptedData)
}

func (v *Vault) Get(keyName string) (string, error) {
	// 1. read keys file
	encryptedData, err := v.readFile()
	if err != nil {
		return "", err
	}
	if encryptedData == "" {
		return "", ErrKeyNotFound
	}
	// 2. decrypt
	data, err := v.decryptFile(encryptedData)
	if err != nil {
		return "", err
	}
	// 3. find key
	kv := v.decodeFile(data)
	keyValue, ok := kv[keyName]
	if !ok {
		return "", ErrKeyNotFound
	}
	// 4. return value
	return keyValue, nil
}

func (v *Vault) List() ([]string, error) {
	// 1. read keys file
	encryptedData, err := v.readFile()
	if err != nil {
		return nil, err
	}
	if encryptedData == "" {
		return nil, nil
	}
	// 2. decrypt
	data, err := v.decryptFile(encryptedData)
	if err != nil {
		return nil, err
	}
	// 3. populate keys
	kv := v.decodeFile(data)
	var keyNames []string
	for keyName := range kv {
		keyNames = append(keyNames, keyName)
	}
	return keyNames, nil
}

func (v *Vault) Delete(keyName string) error {
	if _, err := v.Get(keyName); err != nil {
		return err
	}
	return v.Set(keyName, "")
}

func (v *Vault) readFile() (string, error) {
	f, err := os.OpenFile(v.keysPath, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return "", err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (v *Vault) writeFile(data string) error {
	return ioutil.WriteFile(v.keysPath, []byte(data), 0600)
}

func (v *Vault) decryptFile(data string) (string, error) {
	return encrypt.Decrypt(v.encodingKey, data)
}

func (v *Vault) encryptFile(data string) (string, error) {
	return encrypt.Encrypt(v.encodingKey, data)
}

func (v *Vault) decodeFile(data string) map[string]string {
	kv := map[string]string{}
	sc := bufio.NewScanner(strings.NewReader(data))
	for sc.Scan() {
		line := strings.Split(sc.Text(), "=")
		if len(line) < 2 {
			continue
		}
		keyName, keyValue := line[0], line[1]
		kv[keyName] = keyValue
	}
	return kv
}

func (v *Vault) encodeFile(kv map[string]string) string {
	var data []string
	for keyName, keyValue := range kv {
		data = append(data, fmt.Sprintf("%s=%s", keyName, keyValue))
	}
	return strings.Join(data, "\n")
}
