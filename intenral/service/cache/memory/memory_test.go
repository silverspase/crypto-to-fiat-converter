package memory

import (
	"testing"

	price "crypto-to-fiat-converter/intenral/service/price_provider"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetTokenList(t *testing.T) {
	s := &service{}

	var testCases = []struct {
		testName       string
		tokenList      []price.TokenListItem
		pageToken      int32
		pageSize       int32
		expectedOutput []price.TokenListItem
		nextPageToken  int32
		error          error
	}{
		{
			testName:       "success",
			tokenList:      mockTokenList(),
			pageToken:      0,
			pageSize:       2,
			expectedOutput: mockTokenList()[0:2],
			nextPageToken:  1,
			error:          nil,
		},
		{
			testName:       "the_last_page_is_not_full.success",
			tokenList:      mockTokenList(),
			pageToken:      2,
			pageSize:       2,
			expectedOutput: mockTokenList()[4:],
			nextPageToken:  -1,
			error:          nil,
		},
		{
			testName:       "page_token_is_too_big.err",
			tokenList:      mockTokenList(),
			pageToken:      100,
			pageSize:       2,
			expectedOutput: nil,
			nextPageToken:  -1,
			error:          status.Errorf(codes.OutOfRange, "page token is too big"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			s.tokenList = tc.tokenList
			res, nextPageToken, err := s.GetTokenList(tc.pageToken, tc.pageSize)

			assert.Equal(t, tc.expectedOutput, res)
			assert.Equal(t, tc.nextPageToken, nextPageToken)
			if tc.error != nil {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func mockTokenList() []price.TokenListItem {
	return []price.TokenListItem{
		{
			ID:     "1",
			Name:   "first",
			Symbol: "F",
		},
		{
			ID:     "2",
			Name:   "second",
			Symbol: "S",
		},
		{
			ID:     "3",
			Name:   "third",
			Symbol: "T",
		},
		{
			ID:     "4",
			Name:   "fourth",
			Symbol: "FR",
		},
		{
			ID:     "5",
			Name:   "fifth",
			Symbol: "FF",
		},
	}
}
