package service

import (
	"github.com/Surafeljava/Court-Case-Management-System/caseUse"
	"github.com/Surafeljava/Court-Case-Management-System/entity"
)

type CaseServiceImpl struct {
	caseRepo caseUse.CaseRepository
}

func NewCaseServiceImpl(csRepo caseUse.CaseRepository) *CaseServiceImpl {
	return &CaseServiceImpl{caseRepo: csRepo}
}

//get all the cases from the database and return
func (csi *CaseServiceImpl) Cases() ([]entity.Case, error) {
	cases, errs := csi.caseRepo.Cases()
	return cases, errs
}

//get a single case
func (csi *CaseServiceImpl) Case(id int) (*entity.Case, []error) {
	css, _ := csi.caseRepo.Case(id)
	// if errs != nil {
	// 	return css, errs
	// }
	return css, nil
}

//add new case to the database
func (csi *CaseServiceImpl) CreateCase(casedoc *entity.Case) []error {
	err1 := csi.caseRepo.CreateCase(casedoc)
	if len(err1) > 0 {
		panic(err1)
	}
	return err1
}

//UpdateCase >> update a case
func (csi *CaseServiceImpl) UpdateCase(casedoc *entity.Case) (*entity.Case, []error) {
	cs, errs := csi.caseRepo.UpdateCase(casedoc)
	if len(errs) > 0 {
		return nil, errs
	}
	return cs, errs
}

//CloseCase Sevice >> close existing case
func (csi *CaseServiceImpl) CloseCase(casedoc entity.Case) error {
	return nil
}

//ExtendCase Service >> extend the case court date of existing case
func (csi *CaseServiceImpl) ExtendCase(casedoc entity.Case) error {
	return nil
}

//DeleteCase Service >> delete existing case
func (csi *CaseServiceImpl) DeleteCase(id int) error {
	return nil
}
