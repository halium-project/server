package user

import (
	"context"
	"fmt"
	"testing"

	"github.com/halium-project/go-server-utils/password"
	"github.com/halium-project/go-server-utils/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

const newUsername = "some@new.email"

func Test_User_Controller_Create(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("FindOneByUsername", ValidUser.Username).Return("", "", nil, nil).Once()
	uuidMock.On("New").Return("some-user-id").Once()
	passwordMock.On("HashWithSalt", "some-password").Return("some-hash", "some-salt", nil).Once()
	storageMock.On("Set", "some-user-id", "", &ValidUser).Return("some-rev", nil).Once()

	userID, err := controller.Create(context.Background(), &CreateCmd{
		Username: ValidUser.Username,
		Password: "some-password",
		Role:     ValidUser.Role,
	})

	assert.NoError(t, err)
	assert.Equal(t, "some-user-id", userID)

	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
}

func Test_User_Controller_Create_with_validation_error(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	userID, err := controller.Create(context.Background(), &CreateCmd{
		Username: ValidUser.Username,
		Password: "some-password",
		Role:     "invalid-role",
	})

	assert.Empty(t, userID)
	assert.JSONEq(t, `{
		"kind":"validationError",
		"errors": {
			"role":"UNEXPECTED_VALUE"
		}
	}`, err.Error())

	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
}

func Test_User_Controller_Create_storage_error(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("FindOneByUsername", ValidUser.Username).Return("", "", nil, nil).Once()
	uuidMock.On("New").Return("some-user-id").Once()
	passwordMock.On("HashWithSalt", "some-password").Return("some-hash", "some-salt", nil).Once()
	storageMock.On("Set", "some-user-id", "", &ValidUser).Return("", fmt.Errorf("some-error")).Once()

	userID, err := controller.Create(context.Background(), &CreateCmd{
		Username: ValidUser.Username,
		Password: "some-password",
		Role:     ValidUser.Role,
	})

	assert.Empty(t, userID)
	assert.JSONEq(t, `{
		"kind":"internalError",
		"message":"failed to save the user",
		"reason": {
			"kind":"internalError",
			"message":"some-error"
		}

	}`, err.Error())

	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
}

func Test_User_Controller_Create_password_hash_error(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("FindOneByUsername", ValidUser.Username).Return("", "", nil, nil).Once()
	uuidMock.On("New").Return("some-user-id").Once()
	passwordMock.On("HashWithSalt", "some-password").Return("", "", fmt.Errorf("some-error")).Once()

	userID, err := controller.Create(context.Background(), &CreateCmd{
		Username: ValidUser.Username,
		Password: "some-password",
		Role:     ValidUser.Role,
	})

	assert.Empty(t, userID)
	assert.JSONEq(t, `{
		"kind":"internalError",
		"message":"failed to hash the password",
		"reason": {
			"kind":"internalError",
			"message":"some-error"
		}

	}`, err.Error())

	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
}

func Test_User_Controller_Get(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("Get", "343c18cd-3bfc-48d0-bcad-180ce34dc948").Return("some-rev", &ValidUser, nil).Once()

	res, err := controller.Get(context.Background(), &GetCmd{
		UserID: "343c18cd-3bfc-48d0-bcad-180ce34dc948",
	})

	assert.NoError(t, err)
	assert.EqualValues(t, &ValidUser, res)

	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
}

func Test_User_Controller_Get_with_validationError(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	res, err := controller.Get(context.Background(), &GetCmd{
		UserID: "not a valid id",
	})

	assert.Nil(t, res)
	assert.JSONEq(t, `{
		"kind":"validationError",
		"errors":{
			"userID":"INVALID_FORMAT"
		}
	}`, err.Error())

	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
}

func Test_User_Controller_Get_driver_error(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("Get", "343c18cd-3bfc-48d0-bcad-180ce34dc948").Return("", nil, fmt.Errorf("some-error")).Once()

	res, err := controller.Get(context.Background(), &GetCmd{
		UserID: "343c18cd-3bfc-48d0-bcad-180ce34dc948",
	})

	assert.Nil(t, res)
	assert.JSONEq(t, `{
		"kind":"internalError",
		"message":"failed to get the user",
		"reason":{
			"kind": "internalError",
			"message": "some-error"
		}
	}`, err.Error())

	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
}

func Test_User_Controller_GetAll(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("GetAll").Return(map[string]User{
		"some-rev":   ValidUser,
		"some-rev-2": ValidUser,
	}, nil).Once()

	res, err := controller.GetAll(context.Background(), &GetAllCmd{})

	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.EqualValues(t, map[string]User{
		"some-rev":   ValidUser,
		"some-rev-2": ValidUser,
	}, res)

	uuidMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
}

func Test_User_Controller_GetAll_with_storage_error(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("GetAll").Return(nil, errors.New("some-error")).Once()

	res, err := controller.GetAll(context.Background(), &GetAllCmd{})

	assert.Nil(t, res)
	assert.JSONEq(t, `{
		"kind":"internalError",
		"message":"failed to get all users",
		"reason":{
			"kind":"internalError",
			"message":"some-error"
		}
	}`, err.Error())

	uuidMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
}

func Test_User_Controller_Validate(t *testing.T) {
	passwordMock := new(password.HashManagerMock)
	uuidMock := new(uuid.ProducerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("FindOneByUsername", "some-email@foo.bar").Return("some-user-id", "some-rev", &ValidUser, nil).Once()

	passwordMock.On("ValidateWithSalt", "some-password", "some-salt", "some-hash").Return(true, nil).Once()

	userID, user, err := controller.Validate(context.Background(), &ValidateCmd{
		Username: "some-email@foo.bar",
		Password: "some-password",
	})

	assert.NoError(t, err)
	assert.Equal(t, "some-user-id", userID)
	assert.EqualValues(t, &ValidUser, user)

	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
}

func Test_User_Controller_Validate_with_credentials_storage_error(t *testing.T) {
	passwordMock := new(password.HashManagerMock)
	uuidMock := new(uuid.ProducerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("FindOneByUsername", "some-email@foo.bar").Return("", "", nil, fmt.Errorf("some-error")).Once()

	userID, user, err := controller.Validate(context.Background(), &ValidateCmd{
		Username: "some-email@foo.bar",
		Password: "some-password",
	})

	assert.Empty(t, userID)
	assert.Nil(t, user)
	assert.JSONEq(t, `{
		"kind":"internalError",
		"message":"failed to get the credentials",
		"reason":{
			"kind":"internalError","message":"some-error"
		}
	}`, err.Error())

	passwordMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
}

func Test_User_Controller_Validate_with_unknown_email(t *testing.T) {
	passwordMock := new(password.HashManagerMock)
	uuidMock := new(uuid.ProducerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("FindOneByUsername", "some-invalid-email").Return("", "", nil, nil).Once()

	userID, user, err := controller.Validate(context.Background(), &ValidateCmd{
		Username: "some-invalid-email",
		Password: "some-password",
	})

	assert.NoError(t, err)
	assert.Nil(t, user)
	assert.Empty(t, userID)

	passwordMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
}

func Test_User_Controller_Validate_with_password_validationError(t *testing.T) {
	passwordMock := new(password.HashManagerMock)
	uuidMock := new(uuid.ProducerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("FindOneByUsername", "some-email@foo.bar").Return("some-user-id", "some-rev", &ValidUser, nil).Once()

	passwordMock.On("ValidateWithSalt", "some-password", "some-salt", "some-hash").Return(false, fmt.Errorf("some-error")).Once()

	userID, user, err := controller.Validate(context.Background(), &ValidateCmd{
		Username: "some-email@foo.bar",
		Password: "some-password",
	})

	assert.Empty(t, userID)
	assert.Nil(t, user)
	assert.JSONEq(t, `{
		"kind":"internalError",
		"message":"failed to compare the password with the hash",
		"reason":{
			"kind":"internalError","message":"some-error"
		}
	}`, err.Error())

	passwordMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
}

func Test_User_Controller_Validate_with_invalid_password(t *testing.T) {
	passwordMock := new(password.HashManagerMock)
	uuidMock := new(uuid.ProducerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("FindOneByUsername", "some-email@foo.bar").Return("some-user-id", "some-rev", &ValidUser, nil).Once()

	passwordMock.On("ValidateWithSalt", "some-invalid-password", "some-salt", "some-hash").Return(false, nil).Once()

	userID, user, err := controller.Validate(context.Background(), &ValidateCmd{
		Username: "some-email@foo.bar",
		Password: "some-invalid-password",
	})

	assert.NoError(t, err)
	assert.Empty(t, userID)
	assert.Nil(t, user)

	passwordMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
}

func Test_User_Controller_GetTotalUserCount(t *testing.T) {
	passwordMock := new(password.HashManagerMock)
	uuidMock := new(uuid.ProducerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("FindTotalUserCount").Return(42, nil).Once()

	count, err := controller.GetTotalUserCount(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, 42, count)

	passwordMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
}

func Test_User_Controller_GetTotalUserCount_with_error(t *testing.T) {
	passwordMock := new(password.HashManagerMock)
	uuidMock := new(uuid.ProducerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("FindTotalUserCount").Return(0, fmt.Errorf("some-error")).Once()

	count, err := controller.GetTotalUserCount(context.Background())

	assert.Equal(t, 0, count)
	assert.JSONEq(t, `{
		"kind":"internalError",
		"message":"failed to retrieve the number of user",
		"reason":{
			"kind":"internalError",
			"message":"some-error"
		}
	}`, err.Error())

	passwordMock.AssertExpectations(t)
	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
}

func Test_User_Controller_Create_with_email_checking_error(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("FindOneByUsername", ValidUser.Username).Return("", "", nil, fmt.Errorf("some-error")).Once()

	userID, err := controller.Create(context.Background(), &CreateCmd{
		Username: ValidUser.Username,
		Password: "some-password",
		Role:     ValidUser.Role,
	})

	assert.Empty(t, userID)
	assert.JSONEq(t, `{
		"kind":"internalError",
		"message":"failed to check if the user email is already taken",
		"reason":{
			"kind":"internalError",
			"message":"some-error"
		}
	}`, err.Error())

	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
}

func Test_User_Controller_Create_with_email_already_taken(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("FindOneByUsername", ValidUser.Username).Return("some-rev", "some-id", &ValidUser, nil).Once()

	userID, err := controller.Create(context.Background(), &CreateCmd{
		Username: ValidUser.Username,
		Password: "some-password",
		Role:     ValidUser.Role,
	})

	assert.Empty(t, userID)
	assert.JSONEq(t, `{
		"kind":"validationError",
		"errors":{
			"email":"ALREADY_USED"
		}
	}`, err.Error())

	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
}

func Test_User_Controller_Update_with_validation_error(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	err := controller.Update(context.Background(), &UpdateCmd{
		UserID:   "not a valid id",
		Username: ValidUser.Username,
		Role:     ValidUser.Role,
	})

	assert.JSONEq(t, `{
		"kind":"validationError",
		"errors":{
			"userID":"INVALID_FORMAT"
		}
	}`, err.Error())

	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
}

func Test_User_Controller_Update_with_storage_error(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("Get", "e16edc95-2063-4fc9-9f46-1431a0ddd6fa").Return("", nil, fmt.Errorf("some-error")).Once()

	err := controller.Update(context.Background(), &UpdateCmd{
		UserID:   "e16edc95-2063-4fc9-9f46-1431a0ddd6fa",
		Username: ValidUser.Username,
		Role:     ValidUser.Role,
	})

	assert.JSONEq(t, `{
		"kind":"internalError",
		"message": "failed to retrieve the user",
		"reason": {
			"kind": "internalError",
			"message": "some-error"
		}
	}`, err.Error())

	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
}

func Test_User_Controller_Update_with_email_validation_error(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("Get", "e16edc95-2063-4fc9-9f46-1431a0ddd6fa").Return("some-rev", &ValidUser, nil).Once()
	storageMock.On("FindOneByUsername", newUsername).Return("", "", nil, fmt.Errorf("some-error")).Once()

	err := controller.Update(context.Background(), &UpdateCmd{
		UserID:   "e16edc95-2063-4fc9-9f46-1431a0ddd6fa",
		Username: newUsername,
		Role:     ValidUser.Role,
	})

	assert.JSONEq(t, `{
		"kind":"internalError",
		"message": "failed to check if the user email is already taken",
		"reason": {
			"kind": "internalError",
			"message": "some-error"
		}
	}`, err.Error())

	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
}

func Test_User_Controller_Update_with_user_not_found(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("Get", "e16edc95-2063-4fc9-9f46-1431a0ddd6fa").Return("", nil, nil).Once()

	err := controller.Update(context.Background(), &UpdateCmd{
		UserID:   "e16edc95-2063-4fc9-9f46-1431a0ddd6fa",
		Username: ValidUser.Username,
		Role:     ValidUser.Role,
	})

	assert.JSONEq(t, `{
		"kind":"notFound"
	}`, err.Error())

	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
}

func Test_User_Controller_Update(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("Get", "e16edc95-2063-4fc9-9f46-1431a0ddd6fa").Return("some-rev", &ValidUser, nil).Once()
	storageMock.On("FindOneByUsername", newUsername).Return("", "", nil, nil).Once()

	newUser := ValidUser
	newUser.Username = newUsername
	storageMock.On("Set", "e16edc95-2063-4fc9-9f46-1431a0ddd6fa", "some-rev", &newUser).Return("some-new-rev", nil).Once()

	err := controller.Update(context.Background(), &UpdateCmd{
		UserID:   "e16edc95-2063-4fc9-9f46-1431a0ddd6fa",
		Username: newUsername,
		Role:     ValidUser.Role,
	})

	assert.NoError(t, err)

	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
}

func Test_User_Controller_Update_with_set_error(t *testing.T) {
	uuidMock := new(uuid.ProducerMock)
	passwordMock := new(password.HashManagerMock)
	storageMock := new(StorageMock)
	controller := NewController(uuidMock, passwordMock, storageMock)

	storageMock.On("Get", "e16edc95-2063-4fc9-9f46-1431a0ddd6fa").Return("some-rev", &ValidUser, nil).Once()
	storageMock.On("FindOneByUsername", newUsername).Return("", "", nil, nil).Once()

	newUser := ValidUser
	newUser.Username = newUsername
	storageMock.On("Set", "e16edc95-2063-4fc9-9f46-1431a0ddd6fa", "some-rev", &newUser).Return("", fmt.Errorf("some-error")).Once()

	err := controller.Update(context.Background(), &UpdateCmd{
		UserID:   "e16edc95-2063-4fc9-9f46-1431a0ddd6fa",
		Username: newUsername,
		Role:     ValidUser.Role,
	})

	assert.JSONEq(t, `{
		"kind":"internalError",
		"message": "failed to save the user",
		"reason": {
			"kind": "internalError",
			"message": "some-error"
		}
	}`, err.Error())

	storageMock.AssertExpectations(t)
	uuidMock.AssertExpectations(t)
	passwordMock.AssertExpectations(t)
}
