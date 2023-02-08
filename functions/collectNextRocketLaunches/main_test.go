package main

import (
	"testing"
)

//func Test_getNext5Launches(t *testing.T) {
//	tests := []struct {
//		name string
//		want models.UpcomingRocketLaunchesApiResponse
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := getNext5Launches(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("getNext5Launches() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

func TestRead(t *testing.T) {
	value := getNext5Launches()

	if value.Count <= 0 {
		t.Errorf("Expected value to be greater than 0")
	}
}
