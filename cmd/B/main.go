/*
	suppose we have slice of integers which every element is >=1 and all the numbers
	are repeated twice except for 1 number for example : {2,2,5,6,5} we want a way to
	find the exception number with only 1 iteration (loop) in this case the answer is 6
*/
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	/*
		key idea: i:=0
		if you see 2 in loop then:i+= 2
		then if you see it again i-=2
		then after sum and sub of seeing repeated number the value of i would not change unless we see exception number
		that there is no sub afterward. so the final value of i will be the exception number

		another solution:
		sort the array with sort.Ints(inputSliceOfIntegers) and check if two consecutive numbers are the same or not.
		If not, we check one more step and recognize the excep number.
	*/

	lookUpMap := make(map[int64]int, 0)
	var totalSum int64
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	inputSliceOfIntegers := strings.Split(scanner.Text(), " ")
	//inputSliceOfIntegers := []int64 {2, 2, 5, 6, 5}

	for _, val := range inputSliceOfIntegers {
		v, _ := strconv.ParseInt(val, 10, 64)
		if _, ok := lookUpMap[v]; !ok {
			lookUpMap[v] = 1
			totalSum += v
		} else {
			totalSum -= v
			delete(lookUpMap, v)
		}
	}

	fmt.Println("the exception number is: ", totalSum)
}
