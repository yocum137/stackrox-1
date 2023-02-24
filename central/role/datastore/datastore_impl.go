package datastore

import (
	"context"

	"github.com/pkg/errors"
	rolePkg "github.com/stackrox/rox/central/role"
	"github.com/stackrox/rox/central/role/resources"
	rocksDBStore "github.com/stackrox/rox/central/role/store"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/auth/permissions"
	"github.com/stackrox/rox/pkg/declarativeconfig"
	"github.com/stackrox/rox/pkg/errox"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/postgres/pgutils"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/sync"
)

var (
	// TODO: ROX-14398 Replace Role with Access
	roleSAC = sac.ForResource(resources.Role)

	log = logging.LoggerForModule()
)

type dataStoreImpl struct {
	roleStorage          rocksDBStore.RoleStore
	permissionSetStorage rocksDBStore.PermissionSetStore
	accessScopeStorage   rocksDBStore.SimpleAccessScopeStore

	lock sync.RWMutex
}

func (ds *dataStoreImpl) UpsertRole(ctx context.Context, newRole *storage.Role) error {
	if err := sac.VerifyAuthzOK(roleSAC.WriteAllowed(ctx)); err != nil {
		return err
	}
	if err := rolePkg.ValidateRole(newRole); err != nil {
		return errors.Wrap(errox.InvalidArgs, err.Error())
	}

	ds.lock.Lock()
	defer ds.lock.Unlock()

	oldRole, exists, err := ds.roleStorage.Get(ctx, newRole.GetName())
	if err != nil {
		return err
	}
	if exists {
		if err := verifyRoleOriginMatches(ctx, oldRole); err != nil {
			return err
		}
	}
	if err := verifyRoleOriginMatches(ctx, newRole); err != nil {
		return err
	}

	permissionSet, accessScope, err := ds.verifyRoleReferencesExist(ctx, newRole)
	if err != nil {
		return err
	}
	if err := verifyPermissionSetOriginMatches(ctx, permissionSet); err != nil {
		return err
	}
	if err := verifyAccessScopeOriginMatches(ctx, accessScope); err != nil {
		return err
	}

	// Constraints ok, write the object. We expect the underlying store to
	// verify there is no role with the same name.
	if err := ds.roleStorage.Upsert(ctx, newRole); err != nil {
		return err
	}

	return nil
}

func (ds *dataStoreImpl) UpsertPermissionSet(ctx context.Context, newPS *storage.PermissionSet) error {
	if err := sac.VerifyAuthzOK(roleSAC.WriteAllowed(ctx)); err != nil {
		return err
	}
	if err := rolePkg.ValidatePermissionSet(newPS); err != nil {
		return errors.Wrap(errox.InvalidArgs, err.Error())
	}

	ds.lock.Lock()
	defer ds.lock.Unlock()

	oldPS, exists, err := ds.permissionSetStorage.Get(ctx, newPS.GetId())
	if err != nil {
		return err
	}
	if exists {
		if err := verifyPermissionSetOriginMatches(ctx, oldPS); err != nil {
			return err
		}
	}
	if err := verifyPermissionSetOriginMatches(ctx, newPS); err != nil {
		return err
	}

	// Constraints ok, write the object. We expect the underlying store to
	// verify there is no permission set with the same name.
	if err := ds.permissionSetStorage.Upsert(ctx, newPS); err != nil {
		return err
	}

	return nil
}

func (ds *dataStoreImpl) UpsertAccessScope(ctx context.Context, newScope *storage.SimpleAccessScope) error {
	if err := sac.VerifyAuthzOK(roleSAC.WriteAllowed(ctx)); err != nil {
		return err
	}
	if err := rolePkg.ValidateSimpleAccessScope(newScope); err != nil {
		return errors.Wrap(errox.InvalidArgs, err.Error())
	}

	ds.lock.Lock()
	defer ds.lock.Unlock()

	oldScope, exists, err := ds.accessScopeStorage.Get(ctx, newScope.GetId())
	if err != nil {
		return err
	}
	if exists {
		if err := verifyAccessScopeOriginMatches(ctx, oldScope); err != nil {
			return err
		}
		if err := verifyAccessScopeOriginMatches(ctx, newScope); err != nil {
			return err
		}
		return err
	}

	// Constraints ok, write the object. We expect the underlying store to
	// verify there is no access scope with the same name.
	if err := ds.accessScopeStorage.Upsert(ctx, newScope); err != nil {
		return err
	}

	return nil
}

