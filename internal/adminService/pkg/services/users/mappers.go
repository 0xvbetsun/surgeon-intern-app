package users

import (
	qlmodel "github.com/vbetsun/surgeon-intern-app/internal/adminService/graph/model"
	commonModel "github.com/vbetsun/surgeon-intern-app/internal/graph/model"
	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/users"
)

func (u *UserService) FromUsersWithUnitRoles(usersWithUnitRoles []*users.UserWithUnitRoles) ([]*qlmodel.DetailedUser, error) {
	detailedUsers := make([]*qlmodel.DetailedUser, 0)
	for _, userWithUnitRoles := range usersWithUnitRoles {
		detailedUser, err := u.FromUserWithUnitRoles(userWithUnitRoles)
		if err != nil {
			continue
		}
		detailedUsers = append(detailedUsers, detailedUser)
	}

	return detailedUsers, nil
}

func (u *UserService) FromUserWithUnitRoles(userWithUnitRoles *users.UserWithUnitRoles) (*qlmodel.DetailedUser, error) {
	detailedUser := &qlmodel.DetailedUser{}
	if userWithUnitRoles.User.UserExternalID.IsZero() {
		detailedUser.UserID = userWithUnitRoles.User.ID
		detailedUser.Email = userWithUnitRoles.User.Email
		detailedUser.Activated = userWithUnitRoles.User.Activated
		detailedUser.ClinicRoles = nil
	} else {
		auth0User, err := u.auth0.Management.User.Read(userWithUnitRoles.User.UserExternalID.String)
		if err != nil {
			return nil, err
		}
		detailedUser.UserID = userWithUnitRoles.User.ID
		detailedUser.Email = auth0User.GetEmail()
		detailedUser.Username = auth0User.GetUsername()
		detailedUser.FamilyName = auth0User.GetFamilyName()
		detailedUser.GivenName = auth0User.GetGivenName()
		detailedUser.ClinicRoles = nil
		detailedUser.Activated = userWithUnitRoles.User.Activated
	}

	clinicRoles := make([]*commonModel.ClinicRole, 0)
	for _, unitRole := range userWithUnitRoles.UnitRoles {
		clinicRole := &commonModel.ClinicRole{
			Clinic: &commonModel.Clinic{
				ClinicID:    unitRole.Unit.ID,
				DisplayName: unitRole.Unit.DisplayName,
				HospitalID:  unitRole.Unit.ParentID.Ptr(),
			},
			Role: &commonModel.Role{
				RoleIdentifier: unitRole.Role.Name,
				DisplayName:    unitRole.Role.DisplayName,
			},
		}
		if unitRole.Unit.R.Parent != nil {
			clinicRole.Clinic.Hospital = &commonModel.Hospital{
				HospitalID:  unitRole.Unit.ParentID.String,
				DisplayName: unitRole.Unit.R.Parent.DisplayName,
			}
		}
		clinicRoles = append(clinicRoles, clinicRole)
	}
	detailedUser.ClinicRoles = clinicRoles
	return detailedUser, nil
}

func (u *UserService) FromUser(user *models.User) (*qlmodel.DetailedUser, error) {
	detailedUser := &qlmodel.DetailedUser{}
	if user.UserExternalID.IsZero() {
		detailedUser.UserID = user.ID
		detailedUser.Email = user.Email
		detailedUser.Activated = user.Activated
		detailedUser.ClinicRoles = nil
	} else {
		auth0User, err := u.auth0.Management.User.Read(user.UserExternalID.String)
		if err != nil {
			return nil, err
		}
		detailedUser.UserID = user.ID
		detailedUser.Email = auth0User.GetEmail()
		detailedUser.Username = auth0User.GetUsername()
		detailedUser.FamilyName = auth0User.GetFamilyName()
		detailedUser.GivenName = auth0User.GetGivenName()
		detailedUser.ClinicRoles = nil
		detailedUser.Activated = user.Activated
	}

	clinicRoles := make([]*commonModel.ClinicRole, 0)
	if user.R != nil && user.R.UserOrganizationalUnitRoles != nil {
		for _, uor := range user.R.UserOrganizationalUnitRoles {
			clinicRole := &commonModel.ClinicRole{
				Clinic: &commonModel.Clinic{
					ClinicID:    uor.R.Unit.ID,
					DisplayName: uor.R.Unit.DisplayName,
					HospitalID:  uor.R.Unit.ParentID.Ptr(),
				},
				Role: &commonModel.Role{
					RoleIdentifier: uor.R.Role.Name,
					DisplayName:    uor.R.Role.DisplayName,
				},
			}
			if uor.R.Unit.R != nil && uor.R.Unit.R.Parent != nil {
				clinicRole.Clinic.Hospital = &commonModel.Hospital{
					HospitalID:  uor.R.Unit.ParentID.String,
					DisplayName: uor.R.Unit.R.Parent.DisplayName,
				}
			}
			clinicRoles = append(clinicRoles, clinicRole)
		}
	}
	detailedUser.ClinicRoles = clinicRoles
	return detailedUser, nil
}
