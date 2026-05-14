package internal

import "github.com/MatthewSH/pdq"

const n = 256
const half = (n + 1) / 2

func TorbenMedian(m []float32) (float32, error) {
	if len(m) != n {
		return 0, pdq.ErrTorbenElementLength
	}

	min, max := m[0], m[0]

	for _, v := range m[1:] {
		if v < min {
			min = v
		}

		if v > max {
			max = v
		}
	}

	for {
		guess := (min + max) / 2

		var less, greater, equal int
		maxLTGuess := min
		minGTGuess := max

		for _, v := range m {
			if v < guess {
				less++

				if v > maxLTGuess {
					maxLTGuess = v
				}
			} else if v > guess {
				greater++

				if v < minGTGuess {
					minGTGuess = v
				}
			} else {
				equal++
			}

		}

		if less <= half && greater <= half {
			if less >= half {
				return maxLTGuess, nil
			}

			if less+equal >= half {
				return guess, nil
			}

			return minGTGuess, nil
		}

		if less > greater {
			max = maxLTGuess
		} else {
			min = minGTGuess
		}
	}
}
