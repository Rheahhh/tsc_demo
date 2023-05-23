package monitor

import (
	"testing"
	"tsc_demo/backend/common"
	"tsc_demo/backend/storage"

	"tsc_demo/backend/mocks"

	"github.com/stretchr/testify/assert"
)

func TestBrowserHistory(t *testing.T) {
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

	// Call the BrowserHistory method
	records, err := getter.BrowserHistory([]string{"http://example.com"})

	// Assert that there was no error and the records were correctly returned
	assert.NoError(t, err)
	assert.Equal(t, 1, len(records))
	assert.True(t, records[0].IsInBlacklist)

	// Assert that the expectations were met
	mockGetter.AssertExpectations(t)
}
