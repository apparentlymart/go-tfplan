package tfplan

import (
	"os"
	"testing"

	"github.com/apparentlymart/go-tfplan/tfplan/v1"
	"github.com/apparentlymart/go-tfplan/tfplan/v2"
)

func TestLoadPlanV1(t *testing.T) {
	f, err := os.Open("test-fixtures/v1-valid.tfplan")
	if err != nil {
		t.Fatal(err)
	}

	planI, err := LoadPlan(f)
	if err != nil {
		t.Fatal(err)
	}

	plan, isPlanV1 := planI.(*v1.Plan)
	if !isPlanV1 {
		t.Fatalf("result is %T; want *v1.Plan", planI)
	}

	if plan.Diff == nil {
		t.Fatalf("plan has nil diff")
	}
	if plan.Module == nil {
		t.Fatalf("plan has nil config")
	}
	if plan.State == nil {
		t.Fatalf("plan has nil state")
	}

	if got, want := plan.State.TFVersion, "0.9.10"; got != want {
		t.Errorf("wrong terraform version %q; want %q", got, want)
	}
}

func TestLoadPlanV2(t *testing.T) {
	f, err := os.Open("test-fixtures/v2-valid.tfplan")
	if err != nil {
		t.Fatal(err)
	}

	planI, err := LoadPlan(f)
	if err != nil {
		t.Fatal(err)
	}

	plan, isPlanV2 := planI.(*v2.Plan)
	if !isPlanV2 {
		t.Fatalf("result is %T; want *v2.Plan", planI)
	}

	if plan.Diff == nil {
		t.Fatalf("plan has nil diff")
	}
	if plan.Module == nil {
		t.Fatalf("plan has nil config")
	}
	if plan.State == nil {
		t.Fatalf("plan has nil state")
	}

	if got, want := plan.TerraformVersion, "0.10.9-dev"; got != want {
		t.Errorf("wrong terraform version %q; want %q", got, want)
	}
}
