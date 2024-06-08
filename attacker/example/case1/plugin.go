package main

import (
	"fmt"
	"github.com/golang/groupcache/lru"
	ethpb "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	attaggregation "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1/attestation/aggregation/attestations"
	log "github.com/sirupsen/logrus"
	"github.com/tsinghua-cel/attacker-service/common"
	"github.com/tsinghua-cel/attacker-service/plugins"
	"github.com/tsinghua-cel/attacker-service/types"
	"strconv"
	"time"
)

type PluginCaseV1 struct {
	blockCacheContent *lru.Cache
}

func NewPluginCaseV1() plugins.AttackerPlugin {
	return &PluginCaseV1{
		blockCacheContent: lru.New(1000),
	}
}

var _ plugins.AttackerPlugin = &PluginCaseV1{}

func (c *PluginCaseV1) AttestBeforeBroadCast(ctx plugins.PluginContext, slot uint64) plugins.PluginResponse {
	return plugins.PluginResponse{
		Cmd: types.CMD_NULL,
	}
}

func (c *PluginCaseV1) AttestAfterBroadCast(ctx plugins.PluginContext, slot uint64) plugins.PluginResponse {
	return plugins.PluginResponse{
		Cmd: types.CMD_NULL,
	}
}

func (c *PluginCaseV1) AttestBeforeSign(ctx plugins.PluginContext, slot uint64, pubkey string, attestData *ethpb.AttestationData) plugins.PluginResponse {
	return plugins.PluginResponse{
		Cmd:    types.CMD_NULL,
		Result: attestData,
	}
}

func (c *PluginCaseV1) AttestAfterSign(ctx plugins.PluginContext, slot uint64, pubkey string, attest *ethpb.Attestation) plugins.PluginResponse {
	if role := ctx.Backend.GetValidatorRoleByPubkey(int(slot), pubkey); role == types.NormalRole {
		return plugins.PluginResponse{
			Cmd:    types.CMD_NULL,
			Result: attest,
		}
	}
	logger := ctx.Logger
	logger.WithFields(log.Fields{
		"slot":   slot,
		"pubkey": pubkey,
	}).Debug("receive and cache signed attest")

	ctx.Backend.AddSignedAttestation(slot, pubkey, attest)

	return plugins.PluginResponse{
		Cmd:    types.CMD_NULL,
		Result: attest,
	}
}

func (c *PluginCaseV1) AttestBeforePropose(ctx plugins.PluginContext, slot uint64, pubkey string, attest *ethpb.Attestation) plugins.PluginResponse {
	isAttacker := false
	if ctx.Backend.GetValidatorRoleByPubkey(int(slot), pubkey) == types.AttackerRole {
		isAttacker = true
	}

	if isAttacker {
		// don't propose attestation for attacker.
		log.WithFields(log.Fields{}).Debug("this is attacker, not broadcast attest")
		return plugins.PluginResponse{
			Cmd:    types.CMD_RETURN,
			Result: attest,
		}
	} else {
		// do nothing.
		return plugins.PluginResponse{
			Cmd:    types.CMD_NULL,
			Result: attest,
		}
	}
}

func (c *PluginCaseV1) AttestAfterPropose(ctx plugins.PluginContext, slot uint64, pubkey string, attest *ethpb.Attestation) plugins.PluginResponse {
	return plugins.PluginResponse{
		Cmd:    types.CMD_NULL,
		Result: attest,
	}
}

func (c *PluginCaseV1) BlockDelayForBroadCast(ctx plugins.PluginContext) plugins.PluginResponse {
	bs := ctx.Backend.GetStrategy().Block
	if !bs.DelayEnable {
		return plugins.PluginResponse{
			Cmd: types.CMD_NULL,
		}
	}
	time.Sleep(time.Millisecond * time.Duration(bs.BroadCastDelay))
	return plugins.PluginResponse{
		Cmd: types.CMD_NULL,
	}
}

func (c *PluginCaseV1) BlockDelayForReceiveBlock(ctx plugins.PluginContext, slot uint64) plugins.PluginResponse {
	backend := ctx.Backend
	ret := plugins.PluginResponse{
		Cmd: types.CMD_NULL,
	}
	valIdx, err := backend.GetValidatorByProposeSlot(slot)
	if err != nil {
		return ret
	}
	{
		valRole := backend.GetValidatorRole(int(slot), valIdx)
		if valRole != types.AttackerRole {
			return ret
		}

		duties, err := backend.GetCurrentEpochProposeDuties()
		if err != nil {
			return ret
		}

		latestAttackerVal := int64(-1)
		for _, duty := range duties {
			dutySlot, _ := strconv.Atoi(duty.Slot)
			dutyValIdx, _ := strconv.Atoi(duty.ValidatorIndex)
			if backend.GetValidatorRole(dutySlot, dutyValIdx) == types.AttackerRole {
				latestAttackerVal = int64(dutyValIdx)
			}
		}
		if valIdx != int(latestAttackerVal) {
			// 不是最后一个出块的恶意节点，不出块
			ret.Cmd = types.CMD_RETURN
			return ret
		}
	}
	// 当前是最后一个出块的恶意节点，进行延时

	epochSlots := backend.GetSlotsPerEpoch()
	seconds := backend.GetIntervalPerSlot()
	delay := (epochSlots - int(slot%uint64(epochSlots))) * seconds
	time.Sleep(time.Second * time.Duration(delay))
	key := fmt.Sprintf("delay_%d_%d", slot, valIdx)
	c.blockCacheContent.Add(key, delay)
	ctx.Logger.WithFields(log.Fields{
		"slot":     slot,
		"validx":   valIdx,
		"duration": delay,
	}).Info("delay for receive block")

	return ret
}

