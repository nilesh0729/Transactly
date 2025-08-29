package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	mockDB "github.com/nilesh0729/OrdinaryBank/db/Mock"
	Anuskh "github.com/nilesh0729/OrdinaryBank/db/Result"
	"github.com/nilesh0729/OrdinaryBank/util"
	"github.com/stretchr/testify/require"
)

type eqCreateParamsMatcher struct {
	arg      Anuskh.CreateUserParams
	password string
}

func (e eqCreateParamsMatcher) Matches(x interface{}) bool {
	// In case, some value is nil
	arg, ok := x.(Anuskh.CreateUserParams)
	if !ok {
		return false
	}
	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}
	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg Anuskh.CreateUserParams, password string) gomock.Matcher {
	return eqCreateParamsMatcher{arg, password}
}

func TestCreateUserAPI(t *testing.T) {

	password, user := RandomUser(t)

	testcases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockDB.MockStore)
		CheckResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockDB.MockStore) {
				arg := Anuskh.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}

				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(user, nil)
			},
			CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				RequireBodyMatching(t, recorder.Body, user)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockDB.MockStore) {

				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(Anuskh.User{}, sql.ErrConnDone)
			},
			CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusInternalServerError)
			},
		},
		{
			name: "InvalidUser",
			body: gin.H{
				"username":  "Invalid-user#1",
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},

			buildStubs: func(store *mockDB.MockStore) {

				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
			},
		},
		{
			name: "DuplicateUsername",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(Anuskh.User{}, &pq.Error{Code: "23505"})
			},
			CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "TooShortPassword",
			body: gin.H{
				"username":  user.Username,
				"password":  "1234567",
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testcases {

		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockDB.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/user"

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.CheckResponse(t, recorder)
		})
	}
}

func RequireBodyMatching(t *testing.T, body *bytes.Buffer, user Anuskh.User) {

	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var GotUser Anuskh.User

	err = json.Unmarshal(data, &GotUser)
	require.NoError(t, err)

	require.Equal(t, GotUser.Username, user.Username)
	// require.Equal(t, GotUser.FullName, user.FullName)
	require.Equal(t, GotUser.Email, user.Email)
	require.Empty(t, GotUser.HashedPassword)
}

func RandomUser(t *testing.T) (password string, user Anuskh.User) {
	password = util.RandomString(8)

	hashedPassword, err := util.HashedPassword(password)
	require.NoError(t, err)

	user = Anuskh.User{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	return
}
