// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/mattermost/mattermost/server/v8/boards/model"
	"github.com/mattermost/mattermost/server/v8/boards/services/audit"

	"github.com/mattermost/mattermost/server/public/shared/mlog"
)

func (a *API) registerBoardsAndBlocksRoutes(r *mux.Router) {
	// BoardsAndBlocks APIs
	r.HandleFunc("/boards-and-blocks", a.sessionRequired(a.handleCreateBoardsAndBlocks)).Methods("POST")
	r.HandleFunc("/boards-and-blocks", a.sessionRequired(a.handlePatchBoardsAndBlocks)).Methods("PATCH")
	r.HandleFunc("/boards-and-blocks", a.sessionRequired(a.handleDeleteBoardsAndBlocks)).Methods("DELETE")
}

func (a *API) handleCreateBoardsAndBlocks(w http.ResponseWriter, r *http.Request) {
	// swagger:operation POST /boards-and-blocks insertBoardsAndBlocks
	//
	// Creates new boards and blocks
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: Body
	//   in: body
	//   description: the boards and blocks to create
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/BoardsAndBlocks"
	// security:
	// - BearerAuth: []
	// responses:
	//   '200':
	//     description: success
	//     schema:
	//       $ref: '#/definitions/BoardsAndBlocks'
	//   default:
	//     description: internal error
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"

	userID := getUserID(r)

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		a.errorResponse(w, r, err)
		return
	}

	var newBab *model.BoardsAndBlocks
	if err = json.Unmarshal(requestBody, &newBab); err != nil {
		a.errorResponse(w, r, err)
		return
	}

	if len(newBab.Boards) == 0 {
		a.errorResponse(w, r, model.NewErrBadRequest("at least one board is required"))
		return
	}

	teamID := ""
	boardIDs := map[string]bool{}
	for _, board := range newBab.Boards {
		boardIDs[board.ID] = true

		if teamID == "" {
			teamID = board.TeamID
			continue
		}

		if teamID != board.TeamID {
			a.errorResponse(w, r, model.NewErrBadRequest("cannot create boards for multiple teams"))
			return
		}

		if board.ID == "" {
			a.errorResponse(w, r, model.NewErrBadRequest("boards need an ID to be referenced from the blocks"))
			return
		}
	}

	if !a.permissions.HasPermissionToTeam(userID, teamID, model.PermissionViewTeam) {
		a.errorResponse(w, r, model.NewErrPermission("access denied to board template"))
		return
	}

	isGuest, err := a.userIsGuest(userID)
	if err != nil {
		a.errorResponse(w, r, err)
		return
	}
	if isGuest {
		a.errorResponse(w, r, model.NewErrPermission("access denied to create board"))
		return
	}

	for _, block := range newBab.Blocks {
		// Error checking
		if len(block.Type) < 1 {
			message := fmt.Sprintf("missing type for block id %s", block.ID)
			a.errorResponse(w, r, model.NewErrBadRequest(message))
			return
		}

		if block.CreateAt < 1 {
			message := fmt.Sprintf("invalid createAt for block id %s", block.ID)
			a.errorResponse(w, r, model.NewErrBadRequest(message))
			return
		}

		if block.UpdateAt < 1 {
			message := fmt.Sprintf("invalid UpdateAt for block id %s", block.ID)
			a.errorResponse(w, r, model.NewErrBadRequest(message))
			return
		}

		if !boardIDs[block.BoardID] {
			message := fmt.Sprintf("invalid BoardID %s (not exists in the created boards)", block.BoardID)
			a.errorResponse(w, r, model.NewErrBadRequest(message))
			return
		}
	}

	// IDs of boards and blocks are used to confirm that they're
	// linked and then regenerated by the server
	newBab, err = model.GenerateBoardsAndBlocksIDs(newBab, a.logger)
	if err != nil {
		a.errorResponse(w, r, model.NewErrBadRequest(err.Error()))
		return
	}

	auditRec := a.makeAuditRecord(r, "createBoardsAndBlocks", audit.Fail)
	defer a.audit.LogRecord(audit.LevelModify, auditRec)
	auditRec.AddMeta("teamID", teamID)
	auditRec.AddMeta("userID", userID)
	auditRec.AddMeta("boardsCount", len(newBab.Boards))
	auditRec.AddMeta("blocksCount", len(newBab.Blocks))

	// create boards and blocks
	bab, err := a.app.CreateBoardsAndBlocks(newBab, userID, true)
	if err != nil {
		a.errorResponse(w, r, err)
		return
	}

	a.logger.Debug("CreateBoardsAndBlocks",
		mlog.String("teamID", teamID),
		mlog.String("userID", userID),
		mlog.Int("boardCount", len(bab.Boards)),
		mlog.Int("blockCount", len(bab.Blocks)),
	)

	data, err := json.Marshal(bab)
	if err != nil {
		a.errorResponse(w, r, err)
		return
	}

	// response
	jsonBytesResponse(w, http.StatusOK, data)

	auditRec.Success()
}

