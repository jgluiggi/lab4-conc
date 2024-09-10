package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

type CrazyArray map[byte]int

// read a file from a filepath and return a slice of bytes
func readFile(filePath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file %s: %v", filePath, err)
		return nil, err
	}
	return data, nil
}

// sum all bytes of a file
func sum(filePath string, c chan int, c2 chan CrazyArray)  {
	data, err := readFile(filePath)
	if err != nil {
	}

    crazy := make(CrazyArray)

	_sum := 0
	for _, b := range data {
        _, existe := crazy[b]
        if existe {
            crazy[b] += 1
        } else {
            crazy[b] = 1
        }
		_sum += int(b)
	}

    c2 <- crazy
    c <- _sum
}

func similarity(map1 CrazyArray, map2 CrazyArray) float64 {
    counter := 0
    baseSize := 0

    for chave, valor1 := range map1 {
        valor2, existe := map2[chave]
        baseSize += valor1
        if existe {
            if valor1 < valor2 {
                counter += valor1
            } else {
                counter += valor2
            }
        }
    }
    return float64(counter) / float64(baseSize)
}

// print the totalSum for all files and the files with equal sum
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <file1> <file2> ...")
		return
	}
    c := make(chan int, 8)
    c2 := make(chan CrazyArray, 8)
	for _, path := range os.Args[1:] {
        go sum(path, c, c2)
    }

	var totalSum int64
	sums := make(map[int][]string)
	counts := make([]CrazyArray, len(os.Args))
	for i, path := range os.Args[1:] {
        _sum := <- c
		totalSum += int64(_sum)

		sums[_sum] = append(sums[_sum], path)

        counts[i] = <- c2
	}

	fmt.Println(totalSum)

	for sum, files := range sums {
		if len(files) > 1 {
			fmt.Printf("Sum %d: %v\n", sum, files)
		}
	}

    // partial similarity in motion
    for i := range os.Args[1:] {
        for j := i + 1; j < len(os.Args); j++ {
            similar := similarity(counts[i], counts[j])
            fmt.Printf("Similaridade entre %v, %v, %.2f/100  \n",os.Args[i], os.Args[j], similar*100)
        }
    }
}
