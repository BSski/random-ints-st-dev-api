package randomintsstdev

import "github.com/montanaflynn/stats"

func getRoundedStDev(intSeq []int) (roundedStDev float64, err error) {
	data := stats.LoadRawData(intSeq)
	stDev, err := stats.StandardDeviation(data)
	if err != nil {
		return
	}
	roundedStDev, err = stats.Round(stDev, 3)
	if err != nil {
		return
	}
	return
}

func getIntSeqsSum(intSeqs [][]int) (intSeqsSum []int) {
	intSeqsSum = make([]int, len(intSeqs)*len(intSeqs[0]))
	for i, seq := range intSeqs {
		for j, x := range seq {
			intSeqsSum[i*len(intSeqs[0])+j] = x
		}
	}
	return
}
