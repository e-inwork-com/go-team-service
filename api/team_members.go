package api

import (
	"errors"
	"net/http"

	"github.com/e-inwork-com/go-team-service/internal/data"
	"github.com/e-inwork-com/go-team-service/internal/validator"
	"github.com/google/uuid"
)

func (app *Application) createTeamMemberHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TeamMemberTeam uuid.UUID `json:"team_member_team"`
		TeamMemberUser uuid.UUID `json:"team_member_user"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Check team exist
	team, err := app.Models.Teams.GetByID(input.TeamMemberTeam)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Check user exist
	_, err = app.Models.Users.GetByID(input.TeamMemberUser)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Only team's owner can add a member user
	user := app.contextGetUser(r)
	if team.TeamUser != user.ID {
		app.notPermittedResponse(w, r)
		return
	}

	teamMember := &data.TeamMember{
		TeamMemberTeam: input.TeamMemberTeam,
		TeamMemberUser: input.TeamMemberUser,
	}

	v := validator.New()
	if data.ValidateTeamMember(v, teamMember); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.Models.TeamMembers.Insert(teamMember)
	if err != nil {
		switch {
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Get the Team Member just created
	teamMember, err = app.Models.TeamMembers.GetByID(teamMember.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"team_member": teamMember}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) deleteTeamMemberHandler(w http.ResponseWriter, r *http.Request) {
	// Get ID from the request parameters
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Get Team Member from the database
	teamMember, err := app.Models.TeamMembers.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Get a Team from the database
	team, err := app.Models.Teams.GetByID(teamMember.TeamMemberTeam)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Get the current user
	user := app.contextGetUser(r)

	// Check if the record has a related to the current user
	if team.TeamUser != user.ID {
		app.notPermittedResponse(w, r)
		return
	}

	// Delete Team Member
	err = app.Models.TeamMembers.Delete(teamMember)
	if err != nil {
		switch {
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Send a request response
	err = app.writeJSON(w, http.StatusOK, nil, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) listTeamMembersByOwnerHandler(w http.ResponseWriter, r *http.Request) {
	// Get the current user
	user := app.contextGetUser(r)

	// Get a Team from the database
	team, err := app.Models.Teams.GetByTeamUser(user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Get list
	teamMembers, err := app.Models.TeamMembers.ListByOwner(team.ID)
	if err != nil {
		switch {
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Response
	err = app.writeJSON(w, http.StatusOK, envelope{"team_members": teamMembers}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) getTeamMemberHandler(w http.ResponseWriter, r *http.Request) {
	// Get ID from the request parameters
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Get Team Member from the database
	teamMember, err := app.Models.TeamMembers.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Get a Team
	team, err := app.Models.Teams.GetByID(teamMember.TeamMemberTeam)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Get the current user
	user := app.contextGetUser(r)

	// Check if the record has a related to the current user
	if teamMember.TeamMemberUser != user.ID && team.TeamUser != user.ID {
		app.notPermittedResponse(w, r)
		return
	}

	// Send a request response
	err = app.writeJSON(w, http.StatusOK, envelope{"team_member": teamMember}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
