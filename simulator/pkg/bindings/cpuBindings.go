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

/////////// Copy Pastes for updating bindings ///////////

// bindings.MdrIsInstructionBinding.Set(cpu.Registers.MDR.IsInstruction)
// bindings.MdrInstructionOpTypeBinding.Set(cpu.Registers.MDR.Instruction.OpType)
// bindings.MdrInstructionOpCodeBinding.Set(cpu.Registers.MDR.Instruction.Opcode)
// bindings.MdrInstructionOperandBinding.Set(cpu.Registers.MDR.Instruction.Operand)
// bindings.MdrDataBinding.Set(cpu.Registers.MDR.Data)

// bindings.InstructionOpTypeBinding.Set(cpu.Registers.IR.OpType)
// bindings.InstructionOpCodeBinding.Set(cpu.Registers.IR.Opcode)
// bindings.InstructionOperandBinding.Set(cpu.Registers.IR.Operand)

// bindings.MarBinding.Set(cpu.Registers.MAR)

// bindings.AcBinding.Set(cpu.Registers.AC)

// bindings.PcBinding.Set(cpu.Registers.PC)
