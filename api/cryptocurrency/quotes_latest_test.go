package cryptocurrency_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/Mikhalevich/coinmarketcap/api/cryptocurrency"
	"github.com/Mikhalevich/coinmarketcap/currency"
)

func TestQuotesLatestExecutorError(t *testing.T) {
	t.Parallel()

	const (
		btcID = 1
		ltcID = 2
		usdID = 2781
	)

	var (
		ctrl         = gomock.NewController(t)
		mockExecutor = cryptocurrency.NewMockExecutor(ctrl)
		cryptoc      = cryptocurrency.New(mockExecutor)
	)

	mockExecutor.EXPECT().
		Get(t.Context(), "/v2/cryptocurrency/quotes/latest", gomock.Any(), gomock.Any()).
		Return(errors.New("some executor error"))

	quotesRsp, err := cryptoc.QuotesLatest(
		t.Context(),
		[]currency.Currency{currency.ID(btcID), currency.ID(ltcID)},
		[]currency.Currency{currency.ID(usdID)},
	)

	require.Nil(t, quotesRsp)
	require.EqualError(t, err, "execute get request: some executor error")
}
