package submission

import (
	"github.com/LightningFootball/backend/base/utils"
	"github.com/pkg/errors"
)

func UpdateGrade(r EventArgs) EventRst {
	err := utils.UpdateGrade(r)
	return errors.Wrap(err, "could not update grade")
}
