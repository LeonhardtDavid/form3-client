package accounts

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/LeonhardtDavid/form3-client/models"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAccountClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Account Client Suite")
}

var _ = Describe("Account Client", Ordered, func() {

	var accountClient AccountClient
	var account models.AccountData

	var httpServer *httptest.Server

	idToGet := uuid.New()
	idToDelete := uuid.New()
	idToFail := uuid.New()

	BeforeAll(func() {
		httpServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const path = "/v1/organisation/accounts"
			if r.URL.Path == path && r.Method == http.MethodPost {
				var res map[string]string
				_ = json.NewDecoder(r.Body).Decode(&res)
				if res["id"] == idToFail.String() {
					w.WriteHeader(http.StatusBadRequest)
					errorResponse := errorBody{ErrorMessage: "some invalid field"}
					jsonResponse, _ := json.Marshal(errorResponse)
					w.Write(jsonResponse)
				} else {
					w.WriteHeader(http.StatusCreated)
					accountJson, _ := json.Marshal(accountBody{account})
					w.Write(accountJson)
				}
			} else if r.URL.Path == fmt.Sprintf("%s/%s", path, idToGet) && r.Method == http.MethodGet {
				w.WriteHeader(http.StatusOK)
				accountJson, _ := json.Marshal(accountBody{account})
				w.Write(accountJson)
			} else if r.URL.Path == fmt.Sprintf("%s/%s", path, idToDelete) && r.Method == http.MethodDelete {
				if r.URL.Query().Get("version") == "0" {
					w.WriteHeader(http.StatusNoContent)
				} else {
					w.WriteHeader(http.StatusConflict)
					errorResponse := errorBody{ErrorMessage: "no version"}
					jsonResponse, _ := json.Marshal(errorResponse)
					w.Write(jsonResponse)
				}
			} else {
				w.WriteHeader(http.StatusNotFound)
				errorResponse := errorBody{ErrorMessage: "not found"}
				jsonResponse, _ := json.Marshal(errorResponse)
				w.Write(jsonResponse)
			}
		}))

		accountClient = NewAccountClient(
			AccountClientOptions{
				BaseURL: httpServer.URL,
			},
		)
	})

	AfterAll(func() {
		httpServer.Close()
	})

	BeforeEach(func() {
		var classification models.AccountClassification = models.Personal
		country := "GB"

		account = models.AccountData{
			ID:             uuid.New(),
			OrganisationID: uuid.New(),
			Type:           models.Accounts,
			Attributes: &models.AccountAttributes{
				AccountClassification: &classification,
				BankID:                "123456",
				BankIDCode:            "GBDSC",
				BaseCurrency:          "GBP",
				Bic:                   "EXMPLGB2XXX",
				Country:               &country,
				Name:                  []string{"David", "Test"},
			},
		}
	})

	Describe("Creating an account", func() {
		When("the returned status is Created", func() {
			It("should return the account without errors", func() {
				created, err := accountClient.Create(context.Background(), account)

				Expect(err).To(BeNil())
				Expect(*created).To(Equal(account))
			})
		})

		When("the returned status is Bad Request", func() {
			It("should return an error", func() {
				account.ID = idToFail
				created, err := accountClient.Create(context.Background(), account)

				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("some invalid field"))
				Expect(created).To(BeNil())
			})
		})
	})

	Describe("Fetching an account", func() {
		When("the returned status is Ok", func() {
			It("should return the account without errors", func() {
				fetched, err := accountClient.Fetch(context.Background(), idToGet)

				Expect(err).To(BeNil())
				Expect(*fetched).To(Equal(account))
			})
		})

		When("the returned status is Not Found", func() {
			It("should return an error", func() {
				fetched, err := accountClient.Fetch(context.Background(), uuid.New())

				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("not found"))
				Expect(fetched).To(BeNil())
			})
		})
	})

	Describe("Deleting an account", func() {
		When("the returned status is No Content", func() {
			It("should return no errors", func() {
				err := accountClient.Delete(context.Background(), idToDelete, 0)

				Expect(err).To(BeNil())
			})
		})

		When("the returned status is Conflict", func() {
			It("should return an error", func() {
				err := accountClient.Delete(context.Background(), idToDelete, 2)

				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("no version"))
			})
		})

		When("the returned status is Not Found", func() {
			It("should return an error", func() {
				err := accountClient.Delete(context.Background(), uuid.New(), 0)

				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("not found"))
			})
		})
	})

})
