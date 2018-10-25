// Copyright © 2018 Jesse Swidler <jswidler@gmail.com>
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

package content

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/jswidler/lockgit/pkg/log"
	"github.com/nu7hatch/gouuid"
	"github.com/pkg/errors"
	"github.com/sethvargo/go-diceware/diceware"
)

type LgConfig struct {
	Ver  int
	Id   string
	Name string
}

func NewLgConfig() LgConfig {
	id, err := uuid.NewV4()
	log.FatalPanic(err)
	tokens, err := diceware.Generate(2)
	log.FatalPanic(err)
	name := strings.Join(tokens, "-")
	return LgConfig{
		Ver:  1,
		Id:   id.String(),
		Name: name,
	}
}

func (config LgConfig) Write(path string) {
	filedata, err := json.Marshal(config)
	log.FatalPanic(err)
	err = ioutil.WriteFile(path, filedata, 0644)
	log.FatalPanic(err)
}

func ReadConfig(ctx Context) (LgConfig, error) {
	config := LgConfig{}
	filedata, err := ioutil.ReadFile(ctx.ConfigPath)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(filedata, &config)
	if err != nil {
		return config, err
	}

	// Validate config expectations
	if config.Id == "" {
		return config, errors.New("no vault id found in ldconfig")
	} else if config.Name == "" {
		return config, errors.New("no vault name found in ldconfig")
	}

	return config, nil
}
