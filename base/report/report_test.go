package report

import (
	"bytes"
	"fmt"
	"github.com/LightningFootball/backend/base"
	"github.com/LightningFootball/backend/base/exit"
	"github.com/LightningFootball/backend/database"
	"github.com/LightningFootball/backend/database/models"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"os"
	"testing"
	"time"
)

// 请先填入测试邮箱后再进行测试，需要修改的位置位于const常量与func TestMain中

var applyAdminUser headerOption
var applyNormalUser headerOption

type headerOption map[string][]string

//由于结构比较复杂，需要赋值的地方繁多，所以设置静态变量来理清数据库变量
const (
	//一个班级，两个管理员/老师，两个学生，一个作业（报告是对于指定作业的），两个题目，四个提交
	classID = 1

	managerID1      = 1
	managerID2      = 2
	studentID1 uint = 3
	studentID2 uint = 4

	problemSetID1 = 1

	problemID1 = 1
	problemID2 = 2

	submissionID1 = 1
	submissionID2 = 2
	submissionID3 = 3
	submissionID4 = 4

	problemName1 = "problem1"
	problemName2 = "problem2"

	// 请在此处填入测试用目标邮箱，需要修改项为：managerEmail1 managerEmail2
	managerUserName1 = "ManagerUserName1"
	managerNickName1 = "ManagerNickName1"
	managerEmail1    = "test1@test.com"
	managerUserName2 = "ManagerUserName2"
	managerNickName2 = "ManagerNickName2"
	managerEmail2    = "test2@test.com"

	studentEmail1 = "1@2.com"
	studentEmail2 = "2@2.com"
)

