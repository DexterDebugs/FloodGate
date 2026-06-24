//Used for debugging purposes

package main

import (
	"fmt"
	"github.com/DexterDebugs/FloodGate/internal/control"
)

func main() {
	fmt.Println("=== Pure P (Kp=0.5, Ki=0, Kd=0) ===")
	pidP := control.NewPIDController(0.5,0,0)
	result := pidP.Compute(100)
	fmt.Println(result)
	result = pidP.Compute(100)
	fmt.Println(result)

	fmt.Println("=== Pure I (Kp=0, Ki=0.1, Kd=0) ===")
	pidI := control.NewPIDController(0,0.1,0)
	result = pidI.Compute(10)
	fmt.Println(result)
	result = pidI.Compute(10)
	fmt.Println(result)
	result = pidI.Compute(10)
	fmt.Println(result)

	fmt.Println("=== Pure D (Kp=0, Ki=0, Kd=0.2) ===")
	pidD := control.NewPIDController(0,0,0.2)
	result = pidD.Compute(0)
	fmt.Println(result)
	result = pidD.Compute(100)
	fmt.Println(result)
	result = pidD.Compute(100)
	fmt.Println(result)

}