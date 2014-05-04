package gerber_rs274x

import "fmt"

type MacroAperture struct {
	macroName string
}

func (macro* MacroAperture) AperturePlaceholder() {

}

func (macro* MacroAperture) GetHole() Hole {
	return nil
}

func (macro* MacroAperture) SetHole(hole Hole) {
	
}

func (macro* MacroAperture) String() string {
	return fmt.Sprintf("{MA, Name: %s}", macro.macroName)
}