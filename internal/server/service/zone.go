package service

import (
	"context"
	"errors"

	"git.ocn.com.vn/ocn/common/httpbase"
	"git.ocn.com.vn/ocn/common/httpbase/ierror"
	"git.ocn.com.vn/ocn/common/ormbase"
	"github.com/quanNaNbk03/monitor-demo/internal/server/dto"
	"github.com/quanNaNbk03/monitor-demo/internal/server/mapper"
	"github.com/quanNaNbk03/monitor-demo/internal/server/model"
	"github.com/quanNaNbk03/monitor-demo/internal/server/repository"
	"github.com/quanNaNbk03/monitor-demo/pkg/api/zone"
	"github.com/quanNaNbk03/monitor-demo/pkg/common/errlist"
)

type ZoneService interface {
	Get(ctx context.Context, id uint64) (*zone.Zone, *ierror.Error)
	Create(ctx context.Context, in *zone.CreateZoneRequest) *ierror.Error
	List(ctx context.Context, in *dto.ListZonesJSONBody) (*dto.ListZoneResponse, *ierror.Error)
	Delete(ctx context.Context, id uint64) *ierror.Error
	Update(ctx context.Context, id uint64, in *dto.UpdateZoneRequest) *ierror.Error
}

type zoneService struct {
	zoneRepo repository.ZoneRepository
	orgRepo  repository.OrgRepository
}

func NewZoneService(zoneRepo repository.ZoneRepository, orgRepo repository.OrgRepository) ZoneService {
	return &zoneService{
		zoneRepo: zoneRepo,
		orgRepo:  orgRepo,
	}
}

const (
	msgGetZoneFailed = "failed to get zone by id"
)

func (s *zoneService) Get(ctx context.Context, id uint64) (*zone.Zone, *ierror.Error) {
	re, coreErr := s.zoneRepo.GetZoneByID(ctx, id, false)
	if coreErr != nil {
		if errors.Is(coreErr, errlist.ErrZoneNotFound) {
			return nil, httpbase.ErrNotFound(ctx, "zone not found")
		}
		return nil, httpbase.ErrInternal(ctx, msgGetZoneFailed).SetSubError(coreErr)
	}
	if re.IsTemplate {
		return nil, httpbase.ErrBadRequest(ctx, "this is zone template")
	}
	return mapper.ToZoneDTO(re), nil
}

func (s *zoneService) List(ctx context.Context, in *dto.ListZonesJSONBody) (*dto.ListZoneResponse, *ierror.Error) {
	param := mapper.ToListZonesParams(in)
	param.Searches = append(param.Searches, &ormbase.Search{
		Field: "is_template",
		Op:    "eq",
		Value: false,
	})
	zones, total, coreErr := s.zoneRepo.ListZones(ctx, *param, ormbase.WithPage(param.Paging))
	if coreErr != nil {
		return nil, httpbase.ErrInternal(ctx, "list zones failed").SetSubError(coreErr)
	}
	pagingDTO := mapper.ToPagingDTO(param.Paging.GetPage(), param.Paging.GetLimit(), int(total))
	return mapper.ToListZonesResponse(zones, pagingDTO), nil
}

func (s *zoneService) Create(ctx context.Context, in *zone.CreateZoneRequest) *ierror.Error {
	exist, coreErr := s.zoneRepo.GetZoneByName(ctx, in.Name, false)
	if coreErr != nil && !errors.Is(coreErr, errlist.ErrZoneNotFound) {
		return httpbase.ErrInternal(ctx, msgGetZoneFailed).SetSubError(coreErr)
	}
	if exist != nil {
		return httpbase.ErrBadRequest(ctx, "zone already exists")
	}
	existOrg, coreErr := s.orgRepo.GetOrgByID(ctx, in.OrganizationID)
	if coreErr != nil {
		return httpbase.ErrInternal(ctx, "failed to get org by id").SetSubError(coreErr)
	}
	if existOrg == nil {
		return httpbase.ErrBadRequest(ctx, "organization does not exist")
	}
	req := &model.Zone{
		Name:           in.Name,
		Type:           model.ParseType(in.Type),
		OrganizationID: in.OrganizationID,
		SOAType:        model.ParseSOAType(in.SoaType),
		IsTemplate:     false,
	}
	if coreErr := s.zoneRepo.CreateZone(ctx, req); coreErr != nil {
		if errors.Is(coreErr, errlist.ErrZoneAlreadyExists) {
			return httpbase.ErrBadRequest(ctx, "zone already exists")
		}
		return httpbase.ErrInternal(ctx, "failed to create zone").SetSubError(coreErr)
	}
	return nil
}

func (s *zoneService) Update(ctx context.Context, id uint64, in *dto.UpdateZoneRequest) *ierror.Error {
	exist, coreErr := s.zoneRepo.GetZoneByID(ctx, id, false)
	if coreErr != nil {
		return httpbase.ErrInternal(ctx, msgGetZoneFailed).SetSubError(coreErr)
	}

	updateFields := make([]string, 0, 4)
	if in.DnsSec != nil {
		exist.DNSSEC = *in.DnsSec
		updateFields = append(updateFields, "dnssec")
	}
	if in.SoaType != nil {
		exist.SOAType = model.ParseSOAType(*in.SoaType)
		updateFields = append(updateFields, "soa_type")
	}
	if in.Type != nil {
		exist.Type = model.ParseType(*in.Type)
		updateFields = append(updateFields, "type")
	}
	if err := s.zoneRepo.UpdateZoneByID(ctx, id, exist, updateFields); err != nil {
		return httpbase.ErrInternal(ctx, "update zone error").SetSubError(err)
	}
	return nil
}

func (s *zoneService) Delete(ctx context.Context, id uint64) *ierror.Error {
	if ierr := s.zoneRepo.DeleteZoneByID(ctx, id); ierr != nil {
		return httpbase.ErrInternal(ctx, "delete zone error").SetSubError(ierr)
	}
	return nil
}
