package main

import (
	"fmt"
	"math/rand"
	"time"

	mpi "github.com/mvneves/gompi"
	"github.com/mvneves/gompi/comm"
)

// *******************************

const MATRIX_SIZE  = 5000;
const SIZE  	   = MATRIX_SIZE * MATRIX_SIZE;
const MASTER       = 0;
const SEED		   = 42;
const MATRIX_A_TAG = 0;
const MATRIX_B_TAG = 1;
const MATRIX_C_TAG = 2;
const ROWS_TAG     = 3

// *******************************

func main() {
	fmt.Printf("Multiplicação de matrizes %dx%d (paralell)\n", MATRIX_SIZE, MATRIX_SIZE)
	
	mpi.Init();
	defer mpi.Finalize();

	var com     comm.Communicator =  mpi.NewComm( true );
	var rank    int 			  =  com.GetRank();
	var numProc int			      =  com.GetSize();
	var src     rand.Source    	  =  rand.NewSource( SEED );
	var rng     rand.Rand         = *rand.New( src );

	if ( numProc < 2 ) {
		fmt.Println( "Ao menos dois processos são necessários" );
		return;
	}

	if ( rank == MASTER ) {
		M_Init( &com, numProc - 1, &rng );
		return;
	}

	S_Init( &com, numProc - 1 );
}

/*
=============================================

Master_Init()

Inicializa o Mestre, que envia ( com.Send ) 
um bloco de linhas para cada escravo processar
e coleta ( com.Recv ) para guardar na matrixC

=============================================
*/

func M_Init( com *comm.Communicator, numProc int, rng *rand.Rand ) {	
	matrixA   := make( []float64, SIZE );
	matrixB   := make( []float64, SIZE );
	matrixC   := make( []float64, SIZE );

	rowsBase := MATRIX_SIZE / numProc;

	mod := MATRIX_SIZE % numProc

	offset := 0

	InitializeArray( matrixA, rng );
	InitializeArray( matrixB, rng );

	var time0 time.Time = time.Now();

	for i := 0; i < numProc; i++ {

		rows := rowsBase
		if i == numProc-1 {
			rows += mod
		}

		blkSize := rows * MATRIX_SIZE
		slaveRank := i + 1

    	com.Send( []int{rows}, slaveRank, ROWS_TAG )

		com.Send( matrixA[ offset: offset+blkSize ], slaveRank, MATRIX_A_TAG )
    	com.Send( matrixB, slaveRank, MATRIX_B_TAG )

		offset += blkSize
	}

	offset = 0

	for i := 0; i < numProc; i++ {
		rows := rowsBase
		if i == numProc-1 {
			rows += mod
		}

		blkSize := rows * MATRIX_SIZE
		slaveRank := i + 1

		com.Recv( matrixC[offset: offset+blkSize], slaveRank, MATRIX_C_TAG )

		offset += blkSize
	}

	CalculateCheckSum( matrixA, matrixB, matrixC, time0 );
}

/*
=============================================

Slave_Init()

Inicializa um Slave, que recebe um bloco de
linhas da matrixA, a matrixB inteira e após
calcular o produto escalar entre elas, guarda
na matrixC

=============================================
*/

func S_Init( com *comm.Communicator, numProc int ) {

	var rows []int = make([]int, 1)

	com.Recv(rows, MASTER, ROWS_TAG)

	rowsProcs := rows[0]
	blkSize := rowsProcs * MATRIX_SIZE

	matrixA   := make( []float64, blkSize );
	matrixB   := make( []float64, SIZE );
	matrixC   := make( []float64, blkSize );

	com.Recv( matrixA, MASTER, MATRIX_A_TAG );
	com.Recv( matrixB, MASTER, MATRIX_B_TAG );

	for i := 0; i < rowsProcs; i++ {
		for j := 0; j < MATRIX_SIZE; j++ {
			var sum float64 = 0.0;

			for k := 0; k < MATRIX_SIZE; k++ {
				sum += matrixA[ i * MATRIX_SIZE + k ] * matrixB[ k * MATRIX_SIZE + j ];
			}

			matrixC[ i * MATRIX_SIZE + j ] = sum;
		}
	}

	com.Send( matrixC, MASTER, MATRIX_C_TAG );
}

func CalculateCheckSum( matrixA []float64, matrixB []float64, matrixC []float64, time0 time.Time ) {

	var totalTime time.Duration = time.Since( time0 );
	var checksum  float64     	= 0;

	fmt.Printf( "\nTempo total: %v\n", totalTime );

	fmt.Printf( "\nVerificação (valores nos cantos da matriz C):\n" );

	fmt.Printf( "  C[0][0]       = %.15f\n", matrixC[ 0 ] )
	fmt.Printf( "  C[0][N-1]     = %.15f\n", matrixC[ MATRIX_SIZE - 1 ] );
	fmt.Printf( "  C[N-1][0]     = %.15f\n", matrixC[ ( MATRIX_SIZE - 1 ) * MATRIX_SIZE ] );
	fmt.Printf( "  C[N-1][N-1]   = %.15f\n", matrixC[ ( MATRIX_SIZE - 1 ) * MATRIX_SIZE + ( MATRIX_SIZE -1) ] );

	for i := range matrixC {
		checksum += matrixC[ i ]
	}
	
	fmt.Printf( "  Checksum(C)   = %.15f\n", checksum )
}

func InitializeArray( matrix []float64, rng *rand.Rand ) {
	for i := range matrix {
		matrix[ i ] = rng.Float64();
	}
}