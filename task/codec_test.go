package task_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"cthulhu/task"
)

func TestMarshalUnmarshalWMTaskContext(t *testing.T) {
	var ctx *task.WelcomeMessageTaskContext = &task.WelcomeMessageTaskContext{
		UsersProcessed: task.UserSet{
			"1": struct{}{},
			"2": struct{}{},
		},
	}

	encoded := task.MarshalWMTaskContext(ctx)
	res, err := task.UnmarshalWMTaskContext(encoded)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(ctx, res) {
		t.Fatalf("error during marshal/unmarshal: want %v, got %v", ctx, res)
	}
}
