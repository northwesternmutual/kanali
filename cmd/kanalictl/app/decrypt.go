// Copyright (c) 2018 Northwestern Mutual.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package app

import (
	"bufio"
	"bytes"
	"crypto/rsa"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ghodss/yaml"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	yamlReader "k8s.io/apimachinery/pkg/util/yaml"

	"github.com/northwesternmutual/kanali/cmd/kanalictl/app/options"
	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	rsautils "github.com/northwesternmutual/kanali/pkg/rsa"
	"github.com/northwesternmutual/kanali/pkg/utils"
)

type result struct {
	name string
	data []revision
}

type revision struct {
	status v2.RevisionStatus
	data   string
}

func Decrypt(stdout, stderr io.Writer) error {
	decryptionKey, err := rsautils.LoadDecryptionKey(viper.GetString(options.FlagRSAPrivateKeyFile.GetLong()))
	if err != nil {
		utils.Write(stderr, err.Error())
		return err
	}

	fileList, err := discoverFiles(viper.GetString(options.FlagKeyInFile.GetLong()))
	if err != nil {
		utils.Write(stderr, err.Error())
		return err
	}

	data, err := renderResults(decryptFiles(fileList, decryptionKey))
	if err != nil {
		utils.Write(stderr, err.Error())
		return err
	}

	return utils.Write(stderr, data)
}

func renderResults(results []result) (string, error) {
	if len(results) < 1 {
		return "no keys found", nil
	}

	buffer := bytes.NewBuffer(nil)
	table := tablewriter.NewWriter(buffer)
	table.SetAutoWrapText(false)
	table.SetHeader([]string{"Name", "Revisions"})

	for _, r := range results {
		tmp := []string{}
		for _, rev := range r.data {
			tmp = append(tmp, fmt.Sprintf("%-8s - %s", rev.status, rev.data))
		}
		table.Append([]string{r.name, strings.Join(tmp, "\n")})
	}

	table.Render()
	return buffer.String(), nil
}

func decryptFiles(fileList []string, key *rsa.PrivateKey) []result {
	allResults := []result{}
	resultsChan := make(chan []result)

	wg := sync.WaitGroup{}
	wg.Add(len(fileList))

	for _, file := range fileList {
		go func(file string) {
			resultsChan <- decryptFile(file, key)
		}(file)
	}
	go func() {
		for r := range resultsChan {
			allResults = append(allResults, r...)
			wg.Done()
		}
	}()

	wg.Wait()
	return allResults
}

func decryptFile(file string, key *rsa.PrivateKey) []result {
	fileResults := []result{}
	fileResultsChan := make(chan result)
	wg := sync.WaitGroup{}

	fileData, err := ioutil.ReadFile(file)
	if err != nil {
		return fileResults
	}

	reader := yamlReader.NewYAMLReader(bufio.NewReader(bytes.NewReader(fileData)))

	for {
		doc, err := reader.Read()
		if err != nil && err != io.EOF {
			return fileResults
		} else if err != nil && err == io.EOF {
			break
		}

		wg.Add(1)
		go func(doc []byte) {
			currResult, err := decryptKey(doc, key)
			if err != nil {
				wg.Done()
			} else {
				fileResultsChan <- *currResult
			}
		}(doc)
		go func() {
			for r := range fileResultsChan {
				fileResults = append(fileResults, r)
				wg.Done()
			}
		}()
	}

	wg.Wait()
	return fileResults
}

func decryptKey(data []byte, key *rsa.PrivateKey) (*result, error) {
	var meta metav1.TypeMeta

	err := yaml.Unmarshal(data, &meta)
	if err != nil {
		return nil, errors.New("yaml document did not contain Kubernetes resource")
	}

	if meta.Kind != "ApiKey" {
		return nil, errors.New("yaml document was not an ApiKey")
	}

	var apikey v2.ApiKey
	if err := yaml.Unmarshal(data, &apikey); err != nil {
		return nil, err
	}

	arr := make([]revision, len(apikey.Spec.Revisions))

	for i := 0; i < len(apikey.Spec.Revisions); i++ {
		arr[i].status = apikey.Spec.Revisions[i].Status
		unecryptedData, err := rsautils.Decrypt([]byte(apikey.Spec.Revisions[i].Data), key, rsautils.Base64Decode(), rsautils.WithEncryptionLabel(rsautils.EncryptionLabel))
		if err != nil {
			arr[i].data = err.Error()
		} else {
			arr[i].data = string(unecryptedData)
		}
	}

	return &result{
		name: apikey.GetName(),
		data: arr,
	}, nil
}

func discoverFiles(inFilePath string) ([]string, error) {
	fileList := []string{}

	err := filepath.Walk(inFilePath, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		switch filepath.Ext(path) {
		case ".yaml", ".yml", ".json":
			fileList = append(fileList, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return fileList, nil
}
