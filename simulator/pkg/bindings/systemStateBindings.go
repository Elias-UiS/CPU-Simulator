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
