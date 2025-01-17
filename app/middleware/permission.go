package middleware

import (
	"github.com/LightningFootball/backend/app/response"
	"github.com/LightningFootball/backend/base"
	"github.com/LightningFootball/backend/base/log"
	"github.com/LightningFootball/backend/base/utils"
	"github.com/LightningFootball/backend/database/models"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
)

type hasRole struct {
	ID   uint
	Name string
}

func (h *hasRole) GetID() uint {
	return h.ID
}
func (h *hasRole) TypeName() string {
	return h.Name
}

type PermissionOption interface {
	Check(ctx echo.Context) bool
}

type UnscopedPermission struct {
	P string
}

func (p UnscopedPermission) Check(c echo.Context) bool {
	u := c.Get("user").(models.User)
	return u.Can(p.P)
}

type ScopedPermission struct {
	P           string
	IdFieldName string
	T           string
}

func (p ScopedPermission) Check(c echo.Context) bool {
	idFieldName := "id"
	if p.IdFieldName != "" {
		idFieldName = p.IdFieldName
	}
	u := c.Get("user").(models.User)
	id, err := strconv.ParseUint(c.Param(idFieldName), 10, strconv.IntSize)
	if err != nil {
		return false
	}
	return u.Can(p.P, &hasRole{
		ID:   uint(id),
		Name: p.T,
	})
}

type OrPermission struct {
	A PermissionOption
	B PermissionOption
}

func (p OrPermission) Check(c echo.Context) bool {
	return p.A.Check(c) || p.B.Check(c)
}

type AndPermission struct {
	A PermissionOption
	B PermissionOption
}

func (p AndPermission) Check(c echo.Context) bool {
	return p.A.Check(c) && p.B.Check(c)
}

type CustomPermission struct {
	F func(c echo.Context) bool
}

func (p CustomPermission) Check(c echo.Context) bool {
	return p.F(c)
}

func HasPermission(p PermissionOption) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if p.Check(c) {
				return next(c)
			}
			return c.JSON(http.StatusForbidden, response.ErrorResp("PERMISSION_DENIED", nil))
		}
	}
}

func IsTestCaseSample(c echo.Context) (result bool) {
	user := c.Get("user").(models.User)
	testCase, problem, err := utils.FindTestCase(c.Param("id"), c.Param("test_case_id"), &user)
	if err == nil {
		result = testCase.Sample
	}
	c.Set("test_case", testCase)
	c.Set("problem", problem)
	c.Set("find_test_case_err", err)
	return
}

func IsTestCaseSampleProblemSet(c echo.Context) (result bool) {
	testCase, problem, err := utils.FindTestCase(c.Param("id"), c.Param("test_case_id"), nil)
	if err == nil {
		result = testCase.Sample
	}
	c.Set("test_case", testCase)
	c.Set("problem", problem)
	c.Set("find_test_case_err", err)
	return
}

func ProblemSetStarted(c echo.Context) (result bool) {
	problemSet := models.ProblemSet{}
	err := base.DB.First(&problemSet, "class_id = ? and id = ?", c.Param("class_id"), c.Param("problem_set_id")).Error
	if c.Param("problem_set_id") == "" {
		log.Infof("%+v\n%s\n", c.Request(), c.Param("problem_set_id"))
	}
	c.Set("problem_set", &problemSet)
	c.Set("find_problem_set_error", err)
	if err == nil {
		return time.Now().After(problemSet.StartTime)
	} else {
		log.Info(c.Request())
		return true
	}
}
