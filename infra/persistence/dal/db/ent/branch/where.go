// Code generated by ent, DO NOT EDIT.

package branch

import (
	"purchase/infra/persistence/dal/db/ent/predicate"
	"time"

	"entgo.io/ent/dialect/sql"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.Branch {
	return predicate.Branch(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Branch {
	return predicate.Branch(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Branch {
	return predicate.Branch(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Branch {
	return predicate.Branch(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Branch {
	return predicate.Branch(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Branch {
	return predicate.Branch(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Branch {
	return predicate.Branch(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Branch {
	return predicate.Branch(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Branch {
	return predicate.Branch(sql.FieldLTE(FieldID, id))
}

// Code applies equality check predicate on the "code" field. It's identical to CodeEQ.
func Code(v string) predicate.Branch {
	return predicate.Branch(sql.FieldEQ(FieldCode, v))
}

// TransID applies equality check predicate on the "trans_id" field. It's identical to TransIDEQ.
func TransID(v int) predicate.Branch {
	return predicate.Branch(sql.FieldEQ(FieldTransID, v))
}

// Type applies equality check predicate on the "type" field. It's identical to TypeEQ.
func Type(v string) predicate.Branch {
	return predicate.Branch(sql.FieldEQ(FieldType, v))
}

// State applies equality check predicate on the "state" field. It's identical to StateEQ.
func State(v string) predicate.Branch {
	return predicate.Branch(sql.FieldEQ(FieldState, v))
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.Branch {
	return predicate.Branch(sql.FieldEQ(FieldName, v))
}

// Action applies equality check predicate on the "action" field. It's identical to ActionEQ.
func Action(v string) predicate.Branch {
	return predicate.Branch(sql.FieldEQ(FieldAction, v))
}

// Compensate applies equality check predicate on the "compensate" field. It's identical to CompensateEQ.
func Compensate(v string) predicate.Branch {
	return predicate.Branch(sql.FieldEQ(FieldCompensate, v))
}

// ActionPayload applies equality check predicate on the "action_payload" field. It's identical to ActionPayloadEQ.
func ActionPayload(v string) predicate.Branch {
	return predicate.Branch(sql.FieldEQ(FieldActionPayload, v))
}

// CompensatePayload applies equality check predicate on the "compensate_payload" field. It's identical to CompensatePayloadEQ.
func CompensatePayload(v string) predicate.Branch {
	return predicate.Branch(sql.FieldEQ(FieldCompensatePayload, v))
}

// ActionDepend applies equality check predicate on the "action_depend" field. It's identical to ActionDependEQ.
func ActionDepend(v string) predicate.Branch {
	return predicate.Branch(sql.FieldEQ(FieldActionDepend, v))
}

// CompensateDepend applies equality check predicate on the "compensate_depend" field. It's identical to CompensateDependEQ.
func CompensateDepend(v string) predicate.Branch {
	return predicate.Branch(sql.FieldEQ(FieldCompensateDepend, v))
}

// IsDead applies equality check predicate on the "is_dead" field. It's identical to IsDeadEQ.
func IsDead(v bool) predicate.Branch {
	return predicate.Branch(sql.FieldEQ(FieldIsDead, v))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.Branch {
	return predicate.Branch(sql.FieldEQ(FieldCreatedAt, v))
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.Branch {
	return predicate.Branch(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedBy applies equality check predicate on the "updated_by" field. It's identical to UpdatedByEQ.
func UpdatedBy(v string) predicate.Branch {
	return predicate.Branch(sql.FieldEQ(FieldUpdatedBy, v))
}

// CreatedBy applies equality check predicate on the "created_by" field. It's identical to CreatedByEQ.
func CreatedBy(v string) predicate.Branch {
	return predicate.Branch(sql.FieldEQ(FieldCreatedBy, v))
}

// CodeEQ applies the EQ predicate on the "code" field.
func CodeEQ(v string) predicate.Branch {
	return predicate.Branch(sql.FieldEQ(FieldCode, v))
}

// CodeNEQ applies the NEQ predicate on the "code" field.
func CodeNEQ(v string) predicate.Branch {
	return predicate.Branch(sql.FieldNEQ(FieldCode, v))
}

// CodeIn applies the In predicate on the "code" field.
func CodeIn(vs ...string) predicate.Branch {
	return predicate.Branch(sql.FieldIn(FieldCode, vs...))
}

// CodeNotIn applies the NotIn predicate on the "code" field.
func CodeNotIn(vs ...string) predicate.Branch {
	return predicate.Branch(sql.FieldNotIn(FieldCode, vs...))
}

// CodeGT applies the GT predicate on the "code" field.
func CodeGT(v string) predicate.Branch {
	return predicate.Branch(sql.FieldGT(FieldCode, v))
}

// CodeGTE applies the GTE predicate on the "code" field.
func CodeGTE(v string) predicate.Branch {
	return predicate.Branch(sql.FieldGTE(FieldCode, v))
}

// CodeLT applies the LT predicate on the "code" field.
func CodeLT(v string) predicate.Branch {
	return predicate.Branch(sql.FieldLT(FieldCode, v))
}

// CodeLTE applies the LTE predicate on the "code" field.
func CodeLTE(v string) predicate.Branch {
	return predicate.Branch(sql.FieldLTE(FieldCode, v))
}

// CodeContains applies the Contains predicate on the "code" field.
func CodeContains(v string) predicate.Branch {
	return predicate.Branch(sql.FieldContains(FieldCode, v))
}

// CodeHasPrefix applies the HasPrefix predicate on the "code" field.
func CodeHasPrefix(v string) predicate.Branch {
	return predicate.Branch(sql.FieldHasPrefix(FieldCode, v))
}

// CodeHasSuffix applies the HasSuffix predicate on the "code" field.
func CodeHasSuffix(v string) predicate.Branch {
	return predicate.Branch(sql.FieldHasSuffix(FieldCode, v))
}

// CodeEqualFold applies the EqualFold predicate on the "code" field.
func CodeEqualFold(v string) predicate.Branch {
	return predicate.Branch(sql.FieldEqualFold(FieldCode, v))
}

// CodeContainsFold applies the ContainsFold predicate on the "code" field.
func CodeContainsFold(v string) predicate.Branch {
	return predicate.Branch(sql.FieldContainsFold(FieldCode, v))
}

// TransIDEQ applies the EQ predicate on the "trans_id" field.
func TransIDEQ(v int) predicate.Branch {
	return predicate.Branch(sql.FieldEQ(FieldTransID, v))
}

// TransIDNEQ applies the NEQ predicate on the "trans_id" field.
func TransIDNEQ(v int) predicate.Branch {
	return predicate.Branch(sql.FieldNEQ(FieldTransID, v))
}

// TransIDIn applies the In predicate on the "trans_id" field.
func TransIDIn(vs ...int) predicate.Branch {
	return predicate.Branch(sql.FieldIn(FieldTransID, vs...))
}

// TransIDNotIn applies the NotIn predicate on the "trans_id" field.
func TransIDNotIn(vs ...int) predicate.Branch {
	return predicate.Branch(sql.FieldNotIn(FieldTransID, vs...))
}

// TransIDGT applies the GT predicate on the "trans_id" field.
func TransIDGT(v int) predicate.Branch {
	return predicate.Branch(sql.FieldGT(FieldTransID, v))
}

// TransIDGTE applies the GTE predicate on the "trans_id" field.
func TransIDGTE(v int) predicate.Branch {
	return predicate.Branch(sql.FieldGTE(FieldTransID, v))
}

// TransIDLT applies the LT predicate on the "trans_id" field.
func TransIDLT(v int) predicate.Branch {
	return predicate.Branch(sql.FieldLT(FieldTransID, v))
}

// TransIDLTE applies the LTE predicate on the "trans_id" field.
func TransIDLTE(v int) predicate.Branch {
	return predicate.Branch(sql.FieldLTE(FieldTransID, v))
}

// TypeEQ applies the EQ predicate on the "type" field.
func TypeEQ(v string) predicate.Branch {
	return predicate.Branch(sql.FieldEQ(FieldType, v))
}

// TypeNEQ applies the NEQ predicate on the "type" field.
func TypeNEQ(v string) predicate.Branch {
	return predicate.Branch(sql.FieldNEQ(FieldType, v))
}

// TypeIn applies the In predicate on the "type" field.
func TypeIn(vs ...string) predicate.Branch {
	return predicate.Branch(sql.FieldIn(FieldType, vs...))
}

// TypeNotIn applies the NotIn predicate on the "type" field.
func TypeNotIn(vs ...string) predicate.Branch {
	return predicate.Branch(sql.FieldNotIn(FieldType, vs...))
}

// TypeGT applies the GT predicate on the "type" field.
func TypeGT(v string) predicate.Branch {
	return predicate.Branch(sql.FieldGT(FieldType, v))
}

// TypeGTE applies the GTE predicate on the "type" field.
func TypeGTE(v string) predicate.Branch {
	return predicate.Branch(sql.FieldGTE(FieldType, v))
}

// TypeLT applies the LT predicate on the "type" field.
func TypeLT(v string) predicate.Branch {
	return predicate.Branch(sql.FieldLT(FieldType, v))
}

// TypeLTE applies the LTE predicate on the "type" field.
func TypeLTE(v string) predicate.Branch {
	return predicate.Branch(sql.FieldLTE(FieldType, v))
}

// TypeContains applies the Contains predicate on the "type" field.
func TypeContains(v string) predicate.Branch {
	return predicate.Branch(sql.FieldContains(FieldType, v))
}

// TypeHasPrefix applies the HasPrefix predicate on the "type" field.
func TypeHasPrefix(v string) predicate.Branch {
	return predicate.Branch(sql.FieldHasPrefix(FieldType, v))
}

// TypeHasSuffix applies the HasSuffix predicate on the "type" field.
func TypeHasSuffix(v string) predicate.Branch {
	return predicate.Branch(sql.FieldHasSuffix(FieldType, v))
}

// TypeEqualFold applies the EqualFold predicate on the "type" field.
func TypeEqualFold(v string) predicate.Branch {
	return predicate.Branch(sql.FieldEqualFold(FieldType, v))
}

// TypeContainsFold applies the ContainsFold predicate on the "type" field.
func TypeContainsFold(v string) predicate.Branch {
	return predicate.Branch(sql.FieldContainsFold(FieldType, v))
}

// StateEQ applies the EQ predicate on the "state" field.
func StateEQ(v string) predicate.Branch {
	return predicate.Branch(sql.FieldEQ(FieldState, v))
}

// StateNEQ applies the NEQ predicate on the "state" field.
func StateNEQ(v string) predicate.Branch {
	return predicate.Branch(sql.FieldNEQ(FieldState, v))
}

// StateIn applies the In predicate on the "state" field.
func StateIn(vs ...string) predicate.Branch {
	return predicate.Branch(sql.FieldIn(FieldState, vs...))
}

// StateNotIn applies the NotIn predicate on the "state" field.
func StateNotIn(vs ...string) predicate.Branch {
	return predicate.Branch(sql.FieldNotIn(FieldState, vs...))
}

// StateGT applies the GT predicate on the "state" field.
func StateGT(v string) predicate.Branch {
	return predicate.Branch(sql.FieldGT(FieldState, v))
}

// StateGTE applies the GTE predicate on the "state" field.
func StateGTE(v string) predicate.Branch {
	return predicate.Branch(sql.FieldGTE(FieldState, v))
}

// StateLT applies the LT predicate on the "state" field.
func StateLT(v string) predicate.Branch {
	return predicate.Branch(sql.FieldLT(FieldState, v))
}

// StateLTE applies the LTE predicate on the "state" field.
func StateLTE(v string) predicate.Branch {
	return predicate.Branch(sql.FieldLTE(FieldState, v))
}

// StateContains applies the Contains predicate on the "state" field.
func StateContains(v string) predicate.Branch {
	return predicate.Branch(sql.FieldContains(FieldState, v))
}

// StateHasPrefix applies the HasPrefix predicate on the "state" field.
func StateHasPrefix(v string) predicate.Branch {
	return predicate.Branch(sql.FieldHasPrefix(FieldState, v))
}

// StateHasSuffix applies the HasSuffix predicate on the "state" field.
func StateHasSuffix(v string) predicate.Branch {
	return predicate.Branch(sql.FieldHasSuffix(FieldState, v))
}

// StateEqualFold applies the EqualFold predicate on the "state" field.
func StateEqualFold(v string) predicate.Branch {
	return predicate.Branch(sql.FieldEqualFold(FieldState, v))
}

// StateContainsFold applies the ContainsFold predicate on the "state" field.
func StateContainsFold(v string) predicate.Branch {
	return predicate.Branch(sql.FieldContainsFold(FieldState, v))
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.Branch {
	return predicate.Branch(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.Branch {
	return predicate.Branch(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.Branch {
	return predicate.Branch(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.Branch {
	return predicate.Branch(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.Branch {
	return predicate.Branch(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.Branch {
	return predicate.Branch(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.Branch {
	return predicate.Branch(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.Branch {
	return predicate.Branch(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.Branch {
	return predicate.Branch(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.Branch {
	return predicate.Branch(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.Branch {
	return predicate.Branch(sql.FieldHasSuffix(FieldName, v))
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.Branch {
	return predicate.Branch(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.Branch {
	return predicate.Branch(sql.FieldContainsFold(FieldName, v))
}

// ActionEQ applies the EQ predicate on the "action" field.
func ActionEQ(v string) predicate.Branch {
	return predicate.Branch(sql.FieldEQ(FieldAction, v))
}

// ActionNEQ applies the NEQ predicate on the "action" field.
func ActionNEQ(v string) predicate.Branch {
	return predicate.Branch(sql.FieldNEQ(FieldAction, v))
}

// ActionIn applies the In predicate on the "action" field.
func ActionIn(vs ...string) predicate.Branch {
	return predicate.Branch(sql.FieldIn(FieldAction, vs...))
}

// ActionNotIn applies the NotIn predicate on the "action" field.
func ActionNotIn(vs ...string) predicate.Branch {
	return predicate.Branch(sql.FieldNotIn(FieldAction, vs...))
}

// ActionGT applies the GT predicate on the "action" field.
func ActionGT(v string) predicate.Branch {
	return predicate.Branch(sql.FieldGT(FieldAction, v))
}

// ActionGTE applies the GTE predicate on the "action" field.
func ActionGTE(v string) predicate.Branch {
	return predicate.Branch(sql.FieldGTE(FieldAction, v))
}

// ActionLT applies the LT predicate on the "action" field.
func ActionLT(v string) predicate.Branch {
	return predicate.Branch(sql.FieldLT(FieldAction, v))
}

// ActionLTE applies the LTE predicate on the "action" field.
func ActionLTE(v string) predicate.Branch {
	return predicate.Branch(sql.FieldLTE(FieldAction, v))
}

// ActionContains applies the Contains predicate on the "action" field.
func ActionContains(v string) predicate.Branch {
	return predicate.Branch(sql.FieldContains(FieldAction, v))
}

// ActionHasPrefix applies the HasPrefix predicate on the "action" field.
func ActionHasPrefix(v string) predicate.Branch {
	return predicate.Branch(sql.FieldHasPrefix(FieldAction, v))
}

// ActionHasSuffix applies the HasSuffix predicate on the "action" field.
func ActionHasSuffix(v string) predicate.Branch {
	return predicate.Branch(sql.FieldHasSuffix(FieldAction, v))
}

// ActionEqualFold applies the EqualFold predicate on the "action" field.
func ActionEqualFold(v string) predicate.Branch {
	return predicate.Branch(sql.FieldEqualFold(FieldAction, v))
}

// ActionContainsFold applies the ContainsFold predicate on the "action" field.
func ActionContainsFold(v string) predicate.Branch {
	return predicate.Branch(sql.FieldContainsFold(FieldAction, v))
}

// CompensateEQ applies the EQ predicate on the "compensate" field.
func CompensateEQ(v string) predicate.Branch {
	return predicate.Branch(sql.FieldEQ(FieldCompensate, v))
}

// CompensateNEQ applies the NEQ predicate on the "compensate" field.
func CompensateNEQ(v string) predicate.Branch {
	return predicate.Branch(sql.FieldNEQ(FieldCompensate, v))
}

// CompensateIn applies the In predicate on the "compensate" field.
func CompensateIn(vs ...string) predicate.Branch {
	return predicate.Branch(sql.FieldIn(FieldCompensate, vs...))
}

// CompensateNotIn applies the NotIn predicate on the "compensate" field.
func CompensateNotIn(vs ...string) predicate.Branch {
	return predicate.Branch(sql.FieldNotIn(FieldCompensate, vs...))
}

// CompensateGT applies the GT predicate on the "compensate" field.
func CompensateGT(v string) predicate.Branch {
	return predicate.Branch(sql.FieldGT(FieldCompensate, v))
}

// CompensateGTE applies the GTE predicate on the "compensate" field.
func CompensateGTE(v string) predicate.Branch {
	return predicate.Branch(sql.FieldGTE(FieldCompensate, v))
}

// CompensateLT applies the LT predicate on the "compensate" field.
func CompensateLT(v string) predicate.Branch {
	return predicate.Branch(sql.FieldLT(FieldCompensate, v))
}

// CompensateLTE applies the LTE predicate on the "compensate" field.
func CompensateLTE(v string) predicate.Branch {
	return predicate.Branch(sql.FieldLTE(FieldCompensate, v))
}

// CompensateContains applies the Contains predicate on the "compensate" field.
func CompensateContains(v string) predicate.Branch {
	return predicate.Branch(sql.FieldContains(FieldCompensate, v))
}

// CompensateHasPrefix applies the HasPrefix predicate on the "compensate" field.
func CompensateHasPrefix(v string) predicate.Branch {
	return predicate.Branch(sql.FieldHasPrefix(FieldCompensate, v))
}

// CompensateHasSuffix applies the HasSuffix predicate on the "compensate" field.
func CompensateHasSuffix(v string) predicate.Branch {
	return predicate.Branch(sql.FieldHasSuffix(FieldCompensate, v))
}

// CompensateEqualFold applies the EqualFold predicate on the "compensate" field.
func CompensateEqualFold(v string) predicate.Branch {
	return predicate.Branch(sql.FieldEqualFold(FieldCompensate, v))
}

// CompensateContainsFold applies the ContainsFold predicate on the "compensate" field.
func CompensateContainsFold(v string) predicate.Branch {
	return predicate.Branch(sql.FieldContainsFold(FieldCompensate, v))
}

// ActionPayloadEQ applies the EQ predicate on the "action_payload" field.
func ActionPayloadEQ(v string) predicate.Branch {
	return predicate.Branch(sql.FieldEQ(FieldActionPayload, v))
}

// ActionPayloadNEQ applies the NEQ predicate on the "action_payload" field.
func ActionPayloadNEQ(v string) predicate.Branch {
	return predicate.Branch(sql.FieldNEQ(FieldActionPayload, v))
}

// ActionPayloadIn applies the In predicate on the "action_payload" field.
func ActionPayloadIn(vs ...string) predicate.Branch {
	return predicate.Branch(sql.FieldIn(FieldActionPayload, vs...))
}

// ActionPayloadNotIn applies the NotIn predicate on the "action_payload" field.
func ActionPayloadNotIn(vs ...string) predicate.Branch {
	return predicate.Branch(sql.FieldNotIn(FieldActionPayload, vs...))
}

// ActionPayloadGT applies the GT predicate on the "action_payload" field.
func ActionPayloadGT(v string) predicate.Branch {
	return predicate.Branch(sql.FieldGT(FieldActionPayload, v))
}

// ActionPayloadGTE applies the GTE predicate on the "action_payload" field.
func ActionPayloadGTE(v string) predicate.Branch {
	return predicate.Branch(sql.FieldGTE(FieldActionPayload, v))
}

// ActionPayloadLT applies the LT predicate on the "action_payload" field.
func ActionPayloadLT(v string) predicate.Branch {
	return predicate.Branch(sql.FieldLT(FieldActionPayload, v))
}

// ActionPayloadLTE applies the LTE predicate on the "action_payload" field.
func ActionPayloadLTE(v string) predicate.Branch {
	return predicate.Branch(sql.FieldLTE(FieldActionPayload, v))
}

// ActionPayloadContains applies the Contains predicate on the "action_payload" field.
func ActionPayloadContains(v string) predicate.Branch {
	return predicate.Branch(sql.FieldContains(FieldActionPayload, v))
}

// ActionPayloadHasPrefix applies the HasPrefix predicate on the "action_payload" field.
func ActionPayloadHasPrefix(v string) predicate.Branch {
	return predicate.Branch(sql.FieldHasPrefix(FieldActionPayload, v))
}

// ActionPayloadHasSuffix applies the HasSuffix predicate on the "action_payload" field.
func ActionPayloadHasSuffix(v string) predicate.Branch {
	return predicate.Branch(sql.FieldHasSuffix(FieldActionPayload, v))
}

// ActionPayloadEqualFold applies the EqualFold predicate on the "action_payload" field.
func ActionPayloadEqualFold(v string) predicate.Branch {
	return predicate.Branch(sql.FieldEqualFold(FieldActionPayload, v))
}

// ActionPayloadContainsFold applies the ContainsFold predicate on the "action_payload" field.
func ActionPayloadContainsFold(v string) predicate.Branch {
	return predicate.Branch(sql.FieldContainsFold(FieldActionPayload, v))
}

// CompensatePayloadEQ applies the EQ predicate on the "compensate_payload" field.
func CompensatePayloadEQ(v string) predicate.Branch {
	return predicate.Branch(sql.FieldEQ(FieldCompensatePayload, v))
}

// CompensatePayloadNEQ applies the NEQ predicate on the "compensate_payload" field.
func CompensatePayloadNEQ(v string) predicate.Branch {
	return predicate.Branch(sql.FieldNEQ(FieldCompensatePayload, v))
}

// CompensatePayloadIn applies the In predicate on the "compensate_payload" field.
func CompensatePayloadIn(vs ...string) predicate.Branch {
	return predicate.Branch(sql.FieldIn(FieldCompensatePayload, vs...))
}

// CompensatePayloadNotIn applies the NotIn predicate on the "compensate_payload" field.
func CompensatePayloadNotIn(vs ...string) predicate.Branch {
	return predicate.Branch(sql.FieldNotIn(FieldCompensatePayload, vs...))
}

// CompensatePayloadGT applies the GT predicate on the "compensate_payload" field.
func CompensatePayloadGT(v string) predicate.Branch {
	return predicate.Branch(sql.FieldGT(FieldCompensatePayload, v))
}

// CompensatePayloadGTE applies the GTE predicate on the "compensate_payload" field.
func CompensatePayloadGTE(v string) predicate.Branch {
	return predicate.Branch(sql.FieldGTE(FieldCompensatePayload, v))
}

// CompensatePayloadLT applies the LT predicate on the "compensate_payload" field.
func CompensatePayloadLT(v string) predicate.Branch {
	return predicate.Branch(sql.FieldLT(FieldCompensatePayload, v))
}

// CompensatePayloadLTE applies the LTE predicate on the "compensate_payload" field.
func CompensatePayloadLTE(v string) predicate.Branch {
	return predicate.Branch(sql.FieldLTE(FieldCompensatePayload, v))
}

// CompensatePayloadContains applies the Contains predicate on the "compensate_payload" field.
func CompensatePayloadContains(v string) predicate.Branch {
	return predicate.Branch(sql.FieldContains(FieldCompensatePayload, v))
}

// CompensatePayloadHasPrefix applies the HasPrefix predicate on the "compensate_payload" field.
func CompensatePayloadHasPrefix(v string) predicate.Branch {
	return predicate.Branch(sql.FieldHasPrefix(FieldCompensatePayload, v))
}

// CompensatePayloadHasSuffix applies the HasSuffix predicate on the "compensate_payload" field.
func CompensatePayloadHasSuffix(v string) predicate.Branch {
	return predicate.Branch(sql.FieldHasSuffix(FieldCompensatePayload, v))
}

// CompensatePayloadEqualFold applies the EqualFold predicate on the "compensate_payload" field.
func CompensatePayloadEqualFold(v string) predicate.Branch {
	return predicate.Branch(sql.FieldEqualFold(FieldCompensatePayload, v))
}

// CompensatePayloadContainsFold applies the ContainsFold predicate on the "compensate_payload" field.
func CompensatePayloadContainsFold(v string) predicate.Branch {
	return predicate.Branch(sql.FieldContainsFold(FieldCompensatePayload, v))
}

// ActionDependEQ applies the EQ predicate on the "action_depend" field.
func ActionDependEQ(v string) predicate.Branch {
	return predicate.Branch(sql.FieldEQ(FieldActionDepend, v))
}

// ActionDependNEQ applies the NEQ predicate on the "action_depend" field.
func ActionDependNEQ(v string) predicate.Branch {
	return predicate.Branch(sql.FieldNEQ(FieldActionDepend, v))
}

// ActionDependIn applies the In predicate on the "action_depend" field.
func ActionDependIn(vs ...string) predicate.Branch {
	return predicate.Branch(sql.FieldIn(FieldActionDepend, vs...))
}

// ActionDependNotIn applies the NotIn predicate on the "action_depend" field.
func ActionDependNotIn(vs ...string) predicate.Branch {
	return predicate.Branch(sql.FieldNotIn(FieldActionDepend, vs...))
}

// ActionDependGT applies the GT predicate on the "action_depend" field.
func ActionDependGT(v string) predicate.Branch {
	return predicate.Branch(sql.FieldGT(FieldActionDepend, v))
}

// ActionDependGTE applies the GTE predicate on the "action_depend" field.
func ActionDependGTE(v string) predicate.Branch {
	return predicate.Branch(sql.FieldGTE(FieldActionDepend, v))
}

// ActionDependLT applies the LT predicate on the "action_depend" field.
func ActionDependLT(v string) predicate.Branch {
	return predicate.Branch(sql.FieldLT(FieldActionDepend, v))
}

// ActionDependLTE applies the LTE predicate on the "action_depend" field.
func ActionDependLTE(v string) predicate.Branch {
	return predicate.Branch(sql.FieldLTE(FieldActionDepend, v))
}

// ActionDependContains applies the Contains predicate on the "action_depend" field.
func ActionDependContains(v string) predicate.Branch {
	return predicate.Branch(sql.FieldContains(FieldActionDepend, v))
}

// ActionDependHasPrefix applies the HasPrefix predicate on the "action_depend" field.
func ActionDependHasPrefix(v string) predicate.Branch {
	return predicate.Branch(sql.FieldHasPrefix(FieldActionDepend, v))
}

// ActionDependHasSuffix applies the HasSuffix predicate on the "action_depend" field.
func ActionDependHasSuffix(v string) predicate.Branch {
	return predicate.Branch(sql.FieldHasSuffix(FieldActionDepend, v))
}

// ActionDependEqualFold applies the EqualFold predicate on the "action_depend" field.
func ActionDependEqualFold(v string) predicate.Branch {
	return predicate.Branch(sql.FieldEqualFold(FieldActionDepend, v))
}

// ActionDependContainsFold applies the ContainsFold predicate on the "action_depend" field.
func ActionDependContainsFold(v string) predicate.Branch {
	return predicate.Branch(sql.FieldContainsFold(FieldActionDepend, v))
}

// CompensateDependEQ applies the EQ predicate on the "compensate_depend" field.
func CompensateDependEQ(v string) predicate.Branch {
	return predicate.Branch(sql.FieldEQ(FieldCompensateDepend, v))
}

// CompensateDependNEQ applies the NEQ predicate on the "compensate_depend" field.
func CompensateDependNEQ(v string) predicate.Branch {
	return predicate.Branch(sql.FieldNEQ(FieldCompensateDepend, v))
}

// CompensateDependIn applies the In predicate on the "compensate_depend" field.
func CompensateDependIn(vs ...string) predicate.Branch {
	return predicate.Branch(sql.FieldIn(FieldCompensateDepend, vs...))
}

// CompensateDependNotIn applies the NotIn predicate on the "compensate_depend" field.
func CompensateDependNotIn(vs ...string) predicate.Branch {
	return predicate.Branch(sql.FieldNotIn(FieldCompensateDepend, vs...))
}

// CompensateDependGT applies the GT predicate on the "compensate_depend" field.
func CompensateDependGT(v string) predicate.Branch {
	return predicate.Branch(sql.FieldGT(FieldCompensateDepend, v))
}

// CompensateDependGTE applies the GTE predicate on the "compensate_depend" field.
func CompensateDependGTE(v string) predicate.Branch {
	return predicate.Branch(sql.FieldGTE(FieldCompensateDepend, v))
}

// CompensateDependLT applies the LT predicate on the "compensate_depend" field.
func CompensateDependLT(v string) predicate.Branch {
	return predicate.Branch(sql.FieldLT(FieldCompensateDepend, v))
}

// CompensateDependLTE applies the LTE predicate on the "compensate_depend" field.
func CompensateDependLTE(v string) predicate.Branch {
	return predicate.Branch(sql.FieldLTE(FieldCompensateDepend, v))
}

// CompensateDependContains applies the Contains predicate on the "compensate_depend" field.
func CompensateDependContains(v string) predicate.Branch {
	return predicate.Branch(sql.FieldContains(FieldCompensateDepend, v))
}

// CompensateDependHasPrefix applies the HasPrefix predicate on the "compensate_depend" field.
func CompensateDependHasPrefix(v string) predicate.Branch {
	return predicate.Branch(sql.FieldHasPrefix(FieldCompensateDepend, v))
}

// CompensateDependHasSuffix applies the HasSuffix predicate on the "compensate_depend" field.
func CompensateDependHasSuffix(v string) predicate.Branch {
	return predicate.Branch(sql.FieldHasSuffix(FieldCompensateDepend, v))
}

// CompensateDependEqualFold applies the EqualFold predicate on the "compensate_depend" field.
func CompensateDependEqualFold(v string) predicate.Branch {
	return predicate.Branch(sql.FieldEqualFold(FieldCompensateDepend, v))
}

// CompensateDependContainsFold applies the ContainsFold predicate on the "compensate_depend" field.
func CompensateDependContainsFold(v string) predicate.Branch {
	return predicate.Branch(sql.FieldContainsFold(FieldCompensateDepend, v))
}

// IsDeadEQ applies the EQ predicate on the "is_dead" field.
func IsDeadEQ(v bool) predicate.Branch {
	return predicate.Branch(sql.FieldEQ(FieldIsDead, v))
}

// IsDeadNEQ applies the NEQ predicate on the "is_dead" field.
func IsDeadNEQ(v bool) predicate.Branch {
	return predicate.Branch(sql.FieldNEQ(FieldIsDead, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.Branch {
	return predicate.Branch(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.Branch {
	return predicate.Branch(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.Branch {
	return predicate.Branch(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.Branch {
	return predicate.Branch(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.Branch {
	return predicate.Branch(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.Branch {
	return predicate.Branch(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.Branch {
	return predicate.Branch(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.Branch {
	return predicate.Branch(sql.FieldLTE(FieldCreatedAt, v))
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.Branch {
	return predicate.Branch(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.Branch {
	return predicate.Branch(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.Branch {
	return predicate.Branch(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.Branch {
	return predicate.Branch(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.Branch {
	return predicate.Branch(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.Branch {
	return predicate.Branch(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.Branch {
	return predicate.Branch(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.Branch {
	return predicate.Branch(sql.FieldLTE(FieldUpdatedAt, v))
}

// UpdatedByEQ applies the EQ predicate on the "updated_by" field.
func UpdatedByEQ(v string) predicate.Branch {
	return predicate.Branch(sql.FieldEQ(FieldUpdatedBy, v))
}

// UpdatedByNEQ applies the NEQ predicate on the "updated_by" field.
func UpdatedByNEQ(v string) predicate.Branch {
	return predicate.Branch(sql.FieldNEQ(FieldUpdatedBy, v))
}

// UpdatedByIn applies the In predicate on the "updated_by" field.
func UpdatedByIn(vs ...string) predicate.Branch {
	return predicate.Branch(sql.FieldIn(FieldUpdatedBy, vs...))
}

// UpdatedByNotIn applies the NotIn predicate on the "updated_by" field.
func UpdatedByNotIn(vs ...string) predicate.Branch {
	return predicate.Branch(sql.FieldNotIn(FieldUpdatedBy, vs...))
}

// UpdatedByGT applies the GT predicate on the "updated_by" field.
func UpdatedByGT(v string) predicate.Branch {
	return predicate.Branch(sql.FieldGT(FieldUpdatedBy, v))
}

// UpdatedByGTE applies the GTE predicate on the "updated_by" field.
func UpdatedByGTE(v string) predicate.Branch {
	return predicate.Branch(sql.FieldGTE(FieldUpdatedBy, v))
}

// UpdatedByLT applies the LT predicate on the "updated_by" field.
func UpdatedByLT(v string) predicate.Branch {
	return predicate.Branch(sql.FieldLT(FieldUpdatedBy, v))
}

// UpdatedByLTE applies the LTE predicate on the "updated_by" field.
func UpdatedByLTE(v string) predicate.Branch {
	return predicate.Branch(sql.FieldLTE(FieldUpdatedBy, v))
}

// UpdatedByContains applies the Contains predicate on the "updated_by" field.
func UpdatedByContains(v string) predicate.Branch {
	return predicate.Branch(sql.FieldContains(FieldUpdatedBy, v))
}

// UpdatedByHasPrefix applies the HasPrefix predicate on the "updated_by" field.
func UpdatedByHasPrefix(v string) predicate.Branch {
	return predicate.Branch(sql.FieldHasPrefix(FieldUpdatedBy, v))
}

// UpdatedByHasSuffix applies the HasSuffix predicate on the "updated_by" field.
func UpdatedByHasSuffix(v string) predicate.Branch {
	return predicate.Branch(sql.FieldHasSuffix(FieldUpdatedBy, v))
}

// UpdatedByEqualFold applies the EqualFold predicate on the "updated_by" field.
func UpdatedByEqualFold(v string) predicate.Branch {
	return predicate.Branch(sql.FieldEqualFold(FieldUpdatedBy, v))
}

// UpdatedByContainsFold applies the ContainsFold predicate on the "updated_by" field.
func UpdatedByContainsFold(v string) predicate.Branch {
	return predicate.Branch(sql.FieldContainsFold(FieldUpdatedBy, v))
}

// CreatedByEQ applies the EQ predicate on the "created_by" field.
func CreatedByEQ(v string) predicate.Branch {
	return predicate.Branch(sql.FieldEQ(FieldCreatedBy, v))
}

// CreatedByNEQ applies the NEQ predicate on the "created_by" field.
func CreatedByNEQ(v string) predicate.Branch {
	return predicate.Branch(sql.FieldNEQ(FieldCreatedBy, v))
}

// CreatedByIn applies the In predicate on the "created_by" field.
func CreatedByIn(vs ...string) predicate.Branch {
	return predicate.Branch(sql.FieldIn(FieldCreatedBy, vs...))
}

// CreatedByNotIn applies the NotIn predicate on the "created_by" field.
func CreatedByNotIn(vs ...string) predicate.Branch {
	return predicate.Branch(sql.FieldNotIn(FieldCreatedBy, vs...))
}

// CreatedByGT applies the GT predicate on the "created_by" field.
func CreatedByGT(v string) predicate.Branch {
	return predicate.Branch(sql.FieldGT(FieldCreatedBy, v))
}

// CreatedByGTE applies the GTE predicate on the "created_by" field.
func CreatedByGTE(v string) predicate.Branch {
	return predicate.Branch(sql.FieldGTE(FieldCreatedBy, v))
}

// CreatedByLT applies the LT predicate on the "created_by" field.
func CreatedByLT(v string) predicate.Branch {
	return predicate.Branch(sql.FieldLT(FieldCreatedBy, v))
}

// CreatedByLTE applies the LTE predicate on the "created_by" field.
func CreatedByLTE(v string) predicate.Branch {
	return predicate.Branch(sql.FieldLTE(FieldCreatedBy, v))
}

// CreatedByContains applies the Contains predicate on the "created_by" field.
func CreatedByContains(v string) predicate.Branch {
	return predicate.Branch(sql.FieldContains(FieldCreatedBy, v))
}

// CreatedByHasPrefix applies the HasPrefix predicate on the "created_by" field.
func CreatedByHasPrefix(v string) predicate.Branch {
	return predicate.Branch(sql.FieldHasPrefix(FieldCreatedBy, v))
}

// CreatedByHasSuffix applies the HasSuffix predicate on the "created_by" field.
func CreatedByHasSuffix(v string) predicate.Branch {
	return predicate.Branch(sql.FieldHasSuffix(FieldCreatedBy, v))
}

// CreatedByEqualFold applies the EqualFold predicate on the "created_by" field.
func CreatedByEqualFold(v string) predicate.Branch {
	return predicate.Branch(sql.FieldEqualFold(FieldCreatedBy, v))
}

// CreatedByContainsFold applies the ContainsFold predicate on the "created_by" field.
func CreatedByContainsFold(v string) predicate.Branch {
	return predicate.Branch(sql.FieldContainsFold(FieldCreatedBy, v))
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Branch) predicate.Branch {
	return predicate.Branch(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Branch) predicate.Branch {
	return predicate.Branch(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Branch) predicate.Branch {
	return predicate.Branch(sql.NotPredicates(p))
}
