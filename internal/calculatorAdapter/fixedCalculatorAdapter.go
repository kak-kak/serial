package calculatorAdapter

type FixedCalculatorAdapter struct {
	fixedEstimate Estimate
}

func NewFixedCalculatorAdapter(fixedEstimate Estimate) *FixedCalculatorAdapter {
	return &FixedCalculatorAdapter{
		fixedEstimate: fixedEstimate,
	}
}

func (f *FixedCalculatorAdapter) Calculate(estimates chan<- Estimate, data []byte) error {
	estimates <- f.fixedEstimate
	return nil
}

func (f *FixedCalculatorAdapter) Close() {
}
