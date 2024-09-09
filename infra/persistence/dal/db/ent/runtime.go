// Code generated by ent, DO NOT EDIT.

package ent

import (
	"purchase/infra/persistence/dal/db/ent/asynctask"
	"purchase/infra/persistence/dal/db/ent/branch"
	"purchase/infra/persistence/dal/db/ent/pahead"
	"purchase/infra/persistence/dal/db/ent/parow"
	"purchase/infra/persistence/dal/db/ent/schema"
	"purchase/infra/persistence/dal/db/ent/trans"
	"time"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	asynctaskFields := schema.AsyncTask{}.Fields()
	_ = asynctaskFields
	// asynctaskDescCreatedAt is the schema descriptor for created_at field.
	asynctaskDescCreatedAt := asynctaskFields[7].Descriptor()
	// asynctask.DefaultCreatedAt holds the default value on creation for the created_at field.
	asynctask.DefaultCreatedAt = asynctaskDescCreatedAt.Default.(func() time.Time)
	// asynctaskDescUpdatedAt is the schema descriptor for updated_at field.
	asynctaskDescUpdatedAt := asynctaskFields[8].Descriptor()
	// asynctask.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	asynctask.DefaultUpdatedAt = asynctaskDescUpdatedAt.Default.(func() time.Time)
	branchFields := schema.Branch{}.Fields()
	_ = branchFields
	// branchDescFinishedAt is the schema descriptor for finished_at field.
	branchDescFinishedAt := branchFields[9].Descriptor()
	// branch.DefaultFinishedAt holds the default value on creation for the finished_at field.
	branch.DefaultFinishedAt = branchDescFinishedAt.Default.(func() time.Time)
	// branchDescCreatedAt is the schema descriptor for created_at field.
	branchDescCreatedAt := branchFields[11].Descriptor()
	// branch.DefaultCreatedAt holds the default value on creation for the created_at field.
	branch.DefaultCreatedAt = branchDescCreatedAt.Default.(func() time.Time)
	// branchDescUpdatedAt is the schema descriptor for updated_at field.
	branchDescUpdatedAt := branchFields[12].Descriptor()
	// branch.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	branch.DefaultUpdatedAt = branchDescUpdatedAt.Default.(func() time.Time)
	paheadFields := schema.PAHead{}.Fields()
	_ = paheadFields
	// paheadDescCreatedAt is the schema descriptor for created_at field.
	paheadDescCreatedAt := paheadFields[10].Descriptor()
	// pahead.DefaultCreatedAt holds the default value on creation for the created_at field.
	pahead.DefaultCreatedAt = paheadDescCreatedAt.Default.(func() time.Time)
	// paheadDescUpdatedAt is the schema descriptor for updated_at field.
	paheadDescUpdatedAt := paheadFields[11].Descriptor()
	// pahead.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	pahead.DefaultUpdatedAt = paheadDescUpdatedAt.Default.(func() time.Time)
	parowFields := schema.PARow{}.Fields()
	_ = parowFields
	// parowDescCreatedAt is the schema descriptor for created_at field.
	parowDescCreatedAt := parowFields[7].Descriptor()
	// parow.DefaultCreatedAt holds the default value on creation for the created_at field.
	parow.DefaultCreatedAt = parowDescCreatedAt.Default.(func() time.Time)
	// parowDescUpdatedAt is the schema descriptor for updated_at field.
	parowDescUpdatedAt := parowFields[8].Descriptor()
	// parow.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	parow.DefaultUpdatedAt = parowDescUpdatedAt.Default.(func() time.Time)
	transFields := schema.Trans{}.Fields()
	_ = transFields
	// transDescFinishedAt is the schema descriptor for finished_at field.
	transDescFinishedAt := transFields[3].Descriptor()
	// trans.DefaultFinishedAt holds the default value on creation for the finished_at field.
	trans.DefaultFinishedAt = transDescFinishedAt.Default.(func() time.Time)
	// transDescCreatedAt is the schema descriptor for created_at field.
	transDescCreatedAt := transFields[4].Descriptor()
	// trans.DefaultCreatedAt holds the default value on creation for the created_at field.
	trans.DefaultCreatedAt = transDescCreatedAt.Default.(func() time.Time)
	// transDescUpdatedAt is the schema descriptor for updated_at field.
	transDescUpdatedAt := transFields[5].Descriptor()
	// trans.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	trans.DefaultUpdatedAt = transDescUpdatedAt.Default.(func() time.Time)
}
