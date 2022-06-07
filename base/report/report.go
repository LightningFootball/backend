package report

import (
	"bytes"
	"github.com/LightningFootball/backend/base"
	"github.com/LightningFootball/backend/base/utils"
	"github.com/LightningFootball/backend/database/models"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"strconv"
	"strings"
	"text/template"
)

// reporter应该具有什么功能：根据作业，向不同的人发送不同的report

const teacherMailSubject = "作业提交报告（教师版本）"
const studentMainSubject = "作业提交报告（学生版本）"
const teacherTestMailBody = "Teacher Report."
const studentMailBody = "Student Report."

func Reporter(problemSet models.ProblemSet) {
	// 根据开始时间与截止时间，设置一个cron周期性发送report，设置一个cron或者应该是timer，到截止日期，发送最终report并自动删除cron
	// 学生邮件先不开发，优先将教师邮件开发完成

	//获取老师
	//发送给老师的邮件
	//获取学生
	//发送给学生的邮件

	var teacherEmailAddress []string
	var studentEmailAddress []string

	//获取班级
	class := models.Class{}
	queryClass := base.DB.Model(&models.Class{}).Where("id = ?", problemSet.ClassID)
	err := queryClass.First(&class).Error
	if err != nil {
		err = errors.Wrap(err, "query class error")
	}

	//获取班级所有管理员，相当于所有老师
	var classManager []models.User
	queryManagers := base.DB.Model(&class).Association("Managers")
	err = queryManagers.Find(&classManager)
	if err != nil {
		err = errors.Wrap(err, "query class manager failed")
	}

	//获取所有老师的邮箱
	for _, user := range classManager {
		teacherEmailAddress = append(teacherEmailAddress, user.Email)
	}

	//todo:后续需要添加根据problemSet时间修改而更改cron的功能
	if teacherEmailAddress != nil {
		// 每日早八点，向老师发送报告
		// 创建定时任务
		teacherCronEntryID := utils.AddCron(viper.GetString("mail.teacher_cron"), func() { teacherReporter(problemSet) })
		// 到期自动删除定时任务
		utils.RemoveCronByTime(teacherCronEntryID, problemSet.EndTime)
	} else {
		err = errors.New("Unable to create cron job for send email to teacher, no teacher email found")
	}

	if studentEmailAddress != nil {
		// 每日早八点，向学生发送报告
		studentCronEntryID := utils.AddCron("00 08 * * *", func() { studentReporter(studentEmailAddress...) })
		// 到期自动删除定时任务
		utils.RemoveCronByTime(studentCronEntryID, problemSet.EndTime)
	} else {
		err = errors.New("Unable to create cron job for send email to student, no student email found")
	}
}

func teacherReporter(problemSet models.ProblemSet) {
	// 更新作业信息，以获取作业所拥有什么题目的最新情况
	var problemSetUpdate models.ProblemSet
	queryProblemSetUpdate := base.DB.Model(&models.ProblemSet{}).Preload("Problems").Where("id = ?", problemSet.ID)
	err := queryProblemSetUpdate.Find(&problemSetUpdate).Error
	if err != nil {
		err = errors.Wrap(err, "update problem set error")
	}

	// 获取班级信息，包括班级本身、教师、学生
	var managers []models.User
	var students []models.User
	var class models.Class
	queryStudents := base.DB.Model(&models.Class{}).Where("id = ?", problemSetUpdate.ClassID).Preload("Students").Preload("Managers")
	err = queryStudents.Find(&class).Error
	if err != nil {
		err = errors.Wrap(err, "query student failed")
	}
	managers = getUserFromPointerSlice(class.Managers)
	students = getUserFromPointerSlice(class.Students)

	//todo:其实可以不一次把这个作业所有的题目的提交全部获取，可以先获取其中一道题目的，然后填充，然后再获取其它题目的提交，再填充
	//todo:这样可以避免下面填充提交分数的时候多层for嵌套，不过如果改掉report相关结构体，从struct改成map，应该一次性获取然后key/value模式填充也没啥问题
	//todo:不过多次访问数据库应该比在这逻辑判断啥的慢，且会加重数据库负端，不过这数据库也没啥负担。。

	// 获取该作业的所有提交
	var submissions []models.Submission
	querySubmissions := base.DB.Model(&models.Submission{}).Where("problem_set_id = ?", problemSet.ID).Order("user_id asc")
	err = querySubmissions.Find(&submissions).Error
	if err != nil {
		err = errors.Wrap(err, "query submission failed")
	}

	var reports []models.ProblemSetStudentReport
	// 填充学生信息
	for _, student := range students {
		reports = append(reports, models.ProblemSetStudentReport{
			UserID:   student.ID,
			UserName: student.Username,
			NickName: student.Nickname,
			Problem:  nil,
		})
	}

	//todo:for range在此处不能用，会导致reports.problem的值append不上
	//todo:https://blog.csdn.net/wYc037/article/details/108824049，详情看这篇帖子
	//todo:简述就是for中的每次实例，比如如下的report，都是那固定的一个，只不过每次传值给他修改它指定的值，相当于实例其实是指针
	//for _, problem := range problemSetUpdate.Problems {
	//	for _, report := range reports {
	//		sp := models.SimpleProblem{
	//			ProblemID:    problem.ID,
	//			HighestScore: 0,
	//		}
	//		report.Problem = append(report.Problem, sp)
	//	}
	//}

	// 计算分析信息
	for i := 0; i < len(problemSetUpdate.Problems); i++ {
		for n := 0; n < len(reports); n++ {
			sp := models.SimpleProblem{
				ProblemID:    problemSetUpdate.Problems[i].ID,
				HighestScore: 0,
			}
			reports[n].Problem = append(reports[n].Problem, sp)
		}
	}
	// 遍历所有提交
	for s := 0; s < len(submissions); s++ {
	C:
		for r := 0; r < len(reports); r++ {
			if reports[r].UserID == submissions[s].UserID {
				problems := reports[r].Problem
				for p := 0; p < len(problems); p++ {
					if problems[p].ProblemID == submissions[s].ProblemID {
						problems[p].HighestScore = submissions[s].Score
						break C
					}
				}
			}
		}
	}
	teacherMailBody := mailBodyExecute(managers, problemSetUpdate, reports)
	for i := 0; i < len(teacherMailBody); i++ {
		if reports != nil {
			SendMail(teacherMailSubject, teacherMailBody[i], managers[i].Email)
		} else {
			panic("mail body is empty, did not send")
		}
	}
}

