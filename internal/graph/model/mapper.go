package commonModel

import (
	"github.com/thoas/go-funk"
	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/users"
)

func RolesFromRepoTypes(roles []*models.Role) []*Role {
	var qlRoles = make([]*Role, 0)
	for _, role := range roles {
		qlRoles = append(qlRoles, RoleFromRepoType(role))
	}
	return qlRoles
}

func RoleFromRepoType(rm *models.Role) *Role {
	return &Role{
		DisplayName:    rm.DisplayName,
		RoleIdentifier: rm.Name,
	}
}

func ClinicFromRepoType(ou *models.OrganizationalUnit) *Clinic {
	return &Clinic{
		ClinicID:    ou.ID,
		DisplayName: ou.DisplayName,
		HospitalID:  ou.ParentID.Ptr(),
	}
}

func HospitalFromRepoType(ou *models.OrganizationalUnit) *Hospital {
	return &Hospital{
		HospitalID:  ou.ID,
		DisplayName: ou.DisplayName,
		Clinics:     nil,
	}
}

func PracticalActivityTypesFromRepoType(pt []*models.PracticalActivityType) []*PracticalActivityType {
	var practicalActivityTypes = make([]*PracticalActivityType, 0)
	for _, activity := range pt {
		practicalActivityTypes = append(practicalActivityTypes, PracticalActivityTypeFromRepoType(activity))
	}
	return practicalActivityTypes
}

func PracticalActivityTypeFromRepoType(pt *models.PracticalActivityType) *PracticalActivityType {
	return &PracticalActivityType{
		PracticalActivityTypeID: pt.ID,
		DisplayName:             pt.DisplayName,
		Name:                    LogbookEntryType(pt.Name),
	}
}

func UsersFromUserUnitRoles(userUnitRoles users.UserUnitRoles) []*User {
	uniqueUsers := make([]*models.User, 0)
	userIds := make([]string, 0)
	for _, userUnitRole := range userUnitRoles {
		if !funk.ContainsString(userIds, userUnitRole.User.ID) {
			userIds = append(userIds, userUnitRole.User.ID)
			uniqueUsers = append(uniqueUsers, &userUnitRole.User)
		}
	}
	qlUsers := make([]*User, 0)
	for _, user := range uniqueUsers {
		clinicRoles := make([]*ClinicRole, 0)
		clinicIds := make([]string, 0)
		for _, userUnitRole := range userUnitRoles {
			if userUnitRole.User.ID == user.ID {
				clinicRoles = append(clinicRoles, ClinicRoleFromRepoTypes(&userUnitRole.OrganizationalUnit, &userUnitRole.Role))
				clinicIds = append(clinicIds, userUnitRole.OrganizationalUnit.ID)
			}
		}
		qlUser := UserFromRepoType(user)
		qlUser.ClinicRoles = clinicRoles
		qlUser.ClinicIds = clinicIds
		qlUsers = append(qlUsers, qlUser)
	}
	return qlUsers
}

func UserFromRepoType(u *models.User) *User {
	return &User{
		UserID:      u.ID,
		DisplayName: u.DisplayName,
		ClinicRoles: nil,
		ClinicIds:   nil,
	}
}

func ClinicRoleFromRepoTypes(ou *models.OrganizationalUnit, r *models.Role) *ClinicRole {
	return &ClinicRole{
		Clinic: &Clinic{
			ClinicID:    ou.ID,
			DisplayName: ou.DisplayName,
			HospitalID:  ou.ParentID.Ptr(),
			Hospital:    nil,
		},
		Role: &Role{
			RoleIdentifier: r.Name,
			DisplayName:    r.DisplayName,
		},
	}
}
