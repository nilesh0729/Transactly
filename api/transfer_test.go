package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockDB "github.com/nilesh0729/OrdinaryBank/db/Mock"
	Anuskh "github.com/nilesh0729/OrdinaryBank/db/Result"
	"github.com/nilesh0729/OrdinaryBank/util"
	"github.com/stretchr/testify/require"
)

func TestTransfersAPI(t *testing.T) {
	account1 := randomAccount()
	account2 := randomAccount()
	account3 := randomAccount()
	amount := int64(10)

	account1.Currency = util.INR
	account2.Currency = util.INR
	account3.Currency = util.USD

	testcases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockDB.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        util.INR,
			},
			buildStubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetAccounts(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(account1, nil)

				store.EXPECT().
					GetAccounts(gomock.Any(), gomock.Eq(account2.ID)).
					Times(1).
					Return(account2, nil)

				arg := Anuskh.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}
				store.EXPECT().
				TransferTx(gomock.Any(), gomock.Eq(arg)).
				Times(1)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name : "FromAccountNotFound",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id": account2.ID,
				"amount": amount,
				"currency": util.INR,
			},
			buildStubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetAccounts(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(Anuskh.Account{}, sql.ErrNoRows)

				store.EXPECT().
					GetAccounts(gomock.Any(), gomock.Eq(account2.ID)).
					Times(0)
					
				arg := Anuskh.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}
				store.EXPECT().
				TransferTx(gomock.Any(), gomock.Eq(arg)).
				Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},

		{
			name : "ToAccountNotFound",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id": account2.ID,
				"amount": amount,
				"currency": util.INR,
			},
			buildStubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetAccounts(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(account1, nil)

				store.EXPECT().
					GetAccounts(gomock.Any(), gomock.Eq(account2.ID)).
					Times(1).
					Return(Anuskh.Account{},sql.ErrNoRows)
					
				arg := Anuskh.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}
				store.EXPECT().
				TransferTx(gomock.Any(), gomock.Eq(arg)).
				Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},

		{
			name: "FromAccountCurrencyMismatched",
			body: gin.H{
				"from_account_id": account3.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        util.INR,
			},
			buildStubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetAccounts(gomock.Any(), gomock.Eq(account3.ID)).
					Times(1).
					Return(account3, nil)

				store.EXPECT().
					GetAccounts(gomock.Any(), gomock.Eq(account1.ID)).
					Times(0)

				arg := Anuskh.TransferTxParams{
					FromAccountID: account3.ID,
					ToAccountID:   account1.ID,
					Amount:        amount,
				}
				store.EXPECT().
				TransferTx(gomock.Any(), gomock.Eq(arg)).
				Times(0)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		{
			name: "ToAccountCurrencyMismatched",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account3.ID,
				"amount":          amount,
				"currency":        util.INR,
			},
			buildStubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetAccounts(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(account3, nil)

				store.EXPECT().
					GetAccounts(gomock.Any(), gomock.Eq(account3.ID)).
					Times(0)

				arg := Anuskh.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account3.ID,
					Amount:        amount,
				}
				store.EXPECT().
				TransferTx(gomock.Any(), gomock.Eq(arg)).
				Times(0)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		{
			name: "InvalidCurrency",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        "xyz",
			},
			buildStubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetAccounts(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					GetAccounts(gomock.Any(), gomock.Any()).
					Times(0)

				arg := Anuskh.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}
				store.EXPECT().
				TransferTx(gomock.Any(), gomock.Eq(arg)).
				Times(0)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},	
		},

		{
			name: "NegativeAmount",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          -amount,
				"currency":        util.INR,
			},
			buildStubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetAccounts(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					GetAccounts(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
				TransferTx(gomock.Any(), gomock.Any()).
				Times(0)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},

		{
			name: "GetAccountError",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        util.INR,
			},
			buildStubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetAccounts(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(Anuskh.Account{}, sql.ErrConnDone)

				store.EXPECT().
					GetAccounts(gomock.Any(), gomock.Eq(account2.ID)).
					Times(0)

				arg := Anuskh.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}
				store.EXPECT().
				TransferTx(gomock.Any(), gomock.Eq(arg)).
				Times(0)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},

		{
			name: "transferTxError",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        util.INR,
			},
			buildStubs: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetAccounts(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(account1, nil)

				store.EXPECT().
					GetAccounts(gomock.Any(), gomock.Eq(account2.ID)).
					Times(1).
					Return(account2, nil)

				arg := Anuskh.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}
				store.EXPECT().
				TransferTx(gomock.Any(), gomock.Eq(arg)).
				Times(1).
				Return(Anuskh.TransferTxResult{},sql.ErrConnDone)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
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

			
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/transfers"

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			Server.router.ServeHTTP(recorder, request)

			tc.checkResponse(recorder)
		})
	}

}
