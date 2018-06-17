package main

// Distinct returns the how many distinct urls are at the given date, where c contains the most significants components of the date
// [Year, Month, Day, Hour, Minute, Seccond]. For example in 2012 c is [2012]; in 2012-12 c is [2012, 12]
func Distinct(t Trie, c []int) int {
	node := t.Get(c)
	if node != nil && node.sortedUrls != nil {
		return len(*node.sortedUrls)
	}
	return 0
}
