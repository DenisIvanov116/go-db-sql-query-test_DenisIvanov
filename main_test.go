package main

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestSelectClientWhenOk(t *testing.T) {
	db, err := sql.Open("sqlite", "demo.db")
	require.NoError(t, err)
	defer db.Close()

	clientID := 1
	cl, err := selectClient(db, clientID)
	require.NoError(t, err)

	assert.Equal(t, clientID, cl.ID)
	assert.NotEmpty(t, cl.FIO)
	assert.NotEmpty(t, cl.Login)
	assert.NotEmpty(t, cl.Birthday)
	assert.NotEmpty(t, cl.Email)
}

func TestSelectClientWhenNoClient(t *testing.T) {
	db, err := sql.Open("sqlite", "demo.db")
	require.NoError(t, err)
	defer db.Close()

	missingID := -1
	cl, err := selectClient(db, missingID)
	require.Equal(t, sql.ErrNoRows, err)

	assert.Zero(t, cl.ID)
	assert.Empty(t, cl.FIO)
	assert.Empty(t, cl.Login)
	assert.Empty(t, cl.Birthday)
	assert.Empty(t, cl.Email)
}

func TestInsertClientThenSelectAndCheck(t *testing.T) {
	db, err := sql.Open("sqlite", "demo.db")
	require.NoError(t, err)
	defer db.Close()

	newClient := Client{
		FIO:      "Test User",
		Login:    "testuser",
		Birthday: "19900101",
		Email:    "test@example.com",
	}
	id, err := insertClient(db, newClient)
	require.NoError(t, err)
	require.NotZero(t, id)

	stored, err := selectClient(db, id)
	require.NoError(t, err)
	assert.Equal(t, newClient.FIO, stored.FIO)
	assert.Equal(t, newClient.Login, stored.Login)
	assert.Equal(t, newClient.Birthday, stored.Birthday)
	assert.Equal(t, newClient.Email, stored.Email)
}

func TestInsertClientDeleteClientThenCheck(t *testing.T) {
	db, err := sql.Open("sqlite", "demo.db")
	require.NoError(t, err)
	defer db.Close()

	cl := Client{
		FIO:      "To Be Deleted",
		Login:    "tobedeleted",
		Birthday: "20000101",
		Email:    "del@example.com",
	}
	id, err := insertClient(db, cl)
	require.NoError(t, err)
	require.NotZero(t, id)

	// проверяем, что клиент есть
	_, err = selectClient(db, id)
	require.NoError(t, err)

	// удаляем и проверяем, что его больше нет
	require.NoError(t, deleteClient(db, id))
	_, err = selectClient(db, id)
	require.Equal(t, sql.ErrNoRows, err)
}
