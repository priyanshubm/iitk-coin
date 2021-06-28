package jwtTokens

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/priyanshubm/iitk-coin/database"
)

var Options = sql.TxOptions{
	Isolation: sql.LevelSerializable,
}

func WriteUserToDb(name string, rollno string, password string, account_type string) error {
	statement, _ :=
		database.Db.Prepare("CREATE TABLE IF NOT EXISTS user (name TEXT,rollno TEXT PRIMARY KEY,password TEXT)")

	statement.Exec()

	statement, _ =
		database.Db.Prepare("INSERT INTO user (name,rollno,password,account_type) VALUES (?, ?, ?, ?)")
	_, err := statement.Exec(name, rollno, password, account_type)
	if err != nil {
		return err
	}
	err = InitializeCoins(rollno)
	if err != nil {
		return err
	}

	return nil

}

func InitializeCoins(rollno string) error {
	Db, _ :=
		sql.Open("sqlite3", "./database/user.db")

	statement, _ :=
		Db.Prepare(`INSERT INTO bank (rollno,coins) VALUES ($1,$2); `)

	_, err := statement.Exec(rollno, 0.00)
	if err != nil {
		return err
	}

	return nil

}

func WriteCoinsToDb(rollno string, numberOfCoins string, remarks string) (error, string) { // cpnvert this into a transaction

	coins_number, e := strconv.ParseFloat(numberOfCoins, 32)
	if e != nil {
		return e, "Coins not valid "
	}
	_, _, err := GetUserFromRollNo(rollno)
	if err != nil {
		return err, "User not present "
	}

	tx, _ := database.Db.BeginTx(context.Background(), &Options)

	res, execErr := tx.Exec(`UPDATE bank SET coins = coins + ? WHERE rollno= ? AND coins + ?<= ?;`, coins_number, rollno, coins_number, database.MaxCoins)
	rowsAffected, _ := res.RowsAffected()
	if execErr != nil || rowsAffected != 1 {
		_ = tx.Rollback()
		if execErr != nil {
			return execErr, ""
		}
		overflowError := errors.New("Balance cannot exceed " + fmt.Sprintf("%f", database.MaxCoins))
		return overflowError, ""
	}
	_, err = tx.Exec(`INSERT INTO rewards (user,amount,remarks,time) VALUES (?,?,?,?)`, rollno, coins_number, remarks, time.Now())
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		return err, "Some error occured in the transaction, please try again later "
	}

	if err = tx.Commit(); err != nil {
		return err, ""
	}

	return nil, "Coins added sucessfully "
}

func TransferCoinDb(firstRollno string, secondRollno string, transferAmount float64) (error, float64) {
	if firstRollno == secondRollno {
		return nil, 0
	}
	_, _, err := GetUserFromRollNo(firstRollno)
	if err != nil {
		return errors.New("user " + firstRollno + " not present "), 0
	}
	_, _, err = GetUserFromRollNo(secondRollno)
	if err != nil {
		return errors.New("user " + secondRollno + " not present "), 0
	}

	var options = sql.TxOptions{
		Isolation: sql.LevelSerializable,
	}
	tx, err := database.Db.BeginTx(context.Background(), &options)
	if err != nil {
		_ = tx.Rollback()
		log.Fatal(err)
		return err, 0
	}

	batch1 := firstRollno[0:3]
	batch2 := secondRollno[0:3]
	var taxRate float32 = 0.02
	if batch1 != batch2 {
		taxRate = 0.33
	}
	taxAmount := taxRate * float32(transferAmount)
	res, execErr := tx.Exec("UPDATE bank SET coins = coins - (?+?) WHERE rollno=? AND  coins - (?+?) >= 0 ", transferAmount, taxAmount, firstRollno, transferAmount, taxAmount)
	rowsAffected, _ := res.RowsAffected()
	if execErr != nil || rowsAffected != 1 {
		_ = tx.Rollback()
		if execErr != nil {
			return err, 0
		}

		balanceError := errors.New("not enough balance  ")
		return balanceError, 0

	}

	res, execErr = tx.Exec("UPDATE bank SET coins = coins + ? WHERE rollno=? AND coins + ? <= ?", transferAmount, secondRollno, transferAmount, database.MaxCoins)

	rowsAffected, _ = res.RowsAffected()
	if execErr != nil || rowsAffected != 1 {
		_ = tx.Rollback()
		if execErr != nil {
			return execErr, 0
		}
		overflowError := errors.New("Balance cannot exceed " + fmt.Sprintf("%f", database.MaxCoins))
		return overflowError, 0
	}

	_, execErr = tx.Exec(`INSERT INTO transfers (TransferFrom,TransferTo,amount,tax,time) VALUES (?,?,?,?,?)`, firstRollno, secondRollno, transferAmount, taxAmount, time.Now())
	if execErr != nil {
		_ = tx.Rollback()
		return execErr, 0
	}
	if err = tx.Commit(); err != nil {
		return err, 0
	}

	return nil, float64(taxAmount)
}

func RedeemCoinsDb(roll_no string, item_id int) (float64, error) { //convert this into a transaction and add eror handling
	// Check if the item id is valid and obtain the coist of the item

	numEvents, _ := GetNumEvents(roll_no)
	if numEvents < database.MinEvents {
		return 0, errors.New("You need to participate in at least " + strconv.Itoa(database.MinEvents) + " events to clam a reward ")
	}
	cost, available, err := getItemFromId(item_id)
	if err != nil {
		return 0, err
	}
	if available == 0 {
		return 0, errors.New("item not available, please select another item ")
	}

	_, _, err = GetUserFromRollNo(roll_no)
	if err != nil {
		return 0, errors.New("user " + roll_no + " not present ")
	}
	var options = sql.TxOptions{
		Isolation: sql.LevelSerializable,
	}
	tx, err := database.Db.BeginTx(context.Background(), &options)
	if err != nil {
		_ = tx.Rollback()
		log.Fatal(err)
		return 0, err
	}
	res, err := tx.Exec(`UPDATE bank SET coins = coins - ? WHERE rollno= ? AND coins - ? >=0 `, cost, roll_no, cost)
	rowsAffected, _ := res.RowsAffected()
	if err != nil || rowsAffected != 1 {
		tx.Rollback()
		if err != nil {
			return 0, err
		}
		return 0, errors.New("insufficient coins to claim this item ")
	}
	res, err = tx.Exec(`UPDATE items SET available = available -1 WHERE id = ? `, item_id)
	rowsAffected, _ = res.RowsAffected()
	if err != nil || rowsAffected != 1 {
		tx.Rollback()
		if err != nil {
			return 0, err
		}
		return 0, errors.New("error occured while transaction please try later ")
	}
	_, err = tx.Exec(`INSERT INTO redeems (user,item,time) VALUES (?,?,?)`, roll_no, item_id, time.Now())
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		return 0, err
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	coins, _ := GetCoinsFromRollNo(roll_no)

	return coins, err
}

func WriteItemsToDb(item_id int, cost string, number int) (string, error) { // cpnvert this into a transaction

	cost_number, e := strconv.ParseFloat(cost, 32)
	if e != nil {
		return "Coins not valid ", e
	}

	tx, _ := database.Db.BeginTx(context.Background(), &Options)

	_, err := tx.Exec(`INSERT INTO items (id,cost,available) VALUES (?,?,?) ON CONFLICT(id) DO UPDATE SET available = available + ? ;`, item_id, cost_number, number, number)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		return "Some error occured in the transaction, please try again later ", err
	}

	if err = tx.Commit(); err != nil {
		return "", err
	}
	return "success", e

}
