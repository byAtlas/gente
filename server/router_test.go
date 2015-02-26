package gente

import "testing"

func TestStaticRoute(t *testing.T) {
	router := RouterBuilder{}

	testPath := "aergaervaereerag"

	timesCalled := 0

	method := func(_ interface{}) (interface{}, error) {
		timesCalled++

		return testPath, nil
	}

	err := router.AddRoute(testPath, method)

	if err != nil {
		t.Error(err)
	}

	_, err = router.Finalize().Route(testPath, "string")

	if err != nil {
		t.Error(err)
	}

	if timesCalled != 1 {
		t.Errorf("Times called was %d, expected %d", timesCalled, 1)
	}
}

func TestRouterErrorsWhenAddingSamePathTwice(t *testing.T) {
	router := RouterBuilder{}
	testPath := "aergaervaereerag"
	method := func(_ interface{}) (interface{}, error) { return nil, nil }

	err := router.AddRoute(testPath, method)

	if err != nil {
		t.Error(err)
	}

	err = router.AddRoute(testPath, method)

	if err == nil {
		t.Error("Didn't get error from double insert")
	}
}