func (a *API) handlePatchBoardsAndBlocks(w http.ResponseWriter, r *http.Request) {
	// swagger:operation PATCH /boards-and-blocks patchBoardsAndBlocks
	//
	// Patches a set of related boards and blocks
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: Body
	//   in: body
	//   description: the patches for the boards and blocks
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/PatchBoardsAndBlocks"
	// security:
	// - BearerAuth: []
	// responses:
	//   '200':
	//     description: success
	//     schema:
	//       $ref: '#/definitions/BoardsAndBlocks'
	//   default:
	//     description: internal error
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"

	userID := getUserID(r)

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		a.errorResponse(w, r, err)
		return
	}

	var pbab *model.PatchBoardsAndBlocks
	if err = json.Unmarshal(requestBody, &pbab); err != nil {
		a.errorResponse(w, r, err)
		return
	}

	if err = pbab.IsValid(); err != nil {
		a.errorResponse(w, r, model.NewErrBadRequest(err.Error()))
		return
	}

	teamID := ""
	boardIDMap := map[string]bool{}
	for i, boardID := range pbab.BoardIDs {
		boardIDMap[boardID] = true
		patch := pbab.BoardPatches[i]

		if err = patch.IsValid(); err != nil {
			a.errorResponse(w, r, model.NewErrBadRequest(err.Error()))
			return
		}

		if !a.permissions.HasPermissionToBoard(userID, boardID, model.PermissionManageBoardProperties) {
			a.errorResponse(w, r, model.NewErrPermission("access denied to modifying board properties"))
			return
		}

		if patch.Type != nil || patch.MinimumRole != nil {
			if !a.permissions.HasPermissionToBoard(userID, boardID, model.PermissionManageBoardType) {
				a.errorResponse(w, r, model.NewErrPermission("access denied to modifying board type"))
				return
			}
		}

		board, err2 := a.app.GetBoard(boardID)
		if err2 != nil {
			a.errorResponse(w, r, err2)
			return
		}

		if teamID == "" {
			teamID = board.TeamID
		}
		if teamID != board.TeamID {
			a.errorResponse(w, r, model.NewErrBadRequest("mismatched team ID"))
			return
		}
	}

	for _, blockID := range pbab.BlockIDs {
		block, err2 := a.app.GetBlockByID(blockID)
		if err2 != nil {
			a.errorResponse(w, r, err2)
			return
		}

		if _, ok := boardIDMap[block.BoardID]; !ok {
			a.errorResponse(w, r, model.NewErrBadRequest("missing BoardID="+block.BoardID))
			return
		}

		if !a.permissions.HasPermissionToBoard(userID, block.BoardID, model.PermissionManageBoardCards) {
			a.errorResponse(w, r, model.NewErrPermission("access denied to modifying cards"))
			return
		}
	}

	auditRec := a.makeAuditRecord(r, "patchBoardsAndBlocks", audit.Fail)
	defer a.audit.LogRecord(audit.LevelModify, auditRec)
	auditRec.AddMeta("boardsCount", len(pbab.BoardIDs))
	auditRec.AddMeta("blocksCount", len(pbab.BlockIDs))

	bab, err := a.app.PatchBoardsAndBlocks(pbab, userID)
	if err != nil {
		a.errorResponse(w, r, err)
		return
	}

	a.logger.Debug("PATCH BoardsAndBlocks",
		mlog.Int("boardsCount", len(pbab.BoardIDs)),
		mlog.Int("blocksCount", len(pbab.BlockIDs)),
	)

	data, err := json.Marshal(bab)
	if err != nil {
		a.errorResponse(w, r, err)
		return
	}

	// response
	jsonBytesResponse(w, http.StatusOK, data)

	auditRec.Success()
}

