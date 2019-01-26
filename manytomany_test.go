package otira

import (
	//"log"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

func TestManyToManyCreateJoinTable(t *testing.T) {
	pers, team, person, _, err := simpleManyToMany()

	if err != nil {
		t.Fatal(err)
	}

	err = pers.CreateTables(team, person)
	if err != nil {
		t.Fatal(err)
	}
	err = pers.Done()
	if err != nil {
		t.Fatal(err)
	}
}

func TestManyToManyInsertSingleRecordWithNoRelationRecord(t *testing.T) {
	pers, team, person, _, err := simpleManyToMany()

	if err != nil {
		t.Fatal(err)
	}

	err = pers.CreateTables(team, person)
	if err != nil {
		t.Fatal(err)
	}

	teamRecord, err := makeTeamRecord(team, 44, "Leafs")
	if err != nil {
		t.Fatal(err)
	}
	err = pers.Save(teamRecord)
	if err != nil {
		t.Fatal(err)
	}
	err = pers.Done()
	if err != nil {
		t.Fatal(err)
	}
}

func TestManyToManyInsertSingleRecordWithOneRelationRecord(t *testing.T) {
	pers, team, person, m2m, err := simpleManyToMany()

	if err != nil {
		t.Fatal(err)
	}

	err = pers.CreateTables(team, person)
	if err != nil {
		t.Fatal(err)
	}
	teamRecord, err := makeTeamRecord(team, 44, "Leafs")
	if err != nil {
		t.Fatal(err)
	}

	personRecord, err := makePersonRecord(person, 323, "Bill Smith")

	teamRecord.AddRelationRecord(m2m, personRecord)
	if err != nil {
		t.Fatal(err)
	}

	err = pers.Save(teamRecord)
	if err != nil {
		t.Fatal(err)
	}
	err = pers.Done()
	if err != nil {
		t.Fatal(err)
	}
}

func TestManyToManyInsertTwoRecordsWithSameRelationRecord(t *testing.T) {
	pers, team, person, m2m, err := simpleManyToMany()

	if err != nil {
		t.Fatal(err)
	}

	err = pers.CreateTables(team, person)
	if err != nil {
		t.Fatal(err)
	}
	teamRecord1, err := makeTeamRecord(team, 44, "Leafs")
	if err != nil {
		t.Fatal(err)
	}

	teamRecord2, err := makeTeamRecord(team, 1, "Canadiens")
	if err != nil {
		t.Fatal(err)
	}

	personRecord, err := makePersonRecord(person, 323, "Bill Smith")

	teamRecord1.AddRelationRecord(m2m, personRecord)
	if err != nil {
		t.Fatal(err)
	}

	err = pers.Save(teamRecord1)
	if err != nil {
		t.Fatal(err)
	}

	teamRecord2.AddRelationRecord(m2m, personRecord)
	if err != nil {
		t.Fatal(err)
	}

	err = pers.Save(teamRecord2)
	if err != nil {
		t.Fatal(err)
	}
	err = pers.Done()
	if err != nil {
		t.Fatal(err)
	}
}

func TestManyToManyInsertOneRecordWithTwoDifferentRelationRecords(t *testing.T) {
	pers, team, person, m2m, err := simpleManyToMany()

	if err != nil {
		t.Fatal(err)
	}

	err = pers.CreateTables(team, person)
	if err != nil {
		t.Fatal(err)
	}
	teamRecord, err := makeTeamRecord(team, 44, "Leafs")
	if err != nil {
		t.Fatal(err)
	}

	personRecord1, err := makePersonRecord(person, 323, "Bill Smith")
	if err != nil {
		t.Fatal(err)
	}

	teamRecord.AddRelationRecord(m2m, personRecord1)
	if err != nil {
		t.Fatal(err)
	}

	personRecord2, err := makePersonRecord(person, 8343, "Bobby Orr")
	if err != nil {
		t.Fatal(err)
	}

	teamRecord.AddRelationRecord(m2m, personRecord2)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("A1")
	err = pers.Save(teamRecord)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("A2")
	err = pers.Done()
	if err != nil {
		t.Fatal(err)
	}
}

func TestManyToManyInsertManyRecordsWithTwoRelationRecords(t *testing.T) {
	pers, team, person, m2m, err := simpleManyToMany()

	if err != nil {
		t.Fatal(err)
	}

	err = pers.CreateTables(team, person)
	if err != nil {
		t.Fatal(err)
	}

	personRecord1, err := makePersonRecord(person, uint64(2), "Bill Smith")
	if err != nil {
		t.Fatal(err)
	}

	for i := 11; i < 1000; i++ {

		q := toString(i)
		teamRecord, err := makeTeamRecord(team, uint64(i), "Leafs_"+q)
		if err != nil {
			t.Fatal(err)
		}

		teamRecord.AddRelationRecord(m2m, personRecord1)
		if err != nil {
			t.Fatal(err)
		}

		personRecord2, err := makePersonRecord(person, uint64(i*30000), "Bobby Orr"+q)
		if err != nil {
			t.Fatal(err)
		}

		teamRecord.AddRelationRecord(m2m, personRecord2)
		if err != nil {
			t.Fatal(err)
		}

		err = pers.Save(teamRecord)
		if err != nil {
			t.Fatal(err)
		}
	}
	err = pers.Done()
	if err != nil {
		t.Fatal(err)
	}
}

func TestManyToManyInsertTwoRecordsWithDifferentRelationRecord(t *testing.T) {
	pers, team, person, _, err := simpleManyToMany()

	if err != nil {
		t.Fatal(err)
	}

	err = pers.CreateTables(team, person)
	if err != nil {
		t.Fatal(err)
	}
	// TODO
	//t.Fatal(err)
	err = pers.Done()
	if err != nil {
		t.Fatal(err)
	}
}

const TableNamePerson = "person"

func simpleManyToMany() (*Persister, *TableDef, *TableDef, *ManyToMany, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	//db, err := sql.Open("sqlite3", "db.sqlite3")

	if err != nil {
		return nil, nil, nil, nil, err
	}
	pers, err := NewPersister(db, NewDialectSqlite3(nil, false))
	if err != nil {
		return nil, nil, nil, nil, err
	}

	team, person, m2m, err := newManyToManyDefaultTables()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return pers, team, person, m2m, nil
}

func newManyToManyDefaultTables() (*TableDef, *TableDef, *ManyToMany, error) {
	personTable, err := makePersonTable()
	if err != nil {
		return nil, nil, nil, err
	}
	teamTable, err := makeTeamTable()
	if err != nil {
		return nil, nil, nil, err
	}

	m2m := NewManyToMany()
	m2m.LeftTable = teamTable
	m2m.RightTable = personTable
	teamTable.AddManyToMany(m2m)

	return teamTable, personTable, m2m, err

}

func makePersonTable() (*TableDef, error) {
	personTable, err := NewTableDef(TableNamePerson)
	if err != nil {
		return nil, err
	}
	personTable.UseRecordPrimaryKeys = true
	id := new(FieldDefUint64)
	id.SetName(pk)
	id.SetUnique(true)
	err = personTable.Add(id)
	if err != nil {
		return nil, err
	}

	nameField := new(FieldDefString)
	nameField.SetName(NAME)
	nameField.SetFixed(true)
	nameField.SetLength(24)
	err = personTable.Add(nameField)
	if err != nil {
		return nil, err
	}
	personTable.SetJoinDiscrimFields(nameField)

	personTable.SetDone()

	return personTable, nil
}

func makePersonRecord(t *TableDef, id uint64, name string) (*Record, error) {
	rec, err := t.NewRecord()
	if err != nil {
		return nil, err
	}

	err = rec.SetByName(pk, id)
	if err != nil {
		return nil, err
	}
	rec.SetByName(NAME, name)

	return rec, nil
}

const TableNameTeam = "team"
const TeamNameField = "team_name"

func makeTeamTable() (*TableDef, error) {
	teamTable, err := NewTableDef(TableNameTeam)
	if err != nil {
		return nil, err
	}
	id := new(FieldDefUint64)
	id.SetName(pk)
	id.SetUnique(true)
	err = teamTable.Add(id)
	if err != nil {
		return nil, err
	}

	nameField := new(FieldDefString)
	nameField.SetName(TeamNameField)
	nameField.SetFixed(true)
	nameField.SetLength(24)
	err = teamTable.Add(nameField)
	if err != nil {
		return nil, err
	}
	teamTable.SetDone()

	return teamTable, nil
}

func makeTeamRecord(t *TableDef, id uint64, name string) (*Record, error) {
	rec, err := t.NewRecord()
	if err != nil {
		return nil, err
	}

	err = rec.SetByName(pk, id)
	if err != nil {
		return nil, err
	}
	rec.SetByName(TeamNameField, name)

	return rec, nil
}
