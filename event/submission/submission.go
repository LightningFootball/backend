package submission

import "github.com/LightningFootball/backend/database/models"

// EventArgs is the arguments of "submission" event.
type EventArgs = *models.Submission

// EventRst is the result of "submission" event.
type EventRst error
