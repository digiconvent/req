package req

type Stage string

const (
	StageBuildRequest Stage = "build request"
	StageSendRequest  Stage = "send request"
	StageReadBody     Stage = "read body"
	StageStatus       Stage = "status"
)
