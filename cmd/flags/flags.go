package flags

import (
	"fmt"

	"github.com/spf13/pflag"
)

// MustGetDefinedString attempts to get a non-empty string flag from the provided flag set or panic
func MustGetDefinedString(flagName string, flags *pflag.FlagSet) string {
	flagVal := MustGetString(flagName, flags)
	if flagVal == "" {
		fmt.Println(undefinedValueMessage(flagName))
	}
	return flagVal
}

// MustGetString attempts to get a string flag from the provided flag set or panic
func MustGetString(flagName string, flags *pflag.FlagSet) string {
	flagVal, err := flags.GetString(flagName)
	if err != nil {
		fmt.Println(notFoundMessage(flagName, err))
	}
	return flagVal
}

// MustGetBool attempts to get a boolean flag from the provided flag set or panic
func MustGetBool(flagName string, flags *pflag.FlagSet) bool {
	flagVal, err := flags.GetBool(flagName)
	if err != nil {
		fmt.Println(notFoundMessage(flagName, err))
	}
	return flagVal
}

func MustGetInt(flagName string, flags *pflag.FlagSet) int {
	flagVal, err := flags.GetInt(flagName)
	if err != nil {
		fmt.Println(notFoundMessage(flagName, err))
	}
	return flagVal
}

func undefinedValueMessage(flagName string) string {
	return fmt.Sprintf("flag %s has undefined value", flagName)
}

func notFoundMessage(flagName string, err error) string {
	return fmt.Sprintf("could not get flag %s from flag set: %s", flagName, err.Error())
}
