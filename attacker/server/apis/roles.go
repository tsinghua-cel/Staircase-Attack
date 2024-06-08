package apis

import "github.com/tsinghua-cel/attacker-service/plugins"

// RoleAPI offers and API for role operations.
type AdminAPI struct {
	b      Backend
	plugin plugins.AttackerPlugin
}

// NewRoleAPI creates a new tx pool service that gives information about the transaction pool.
func NewAdminAPI(b Backend, plugin plugins.AttackerPlugin) *AdminAPI {
	return &AdminAPI{b, plugin}
}

func (s *AdminAPI) SetRoleAttacker(valIndex int) {
	//valSet := s.b.GetValidatorDataSet()
	//valSet.SetValidatorRole(valIndex, types.AttackerRole)
}

func (s *AdminAPI) SetRoleNormal(valIndex int) {
	//valSet := s.b.GetValidatorDataSet()
	//valSet.SetValidatorRole(valIndex, types.NormalRole)
}
