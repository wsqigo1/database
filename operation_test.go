package redis

import "testing"

func TestOperatorType(t *testing.T) {
	OperatorType()
}

func TestCreateBasic(t *testing.T) {
	CreateBasic()
}

func TestCreateMulti(t *testing.T) {
	CreateMulti()
}

func TestCreateInBatches(t *testing.T) {
	CreateInBatches()
}

func TestUpsert(t *testing.T) {
	Upsert()
}

func TestDefaultValue(t *testing.T) {
	DefaultValue()
}

func TestSelectOmit(t *testing.T) {
	SelectOmit()
}

func TestCreateHook(t *testing.T) {
	CreateHook()
}