func studentReporter(mail ...string) {
	SendMail(studentMainSubject, studentMailBody, mail...)
}

func mailBodyExecute(user []models.User, problemSet models.ProblemSet, report []models.ProblemSetStudentReport) []string {
	/*
		邮件数据部分所含内容：将所有同学的作业信息按表格打印出来
	*/
	var data bytes.Buffer
	for _, r := range report {
		data.WriteString("用户名称: ")
		data.WriteString(r.UserName)
		data.WriteString("<br/>用户昵称:  ")
		data.WriteString(r.NickName)
		data.WriteString("<br/>")
		for _, p := range r.Problem {
			// 些许丑陋，后续看看能不能改进此处.
			data.WriteString("&nbsp;&nbsp;问题ID: ")
			data.WriteString(strconv.Itoa(int(p.ProblemID)))
			data.WriteString("<br>&nbsp;&nbsp;当前最高分: ")
			data.WriteString(strconv.Itoa(int(p.HighestScore)))
			data.WriteString("<br/>")
		}
		data.WriteString("<br/>")
	}

	// 找出未完成作业的学生
	var unfinishedUser []string
	for i := 0; i < len(report); i++ {
		for k := 0; k < len(report[i].Problem); k++ {
			if report[i].Problem[k].HighestScore != 100 {
				unfinishedUser = append(unfinishedUser, report[i].UserName)
				break
			}
		}
	}
	unfinishedUserString := strings.Join(unfinishedUser, ",")

	// 将数据填充到静态文件中
	var bodies []bytes.Buffer
	for i := 0; i < len(user); i++ {
		// 解析静态文件
		t, err := template.ParseFiles("resource/template/cron_teacher_mail.html")
		//单元测试获取不到这个地址，所以需要做一下判断
		//如果是单元测试，那么此时的路径就是xxx_test.go文件所在的路径，所以要根据相对路径调整下
		if t == nil {
			t, err = template.ParseFiles("../../resource/template/cron_teacher_mail.html")
		}
		if err != nil {
			err = errors.Wrap(err, "can not parse files")
		}

		// 填充数据到静态文件中
		var body bytes.Buffer
		err = t.Execute(&body, struct {
			/*注意此处需要大写，大写才能从html导出，卡在这里很久很久*/
			UnfinishedUser string
			NickName       string
			ProblemSetName string
			Data           string
		}{
			UnfinishedUser: unfinishedUserString,
			NickName:       user[i].Nickname,
			ProblemSetName: problemSet.Name,
			Data:           data.String(),
		})
		if err != nil {
			err = errors.Wrap(err, "can not execute template")
		}
		bodies = append(bodies, body)
	}
	bodiesString := make([]string, len(bodies))
	for i := 0; i < len(bodies); i++ {
		bodiesString[i] = bodies[i].String()
	}
	return bodiesString
}

// PointerSlice 转换

func getUserPointerSlice(userSlice []models.User) []*models.User {
	var userPointerSlice []*models.User
	//此处一样不能用for range
	//for _, user := range userSlice {
	//	userPointerSlice = append(userPointerSlice, &user)
	//}
	for i := 0; i < len(userSlice); i++ {
		userPointerSlice = append(userPointerSlice, &userSlice[i])
	}
	return userPointerSlice
}

func getProblemPointerSlice(problemSlice []models.Problem) []*models.Problem {
	var problemPointerSlice []*models.Problem
	for i := 0; i < len(problemSlice); i++ {
		problemPointerSlice = append(problemPointerSlice, &problemSlice[i])
	}
	return problemPointerSlice
}

func getUserFromPointerSlice(userPointerSlice []*models.User) []models.User {
	var userSlice []models.User
	for i := 0; i < len(userPointerSlice); i++ {
		userSlice = append(userSlice, *userPointerSlice[i])
	}
	return userSlice
}
