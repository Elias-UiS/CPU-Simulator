package bindings

import "fyne.io/fyne/v2/data/binding"

// Register bindings:
var SystemStatePcBinding binding.Int = binding.NewInt()
var SystemStateAcBinding binding.Int = binding.NewInt()
var SystemStateMarBinding binding.Int = binding.NewInt()
var SystemStateSpBinding binding.Int = binding.NewInt()

var SystemStateInstructionOpTypeBinding binding.Int = binding.NewInt()
var SystemStateInstructionOpCodeBinding binding.Int = binding.NewInt()
var SystemStateInstructionOperandBinding binding.Int = binding.NewInt()

var SystemStateMdrIsInstructionBinding binding.Bool = binding.NewBool()
var SystemStateMdrInstructionOpTypeBinding binding.Int = binding.NewInt()
var SystemStateMdrInstructionOpCodeBinding binding.Int = binding.NewInt()
var SystemStateMdrInstructionOperandBinding binding.Int = binding.NewInt()
var SystemStateMdrDataBinding binding.Int = binding.NewInt()

var SystemStateInstructionCount binding.Int = binding.NewInt()

var SystemStateNameBinding binding.String = binding.NewString()
