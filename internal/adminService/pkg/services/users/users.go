package users

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	qlmodel "github.com/vbetsun/surgeon-intern-app/internal/adminService/graph/model"
	"github.com/vbetsun/surgeon-intern-app/internal/adminService/pkg/services/encrypt"
	"github.com/vbetsun/surgeon-intern-app/internal/pkg/commonErrors"
	"github.com/vbetsun/surgeon-intern-app/internal/pkg/thirdparty"
	"github.com/vbetsun/surgeon-intern-app/pkg/models"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/dbexecutor"
	"github.com/vbetsun/surgeon-intern-app/pkg/repo/users"
	"github.com/volatiletech/null/v8"
	"gopkg.in/auth0.v5"
	"gopkg.in/auth0.v5/management"
)

type (
	UserService struct {
		auth0    *thirdparty.Auth0
		mailgun  *thirdparty.Mailgun
		userRepo users.IRepo
		dbexecutor.IDBExecutor
		auth0UserDbParameters *Auth0UserDbParameters
		hmac                  *encrypt.Hmac
		inviteLinkUrl         string
	}
	Auth0UserDbParameters struct {
		DbConnectionName string
		DbConnectionId   string
	}
	Filter struct {
	}
)

func NewUserService(auth0 *thirdparty.Auth0, userRepo users.IRepo, dbExecutor dbexecutor.IDBExecutor, mailgun *thirdparty.Mailgun, auth0UserDbParameters *Auth0UserDbParameters, hmac *encrypt.Hmac) *UserService {
	inviteLinkUrl := os.Getenv("USER_INVITE_LINK_URL")
	return &UserService{userRepo: userRepo, auth0: auth0, IDBExecutor: dbExecutor, mailgun: mailgun, auth0UserDbParameters: auth0UserDbParameters, hmac: hmac, inviteLinkUrl: inviteLinkUrl}
}

func (u *UserService) GetByFilter(ctx context.Context, filter Filter) ([]*qlmodel.DetailedUser, error) {
	usersWithUnitRoles, err := u.userRepo.ListByFilter(ctx, users.UserRepoFilter{})
	if err != nil {
		return nil, err
	}
	// TODO: fix error propagation ( Yes do it! )
	detailedUsers, err := u.FromUsersWithUnitRoles(usersWithUnitRoles)
	if err != nil {
		return nil, err
	}
	return detailedUsers, nil

}

func (u *UserService) GetByUserAndFilter(authorizedUserID int, filter Filter) ([]*qlmodel.DetailedUser, error) {
	return make([]*qlmodel.DetailedUser, 0), nil
}

func (u *UserService) UpdateUser(ctx context.Context, input *qlmodel.UpdateUserInput) (*qlmodel.DetailedUser, error) {
	retUser, err := u.userRepo.GetByID(ctx, input.UserID, true)
	if err != nil {
		return nil, err
	}
	if retUser.UserExternalID.IsZero() {
		return nil, commonErrors.NewInvalidInputError("User does not exist at the IDP")
	}
	err = u.auth0.Management.User.Update(retUser.UserExternalID.String, &management.User{FamilyName: &input.FamilyName, GivenName: &input.GivenName})
	if err != nil {
		return nil, err
	}
	retDetailedUser, err := u.FromUser(retUser)
	if err != nil {
		return nil, err
	}
	return retDetailedUser, nil
}

func (u *UserService) DeleteUser(ctx context.Context, userID string) error {
	dbUser, err := u.userRepo.GetByID(ctx, userID, false)
	if err != nil {
		return err
	}
	if !dbUser.UserExternalID.IsZero() {
		auth0User, err := u.auth0.Management.User.Read(dbUser.UserExternalID.String)
		if err != nil {
			return err
		}
		err = u.auth0.Management.User.Delete(auth0User.GetID())
		if err != nil {
			return err
		}
	}
	err = u.userRepo.Delete(ctx, dbUser.ID)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) GetByID(ctx context.Context, userID string) (*qlmodel.DetailedUser, error) {
	user, err := u.userRepo.GetByID(ctx, userID, true)
	if err != nil {
		return nil, err
	}
	detailedUser, err := u.FromUser(user)
	if err != nil {
		return nil, err
	}
	return detailedUser, nil
}

