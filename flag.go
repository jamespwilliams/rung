package rung

type BoolFlag struct {
	Value  bool
	WasSet bool
}

type StringFlag struct {
	Value  string
	WasSet bool
}

type Int64Flag struct {
	Value  int64
	WasSet bool
}

type Uint64Flag struct {
	Value  uint64
	WasSet bool
}

type IntFlag struct {
	Value  int
	WasSet bool
}

type UintFlag struct {
	Value  uint
	WasSet bool
}

type Float64Flag struct {
	Value  float64
	WasSet bool
}
