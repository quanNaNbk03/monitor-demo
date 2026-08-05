package repository

import (
	"context"
	"errors"
	"fmt"

	"git.ocn.com.vn/ocn/common/httpbase/ierror"
	"git.ocn.com.vn/ocn/common/ormbase"
	"github.com/quanNaNbk03/monitor-demo/internal/server/model"
	"github.com/quanNaNbk03/monitor-demo/pkg/common/errlist"
	"gorm.io/gorm"
)

type OrgRepository interface {
	CreateOrg(ctx context.Context, req *model.Organization, opts ...ormbase.Option) *ierror.CoreError
	GetOrgByID(ctx context.Context, id uint64, opts ...ormbase.Option) (*model.Organization, *ierror.CoreError)
	GetOrgByName(ctx context.Context, name string) (*model.Organization, *ierror.CoreError)
	ListOrgs(ctx context.Context, param model.ListOrgParams, opts ...ormbase.Option) ([]*model.Organization, int64, *ierror.CoreError)
	DeleteOrgByID(ctx context.Context, id uint64, opts ...ormbase.Option) *ierror.CoreError
	UpdateOrgByID(ctx context.Context, id uint64, organization *model.Organization, updatedFields []string, opts ...ormbase.Option) *ierror.CoreError
}

type orgRepository struct {
	db *gorm.DB
}

func NewOrgRepository(db *gorm.DB) OrgRepository {
	return &orgRepository{
		db: db,
	}
}

func (r *orgRepository) CreateOrg(ctx context.Context, req *model.Organization, opts ...ormbase.Option) *ierror.CoreError {
	tx := ormbase.NewConfig(opts...).ToGormTx(r.db).WithContext(ctx)
	if err := tx.Create(req).Error; err != nil {
		if ormbase.IsMySQLDuplicateKey(err) {
			return errlist.ErrOrgAlreadyExists.WithChild(err)
		}
		return errlist.ErrDatabase.WithChild(err)
	}

	return nil
}

func (r *orgRepository) GetOrgByID(ctx context.Context, id uint64, opts ...ormbase.Option) (*model.Organization, *ierror.CoreError) {
	tx := ormbase.NewConfig(opts...).ToGormTx(r.db).WithContext(ctx)
	var rs model.Organization
	if err := tx.Take(&rs, "id=?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errlist.ErrOrgNotFound
		}
		return nil, errlist.ErrDatabase.WithChild(err)
	}
	return &rs, nil
}

func (r *orgRepository) GetOrgByName(ctx context.Context, name string) (*model.Organization, *ierror.CoreError) {
	db := r.db.WithContext(ctx)
	var rs model.Organization
	if err := db.Where("name = ?", name).First(&rs).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errlist.ErrOrgNotFound
		}
		return nil, errlist.ErrDatabase.WithChild(err)
	}
	return &rs, nil
}

func (r *orgRepository) DeleteOrgByID(ctx context.Context, id uint64, opts ...ormbase.Option) *ierror.CoreError {
	tx := ormbase.NewConfig(opts...).ToGormTx(r.db).WithContext(ctx)

	if err := tx.Delete(&model.Organization{}, "id=?", id).Error; err != nil {
		return errlist.ErrDatabase.WithChild(fmt.Errorf("failed to delete organization %w", err))
	}

	return nil
}

func (r *orgRepository) UpdateOrgByID(ctx context.Context, id uint64, organization *model.Organization, updatedFields []string, opts ...ormbase.Option) *ierror.CoreError {
	tx := ormbase.NewConfig(opts...).ToGormTx(r.db).WithContext(ctx)

	if len(updatedFields) > 0 {
		tx = tx.Select(updatedFields)
	}
	if err := tx.Model(&model.Organization{}).Where("id = ?", id).Updates(organization).Error; err != nil {
		return errlist.ErrDatabase.WithChild(err)
	}
	return nil
}

func (r *orgRepository) ListOrgs(ctx context.Context, param model.ListOrgParams, opts ...ormbase.Option) ([]*model.Organization, int64, *ierror.CoreError) {
	var res []*model.Organization
	db := r.db.WithContext(ctx)
	db = ormbase.ApplyFilters(db, model.Organization{}.AvailableSearchFields(), param.Searches)

	var total int64
	if err := db.Model(&model.Organization{}).Count(&total).Error; err != nil {
		return nil, 0, errlist.ErrDatabase.WithChild(err)
	}

	ormCfg := ormbase.NewConfig(opts...)
	if ormCfg.Tx == nil {
		db = ormCfg.ToGormTx(db)
	} else {
		db = ormCfg.Tx.WithContext(ctx)
		db = ormbase.ApplyFilters(db, model.Organization{}.AvailableSearchFields(), param.Searches)
	}

	if err := db.Find(&res).Error; err != nil {
		return nil, 0, errlist.ErrDatabase.WithChild(err)
	}
	return res, total, nil
}