func (u *UserService) InviteUser(ctx context.Context, input qlmodel.InviteUserInput) (*qlmodel.DetailedUser, error) {
	exists, err := u.auth0ContainsUserWithEmail(input.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, commonErrors.NewInvalidInputError(fmt.Sprintf("User with email: %s already exists", input.Email))
	}
	var dbUser *models.User
	err = u.RunWithTX(ctx, func(tx *sql.Tx) error {
		newUser := &models.User{
			Activated: false,
			CreatedAt: time.Now(),
			Email:     input.Email,
		}
		newUser.Activationcode = null.StringFrom(u.getActivationCode(newUser))
		dbUser, err = u.userRepo.Add(ctx, newUser)
		if err != nil {
			return err
		}
		activationLink := u.getActivationLink(dbUser.Email, dbUser.Activationcode.String)
		err = u.mailgun.SendInviteMail(dbUser.Email, activationLink)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	detailedUser, err := u.FromUser(dbUser)
	return detailedUser, nil
}

func (u *UserService) ActivateUser(ctx context.Context, input qlmodel.ActivateUserInput) (*qlmodel.DetailedUser, error) {
	dbUser, err := u.userRepo.GetInactiveByActivationCode(ctx, input.ActivationCode, true)
	if err != nil {
		return nil, commonErrors.NewUnauthorizedAccessError("")
	}
	if !u.verifyUserActivationCode(dbUser, input.Email) {
		return nil, commonErrors.NewUnauthorizedAccessError("")
	}

	exists, err := u.auth0ContainsUserWithEmail(dbUser.Email)
	if err != nil {
		return nil, commonErrors.NewUnauthorizedAccessError("")
	}
	if exists {
		return nil, commonErrors.NewUnauthorizedAccessError("")
	}

	err = u.RunWithTX(ctx, func(tx *sql.Tx) error {
		newUser := &management.User{
			Connection:    auth0.String(u.auth0UserDbParameters.DbConnectionName),
			FamilyName:    auth0.String(input.FamilyName),
			GivenName:     auth0.String(input.GivenName),
			Name:          auth0.String(fmt.Sprintf("%s %s", input.GivenName, input.FamilyName)),
			Email:         auth0.String(input.Email),
			Password:      auth0.String(input.Password),
			EmailVerified: auth0.Bool(true),
		}
		err := u.auth0.Management.User.Create(newUser)
		if err != nil {
			return err
		}
		dbUser.UserExternalID = null.StringFromPtr(newUser.ID)
		dbUser.Activated = true
		dbUser.DisplayName = fmt.Sprintf("%s %s", input.GivenName, input.FamilyName)
		dbUser.Activationcode = null.StringFromPtr(nil)
		dbUser, err = u.userRepo.Update(ctx, dbUser)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	detailedUser, err := u.FromUser(dbUser)
	return detailedUser, nil
}

func (u *UserService) ValidateActivation(ctx context.Context, input *qlmodel.ActivationVerificationInput) (qlmodel.ActivationVerificationStatus, error) {
	dbUser, err := u.userRepo.GetInactiveByActivationCode(ctx, input.ActivationCode, false)
	if err != nil {
		return qlmodel.ActivationVerificationStatusInvalid, nil
	}
	if !u.verifyUserActivationCode(dbUser, input.Email) {
		return qlmodel.ActivationVerificationStatusInvalid, nil
	}
	exists, err := u.auth0ContainsUserWithEmail(dbUser.Email)
	if err != nil {
		return qlmodel.ActivationVerificationStatusInvalid, nil
	}
	if exists {
		return qlmodel.ActivationVerificationStatusInvalid, nil
	}
	return qlmodel.ActivationVerificationStatusActive, nil
}

func (u *UserService) auth0ContainsUserWithEmail(email string) (bool, error) {
	retUsers, err := u.auth0.Management.User.ListByEmail(email)
	if err != nil {
		return false, err
	}
	if len(retUsers) > 0 {
		return true, nil
	}
	return false, nil
}

func (u *UserService) getActivationLink(email string, activationCode string) string {
	return fmt.Sprintf("%s?email=%s&activationCode=%s", u.inviteLinkUrl, email, activationCode)
}

func (u *UserService) getActivationCode(user *models.User) string {
	return u.hmac.EncryptString(user.Email + strconv.FormatInt(user.CreatedAt.Unix(), 10))
}

func (u *UserService) verifyUserActivationCode(user *models.User, email string) bool {
	verified, err := u.hmac.VerifySignature(email+strconv.FormatInt(user.CreatedAt.Unix(), 10), user.Activationcode.String)
	if err != nil {
		return false
	}
	return verified
}
