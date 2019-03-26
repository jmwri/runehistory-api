package validate

import (
	"errors"
	"github.com/runehistory/runehistory-api/internal/domain/account"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStdAccountRules_IDIsPresent(t *testing.T) {
	a := assert.New(t)
	repo := new(account.MockRepository)
	rules := &StdAccountRules{
		AccountRepo: repo,
	}

	acc := &account.Account{
		ID: "present-id",
	}
	err := rules.IDIsPresent(acc)
	a.Nil(err, "not expecting err for present ID")

	acc = &account.Account{
		ID: "",
	}
	err = rules.IDIsPresent(acc)
	a.NotNil(err, "expecting error for vacant ID")
}

func TestStdAccountRules_IDIsCorrectLength(t *testing.T) {
	a := assert.New(t)
	repo := new(account.MockRepository)
	rules := &StdAccountRules{
		AccountRepo: repo,
	}

	acc := &account.Account{
		ID: "uuid-correct-length-1234567890123456",
	}
	err := rules.IDIsCorrectLength(acc)
	a.Nilf(err, "id with length %d should be valid", len(acc.ID))

	acc = &account.Account{
		ID: "uuid-incorrect-length",
	}
	err = rules.IDIsCorrectLength(acc)
	a.NotNilf(err, "id with length %d should be invalid", len(acc.ID))
}

func TestStdAccountRules_IDIsUnique_UniqueID(t *testing.T) {
	a := assert.New(t)

	repo := new(account.MockRepository)
	rules := &StdAccountRules{
		AccountRepo: repo,
	}

	acc := &account.Account{
		ID: "unique-id",
	}
	repo.On("CountId", acc.ID).Return(1, nil).Once()
	err := rules.IDIsUnique(acc)
	a.Nil(err, "expecting ID to be unique", acc.ID)
	repo.AssertExpectations(t)
}

func TestStdAccountRules_IDIsUnique_NonUniqueID(t *testing.T) {
	a := assert.New(t)

	repo := new(account.MockRepository)
	rules := &StdAccountRules{
		AccountRepo: repo,
	}

	acc := &account.Account{
		ID: "non-unique-id",
	}
	repo.On("CountId", acc.ID).Return(2, nil).Once()
	err := rules.IDIsUnique(acc)
	a.NotNil(err, "expecting duplicate ID", acc.ID)
	repo.AssertExpectations(t)
}

func TestStdAccountRules_IDIsUnique_Err(t *testing.T) {
	a := assert.New(t)

	repo := new(account.MockRepository)
	rules := &StdAccountRules{
		AccountRepo: repo,
	}

	acc := &account.Account{
		ID: "id-is-unique-err",
	}
	repo.On("CountId", acc.ID).Return(0, errors.New("expecting failure")).Once()
	err := rules.IDIsUnique(acc)
	a.NotNil(err, "expecting error")
	a.EqualError(err, "expecting failure")
	repo.AssertExpectations(t)
}

func TestStdAccountRules_IDWillBeUnique_Unique(t *testing.T) {
	a := assert.New(t)

	repo := new(account.MockRepository)
	rules := &StdAccountRules{
		AccountRepo: repo,
	}

	acc := &account.Account{
		ID: "will-be-unique-id",
	}
	repo.On("CountId", acc.ID).Return(1, nil).Once()
	err := rules.IDIsUnique(acc)
	a.Nil(err, "expecting id to be unique")
	repo.AssertExpectations(t)
}

func TestStdAccountRules_IDWillBeUnique_NonUnique(t *testing.T) {
	a := assert.New(t)

	repo := new(account.MockRepository)
	rules := &StdAccountRules{
		AccountRepo: repo,
	}

	acc := &account.Account{
		ID: "non-unique-id",
	}
	repo.On("CountId", acc.ID).Return(1, nil).Once()
	err := rules.IDWillBeUnique(acc)
	a.NotNilf(err, "expecting duplicate id: %s", acc.ID)
	repo.AssertExpectations(t)
}

func TestStdAccountRules_IDWillBeUnique_Error(t *testing.T) {
	a := assert.New(t)

	repo := new(account.MockRepository)
	rules := &StdAccountRules{
		AccountRepo: repo,
	}

	acc := &account.Account{
		ID: "id-is-unique-err",
	}
	repo.On("CountId", acc.ID).Return(0, errors.New("expecting failure")).Once()
	err := rules.IDWillBeUnique(acc)
	a.NotNil(err, "expecting error")
	a.EqualError(err, "expecting failure")
	repo.AssertExpectations(t)
}