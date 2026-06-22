package control

type PIDController struct {
	Kp float64	//proportional gain
	Ki float64	//integral gain
	Kd float64	//derivative gain
	integral float64	//accumulated error so far
	prevError float64	//the error from the previous call
}

func NewPIDController(kp, ki, kd float64)	*PIDController {
	return &PIDController{
		Kp: kp,
		Ki: ki,
		Kd: kd,
	}
}

func (pid *PIDController)	Compute (currentError float64) float64 {
	return currentError * pid.Kp
}