func (ds *dataStoreImpl) GetRole(ctx context.Context, name string) (*storage.Role, bool, error) {
	if ok, err := roleSAC.ReadAllowed(ctx); !ok || err != nil {
		return nil, false, err
	}

	return ds.roleStorage.Get(ctx, name)
}

func (ds *dataStoreImpl) GetAllRoles(ctx context.Context) ([]*storage.Role, error) {
	if ok, err := roleSAC.ReadAllowed(ctx); !ok || err != nil {
		return nil, err
	}

	return ds.getAllRolesNoScopeCheck(ctx)
}

func (ds *dataStoreImpl) CountRoles(ctx context.Context) (int, error) {
	if ok, err := roleSAC.ReadAllowed(ctx); !ok || err != nil {
		return 0, err
	}

	return ds.roleStorage.Count(ctx)
}

func (ds *dataStoreImpl) getAllRolesNoScopeCheck(ctx context.Context) ([]*storage.Role, error) {
	var roles []*storage.Role
	walkFn := func() error {
		roles = roles[:0]
		return ds.roleStorage.Walk(ctx, func(role *storage.Role) error {
			roles = append(roles, role)
			return nil
		})
	}
	if err := pgutils.RetryIfPostgres(walkFn); err != nil {
		return nil, err
	}

	return roles, nil
}

func (ds *dataStoreImpl) AddRole(ctx context.Context, role *storage.Role) error {
	if err := sac.VerifyAuthzOK(roleSAC.WriteAllowed(ctx)); err != nil {
		return err
	}
	if err := rolePkg.ValidateRole(role); err != nil {
		return errors.Wrap(errox.InvalidArgs, err.Error())
	}
	if err := verifyNotDefaultRole(role); err != nil {
		return err
	}

	// protect against TOCTOU race condition
	ds.lock.Lock()
	defer ds.lock.Unlock()

	// Verify storage constraints.
	if err := ds.verifyRoleNameDoesNotExist(ctx, role.GetName()); err != nil {
		return err
	}
	if _, _, err := ds.verifyRoleReferencesExist(ctx, role); err != nil {
		return err
	}

	return ds.roleStorage.Upsert(ctx, role)
}

func (ds *dataStoreImpl) UpdateRole(ctx context.Context, role *storage.Role) error {
	if err := sac.VerifyAuthzOK(roleSAC.WriteAllowed(ctx)); err != nil {
		return err
	}
	if err := rolePkg.ValidateRole(role); err != nil {
		return errors.Wrap(errox.InvalidArgs, err.Error())
	}
	if err := verifyNotDefaultRole(role); err != nil {
		return err
	}

	// protect against TOCTOU race condition
	ds.lock.Lock()
	defer ds.lock.Unlock()

	// Verify storage constraints.
	existingRole, err := ds.verifyRoleNameExists(ctx, role.GetName())
	if err != nil {
		return err
	}
	if err = verifyRoleOriginMatches(ctx, existingRole); err != nil {
		return err
	}
	if _, _, err = ds.verifyRoleReferencesExist(ctx, role); err != nil {
		return err
	}

	return ds.roleStorage.Upsert(ctx, role)
}

func (ds *dataStoreImpl) RemoveRole(ctx context.Context, name string) error {
	if err := sac.VerifyAuthzOK(roleSAC.WriteAllowed(ctx)); err != nil {
		return err
	}

	if err := ds.verifyRoleForDeletion(ctx, name); err != nil {
		return err
	}

	return ds.roleStorage.Delete(ctx, name)
}

