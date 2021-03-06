package config

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	cmn "github.com/tendermint/tendermint/libs/common"
	"path/filepath"
	"text/template"
)

const defaultConfigTemplate = `# This is a TOML config file.
# For more information, see https://github.com/toml-lang/toml

##### ethereum config options #####
# Ethereum plasma contract address
ethereum_plasma_contract_address = "{{ .EthPlasmaContractAddr }}"

# Node URL for eth client
ethereum_nodeurl = "{{ .EthNodeURL }}"

# Number of Ethereum blocks until a transaction is considered final
ethereum_finality = "{{ .EthBlockFinality }}"`

// Must match the above defaultConfigTemplate
type PlasmaConfig struct {
	EthPlasmaContractAddr string `mapstructure:"ethereum_plasma_contract_address"`
	EthNodeURL            string `mapstructure:"ethereum_nodeurl"`
	EthBlockFinality      string `mapstructure:"ethereum_finality"`
}

var configTemplate *template.Template

func init() {
	var err error
	tmpl := template.New("plasmaConfigFileTemplate")
	if configTemplate, err = tmpl.Parse(defaultConfigTemplate); err != nil {
		panic(err)
	}
}

func DefaultPlasmaConfig() PlasmaConfig {
	return PlasmaConfig{"", "http://localhost:8545", "0"}
}

// parses the plasma.toml file and unmarshals it into a Config struct
func ParsePlasmaConfigFromViper() (PlasmaConfig, error) {
	config := DefaultPlasmaConfig()
	err := viper.Unmarshal(&config)
	return config, err
}

// WriteConfigFile renders config using the template and writes it to configFilePath.
func WritePlasmaConfigFile(configFilePath string, config PlasmaConfig) {
	var buffer bytes.Buffer

	if err := configTemplate.Execute(&buffer, &config); err != nil {
		panic(err)
	}

	if err := cmn.EnsureDir(filepath.Dir(configFilePath), 0600); err != nil {
		fmt.Printf("ERROR: failed to create directory: %s, recieved error: { %s }", filepath.Dir(configFilePath), err)
	}

	// 0600 for owner only read+write permissions
	cmn.MustWriteFile(configFilePath, buffer.Bytes(), 0600)
}
