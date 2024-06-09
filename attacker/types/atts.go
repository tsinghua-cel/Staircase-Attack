package types

import (
	"context"
	"github.com/prysmaticlabs/go-bitfield"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/core/blocks"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/state"
	"github.com/prysmaticlabs/prysm/v4/config/params"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	ethpb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1/attestation/aggregation"
	"sort"
)

type ProposerAtts []*ethpb.Attestation

// Filter separates attestation list into two groups: valid and invalid attestations.
// The first group passes the all the required checks for attestation to be considered for proposing.
// And attestations from the second group should be deleted.
func (a ProposerAtts) Filter(ctx context.Context, st state.BeaconState) (ProposerAtts, ProposerAtts) {
	validAtts := make([]*ethpb.Attestation, 0, len(a))
	invalidAtts := make([]*ethpb.Attestation, 0, len(a))

	for _, att := range a {
		if err := blocks.VerifyAttestationNoVerifySignature(ctx, st, att); err == nil {
			validAtts = append(validAtts, att)
			continue
		}
		invalidAtts = append(invalidAtts, att)
	}
	return validAtts, invalidAtts
}

// SortByProfitability orders attestations by highest slot and by highest aggregation bit count.
func (a ProposerAtts) SortByProfitability() (ProposerAtts, error) {
	if len(a) < 2 {
		return a, nil
	}
	return a.SortByProfitabilityUsingMaxCover()
}

// SortByProfitabilityUsingMaxCover orders attestations by highest slot and by highest aggregation bit count.
// Duplicate bits are counted only once, using max-cover algorithm.
func (a ProposerAtts) SortByProfitabilityUsingMaxCover() (ProposerAtts, error) {
	// Separate attestations by slot, as slot number takes higher precedence when sorting.
	var slots []primitives.Slot
	attsBySlot := map[primitives.Slot]ProposerAtts{}
	for _, att := range a {
		if _, ok := attsBySlot[att.Data.Slot]; !ok {
			slots = append(slots, att.Data.Slot)
		}
		attsBySlot[att.Data.Slot] = append(attsBySlot[att.Data.Slot], att)
	}

	selectAtts := func(atts ProposerAtts) (ProposerAtts, error) {
		if len(atts) < 2 {
			return atts, nil
		}
		candidates := make([]*bitfield.Bitlist64, len(atts))
		for i := 0; i < len(atts); i++ {
			var err error
			candidates[i], err = atts[i].AggregationBits.ToBitlist64()
			if err != nil {
				return nil, err
			}
		}
		// Add selected candidates on top, those that are not selected - append at bottom.
		selectedKeys, _, err := aggregation.MaxCover(candidates, len(candidates), true /* allowOverlaps */)
		if err == nil {
			// Pick selected attestations first, leftover attestations will be appended at the end.
			// Both lists will be sorted by number of bits set.
			selectedAtts := make(ProposerAtts, selectedKeys.Count())
			leftoverAtts := make(ProposerAtts, selectedKeys.Not().Count())
			for i, key := range selectedKeys.BitIndices() {
				selectedAtts[i] = atts[key]
			}
			for i, key := range selectedKeys.Not().BitIndices() {
				leftoverAtts[i] = atts[key]
			}
			sort.Slice(selectedAtts, func(i, j int) bool {
				return selectedAtts[i].AggregationBits.Count() > selectedAtts[j].AggregationBits.Count()
			})
			sort.Slice(leftoverAtts, func(i, j int) bool {
				return leftoverAtts[i].AggregationBits.Count() > leftoverAtts[j].AggregationBits.Count()
			})
			return append(selectedAtts, leftoverAtts...), nil
		}
		return atts, nil
	}

	// Select attestations. Slots are sorted from higher to lower at this point. Within slots attestations
	// are sorted to maximize profitability (greedily selected, with previous attestations' bits
	// evaluated before including any new attestation).
	var sortedAtts ProposerAtts
	sort.Slice(slots, func(i, j int) bool {
		return slots[i] > slots[j]
	})
	for _, slot := range slots {
		selected, err := selectAtts(attsBySlot[slot])
		if err != nil {
			return nil, err
		}
		sortedAtts = append(sortedAtts, selected...)
	}

	return sortedAtts, nil
}

// LimitToMaxAttestations limits attestations to maximum attestations per block.
func (a ProposerAtts) LimitToMaxAttestations() ProposerAtts {
	if uint64(len(a)) > params.BeaconConfig().MaxAttestations {
		return a[:params.BeaconConfig().MaxAttestations]
	}
	return a
}

// Dedup removes duplicate attestations (ones with the same bits set on).
// Important: not only exact duplicates are removed, but proper subsets are removed too
// (their known bits are redundant and are already contained in their supersets).
func (a ProposerAtts) Dedup() (ProposerAtts, error) {
	if len(a) < 2 {
		return a, nil
	}
	attsByDataRoot := make(map[[32]byte][]*ethpb.Attestation, len(a))
	for _, att := range a {
		attDataRoot, err := att.Data.HashTreeRoot()
		if err != nil {
			continue
		}
		attsByDataRoot[attDataRoot] = append(attsByDataRoot[attDataRoot], att)
	}

	uniqAtts := make([]*ethpb.Attestation, 0, len(a))
	for _, atts := range attsByDataRoot {
		for i := 0; i < len(atts); i++ {
			a := atts[i]
			for j := i + 1; j < len(atts); j++ {
				b := atts[j]
				if c, err := a.AggregationBits.Contains(b.AggregationBits); err != nil {
					return nil, err
				} else if c {
					// a contains b, b is redundant.
					atts[j] = atts[len(atts)-1]
					atts[len(atts)-1] = nil
					atts = atts[:len(atts)-1]
					j--
				} else if c, err := b.AggregationBits.Contains(a.AggregationBits); err != nil {
					return nil, err
				} else if c {
					// b contains a, a is redundant.
					atts[i] = atts[len(atts)-1]
					atts[len(atts)-1] = nil
					atts = atts[:len(atts)-1]
					i--
					break
				}
			}
		}
		uniqAtts = append(uniqAtts, atts...)
	}

	return uniqAtts, nil
}
