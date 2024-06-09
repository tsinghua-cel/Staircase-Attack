package apis

// RoleAPI offers and API for role operations.
type AdminAPI struct {
	b Backend
}

// NewRoleAPI creates a new tx pool service that gives information about the transaction pool.
func NewAdminAPI(b Backend) *AdminAPI {
	return &AdminAPI{b}
}

func (s *AdminAPI) SetRoleAttacker(valIndex int) {
	//valSet := s.b.GetValidatorDataSet()
	//valSet.SetValidatorRole(valIndex, types.AttackerRole)
}

func (s *AdminAPI) SetRoleNormal(valIndex int) {
	//valSet := s.b.GetValidatorDataSet()
	//valSet.SetValidatorRole(valIndex, types.NormalRole)
}
