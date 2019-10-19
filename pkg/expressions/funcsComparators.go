package expressions

import (
	"strconv"
	"strings"
)

func stringComparator(equation func(string, string) string) KeyBuilderFunction {
	return KeyBuilderFunction(func(args []KeyBuilderStage) KeyBuilderStage {
		if len(args) < 2 {
			return stageError(ErrorArgCount)
		}
		return KeyBuilderStage(func(context KeyBuilderContext) string {
			val := args[0](context)
			for i := 1; i < len(args); i++ {
				val = equation(val, args[i](context))
			}

			return val
		})
	})
}

// Checks equality, and returns truthy if equals, and empty if not
func arithmaticEqualityHelper(test func(int, int) bool) KeyBuilderFunction {
	return KeyBuilderFunction(func(args []KeyBuilderStage) KeyBuilderStage {
		if len(args) != 2 {
			return stageError(ErrorArgCount)
		}
		return KeyBuilderStage(func(context KeyBuilderContext) string {
			left, err := strconv.Atoi(args[0](context))
			if err != nil {
				return ErrorType
			}
			right, err := strconv.Atoi(args[1](context))
			if err != nil {
				return ErrorType
			}

			if test(left, right) {
				return "1"
			}
			return ""
		})
	})
}

func kfNot(args []KeyBuilderStage) KeyBuilderStage {
	if len(args) != 1 {
		return stageError(ErrorArgCount)
	}
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		if Truthy(args[0](context)) {
			return ""
		}
		return "1"
	})
}

// {and a b c ...}
func kfAnd(args []KeyBuilderStage) KeyBuilderStage {
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		for _, arg := range args {
			if arg(context) == "" {
				return ""
			}
		}
		return "1"
	})
}

// {or a b c ...}
func kfOr(args []KeyBuilderStage) KeyBuilderStage {
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		for _, arg := range args {
			if arg(context) != "" {
				return "1"
			}
		}
		return ""
	})
}

// {like string contains}
func kfLike(args []KeyBuilderStage) KeyBuilderStage {
	if len(args) != 2 {
		return stageError(ErrorArgCount)
	}
	return KeyBuilderStage(func(context KeyBuilderContext) string {
		val := args[0](context)
		contains := args[1](context)

		if strings.Contains(val, contains) {
			return val
		}
		return ""
	})
}