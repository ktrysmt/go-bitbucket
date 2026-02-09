package mock_tests

import (
	"errors"
	"testing"

	go_bitbucket "github.com/ktrysmt/go-bitbucket"
	"github.com/ktrysmt/go-bitbucket/mockgen"
	"github.com/stretchr/testify/assert"

	"go.uber.org/mock/gomock"
)

func TestMockRepositoryPipelineVariables_List_Success(t *testing.T) {
	t.Parallel()
	// 1. Create a gomock controller
	// The controller manages the mock objects and their expectations
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // Assert that all expected calls were made

	// 2. Create an instance of the mock
	mockRepositoryInst := mockgen.NewMockrepository(ctrl)

	// 3. Set up expectations on the mock
	listPipelineVarsOpts := go_bitbucket.RepositoryPipelineVariablesOptions{
		Owner:    "testworkspace",
		RepoSlug: "testrepo",
	}

	outPipelineVar := go_bitbucket.PipelineVariable{
		Type:    "pipeline_variable",
		Uuid:    "6b98a093-21e3-4e15-ad48-f06aad1d2399",
		Key:     "test-key",
		Value:   "test-value",
		Secured: false,
	}

	var outPipelineVarList []go_bitbucket.PipelineVariable
	outPipelineVarList = append(outPipelineVarList, outPipelineVar)
	outPipelineVars := &go_bitbucket.PipelineVariables{
		Page:      1,
		Size:      1,
		MaxDepth:  1,
		Next:      "test-next",
		Variables: outPipelineVarList,
	}

	mockRepositoryInst.EXPECT().
		ListPipelineVariables(listPipelineVarsOpts). // Input call with opt
		Return(outPipelineVars, nil).                // And return these values
		Times(1)                                     // Expect it to be called only once

	actualPipelineVars, actualErr := mockRepositoryInst.ListPipelineVariables(listPipelineVarsOpts)

	assert.Nil(t, actualErr, "No errors should be thrown, but got: %v", actualErr)
	assert.GreaterOrEqual(t, actualPipelineVars.Size, 1)
	assert.GreaterOrEqual(t, len(actualPipelineVars.Variables), 1)
	assert.Equal(t, outPipelineVar, actualPipelineVars.Variables[0])

	// ctrl.Finish() in the defer statement will automatically verify mock expectations.
}

func TestMockRepositoryPipelineVariables_List_Error(t *testing.T) {
	t.Parallel()
	// 1. Create a gomock controller
	// The controller manages the mock objects and their expectations
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // Assert that all expected calls were made

	// 2. Create an instance of the mock
	mockRepositoryInst := mockgen.NewMockrepository(ctrl)
	expectedMockError := errors.New("Not Found")

	// 3. Set up expectations on the mock

	listPipelineVarsOpts := go_bitbucket.RepositoryPipelineVariablesOptions{
		Owner:    "testworkspace-not-found",
		RepoSlug: "testrepo-not-found",
	}

	mockRepositoryInst.EXPECT().
		ListPipelineVariables(listPipelineVarsOpts). // Input call with opt
		Return(nil, expectedMockError).              // And return these values
		Times(1)                                     // Expect it to be called only once

	actualPipelineVars, actualErr := mockRepositoryInst.ListPipelineVariables(listPipelineVarsOpts)

	assert.Nil(t, actualPipelineVars, "actualPipelineVars should have been nil, but got: %v", actualPipelineVars)
	assert.NotNil(t, actualErr)
	assert.Equal(t, expectedMockError.Error(), actualErr.Error())
}

