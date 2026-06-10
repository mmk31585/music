package creator

type OverviewResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type DailyStatsQuery struct {
	From  string `form:"from"`
	To    string `form:"to"`
	Limit int    `form:"limit,default=30"`
}