func (c *PluginCaseV1) BlockBeforeBroadCast(ctx plugins.PluginContext, slot uint64) plugins.PluginResponse {
	backend := ctx.Backend
	ret := plugins.PluginResponse{
		Cmd: types.CMD_NULL,
	}
	valIdx, err := backend.GetValidatorByProposeSlot(slot)
	if err != nil {
		return ret
	}
	{
		valRole := backend.GetValidatorRole(int(slot), valIdx)
		//val := s.b.GetValidatorDataSet().GetValidatorByIndex(valIdx)
		if valRole != types.AttackerRole {
			return ret
		}

		duties, err := backend.GetCurrentEpochProposeDuties()
		if err != nil {
			return ret
		}

		latestAttackerVal := int64(-1)
		for _, duty := range duties {
			dutySlot, _ := strconv.Atoi(duty.Slot)
			dutyValIdx, _ := strconv.Atoi(duty.ValidatorIndex)
			if backend.GetValidatorRole(dutySlot, dutyValIdx) == types.AttackerRole {
				latestAttackerVal = int64(dutyValIdx)
			}
		}
		if valIdx != int(latestAttackerVal) {
			// 不是最后一个出块的恶意节点，不出块
			ret.Cmd = types.CMD_RETURN
			return ret
		}
	}
	// 当前是最后一个出块的恶意节点，进行延时
	seconds := backend.GetIntervalPerSlot()
	n2delay := 12 * seconds
	total := n2delay
	time.Sleep(time.Second * time.Duration(total))

	ctx.Logger.WithFields(log.Fields{
		"slot":     slot,
		"validx":   valIdx,
		"duration": total,
	}).Info("delay for beforeBroadcastBlock")

	return ret
}

func (c *PluginCaseV1) BlockAfterBroadCast(ctx plugins.PluginContext, slot uint64) plugins.PluginResponse {
	return plugins.PluginResponse{
		Cmd: types.CMD_NULL,
	}
}

