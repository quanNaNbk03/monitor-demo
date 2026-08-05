package mapper

import (
	"git.ocn.com.vn/ocn/common/ormbase"
	"github.com/quanNaNbk03/monitor-demo/internal/server/dto"
	"github.com/quanNaNbk03/monitor-demo/internal/server/model"
	org "github.com/quanNaNbk03/monitor-demo/pkg/api/organization"
)

func ToOrgDTO(o *model.Organization) *org.Organization {
	return &org.Organization{
		Id:          o.ID,
		Name:        o.Name,
		Description: o.Description,
		Email:       o.Email,
		Contact:     o.ContactName,
		UpdatedAt:   o.UpdatedAt,
		CreatedAt:   o.CreatedAt,
	}
}

func ToListOrgsParams(req *org.ListOrgsJSONBody) *model.ListOrgParams {
	searches := make([]*ormbase.Search, 0, len(req.Searches))
	for _, search := range req.Searches {
		searches = append(searches, &ormbase.Search{
			Field: search.Field,
			Op:    search.Op,
			Value: search.Value,
		})
	}
	return &model.ListOrgParams{
		Paging:   ormbase.NewPaging(req.Paging.Page, req.Paging.Limit),
		Searches: searches,
	}
}

func ToListOrgResponse(orgs []*model.Organization, pagingDTO dto.Paging) *org.ListOrganizationResponse {
	return &org.ListOrganizationResponse{
		Organizations: ToOrgsDTO(orgs),
		Paging:        pagingDTO,
	}
}

func ToOrgsDTO(orgs []*model.Organization) []org.Organization {
	dtos := make([]org.Organization, 0, len(orgs))
	for _, n := range orgs {
		res := ToOrgDTO(n)
		if res != nil {
			dtos = append(dtos, *res)
		} else {
			dtos = append(dtos, org.Organization{})
		}
	}
	return dtos
}
