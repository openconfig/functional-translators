package registrar

import (
	"testing"
)

func TestFTMetadataConsistency(t *testing.T) {
	// All functional translators must have "some metadata" and in all provided metadata, EITHER
	// SoftwareVersion or SoftwareVersionRange can be set, but not both.
	for _, ft := range FunctionalTranslatorRegistry {
		if len(ft.Metadata) == 0 {
			t.Errorf("Functional translator %s has no metadata.", ft.ID)
		}
		for _, m := range ft.Metadata {
			if m.SoftwareVersion != "" && m.SoftwareVersionRange != nil {
				t.Errorf("Functional translator %s has metadata %v with both SoftwareVersion and SoftwareVersionRange set.", ft.ID, m)
			}
			// Confirm min < max in software ranges; so, range is non-empty.
			swRange := m.SoftwareVersionRange
			if swRange != nil {
				if !swRange.Contains(m.SoftwareVersionRange.InclusiveMin) {
					t.Errorf("Functional translator %s has metadata %v with min version %s not contained in range %v", ft.ID, m, m.SoftwareVersionRange.InclusiveMin, m.SoftwareVersionRange)
				}
				if swRange.Contains(m.SoftwareVersionRange.ExclusiveMax) {
					t.Errorf("Functional translator %s has metadata %v with max version %s contained in range %v", ft.ID, m, m.SoftwareVersionRange.ExclusiveMax, m.SoftwareVersionRange)
				}
			}
		}
	}
}
