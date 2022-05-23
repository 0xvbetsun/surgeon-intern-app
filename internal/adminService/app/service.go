package adminService

import (
	"github.com/google/wire"
	"github.com/vbetsun/surgeon-intern-app/internal/adminService/pkg/services/clinics"
	"github.com/vbetsun/surgeon-intern-app/internal/adminService/pkg/services/encrypt"
	"github.com/vbetsun/surgeon-intern-app/internal/adminService/pkg/services/hospitals"
	"github.com/vbetsun/surgeon-intern-app/internal/adminService/pkg/services/roles"
	adminUsers "github.com/vbetsun/surgeon-intern-app/internal/adminService/pkg/services/users"
	"github.com/vbetsun/surgeon-intern-app/internal/pkg/services/authorization"
	"github.com/vbetsun/surgeon-intern-app/internal/pkg/services/users"
	"github.com/vbetsun/surgeon-intern-app/internal/pkg/thirdparty"
)

type (
	Service struct {
		AuthorizationService authorization.IService
		UserService          users.IService
		adminUserService     *adminUsers.UserService
		clinicService        *clinics.ClinicService
		hospitalService      *hospitals.HospitalService
		roleService          *roles.RoleService
	}
)

func NewService(
	authorizationService authorization.IService,
	userService users.IService,
	adminUserService *adminUsers.UserService,
	clinicService *clinics.ClinicService,
	roleService *roles.RoleService,
	hospitalService *hospitals.HospitalService) (*Service, error) {

	return &Service{
		AuthorizationService: authorizationService,
		UserService:          userService,
		adminUserService:     adminUserService,
		clinicService:        clinicService,
		hospitalService:      hospitalService,
		roleService:          roleService,
	}, nil
}

var ServiceSet = wire.NewSet(NewService,
	users.NewService,
	adminUsers.NewUserService,
	authorization.NewService,
	clinics.NewClinicService,
	hospitals.NewHospitalService,
	roles.NewRoleService,
	encrypt.NewHmac,
	thirdparty.NewAuth0ManagementApi,
	thirdparty.NewMailgun)
