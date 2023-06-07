package test

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"kayak-backend/api"
	"kayak-backend/model"
	"kayak-backend/utils"
	"math/rand"
	"strconv"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))
var randomUsers []model.User
var randomGroups []model.Group
var randomProblems []model.ProblemType
var randomProblemSets []model.ProblemSet
var randomNotes []model.Note
var randomNoteReviews []model.NoteReview
var randomDiscussions []model.Discussion

func randomInitUser(tx *sqlx.Tx) error {
	sqlString := `INSERT INTO "user" (name, created_at, password, nick_name, email) VALUES ($1, now(), $2, $3, $4)`
	for i := 0; i < 5; i++ {
		user := model.User{
			ID:       i + 1,
			Name:     fmt.Sprintf("user%d", i),
			NickName: fmt.Sprintf("user%d", i),
			Email:    fmt.Sprintf("user%d@boat4study.com", i),
			Password: fmt.Sprintf("%d-pwd", i),
		}
		randomUsers = append(randomUsers, user)
		encryptPassword, _ := utils.EncryptPassword(user.Password)
		if _, err := tx.Exec(sqlString, user.Name, encryptPassword, user.NickName, user.Email); err != nil {
			return err
		}
	}
	return nil
}

func randomInitGroup(tx *sqlx.Tx) error {
	sqlString := `INSERT INTO "group" (name, description, invitation, created_at, user_id, area_id) VALUES ($1, $2, $3, now(), $4, $5)`
	for i := 0; i < 3; i++ {
		user := randomUser()
		group := model.Group{
			Id:          i + 1,
			Name:        fmt.Sprintf("group%d", i),
			Description: fmt.Sprintf("group%d", i),
			Invitation:  "0",
			UserId:      user.ID,
			AreaId:      r.Intn(20) + 1,
		}
		randomGroups = append(randomGroups, group)
		if _, err := tx.Exec(sqlString, group.Name, group.Description, group.Invitation, group.UserId, group.AreaId); err != nil {
			return err
		}
		_sqlString := `INSERT INTO group_member (group_id, user_id, created_at, is_admin, is_owner) VALUES ($1, $2, now(), $3, $4)`
		if _, err := tx.Exec(_sqlString, group.Id, group.UserId, true, true); err != nil {
			return err
		}
		for j := 0; j < 5; j++ {
			_sqlString = `INSERT INTO group_member (group_id, user_id, created_at, is_admin, is_owner) VALUES ($1, $2, now(), $3, $4) ON CONFLICT DO NOTHING`
			user = randomUser()
			if _, err := tx.Exec(_sqlString, group.Id, user.ID, false, false); err != nil {
				return err
			}
		}
	}
	return nil
}

func randomInitNote(tx *sqlx.Tx) error {
	sqlString := `INSERT INTO note (title, content, created_at, updated_at, user_id, is_public) VALUES ($1, $2, now(), now(), $3, $4)`
	for i := 0; i < 10; i++ {
		user := randomUser()
		note := model.Note{
			ID:       i + 1,
			Title:    fmt.Sprintf("note%d", i),
			Content:  fmt.Sprintf("note%d", i),
			UserId:   user.ID,
			IsPublic: r.Intn(2) == 0,
		}
		randomNotes = append(randomNotes, note)
		if _, err := tx.Exec(sqlString, note.Title, note.Content, note.UserId, note.IsPublic); err != nil {
			return err
		}
	}
	return nil
}

func randomInitNoteReview(tx *sqlx.Tx) error {
	sqlString := `INSERT INTO note_review (title, content, created_at, updated_at, user_id, note_id) VALUES ($1, $2, now(), now(), $3, $4)`
	for i := 0; i < 20; i++ {
		user := randomUser()
		note := randomNote()
		noteReview := model.NoteReview{
			ID:      i + 1,
			Title:   fmt.Sprintf("noteReview%d", i),
			Content: fmt.Sprintf("noteReview%d", i),
			UserId:  user.ID,
			NoteId:  note.ID,
		}
		randomNoteReviews = append(randomNoteReviews, noteReview)
		if _, err := tx.Exec(sqlString, noteReview.Title, noteReview.Content, noteReview.UserId, noteReview.NoteId); err != nil {
			return err
		}
	}
	return nil
}

func randomInitProblemSet(tx *sqlx.Tx) error {
	sqlString := `INSERT INTO problem_set (name, description, created_at, updated_at, 
        user_id, is_public, area_id, group_id) VALUES ($1, $2, now(), now(), $3, $4, $5, $6)`
	for i := 0; i < 5; i++ {
		user := randomUser()
		problemSet := model.ProblemSet{
			ID:          i + 1,
			Name:        fmt.Sprintf("problemSet%d", i),
			Description: fmt.Sprintf("problemSet%d", i),
			UserId:      user.ID,
			IsPublic:    r.Intn(2) == 0,
			AreaId:      r.Intn(20) + 1,
			GroupId:     0,
		}
		randomProblemSets = append(randomProblemSets, problemSet)
		if _, err := tx.Exec(sqlString, problemSet.Name, problemSet.Description, problemSet.UserId,
			problemSet.IsPublic, problemSet.AreaId, problemSet.GroupId); err != nil {
			return err
		}
	}
	for i := 0; i < 3; i++ {
		group := randomGroup()
		problemSet := model.ProblemSet{
			ID:          i + 6,
			Name:        fmt.Sprintf("problemSet%d", i),
			Description: fmt.Sprintf("problemSet%d", i),
			UserId:      group.UserId,
			IsPublic:    r.Intn(2) == 0,
			AreaId:      r.Intn(20) + 1,
			GroupId:     group.Id,
		}
		randomProblemSets = append(randomProblemSets, problemSet)
		if _, err := tx.Exec(sqlString, problemSet.Name, problemSet.Description, problemSet.UserId,
			problemSet.IsPublic, problemSet.AreaId, problemSet.GroupId); err != nil {
			return err
		}
	}
	return nil
}