func (a *API) handleDeleteBoardsAndBlocks(w http.ResponseWriter, r *http.Request) {
	// swagger:operation DELETE /boards-and-blocks deleteBoardsAndBlocks
	//
	// Deletes boards and blocks
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: Body
	//   in: body
	//   description: the boards and blocks to delete
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/DeleteBoardsAndBlocks"
	// security:
	// - BearerAuth: []
	// responses:
	//   '200':
	//     description: success
	//   default:
	//     description: internal error
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"

	userID := getUserID(r)

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		a.errorResponse(w, r, err)
		return
	}

	var dbab *model.DeleteBoardsAndBlocks
	if err = json.Unmarshal(requestBody, &dbab); err != nil {
		a.errorResponse(w, r, model.NewErrBadRequest(err.Error()))
		return
	}

	// user must have permission to delete all the boards, and that
	// would include the permission to manage their blocks
	teamID := ""
	boardIDMap := map[string]bool{}
	for _, boardID := range dbab.Boards {
		boardIDMap[boardID] = true
		// all boards in the request should belong to the same team
		board, err := a.app.GetBoard(boardID)
		if err != nil {
			a.errorResponse(w, r, err)
			return
		}
		if teamID == "" {
			teamID = board.TeamID
		}
		if teamID != board.TeamID {
			a.errorResponse(w, r, model.NewErrBadRequest("all boards should belong to the same team"))
			return
		}

		// permission check
		if !a.permissions.HasPermissionToBoard(userID, boardID, model.PermissionDeleteBoard) {
			a.errorResponse(w, r, model.NewErrPermission("access denied to delete board"))
			return
		}
	}

	for _, blockID := range dbab.Blocks {
		block, err2 := a.app.GetBlockByID(blockID)
		if err2 != nil {
			a.errorResponse(w, r, err2)
			return
		}

		if _, ok := boardIDMap[block.BoardID]; !ok {
			a.errorResponse(w, r, model.NewErrBadRequest("missing BoardID="+block.BoardID))
			return
		}

		if !a.permissions.HasPermissionToBoard(userID, block.BoardID, model.PermissionManageBoardCards) {
			a.errorResponse(w, r, model.NewErrPermission("access denied to modifying cards"))
			return
		}
	}

	if err := dbab.IsValid(); err != nil {
		a.errorResponse(w, r, model.NewErrBadRequest(err.Error()))
		return
	}

	auditRec := a.makeAuditRecord(r, "deleteBoardsAndBlocks", audit.Fail)
	defer a.audit.LogRecord(audit.LevelModify, auditRec)
	auditRec.AddMeta("boardsCount", len(dbab.Boards))
	auditRec.AddMeta("blocksCount", len(dbab.Blocks))

	if err := a.app.DeleteBoardsAndBlocks(dbab, userID); err != nil {
		a.errorResponse(w, r, err)
		return
	}

	a.logger.Debug("DELETE BoardsAndBlocks",
		mlog.Int("boardsCount", len(dbab.Boards)),
		mlog.Int("blocksCount", len(dbab.Blocks)),
	)

	// response
	jsonStringResponse(w, http.StatusOK, "{}")
	auditRec.Success()
}
