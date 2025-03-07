package main

import (
	"fmt"
	"hash"
	"math"
	"math/rand"
	"time"

	"github.com/spaolacci/murmur3"
)

// BloomFilter estrutura do filtro
type BloomFilter struct {
	bitset   []bool
	size     uint
	hashes   uint
	hashFunc hash.Hash32
}

// NewBloomFilter cria um novo Bloom Filter
func NewBloomFilter(n uint, p float64) *BloomFilter {
	m := uint(math.Ceil(-float64(n) * math.Log(p) / (math.Ln2 * math.Ln2)))
	k := uint(math.Ceil(math.Ln2 * float64(m) / float64(n)))
	return &BloomFilter{
		bitset:   make([]bool, m),
		size:     m,
		hashes:   k,
		hashFunc: murmur3.New32(),
	}
}

// Add insere um item no filtro
func (bf *BloomFilter) Add(item string) {
	for i := uint(0); i < bf.hashes; i++ {
		index := bf.hash(item, i) % bf.size
		bf.bitset[index] = true
	}
}

// Contains verifica se um item pode estar no filtro
func (bf *BloomFilter) Contains(item string) bool {
	for i := uint(0); i < bf.hashes; i++ {
		index := bf.hash(item, i) % bf.size
		if !bf.bitset[index] {
			return false
		}
	}
	return true
}

// hash calcula o hash baseado em Murmur3
func (bf *BloomFilter) hash(item string, seed uint) uint {
	bf.hashFunc.Reset()
	bf.hashFunc.Write([]byte(item))
	h1 := bf.hashFunc.Sum32()
	return uint(h1) + seed*0x5bd1e995
}

func main() {
	// Parâmetros: 1000 registros, taxa de falso positivo de 1%
	bf := NewBloomFilter(1000, 0.01)

	// Criar alguns e-mails de exemplo
	rand.Seed(time.Now().UnixNano())
	emails := []string{"user1@example.com", "user2@example.com", "user3@example.com"}

	// Adicionando ao filtro
	for _, email := range emails {
		bf.Add(email)
	}

	// Testando se os e-mails estão no filtro
	testEmails := []string{"user1@example.com", "user4@example.com"}
	for _, email := range testEmails {
		if bf.Contains(email) {
			fmt.Printf("Provavelmente contém: %s\n", email)
		} else {
			fmt.Printf("Com certeza não contém: %s\n", email)
		}
	}
}
