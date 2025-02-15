package cmd

import (
	"github.com/AlecAivazis/survey/v2"
)

// The following function variables allow us to override their default implementations during testing.
// This is useful for simulating user interactions and configuration behavior without relying on external
// dependencies or actual user input.
//   - surveyAskOne: References survey.AskOne, which can be replaced with a mock function to simulate user responses.
//   - GetConfigFunc and UpdateConfigFunc: Reference the actual GetConfig and UpdateConfig functions,
//     allowing tests to substitute them with in-memory versions or mocks.
var surveyAskOne = survey.AskOne

var GetConfigFunc = GetConfig
var UpdateConfigFunc = UpdateConfig
