package trans

import (
	"fmt"
	"testing"
)

func TestTrans(t *testing.T) {
	result, err := Trans("Malicious Internet Entity Detection Using Local Graph Inference")
	if err != nil {
		fmt.Printf("[!] Error occurd: %v", err)
		return
	}
	fmt.Println(result)
}
