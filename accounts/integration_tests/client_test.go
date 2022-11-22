package accounts

import (
	"context"
	"fmt"
	"github.com/LeonhardtDavid/form3-client/accounts"
	"github.com/LeonhardtDavid/form3-client/models"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
	"testing"
)

func TestAccountClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Account Client Integration Suite")
}

var _ = Describe("Account Client", Ordered, func() {

	var accountClient accounts.AccountClient
	var account models.AccountData

	BeforeAll(func() {
		baseUrl := os.Getenv("FORM3_URL")

		Expect(baseUrl).ToNot(BeEmpty())

		accountClient = accounts.NewAccountClient(
			accounts.AccountClientOptions{
				BaseURL: baseUrl,
			},
		)
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

	Describe("Lifecycle", func() {
		When("the account is valid", func() {
			It("should be possible to fetch and delete the account", func() {
				var version int64 = 0
				account.Version = &version
				created, err := accountClient.Create(context.Background(), account)

				Expect(err).To(BeNil())
				Expect(*created).To(Equal(account))

				fetched, err := accountClient.Fetch(context.Background(), account.ID)

				Expect(err).To(BeNil())
				Expect(*fetched).To(Equal(account))

				err = accountClient.Delete(context.Background(), account.ID, 0)

				Expect(err).To(BeNil())
			})
		})

		When("the account is invalid", func() {
			It("should return an error", func() {
				account.Attributes = nil
				created, err := accountClient.Create(context.Background(), account)

				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("validation failure list:\nvalidation failure list:\nattributes in body is required"))
				Expect(created).To(BeNil())
			})
		})
	})

	Describe("Fetching an account", func() {
		When("the account doesn't exist", func() {
			It("should return an error", func() {
				id := uuid.New()
				fetched, err := accountClient.Fetch(context.Background(), id)

				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal(fmt.Sprintf("record %s does not exist", id)))
				Expect(fetched).To(BeNil())
			})
		})
	})

	Describe("Deleting an account", func() {

		When("the version is invalid", func() {
			It("should return an error", func() {
				_, err := accountClient.Create(context.Background(), account)
				err = accountClient.Delete(context.Background(), account.ID, 2)

				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("invalid version"))
			})
		})

		When("the account doesn't exist", func() {
			It("should return an error", func() {
				id := uuid.New()
				err := accountClient.Delete(context.Background(), id, 0)

				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("EOF")) // TODO this is because it return an empty response for 404.
			})
		})
	})

})
