package errlist

import "git.ocn.com.vn/ocn/common/httpbase/ierror"

var (
	ErrDatabase          = ierror.NewCoreError("err_database", "")
	ErrOrgAlreadyExists  = ierror.NewCoreError("err_organization_already_exists", "")
	ErrZoneAlreadyExists = ierror.NewCoreError("err_zone_already_exists", "")
	ErrZoneNotFound      = ierror.NewCoreError("err_zone_not_found", "")
	ErrOrgNotFound       = ierror.NewCoreError("err_org_not_found", "")
)
