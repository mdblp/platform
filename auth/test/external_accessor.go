package test

import (
	"context"

	"github.com/tidepool-org/platform/request"
)

type ServerSessionTokenOutput struct {
	Token string
	Error error
}

type ValidateSessionTokenInput struct {
	Context context.Context
	Token   string
}

type ValidateSessionTokenOutput struct {
	Details request.Details
	Error   error
}

type EnsureAuthorizedUserInput struct {
	Context              context.Context
	TargetUserID         string
	AuthorizedPermission string
}

type EnsureAuthorizedUserOutput struct {
	AuthorizedUserID string
	Error            error
}

type ExternalAccessor struct {
	ServerSessionTokenInvocations      int
	ServerSessionTokenStub             func() (string, error)
	ServerSessionTokenOutputs          []ServerSessionTokenOutput
	ServerSessionTokenOutput           *ServerSessionTokenOutput
	ValidateSessionTokenInvocations    int
	ValidateSessionTokenInputs         []ValidateSessionTokenInput
	ValidateSessionTokenStub           func(ctx context.Context, token string) (request.Details, error)
	ValidateSessionTokenOutputs        []ValidateSessionTokenOutput
	ValidateSessionTokenOutput         *ValidateSessionTokenOutput
	EnsureAuthorizedInvocations        int
	EnsureAuthorizedInputs             []context.Context
	EnsureAuthorizedStub               func(ctx context.Context) error
	EnsureAuthorizedOutputs            []error
	EnsureAuthorizedOutput             *error
	EnsureAuthorizedServiceInvocations int
	EnsureAuthorizedServiceInputs      []context.Context
	EnsureAuthorizedServiceStub        func(ctx context.Context) error
	EnsureAuthorizedServiceOutputs     []error
	EnsureAuthorizedServiceOutput      *error
	EnsureAuthorizedUserInvocations    int
	EnsureAuthorizedUserInputs         []EnsureAuthorizedUserInput
	EnsureAuthorizedUserStub           func(ctx context.Context, targetUserID string, authorizedPermission string) (string, error)
	EnsureAuthorizedUserOutputs        []EnsureAuthorizedUserOutput
	EnsureAuthorizedUserOutput         *EnsureAuthorizedUserOutput
}

func NewExternalAccessor() *ExternalAccessor {
	return &ExternalAccessor{}
}

func (e *ExternalAccessor) ServerSessionToken() (string, error) {
	e.ServerSessionTokenInvocations++
	if e.ServerSessionTokenStub != nil {
		return e.ServerSessionTokenStub()
	}
	if len(e.ServerSessionTokenOutputs) > 0 {
		output := e.ServerSessionTokenOutputs[0]
		e.ServerSessionTokenOutputs = e.ServerSessionTokenOutputs[1:]
		return output.Token, output.Error
	}
	if e.ServerSessionTokenOutput != nil {
		return e.ServerSessionTokenOutput.Token, e.ServerSessionTokenOutput.Error
	}
	panic("ServerSessionToken has no output")
}

func (e *ExternalAccessor) ValidateSessionToken(ctx context.Context, token string) (request.Details, error) {
	e.ValidateSessionTokenInvocations++
	e.ValidateSessionTokenInputs = append(e.ValidateSessionTokenInputs, ValidateSessionTokenInput{Context: ctx, Token: token})
	if e.ValidateSessionTokenStub != nil {
		return e.ValidateSessionTokenStub(ctx, token)
	}
	if len(e.ValidateSessionTokenOutputs) > 0 {
		output := e.ValidateSessionTokenOutputs[0]
		e.ValidateSessionTokenOutputs = e.ValidateSessionTokenOutputs[1:]
		return output.Details, output.Error
	}
	if e.ValidateSessionTokenOutput != nil {
		return e.ValidateSessionTokenOutput.Details, e.ValidateSessionTokenOutput.Error
	}
	panic("ValidateSessionToken has no output")
}

func (e *ExternalAccessor) EnsureAuthorized(ctx context.Context) error {
	e.EnsureAuthorizedInvocations++
	e.EnsureAuthorizedInputs = append(e.EnsureAuthorizedInputs, ctx)
	if e.EnsureAuthorizedStub != nil {
		return e.EnsureAuthorizedStub(ctx)
	}
	if len(e.EnsureAuthorizedOutputs) > 0 {
		output := e.EnsureAuthorizedOutputs[0]
		e.EnsureAuthorizedOutputs = e.EnsureAuthorizedOutputs[1:]
		return output
	}
	if e.EnsureAuthorizedOutput != nil {
		return *e.EnsureAuthorizedOutput
	}
	panic("EnsureAuthorized has no output")
}

func (e *ExternalAccessor) EnsureAuthorizedService(ctx context.Context) error {
	e.EnsureAuthorizedServiceInvocations++
	e.EnsureAuthorizedServiceInputs = append(e.EnsureAuthorizedServiceInputs, ctx)
	if e.EnsureAuthorizedServiceStub != nil {
		return e.EnsureAuthorizedServiceStub(ctx)
	}
	if len(e.EnsureAuthorizedServiceOutputs) > 0 {
		output := e.EnsureAuthorizedServiceOutputs[0]
		e.EnsureAuthorizedServiceOutputs = e.EnsureAuthorizedServiceOutputs[1:]
		return output
	}
	if e.EnsureAuthorizedServiceOutput != nil {
		return *e.EnsureAuthorizedServiceOutput
	}
	panic("EnsureAuthorizedService has no output")
}

func (e *ExternalAccessor) EnsureAuthorizedUser(ctx context.Context, targetUserID string, authorizedPermission string) (string, error) {
	e.EnsureAuthorizedUserInvocations++
	e.EnsureAuthorizedUserInputs = append(e.EnsureAuthorizedUserInputs, EnsureAuthorizedUserInput{Context: ctx, TargetUserID: targetUserID, AuthorizedPermission: authorizedPermission})
	if e.EnsureAuthorizedUserStub != nil {
		return e.EnsureAuthorizedUserStub(ctx, targetUserID, authorizedPermission)
	}
	if len(e.EnsureAuthorizedUserOutputs) > 0 {
		output := e.EnsureAuthorizedUserOutputs[0]
		e.EnsureAuthorizedUserOutputs = e.EnsureAuthorizedUserOutputs[1:]
		return output.AuthorizedUserID, output.Error
	}
	if e.EnsureAuthorizedUserOutput != nil {
		return e.EnsureAuthorizedUserOutput.AuthorizedUserID, e.EnsureAuthorizedUserOutput.Error
	}
	panic("EnsureAuthorizedUser has no output")
}

func (e *ExternalAccessor) AssertOutputsEmpty() {
	if len(e.ServerSessionTokenOutputs) > 0 {
		panic("ServerSessionTokenOutputs is not empty")
	}
	if len(e.ValidateSessionTokenOutputs) > 0 {
		panic("ValidateSessionTokenOutputs is not empty")
	}
	if len(e.EnsureAuthorizedOutputs) > 0 {
		panic("EnsureAuthorizedOutputs is not empty")
	}
	if len(e.EnsureAuthorizedServiceOutputs) > 0 {
		panic("EnsureAuthorizedServiceOutputs is not empty")
	}
	if len(e.EnsureAuthorizedUserOutputs) > 0 {
		panic("EnsureAuthorizedUserOutputs is not empty")
	}
}
