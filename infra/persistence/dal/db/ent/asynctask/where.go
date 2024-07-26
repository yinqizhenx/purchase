// Code generated by ent, DO NOT EDIT.

package asynctask

import (
	"purchase/infra/persistence/dal/db/ent/predicate"
	"time"

	"entgo.io/ent/dialect/sql"
)

// ID filters vertices based on their ID field.
func ID(id int64) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int64) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int64) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int64) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int64) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int64) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int64) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int64) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int64) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldLTE(FieldID, id))
}

// TaskID applies equality check predicate on the "task_id" field. It's identical to TaskIDEQ.
func TaskID(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldEQ(FieldTaskID, v))
}

// TaskType applies equality check predicate on the "task_type" field. It's identical to TaskTypeEQ.
func TaskType(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldEQ(FieldTaskType, v))
}

// TaskName applies equality check predicate on the "task_name" field. It's identical to TaskNameEQ.
func TaskName(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldEQ(FieldTaskName, v))
}

// BizID applies equality check predicate on the "biz_id" field. It's identical to BizIDEQ.
func BizID(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldEQ(FieldBizID, v))
}

// TaskData applies equality check predicate on the "task_data" field. It's identical to TaskDataEQ.
func TaskData(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldEQ(FieldTaskData, v))
}

// State applies equality check predicate on the "state" field. It's identical to StateEQ.
func State(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldEQ(FieldState, v))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldEQ(FieldCreatedAt, v))
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldEQ(FieldUpdatedAt, v))
}

// TaskIDEQ applies the EQ predicate on the "task_id" field.
func TaskIDEQ(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldEQ(FieldTaskID, v))
}

// TaskIDNEQ applies the NEQ predicate on the "task_id" field.
func TaskIDNEQ(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldNEQ(FieldTaskID, v))
}

// TaskIDIn applies the In predicate on the "task_id" field.
func TaskIDIn(vs ...string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldIn(FieldTaskID, vs...))
}

// TaskIDNotIn applies the NotIn predicate on the "task_id" field.
func TaskIDNotIn(vs ...string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldNotIn(FieldTaskID, vs...))
}

// TaskIDGT applies the GT predicate on the "task_id" field.
func TaskIDGT(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldGT(FieldTaskID, v))
}

// TaskIDGTE applies the GTE predicate on the "task_id" field.
func TaskIDGTE(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldGTE(FieldTaskID, v))
}

// TaskIDLT applies the LT predicate on the "task_id" field.
func TaskIDLT(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldLT(FieldTaskID, v))
}

// TaskIDLTE applies the LTE predicate on the "task_id" field.
func TaskIDLTE(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldLTE(FieldTaskID, v))
}

// TaskIDContains applies the Contains predicate on the "task_id" field.
func TaskIDContains(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldContains(FieldTaskID, v))
}

// TaskIDHasPrefix applies the HasPrefix predicate on the "task_id" field.
func TaskIDHasPrefix(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldHasPrefix(FieldTaskID, v))
}

// TaskIDHasSuffix applies the HasSuffix predicate on the "task_id" field.
func TaskIDHasSuffix(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldHasSuffix(FieldTaskID, v))
}

// TaskIDEqualFold applies the EqualFold predicate on the "task_id" field.
func TaskIDEqualFold(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldEqualFold(FieldTaskID, v))
}

// TaskIDContainsFold applies the ContainsFold predicate on the "task_id" field.
func TaskIDContainsFold(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldContainsFold(FieldTaskID, v))
}

// TaskTypeEQ applies the EQ predicate on the "task_type" field.
func TaskTypeEQ(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldEQ(FieldTaskType, v))
}

// TaskTypeNEQ applies the NEQ predicate on the "task_type" field.
func TaskTypeNEQ(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldNEQ(FieldTaskType, v))
}

// TaskTypeIn applies the In predicate on the "task_type" field.
func TaskTypeIn(vs ...string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldIn(FieldTaskType, vs...))
}

// TaskTypeNotIn applies the NotIn predicate on the "task_type" field.
func TaskTypeNotIn(vs ...string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldNotIn(FieldTaskType, vs...))
}

// TaskTypeGT applies the GT predicate on the "task_type" field.
func TaskTypeGT(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldGT(FieldTaskType, v))
}

// TaskTypeGTE applies the GTE predicate on the "task_type" field.
func TaskTypeGTE(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldGTE(FieldTaskType, v))
}

// TaskTypeLT applies the LT predicate on the "task_type" field.
func TaskTypeLT(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldLT(FieldTaskType, v))
}

