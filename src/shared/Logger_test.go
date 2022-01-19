package shared_test

import (
	"fmt"
	"testing"

	"github.com/alikarimii/go_starter/src/shared"
)

func TestNewStandardLogger(t *testing.T) {
	fmt.Println("When a new standard logger is created")
	logger := shared.NewStandardLogger()
	if logger.Verbose() != true {
		t.Errorf("standard logger failed")
	}
	fmt.Println("When a new nil logger is created")
	logger2 := shared.NewNilLogger()
	if logger2.Verbose() != true {
		t.Errorf("nil logger failed")
	}

}
