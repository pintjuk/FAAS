package function

import (
	"testing"
)

func fac(val int8) (int8, error) {
	switch val {
	case 0:
		return 1, nil
	default:
		res, err := fac(val - 1)
		return res * val, err
	}
}

func TestRuning(t *testing.T) {
	RunFunc([]string{"value"}, fac)
}