func verifyRoleOriginMatches(ctx context.Context, role *storage.Role) error {
	if !declarativeconfig.CanModifyResource(ctx, role) {
		return errors.Wrapf(errox.NotAuthorized, "role %q's origin is %s, cannot be modified or deleted with the current permission",
			role.GetName(), role.GetTraits().GetOrigin())
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// Permission sets                                                            //
//                                                                            //

func (ds *dataStoreImpl) GetPermissionSet(ctx context.Context, id string) (*storage.PermissionSet, bool, error) {
	if ok, err := roleSAC.ReadAllowed(ctx); !ok || err != nil {
		return nil, false, err
	}

	return ds.permissionSetStorage.Get(ctx, id)
}

func (ds *dataStoreImpl) GetAllPermissionSets(ctx context.Context) ([]*storage.PermissionSet, error) {
	if ok, err := roleSAC.ReadAllowed(ctx); !ok || err != nil {
		return nil, err
	}

	var permissionSets []*storage.PermissionSet
	walkFn := func() error {
		permissionSets = permissionSets[:0]
		return ds.permissionSetStorage.Walk(ctx, func(permissionSet *storage.PermissionSet) error {
			permissionSets = append(permissionSets, permissionSet)
			return nil
		})
	}
	if err := pgutils.RetryIfPostgres(walkFn); err != nil {
		return nil, err
	}

	return permissionSets, nil
}

func (ds *dataStoreImpl) CountPermissionSets(ctx context.Context) (int, error) {
	if ok, err := roleSAC.ReadAllowed(ctx); !ok || err != nil {
		return 0, err
	}

	return ds.permissionSetStorage.Count(ctx)
}

func (ds *dataStoreImpl) AddPermissionSet(ctx context.Context, permissionSet *storage.PermissionSet) error {
	if err := sac.VerifyAuthzOK(roleSAC.WriteAllowed(ctx)); err != nil {
		return err
	}
	if err := rolePkg.ValidatePermissionSet(permissionSet); err != nil {
		return errors.Wrap(errox.InvalidArgs, err.Error())
	}
	if err := verifyNotDefaultPermissionSet(permissionSet); err != nil {
		return err
	}

	ds.lock.Lock()
	defer ds.lock.Unlock()

	// Verify storage constraints.
	if err := ds.verifyPermissionSetIDDoesNotExist(ctx, permissionSet.GetId()); err != nil {
		return err
	}

	// Constraints ok, write the object. We expect the underlying store to
	// verify there is no permission set with the same name.
	if err := ds.permissionSetStorage.Upsert(ctx, permissionSet); err != nil {
		return err
	}

	return nil
}

func (ds *dataStoreImpl) UpdatePermissionSet(ctx context.Context, permissionSet *storage.PermissionSet) error {
	if err := sac.VerifyAuthzOK(roleSAC.WriteAllowed(ctx)); err != nil {
		return err
	}
	if err := rolePkg.ValidatePermissionSet(permissionSet); err != nil {
		return errors.Wrap(errox.InvalidArgs, err.Error())
	}
	if err := verifyNotDefaultPermissionSet(permissionSet); err != nil {
		return err
	}

	ds.lock.Lock()
	defer ds.lock.Unlock()

	// Verify storage constraints.
	existingPermissionSet, err := ds.verifyPermissionSetIDExists(ctx, permissionSet.GetId())
	if err != nil {
		return err
	}
	if err := verifyPermissionSetOriginMatches(ctx, existingPermissionSet); err != nil {
		return err
	}

	// Constraints ok, write the object. We expect the underlying store to
	// verify there is no permission set with the same name.
	if err := ds.permissionSetStorage.Upsert(ctx, permissionSet); err != nil {
		return err
	}

	return nil
}

func (ds *dataStoreImpl) RemovePermissionSet(ctx context.Context, id string) error {
	if err := sac.VerifyAuthzOK(roleSAC.WriteAllowed(ctx)); err != nil {
		return err
	}

	ds.lock.Lock()
	defer ds.lock.Unlock()

	permissionSet, found, err := ds.permissionSetStorage.Get(ctx, id)
	if err != nil {
		return err
	}
	if !found {
		return errors.Wrapf(errox.NotFound, "id = %s", id)
	}
	if err := verifyNotDefaultPermissionSet(permissionSet); err != nil {
		return err
	}
	if err := verifyPermissionSetOriginMatches(ctx, permissionSet); err != nil {
		return err
	}

	// Ensure this PermissionSet isn't in use by any Role.
	roles, err := ds.getAllRolesNoScopeCheck(ctx)
	if err != nil {
		return err
	}
	for _, role := range roles {
		if role.GetPermissionSetId() == id {
			return errors.Wrapf(errox.ReferencedByAnotherObject, "cannot delete permission set in use by role %q", role.GetName())
		}
	}

	// Constraints ok, delete the object.
	if err := ds.permissionSetStorage.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

func verifyPermissionSetOriginMatches(ctx context.Context, ps *storage.PermissionSet) error {
	if !declarativeconfig.CanModifyResource(ctx, ps) {
		return errors.Wrapf(errox.NotAuthorized, "permission set %q's origin is %s, cannot be modified or deleted with the current permission",
			ps.GetName(), ps.GetTraits().GetOrigin())
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// Access scopes                                                              //
//                                                                            //

func (ds *dataStoreImpl) GetAccessScope(ctx context.Context, id string) (*storage.SimpleAccessScope, bool, error) {
	if ok, err := roleSAC.ReadAllowed(ctx); !ok || err != nil {
		return nil, false, err
	}

	return ds.accessScopeStorage.Get(ctx, id)
}

func (ds *dataStoreImpl) GetAllAccessScopes(ctx context.Context) ([]*storage.SimpleAccessScope, error) {
	if ok, err := roleSAC.ReadAllowed(ctx); !ok || err != nil {
		return nil, err
	}

	var scopes []*storage.SimpleAccessScope
	walkFn := func() error {
		scopes = scopes[:0]
		return ds.accessScopeStorage.Walk(ctx, func(scope *storage.SimpleAccessScope) error {
			scopes = append(scopes, scope)
			return nil
		})
	}
	if err := pgutils.RetryIfPostgres(walkFn); err != nil {
		return nil, err
	}

	return scopes, nil
}

func (ds *dataStoreImpl) CountAccessScopes(ctx context.Context) (int, error) {
	if ok, err := roleSAC.ReadAllowed(ctx); !ok || err != nil {
		return 0, err
	}

	return ds.accessScopeStorage.Count(ctx)
}

func (ds *dataStoreImpl) AddAccessScope(ctx context.Context, scope *storage.SimpleAccessScope) error {
	if err := sac.VerifyAuthzOK(roleSAC.WriteAllowed(ctx)); err != nil {
		return err
	}
	if err := rolePkg.ValidateSimpleAccessScope(scope); err != nil {
		return errors.Wrap(errox.InvalidArgs, err.Error())
	}
	if err := verifyNotDefaultAccessScope(scope); err != nil {
		return err
	}

	ds.lock.Lock()
	defer ds.lock.Unlock()

	// Verify storage constraints.
	if err := ds.verifyAccessScopeIDDoesNotExist(ctx, scope.GetId()); err != nil {
		return err
	}

	// Constraints ok, write the object. We expect the underlying store to
	// verify there is no access scope with the same name.
	if err := ds.accessScopeStorage.Upsert(ctx, scope); err != nil {
		return err
	}

	return nil
}

func (ds *dataStoreImpl) UpdateAccessScope(ctx context.Context, scope *storage.SimpleAccessScope) error {
	if err := sac.VerifyAuthzOK(roleSAC.WriteAllowed(ctx)); err != nil {
		return err
	}
	if err := rolePkg.ValidateSimpleAccessScope(scope); err != nil {
		return errors.Wrap(errox.InvalidArgs, err.Error())
	}
	if err := verifyNotDefaultAccessScope(scope); err != nil {
		return err
	}

	ds.lock.Lock()
	defer ds.lock.Unlock()

	// Verify storage constraints.
	as, err := ds.verifyAccessScopeIDExists(ctx, scope.GetId())
	if err != nil {
		return err
	}
	if err := verifyAccessScopeOriginMatches(ctx, as); err != nil {
		return err
	}

	// Constraints ok, write the object. We expect the underlying store to
	// verify there is no access scope with the same name.
	if err := ds.accessScopeStorage.Upsert(ctx, scope); err != nil {
		return err
	}

	return nil
}

func (ds *dataStoreImpl) RemoveAccessScope(ctx context.Context, id string) error {
	if err := sac.VerifyAuthzOK(roleSAC.WriteAllowed(ctx)); err != nil {
		return err
	}

	ds.lock.Lock()
	defer ds.lock.Unlock()

	// Verify storage constraints.
	accessScope, found, err := ds.accessScopeStorage.Get(ctx, id)
	if err != nil {
		return err
	}
	if !found {
		return errors.Wrapf(errox.NotFound, "id = %s", id)
	}
	if err := verifyNotDefaultAccessScope(accessScope); err != nil {
		return err
	}
	if err := verifyAccessScopeOriginMatches(ctx, accessScope); err != nil {
		return err
	}

	// Ensure this AccessScope isn't in use by any Role.
	roles, err := ds.getAllRolesNoScopeCheck(ctx)
	if err != nil {
		return err
	}
	for _, role := range roles {
		if role.GetAccessScopeId() == id {
			return errors.Wrapf(errox.ReferencedByAnotherObject, "cannot delete access scope in use by role %q", role.GetName())
		}
	}

	// Constraints ok, delete the object.
	if err := ds.accessScopeStorage.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

func (ds *dataStoreImpl) GetAndResolveRole(ctx context.Context, name string) (permissions.ResolvedRole, error) {
	if ok, err := roleSAC.ReadAllowed(ctx); !ok || err != nil {
		return nil, err
	}

	ds.lock.RLock()
	defer ds.lock.RUnlock()

	// No need to continue if the role does not exist.
	role, found, err := ds.roleStorage.Get(ctx, name)
	if err != nil || !found {
		return nil, err
	}

	permissionSet, err := ds.getRolePermissionSetOrError(ctx, role)
	if err != nil {
		return nil, err
	}

	accessScope, err := ds.getRoleAccessScopeOrError(ctx, role)
	if err != nil {
		return nil, err
	}

	resolvedRole := &resolvedRoleImpl{
		role:          role,
		permissionSet: permissionSet,
		accessScope:   accessScope,
	}

	return resolvedRole, nil
}

func verifyAccessScopeOriginMatches(ctx context.Context, as *storage.SimpleAccessScope) error {
	if !declarativeconfig.CanModifyResource(ctx, as) {
		return errors.Wrapf(errox.NotAuthorized, "access scope %q's origin is %s, cannot be modified or deleted with the current permission",
			as.GetName(), as.GetTraits().GetOrigin())
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// Storage constraints                                                        //
//                                                                            //
// Uniqueness of the 'name' field is expected to be verified by the           //
// underlying store, see its `--uniq-key-func` flag                           //

func (ds *dataStoreImpl) verifyRoleReferencesExist(ctx context.Context, role *storage.Role) (*storage.PermissionSet, *storage.SimpleAccessScope, error) {
	// Verify storage constraints.
	permissionSet, err := ds.verifyPermissionSetIDExists(ctx, role.GetPermissionSetId())
	if err != nil {
		return nil, nil, errors.Wrapf(errox.InvalidArgs, "referenced permission set %s does not exist", role.GetPermissionSetId())
	}
	accessScope, err := ds.verifyAccessScopeIDExists(ctx, role.GetAccessScopeId())
	if err != nil {
		return nil, nil, errors.Wrapf(errox.InvalidArgs, "referenced access scope %s does not exist", role.GetAccessScopeId())
	}
	return permissionSet, accessScope, nil
}

// Returns errox.InvalidArgs if the given role is a default one.
func verifyNotDefaultRole(role *storage.Role) error {
	if rolePkg.IsDefaultRole(role) {
		return errors.Wrapf(errox.InvalidArgs, "default role %q cannot be modified or deleted", role.GetName())
	}
	return nil
}

// Returns errox.NotFound if there is no permission set with the supplied ID.
func (ds *dataStoreImpl) verifyPermissionSetIDExists(ctx context.Context, id string) (*storage.PermissionSet, error) {
	ps, found, err := ds.permissionSetStorage.Get(ctx, id)

	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errors.Wrapf(errox.NotFound, "id = %s", id)
	}
	return ps, nil
}

// Returns errox.AlreadyExists if there is a permission set with the same ID.
func (ds *dataStoreImpl) verifyPermissionSetIDDoesNotExist(ctx context.Context, id string) error {
	_, found, err := ds.permissionSetStorage.Get(ctx, id)

	if err != nil {
		return err
	}
	if found {
		return errors.Wrapf(errox.AlreadyExists, "id = %s", id)
	}
	return nil
}

// Returns errox.InvalidArgs if the given permission set is a default
// one. Note that IsDefaultRole() is reused due to the name sameness.
func verifyNotDefaultPermissionSet(permissionSet *storage.PermissionSet) error {
	if rolePkg.IsDefaultPermissionSet(permissionSet) {
		return errors.Wrapf(errox.InvalidArgs, "default permission set %q cannot be modified or deleted",
			permissionSet.GetName())
	}
	return nil
}

// Returns errox.NotFound if there is no access scope with the supplied ID.
func (ds *dataStoreImpl) verifyAccessScopeIDExists(ctx context.Context, id string) (*storage.SimpleAccessScope, error) {
	as, found, err := ds.accessScopeStorage.Get(ctx, id)

	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errors.Wrapf(errox.NotFound, "id = %s", id)
	}
	return as, nil
}

// Returns errox.AlreadyExists if there is an access scope with the same ID.
func (ds *dataStoreImpl) verifyAccessScopeIDDoesNotExist(ctx context.Context, id string) error {
	_, found, err := ds.accessScopeStorage.Get(ctx, id)

	if err != nil {
		return err
	}
	if found {
		return errors.Wrapf(errox.AlreadyExists, "id = %s", id)
	}
	return nil
}

// Returns errox.AlreadyExists if there is a role with the same name.
func (ds *dataStoreImpl) verifyRoleNameDoesNotExist(ctx context.Context, name string) error {
	_, found, err := ds.roleStorage.Get(ctx, name)

	if err != nil {
		return err
	}
	if found {
		return errors.Wrapf(errox.AlreadyExists, "name = %q", name)
	}
	return nil
}

// Returns errox.NotFound if there is no role with the supplied name.
func (ds *dataStoreImpl) verifyRoleNameExists(ctx context.Context, name string) (*storage.Role, error) {
	role, found, err := ds.roleStorage.Get(ctx, name)

	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errors.Wrapf(errox.NotFound, "name = %q", name)
	}
	return role, nil
}

// verifyRoleForDeletion verifies the storage constraints for deleting a role.
// It will:
// - verify that the role is not a default role
// - verify that the role exists
func (ds *dataStoreImpl) verifyRoleForDeletion(ctx context.Context, name string) error {
	role, found, err := ds.roleStorage.Get(ctx, name)

	if err != nil {
		return err
	}
	if !found {
		return errors.Wrapf(errox.NotFound, "name = %q", name)
	}
	if err = verifyRoleOriginMatches(ctx, role); err != nil {
		return err
	}

	return verifyNotDefaultRole(role)
}

// Returns errox.InvalidArgs if the given scope is a default one.
func verifyNotDefaultAccessScope(scope *storage.SimpleAccessScope) error {
	if rolePkg.IsDefaultAccessScope(scope) {
		return errors.Wrapf(errox.InvalidArgs, "default access scope %q cannot be modified or deleted", scope.GetName())
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// Helpers                                                                    //
//                                                                            //

// Finds the permission set associated with the given role. Every stored role
// must reference an existing permission set.
func (ds *dataStoreImpl) getRolePermissionSetOrError(ctx context.Context, role *storage.Role) (*storage.PermissionSet, error) {
	permissionSet, found, err := ds.permissionSetStorage.Get(ctx, role.GetPermissionSetId())
	if err != nil {
		return nil, err
	} else if !found || permissionSet == nil {
		log.Errorf("Failed to fetch permission set %s for the existing role %q", role.GetPermissionSetId(), role.GetName())
		return nil, errors.Wrapf(errox.InvariantViolation, "permission set %s for role %q is missing", role.GetPermissionSetId(), role.GetName())
	}
	return permissionSet, nil
}

// Finds the access scope associated with the given role. Every stored role must
// reference an existing access scope.
func (ds *dataStoreImpl) getRoleAccessScopeOrError(ctx context.Context, role *storage.Role) (*storage.SimpleAccessScope, error) {
	accessScope, found, err := ds.accessScopeStorage.Get(ctx, role.GetAccessScopeId())
	if err != nil {
		return nil, err
	} else if !found || accessScope == nil {
		log.Errorf("Failed to fetch access scope %s for the existing role %q", role.GetAccessScopeId(), role.GetName())
		return nil, errors.Wrapf(errox.InvariantViolation, "access scope %s for role %q is missing", role.GetAccessScopeId(), role.GetName())
	}
	return accessScope, nil
}
