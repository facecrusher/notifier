package notifier

import (
	"errors"
	"testing"
	"time"

	"github.com/facecrusher/notifier/message/processor/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_Notify(t *testing.T) {
	ctrl := gomock.NewController(t)
	testDispatcher := mock.NewMockDispatcher(ctrl)
	url := "http://www.test.com"
	interval := time.Duration(time.Second)
	message := "test message"
	testNotifier := Notifier{
		url:               url,
		interval:          interval,
		messageDispatcher: testDispatcher,
	}

	testDispatcher.EXPECT().Dispatch(gomock.Any()).Return(nil).AnyTimes()
	err := testNotifier.Notify(message)
	assert.Nil(t, err)

	testDispatcher.EXPECT().Dispatch(gomock.Any()).Return(errors.New("")).AnyTimes()
	err = testNotifier.Notify(message)
	assert.Nil(t, err)
}