func TestMockRepositoryPipelineVariable_Get_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepositoryInst := mockgen.NewMockrepository(ctrl)

	inGetPipelineVarOpts := go_bitbucket.RepositoryPipelineVariableOptions{
		Owner:    "testworkspace",
		RepoSlug: "testrepo",
		Uuid:     "6b98a093-21e3-4e15-ad48-f06aad1d2399",
	}

	expectedPipelineVar := &go_bitbucket.PipelineVariable{
		Type:  "pipeline_variable",
		Uuid:  "6b98a093-21e3-4e15-ad48-f06aad1d2399",
		Key:   "test-key",
		Value: "test-value",
	}

	mockRepositoryInst.EXPECT().
		GetPipelineVariable(inGetPipelineVarOpts).
		Return(expectedPipelineVar, nil).
		Times(1)

	actualPipelineVar, actualErr := mockRepositoryInst.GetPipelineVariable(inGetPipelineVarOpts)

	assert.Nil(t, actualErr, "results should have been nil: %v", actualPipelineVar)
	assert.Equal(t, expectedPipelineVar, actualPipelineVar)
}

func TestMockRepositoryPipelineVariable_Get_Error(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepositoryInst := mockgen.NewMockrepository(ctrl)
	expectedMockError := errors.New("Not Found")

	inGetPipelineVarOpts := go_bitbucket.RepositoryPipelineVariableOptions{
		Owner:    "testworkspace",
		RepoSlug: "testrepo",
		Uuid:     "6b98a093-21e3-4e15-ad48-f06aad1d2369",
	}

	mockRepositoryInst.EXPECT().
		GetPipelineVariable(inGetPipelineVarOpts).
		Return(nil, expectedMockError).
		Times(1)

	actualPipelineVar, actualErr := mockRepositoryInst.GetPipelineVariable(inGetPipelineVarOpts)

	assert.Nil(t, actualPipelineVar, "results should have been nil: %v", actualPipelineVar)
	assert.NotNil(t, actualErr)
}

func TestMockRepositoryPipelineVariable_Update_Success(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepositoryInst := mockgen.NewMockrepository(ctrl)

	inPipelineVarUpdateOpts := go_bitbucket.RepositoryPipelineVariableOptions{
		Owner:    "testworkspace",
		RepoSlug: "testrepo",
		Uuid:     "6b98a093-21e3-4e15-ad48-f06aad1d2399",
	}

	expectedNewPipelineUpdateVar := &go_bitbucket.PipelineVariable{
		Type:  "pipeline_variable",
		Uuid:  "6b98a093-21e3-4e15-ad48-f06aad1d2399",
		Key:   "test-key",
		Value: "test-value-new",
	}

	mockRepositoryInst.EXPECT().
		UpdatePipelineVariable(inPipelineVarUpdateOpts).
		Return(expectedNewPipelineUpdateVar, nil).
		Times(1)

	actualPipelineUpdateVar, actualUpdateErr := mockRepositoryInst.UpdatePipelineVariable(inPipelineVarUpdateOpts)

	assert.Nil(t, actualUpdateErr, "No error should have been thrown, but got: %v", actualUpdateErr)
	assert.Equal(t, expectedNewPipelineUpdateVar, actualPipelineUpdateVar)
}

func TestMockRepositoryPipelineVariable_Update_Error(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepositoryInst := mockgen.NewMockrepository(ctrl)
	expectedMockUpdateError := errors.New("Not Found")

	inPipelineVarUpdateOpts := go_bitbucket.RepositoryPipelineVariableOptions{
		Owner:    "testworkspace",
		RepoSlug: "testrepo",
		Uuid:     "6b98a093-21e3-4e15-ad48-f06aad1d2399",
	}

	mockRepositoryInst.EXPECT().
		UpdatePipelineVariable(inPipelineVarUpdateOpts).
		Return(nil, expectedMockUpdateError).
		Times(1)

	actualPipelineUpdateVar, actualUpdateErr := mockRepositoryInst.UpdatePipelineVariable(inPipelineVarUpdateOpts)

	assert.NotNil(t, actualUpdateErr, "An error should been thrown")
	assert.Nil(t, actualPipelineUpdateVar, "An error should have been thrown, but got: %v", actualPipelineUpdateVar)
}
