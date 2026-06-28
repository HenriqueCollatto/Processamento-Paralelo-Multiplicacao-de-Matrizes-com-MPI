package main

import (
	"fmt"
	"math/rand"
	"time"
)

const N = 5000

func main() {
	fmt.Printf("Multiplicação de matrizes %dx%d (sequencial)\n", N, N)

	src := rand.NewSource(42)
	rng := rand.New(src)

	A := make([]float64, N*N)
	B := make([]float64, N*N)
	C := make([]float64, N*N)

	for i := range A {
		A[i] = rng.Float64()
	}
	for i := range B {
		B[i] = rng.Float64()
	}

	fmt.Println("Matrizes geradas. Iniciando multiplicação...")
	start := time.Now()

	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			sum := 0.0
			for k := 0; k < N; k++ {
				sum += A[i*N+k] * B[k*N+j]
			}
			C[i*N+j] = sum
		}
		if (i+1)%500 == 0 {
			fmt.Printf("  Linha %d/%d concluída (%.1fs)\n", i+1, N, time.Since(start).Seconds())
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("\nTempo total: %v\n", elapsed)

	fmt.Printf("\nVerificação (valores nos cantos da matriz C):\n")
	fmt.Printf("  C[0][0]       = %.15f\n", C[0])
	fmt.Printf("  C[0][N-1]     = %.15f\n", C[N-1])
	fmt.Printf("  C[N-1][0]     = %.15f\n", C[(N-1)*N])
	fmt.Printf("  C[N-1][N-1]   = %.15f\n", C[(N-1)*N+(N-1)])

	checksum := 0.0
	for i := range C {
		checksum += C[i]
	}
	fmt.Printf("  Checksum(C)   = %.15f\n", checksum)
}