// TaskTypeLTE applies the LTE predicate on the "task_type" field.
func TaskTypeLTE(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldLTE(FieldTaskType, v))
}

// TaskTypeContains applies the Contains predicate on the "task_type" field.
func TaskTypeContains(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldContains(FieldTaskType, v))
}

// TaskTypeHasPrefix applies the HasPrefix predicate on the "task_type" field.
func TaskTypeHasPrefix(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldHasPrefix(FieldTaskType, v))
}

// TaskTypeHasSuffix applies the HasSuffix predicate on the "task_type" field.
func TaskTypeHasSuffix(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldHasSuffix(FieldTaskType, v))
}

// TaskTypeEqualFold applies the EqualFold predicate on the "task_type" field.
func TaskTypeEqualFold(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldEqualFold(FieldTaskType, v))
}

// TaskTypeContainsFold applies the ContainsFold predicate on the "task_type" field.
func TaskTypeContainsFold(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldContainsFold(FieldTaskType, v))
}

// TaskNameEQ applies the EQ predicate on the "task_name" field.
func TaskNameEQ(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldEQ(FieldTaskName, v))
}

// TaskNameNEQ applies the NEQ predicate on the "task_name" field.
func TaskNameNEQ(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldNEQ(FieldTaskName, v))
}

// TaskNameIn applies the In predicate on the "task_name" field.
func TaskNameIn(vs ...string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldIn(FieldTaskName, vs...))
}

// TaskNameNotIn applies the NotIn predicate on the "task_name" field.
func TaskNameNotIn(vs ...string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldNotIn(FieldTaskName, vs...))
}

// TaskNameGT applies the GT predicate on the "task_name" field.
func TaskNameGT(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldGT(FieldTaskName, v))
}

// TaskNameGTE applies the GTE predicate on the "task_name" field.
func TaskNameGTE(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldGTE(FieldTaskName, v))
}

// TaskNameLT applies the LT predicate on the "task_name" field.
func TaskNameLT(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldLT(FieldTaskName, v))
}

// TaskNameLTE applies the LTE predicate on the "task_name" field.
func TaskNameLTE(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldLTE(FieldTaskName, v))
}

// TaskNameContains applies the Contains predicate on the "task_name" field.
func TaskNameContains(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldContains(FieldTaskName, v))
}

// TaskNameHasPrefix applies the HasPrefix predicate on the "task_name" field.
func TaskNameHasPrefix(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldHasPrefix(FieldTaskName, v))
}

// TaskNameHasSuffix applies the HasSuffix predicate on the "task_name" field.
func TaskNameHasSuffix(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldHasSuffix(FieldTaskName, v))
}

// TaskNameEqualFold applies the EqualFold predicate on the "task_name" field.
func TaskNameEqualFold(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldEqualFold(FieldTaskName, v))
}

// TaskNameContainsFold applies the ContainsFold predicate on the "task_name" field.
func TaskNameContainsFold(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldContainsFold(FieldTaskName, v))
}

// BizIDEQ applies the EQ predicate on the "biz_id" field.
func BizIDEQ(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldEQ(FieldBizID, v))
}

// BizIDNEQ applies the NEQ predicate on the "biz_id" field.
func BizIDNEQ(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldNEQ(FieldBizID, v))
}

// BizIDIn applies the In predicate on the "biz_id" field.
func BizIDIn(vs ...string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldIn(FieldBizID, vs...))
}

// BizIDNotIn applies the NotIn predicate on the "biz_id" field.
func BizIDNotIn(vs ...string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldNotIn(FieldBizID, vs...))
}

// BizIDGT applies the GT predicate on the "biz_id" field.
func BizIDGT(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldGT(FieldBizID, v))
}

// BizIDGTE applies the GTE predicate on the "biz_id" field.
func BizIDGTE(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldGTE(FieldBizID, v))
}

// BizIDLT applies the LT predicate on the "biz_id" field.
func BizIDLT(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldLT(FieldBizID, v))
}

// BizIDLTE applies the LTE predicate on the "biz_id" field.
func BizIDLTE(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldLTE(FieldBizID, v))
}

// BizIDContains applies the Contains predicate on the "biz_id" field.
func BizIDContains(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldContains(FieldBizID, v))
}

// BizIDHasPrefix applies the HasPrefix predicate on the "biz_id" field.
func BizIDHasPrefix(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldHasPrefix(FieldBizID, v))
}

