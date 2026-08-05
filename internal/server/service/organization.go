package service

import (
	"context"
	"errors"

	"git.ocn.com.vn/ocn/common/httpbase"
	"git.ocn.com.vn/ocn/common/httpbase/ierror"
	"git.ocn.com.vn/ocn/common/ormbase"
	"github.com/quanNaNbk03/monitor-demo/internal/server/mapper"
	"github.com/quanNaNbk03/monitor-demo/internal/server/model"
	"github.com/quanNaNbk03/monitor-demo/internal/server/repository"
	org "github.com/quanNaNbk03/monitor-demo/pkg/api/organization"
	"github.com/quanNaNbk03/monitor-demo/pkg/common/errlist"
)

type OrgService interface {
	Create(ctx context.Context, in *org.CreateOrgRequest) *ierror.Error
	Get(ctx context.Context, id uint64) (*org.Organization, *ierror.Error)
	List(ctx context.Context, in *org.ListOrgsJSONBody) (*org.ListOrganizationResponse, *ierror.Error)
	Delete(ctx context.Context, id uint64) *ierror.Error
	Update(ctx context.Context, id uint64, in *org.UpdateOrgRequest) *ierror.Error
}

type orgService struct {
	orgRepo repository.OrgRepository
}

func NewOrgService(orgRepo repository.OrgRepository) OrgService {
	return &orgService{
		orgRepo: orgRepo,
	}
}

func (s *orgService) Create(ctx context.Context, in *org.CreateOrgRequest) *ierror.Error {
	exist, coreErr := s.orgRepo.GetOrgByName(ctx, in.Name)
	if coreErr != nil && !errors.Is(coreErr, errlist.ErrOrgNotFound) {
		return httpbase.ErrInternal(ctx, "get org by name error").SetSubError(coreErr)
	}
	if exist != nil {
		return httpbase.ErrBadRequest(ctx, "org already exists")
	}
	req := &model.Organization{
		ContactName: in.Contact,
		Description: in.Description,
		Name:        in.Name,
		Email:       in.Email,
	}
	if coreErr := s.orgRepo.CreateOrg(ctx, req); coreErr != nil {
		if errors.Is(coreErr, errlist.ErrOrgAlreadyExists) {
			return httpbase.ErrBadRequest(ctx, "org already exists")
		}
		return httpbase.ErrInternal(ctx, "create org error").SetSubError(coreErr)
	}
	return nil
}

func (s *orgService) Get(ctx context.Context, id uint64) (*org.Organization, *ierror.Error) {
	res, coreErr := s.orgRepo.GetOrgByID(ctx, id)
	if coreErr != nil {
		if errors.Is(coreErr, errlist.ErrOrgNotFound) {
			return nil, httpbase.ErrNotFound(ctx, "org not found")
		}
		return nil, httpbase.ErrInternal(ctx, "get org error").SetSubError(coreErr)
	}
	return mapper.ToOrgDTO(res), nil
}

func (s *orgService) Delete(ctx context.Context, id uint64) *ierror.Error {
	if ierr := s.orgRepo.DeleteOrgByID(ctx, id); ierr != nil {
		return httpbase.ErrInternal(ctx, "delete org error").SetSubError(ierr)
	}
	return nil
}

func (s *orgService) Update(ctx context.Context, id uint64, in *org.UpdateOrgRequest) *ierror.Error {
	exist, coreErr := s.orgRepo.GetOrgByID(ctx, id)
	if coreErr != nil && !errors.Is(coreErr, errlist.ErrOrgNotFound) {
		return httpbase.ErrInternal(ctx, "get org by id error").SetSubError(coreErr)
	}
	updateFields := make([]string, 0, 3)
	if in.Contact != nil {
		exist.ContactName = *in.Contact
		updateFields = append(updateFields, "contact_name")
	}
	if in.Description != nil {
		exist.Description = *in.Description
		updateFields = append(updateFields, "description")
	}
	if in.Email != nil {
		exist.Email = *in.Email
		updateFields = append(updateFields, "email")
	}
	if err := s.orgRepo.UpdateOrgByID(ctx, id, exist, updateFields); err != nil {
		return httpbase.ErrInternal(ctx, "update org error").SetSubError(err)
	}
	return nil
}

func (s *orgService) List(ctx context.Context, in *org.ListOrgsJSONBody) (*org.ListOrganizationResponse, *ierror.Error) {
	param := mapper.ToListOrgsParams(in)
	orgs, total, coreErr := s.orgRepo.ListOrgs(ctx, *param, ormbase.WithPage(param.Paging))
	if coreErr != nil {
		return nil, httpbase.ErrInternal(ctx, "list orgs failed").SetSubError(coreErr)
	}
	pagingDTO := mapper.ToPagingDTO(param.Paging.GetPage(), param.Paging.GetLimit(), int(total))
	return mapper.ToListOrgResponse(orgs, pagingDTO), nil
}
