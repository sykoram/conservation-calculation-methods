package main

type MsaColumn string
type MethodFunc func(MsaColumn) float64

var Methods = map[string]MethodFunc{
	"zero": Zero,
}

func GetMethodNames() []string {
	keys := make([]string, len(Methods))
	i := 0
	for k := range Methods {
		keys[i] = k
		i++
	}
	return keys
}

func Zero(msa MsaColumn) float64 {
	return 0
}
