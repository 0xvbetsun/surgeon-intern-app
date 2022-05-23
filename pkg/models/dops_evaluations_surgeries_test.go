// Code generated by SQLBoiler 4.8.3 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"bytes"
	"context"
	"reflect"
	"testing"

	"github.com/volatiletech/randomize"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/strmangle"
)

var (
	// Relationships sometimes use the reflection helper queries.Equal/queries.Assign
	// so force a package dependency in case they don't.
	_ = queries.Equal
)

func testDopsEvaluationsSurgeries(t *testing.T) {
	t.Parallel()

	query := DopsEvaluationsSurgeries()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testDopsEvaluationsSurgeriesDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &DopsEvaluationsSurgery{}
	if err = randomize.Struct(seed, o, dopsEvaluationsSurgeryDBTypes, true, dopsEvaluationsSurgeryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DopsEvaluationsSurgery struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := o.Delete(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := DopsEvaluationsSurgeries().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testDopsEvaluationsSurgeriesQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &DopsEvaluationsSurgery{}
	if err = randomize.Struct(seed, o, dopsEvaluationsSurgeryDBTypes, true, dopsEvaluationsSurgeryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DopsEvaluationsSurgery struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := DopsEvaluationsSurgeries().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := DopsEvaluationsSurgeries().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testDopsEvaluationsSurgeriesSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &DopsEvaluationsSurgery{}
	if err = randomize.Struct(seed, o, dopsEvaluationsSurgeryDBTypes, true, dopsEvaluationsSurgeryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DopsEvaluationsSurgery struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := DopsEvaluationsSurgerySlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := DopsEvaluationsSurgeries().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testDopsEvaluationsSurgeriesExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &DopsEvaluationsSurgery{}
	if err = randomize.Struct(seed, o, dopsEvaluationsSurgeryDBTypes, true, dopsEvaluationsSurgeryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DopsEvaluationsSurgery struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := DopsEvaluationsSurgeryExists(ctx, tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if DopsEvaluationsSurgery exists: %s", err)
	}
	if !e {
		t.Errorf("Expected DopsEvaluationsSurgeryExists to return true, but got false.")
	}
}

func testDopsEvaluationsSurgeriesFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &DopsEvaluationsSurgery{}
	if err = randomize.Struct(seed, o, dopsEvaluationsSurgeryDBTypes, true, dopsEvaluationsSurgeryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DopsEvaluationsSurgery struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	dopsEvaluationsSurgeryFound, err := FindDopsEvaluationsSurgery(ctx, tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if dopsEvaluationsSurgeryFound == nil {
		t.Error("want a record, got nil")
	}
}

func testDopsEvaluationsSurgeriesBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &DopsEvaluationsSurgery{}
	if err = randomize.Struct(seed, o, dopsEvaluationsSurgeryDBTypes, true, dopsEvaluationsSurgeryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DopsEvaluationsSurgery struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = DopsEvaluationsSurgeries().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testDopsEvaluationsSurgeriesOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &DopsEvaluationsSurgery{}
	if err = randomize.Struct(seed, o, dopsEvaluationsSurgeryDBTypes, true, dopsEvaluationsSurgeryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DopsEvaluationsSurgery struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := DopsEvaluationsSurgeries().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testDopsEvaluationsSurgeriesAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	dopsEvaluationsSurgeryOne := &DopsEvaluationsSurgery{}
	dopsEvaluationsSurgeryTwo := &DopsEvaluationsSurgery{}
	if err = randomize.Struct(seed, dopsEvaluationsSurgeryOne, dopsEvaluationsSurgeryDBTypes, false, dopsEvaluationsSurgeryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DopsEvaluationsSurgery struct: %s", err)
	}
	if err = randomize.Struct(seed, dopsEvaluationsSurgeryTwo, dopsEvaluationsSurgeryDBTypes, false, dopsEvaluationsSurgeryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DopsEvaluationsSurgery struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = dopsEvaluationsSurgeryOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = dopsEvaluationsSurgeryTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := DopsEvaluationsSurgeries().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testDopsEvaluationsSurgeriesCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	dopsEvaluationsSurgeryOne := &DopsEvaluationsSurgery{}
	dopsEvaluationsSurgeryTwo := &DopsEvaluationsSurgery{}
	if err = randomize.Struct(seed, dopsEvaluationsSurgeryOne, dopsEvaluationsSurgeryDBTypes, false, dopsEvaluationsSurgeryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DopsEvaluationsSurgery struct: %s", err)
	}
	if err = randomize.Struct(seed, dopsEvaluationsSurgeryTwo, dopsEvaluationsSurgeryDBTypes, false, dopsEvaluationsSurgeryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DopsEvaluationsSurgery struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = dopsEvaluationsSurgeryOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = dopsEvaluationsSurgeryTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := DopsEvaluationsSurgeries().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func dopsEvaluationsSurgeryBeforeInsertHook(ctx context.Context, e boil.ContextExecutor, o *DopsEvaluationsSurgery) error {
	*o = DopsEvaluationsSurgery{}
	return nil
}

func dopsEvaluationsSurgeryAfterInsertHook(ctx context.Context, e boil.ContextExecutor, o *DopsEvaluationsSurgery) error {
	*o = DopsEvaluationsSurgery{}
	return nil
}

func dopsEvaluationsSurgeryAfterSelectHook(ctx context.Context, e boil.ContextExecutor, o *DopsEvaluationsSurgery) error {
	*o = DopsEvaluationsSurgery{}
	return nil
}

func dopsEvaluationsSurgeryBeforeUpdateHook(ctx context.Context, e boil.ContextExecutor, o *DopsEvaluationsSurgery) error {
	*o = DopsEvaluationsSurgery{}
	return nil
}

func dopsEvaluationsSurgeryAfterUpdateHook(ctx context.Context, e boil.ContextExecutor, o *DopsEvaluationsSurgery) error {
	*o = DopsEvaluationsSurgery{}
	return nil
}

func dopsEvaluationsSurgeryBeforeDeleteHook(ctx context.Context, e boil.ContextExecutor, o *DopsEvaluationsSurgery) error {
	*o = DopsEvaluationsSurgery{}
	return nil
}

func dopsEvaluationsSurgeryAfterDeleteHook(ctx context.Context, e boil.ContextExecutor, o *DopsEvaluationsSurgery) error {
	*o = DopsEvaluationsSurgery{}
	return nil
}

func dopsEvaluationsSurgeryBeforeUpsertHook(ctx context.Context, e boil.ContextExecutor, o *DopsEvaluationsSurgery) error {
	*o = DopsEvaluationsSurgery{}
	return nil
}

func dopsEvaluationsSurgeryAfterUpsertHook(ctx context.Context, e boil.ContextExecutor, o *DopsEvaluationsSurgery) error {
	*o = DopsEvaluationsSurgery{}
	return nil
}

func testDopsEvaluationsSurgeriesHooks(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.Background()
	empty := &DopsEvaluationsSurgery{}
	o := &DopsEvaluationsSurgery{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, dopsEvaluationsSurgeryDBTypes, false); err != nil {
		t.Errorf("Unable to randomize DopsEvaluationsSurgery object: %s", err)
	}

	AddDopsEvaluationsSurgeryHook(boil.BeforeInsertHook, dopsEvaluationsSurgeryBeforeInsertHook)
	if err = o.doBeforeInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	dopsEvaluationsSurgeryBeforeInsertHooks = []DopsEvaluationsSurgeryHook{}

	AddDopsEvaluationsSurgeryHook(boil.AfterInsertHook, dopsEvaluationsSurgeryAfterInsertHook)
	if err = o.doAfterInsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	dopsEvaluationsSurgeryAfterInsertHooks = []DopsEvaluationsSurgeryHook{}

	AddDopsEvaluationsSurgeryHook(boil.AfterSelectHook, dopsEvaluationsSurgeryAfterSelectHook)
	if err = o.doAfterSelectHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	dopsEvaluationsSurgeryAfterSelectHooks = []DopsEvaluationsSurgeryHook{}

	AddDopsEvaluationsSurgeryHook(boil.BeforeUpdateHook, dopsEvaluationsSurgeryBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	dopsEvaluationsSurgeryBeforeUpdateHooks = []DopsEvaluationsSurgeryHook{}

	AddDopsEvaluationsSurgeryHook(boil.AfterUpdateHook, dopsEvaluationsSurgeryAfterUpdateHook)
	if err = o.doAfterUpdateHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	dopsEvaluationsSurgeryAfterUpdateHooks = []DopsEvaluationsSurgeryHook{}

	AddDopsEvaluationsSurgeryHook(boil.BeforeDeleteHook, dopsEvaluationsSurgeryBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	dopsEvaluationsSurgeryBeforeDeleteHooks = []DopsEvaluationsSurgeryHook{}

	AddDopsEvaluationsSurgeryHook(boil.AfterDeleteHook, dopsEvaluationsSurgeryAfterDeleteHook)
	if err = o.doAfterDeleteHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	dopsEvaluationsSurgeryAfterDeleteHooks = []DopsEvaluationsSurgeryHook{}

	AddDopsEvaluationsSurgeryHook(boil.BeforeUpsertHook, dopsEvaluationsSurgeryBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	dopsEvaluationsSurgeryBeforeUpsertHooks = []DopsEvaluationsSurgeryHook{}

	AddDopsEvaluationsSurgeryHook(boil.AfterUpsertHook, dopsEvaluationsSurgeryAfterUpsertHook)
	if err = o.doAfterUpsertHooks(ctx, nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	dopsEvaluationsSurgeryAfterUpsertHooks = []DopsEvaluationsSurgeryHook{}
}

func testDopsEvaluationsSurgeriesInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &DopsEvaluationsSurgery{}
	if err = randomize.Struct(seed, o, dopsEvaluationsSurgeryDBTypes, true, dopsEvaluationsSurgeryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DopsEvaluationsSurgery struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := DopsEvaluationsSurgeries().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testDopsEvaluationsSurgeriesInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &DopsEvaluationsSurgery{}
	if err = randomize.Struct(seed, o, dopsEvaluationsSurgeryDBTypes, true); err != nil {
		t.Errorf("Unable to randomize DopsEvaluationsSurgery struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(dopsEvaluationsSurgeryColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := DopsEvaluationsSurgeries().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testDopsEvaluationsSurgeryToOneDopsEvaluationUsingDopsEvaluation(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local DopsEvaluationsSurgery
	var foreign DopsEvaluation

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, dopsEvaluationsSurgeryDBTypes, false, dopsEvaluationsSurgeryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DopsEvaluationsSurgery struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, dopsEvaluationDBTypes, false, dopsEvaluationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DopsEvaluation struct: %s", err)
	}

	if err := foreign.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	local.DopsEvaluationID = foreign.ID
	if err := local.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.DopsEvaluation().One(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := DopsEvaluationsSurgerySlice{&local}
	if err = local.L.LoadDopsEvaluation(ctx, tx, false, (*[]*DopsEvaluationsSurgery)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.DopsEvaluation == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.DopsEvaluation = nil
	if err = local.L.LoadDopsEvaluation(ctx, tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.DopsEvaluation == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testDopsEvaluationsSurgeryToOneSurgeryUsingSurgery(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local DopsEvaluationsSurgery
	var foreign Surgery

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, dopsEvaluationsSurgeryDBTypes, false, dopsEvaluationsSurgeryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DopsEvaluationsSurgery struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, surgeryDBTypes, false, surgeryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Surgery struct: %s", err)
	}

	if err := foreign.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	local.SurgeryID = foreign.ID
	if err := local.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.Surgery().One(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := DopsEvaluationsSurgerySlice{&local}
	if err = local.L.LoadSurgery(ctx, tx, false, (*[]*DopsEvaluationsSurgery)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Surgery == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Surgery = nil
	if err = local.L.LoadSurgery(ctx, tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Surgery == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testDopsEvaluationsSurgeryToOneSetOpDopsEvaluationUsingDopsEvaluation(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a DopsEvaluationsSurgery
	var b, c DopsEvaluation

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, dopsEvaluationsSurgeryDBTypes, false, strmangle.SetComplement(dopsEvaluationsSurgeryPrimaryKeyColumns, dopsEvaluationsSurgeryColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, dopsEvaluationDBTypes, false, strmangle.SetComplement(dopsEvaluationPrimaryKeyColumns, dopsEvaluationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, dopsEvaluationDBTypes, false, strmangle.SetComplement(dopsEvaluationPrimaryKeyColumns, dopsEvaluationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*DopsEvaluation{&b, &c} {
		err = a.SetDopsEvaluation(ctx, tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.DopsEvaluation != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.DopsEvaluationsSurgeries[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.DopsEvaluationID != x.ID {
			t.Error("foreign key was wrong value", a.DopsEvaluationID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.DopsEvaluationID))
		reflect.Indirect(reflect.ValueOf(&a.DopsEvaluationID)).Set(zero)

		if err = a.Reload(ctx, tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.DopsEvaluationID != x.ID {
			t.Error("foreign key was wrong value", a.DopsEvaluationID, x.ID)
		}
	}
}
func testDopsEvaluationsSurgeryToOneSetOpSurgeryUsingSurgery(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a DopsEvaluationsSurgery
	var b, c Surgery

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, dopsEvaluationsSurgeryDBTypes, false, strmangle.SetComplement(dopsEvaluationsSurgeryPrimaryKeyColumns, dopsEvaluationsSurgeryColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, surgeryDBTypes, false, strmangle.SetComplement(surgeryPrimaryKeyColumns, surgeryColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, surgeryDBTypes, false, strmangle.SetComplement(surgeryPrimaryKeyColumns, surgeryColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*Surgery{&b, &c} {
		err = a.SetSurgery(ctx, tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Surgery != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.DopsEvaluationsSurgeries[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.SurgeryID != x.ID {
			t.Error("foreign key was wrong value", a.SurgeryID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.SurgeryID))
		reflect.Indirect(reflect.ValueOf(&a.SurgeryID)).Set(zero)

		if err = a.Reload(ctx, tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.SurgeryID != x.ID {
			t.Error("foreign key was wrong value", a.SurgeryID, x.ID)
		}
	}
}

func testDopsEvaluationsSurgeriesReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &DopsEvaluationsSurgery{}
	if err = randomize.Struct(seed, o, dopsEvaluationsSurgeryDBTypes, true, dopsEvaluationsSurgeryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DopsEvaluationsSurgery struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = o.Reload(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testDopsEvaluationsSurgeriesReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &DopsEvaluationsSurgery{}
	if err = randomize.Struct(seed, o, dopsEvaluationsSurgeryDBTypes, true, dopsEvaluationsSurgeryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DopsEvaluationsSurgery struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := DopsEvaluationsSurgerySlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testDopsEvaluationsSurgeriesSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &DopsEvaluationsSurgery{}
	if err = randomize.Struct(seed, o, dopsEvaluationsSurgeryDBTypes, true, dopsEvaluationsSurgeryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DopsEvaluationsSurgery struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := DopsEvaluationsSurgeries().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	dopsEvaluationsSurgeryDBTypes = map[string]string{`ID`: `uuid`, `SurgeryID`: `uuid`, `DopsEvaluationID`: `uuid`}
	_                             = bytes.MinRead
)

func testDopsEvaluationsSurgeriesUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(dopsEvaluationsSurgeryPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(dopsEvaluationsSurgeryAllColumns) == len(dopsEvaluationsSurgeryPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &DopsEvaluationsSurgery{}
	if err = randomize.Struct(seed, o, dopsEvaluationsSurgeryDBTypes, true, dopsEvaluationsSurgeryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DopsEvaluationsSurgery struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := DopsEvaluationsSurgeries().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, dopsEvaluationsSurgeryDBTypes, true, dopsEvaluationsSurgeryPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize DopsEvaluationsSurgery struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testDopsEvaluationsSurgeriesSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(dopsEvaluationsSurgeryAllColumns) == len(dopsEvaluationsSurgeryPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &DopsEvaluationsSurgery{}
	if err = randomize.Struct(seed, o, dopsEvaluationsSurgeryDBTypes, true, dopsEvaluationsSurgeryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize DopsEvaluationsSurgery struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := DopsEvaluationsSurgeries().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, dopsEvaluationsSurgeryDBTypes, true, dopsEvaluationsSurgeryPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize DopsEvaluationsSurgery struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(dopsEvaluationsSurgeryAllColumns, dopsEvaluationsSurgeryPrimaryKeyColumns) {
		fields = dopsEvaluationsSurgeryAllColumns
	} else {
		fields = strmangle.SetComplement(
			dopsEvaluationsSurgeryAllColumns,
			dopsEvaluationsSurgeryPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	typ := reflect.TypeOf(o).Elem()
	n := typ.NumField()

	updateMap := M{}
	for _, col := range fields {
		for i := 0; i < n; i++ {
			f := typ.Field(i)
			if f.Tag.Get("boil") == col {
				updateMap[col] = value.Field(i).Interface()
			}
		}
	}

	slice := DopsEvaluationsSurgerySlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testDopsEvaluationsSurgeriesUpsert(t *testing.T) {
	t.Parallel()

	if len(dopsEvaluationsSurgeryAllColumns) == len(dopsEvaluationsSurgeryPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := DopsEvaluationsSurgery{}
	if err = randomize.Struct(seed, &o, dopsEvaluationsSurgeryDBTypes, true); err != nil {
		t.Errorf("Unable to randomize DopsEvaluationsSurgery struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert DopsEvaluationsSurgery: %s", err)
	}

	count, err := DopsEvaluationsSurgeries().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, dopsEvaluationsSurgeryDBTypes, false, dopsEvaluationsSurgeryPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize DopsEvaluationsSurgery struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert DopsEvaluationsSurgery: %s", err)
	}

	count, err = DopsEvaluationsSurgeries().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}