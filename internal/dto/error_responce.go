package dto

type ErrorResp struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

const (
	CodeTeamExists      = "TEAM_EXISTS"
	CodePRExists        = "PR_EXISTS"
	CodePRMerged        = "PR_MERGED"
	CodeNotAssigned     = "NOT_ASSIGNED"
	CodeNoCandidate     = "NO_CANDIDATE"
	CodeNotFound        = "NOT_FOUND"
	CodeValidationError = "VALIDATION_ERROR"
	CodeInternalError   = "INTERNAL_ERROR"
)

const (
	MsgTeamExists        = "team_name already exists"
	MsgPRExists          = "PR id already exists"
	MsgPRMerged          = "cannot reassign on merged PR"
	MsgNotFound          = "resource not found"
	MsgValidationError   = "invalid request body"
	MsgInvalidQueryError = "invalid request query"
	MsgInternalError     = "internal error"
)
