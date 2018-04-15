package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/spf13/viper"

	"github.com/northwesternmutual/kanali/cmd/kanalictl/app/options"
	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	"github.com/northwesternmutual/kanali/pkg/rsa"
	"github.com/northwesternmutual/kanali/pkg/utils"
	"github.com/northwesternmutual/kanali/test/builder"
)

const (
	keyNameRegex = "^[a-z0-9]([-a-z0-9]*[a-z0-9])?$"
	keyDataRegex = "^[0-9a-zA-Z]+$"
)

func Generate(stdout, stderr io.Writer) error {
	outFile, err := getOutFile(viper.GetString(options.FlagKeyOutFile.GetLong()))
	if err != nil {
		utils.Write(stderr, err.Error())
		return err
	}

	publicKey, err := rsa.LoadPublicKey(viper.GetString(options.FlagRSAPublicKeyFile.GetLong()))
	if err != nil {
		utils.Write(stderr, err.Error())
		return err
	}

	unencrypted, err := getUnencryptedData(viper.GetString(options.FlagKeyData.GetLong()), viper.GetInt(options.FlagKeyLength.GetLong()))
	if err != nil {
		utils.Write(stderr, err.Error())
		return err
	}

	ciphertext, err := rsa.Encrypt(unencrypted, publicKey, rsa.Base64Encode(), rsa.WithEncryptionLabel(rsa.EncryptionLabel))
	if err != nil {
		utils.Write(stderr, err.Error())
		return err
	}

	if len(viper.GetString(options.FlagKeyName.GetLong())) < 1 {
		err := errors.New("must specify name")
		utils.Write(stderr, err.Error())
		return err
	}

	if !regexp.MustCompile(keyNameRegex).MatchString(viper.GetString(options.FlagKeyName.GetLong())) {
		err := errors.New("name must conform to the regex " + keyNameRegex)
		utils.Write(stderr, err.Error())
		return err
	}

	apikey := builder.NewApiKey(viper.GetString(options.FlagKeyName.GetLong())).WithRevision(v2.RevisionStatusActive, ciphertext).NewOrDie()

	if err := utils.Write(stdout, fmt.Sprintf("Here is your api key (you will only see this once): %s", string(unencrypted))); err != nil {
		return err
	}

	if len(outFile) < 1 {
		yamlData, err := yaml.Marshal(apikey)
		if err != nil {
			utils.Write(stderr, err.Error())
			return err
		}
		return utils.Write(stdout, fmt.Sprintf("\n---\n%s", string(yamlData)))
	}

	return writeApiKey(outFile, apikey)
}

func getUnencryptedData(existing string, length int) ([]byte, error) {
	if len(existing) > 0 {
		if !regexp.MustCompile(keyDataRegex).MatchString(existing) {
			return nil, fmt.Errorf("key data must conform to the pattern %s", keyDataRegex)
		}

		return []byte(existing), nil
	}
	return rsa.GenerateRandomBytes(length)
}

func getOutFile(f string) (string, error) {
	if len(f) < 1 {
		return f, nil
	}

	if len(strings.Split(f, ".")[0]) < 1 {
		return "", errors.New("out file must have a name")
	}

	switch filepath.Ext(f) {
	case "":
		return f + ".yaml", nil
	case ".yaml", ".yml", ".json":
	default:
		return "", errors.New("out file just be either json or yaml format")
	}

	return f, nil
}

func writeApiKey(outFileName string, keyCRD *v2.ApiKey) error {
	if len(outFileName) < 1 {
		return nil
	}

	_, err := os.Stat(outFileName)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if !os.IsNotExist(err) && shouldNotOverWrite(outFileName) {
		return nil
	}

	var fileData []byte

	if filepath.Ext(outFileName) != ".json" {
		yamlData, err := yaml.Marshal(keyCRD)
		if err != nil {
			return err
		}
		fileData = yamlData
	} else {
		jsonData, err := json.MarshalIndent(keyCRD, "", "   ")
		if err != nil {
			return err
		}
		fileData = jsonData
	}

	if err := ioutil.WriteFile(outFileName, fileData, 0644); err != nil {
		return err
	}

	fmt.Printf("Corresponding Kubernetes config written to %s\n", outFileName)
	return nil
}

func shouldNotOverWrite(outFileName string) bool {
	var input string

	for {
		fmt.Printf("%s exists - do you want to override it? (Y/n) ", outFileName)
		fmt.Scanln(&input)
		switch input {
		case "Y":
			return false
		case "n":
			return true
		}
	}
}
