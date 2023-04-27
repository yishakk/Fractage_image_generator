package math

import (
	"errors"
	"math/cmplx"
	"strconv"
	"strings"
	"unicode"
)

const (
	state_GET_COEFFICIENT  = 0
	state_GET_POWER        = 1
	state_READ_COEFFICIENT = 2
	state_READ_VARIABLE    = 3
	state_READ_POWER       = 4

	term_num_sign      = 0
	term_pos_num_start = 1
	term_pos_num_end   = 2
	term_pos_variable  = 3
	term_pow_sign      = 4
	term_pos_pow_start = 5
	term_pos_pow_end   = 6
)

var (
	NIL_CMPLX_POLYNOMIAL = CmplxPolynomial{}
)

// Represents a term of a complex polynomial.
type PolynomialTerm struct {
	Coefficient float64
	Power       int
}

// Represents a complex polynomial
type CmplxPolynomial struct {
	Terms    []PolynomialTerm
	Variable rune
}

// Computes the first derivative of a complex polynomial.
func (polynomial *CmplxPolynomial) FirstDerivative() CmplxPolynomial {
	n := len(polynomial.Terms)
	for _, term := range polynomial.Terms {
		if term.Power == 0 {
			n--
		}
	}
	derivTerms := make([]PolynomialTerm, n)
	i := 0
	for _, term := range polynomial.Terms {
		if term.Power != 0 {
			derivTerms[i] = PolynomialTerm{
				Coefficient: term.Coefficient * float64(term.Power),
				Power:       term.Power - 1,
			}
			i++
		}
	}
	return CmplxPolynomial{Terms: derivTerms, Variable: polynomial.Variable}
}

// Evaluates the value of a complex polynomial for a given z.
func (polynomial *CmplxPolynomial) Evaluate(z complex128) complex128 {
	value := 0.0 + 0i
	for i := 0; i < len(polynomial.Terms); i++ {
		a := complex(polynomial.Terms[i].Coefficient, 0)
		if polynomial.Terms[i].Power == 0 {
			value += a
		} else {
			m := complex(float64(polynomial.Terms[i].Power), 0)
			value += a * cmplx.Pow(z, m)
		}
	}
	return value
}

// Converts a CmplxPolynomial type to its string representation.
func (polynomial *CmplxPolynomial) ToString() string {
	var sb strings.Builder
	for i, term := range polynomial.Terms {
		if i > 0 && term.Coefficient > 0 {
			sb.WriteRune('+')
		}
		sb.WriteString(strconv.FormatFloat(term.Coefficient, byte('f'), 4, 64))
		if term.Power != 0 {
			sb.WriteRune(polynomial.Variable)
			if term.Power != 1 {
				sb.WriteRune('^')
				sb.WriteString(strconv.FormatInt(int64(term.Power), 32))
			}
		}
	}
	return sb.String()
}

// Adds a term of a polynomial to the slice of polynomial terms.
func add_term(terms *[]PolynomialTerm, txt *[]rune, termInfo *[]int) error {
	var err error = nil
	coefficient := 1.0
	power := 0

	if (*termInfo)[term_pos_num_start] < 0 {
		coefficient = 1
	} else {
		fltTxt := string((*txt)[(*termInfo)[term_pos_num_start]:(*termInfo)[term_pos_num_end]])
		coefficient, err = strconv.ParseFloat(fltTxt, 32)
		if err != nil {
			return err
		}
		coefficient *= float64((*termInfo)[term_num_sign])
	}
	if (*termInfo)[term_pos_variable] < 0 {
		power = 0
	} else if (*termInfo)[term_pos_pow_start] < 0 {
		power = 1
	} else if (*termInfo)[term_pos_pow_start] > 0 {
		intTxt := string((*txt)[(*termInfo)[term_pos_pow_start]:(*termInfo)[term_pos_pow_end]])
		power, err = strconv.Atoi(intTxt)
		if err != nil {
			return err
		}
		power *= (*termInfo)[term_pow_sign]
	}
	if (*termInfo)[term_pos_num_start] > -1 || (*termInfo)[term_pos_variable] > -1 {
		*terms = append(*terms, PolynomialTerm{Coefficient: coefficient, Power: power})
	}
	// reset positions
	(*termInfo)[term_num_sign] = 1
	(*termInfo)[term_pos_num_start] = -1
	(*termInfo)[term_pos_num_end] = -1
	(*termInfo)[term_pos_variable] = -1
	(*termInfo)[term_pow_sign] = 1
	(*termInfo)[term_pos_pow_start] = -1
	(*termInfo)[term_pos_pow_end] = -1
	return err
}

// Constructs a CmplxPolynomial type from a mathematical expression.
func ParseCmplxPolynomial(txt string) (CmplxPolynomial, error) {
	var err error = nil
	var terms []PolynomialTerm
	var variable rune = ' '
	chars := []rune(txt)
	state := state_GET_COEFFICIENT
	// termItemPositions -> num_sign-num_start-num_end-variable_pos-pow_sign-pow_start-pow_end
	termItemPositions := []int{1, -1, -1, -1, 1, -1, -1}
	i := 0
	for {
		if i >= len(chars) {
			if state >= state_READ_COEFFICIENT {
				err = add_term(&terms, &chars, &termItemPositions)
				if err != nil {
					return NIL_CMPLX_POLYNOMIAL, err
				}
				break
			} else {
				return NIL_CMPLX_POLYNOMIAL, errors.New("Invalid polynomial")
			}
		}
		if unicode.IsLetter(chars[i]) {
			if variable != ' ' && chars[i] != variable {
				return NIL_CMPLX_POLYNOMIAL, errors.New("Multiple variables in polynomial")
			}
			if state == state_READ_VARIABLE {
				return NIL_CMPLX_POLYNOMIAL, errors.New("Variables must be a single character")
			}
			variable = chars[i]
			state = state_READ_VARIABLE
			termItemPositions[term_pos_variable] = i
		} else if unicode.IsDigit(chars[i]) {
			a := i
			for i < len(chars) {
				if unicode.IsDigit(chars[i]) || chars[i] == '.' {
					i++
				} else {
					break
				}
			}
			if state == state_GET_POWER {
				termItemPositions[term_pos_pow_start] = a
				termItemPositions[term_pos_pow_end] = i
				state = state_READ_POWER
			} else if state == state_GET_COEFFICIENT {
				termItemPositions[term_pos_num_start] = a
				termItemPositions[term_pos_num_end] = i
				state = state_READ_COEFFICIENT
			}
			continue
		} else if chars[i] == '^' {
			if state != state_READ_VARIABLE {
				return NIL_CMPLX_POLYNOMIAL, errors.New("Invalid polynomial")
			}
			state = state_GET_POWER
		} else if chars[i] == '+' || chars[i] == '-' {
			var sign int = 1
			if chars[i] == '-' {
				sign = -1
			}
			if state == state_GET_COEFFICIENT {
				termItemPositions[term_num_sign] *= sign
			} else if state == state_GET_POWER {
				termItemPositions[term_pow_sign] *= sign
			} else if state >= state_READ_COEFFICIENT {
				err = add_term(&terms, &chars, &termItemPositions)
				if err != nil {
					return NIL_CMPLX_POLYNOMIAL, err
				}
				termItemPositions[term_num_sign] *= sign
				state = state_GET_COEFFICIENT
			}
		} else if chars[i] != ' ' {
			return NIL_CMPLX_POLYNOMIAL, errors.New("Invalid polynomial")
		}
		i++
	}
	return CmplxPolynomial{Terms: terms, Variable: variable}, nil
}
