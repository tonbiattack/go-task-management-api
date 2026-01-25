package repository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tonbiattack/go-task-management-api/pkg/model"
	"github.com/tonbiattack/go-task-management-api/pkg/testutils"
)

func TestTaskRepository_CreateTask(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := NewTaskRepository(db)

	status := "pending"
	task := &model.Task{
		ID:          "test-task-1",
		Title:       "テストタスク",
		Description: "テスト用のタスクです",
		Status:      &status,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := repo.CreateTask(task)
	assert.NoError(t, err)

	// データベースから取得して確認
	var savedTask model.Task
	err = db.First(&savedTask, "id = ?", "test-task-1").Error
	assert.NoError(t, err)
	assert.Equal(t, "テストタスク", savedTask.Title)
	assert.Equal(t, "テスト用のタスクです", savedTask.Description)
}

func TestTaskRepository_GetTaskByID(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := NewTaskRepository(db)

	// テストデータを作成
	testutils.CreateTestTask(db, "test-task-1", "テストタスク", "テスト用のタスクです")

	// タスクを取得
	task, err := repo.GetTaskByID("test-task-1")
	assert.NoError(t, err)
	assert.Equal(t, "test-task-1", task.ID)
	assert.Equal(t, "テストタスク", task.Title)

	// 存在しないタスクの場合
	_, err = repo.GetTaskByID("non-existent")
	assert.Error(t, err)
}

func TestTaskRepository_GetAllTasks(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := NewTaskRepository(db)

	// テストデータを作成
	testutils.CreateTestTask(db, "task-1", "タスク1", "説明1")
	testutils.CreateTestTask(db, "task-2", "タスク2", "説明2")
	testutils.CreateTestTask(db, "task-3", "タスク3", "説明3")

	// すべてのタスクを取得
	tasks, err := repo.GetAllTasks()
	assert.NoError(t, err)
	assert.Len(t, tasks, 3)

	// タスクの順序は保証されないので、IDでソートして確認
	taskIDs := make([]string, len(tasks))
	for i, task := range tasks {
		taskIDs[i] = task.ID
	}
	assert.Contains(t, taskIDs, "task-1")
	assert.Contains(t, taskIDs, "task-2")
	assert.Contains(t, taskIDs, "task-3")
}

func TestTaskRepository_UpdateTask(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := NewTaskRepository(db)

	// テストデータを作成
	originalTask := testutils.CreateTestTask(db, "test-task-1", "元のタイトル", "元の説明")

	// タスクを更新
	newStatus := "completed"
	originalTask.Title = "更新されたタイトル"
	originalTask.Description = "更新された説明"
	originalTask.Status = &newStatus

	err := repo.UpdateTask(originalTask)
	assert.NoError(t, err)

	// データベースから取得して確認
	updatedTask, err := repo.GetTaskByID("test-task-1")
	assert.NoError(t, err)
	assert.Equal(t, "更新されたタイトル", updatedTask.Title)
	assert.Equal(t, "更新された説明", updatedTask.Description)
	assert.Equal(t, "completed", *updatedTask.Status)
}

func TestTaskRepository_DeleteTask(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := NewTaskRepository(db)

	// テストデータを作成
	testutils.CreateTestTask(db, "test-task-1", "削除されるタスク", "削除テスト")

	// タスクが存在することを確認
	_, err := repo.GetTaskByID("test-task-1")
	assert.NoError(t, err)

	// タスクを削除
	err = repo.DeleteTask("test-task-1")
	assert.NoError(t, err)

	// タスクが削除されたことを確認
	_, err = repo.GetTaskByID("test-task-1")
	assert.Error(t, err)
}

func TestTaskRepository_DeleteNonExistentTask(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := NewTaskRepository(db)

	// 存在しないタスクを削除しようとする
	err := repo.DeleteTask("non-existent")
	// GORMでは存在しないレコードの削除はエラーにならない
	assert.NoError(t, err)
}

func TestTaskRepository_CreateTaskWithNilStatus(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := NewTaskRepository(db)

	task := &model.Task{
		ID:          "test-task-nil-status",
		Title:       "ステータスなしタスク",
		Description: "ステータスがnilのタスク",
		Status:      nil, // nilステータス
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := repo.CreateTask(task)
	assert.NoError(t, err)

	// データベースから取得して確認
	savedTask, err := repo.GetTaskByID("test-task-nil-status")
	require.NoError(t, err)

	// デフォルト値が設定されているか確認
	if savedTask.Status != nil {
		assert.Equal(t, "pending", *savedTask.Status)
	}
}
