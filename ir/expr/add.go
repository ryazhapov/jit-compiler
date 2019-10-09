package expr

import (
	"errors"
	"fmt"

	"github.com/bspaans/jit/asm"
	. "github.com/bspaans/jit/ir/shared"
)

type IR_Add struct {
	*BaseIRExpression
	Op1 IRExpression
	Op2 IRExpression
}

func NewIR_Add(op1, op2 IRExpression) *IR_Add {
	return &IR_Add{
		BaseIRExpression: NewBaseIRExpression(Add),
		Op1:              op1,
		Op2:              op2,
	}
}

func (i *IR_Add) ReturnType(ctx *IR_Context) Type {
	return i.Op1.ReturnType(ctx)
}

func (i *IR_Add) Encode(ctx *IR_Context, target *asm.Register) ([]asm.Instruction, error) {
	result := []asm.Instruction{}
	returnType1, returnType2 := i.Op1.ReturnType(ctx), i.Op2.ReturnType(ctx)
	if returnType1 == returnType2 && (returnType1 == TFloat64 || returnType1 == TUint64) {
		expr, err := i.Op1.Encode(ctx, target)
		if err != nil {
			return nil, err
		}
		for _, code := range expr {
			result = append(result, code)
		}

		if i.Op2.Type() == Variable {
			variable := i.Op2.(*IR_Variable).Value
			reg := ctx.VariableMap[variable]
			add := &asm.ADD{reg, target}
			ctx.AddInstruction(add)
			result = append(result, add)
		} else {
			reg := ctx.AllocateRegister(returnType1)
			defer ctx.DeallocateRegister(reg)

			expr, err := i.Op2.Encode(ctx, reg)
			if err != nil {
				return nil, err
			}
			for _, code := range expr {
				result = append(result, code)
			}
			add := &asm.ADD{reg, target}
			ctx.AddInstruction(add)
			result = append(result, add)
		}
	} else {
		return nil, errors.New("Unsupported types in add IR operation" + i.String())
	}
	return result, nil
}

func (i *IR_Add) String() string {
	return fmt.Sprintf("%s + %s", i.Op1.String(), i.Op2.String())
}
