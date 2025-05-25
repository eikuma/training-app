package form

type ProposeTrainingMenu struct {
	TrainingGoal    string   `json:"training_goal" form:"training_goal"`       // 例: "筋肥大", "ダイエット", ...
	TargetParts     []string `json:"target_parts" form:"target_parts"`         // 例: ["胸", "背中", "脚"] ...
	ExperienceLevel string   `json:"experience_level" form:"experience_level"` // 例: "初心者", "中級者", "上級者"
	AvailableTime   int      `json:"available_time" form:"available_time"`
}

func NewProposeTrainingMenu() *ProposeTrainingMenu {
	return &ProposeTrainingMenu{}
}
