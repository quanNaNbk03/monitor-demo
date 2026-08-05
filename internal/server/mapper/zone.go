package mapper

import (
	"git.ocn.com.vn/ocn/common/ormbase"
	"github.com/quanNaNbk03/monitor-demo/internal/server/dto"
	"github.com/quanNaNbk03/monitor-demo/internal/server/model"
	"github.com/quanNaNbk03/monitor-demo/pkg/api/zone"
)

func ToZoneDTOs(zs []*model.Zone) []zone.Zone {
	dtos := make([]zone.Zone, 0, len(zs))
	for _, n := range zs {
		res := ToZoneDTO(n)
		if res != nil {
			dtos = append(dtos, *res)
		} else {
			dtos = append(dtos, zone.Zone{})
		}
	}
	return dtos
}

func ToZoneDTO(z *model.Zone) *zone.Zone {
	return &zone.Zone{
		Id:             z.ID,
		OrganizationID: z.OrganizationID,
		DnsSec:         z.DNSSEC,
		Name:           z.Name,
		Type:           z.Type.String(),
		SoaType:        z.SOAType.String(),
		CreatedAt:      z.UpdatedAt,
		UpdatedAt:      z.CreatedAt,
	}
}

func ToListZonesParams(req *dto.ListZonesJSONBody) *model.ListZoneParams {
	searches := make([]*ormbase.Search, 0, len(req.Searches))
	for _, search := range req.Searches {
		searches = append(searches, &ormbase.Search{
			Field: search.Field,
			Op:    search.Op,
			Value: search.Value,
		})
	}
	return &model.ListZoneParams{
		Paging:   ormbase.NewPaging(req.Paging.Page, req.Paging.Limit),
		Searches: searches,
	}
}

func ToListZonesResponse(zones []*model.Zone, pagingDTO dto.Paging) *dto.ListZoneResponse {
	return &dto.ListZoneResponse{
		Organizations: ToZoneDTOs(zones),
		Paging:        pagingDTO,
	}
}
