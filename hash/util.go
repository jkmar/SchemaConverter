package hash

func powers(x int) (log int, powers []int) {
	for x > 0 {
		if x%2 == 1 {
			powers = append(powers, log)
		}
		log++
		x /= 2
	}
	log--
	return
}
