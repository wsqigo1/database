package redis

import "testing"

func TestMigrate(t *testing.T) {
	Migrate()
}

func TestPointerDiff(t *testing.T) {
	PointerDiff()
}

func TestCustomType(t *testing.T) {
	CustomType()
}

func TestIAndCCreate(t *testing.T) {
	IAndCCreate()
}

func TestServiceCURD(t *testing.T) {
	ServiceCURD()
}

func TestPaperCurd(t *testing.T) {
	PaperCurd()
}

func TestCustomSerializer(t *testing.T) {
	CustomSerializer()
}
