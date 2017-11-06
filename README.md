# go-tfplan

This Go library attempts to support loading Terraform plan files without
depending on Terraform itself, with varying degrees of success depending on
what part of the plan file you were interested in.

Plan files are just a serialization of some of Terraform's internal data
structures, so robust parsing from outside Terraform is not possible but
some basic partial handling can be done by mimicking Terraform's types.

This library is very much a "might work in a pinch" sort of thing, and is not
recommended for any real-world production use. It may or may not continue to
be maintained as new Terraform plan versions are implemented, and it definitely
does _not_ correctly parse the entire data structure. If you have any doubt
at all about using this then you shouldn't use it.

## Install

```
$ go install -u github.com/apparentlymart/go-tfplan/tfplan
```

## Usage

For reference documentation, see [godoc](https://godoc.org/github.com/apparentlymart/go-tfplan/tfplan).