func randomInitProblem(tx *sqlx.Tx) error {
	sqlString := `INSERT INTO problem_type (description, created_at, updated_at, user_id, problem_type_id, is_public) VALUES ($1, now(), now(), $2, $3, $4)`
	for i := 0; i < 20; i++ {
		user := randomUser()
		problemType := model.ProblemType{
			ID:            i + 1,
			Description:   fmt.Sprintf("problem%d", i),
			UserId:        user.ID,
			ProblemTypeId: r.Intn(3),
			IsPublic:      r.Intn(2) == 0,
		}
		randomProblems = append(randomProblems, problemType)
		if _, err := tx.Exec(sqlString, problemType.Description, problemType.UserId,
			problemType.ProblemTypeId, problemType.IsPublic); err != nil {
			return err
		}
		if problemType.ProblemTypeId == api.ChoiceProblemType {
			_sqlString := `INSERT INTO problem_choice (id, choice, description, is_correct) VALUES ($1, $2, $3, $4)`
			for j := 0; j < 4; j++ {
				problemChoice := model.ProblemChoice{
					ID:          i + 1,
					Choice:      strconv.Itoa('A' + j),
					Description: fmt.Sprintf("problem%dchoice%d", i, j),
					IsCorrect:   r.Intn(2) == 0,
				}
				if _, err := tx.Exec(_sqlString, problemChoice.ID, problemChoice.Choice,
					problemChoice.Description, problemChoice.IsCorrect); err != nil {
					return err
				}
			}
		} else if problemType.ProblemTypeId == api.JudgeProblemType {
			_sqlString := `INSERT INTO problem_judge (id, is_correct) VALUES ($1, $2)`
			problemJudge := model.ProblemJudge{
				ID:        i + 1,
				IsCorrect: r.Intn(2) == 0,
			}
			if _, err := tx.Exec(_sqlString, problemJudge.ID, problemJudge.IsCorrect); err != nil {
				return err
			}
		} else if problemType.ProblemTypeId == api.BlankProblemType {
			_sqlString := `INSERT INTO problem_answer (id, answer) VALUES ($1, $2)`
			problemAnswer := model.ProblemAnswer{
				ID:     i + 1,
				Answer: fmt.Sprintf("problem%danswer", i),
			}
			if _, err := tx.Exec(_sqlString, problemAnswer.ID, problemAnswer.Answer); err != nil {
				return err
			}
		}
		_sqlString := `INSERT INTO problem_in_problem_set (problem_set_id, problem_id) VALUES ($1, $2)`
		problemSet := randomProblemSet()
		if _, err := tx.Exec(_sqlString, problemSet.ID, problemType.ID); err != nil {
			return err
		}
	}
	return nil
}

func randomInitDiscussion(tx *sqlx.Tx) error {
	sqlString := `INSERT INTO discussion (title, content, created_at, updated_at, user_id, group_id, is_public)
		VALUES ($1, $2, now(), now(), $3, $4, $5)`
	for i := 0; i < 5; i++ {
		group := randomGroup()
		discussion := model.Discussion{
			ID:       i + 1,
			Title:    fmt.Sprintf("discussion%d", i),
			Content:  fmt.Sprintf("discussion%d", i),
			UserId:   group.UserId,
			GroupId:  group.Id,
			IsPublic: r.Intn(2) == 0,
		}
		randomDiscussions = append(randomDiscussions, discussion)
		if _, err := tx.Exec(sqlString, discussion.Title, discussion.Content, discussion.UserId,
			discussion.GroupId, discussion.IsPublic); err != nil {
			return err
		}
	}
	return nil
}

func randomInitDiscussionReview(tx *sqlx.Tx) error {
	sqlString := `INSERT INTO discussion_review (title, content, created_at, updated_at, user_id, discussion_id)
		VALUES ($1, $2, now(), now(), $3, $4)`
	for i := 0; i < 10; i++ {
		discussion := randomDiscussion()
		discussionReview := model.DiscussionReview{
			ID:           i + 1,
			Title:        fmt.Sprintf("discussionReview%d", i),
			Content:      fmt.Sprintf("discussionReview%d", i),
			UserId:       discussion.UserId,
			DiscussionId: discussion.ID,
		}
		if _, err := tx.Exec(sqlString, discussionReview.Title, discussionReview.Content, discussionReview.UserId,
			discussionReview.DiscussionId); err != nil {
			return err
		}
	}
	return nil
}