func (c *PluginCaseV1) BlockBeforeSign(ctx plugins.PluginContext, slot uint64, pubkey string, block *ethpb.SignedBeaconBlockCapella) plugins.PluginResponse {
	backend := ctx.Backend
	ret := plugins.PluginResponse{
		Cmd:    types.CMD_NULL,
		Result: block,
	}
	slotTool := common.SlotTool{backend.SlotsPerEpoch()}
	// 1. 只有每个epoch最后一个出块的恶意节点出块，其他节点不出快
	valIdx, err := backend.GetValidatorByProposeSlot(slot)
	if err != nil {
		val := backend.GetValidatorDataSet().GetValidatorByPubkey(pubkey)
		if val == nil {
			return ret
		}
		valIdx = int(val.Index)
	}
	role := backend.GetValidatorRole(int(slot), valIdx)
	log.WithFields(log.Fields{
		"slot":   slot,
		"valIdx": valIdx,
		"role":   role,
	}).Info("in modify block, get validator by propose slot")

	if role != types.AttackerRole {
		return ret
	}
	epoch := slotTool.SlotToEpoch(int64(slot))

	duties, err := backend.GetProposeDuties(int(epoch))
	if err != nil {
		return ret
	}

	latestSlotWithAttacker := int64(-1)
	for _, duty := range duties {
		dutySlot, _ := strconv.ParseInt(duty.Slot, 10, 64)
		dutyValIdx, _ := strconv.Atoi(duty.ValidatorIndex)
		if backend.GetValidatorRole(int(slot), dutyValIdx) == types.AttackerRole && dutySlot > latestSlotWithAttacker {
			latestSlotWithAttacker = dutySlot
		}
	}
	log.WithFields(log.Fields{
		"slot":               slot,
		"latestAttackerSlot": latestSlotWithAttacker,
	}).Info("modify block")

	if slot != uint64(latestSlotWithAttacker) {
		// 不是最后一个出块的恶意节点，不出块
		ret.Cmd = types.CMD_RETURN
		return ret
	}

	// 3.出的块的一个字段attestation要包含其他恶意节点的attestation。
	startEpoch := slotTool.EpochStart(epoch)
	endEpoch := slotTool.EpochEnd(epoch)
	attackerAttestations := make([]*ethpb.Attestation, 0)
	validatorSet := backend.GetValidatorDataSet()

	for i := startEpoch; i <= endEpoch; i++ {
		allSlotAttest := backend.GetAttestSet(uint64(i))
		if allSlotAttest == nil {
			continue
		}

		for publicKey, att := range allSlotAttest.Attestations {
			val := validatorSet.GetValidatorByPubkey(publicKey)
			valRole := backend.GetValidatorRole(int(i), int(val.Index))
			if val != nil && valRole == types.AttackerRole {
				log.WithField("pubkey", publicKey).Debug("add attacker attestation to block")
				attackerAttestations = append(attackerAttestations, att)
			}
		}
	}

	allAtt := append(block.Block.Body.Attestations, attackerAttestations...)
	{
		// Remove duplicates from both aggregated/unaggregated attestations. This
		// prevents inefficient aggregates being created.
		atts, _ := types.ProposerAtts(allAtt).Dedup()
		attsByDataRoot := make(map[[32]byte][]*ethpb.Attestation, len(atts))
		for _, att := range atts {
			attDataRoot, err := att.Data.HashTreeRoot()
			if err != nil {
			}
			attsByDataRoot[attDataRoot] = append(attsByDataRoot[attDataRoot], att)
		}

		attsForInclusion := types.ProposerAtts(make([]*ethpb.Attestation, 0))
		for _, as := range attsByDataRoot {
			as, err := attaggregation.Aggregate(as)
			if err != nil {
				continue
			}
			attsForInclusion = append(attsForInclusion, as...)
		}
		deduped, _ := attsForInclusion.Dedup()
		sorted, err := deduped.SortByProfitability()
		if err != nil {
			log.WithError(err).Error("sort attestation failed")
		} else {
			atts = sorted.LimitToMaxAttestations()
		}
		allAtt = atts
	}

	block.Block.Body.Attestations = allAtt

	return plugins.PluginResponse{
		Cmd:    types.CMD_NULL,
		Result: block,
	}
}

func (c *PluginCaseV1) BlockAfterSign(ctx plugins.PluginContext, slot uint64, pubkey string, block *ethpb.SignedBeaconBlockCapella) plugins.PluginResponse {
	// 不是最后一个恶意的出块，不出块
	backend := ctx.Backend
	ret := plugins.PluginResponse{
		Cmd:    types.CMD_NULL,
		Result: block,
	}
	valIdx, err := backend.GetValidatorByProposeSlot(slot)
	if err != nil {
		val := backend.GetValidatorDataSet().GetValidatorByPubkey(pubkey)
		if val == nil {
			return ret
		}
		valIdx = int(val.Index)
	}
	role := backend.GetValidatorRole(int(slot), valIdx)
	ctx.Logger.WithFields(log.Fields{
		"slot":   slot,
		"valIdx": valIdx,
		"role":   role,
	}).Info("in AfterSign, get validator by propose slot")

	if role != types.AttackerRole {
		return ret
	}
	slotTool := common.SlotTool{
		SlotsPerEpoch: backend.GetSlotsPerEpoch(),
	}
	epoch := slotTool.SlotToEpoch(int64(slot))

	duties, err := backend.GetProposeDuties(int(epoch))
	if err != nil {
		return ret
	}

	latestSlotWithAttacker := int64(-1)
	for _, duty := range duties {
		dutySlot, _ := strconv.ParseInt(duty.Slot, 10, 64)
		dutyValIdx, _ := strconv.Atoi(duty.ValidatorIndex)
		ctx.Logger.WithFields(log.Fields{
			"slot":   dutySlot,
			"valIdx": dutyValIdx,
		}).Debug("duty slot")

		if backend.GetValidatorRole(int(slot), dutyValIdx) == types.AttackerRole && dutySlot > latestSlotWithAttacker {
			latestSlotWithAttacker = dutySlot
		}
	}

	if slot != uint64(latestSlotWithAttacker) {
		// 不是最后一个恶意的出块，不出块
		ret.Cmd = types.CMD_RETURN
		return ret
	}

	return ret
}

func (c *PluginCaseV1) BlockBeforePropose(ctx plugins.PluginContext, slot uint64, pubkey string, block *ethpb.SignedBeaconBlockCapella) plugins.PluginResponse {
	return plugins.PluginResponse{
		Cmd:    types.CMD_NULL,
		Result: block,
	}
}

func (c *PluginCaseV1) BlockAfterPropose(ctx plugins.PluginContext, slot uint64, pubkey string, block *ethpb.SignedBeaconBlockCapella) plugins.PluginResponse {
	return plugins.PluginResponse{
		Cmd:    types.CMD_NULL,
		Result: block,
	}
}
