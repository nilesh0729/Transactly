package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockDB "github.com/nilesh0729/Transactly/db/Mock"
	Anuskh "github.com/nilesh0729/Transactly/db/Result"
	"github.com/nilesh0729/Transactly/token"
	"github.com/nilesh0729/Transactly/util"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	_, user := RandomUser(t)
	account := randomAccount(user.Username)

	testCases := []struct {
		name          string
		accountid     int64
		setAuth       func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockDB.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountid: account.ID,
			setAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetAccounts(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchingAccount(t, recorder.Body, account)

			},
		},

		{
			name:      "NotFound",
			accountid: account.ID,
			setAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetAccounts(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(Anuskh.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
				// removed RequireBodyMatching() because there is no body to match or test
			},
		},

		{
			name:      "BadRequest",
			accountid: 0,
			setAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "BadRequest", time.Minute)
			},
			buildStubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

			},
		},

		{
			name:      "InternalServerError",
			accountid: account.ID,
			setAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetAccounts(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(Anuskh.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)

			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockDB.NewMockStore(ctrl)
			tc.buildStubs(store)

			//It starts test Server and sends Requests(like GetAccount)
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountid)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			// check response
			tc.checkResponse(t, recorder)
		})

	}

}

func TestPostAccountAPI(t *testing.T) {

	_, user := RandomUser(t)
	account1 := randomAccount(user.Username)

	arg := Anuskh.CreateAccountsParams{
		Owner:    account1.Owner,
		Currency: account1.Currency,
		Balance:  0,
	}

	testcases := []struct {
		name          string
		body          gin.H
		setAuth       func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockDB.MockStore)
		CheckResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			body: gin.H{
				"owner":    account1.Owner,
				"Currency": account1.Currency,
				"Balance":  account1.Balance,
			},
			setAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					CreateAccounts(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(account1, nil)
			},
			CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchingAccount(t, recorder.Body, account1)
			},
		},

		{
			name: "BadRequest",
			body: gin.H{
				"owner":    util.RandomOwner(),
				"Currency": util.RandomBalance(),
				"Balance":  util.RandomInt(1, 100),
			},
			setAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "BadRequest", time.Minute)
			},
			buildStubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					CreateAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

			},
		},

		{
			name: "InternalServerError",
			body: gin.H{
				"owner":    account1.Owner,
				"Currency": account1.Currency,
				"Balance":  0,
			},
			setAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					CreateAccounts(gomock.Any(), gomock.Any()).
					Times(1).
					Return(Anuskh.Account{}, sql.ErrConnDone)
			},
			CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)

			},
		},
	}

	for i := range testcases {

		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockDB.NewMockStore(ctrl)
			Server := newTestServer(t, store)

			tc.buildStubs(store)

			recorder := httptest.NewRecorder()
			url := "/accounts"

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)
			body := bytes.NewBuffer(data)
			request, err := http.NewRequest(http.MethodPost, url, body)
			require.NoError(t, err)

			tc.setAuth(t, request, Server.tokenMaker)

			Server.router.ServeHTTP(recorder, request)

			tc.CheckResponse(t, recorder)
		})
	}

}

func TestListAccountAPI(t *testing.T) {
	_, user := RandomUser(t)
	n := 5
	accounts := make([]Anuskh.Account, n)
	for i := 0; i < n; i++ {
		accounts[i] = randomAccount(user.Username)
	}

	type Query struct {
		page_size int
		page_id   int
	}

	testcases := []struct {
		name          string
		query         Query
		setAuth       func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockDB.MockStore)
		CheckResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",

			query: Query{
				page_size: n,
				page_id:   1,
			},

			setAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},

			buildStubs: func(store *mockDB.MockStore) {
				arg := Anuskh.ListAccountsParams{
					Owner:  user.Username,
					Limit:  int32(n),
					Offset: 0,
				}
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(accounts, nil)
			},
			CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchingAccounts(t, recorder.Body, accounts)
			},
		},

		{
			name: "BadRequest",
			query: Query{
				page_size: -1,
				page_id:   0,
			},

			setAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "badRequest", time.Minute)
			},
			buildStubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(0)

			},
			CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

			},
		},

		{
			name: "InternalServerError",
			query: Query{
				page_size: n,
				page_id:   1,
			},

			setAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},

			buildStubs: func(store *mockDB.MockStore) {

				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]Anuskh.Account{}, sql.ErrConnDone)

			},
			CheckResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)

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

			Server := newTestServer(t, store)

			recorder := httptest.NewRecorder()

			url := "/accounts"

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			q := request.URL.Query()
			q.Add("page_id", fmt.Sprintf("%d", tc.query.page_id))
			q.Add("page_size", fmt.Sprintf("%d", tc.query.page_size))
			request.URL.RawQuery = q.Encode()

			tc.setAuth(t, request, Server.tokenMaker)

			Server.router.ServeHTTP(recorder, request)

			tc.CheckResponse(t, recorder)
		})
	}

}

func requireBodyMatchingAccounts(t *testing.T, body *bytes.Buffer, expected []Anuskh.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var got []Anuskh.Account
	err = json.Unmarshal(data, &got)
	require.NoError(t, err)

	require.Equal(t, expected, got)
}

func requireBodyMatchingAccount(t *testing.T, body *bytes.Buffer, account Anuskh.Account) {

	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var GotAccount Anuskh.Account

	err = json.Unmarshal(data, &GotAccount)
	require.NoError(t, err)

	require.Equal(t, account, GotAccount)
}

func randomAccount(owner string) Anuskh.Account {
	return Anuskh.Account{
		ID:       util.RandomInt(1, 100),
		Balance:  20000,
		Owner:    owner,
		Currency: util.RandomCurrency(),
	}
}
