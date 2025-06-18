package dao

import (
	"database/sql"
	"github.com/kisara71/WeBook/webook/internal/domain"
)

func userDomainTOentity(dm *domain.User) UserEntity {
	return UserEntity{
		ID:   dm.Id,
		Name: dm.Name,
		Phone: sql.NullString{
			String: dm.Phone,
			Valid:  dm.Phone != "",
		},
		Email: sql.NullString{
			String: dm.Email,
			Valid:  dm.Email != "",
		},
		AboutMe:  dm.AboutMe,
		Birthday: dm.Birthday,
	}
}

func userEntityToDomain(ud *UserEntity) domain.User {
	return domain.User{
		Id:       ud.ID,
		Name:     ud.Name,
		Phone:    ud.Phone.String,
		Email:    ud.Email.String,
		AboutMe:  ud.AboutMe,
		Birthday: ud.Birthday,
	}
}

func oauth2BindingDomainToEntity(dm *domain.Oauth2Binding) Oauth2BindingEntity {
	return Oauth2BindingEntity{
		ID:     dm.ID,
		UserID: dm.UserID,

		Provider: sql.NullString{
			String: dm.Provider,
			Valid:  dm.Provider != "",
		},
		ExternalID: sql.NullString{
			String: dm.ExternalID,
			Valid:  dm.ExternalID != "",
		},
		AccessToken: sql.NullString{
			String: dm.AccessToken,
			Valid:  dm.AccessToken != "",
		},
	}
}

func oauth2BindingEntityToDomain(entity *Oauth2BindingEntity) domain.Oauth2Binding {
	return domain.Oauth2Binding{
		ID:          entity.ID,
		UserID:      entity.UserID,
		ExternalID:  entity.ExternalID.String,
		Provider:    entity.Provider.String,
		AccessToken: entity.AccessToken.String,
	}
}
