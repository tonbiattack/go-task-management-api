package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tonbiattack/go-task-management-api/pkg/model"
	"github.com/tonbiattack/go-task-management-api/pkg/testutils"
)

func TestWorkflowStepRepository_Create(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := NewWorkflowStepRepository(db)

	// 前提条件としてワークフローを作成
	testutils.CreateTestWorkflow(db, "workflow-1", "テストワークフロー", "説明")

	step := &model.WorkflowStep{
		ID:          "step-1",
		WorkflowID:  "workflow-1",
		Name:        "ステップ1",
		Description: "最初のステップ",
		Order:       1,
	}

	err := repo.Create(step)
	assert.NoError(t, err)

	// データベースから取得して確認
	var savedStep model.WorkflowStep
	err = db.First(&savedStep, "id = ?", "step-1").Error
	assert.NoError(t, err)
	assert.Equal(t, "ステップ1", savedStep.Name)
	assert.Equal(t, 1, savedStep.Order)
}

func TestWorkflowStepRepository_FindByID(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := NewWorkflowStepRepository(db)

	// 前提条件としてワークフローを作成
	testutils.CreateTestWorkflow(db, "workflow-1", "テストワークフロー", "説明")
	testutils.CreateTestWorkflowStep(db, "step-1", "workflow-1", "ステップ1", 1)

	// ステップを取得
	step, err := repo.FindByID("step-1")
	assert.NoError(t, err)
	assert.Equal(t, "step-1", step.ID)
	assert.Equal(t, "ステップ1", step.Name)

	// 存在しないステップの場合
	_, err = repo.FindByID("non-existent")
	assert.Error(t, err)
}

func TestWorkflowStepRepository_GetInitialStep(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := NewWorkflowStepRepository(db)

	// 前提条件としてワークフローを作成
	testutils.CreateTestWorkflow(db, "workflow-1", "テストワークフロー", "説明")

	// 複数のステップを作成（順序を意図的に逆にする）
	testutils.CreateTestWorkflowStep(db, "step-3", "workflow-1", "ステップ3", 3)
	testutils.CreateTestWorkflowStep(db, "step-1", "workflow-1", "ステップ1", 1)
	testutils.CreateTestWorkflowStep(db, "step-2", "workflow-1", "ステップ2", 2)

	// 最初のステップを取得
	initialStep, err := repo.GetInitialStep("workflow-1")
	assert.NoError(t, err)
	assert.Equal(t, "step-1", initialStep.ID)
	assert.Equal(t, 1, initialStep.Order)

	// 存在しないワークフローの場合
	_, err = repo.GetInitialStep("non-existent")
	assert.Error(t, err)
}

func TestWorkflowStepRepository_GetNextStep(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := NewWorkflowStepRepository(db)

	// 前提条件としてワークフローを作成
	testutils.CreateTestWorkflow(db, "workflow-1", "テストワークフロー", "説明")

	// 複数のステップを作成
	testutils.CreateTestWorkflowStep(db, "step-1", "workflow-1", "ステップ1", 1)
	testutils.CreateTestWorkflowStep(db, "step-2", "workflow-1", "ステップ2", 2)
	testutils.CreateTestWorkflowStep(db, "step-3", "workflow-1", "ステップ3", 3)

	// step-1の次のステップを取得
	nextStep, err := repo.GetNextStep("step-1")
	assert.NoError(t, err)
	assert.Equal(t, "step-2", nextStep.ID)
	assert.Equal(t, 2, nextStep.Order)

	// step-2の次のステップを取得
	nextStep, err = repo.GetNextStep("step-2")
	assert.NoError(t, err)
	assert.Equal(t, "step-3", nextStep.ID)
	assert.Equal(t, 3, nextStep.Order)

	// 最後のステップの次は存在しない
	_, err = repo.GetNextStep("step-3")
	assert.Error(t, err)
	assert.Equal(t, ErrNoNextStep, err)

	// 存在しないステップの場合
	_, err = repo.GetNextStep("non-existent")
	assert.Error(t, err)
}

func TestWorkflowStepRepository_Update(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := NewWorkflowStepRepository(db)

	// 前提条件としてワークフローを作成
	testutils.CreateTestWorkflow(db, "workflow-1", "テストワークフロー", "説明")
	originalStep := testutils.CreateTestWorkflowStep(db, "step-1", "workflow-1", "元の名前", 1)

	// ステップを更新
	originalStep.Name = "更新された名前"
	originalStep.Description = "更新された説明"
	originalStep.Order = 2

	err := repo.Update(originalStep)
	assert.NoError(t, err)

	// データベースから取得して確認
	updatedStep, err := repo.FindByID("step-1")
	assert.NoError(t, err)
	assert.Equal(t, "更新された名前", updatedStep.Name)
	assert.Equal(t, "更新された説明", updatedStep.Description)
	assert.Equal(t, 2, updatedStep.Order)
}

func TestWorkflowStepRepository_Delete(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := NewWorkflowStepRepository(db)

	// 前提条件としてワークフローを作成
	testutils.CreateTestWorkflow(db, "workflow-1", "テストワークフロー", "説明")
	testutils.CreateTestWorkflowStep(db, "step-1", "workflow-1", "削除されるステップ", 1)

	// ステップが存在することを確認
	_, err := repo.FindByID("step-1")
	assert.NoError(t, err)

	// ステップを削除
	err = repo.Delete("step-1")
	assert.NoError(t, err)

	// ステップが削除されたことを確認
	_, err = repo.FindByID("step-1")
	assert.Error(t, err)
}
