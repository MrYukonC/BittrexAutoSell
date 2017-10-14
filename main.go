package main


import(
	"fmt"
	"time"
	"math"
	"github.com/toorop/go-bittrex"
)


const (
	API_KEY    = ""
	API_SECRET = ""
	MIN_ZEC_THRESH = 0.001
)


func CancelOrder( B *bittrex.Bittrex, OrderUUID string ) {
	if OrderUUID != "" {
		Err := B.CancelOrder( OrderUUID )

		if Err != nil {
			fmt.Println( Err )
		}
	}
}


func main() {

	B := bittrex.New( API_KEY, API_SECRET )

	SellUUID := ""

	for {
		ZecBalance, Err := B.GetBalance( "ZEC" )

		if Err == nil {

			fmt.Printf( "ZEC balance: %f\n", ZecBalance.Available )

			if ZecBalance.Available < MIN_ZEC_THRESH {
				break;
			}
	
			CancelOrder( B, SellUUID )
	
			EthZecOrderBook, Err := B.GetOrderBook( "ETH-ZEC", "buy", 1 )

			if Err != nil {
				time.Sleep( 3 * time.Second )
				continue
			}

			QuantityF64, _ := EthZecOrderBook.Buy[ 0 ].Quantity.Float64()
			RateF64,     _ := EthZecOrderBook.Buy[ 0 ].Rate.Float64()

			//fmt.Printf( "Highest Bid: %.8f ZEC @ %.8f ZEC/ETH\n", QuantityF64, RateF64 )

			FinalQuantity := math.Min( ZecBalance.Available, QuantityF64 )

			fmt.Printf( "Attempt to sell %.8f ZEC @ %.8f ZEC/ETH\n", FinalQuantity, RateF64 )

			UUID, Err := B.SellLimit( "ETH-ZEC", FinalQuantity, RateF64 )
		
			SellUUID = UUID
			
			if Err != nil {
				fmt.Println( Err )
			}

			time.Sleep( 3 * time.Second );
		}
	}
	
	time.Sleep( 3 * time.Second );

	CancelOrder( B, SellUUID )

	EthBalance, Err := B.GetBalance( "ETH" )
	
	if Err == nil {
		fmt.Printf( "ETH balance: %f\n", EthBalance.Available )
	}
}