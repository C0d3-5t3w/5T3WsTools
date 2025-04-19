package flagExt

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Flag represents an extended flag with additional features
type Flag struct {
	Name        string
	Usage       string
	Value       interface{}
	EnvVar      string
	Required    bool
	Validate    func(interface{}) error
	Initialized bool
}

// FlagSet extends the standard flag.FlagSet with additional functionality
type FlagSet struct {
	*flag.FlagSet
	flags         map[string]*Flag
	envPrefix     string
	errorHandling flag.ErrorHandling
}

// NewFlagSet creates a new FlagSet with the specified name and error handling policy
func NewFlagSet(name string, errorHandling flag.ErrorHandling) *FlagSet {
	return &FlagSet{
		FlagSet:       flag.NewFlagSet(name, errorHandling),
		flags:         make(map[string]*Flag),
		errorHandling: errorHandling,
	}
}

// CommandLine is the default set of command-line flags
var CommandLine = NewFlagSet(os.Args[0], flag.ExitOnError)

// String defines a string flag with specified name, default value, and usage string
func (f *FlagSet) String(name string, value string, usage string) *string {
	p := new(string)
	f.StringVar(p, name, value, usage)
	return p
}

// StringVar defines a string flag with specified name, default value, and usage string
// The argument p points to a string variable in which to store the value of the flag
func (f *FlagSet) StringVar(p *string, name string, value string, usage string) {
	f.FlagSet.StringVar(p, name, value, usage)
	f.flags[name] = &Flag{
		Name:  name,
		Usage: usage,
		Value: p,
	}
}

// StringVarWithOptions defines a string flag with additional options
func (f *FlagSet) StringVarWithOptions(p *string, name string, value string, usage string, required bool, envVar string, validate func(string) error) {
	f.StringVar(p, name, value, usage)
	flag := f.flags[name]
	flag.Required = required
	flag.EnvVar = envVar
	if validate != nil {
		flag.Validate = func(v interface{}) error {
			return validate(*v.(*string))
		}
	}
}

// Int defines an int flag with specified name, default value, and usage string
func (f *FlagSet) Int(name string, value int, usage string) *int {
	p := new(int)
	f.IntVar(p, name, value, usage)
	return p
}

// IntVar defines an int flag with specified name, default value, and usage string
func (f *FlagSet) IntVar(p *int, name string, value int, usage string) {
	f.FlagSet.IntVar(p, name, value, usage)
	f.flags[name] = &Flag{
		Name:  name,
		Usage: usage,
		Value: p,
	}
}

// IntVarWithOptions defines an int flag with additional options
func (f *FlagSet) IntVarWithOptions(p *int, name string, value int, usage string, required bool, envVar string, validate func(int) error) {
	f.IntVar(p, name, value, usage)
	flag := f.flags[name]
	flag.Required = required
	flag.EnvVar = envVar
	if validate != nil {
		flag.Validate = func(v interface{}) error {
			return validate(*v.(*int))
		}
	}
}

// Bool defines a bool flag with specified name, default value, and usage string
func (f *FlagSet) Bool(name string, value bool, usage string) *bool {
	p := new(bool)
	f.BoolVar(p, name, value, usage)
	return p
}

// BoolVar defines a bool flag with specified name, default value, and usage string
func (f *FlagSet) BoolVar(p *bool, name string, value bool, usage string) {
	f.FlagSet.BoolVar(p, name, value, usage)
	f.flags[name] = &Flag{
		Name:  name,
		Usage: usage,
		Value: p,
	}
}

// StringSlice defines a string slice flag with specified name, default value, and usage string
func (f *FlagSet) StringSlice(name string, value []string, usage string) *[]string {
	p := new([]string)
	*p = value
	f.Var(newStringSliceValue(p), name, usage)
	f.flags[name] = &Flag{
		Name:  name,
		Usage: usage,
		Value: p,
	}
	return p
}

// SetEnvPrefix sets a prefix for environment variables
func (f *FlagSet) SetEnvPrefix(prefix string) {
	f.envPrefix = prefix
}

// Parse parses the command line arguments and validates required flags
func (f *FlagSet) Parse(arguments []string) error {
	err := f.FlagSet.Parse(arguments)
	if err != nil {
		return err
	}

	// Track which flags have been set
	setFlags := make(map[string]bool)
	f.FlagSet.Visit(func(f *flag.Flag) {
		setFlags[f.Name] = true
	})

	// Handle environment variables and required flags
	for _, flag := range f.flags {
		// Check environment variables
		if flag.EnvVar != "" {
			envName := flag.EnvVar
			if f.envPrefix != "" && !strings.HasPrefix(envName, f.envPrefix) {
				envName = f.envPrefix + envName
			}

			if env := os.Getenv(envName); env != "" {
				switch v := flag.Value.(type) {
				case *string:
					*v = env
				case *int:
					i, err := strconv.Atoi(env)
					if err != nil {
						return fmt.Errorf("invalid value %q for environment variable %s: %v", env, envName, err)
					}
					*v = i
				case *bool:
					b, err := strconv.ParseBool(env)
					if err != nil {
						return fmt.Errorf("invalid value %q for environment variable %s: %v", env, envName, err)
					}
					*v = b
				case *[]string:
					*v = strings.Split(env, ",")
				}
				flag.Initialized = true
			}
		}

		// Validate required flags
		if flag.Required && !flag.Initialized && !setFlags[flag.Name] {
			return fmt.Errorf("flag -%s is required but not provided", flag.Name)
		}

		// Run validation functions
		if flag.Validate != nil {
			if err := flag.Validate(flag.Value); err != nil {
				return fmt.Errorf("validation failed for flag -%s: %v", flag.Name, err)
			}
		}
	}

	return nil
}

// String slice value implementation
type stringSliceValue struct {
	value *[]string
}

func newStringSliceValue(p *[]string) *stringSliceValue {
	return &stringSliceValue{p}
}

func (s *stringSliceValue) Set(val string) error {
	*s.value = append(*s.value, val)
	return nil
}

func (s *stringSliceValue) String() string {
	return strings.Join(*s.value, ",")
}

// Required functions that work with the default CommandLine FlagSet

// StringRequired defines a required string flag
func StringRequired(name string, value string, usage string) *string {
	p := new(string)
	StringVarRequired(p, name, value, usage)
	return p
}

// StringVarRequired defines a required string flag with a pointer
func StringVarRequired(p *string, name string, value string, usage string) {
	CommandLine.StringVarWithOptions(p, name, value, usage, true, "", nil)
}

// StringVarWithEnv defines a string flag with environment variable fallback
func StringVarWithEnv(p *string, name string, value string, usage string, envVar string) {
	CommandLine.StringVarWithOptions(p, name, value, usage, false, envVar, nil)
}

// StringVarRequiredWithEnv defines a required string flag with environment variable fallback
func StringVarRequiredWithEnv(p *string, name string, value string, usage string, envVar string) {
	CommandLine.StringVarWithOptions(p, name, value, usage, true, envVar, nil)
}

// IntRequired defines a required int flag
func IntRequired(name string, value int, usage string) *int {
	p := new(int)
	IntVarRequired(p, name, value, usage)
	return p
}

// IntVarRequired defines a required int flag with a pointer
func IntVarRequired(p *int, name string, value int, usage string) {
	CommandLine.IntVarWithOptions(p, name, value, usage, true, "", nil)
}

// Parse parses command line arguments using the default CommandLine FlagSet
func Parse() error {
	return CommandLine.Parse(os.Args[1:])
}