func TestMain(m *testing.M) {
	defer database.SetupDatabaseForTest()()
	defer exit.SetupExitForTest()
	viper.SetConfigType("yaml")
	// 请填入测试数据，需要修改项：host port username password from teacher_cron
	configFile := []byte(`debug: true
server:
  port: 8080
  origin:
    - http://localhost:8080
mail:
  host: smtp.test.com
  port: 0
  username: test@test.com
  password: test
  encryption: SSL
  from: test@test.com
  teacher_cron: "@every 1m"
`)
	err := viper.ReadConfig(bytes.NewBuffer(configFile))
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func TestReporter(t *testing.T) {
	//从数据库获取ProblemSet用于测试
	//problemSet := createProblemSet(t)
	//createClass(t)
	//Reporter(problemSet)
}

func TestTeacherReporter(t *testing.T) {
	//todo:需要测试能否在每次teacherReporter调用时都会检查是否有ProblemSet更新，及能否看出是否新加了题目
	_, managersPointerSlice := createManagers(t)
	_, studentsPointerSlice := createStudents(t)
	createClass(t, managersPointerSlice, studentsPointerSlice)
	_, problemPointerSlice := createProblems(t)
	problemSet := createProblemSet(t, problemPointerSlice)
	createSubmissions(t)
	teacherReporter(problemSet)
}

func TestMailBodyExecute(t *testing.T) {
	managers, _ := getManagers(t)
	_, problemsPointerSlice := getProblems(t)
	problemSet := getProblemSet(t, problemsPointerSlice)
	body := mailBodyExecute(managers, problemSet, getProblemSetStudentReport(t))
	fmt.Println(body)
}

/*
get函数为仅获取某个数据库模型，通过返回值返回
create函数通过调用get，创建某个数据库模型，并且通过返回值返回
*/

//TestMailBodyExecute专用
func getProblemSetStudentReport(t *testing.T) []models.ProblemSetStudentReport {
	var report = []models.ProblemSetStudentReport{
		{
			UserID:   studentID1,
			UserName: "UserName1",
			NickName: "NickName1",
			Problem: []models.SimpleProblem{
				{
					ProblemID:    problemID1,
					HighestScore: 10,
				},
				{
					ProblemID:    problemID2,
					HighestScore: 20,
				},
			},
		},
		{
			UserID:   studentID2,
			UserName: "UserName2",
			NickName: "NickName2",
			Problem: []models.SimpleProblem{
				{
					ProblemID:    problemID1,
					HighestScore: 30,
				},
				{
					ProblemID:    problemID2,
					HighestScore: 40,
				},
			},
		},
	}
	return report
}

//Problem数据库
func getProblems(t *testing.T) ([]models.Problem, []*models.Problem) {
	var problems = []models.Problem{
		{
			ID:                 problemID1,
			Name:               problemName1,
			Description:        "",
			AttachmentFileName: "",
			Public:             false,
			Privacy:            false,
			MemoryLimit:        0,
			TimeLimit:          0,
			LanguageAllowed:    nil,
			BuildArg:           "",
			CompareScriptName:  "",
			CompareScript:      models.Script{},
			TestCases:          nil,
			Tags:               nil,
			CreatedAt:          time.Time{},
			UpdatedAt:          time.Time{},
			DeletedAt:          gorm.DeletedAt{},
		},
		{
			ID:                 problemID2,
			Name:               problemName2,
			Description:        "",
			AttachmentFileName: "",
			Public:             false,
			Privacy:            false,
			MemoryLimit:        0,
			TimeLimit:          0,
			LanguageAllowed:    nil,
			BuildArg:           "",
			CompareScriptName:  "",
			CompareScript:      models.Script{},
			TestCases:          nil,
			Tags:               nil,
			CreatedAt:          time.Time{},
			UpdatedAt:          time.Time{},
			DeletedAt:          gorm.DeletedAt{},
		},
	}
	return problems, getProblemPointerSlice(problems)
}

func createProblems(t *testing.T) ([]models.Problem, []*models.Problem) {
	problems, problemsPointerSlice := getProblems(t)
	assert.NoError(t, base.DB.Create(&problems).Error)
	return problems, problemsPointerSlice
}

//ProblemSet数据库
func getProblemSet(t *testing.T, problemPointerSlice []*models.Problem) models.ProblemSet {
	problemSet := models.ProblemSet{
		ID:          problemSetID1,
		ClassID:     classID,
		Class:       nil,
		Name:        "ProblemSet Test Name",
		Description: "ProblemSet Test Description",
		Problems:    problemPointerSlice,
		Grades:      nil,
		StartTime:   time.Time{},
		EndTime:     time.Time{},
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		DeletedAt:   gorm.DeletedAt{},
	}
	return problemSet
}

func createProblemSet(t *testing.T, problemsPointerSlice []*models.Problem) models.ProblemSet {
	problemSet := getProblemSet(t, problemsPointerSlice)
	assert.NoError(t, base.DB.Create(&problemSet).Error)
	return problemSet
}

//Class数据库
func getClass(t *testing.T, managersPointerSlice []*models.User, studentsPointerSlice []*models.User) models.Class {
	class := models.Class{
		ID:          classID,
		Name:        "",
		CourseName:  "",
		Description: "",
		InviteCode:  "",
		Managers:    managersPointerSlice,
		Students:    studentsPointerSlice,
		ProblemSets: nil,
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		DeletedAt:   gorm.DeletedAt{},
	}
	return class
}

func createClass(t *testing.T, managersPointerSlice []*models.User, studentsPointerSlice []*models.User) models.Class {
	class := getClass(t, managersPointerSlice, studentsPointerSlice)
	assert.NoError(t, base.DB.Create(&class).Error)
	//assert.NoError(t, class.AddManagers([]uint{managerID1, managerID2}))
	//assert.NoError(t, class.AddStudents([]uint{studentID1, studentID2}))
	//studentsSlice := getUserFromPointerSlice(studentsPointerSlice)
	//managersSlice := getUserFromPointerSlice(managersPointerSlice)
	//assert.NoError(t, base.DB.Model(&class).Association("Managers").Append(managersSlice))
	//assert.NoError(t, base.DB.Model(&class).Association("Students").Append(studentsSlice))
	return class
}

//User数据库
func getManagers(t *testing.T) ([]models.User, []*models.User) {
	managers := []models.User{
		{
			ID:              managerID1,
			Username:        managerUserName1,
			Nickname:        managerNickName1,
			Email:           managerEmail1,
			Password:        "",
			Roles:           nil,
			RoleLoaded:      false,
			ClassesManaging: nil,
			ClassesTaking:   nil,
			Grades:          nil,
			CreatedAt:       time.Time{},
			UpdatedAt:       time.Time{},
			DeletedAt:       gorm.DeletedAt{},
			Credentials:     nil,
		},
		{
			ID:              managerID2,
			Username:        managerUserName2,
			Nickname:        managerNickName2,
			Email:           managerEmail2,
			Password:        "",
			Roles:           nil,
			RoleLoaded:      false,
			ClassesManaging: nil,
			ClassesTaking:   nil,
			Grades:          nil,
			CreatedAt:       time.Time{},
			UpdatedAt:       time.Time{},
			DeletedAt:       gorm.DeletedAt{},
			Credentials:     nil,
		},
	}
	return managers, getUserPointerSlice(managers)
}

func createManagers(t *testing.T) ([]models.User, []*models.User) {
	users, usersPointerSlice := getManagers(t)
	assert.NoError(t, base.DB.Create(&users).Error)
	return users, usersPointerSlice
}

func getStudents(t *testing.T) ([]models.User, []*models.User) {
	student := []models.User{
		{
			ID:              studentID1,
			Username:        "TestStudentUserName1",
			Nickname:        "TestStudentNickName1",
			Email:           studentEmail1,
			Password:        "",
			Roles:           nil,
			RoleLoaded:      false,
			ClassesManaging: nil,
			ClassesTaking:   nil,
			Grades:          nil,
			CreatedAt:       time.Time{},
			UpdatedAt:       time.Time{},
			DeletedAt:       gorm.DeletedAt{},
			Credentials:     nil,
		},
		{
			ID:              studentID2,
			Username:        "TestStudentUserName2",
			Nickname:        "TestStudentNickName2",
			Email:           studentEmail2,
			Password:        "",
			Roles:           nil,
			RoleLoaded:      false,
			ClassesManaging: nil,
			ClassesTaking:   nil,
			Grades:          nil,
			CreatedAt:       time.Time{},
			UpdatedAt:       time.Time{},
			DeletedAt:       gorm.DeletedAt{},
			Credentials:     nil,
		},
	}
	return student, getUserPointerSlice(student)
}

func createStudents(t *testing.T) ([]models.User, []*models.User) {
	users, usersPointerSlice := getStudents(t)
	assert.NoError(t, base.DB.Create(&users).Error)
	return users, usersPointerSlice
}

//Submission数据库
func getSubmissions(t *testing.T) []models.Submission {
	submissions := []models.Submission{
		{
			ID:           submissionID1,
			UserID:       studentID1,
			User:         nil,
			ProblemID:    problemID1,
			Problem:      nil,
			ProblemSetID: problemSetID1,
			ProblemSet:   nil,
			LanguageName: "",
			Language:     nil,
			FileName:     "",
			Priority:     0,
			Judged:       false,
			Score:        10,
			Status:       "",
			Runs:         nil,
			CreatedAt:    time.Time{},
			UpdatedAt:    time.Time{},
			DeletedAt:    gorm.DeletedAt{},
		},
		{
			ID:           submissionID2,
			UserID:       studentID1,
			User:         nil,
			ProblemID:    problemID2,
			Problem:      nil,
			ProblemSetID: problemSetID1,
			ProblemSet:   nil,
			LanguageName: "",
			Language:     nil,
			FileName:     "",
			Priority:     0,
			Judged:       false,
			Score:        20,
			Status:       "",
			Runs:         nil,
			CreatedAt:    time.Time{},
			UpdatedAt:    time.Time{},
			DeletedAt:    gorm.DeletedAt{},
		},
		{
			ID:           submissionID3,
			UserID:       studentID2,
			User:         nil,
			ProblemID:    problemID1,
			Problem:      nil,
			ProblemSetID: problemSetID1,
			ProblemSet:   nil,
			LanguageName: "",
			Language:     nil,
			FileName:     "",
			Priority:     0,
			Judged:       false,
			Score:        30,
			Status:       "",
			Runs:         nil,
			CreatedAt:    time.Time{},
			UpdatedAt:    time.Time{},
			DeletedAt:    gorm.DeletedAt{},
		},
		{
			ID:           submissionID4,
			UserID:       studentID2,
			User:         nil,
			ProblemID:    problemID2,
			Problem:      nil,
			ProblemSetID: problemSetID1,
			ProblemSet:   nil,
			LanguageName: "",
			Language:     nil,
			FileName:     "",
			Priority:     0,
			Judged:       false,
			Score:        40,
			Status:       "",
			Runs:         nil,
			CreatedAt:    time.Time{},
			UpdatedAt:    time.Time{},
			DeletedAt:    gorm.DeletedAt{},
		},
	}
	return submissions
}

func createSubmissions(t *testing.T) []models.Submission {
	submissions := getSubmissions(t)
	assert.NoError(t, base.DB.Create(&submissions).Error)
	return submissions
}
