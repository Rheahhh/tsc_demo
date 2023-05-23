package monitor

import (
	"testing"
	"tsc_demo/backend/common"
	"tsc_demo/backend/mocks"
	"tsc_demo/backend/storage"

	"github.com/stretchr/testify/assert"
)

func TestBrowserHistory(t *testing.T) {
	TestBrowserHistoryEmptyBlacklist(t)
	TestBrowserHistoryEmptyHistory(t)
	TestBrowserHistoryDatabaseQueryFail(t)
	TestBrowserHistoryWithBlacklist(t)
}

func TestBrowserHistoryEmptyBlacklist(t *testing.T) {
	// Create a mock raw history getter
	mockGetter := &mocks.RawHistoryGetter{}

	// Set up the expectation for GetBrowserHistory method
	mockGetter.On("GetBrowserHistory").Return([]common.BrowserHistoryRecord{
		{
			URL:           "http://example.com",
			Title:         "Example",
			VisitCount:    1,
			LastVisitTime: 1000,
			IsInBlacklist: false,
		},
	}, nil)

	// Use the mock getter as the default raw history getter
	storage.SetRawHistoryGetter(mockGetter)

	// Create a BrowserHistoryGetter with the mock raw history getter
	getter := NewBrowserHistoryGetter()

	// Call the BrowserHistory method with an empty blacklist
	records, err := getter.BrowserHistory([]string{})

	// Assert that there was no error and the records were correctly returned
	assert.NoError(t, err)
	assert.Equal(t, 1, len(records))
	assert.False(t, records[0].IsInBlacklist)

	// Assert that the expectations were met
	mockGetter.AssertExpectations(t)
}

func TestBrowserHistoryEmptyHistory(t *testing.T) {
	// Create a mock raw history getter
	mockGetter := &mocks.RawHistoryGetter{}

	// Set up the expectation for GetBrowserHistory method
	mockGetter.On("GetBrowserHistory").Return([]common.BrowserHistoryRecord{}, nil)

	// Use the mock getter as the default raw history getter
	storage.SetRawHistoryGetter(mockGetter)

	// Create a BrowserHistoryGetter with the mock raw history getter
	getter := NewBrowserHistoryGetter()

	// Call the BrowserHistory method with a blacklist containing "http://example.com"
	records, err := getter.BrowserHistory([]string{"http://example.com"})

	// Assert that there was no error and the records slice is empty
	assert.NoError(t, err)
	assert.Empty(t, records)

	// Assert that the expectations were met
	mockGetter.AssertExpectations(t)
}

func TestBrowserHistoryDatabaseQueryFail(t *testing.T) {
	// Create a mock raw history getter
	mockGetter := &mocks.RawHistoryGetter{}

	// Set up the expectation for GetBrowserHistory method
	mockGetter.On("GetBrowserHistory").Return(nil, common.ErrDatabaseQuery)

	// Use the mock getter as the default raw history getter
	storage.SetRawHistoryGetter(mockGetter)

	// Create a BrowserHistoryGetter with the mock raw history getter
	getter := NewBrowserHistoryGetter()

	// Call the BrowserHistory method with a blacklist containing "http://example.com"
	records, err := getter.BrowserHistory([]string{"http://example.com"})

	// Assert that the error is of type DatabaseQueryError
	assert.Error(t, err)
	assert.IsType(t, common.ErrDatabaseQuery, err)

	// Assert that the records slice is nil
	assert.Nil(t, records)

	// Assert that the expectations were met
	mockGetter.AssertExpectations(t)
}

func TestBrowserHistoryWithBlacklist(t *testing.T) {
	// Create a mock raw history getter
	mockGetter := &mocks.RawHistoryGetter{}

	// Set up the expectation for GetBrowserHistory method
	mockGetter.On("GetBrowserHistory").Return([]common.BrowserHistoryRecord{
		{
			URL:           "http://example.com",
			Title:         "Example",
			VisitCount:    1,
			LastVisitTime: 1000,
			IsInBlacklist: true,
		},
	}, nil)

	// Use the mock getter as the default raw history getter
	storage.SetRawHistoryGetter(mockGetter)

	// Create a BrowserHistoryGetter with the mock raw history getter
	getter := NewBrowserHistoryGetter()

	// Call the BrowserHistory method with a blacklist containing "http://example.com"
	records, err := getter.BrowserHistory([]string{"http://example.com"})

	// Assert that there was no error and the records were correctly returned
	assert.NoError(t, err)
	assert.Equal(t, 1, len(records))
	assert.True(t, records[0].IsInBlacklist)

	// Assert that the expectations were met
	mockGetter.AssertExpectations(t)
}
