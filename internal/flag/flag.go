package flag

type FlagVar struct {
	IsSet bool
	Value string
}

func (f *FlagVar) String() string {
	return f.Value
}

func (f *FlagVar) Set(value string) error {
	f.Value = value
	f.IsSet = true
	return nil
}
