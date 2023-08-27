package bloomfilter

import (
	"fmt"
	"hash/fnv"
	"math"
)

type BloomFilter struct {
	bitArray []bool
	hashFunc []hashFunc
}

type hashFunc func(data []byte) uint32

func NewBloomFilter(numItems int, falsePositiveRate float64) *BloomFilter {
	numBits := -(float64(numItems) * math.Log(falsePositiveRate)) / (math.Log(2) * math.Log(2))
	numHashFuncs := int((numBits / float64(numItems)) * math.Log(2))

	return &BloomFilter{
		bitArray: make([]bool, int(numBits)),
		hashFunc: generateHashFuncs(numHashFuncs),
	}
}

func (bf *BloomFilter) Add(data []byte) {
	for _, hf := range bf.hashFunc {
		index := hf(data) % uint32(len(bf.bitArray))
		bf.bitArray[index] = true
	}
}

func (bf *BloomFilter) Contains(data []byte) bool {
	for _, hf := range bf.hashFunc {
		index := hf(data) % uint32(len(bf.bitArray))
		if !bf.bitArray[index] {
			return false
		}
	}
	return true
}

func generateHashFuncs(numFuncs int) []hashFunc {
	hashFuncs := make([]hashFunc, numFuncs)
	for i := 0; i < numFuncs; i++ {
		hashFuncs[i] = generateHashFunc(i)
	}
	return hashFuncs
}

func generateHashFunc(seed int) hashFunc {
	return func(data []byte) uint32 {
		hash := fnv.New32a()
		hash.Write(data)
		return hash.Sum32() + uint32(seed)
	}
}

func main() {
	bloomFilter := NewBloomFilter(1000, 0.01)

	// Add items to Bloom filter
	bloomFilter.Add([]byte("item1"))
	bloomFilter.Add([]byte("item2"))
	bloomFilter.Add([]byte("item3"))

	// Check if items are in Bloom filter
	fmt.Println(bloomFilter.Contains([]byte("item1"))) // true
	fmt.Println(bloomFilter.Contains([]byte("item2"))) // true
	fmt.Println(bloomFilter.Contains([]byte("item3"))) // true
	fmt.Println(bloomFilter.Contains([]byte("item4"))) // false
}
