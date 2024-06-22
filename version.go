package version

import (
	"fmt"
	"github.com/samber/lo"
	"log/slog"
	"strconv"
	"strings"
	"unicode"
)

type Version struct {
	Major VersionPart
	Minor VersionPart
	Patch VersionPart
	Parts []VersionPart
}

func (v Version) IsNonEmpty() bool {
	return len(v.Parts) > 0
}

func (v Version) IsSemverRelease() bool {
	return v.IsNonEmpty() && lo.EveryBy(v.Parts, func(part VersionPart) bool {
		return part.IsInt
	})
}

// IsValid returns true if the version is non-empty and the major part is an int and >= 0
func (v Version) IsValid() bool {
	return v.IsNonEmpty() && v.Major.IsInt && v.Major.IntRepr >= 0
}

type VersionPart struct {
	Raw      string
	IntRepr  int
	IsInt    bool
	SubParts []VersionPart
}

func NewVersion(major int, minor int, patch int) Version {
	res := Version{
		Major: VersionPart{Raw: strconv.Itoa(major), IntRepr: major, IsInt: true},
		Minor: VersionPart{Raw: strconv.Itoa(minor), IntRepr: minor, IsInt: true},
		Patch: VersionPart{Raw: strconv.Itoa(patch), IntRepr: patch, IsInt: true},
	}

	res.Parts = []VersionPart{res.Major, res.Minor, res.Patch}

	return res
}

func ParseVersionPart(part string) VersionPart {

	finalizeNewPart := func(currentPartRaw string, currentIsInt bool) VersionPart {
		if currentIsInt {
			intRepr, err := strconv.Atoi(currentPartRaw)
			if err != nil {
				slog.Error(fmt.Sprintf("failed to parse int from string: %s", currentPartRaw))
				return VersionPart{Raw: currentPartRaw, IsInt: false}
			}
			return VersionPart{Raw: currentPartRaw, IntRepr: intRepr, IsInt: true}
		} else {
			return VersionPart{Raw: currentPartRaw, IsInt: false}
		}
	}

	intRepr, err := strconv.Atoi(part)
	if err != nil {
		parts := make([]VersionPart, 0, 2)
		currentPartRaw := ""
		currentIsInt := false
		for _, c := range part {

			if unicode.IsDigit(c) {
				if currentIsInt || currentPartRaw == "" {
					currentPartRaw += string(c)
					currentIsInt = true
				} else {
					// we are switching from a string part to an int part
					if len(currentPartRaw) > 0 {
						parts = append(parts, finalizeNewPart(currentPartRaw, currentIsInt))
					}
					currentPartRaw = string(c)
					currentIsInt = true
				}
			} else {
				if !currentIsInt || currentPartRaw == "" {
					currentPartRaw += string(c)
					currentIsInt = false
				} else {
					// we are switching from an int part to a string part
					if len(currentPartRaw) > 0 {
						parts = append(parts, finalizeNewPart(currentPartRaw, currentIsInt))
					}
					currentPartRaw = string(c)
					currentIsInt = false
				}
			}
		}

		if len(parts) > 0 {
			parts = append(parts, finalizeNewPart(currentPartRaw, currentIsInt))
		}

		return VersionPart{Raw: part, IsInt: false, SubParts: parts}
	} else {
		return VersionPart{Raw: part, IntRepr: intRepr, IsInt: true}
	}
}

func (p VersionPart) IsGreaterOrEqThan(other VersionPart) bool {
	return p.IsGreaterThan(other) || (!p.IsGreaterThan(other) && !other.IsGreaterThan(p))
}

func (p VersionPart) IsSemanticallyEqualTo(other VersionPart) bool {
	return !p.IsGreaterThan(other) && !other.IsGreaterThan(p)
}

func (p VersionPart) IsLessThan(other VersionPart) bool {
	return !p.IsGreaterOrEqThan(other)
}

