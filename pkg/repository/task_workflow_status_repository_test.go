package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tonbiattack/go-task-management-api/pkg/model"
	"github.com/tonbiattack/go-task-management-api/pkg/testutils"
)

func TestTaskWorkflowStatusRepository_Create(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := NewTaskWorkflowStatusRepository(db)

	// 前提条件としてタスク、ワークフロー、ステップを作成
	testutils.CreateTestTask(db, "task-1", "テストタスク", "説明")
	testutils.CreateTestWorkflow(db, "workflow-1", "テストワークフロー", "説明")
	testutils.CreateTestWorkflowStep(db, "step-1", "workflow-1", "ステップ1", 1)

	status := &model.TaskWorkflowStatus{
		ID:             "status-1",
		TaskID:         "task-1",
		WorkflowStepID: "step-1",
		Status:         "IN_PROGRESS",
	}

	err := repo.Create(status)
	assert.NoError(t, err)

	// データベースから取得して確認
	var savedStatus model.TaskWorkflowStatus
	err = db.First(&savedStatus, "id = ?", "status-1").Error
	assert.NoError(t, err)
	assert.Equal(t, "task-1", savedStatus.TaskID)
	assert.Equal(t, "step-1", savedStatus.WorkflowStepID)
	assert.Equal(t, "IN_PROGRESS", savedStatus.Status)
}

func TestTaskWorkflowStatusRepository_FindByID(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := NewTaskWorkflowStatusRepository(db)

	// 前提条件としてタスク、ワークフロー、ステップを作成
	testutils.CreateTestTask(db, "task-1", "テストタスク", "説明")
	testutils.CreateTestWorkflow(db, "workflow-1", "テストワークフロー", "説明")
	testutils.CreateTestWorkflowStep(db, "step-1", "workflow-1", "ステップ1", 1)
	testutils.CreateTestTaskWorkflowStatus(db, "status-1", "task-1", "step-1", "IN_PROGRESS")

	// ステータスを取得
	status, err := repo.FindByID("status-1")
	assert.NoError(t, err)
	assert.Equal(t, "status-1", status.ID)
	assert.Equal(t, "task-1", status.TaskID)
	assert.Equal(t, "IN_PROGRESS", status.Status)

	// 存在しないステータスの場合
	_, err = repo.FindByID("non-existent")
	assert.Error(t, err)
}

func TestTaskWorkflowStatusRepository_FindByTaskID(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := NewTaskWorkflowStatusRepository(db)

	// 前提条件としてタスク、ワークフロー、ステップを作成
	testutils.CreateTestTask(db, "task-1", "テストタスク", "説明")
	testutils.CreateTestWorkflow(db, "workflow-1", "テストワークフロー", "説明")
	testutils.CreateTestWorkflowStep(db, "step-1", "workflow-1", "ステップ1", 1)
	testutils.CreateTestTaskWorkflowStatus(db, "status-1", "task-1", "step-1", "IN_PROGRESS")

	// タスクIDでステータスを取得
	status, err := repo.FindByTaskID("task-1")
	assert.NoError(t, err)
	assert.Equal(t, "task-1", status.TaskID)
	assert.Equal(t, "IN_PROGRESS", status.Status)

	// 存在しないタスクIDの場合
	_, err = repo.FindByTaskID("non-existent")
	assert.Error(t, err)
}

func TestTaskWorkflowStatusRepository_Update(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := NewTaskWorkflowStatusRepository(db)

	// 前提条件としてタスク、ワークフロー、ステップを作成
	testutils.CreateTestTask(db, "task-1", "テストタスク", "説明")
	testutils.CreateTestWorkflow(db, "workflow-1", "テストワークフロー", "説明")
	testutils.CreateTestWorkflowStep(db, "step-1", "workflow-1", "ステップ1", 1)
	testutils.CreateTestWorkflowStep(db, "step-2", "workflow-1", "ステップ2", 2)
	originalStatus := testutils.CreateTestTaskWorkflowStatus(db, "status-1", "task-1", "step-1", "IN_PROGRESS")

	// ステータスを更新
	originalStatus.WorkflowStepID = "step-2"
	originalStatus.Status = "COMPLETED"

	err := repo.Update(originalStatus)
	assert.NoError(t, err)

	// データベースから取得して確認
	updatedStatus, err := repo.FindByID("status-1")
	assert.NoError(t, err)
	assert.Equal(t, "step-2", updatedStatus.WorkflowStepID)
	assert.Equal(t, "COMPLETED", updatedStatus.Status)
}

func TestTaskWorkflowStatusRepository_Delete(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := NewTaskWorkflowStatusRepository(db)

	// 前提条件としてタスク、ワークフロー、ステップを作成
	testutils.CreateTestTask(db, "task-1", "テストタスク", "説明")
	testutils.CreateTestWorkflow(db, "workflow-1", "テストワークフロー", "説明")
	testutils.CreateTestWorkflowStep(db, "step-1", "workflow-1", "ステップ1", 1)
	testutils.CreateTestTaskWorkflowStatus(db, "status-1", "task-1", "step-1", "IN_PROGRESS")

	// ステータスが存在することを確認
	_, err := repo.FindByID("status-1")
	assert.NoError(t, err)

	// ステータスを削除
	err = repo.Delete("status-1")
	assert.NoError(t, err)

	// ステータスが削除されたことを確認
	_, err = repo.FindByID("status-1")
	assert.Error(t, err)
}
