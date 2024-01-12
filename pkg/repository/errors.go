package repository

import "errors"

// ErrNoNextStep は、ワークフローの次のステップが存在しない場合に使用されるエラーです。
var ErrNoNextStep = errors.New("no next step in workflow")