func (p VersionPart) IsGreaterThan(other VersionPart) bool {
	if p.IsInt && other.IsInt {
		// Both are regular semantic parts. Just compare the ints
		return p.IntRepr > other.IntRepr
	} else if p.IsInt && !other.IsInt {
		// if the other starts with an int part that is lower or same, we are greater (beta, rc, etc)
		if len(other.SubParts) > 0 && other.SubParts[0].IsInt {
			return p.IntRepr >= other.SubParts[0].IntRepr
		} else {
			return true
		}
	} else if !p.IsInt && !other.IsInt {

		if len(p.SubParts) == 0 && len(other.SubParts) == 0 {
			// treat as literal, sort alphabetically, case insensitively
			return strings.ToLower(p.Raw) > strings.ToLower(other.Raw)
		}

		if len(p.SubParts) == 0 {
			return false
		}

		if len(other.SubParts) == 0 {
			return true
		}

		if len(p.SubParts) == len(other.SubParts) {
			for i := range p.SubParts {
				if p.SubParts[i].IsGreaterThan(other.SubParts[i]) {
					return true
				}
				if other.SubParts[i].IsGreaterThan(p.SubParts[i]) {
					return false
				}
				if i+1 == len(p.SubParts) {
					return p.SubParts[i].IsGreaterThan(other.SubParts[i])
				}
			}
		} else if len(p.SubParts) < len(other.SubParts) {
			// up to the equal count, all parts are equal or newer, then our version should be newer
			for i := range p.SubParts {
				if other.SubParts[i].IsGreaterThan(p.SubParts[i]) {
					return false
				}
			}
			return true
		} else if len(p.SubParts) > len(other.SubParts) {
			// up to the equal count, all parts are equal or newer, then our version should be newer
			for i := range other.SubParts {
				if p.SubParts[i].IsGreaterThan(other.SubParts[i]) {
					return true
				}
			}
			return false
		}
	} else {
		return !other.IsGreaterOrEqThan(p)
	}

	return p.Raw > other.Raw
}

func (v Version) String() string {
	return fmt.Sprintf("%s.%s.%s", v.Major.Raw, v.Minor.Raw, v.Patch.Raw)
}

func (v Version) Raw() string {
	result := ""
	for i, part := range v.Parts {
		if i > 0 {
			result += "."
		}
		result += part.Raw
	}
	return result
}

func (v Version) IsGreaterThan(other Version) bool {
	if v.Major.IsGreaterThan(other.Major) {
		return true
	} else if v.Major.IsLessThan(other.Major) {
		return false
	}

	if v.Minor.IsGreaterThan(other.Minor) {
		return true
	} else if v.Minor.IsLessThan(other.Minor) {
		return false
	}

	// TODO: Support longer sequences than major.minor.patch

	return v.Patch.IsGreaterThan(other.Patch)
}

func (v Version) IsGreaterThanOrEqualTo(other Version) bool {
	return v.IsGreaterThan(other) || v.IsSemanticallyIdenticalTo(other)
}

func (a Version) IsSemanticallyIdenticalTo(b Version) bool {
	return !a.IsGreaterThan(b) && !b.IsGreaterThan(a)
}

func (a Version) IsFullyIdenticalTo(b Version) bool {
	return a.Major.IsFullyIdenticalTo(b.Major) &&
		a.Minor.IsFullyIdenticalTo(b.Minor) &&
		a.Patch.IsFullyIdenticalTo(b.Patch)
}

func (v Version) IsLessThan(version Version) bool {
	return !v.IsGreaterThanOrEqualTo(version)
}

func (a VersionPart) IsFullyIdenticalTo(b VersionPart) bool {
	return a.Raw == b.Raw &&
		a.IntRepr == b.IntRepr &&
		a.IsInt == b.IsInt &&
		len(a.SubParts) == len(b.SubParts) &&
		lo.EveryBy(lo.Zip2(a.SubParts, b.SubParts), func(pair lo.Tuple2[VersionPart, VersionPart]) bool {
			return pair.A.IsFullyIdenticalTo(pair.B)
		})
}

func ParseVersion(versionStr string) Version {

	result := Version{}

	if strings.HasPrefix(versionStr, "v") || strings.HasPrefix(versionStr, "V") {
		versionStr = versionStr[1:]
	}

	versionPartsRaw := strings.Split(versionStr, ".")
	result.Parts = lo.Map(versionPartsRaw, func(raw string, _ int) VersionPart {
		return ParseVersionPart(raw)
	})

	if len(result.Parts) >= 1 {
		result.Major = result.Parts[0]
	} else {
		result.Major = VersionPart{Raw: "", IntRepr: 0, IsInt: true}
	}

	if len(result.Parts) >= 2 {
		result.Minor = result.Parts[1]
	} else {
		result.Minor = VersionPart{Raw: "", IntRepr: 0, IsInt: true}
	}

	if len(result.Parts) >= 3 {
		result.Patch = result.Parts[2]
	} else {
		result.Patch = VersionPart{Raw: "", IntRepr: 0, IsInt: true}
	}

	return result
}

func FindLatestVersionBy[T any](items []T, fGetVersion func(T) Version) *T {

	if len(items) == 0 {
		return nil
	}

	latest := &items[0]
	latestVersion := fGetVersion(*latest)

	for _, item := range items {
		itemVersion := fGetVersion(item)
		if itemVersion.IsGreaterThan(latestVersion) {
			latest = &item
			latestVersion = itemVersion
		}
	}

	return latest
}
