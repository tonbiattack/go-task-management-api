package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tonbiattack/go-task-management-api/pkg/model"
	"github.com/tonbiattack/go-task-management-api/pkg/testutils"
)

func TestWorkflowRepository_Create(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := NewWorkflowRepository(db)

	workflow := &model.Workflow{
		ID:          "test-workflow-1",
		Name:        "テストワークフロー",
		Description: "テスト用のワークフローです",
	}

	err := repo.Create(workflow)
	assert.NoError(t, err)

	// データベースから取得して確認
	var savedWorkflow model.Workflow
	err = db.First(&savedWorkflow, "id = ?", "test-workflow-1").Error
	assert.NoError(t, err)
	assert.Equal(t, "テストワークフロー", savedWorkflow.Name)
	assert.Equal(t, "テスト用のワークフローです", savedWorkflow.Description)
}

func TestWorkflowRepository_FindByID(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := NewWorkflowRepository(db)

	// テストデータを作成
	testutils.CreateTestWorkflow(db, "test-workflow-1", "テストワークフロー", "テスト用のワークフローです")

	// ワークフローを取得
	workflow, err := repo.FindByID("test-workflow-1")
	assert.NoError(t, err)
	assert.Equal(t, "test-workflow-1", workflow.ID)
	assert.Equal(t, "テストワークフロー", workflow.Name)

	// 存在しないワークフローの場合
	_, err = repo.FindByID("non-existent")
	assert.Error(t, err)
}

func TestWorkflowRepository_Update(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := NewWorkflowRepository(db)

	// テストデータを作成
	originalWorkflow := testutils.CreateTestWorkflow(db, "test-workflow-1", "元の名前", "元の説明")

	// ワークフローを更新
	originalWorkflow.Name = "更新された名前"
	originalWorkflow.Description = "更新された説明"

	err := repo.Update(originalWorkflow)
	assert.NoError(t, err)

	// データベースから取得して確認
	updatedWorkflow, err := repo.FindByID("test-workflow-1")
	assert.NoError(t, err)
	assert.Equal(t, "更新された名前", updatedWorkflow.Name)
	assert.Equal(t, "更新された説明", updatedWorkflow.Description)
}

func TestWorkflowRepository_Delete(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := NewWorkflowRepository(db)

	// テストデータを作成
	testutils.CreateTestWorkflow(db, "test-workflow-1", "削除されるワークフロー", "削除テスト")

	// ワークフローが存在することを確認
	_, err := repo.FindByID("test-workflow-1")
	assert.NoError(t, err)

	// ワークフローを削除
	err = repo.Delete("test-workflow-1")
	assert.NoError(t, err)

	// ワークフローが削除されたことを確認
	_, err = repo.FindByID("test-workflow-1")
	assert.Error(t, err)
}
