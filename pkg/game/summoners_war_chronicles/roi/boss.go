package roi

type MiniBossName string

const (
	MiniBossBlackAshHarpu = "black_ash_harpu"
)

var (
	MiniBosses = map[MiniBossName]struct{}{
		MiniBossBlackAshHarpu: {},
	}
)