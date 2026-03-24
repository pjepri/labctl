package playground

import (
	"testing"

	"github.com/iximiuz/labctl/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilterTasksByKind(t *testing.T) {
	tasks := map[string]api.PlayTask{
		"init-1":    {Name: "init-1", Init: true, Status: api.PlayTaskStatusCompleted},
		"init-2":    {Name: "init-2", Init: true, Status: api.PlayTaskStatusRunning},
		"helper-1":  {Name: "helper-1", Helper: true, Status: api.PlayTaskStatusCompleted},
		"regular-1": {Name: "regular-1", Status: api.PlayTaskStatusRunning},
		"regular-2": {Name: "regular-2", Status: api.PlayTaskStatusCompleted},
	}

	t.Run("empty kind returns all tasks", func(t *testing.T) {
		result := filterTasksByKind(tasks, "")
		assert.Len(t, result, 5)
	})

	t.Run("kind=init returns only init tasks", func(t *testing.T) {
		result := filterTasksByKind(tasks, "init")
		assert.Len(t, result, 2)
		for _, task := range result {
			assert.True(t, task.Init)
		}
	})

	t.Run("kind=helper returns only helper tasks", func(t *testing.T) {
		result := filterTasksByKind(tasks, "helper")
		assert.Len(t, result, 1)
		for _, task := range result {
			assert.True(t, task.Helper)
		}
	})

	t.Run("kind=regular returns only regular tasks", func(t *testing.T) {
		result := filterTasksByKind(tasks, "regular")
		assert.Len(t, result, 2)
		for _, task := range result {
			assert.False(t, task.Init)
			assert.False(t, task.Helper)
		}
	})
}

func TestCountFinishedTasks(t *testing.T) {
	t.Run("mixed statuses", func(t *testing.T) {
		tasks := map[string]api.PlayTask{
			"t1": {Status: api.PlayTaskStatusCompleted},
			"t2": {Status: api.PlayTaskStatusRunning},
			"t3": {Status: api.PlayTaskStatusFailed},
			"t4": {Status: api.PlayTaskStatusCreated},
		}
		finished, total := countFinishedTasks(tasks)
		assert.Equal(t, 2, finished)
		assert.Equal(t, 4, total)
	})

	t.Run("all completed", func(t *testing.T) {
		tasks := map[string]api.PlayTask{
			"t1": {Status: api.PlayTaskStatusCompleted},
			"t2": {Status: api.PlayTaskStatusCompleted},
		}
		finished, total := countFinishedTasks(tasks)
		assert.Equal(t, 2, finished)
		assert.Equal(t, 2, total)
	})

	t.Run("none finished", func(t *testing.T) {
		tasks := map[string]api.PlayTask{
			"t1": {Status: api.PlayTaskStatusRunning},
			"t2": {Status: api.PlayTaskStatusBlocked},
		}
		finished, total := countFinishedTasks(tasks)
		assert.Equal(t, 0, finished)
		assert.Equal(t, 2, total)
	})

	t.Run("empty map", func(t *testing.T) {
		tasks := map[string]api.PlayTask{}
		finished, total := countFinishedTasks(tasks)
		assert.Equal(t, 0, finished)
		assert.Equal(t, 0, total)
	})
}

func TestTaskIsFinished(t *testing.T) {
	assert.True(t, taskIsFinished(api.PlayTask{Status: api.PlayTaskStatusCompleted}))
	assert.True(t, taskIsFinished(api.PlayTask{Status: api.PlayTaskStatusFailed}))
	assert.False(t, taskIsFinished(api.PlayTask{Status: api.PlayTaskStatusRunning}))
	assert.False(t, taskIsFinished(api.PlayTask{Status: api.PlayTaskStatusCreated}))
	assert.False(t, taskIsFinished(api.PlayTask{Status: api.PlayTaskStatusBlocked}))
	assert.False(t, taskIsFinished(api.PlayTask{Status: api.PlayTaskStatusNone}))
}

func TestFormatTaskStatus(t *testing.T) {
	tests := []struct {
		status   api.PlayTaskStatus
		expected string
	}{
		{api.PlayTaskStatusNone, "none"},
		{api.PlayTaskStatusCreated, "created"},
		{api.PlayTaskStatusBlocked, "blocked"},
		{api.PlayTaskStatusRunning, "running"},
		{api.PlayTaskStatusFailed, "failed"},
		{api.PlayTaskStatusCompleted, "completed"},
		{api.PlayTaskStatus(99), "unknown"},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, formatTaskStatus(tt.status))
	}
}

func TestTasksOptionsValidate(t *testing.T) {
	t.Run("valid defaults", func(t *testing.T) {
		opts := &tasksOptions{output: "table"}
		require.NoError(t, opts.validate())
	})

	t.Run("valid kind=init", func(t *testing.T) {
		opts := &tasksOptions{output: "table", kind: "init"}
		require.NoError(t, opts.validate())
	})

	t.Run("valid kind=helper", func(t *testing.T) {
		opts := &tasksOptions{output: "table", kind: "helper"}
		require.NoError(t, opts.validate())
	})

	t.Run("valid kind=regular", func(t *testing.T) {
		opts := &tasksOptions{output: "table", kind: "regular"}
		require.NoError(t, opts.validate())
	})

	t.Run("invalid kind", func(t *testing.T) {
		opts := &tasksOptions{output: "table", kind: "bogus"}
		assert.Error(t, opts.validate())
	})

	t.Run("invalid output", func(t *testing.T) {
		opts := &tasksOptions{output: "xml"}
		assert.Error(t, opts.validate())
	})

	t.Run("--init sets kind to init", func(t *testing.T) {
		opts := &tasksOptions{output: "table", init: true}
		require.NoError(t, opts.validate())
		assert.Equal(t, "init", opts.kind)
	})

	t.Run("--init and --kind conflict", func(t *testing.T) {
		opts := &tasksOptions{output: "table", init: true, kind: "helper"}
		assert.Error(t, opts.validate())
		assert.Contains(t, opts.validate().Error(), "cannot use both --init and --kind")
	})

	t.Run("all output formats", func(t *testing.T) {
		for _, format := range []string{"table", "json", "name", "none"} {
			opts := &tasksOptions{output: format}
			require.NoError(t, opts.validate(), "format %s should be valid", format)
		}
	})
}
