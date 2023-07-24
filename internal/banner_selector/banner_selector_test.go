package banner_selector

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSelectIndex(t *testing.T) {
	t.Run("find in equal ranks", func(t *testing.T) {
		clicks := []int{3, 3, 3}
		displays := []int{10, 10, 10}
		index, err := SelectBannerIndex(displays, clicks)
		require.NoError(t, err)
		require.Equal(t, index, 0)
	})
	t.Run("select the least known", func(t *testing.T) {
		clicks := []int{8300, 1, 9100}
		displays := []int{16800, 10, 18000}
		bannerIndex, _ := SelectBannerIndex(displays, clicks)
		require.Equal(t, bannerIndex, 1)
	})
	t.Run("select non displayed", func(t *testing.T) {
		clicks := []int{8300, 9100, 0}
		displays := []int{16800, 18000, 0}
		index, err := SelectBannerIndex(displays, clicks)
		require.NoError(t, err)
		require.Equal(t, index, 2)
	})
}

func TestAtLeastOne(t *testing.T) {
	clicks := []int{0, 0, 0}
	displays := []int{0, 0, 0}
	for i := 0; i < 2000; i++ {
		index, _ := SelectBannerIndex(displays, clicks)
		displays[index]++
	}
	require.Greater(t, displays[0], 0)
	require.Greater(t, displays[1], 0)
	require.Greater(t, displays[2], 0)
}

func TestTheMostPopular(t *testing.T) {
	clicks := []int{3, 3, 3}
	displays := []int{50, 50, 50}
	for i := 0; i < 2000; i++ {
		index, _ := SelectBannerIndex(displays, clicks)
		if index == 0 {
			clicks[index]++
		}
		displays[index]++
	}
	require.Greater(t, displays[0], displays[1])
	require.Greater(t, displays[0], displays[2])
}

func TestIncorrectInput(t *testing.T) {
	t.Run("select with different input size", func(t *testing.T) {
		clicks := []int{2, 3}
		displays := []int{4, 5, 6}
		index, err := SelectBannerIndex(displays, clicks)
		require.ErrorIs(t, err, errIncorrectInput)
		require.Equal(t, index, invalidIndex)

		clicks = []int{2, 3, 4}
		displays = []int{4, 5}
		index, err = SelectBannerIndex(displays, clicks)
		require.ErrorIs(t, err, errIncorrectInput)
		require.Equal(t, index, invalidIndex)
	})
	t.Run("select than click greater than display", func(t *testing.T) {
		clicks := []int{4}
		displays := []int{3}
		index, err := SelectBannerIndex(displays, clicks)
		require.ErrorIs(t, err, errIncorrectInput)
		require.Equal(t, index, invalidIndex)
	})
}

func TestSelectWithEmptyInput(t *testing.T) {
	t.Run("select with empty clicks", func(t *testing.T) {
		clicks := []int{}
		displays := []int{3, 4, 5}
		index, err := SelectBannerIndex(displays, clicks)
		require.ErrorIs(t, err, errIncorrectInput)
		require.Equal(t, index, invalidIndex)
	})
	t.Run("select with empty displays", func(t *testing.T) {
		clicks := []int{2, 3, 4}
		displays := []int{}
		index, err := SelectBannerIndex(displays, clicks)
		require.ErrorIs(t, err, errIncorrectInput)
		require.Equal(t, index, invalidIndex)
	})

	t.Run("select with both empty inputs", func(t *testing.T) {
		clicks := []int{}
		displays := []int{}
		index, err := SelectBannerIndex(displays, clicks)
		require.ErrorIs(t, err, errIncorrectInput)
		require.Equal(t, index, invalidIndex)
	})
}

func TestNilInput(t *testing.T) {
	t.Run("select with nil clicks", func(t *testing.T) {
		displays := []int{3, 4, 5}
		index, err := SelectBannerIndex(displays, nil)
		require.ErrorIs(t, err, errIncorrectInput)
		require.Equal(t, index, invalidIndex)
	})

	t.Run("select with nil displays", func(t *testing.T) {
		clicks := []int{2, 3, 4}
		index, err := SelectBannerIndex(nil, clicks)
		require.ErrorIs(t, err, errIncorrectInput)
		require.Equal(t, index, invalidIndex)
	})

	t.Run("select with both nil inputs", func(t *testing.T) {
		index, err := SelectBannerIndex(nil, nil)
		require.ErrorIs(t, err, errIncorrectInput)
		require.Equal(t, index, invalidIndex)
	})
}

func TestEqualInputValues(t *testing.T) {
	t.Run("select when click is equal to display", func(t *testing.T) {
		clicks := []int{4, 3, 1}
		displays := []int{10, 3, 2}
		index, err := SelectBannerIndex(displays, clicks)
		require.NoError(t, err)
		require.Equal(t, index, 1)
	})

	t.Run("select when all clicks is equal to all displays", func(t *testing.T) {
		clicks := []int{10, 1000, 10000000}
		displays := []int{10, 1000, 10000000}
		bannerIndex, err := SelectBannerIndex(displays, clicks)
		require.NoError(t, err)
		require.Equal(t, bannerIndex, 0)
	})
}

func TestNegativeInput(t *testing.T) {

	t.Run("select with negative click", func(t *testing.T) {
		clicks := []int{-1}
		displays := []int{3}
		index, err := SelectBannerIndex(displays, clicks)
		require.ErrorIs(t, err, errIncorrectInput)
		require.Equal(t, index, invalidIndex)
	})

	t.Run("select with negative display", func(t *testing.T) {
		clicks := []int{2}
		displays := []int{-3}
		index, err := SelectBannerIndex(displays, clicks)
		require.ErrorIs(t, err, errIncorrectInput)
		require.Equal(t, index, invalidIndex)
	})
}
