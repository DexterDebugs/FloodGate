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
	proportional := pid.Kp * currentError			//Meaning: Gives high error when the error exceeds wayy past the limit. 
	pid.integral += currentError	//sums up all past errors 
	integralTerm := pid.Ki * pid.integral
	derivative := currentError - pid.prevError
	derivativeTerm := pid.Kd * derivative
	pid.prevError = currentError
	output := proportional + integralTerm + derivativeTerm

	return output
}