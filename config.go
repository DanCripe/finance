package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var cfgFile = "financial.yml"

type config struct {
	StartYear int    `yaml:"start_year"`
	EndYear   int    `yaml:"end_year"`
	Age       int    `yaml:"age"`
	SpouseAge int    `yaml:"spouse_age"`
	Assets    assets `yaml:"assets"`
	Income    income `yaml:"income"`
	// COLA      int64  `yaml:"cola"`
	Expenses     int64 `yaml:"expenses"`
	TaxRate      int64 `yaml:"tax_rate"`
	StdDeduction int64 `yaml:"std_deduction"`
}

type assets struct {
	OrdinaryIncome        int64 `yaml:"ordinary_income"`
	TaxFree               int64 `yaml:"tax_free"`
	CapitalGains          int64 `yaml:"capital_gains"`
	CapitalGainsCostBasis int64 `yaml:"cost_basis"`
}

type income struct {
	SocialSecurity          int64 `yaml:"social_security"`
	SocialSecurityAge       int   `yaml:"social_security_age"`
	SpouseSocialSecurity    int64 `yaml:"spouse_social_security"`
	SpouseSocialSecurityAge int   `yaml:"spouse_social_security_age"`
}

func readConfig() (*config, error) {
	buf, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return nil, err
	}

	var cfg config

	err = yaml.Unmarshal(buf, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func printConfig(cfg *config) {
	data, _ := yaml.Marshal(cfg)

	fmt.Printf("%s", string(data))
}
