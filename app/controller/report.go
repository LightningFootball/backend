package controller

import "github.com/LightningFootball/backend/database/models"

func Reporter(problemSet models.ProblemSet) {
	//根据开始时间与截止时间，设置一个cron周期性发送report，设置一个cron或者应该是timer，到截止日期，发送最终report并自动删除cron
}
