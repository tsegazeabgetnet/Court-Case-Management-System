package caseUse

import "github.com/Surafeljava/Court-Case-Management-System/entity"

type CaseRepository interface {
	Cases() ([]entity.Case, error)
	JudgeCases(juid string) ([]entity.Case, error)
	CaseJudges(case_type string) ([]entity.Judge, error)
	Case(id int) (*entity.Case, []error)
	CreateCase(casedoc *entity.Case) []error
	UpdateCase(casedoc *entity.Case) (*entity.Case, []error)
	CloseCase(casedoc string, decision *entity.Decision) []error
	ExtendCase(casedoc *entity.Case) []error
	DeleteCase(id int) []error
}

type OpponentRepository interface {
	Opponents() ([]entity.Opponent, error)
	Opponent(id int) (*entity.Opponent, []error)
	CreateOpponent(opp *entity.Opponent) (*entity.Opponent, []error)
}

type JudgeRepository interface {
	Judges() ([]entity.Judge, error)
	Judge(id int) (*entity.Judge, []error)
	CreateJudge(judge *entity.Judge) (*entity.Judge, []error)
	UpdateCase(judge *entity.Judge) (*entity.Judge, []error)
	DeleteCase(id int) error
}

type LoginRepository interface {
	CheckLogin(user *entity.UserType) (*entity.UserType, []error)
	CheckAdmin(id string, pwd string) (*entity.Admin, []error)
	CheckJudge(id string, pwd string) (*entity.Judge, []error)
	CheckOpponent(id string, pwd string) (*entity.Opponent, []error)
}

type CaseSearchRepository interface {
	Cases() ([]entity.Case, []error)
	Case(id uint) (*entity.Case, []error)
}
