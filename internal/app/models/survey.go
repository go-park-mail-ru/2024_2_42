package models

type SurveyResponse struct {
	SurveyID  uint64      `json:"survey_id"`
	Title     string      `json:"survey_title"`
	Questions []*Question `json:"survey_questions"`
}

type Survey struct {
	SurveyID uint64 `json:"survey_id"`
	Title    string `json:"survey_title"`
}

type Question struct {
	QuestionID uint64 `json:"question_id"`
	Content    string `json:"question_content"`
}

type Mark struct {
	UserID     uint64 `json:"user_id"`
	SurveyID   uint64 `json:"survey_id"`
	QuestionID uint64 `json:"question_id"`
	Score      int    `json:"score"`
}
