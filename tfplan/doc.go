// Package tfplan can read a subset of the Terraform plan format at various
// versions.
//
// The plan file format is not a documented interface and so it is subject to
// change at any time. This package makes a best effort to snapshot the format
// at various points to enable callers to load plans of different versions,
// but there is no guarantee that the result is complete, and new versions of
// the plan format can be released at any time which this library may or may
// not support.
package tfplan
