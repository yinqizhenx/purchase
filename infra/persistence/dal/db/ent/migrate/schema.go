// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// TAsyncTaskColumns holds the columns for the "t_async_task" table.
	TAsyncTaskColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt64, Increment: true},
		{Name: "task_id", Type: field.TypeString},
		{Name: "task_type", Type: field.TypeString},
		{Name: "task_group", Type: field.TypeString},
		{Name: "task_name", Type: field.TypeString},
		{Name: "biz_id", Type: field.TypeString},
		{Name: "task_data", Type: field.TypeString},
		{Name: "task_state", Type: field.TypeString},
		{Name: "created_at", Type: field.TypeTime, SchemaType: map[string]string{"mysql": "datetime"}},
		{Name: "updated_at", Type: field.TypeTime, SchemaType: map[string]string{"mysql": "datetime"}},
	}
	// TAsyncTaskTable holds the schema information for the "t_async_task" table.
	TAsyncTaskTable = &schema.Table{
		Name:       "t_async_task",
		Columns:    TAsyncTaskColumns,
		PrimaryKey: []*schema.Column{TAsyncTaskColumns[0]},
	}
	// BranchesColumns holds the columns for the "branches" table.
	BranchesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "trans_id", Type: field.TypeInt},
		{Name: "type", Type: field.TypeString},
		{Name: "state", Type: field.TypeString},
		{Name: "name", Type: field.TypeString},
		{Name: "action", Type: field.TypeString},
		{Name: "compensate", Type: field.TypeString},
		{Name: "payload", Type: field.TypeString},
		{Name: "action_depend", Type: field.TypeString},
		{Name: "compensate_depend", Type: field.TypeString},
		{Name: "finished_at", Type: field.TypeTime, SchemaType: map[string]string{"mysql": "datetime"}},
		{Name: "is_dead", Type: field.TypeBool},
		{Name: "created_at", Type: field.TypeTime, SchemaType: map[string]string{"mysql": "datetime"}},
		{Name: "updated_at", Type: field.TypeTime, SchemaType: map[string]string{"mysql": "datetime"}},
		{Name: "updated_by", Type: field.TypeString},
		{Name: "created_by", Type: field.TypeString},
	}
	// BranchesTable holds the schema information for the "branches" table.
	BranchesTable = &schema.Table{
		Name:       "branches",
		Columns:    BranchesColumns,
		PrimaryKey: []*schema.Column{BranchesColumns[0]},
	}
	// PaHeadsColumns holds the columns for the "pa_heads" table.
	PaHeadsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt64, Increment: true},
		{Name: "code", Type: field.TypeString, Unique: true},
		{Name: "state", Type: field.TypeString},
		{Name: "pay_amount", Type: field.TypeString},
		{Name: "applicant", Type: field.TypeString},
		{Name: "department_code", Type: field.TypeString, Nullable: true},
		{Name: "supplier_code", Type: field.TypeString},
		{Name: "is_adv", Type: field.TypeBool},
		{Name: "has_invoice", Type: field.TypeBool},
		{Name: "remark", Type: field.TypeString, Size: 2147483647},
		{Name: "created_at", Type: field.TypeTime, SchemaType: map[string]string{"mysql": "datetime"}},
		{Name: "updated_at", Type: field.TypeTime, SchemaType: map[string]string{"mysql": "datetime"}},
	}
	// PaHeadsTable holds the schema information for the "pa_heads" table.
	PaHeadsTable = &schema.Table{
		Name:       "pa_heads",
		Columns:    PaHeadsColumns,
		PrimaryKey: []*schema.Column{PaHeadsColumns[0]},
	}
	// PaRowsColumns holds the columns for the "pa_rows" table.
	PaRowsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt64, Increment: true},
		{Name: "head_code", Type: field.TypeString},
		{Name: "row_code", Type: field.TypeString},
		{Name: "grn_count", Type: field.TypeInt32},
		{Name: "grn_amount", Type: field.TypeString},
		{Name: "pay_amount", Type: field.TypeString},
		{Name: "description", Type: field.TypeString, Size: 2147483647},
		{Name: "created_at", Type: field.TypeTime, SchemaType: map[string]string{"mysql": "datetime"}},
		{Name: "updated_at", Type: field.TypeTime, SchemaType: map[string]string{"mysql": "datetime"}},
	}
	// PaRowsTable holds the schema information for the "pa_rows" table.
	PaRowsTable = &schema.Table{
		Name:       "pa_rows",
		Columns:    PaRowsColumns,
		PrimaryKey: []*schema.Column{PaRowsColumns[0]},
	}
	// TransColumns holds the columns for the "trans" table.
	TransColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "state", Type: field.TypeString},
		{Name: "name", Type: field.TypeString},
		{Name: "finished_at", Type: field.TypeTime, SchemaType: map[string]string{"mysql": "datetime"}},
		{Name: "created_at", Type: field.TypeTime, SchemaType: map[string]string{"mysql": "datetime"}},
		{Name: "updated_at", Type: field.TypeTime, SchemaType: map[string]string{"mysql": "datetime"}},
		{Name: "updated_by", Type: field.TypeString},
		{Name: "created_by", Type: field.TypeString},
	}
	// TransTable holds the schema information for the "trans" table.
	TransTable = &schema.Table{
		Name:       "trans",
		Columns:    TransColumns,
		PrimaryKey: []*schema.Column{TransColumns[0]},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		TAsyncTaskTable,
		BranchesTable,
		PaHeadsTable,
		PaRowsTable,
		TransTable,
	}
)

func init() {
	TAsyncTaskTable.Annotation = &entsql.Annotation{
		Table: "t_async_task",
	}
}
