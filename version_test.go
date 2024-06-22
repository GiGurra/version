package version

import (
	"fmt"
	"strconv"
	"testing"
)

func regular(major int, minor int, patch int) Version {
	return Version{
		Major: VersionPart{
			Raw:     strconv.Itoa(major),
			IntRepr: major,
			IsInt:   true,
		},
		Minor: VersionPart{
			Raw:     strconv.Itoa(minor),
			IntRepr: minor,
			IsInt:   true,
		},
		Patch: VersionPart{
			Raw:     strconv.Itoa(patch),
			IntRepr: patch,
			IsInt:   true,
		},
	}
}

func TestParseRegularVersion(t *testing.T) {
	if !ParseVersion("1.2.3").IsFullyIdenticalTo(regular(1, 2, 3)) {
		t.Errorf("failed to parse version: %v", "1.2.3")
	}
	if !ParseVersion("1.2.0").IsFullyIdenticalTo(regular(1, 2, 0)) {
		t.Errorf("failed to parse version: %v", "1.2.0")
	}
	if !ParseVersion("1.2").IsSemanticallyIdenticalTo(regular(1, 2, 0)) {
		t.Errorf("failed to parse version: %v", "1.2")
	}
	if !ParseVersion("1").IsSemanticallyIdenticalTo(regular(1, 0, 0)) {
		t.Errorf("failed to parse version: %v", "1")
	}
}

func TestParseVersionRC(t *testing.T) {
	rc1 := ParseVersion("1.2.3-RC1")
	rc2 := ParseVersion("1.2.3-RC2")
	rc11 := ParseVersion("1.2.3-RC11")
	final := ParseVersion("1.2.3")

	// rc1 is greater or equal to rc1
	if !rc1.IsGreaterThanOrEqualTo(rc1) {
		t.Errorf("RC1 should be greater or equal to RC1")
	}

	if rc1.IsGreaterThanOrEqualTo(rc11) {
		t.Errorf("RC1 should not be greater or equal to RC11")
	}

	// Final release is greater than RC
	if !final.IsGreaterThan(rc1) {
		t.Errorf("final version should be greater than RC1")
	}
	if rc1.IsGreaterThan(final) {
		t.Errorf("RC1 should not be greater than final")
	}

	// rc2 is greater than rc1
	if !rc2.IsGreaterThan(rc1) {
		t.Errorf("RC2 should be greater than RC1")
	}
	if rc1.IsGreaterThan(rc2) {
		t.Errorf("RC1 should not be greater than RC2")
	}

	// rc11 is greater than rc1
	if !rc11.IsGreaterThan(rc1) {
		t.Errorf("RC11 should be greater than RC1")
	}
	if rc1.IsGreaterThan(rc11) {
		t.Errorf("RC1 should not be greater than RC11")
	}

}

func TestVersion_CaseSensitivity(t *testing.T) {
	rc1 := ParseVersion("1.2.3-rc1")
	RC2 := ParseVersion("1.2.3-RC1")

	if rc1.IsGreaterThan(RC2) {
		t.Errorf("rc1 should not be greater than RC1")
	}

	if RC2.IsGreaterThan(rc1) {
		t.Errorf("RC2 should not be greater than rc1")
	}

	if !rc1.IsSemanticallyIdenticalTo(RC2) {
		t.Errorf("rc1 should be semantically identical to RC1")
	}

	if rc1.IsFullyIdenticalTo(RC2) {
		t.Errorf("rc1 should not be fully identical to RC1")
	}
}

func TestParseVersionABRC(t *testing.T) {
	alpha2 := ParseVersion("1.2.3-alpha2")
	beta1 := ParseVersion("1.2.3-beta1")
	rc11 := ParseVersion("1.2.3-RC11")
	final := ParseVersion("1.2.3")

	// alpha2 is greater or equal to alpha2
	if !alpha2.IsGreaterThanOrEqualTo(alpha2) {
		t.Errorf("alpha2 should be greater or equal to alpha2")
	}

	if alpha2.IsGreaterThanOrEqualTo(beta1) {
		t.Errorf("alpha2 should not be greater or equal to beta1")
	}

	if alpha2.IsGreaterThanOrEqualTo(rc11) {
		t.Errorf("alpha2 should not be greater or equal to RC11")
	}

	// Final release is greater than alpha
	if !final.IsGreaterThan(alpha2) {
		t.Errorf("final version should be greater than alpha2")
	}

	if alpha2.IsGreaterThan(final) {
		t.Errorf("alpha2 should not be greater than final")
	}

	// beta1 is greater than alpha2
	if !beta1.IsGreaterThan(alpha2) {
		t.Errorf("beta1 should be greater than alpha2")
	}

	if alpha2.IsGreaterThan(beta1) {
		t.Errorf("alpha2 should not be greater than beta1")
	}

	// rc11 is greater than alpha2
	if !rc11.IsGreaterThan(alpha2) {
		t.Errorf("RC11 should be greater than alpha2")
	}

}

func TestParseVersionMixedLiteralAndNumeric(t *testing.T) {
	alpha2 := ParseVersion("1.2.3-alpha2")
	beta1 := ParseVersion("1.2.3-beta1")

	if alpha2.IsGreaterThanOrEqualTo(beta1) {
		t.Errorf("alpha2 should not be greater or equal to beta1")
	}

}

func TestSortByVersion(t *testing.T) {
	unordered := []string{
		"1.2.3-alpha2",
		"1.2.3-RC11",
		"1.2.3",
		"1.2.3-beta1",
		"1.2.1",
		"1.2.2",
		"1.2.4",
		"1.1.1",
		"1.3.1",
		"1.2.3-RC1",
	}

	orderedRef := []string{
		"1.1.1",
		"1.2.1",
		"1.2.2",
		"1.2.3-alpha2",
		"1.2.3-beta1",
		"1.2.3-RC1",
		"1.2.3-RC11",
		"1.2.3",
		"1.2.4",
		"1.3.1",
	}

	sorted := SortByVersion(unordered, func(s string) Version {
		return ParseVersion(s)
	})

	for i, item := range sorted {
		fmt.Printf("%v\n", item)
		if item != orderedRef[i] {
			t.Errorf("failed to sort versions")
		}
	}

}

func TestFindLatestVersionBy(t *testing.T) {
	versions := []string{
		"1.2.3-alpha2",
		"1.2.3-RC11",
		"1.2.3",
		"1.2.3-beta1",
		"1.2.1",
		"1.2.3-RC1",
	}

	latest := FindLatestVersionBy(versions, func(s string) Version {
		return ParseVersion(s)
	})

	if *latest != "1.2.3" {
		t.Errorf("failed to find the latest version")
	}
}