// BizIDHasSuffix applies the HasSuffix predicate on the "biz_id" field.
func BizIDHasSuffix(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldHasSuffix(FieldBizID, v))
}

// BizIDEqualFold applies the EqualFold predicate on the "biz_id" field.
func BizIDEqualFold(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldEqualFold(FieldBizID, v))
}

// BizIDContainsFold applies the ContainsFold predicate on the "biz_id" field.
func BizIDContainsFold(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldContainsFold(FieldBizID, v))
}

// TaskDataEQ applies the EQ predicate on the "task_data" field.
func TaskDataEQ(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldEQ(FieldTaskData, v))
}

// TaskDataNEQ applies the NEQ predicate on the "task_data" field.
func TaskDataNEQ(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldNEQ(FieldTaskData, v))
}

// TaskDataIn applies the In predicate on the "task_data" field.
func TaskDataIn(vs ...string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldIn(FieldTaskData, vs...))
}

// TaskDataNotIn applies the NotIn predicate on the "task_data" field.
func TaskDataNotIn(vs ...string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldNotIn(FieldTaskData, vs...))
}

// TaskDataGT applies the GT predicate on the "task_data" field.
func TaskDataGT(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldGT(FieldTaskData, v))
}

// TaskDataGTE applies the GTE predicate on the "task_data" field.
func TaskDataGTE(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldGTE(FieldTaskData, v))
}

// TaskDataLT applies the LT predicate on the "task_data" field.
func TaskDataLT(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldLT(FieldTaskData, v))
}

// TaskDataLTE applies the LTE predicate on the "task_data" field.
func TaskDataLTE(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldLTE(FieldTaskData, v))
}

// TaskDataContains applies the Contains predicate on the "task_data" field.
func TaskDataContains(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldContains(FieldTaskData, v))
}

// TaskDataHasPrefix applies the HasPrefix predicate on the "task_data" field.
func TaskDataHasPrefix(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldHasPrefix(FieldTaskData, v))
}

// TaskDataHasSuffix applies the HasSuffix predicate on the "task_data" field.
func TaskDataHasSuffix(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldHasSuffix(FieldTaskData, v))
}

// TaskDataEqualFold applies the EqualFold predicate on the "task_data" field.
func TaskDataEqualFold(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldEqualFold(FieldTaskData, v))
}

// TaskDataContainsFold applies the ContainsFold predicate on the "task_data" field.
func TaskDataContainsFold(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldContainsFold(FieldTaskData, v))
}

// StateEQ applies the EQ predicate on the "state" field.
func StateEQ(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldEQ(FieldState, v))
}

// StateNEQ applies the NEQ predicate on the "state" field.
func StateNEQ(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldNEQ(FieldState, v))
}

// StateIn applies the In predicate on the "state" field.
func StateIn(vs ...string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldIn(FieldState, vs...))
}

// StateNotIn applies the NotIn predicate on the "state" field.
func StateNotIn(vs ...string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldNotIn(FieldState, vs...))
}

// StateGT applies the GT predicate on the "state" field.
func StateGT(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldGT(FieldState, v))
}

// StateGTE applies the GTE predicate on the "state" field.
func StateGTE(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldGTE(FieldState, v))
}

// StateLT applies the LT predicate on the "state" field.
func StateLT(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldLT(FieldState, v))
}

// StateLTE applies the LTE predicate on the "state" field.
func StateLTE(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldLTE(FieldState, v))
}

// StateContains applies the Contains predicate on the "state" field.
func StateContains(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldContains(FieldState, v))
}

// StateHasPrefix applies the HasPrefix predicate on the "state" field.
func StateHasPrefix(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldHasPrefix(FieldState, v))
}

// StateHasSuffix applies the HasSuffix predicate on the "state" field.
func StateHasSuffix(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldHasSuffix(FieldState, v))
}

// StateEqualFold applies the EqualFold predicate on the "state" field.
func StateEqualFold(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldEqualFold(FieldState, v))
}

// StateContainsFold applies the ContainsFold predicate on the "state" field.
func StateContainsFold(v string) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldContainsFold(FieldState, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldLTE(FieldCreatedAt, v))
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.AsyncTask {
	return predicate.AsyncTask(sql.FieldLTE(FieldUpdatedAt, v))
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.AsyncTask) predicate.AsyncTask {
	return predicate.AsyncTask(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.AsyncTask) predicate.AsyncTask {
	return predicate.AsyncTask(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.AsyncTask) predicate.AsyncTask {
	return predicate.AsyncTask(sql.NotPredicates(p))
}
