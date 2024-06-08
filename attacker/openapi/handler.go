package openapi

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/tsinghua-cel/attacker-service/dbmodel"
	"github.com/tsinghua-cel/attacker-service/types"
	"net/http"
	"strconv"
)

type apiHandler struct {
	backend types.ServiceBackend
}

// @Summary Get duties by epoch
// @Description get duties by epoch
// @ID get-duties-by-epoch
// @Accept  json
// @Produce  json
// @Param epoch path int true "Epoch"
// @Success 200 {array} types.ProposerDuty
// @Router /duties/{epoch} [get]
func (api apiHandler) GetDutiesByEpoch(c *gin.Context) {
	param := c.Param("epoch")
	epoch, _ := strconv.Atoi(param)
	duties, err := api.backend.GetProposeDuties(epoch)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, duties)
}

// @Summary Get strategy config
// @Description get strategy config
// @ID get-strategy
// @Accept  json
// @Produce  json
// @Success 200 {object} types.Strategy
// @Router /strategy [get]
func (api apiHandler) GetStrategy(c *gin.Context) {
	strategy := api.backend.GetStrategy()
	c.JSON(200, strategy)
}

// @Summary Get reward by epoch
// @Description get reward by epoch
// @ID get-reward-by-epoch
// @Accept  json
// @Produce  json
// @Param epoch path int true "Epoch"
// @Success 200 {array} dbmodel.BlockReward
// @Router /reward/{epoch} [get]
func (api apiHandler) GetRewardByEpoch(c *gin.Context) {
	param := c.Param("epoch")
	epoch, _ := strconv.Atoi(param)
	list := dbmodel.GetRewardListByEpoch(int64(epoch))
	c.JSON(200, list)
}

// @Summary Update strategy
// @Description update strategy
// @ID update-strategy
// @Accept  json
// @Produce  json
// @Param strategy body types.Strategy true "Strategy"
// @Success 200 {string} string
// @Router /update-strategy [post]
func (api apiHandler) UpdateStrategy(c *gin.Context) {
	var req types.Strategy
	err := c.ShouldBindJSON(&req) // 解析req参数
	if err != nil {
		log.WithError(err).Println("UpdateStrategy ctx.ShouldBindJSON error")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	err = api.backend.UpdateStrategy(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, "ok")
}

// @Summary Get reorgs
// @Description get reorgs
// @ID get-reorgs
// @Accept  json
// @Produce  json
// @Success 200 {array} types.ReorgEvent
// @Router /reorgs [get]
func (api apiHandler) GetReorgs(c *gin.Context) {
	reorgs := dbmodel.GetAllReorgList()
	slots := make(map[int64]bool)
	var res []types.ReorgEvent
	for _, reorg := range reorgs {
		if _, ok := slots[reorg.Slot]; ok {
			continue
		}
		res = append(res, types.ReorgEvent{
			Epoch:                 reorg.Epoch,
			Slot:                  reorg.Slot,
			Depth:                 int64(reorg.Depth),
			OldBlockSlot:          reorg.OldBlockSlot,
			NewBlockSlot:          reorg.NewBlockSlot,
			OldBlockProposerIndex: reorg.OldBlockProposerIndex,
			NewBlockProposerIndex: reorg.NewBlockProposerIndex,
			OldHeadState:          reorg.OldHeadState,
			NewHeadState:          reorg.NewHeadState,
		})
	}
	c.JSON(200, res)
}
