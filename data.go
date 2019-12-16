package main

type data struct {
	StartYear    int
	EndYear      int
	Age          int
	SpouseAge    int
	Assets       assetsFloat
	Income       incomeFloat
	Expenses     float64
	TaxRate      float64
	StdDeduction float64
}

type assetsFloat struct {
	OrdinaryIncome        float64
	TaxFree               float64
	CapitalGains          float64
	CapitalGainsCostBasis float64
}

type incomeFloat struct {
	SocialSecurity          float64
	SocialSecurityAge       int
	SpouseSocialSecurity    float64
	SpouseSocialSecurityAge int
}

func convert(cfg *config) data {
	var result data

	result.StartYear = cfg.StartYear
	result.EndYear = cfg.EndYear
	result.Age = cfg.Age
	result.SpouseAge = cfg.SpouseAge
	result.Assets = floatAssets(cfg.Assets)
	result.Income = floatIncome(cfg.Income)
	// result.COLA = float64(cfg.COLA)/100.0 + float64(1.0)
	result.Expenses = float64(cfg.Expenses)
	result.TaxRate = float64(cfg.TaxRate) / 100.0
	result.StdDeduction = float64(cfg.StdDeduction)

	return result
}

func floatAssets(a assets) assetsFloat {
	var result assetsFloat
	result.OrdinaryIncome = float64(a.OrdinaryIncome)
	result.TaxFree = float64(a.TaxFree)
	result.CapitalGains = float64(a.CapitalGains)
	result.CapitalGainsCostBasis = float64(a.CapitalGainsCostBasis)
	return result
}

func floatIncome(i income) incomeFloat {
	var result incomeFloat
	result.SocialSecurityAge = i.SocialSecurityAge
	result.SpouseSocialSecurityAge = i.SpouseSocialSecurityAge
	result.SocialSecurity = float64(i.SocialSecurity * 12)
	result.SpouseSocialSecurity = float64(i.SpouseSocialSecurity * 12)
	return result
}
