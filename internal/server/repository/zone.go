package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"git.ocn.com.vn/ocn/common/httpbase/ierror"
	"git.ocn.com.vn/ocn/common/ormbase"
	"github.com/quanNaNbk03/monitor-demo/internal/server/model"
	"github.com/quanNaNbk03/monitor-demo/pkg/common/errlist"
)

type ZoneRepository interface {
	GetZoneByID(ctx context.Context, id uint64, isTemplate bool, opts ...ormbase.Option) (*model.Zone, *ierror.CoreError)
	CreateZone(ctx context.Context, req *model.Zone, opts ...ormbase.Option) *ierror.CoreError
	ListZones(ctx context.Context, param model.ListZoneParams, opts ...ormbase.Option) ([]*model.Zone, int64, *ierror.CoreError)
	DeleteZoneByID(ctx context.Context, id uint64, opts ...ormbase.Option) *ierror.CoreError
	UpdateZoneByID(ctx context.Context, id uint64, zone *model.Zone, updateFields []string, opts ...ormbase.Option) *ierror.CoreError
	GetZoneByName(ctx context.Context, name string, isTemplate bool) (*model.Zone, *ierror.CoreError)
}

type zoneRepository struct {
	db *gorm.DB
}

func NewZoneRepository(db *gorm.DB) ZoneRepository {
	return &zoneRepository{
		db: db,
	}
}

func (r *zoneRepository) GetZoneByID(ctx context.Context, id uint64, isTemplate bool, opts ...ormbase.Option) (*model.Zone, *ierror.CoreError) {
	tx := ormbase.NewConfig(opts...).ToGormTx(r.db).WithContext(ctx)
	var rs model.Zone
	if err := tx.Take(&rs, "id=?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errlist.ErrZoneNotFound
		}
		return nil, errlist.ErrDatabase.WithChild(err)
	}
	return &rs, nil
}

func (r *zoneRepository) GetZoneByName(ctx context.Context, name string, isTemplate bool) (*model.Zone, *ierror.CoreError) {
	db := r.db.WithContext(ctx)
	var rs model.Zone
	if err := db.Where("name = ?", name).First(&rs).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errlist.ErrZoneNotFound
		}
		return nil, errlist.ErrDatabase.WithChild(err)
	}
	return &rs, nil
}

func (r *zoneRepository) CreateZone(ctx context.Context, req *model.Zone, opts ...ormbase.Option) *ierror.CoreError {
	tx := ormbase.NewConfig(opts...).ToGormTx(r.db).WithContext(ctx)

	if err := tx.Create(req).Error; err != nil {
		if ormbase.IsMySQLDuplicateKey(err) {
			return errlist.ErrZoneAlreadyExists.WithChild(err)
		}
		return errlist.ErrDatabase.WithChild(err)
	}

	return nil
}

func (r *zoneRepository) UpdateZoneByID(ctx context.Context, id uint64, zone *model.Zone, updateFields []string, opts ...ormbase.Option) *ierror.CoreError {
	tx := ormbase.NewConfig(opts...).ToGormTx(r.db).WithContext(ctx)

	if len(updateFields) > 0 {
		tx = tx.Select(updateFields)
	}
	if err := tx.Model(&model.Zone{}).Where("id = ?", id).Updates(zone).Error; err != nil {
		return errlist.ErrDatabase.WithChild(err)
	}
	return nil
}

func (r *zoneRepository) ListZones(ctx context.Context, param model.ListZoneParams, opts ...ormbase.Option) ([]*model.Zone, int64, *ierror.CoreError) {
	var res []*model.Zone
	db := r.db.WithContext(ctx)
	db = ormbase.ApplyFilters(db, model.Zone{}.AvailableSearchFields(), param.Searches)

	var total int64
	if err := db.Model(&model.Zone{}).Count(&total).Error; err != nil {
		return nil, 0, errlist.ErrDatabase.WithChild(err)
	}

	ormCfg := ormbase.NewConfig(opts...)
	if ormCfg.Tx == nil {
		db = ormCfg.ToGormTx(db)
	} else {
		db = ormCfg.Tx.WithContext(ctx)
		db = ormbase.ApplyFilters(db, model.Zone{}.AvailableSearchFields(), param.Searches)
	}

	if err := db.Find(&res).Error; err != nil {
		return nil, 0, errlist.ErrDatabase.WithChild(err)
	}
	return res, total, nil
}

func (r *zoneRepository) DeleteZoneByID(ctx context.Context, id uint64, opts ...ormbase.Option) *ierror.CoreError {
	tx := ormbase.NewConfig(opts...).ToGormTx(r.db).WithContext(ctx)

	if err := tx.Delete(&model.Zone{}, "id=?", id).Error; err != nil {
		return errlist.ErrDatabase.WithChild(fmt.Errorf("failed to delete zone %w", err))
	}

	return nil
}
