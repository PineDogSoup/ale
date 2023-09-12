package types

type LogEvent struct {
	Address    string
	Name       string
	Indexed    []string
	NonIndexed string
}
