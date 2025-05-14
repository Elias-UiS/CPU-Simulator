package bindings

import "fyne.io/fyne/v2/data/binding"

// Register bindings:
var PcBinding binding.Int = binding.NewInt()
var AcBinding binding.Int = binding.NewInt()
var MarBinding binding.Int = binding.NewInt()
var SpBinding binding.Int = binding.NewInt()

var InstructionOpTypeBinding binding.Int = binding.NewInt()
var InstructionOpCodeBinding binding.Int = binding.NewInt()
var InstructionOperandBinding binding.Int = binding.NewInt()

var MdrIsInstructionBinding binding.Bool = binding.NewBool()
var MdrInstructionOpTypeBinding binding.Int = binding.NewInt()
var MdrInstructionOpCodeBinding binding.Int = binding.NewInt()
var MdrInstructionOperandBinding binding.Int = binding.NewInt()
var MdrDataBinding binding.Int = binding.NewInt()

var InstructionCount binding.Int = binding.NewInt()

var NameBinding binding.String = binding.NewString()
