package tfplan

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io"

	"github.com/apparentlymart/go-tfplan/tfplan/v1"
	"github.com/apparentlymart/go-tfplan/tfplan/v2"
)

const planFormatMagic = "tfplan"

// LoadPlan reads a Terraform plan from the given reader and returns it.
//
// Different plan versions use different types, so the value returned may be
// any of the following types depending on the version:
//
//     github.com/apparentlymart/go-tfplan/tfplan/v1 *Plan (version 1)
//     github.com/apparentlymart/go-tfplan/tfplan/v2 *Plan (version 2)
//
// Use a type switch to recognize the versions you wish to support and then
// handle each version separately. New plan types may be returned by future
// versions, so be sure to handle the "default" case of the type switch.
func LoadPlan(src io.Reader) (interface{}, error) {
	var err error
	n := 0

	// Verify the magic bytes
	magic := make([]byte, len(planFormatMagic))
	for n < len(magic) {
		n, err = src.Read(magic[n:])
		if err != nil {
			return nil, fmt.Errorf("error while reading magic bytes: %s", err)
		}
	}
	if string(magic) != planFormatMagic {
		return nil, fmt.Errorf("not a valid plan file")
	}

	// Verify the version is something we can read
	var formatByte [1]byte
	n, err = src.Read(formatByte[:])
	if err != nil {
		return nil, err
	}
	if n != len(formatByte) {
		return nil, errors.New("failed to read plan version byte")
	}

	switch formatByte[0] {
	case 1:
		return loadPlanV1(src)
	case 2:
		return loadPlanV2(src)
	default:
		return nil, fmt.Errorf("unsupported plan file version %d", formatByte[0])
	}
}

func loadPlanV1(src io.Reader) (*v1.Plan, error) {
	var result v1.Plan

	dec := gob.NewDecoder(src)
	if err := dec.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func loadPlanV2(src io.Reader) (*v2.Plan, error) {
	var result v2.Plan

	dec := gob.NewDecoder(src)
	if err := dec.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
