package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type Category struct {
	Catname string
	Sender string
	Receiver string
}

type Task struct {
	Catname    string
	Sender     string
	Taskname   string
	Duedate    string
	Status     string
	Upredicted string
	Spredicted string
	Current    string
}

var Port = ":5558"

func main() {

	http.HandleFunc("/", ServeFiles)
	fmt.Println("Serving @ : ", "http://127.0.0.1"+Port)

	log.Fatal(http.ListenAndServe(Port, nil))
}

func ServeFiles(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "GET":

		path := r.URL.Path

		fmt.Println(path)

		if path == "/" {

			path = "210603_CategoriesTasks.html"
		} else {

			path = "." + path
		}

		http.ServeFile(w, r, path)

	case "POST":

		r.ParseMultipartForm(0)

		message := r.FormValue("accAction")

		//Taking the full string and opening the database
		database, _ := sql.Open("sqlite3", "./Usersdb.db")
		result := strings.SplitAfter(message, "\x1F")
		action := wordTrim(result[3])
		if action == "login" {
			//********************************************************************
			//         USER IS LOGGING IN
			//********************************************************************

			username := wordTrim(result[1])
			password := wordTrim(result[2])
			token, err := GenerateRandomStringURLSafe(32)
			token = "1" + username + token
			timestamp := getTime()
			change, err := getPassword(database, username, password, token, timestamp)
			if change == "1" {
				change = token
				fmt.Fprintf(w, change)
				// respond to client's request
				//    revMsg := Reverse(message)
				//		fmt.Fprintf(w, "Server: %s \n", revMsg+ " | " + time.Now().Format(time.RFC3339))
			}
			if err != nil {
				log.Println("Write failed")
				log.Fatal(err)
			}
			println(change)
			// Display our results.

		}
		if action == "create" {
			//********************************************************************
			//         USER IS CREATING ACCOUNT
			//********************************************************************

			username := wordTrim(result[1])
			existing, err := getUser(database, username)
			if err != nil {
				log.Println("Write failed")
				log.Fatal(err)
			}
			if existing != "" {
			fmt.Fprintf(w, "2")
			return
		}
		password := wordTrim(result[2])

		hash := []byte(password)
		hashpwd := hashAndSalt(hash)
		change := addUser(database, username, hashpwd)
		fmt.Fprintf(w, change)
		return

		 	} else if action == "addCat" {
			//********************************************************************
			//         USER IS ADDING CATEGORY
			//********************************************************************

			username := wordTrim(result[1])
			catID := wordTrim(result[2])
			existing, err := getCat(database, username, catID)
			if err != nil {
				log.Println("Write failed")
				log.Fatal(err)
			}
			if existing != "" {
				fmt.Fprintf(w, "2")
				return
			}
			change := addCat(database, username, catID)
			change = getCatTable(database, username)
			print("Change: ", change)
			fmt.Fprintf(w, change)
		}  else if action == "addExportCat" {
			//********************************************************************
			//         USER IS ADDING CATEGORY
			//********************************************************************

			username := wordTrim(result[1])
			catID := wordTrim(result[2])
			token := wordTrim(result[4])
			timestamp := getTime()
			check := checkTime(database, username, timestamp)
			if check == 1 {
				check2, err := getToken(database, username, token)
				if err != nil {
					log.Println("Write failed")
					log.Fatal(err)
				}
				if check2 == "1" {
			existing, err := getExportCat(database, username, catID)
			if err != nil {
				log.Println("Write failed")
				log.Fatal(err)
			}
			if existing != "" {
				fmt.Fprintf(w, "2")
				return
			}
			for i := 5; i < len(result); i++ {
				receiver := wordTrim(result[i])
				addExportCat(database, username, catID, receiver)
		}
			change := getShareCatTable(database, username)
			print("Change: ", change)
			fmt.Fprintf(w, change)
			} else {
				fmt.Fprintf(w, "2")
				// Display our results.
			}
		} else {
			fmt.Fprintf(w, "2")
		}
		} else if action == "chpwd" {
			//********************************************************************
			//         USER IS CHANGING PASSWORD
			//********************************************************************
			username := wordTrim(result[1])
			password := wordTrim(result[2])
			token := wordTrim(result[4])
			timestamp := getTime()
			check := checkTime(database, username, timestamp)
			if check == 1 {
				check2, err := getToken(database, username, token)
				if err != nil {
					log.Println("Write failed")
					log.Fatal(err)
				}
				if check2 == "1" {
					hash := []byte(password)
					hashpwd := hashAndSalt(hash)
					change := modifyPassword(database, username, hashpwd)
					fmt.Fprintf(w, change)
				} else {
					fmt.Fprintf(w, "2")
					// Display our results.
				}
			} else {
				fmt.Fprintf(w, "2")
			}
		} else if action == "lgout" {
			//********************************************************************
			//         USER IS LOGGING OUT
			//********************************************************************

			username := wordTrim(result[1])
			change := dltToken(database, username)
			fmt.Fprintf(w, change)

			// Display our results.

		} else if action == "catab" {
			//********************************************************************
			//         USER IS CHECKING CATEGORY TABLE
			//********************************************************************

			username := wordTrim(result[1])
			token := wordTrim(result[4])
			timestamp := getTime()
			check := checkTime(database, username, timestamp)
			if check == 1 {
				check2, err := getToken(database, username, token)
				if err != nil {
					log.Println("Write failed")
					log.Fatal(err)
				}
				if check2 == "1" {
					print("We Made it here")
					change := getCatTable(database, username)
					print("gettingCatTable")
					fmt.Fprintf(w, change)
				} else {
					fmt.Fprintf(w, "2")
					// Display our results.
				}
			} else {
				fmt.Fprintf(w, "2")
			}
			} else if action == "pendingcatab" {
				//********************************************************************
				//         USER IS CHECKING PENDING CATEGORY TABLE
				//********************************************************************

				username := wordTrim(result[1])
				token := wordTrim(result[4])
				timestamp := getTime()
				check := checkTime(database, username, timestamp)
				if check == 1 {
					check2, err := getToken(database, username, token)
					if err != nil {
						log.Println("Write failed")
						log.Fatal(err)
					}
					if check2 == "1" {
						change := getPendingCatTable(database, username)
						fmt.Fprintf(w, change)
					} else {
						fmt.Fprintf(w, "2")
						// Display our results.
					}
				} else {
					fmt.Fprintf(w, "2")
				}
		} else if action == "sharecatab" {
			//********************************************************************
			//         USER IS CHECKING CATEGORY TABLE
			//********************************************************************

			username := wordTrim(result[1])
			token := wordTrim(result[2])
			timestamp := getTime()
			check := checkTime(database, username, timestamp)
			if check == 1 {
				check2, err := getToken(database, username, token)
				if err != nil {
					log.Println("Write failed")
					log.Fatal(err)
				}
				if check2 == "1" {
					change := getShareCatTable(database, username)
					fmt.Fprintf(w, change)
				} else {
					fmt.Fprintf(w, "2")
					// Display our results.
				}
			} else {
				fmt.Fprintf(w, "2")
			}
		} else if action == "exportcatab" {
			//********************************************************************
			//         USER IS CHECKING CATEGORY TABLE
			//********************************************************************

			username := wordTrim(result[1])
			token := wordTrim(result[2])
			timestamp := getTime()
			check := checkTime(database, username, timestamp)
			if check == 1 {
				check2, err := getToken(database, username, token)
				if err != nil {
					log.Println("Write failed")
					log.Fatal(err)
				}
				if check2 == "1" {
					change := getExportCatTable(database, username)
					fmt.Fprintf(w, change)
				} else {
					fmt.Fprintf(w, "2")
					// Display our results.
				}
			} else {
				fmt.Fprintf(w, "2")
			}
		} else if action == "pause" {
			//********************************************************************
			//         USER IS UPDATING TIME
			//********************************************************************
			username := wordTrim(result[1])
			catID := wordTrim(result[2])
			token := wordTrim(result[4])
			taskID := wordTrim(result[5])
			time := wordTrim(result[6])
			timestamp := getTime()
			check := checkTime(database, username, timestamp)
			if check == 1 {
				check2, err := getToken(database, username, token)
				if err != nil {
					log.Println("Write failed")
					log.Fatal(err)
				}
				if check2 == "1" {
					fmt.Println("catID: ", catID, "taskID: ", taskID, "time: ", time)
					change := updateTime(database, username, catID, taskID, time)
					fmt.Fprintf(w, change)
				} else {
					fmt.Fprintf(w, "2")
					// Display our results.
				}
			} else {
				fmt.Fprintf(w, "2")
			}
			} else if action == "sharepause" {
				//********************************************************************
				//         USER IS UPDATING TIME
				//********************************************************************
				username := wordTrim(result[1])
				catID := wordTrim(result[2])
				token := wordTrim(result[4])
				taskID := wordTrim(result[5])
				time := wordTrim(result[6])
				sender := wordTrim(result[7])
				timestamp := getTime()
				check := checkTime(database, username, timestamp)
				if check == 1 {
					check2, err := getToken(database, username, token)
					if err != nil {
						log.Println("Write failed")
						log.Fatal(err)
					}
					if check2 == "1" {
						fmt.Println("catID: ", catID, "taskID: ", taskID, "time: ", time)
						change := updateShareTime(database, username, sender, catID, taskID, time)
						fmt.Fprintf(w, change)
					} else {
						fmt.Fprintf(w, "2")
						// Display our results.
					}
				} else {
					fmt.Fprintf(w, "2")
				}
		} else if action == "done" {
			//********************************************************************
			//         USER IS FINISHED
			//********************************************************************
			username := wordTrim(result[1])
			catID := wordTrim(result[2])
			token := wordTrim(result[4])
			taskID := wordTrim(result[5])
			time := wordTrim(result[6])
			timestamp := getTime()
			check := checkTime(database, username, timestamp)
			if check == 1 {
				check2, err := getToken(database, username, token)
				if err != nil {
					log.Println("Write failed")
					log.Fatal(err)
				}
				if check2 == "1" {
					fmt.Println("catID: ", catID, "taskID: ", taskID, "time: ", time)
					change := finishTime(database, username, catID, taskID, time)
					print("Change: ", change)
					fmt.Fprintf(w, change)
				} else {
					fmt.Fprintf(w, "2")
					// Display our results.
				}
			} else {
				fmt.Fprintf(w, "2")
			}
		} else if action == "sharedone" {
			//********************************************************************
			//         USER IS FINISHED
			//********************************************************************
			username := wordTrim(result[1])
			catID := wordTrim(result[2])
			token := wordTrim(result[4])
			taskID := wordTrim(result[5])
			time := wordTrim(result[6])
			sender := wordTrim(result[7])
			timestamp := getTime()
			check := checkTime(database, username, timestamp)
			if check == 1 {
				check2, err := getToken(database, username, token)
				if err != nil {
					log.Println("Write failed")
					log.Fatal(err)
				}
				if check2 == "1" {
					fmt.Println("catID: ", catID, "taskID: ", taskID, "time: ", time)
					change := finishShareTime(database, username, sender, catID, taskID, time)
					print("Change: ", change)
					fmt.Fprintf(w, change)
				} else {
					fmt.Fprintf(w, "2")
					// Display our results.
				}
			} else {
				fmt.Fprintf(w, "2")
			}
		} else if action == "delCat" {
			//********************************************************************
			//         USER IS DELETING CATEGORY
			//********************************************************************
			username := wordTrim(result[1])
			catID := wordTrim(result[2])
			token := wordTrim(result[4])
			timestamp := getTime()
			check := checkTime(database, username, timestamp)
			if check == 1 {
				check2, err := getToken(database, username, token)
				if err != nil {
					log.Println("Write failed")
					log.Fatal(err)
				}
				if check2 == "1" {
					change := delCat(database, username, catID)
					fmt.Fprintf(w, change)
				} else {
					fmt.Fprintf(w, "2")
					// Display our results.
				}
			} else {
				fmt.Fprintf(w, "2")
			}
		} else if action == "delTask" {
			//********************************************************************
			//         USER IS DELETING TASK
			//********************************************************************
			username := wordTrim(result[1])
			catID := wordTrim(result[2])
			token := wordTrim(result[4])
			taskID := wordTrim(result[5])
			timestamp := getTime()
			check := checkTime(database, username, timestamp)
			if check == 1 {
				check2, err := getToken(database, username, token)
				if err != nil {
					log.Println("Write failed")
					log.Fatal(err)
				}
				if check2 == "1" {
					change := delTask(database, username, catID, taskID)
					fmt.Fprintf(w, change)
				} else {
					fmt.Fprintf(w, "2")
					// Display our results.
				}
			} else {
				fmt.Fprintf(w, "2")
			}
		 	} else if action == "rejCat" {
					//********************************************************************
					//         USER IS REJECTING CATEGORY
					//********************************************************************
					username := wordTrim(result[1])
					catID := wordTrim(result[2])
					token := wordTrim(result[4])
					sender := wordTrim(result[5])
					timestamp := getTime()
					check := checkTime(database, username, timestamp)
					if check == 1 {
						check2, err := getToken(database, username, token)
						if err != nil {
							log.Println("Write failed")
							log.Fatal(err)
						}
						if check2 == "1" {
							change := rejCat(database, sender, catID, username)
							fmt.Fprintf(w, change)
						} else {
							fmt.Fprintf(w, "2")
							// Display our results.
						}
					} else {
						fmt.Fprintf(w, "2")
					}
					} else if action == "accCat" {
							//********************************************************************
							//         USER IS ACCEPTING CATEGORY
							//********************************************************************
							username := wordTrim(result[1])
							catID := wordTrim(result[2])
							token := wordTrim(result[4])
							sender := wordTrim(result[5])
							timestamp := getTime()
							check := checkTime(database, username, timestamp)
							if check == 1 {
								check2, err := getToken(database, username, token)
								if err != nil {
									log.Println("Write failed")
									log.Fatal(err)
								}
								if check2 == "1" {
									change := accCat(database, sender, catID, username)
									fmt.Fprintf(w, change)
								} else {
									fmt.Fprintf(w, "2")
									// Display our results.
								}
							} else {
								fmt.Fprintf(w, "2")
							}
		} else if action == "addTask" {
			//********************************************************************
			//         USER IS ADDING TASK
			//********************************************************************
			username := wordTrim(result[1])
			taskID := wordTrim(result[2])
			token := wordTrim(result[4])
			catID := wordTrim(result[5])
			dueDate := wordTrim(result[6])
			status := wordTrim(result[7])
			spredicted := wordTrim(result[8])
			upredicted := wordTrim(result[9])
			current := wordTrim(result[10])
			timestamp := getTime()
			check := checkTime(database, username, timestamp)
			if check == 1 {
				check2, err := getToken(database, username, token)
				if err != nil {
					log.Println("Write failed")
					log.Fatal(err)
				}
				if check2 == "1" {
					change := addTask(database, username, catID, taskID, dueDate, status, upredicted, spredicted, current)
					fmt.Fprintf(w, change)
				} else {
					fmt.Fprintf(w, "2")
					// Display our results.
				}
			} else {
				fmt.Fprintf(w, "2")
			}
		} else if action == "addShareTask" {
			//********************************************************************
			//         USER IS ADDING TASK
			//********************************************************************
			username := wordTrim(result[1])
			taskID := wordTrim(result[2])
			token := wordTrim(result[4])
			catID := wordTrim(result[5])
			dueDate := wordTrim(result[6])
			status := wordTrim(result[7])
			spredicted := wordTrim(result[8])
			upredicted := wordTrim(result[9])
			current := wordTrim(result[10])
			sender := wordTrim(result[11])
			timestamp := getTime()
			check := checkTime(database, username, timestamp)
			if check == 1 {
				check2, err := getToken(database, username, token)
				if err != nil {
					log.Println("Write failed")
					log.Fatal(err)
				}
				if check2 == "1" {
					change := addShareTask(database, username, sender, catID, taskID, dueDate, status, upredicted, spredicted, current)
					fmt.Fprintf(w, change)
				} else {
					fmt.Fprintf(w, "2")
					// Display our results.
				}
			} else {
				fmt.Fprintf(w, "2")
			}
			} else if action == "guess" {
			//********************************************************************
			//         USER IS REQUESTING A HINT
			//********************************************************************
			username := wordTrim(result[1])
			catID := wordTrim(result[2])
			token := wordTrim(result[4])
			upredicted := wordTrim(result[5])
			timestamp := getTime()
			check := checkTime(database, username, timestamp)
			if check == 1 {
				check2, err := getToken(database, username, token)
				if err != nil {
					log.Println("Write failed")
					log.Fatal(err)
				}
				if check2 == "1" {
					change := getHint(database, username, catID, upredicted)
					fmt.Fprintf(w, change)
				} else {
					fmt.Fprintf(w, "2")
					// Display our results.
				}
			} else {
				fmt.Fprintf(w, "2")
			}
		} else if action == "sharecatsk" {
			//********************************************************************
			//         USER IS CHECKING TASK TABLE
			//********************************************************************
			//    appid := wordTrim(result[0])
			username := wordTrim(result[1])
			catID := wordTrim(result[2])
			token := wordTrim(result[4])
			sender := wordTrim(result[5])
			timestamp := getTime()
			check := checkTime(database, username, timestamp)
			if check == 1 {
				check2, err := getToken(database, username, token)
				if err != nil {
					log.Println("Write failed")
					log.Fatal(err)
				}
				if check2 == "1" {
					change := getShareTaskTable(database, username, catID, sender)

					fmt.Fprintf(w, change)
				} else {
					fmt.Fprintf(w, "2")
					// Display our results.
				}
			} else {
				fmt.Fprintf(w, "2")
			}
			} else if action == "exportcatsk" {
				//********************************************************************
				//         USER IS CHECKING TASK TABLE
				//********************************************************************
				//    appid := wordTrim(result[0])
				username := wordTrim(result[1])
				catID := wordTrim(result[2])
				token := wordTrim(result[4])
				timestamp := getTime()
				check := checkTime(database, username, timestamp)
				if check == 1 {
					check2, err := getToken(database, username, token)
					if err != nil {
						log.Println("Write failed")
						log.Fatal(err)
					}
					if check2 == "1" {
						change := getExportTaskTable(database, username, catID)

						fmt.Fprintf(w, change)
					} else {
						fmt.Fprintf(w, "2")
						// Display our results.
					}
				} else {
					fmt.Fprintf(w, "2")
				}

		} else if action == "catsk" {
			//********************************************************************
			//         USER IS CHECKING TASK TABLE
			//********************************************************************
			//    appid := wordTrim(result[0])
			username := wordTrim(result[1])
			catID := wordTrim(result[2])
			token := wordTrim(result[4])
			timestamp := getTime()
			check := checkTime(database, username, timestamp)
			if check == 1 {
				check2, err := getToken(database, username, token)
				if err != nil {
					log.Println("Write failed")
					log.Fatal(err)
				}
				if check2 == "1" {
					change := getTaskTable(database, username, catID)

					fmt.Fprintf(w, change)
				} else {
					fmt.Fprintf(w, "2")
					// Display our results.
				}
			} else {
				fmt.Fprintf(w, "2")
			}

		}
	default:
		fmt.Fprintf(w, "Request type other than GET or POST not supported")
	}
}

//********************************************************************
func GenerateRandomStringURLSafe(n int) (string, error) {
	b, err := GenerateRandomBytes(n)
	return base64.URLEncoding.EncodeToString(b), err
}

//********************************************************************
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

//***********************************

func dltToken(db *sql.DB, name string) string {
	stmt, err := db.Prepare("DELETE FROM user3 WHERE username = ?")
	if err != nil {
		return "2"
	}

	stmt.Exec(name)
	stmt.Close()
	return "1"
}

//***********************************

func delTask(db *sql.DB, name string, catID string, taskID string) string {
	stmt, err := db.Prepare("DELETE FROM apoio WHERE username = ? AND catname = ? AND taskname = ?")
	if err != nil {
		return "2"
	}

	stmt.Exec(name, catID, taskID)
	stmt.Close()
	return "1"
}

//***********************************

func delCat(db *sql.DB, name string, catID string) string {
	stmt, err := db.Prepare("DELETE FROM apoio WHERE username = ? AND catname = ?")
	if err != nil {
		return "2"
	}

	stmt.Exec(name, catID)
	stmt.Close()
	stmt2, err := db.Prepare("DELETE FROM user2 WHERE username = ? AND catname = ?")
	if err != nil {
		return "2"
	}

	stmt2.Exec(name, catID)
	stmt2.Close()
	return "1"
}

//***********************************

func rejCat(db *sql.DB, sender string, catID string, name string) string {
	stmt, err := db.Prepare("DELETE FROM pending WHERE sender = ? AND catname = ? AND receiver = ?")
	if err != nil {
		return "2"
	}

	stmt.Exec(sender, catID, name)
	stmt.Close()
	return "1"
}
//***********************************

func accCat(db *sql.DB, sender string, catID string, name string) string {
	stmt, err := db.Prepare("DELETE FROM pending WHERE sender = ? AND catname = ? AND receiver = ?")
	if err != nil {
		return "2"
	}
	stmt.Exec(sender, catID, name)
	stmt.Close()
	stmt2, err := db.Prepare("INSERT INTO accepted (sender, catname, reciever) VALUES (?,?,?);")
	if err != nil {
		return "2"
	}
	stmt2.Exec(sender, catID, name)
	stmt2.Close()
	return "1"
}

//***********************************

func modifyPassword(db *sql.DB, name string, newPwd string) string {
	stmt, err := db.Prepare("UPDATE USER1 SET password = ? WHERE username = ?")
	if err != nil {
		return "2"
	}

	stmt.Exec(newPwd, name)
	stmt.Close()
	return "1"
}

//***********************************

func getPassword(db *sql.DB, name string, password string, token string, timestamp int64) (string, error) {
	var pwd string
	rtn := ""
	rows, err := db.Query("SELECT PASSWORD FROM USER1 WHERE USERNAME = ?", name)
	if err != nil {
		return "", err
	}
	rows.Next()
	rows.Scan(&pwd)
	rows.Close()
	err = bcrypt.CompareHashAndPassword([]byte(pwd), []byte(password))
	if err != nil {
		log.Println(err)
	} else {
		stmt1, err := db.Prepare("DELETE FROM user3 WHERE username = ?")
		if err != nil {
			return "2", err
		}

		stmt1.Exec(name)
		stmt1.Close()
		stmt, err := db.Prepare("INSERT INTO user3 (username, token, timestamp) VALUES (?,?,?);")
		if err != nil {
			return rtn, err
		}
		stmt.Exec(name, token, timestamp)
		stmt.Close()
		rtn = "1"
		if err != nil {
			return rtn, err
		}
	}
	err = bcrypt.CompareHashAndPassword([]byte(pwd), []byte(password))
	if err != nil {
		log.Println(err)
	} else {
		rtn = "1"
		if err != nil {

			return rtn, err
		}
	}
	return rtn, nil
}

//***********************************

func getUser(db *sql.DB, name string) (string, error) {
	var pwd string
	rows, err := db.Query("SELECT password FROM USER1 WHERE username = ?", name)
	if err != nil {
		return "", err
	}
	rows.Next()
	rows.Scan(&pwd)
	rows.Close()
	return pwd, nil
}
//***********************************

func getCat(db *sql.DB, name string, catID string) (string, error) {
	var pwd string
	rows, err := db.Query("SELECT catname FROM user2 WHERE username = ? and catname = ? ORDER BY catname ASC", name, catID)
	if err != nil {
		return "", err
	}
	rows.Next()
	rows.Scan(&pwd)
	rows.Close()
	return pwd, nil
}
//***********************************

func getExportCat(db *sql.DB, name string, catID string) (string, error) {
	var pwd string
	rows, err := db.Query("SELECT catname FROM accepted WHERE sender = ? and catname = ? ORDER BY catname ASC", name, catID)
	if err != nil {
		return "", err
	}
	rows.Next()
	rows.Scan(&pwd)
	rows.Close()
	return pwd, nil
}
//***********************************

func getToken(db *sql.DB, name string, token string) (string, error) {
	var pwd string
	rows, err := db.Query("SELECT token FROM user3 WHERE USERNAME = ?", name)
	if err != nil {
		return "", err
	}
	rows.Next()
	rows.Scan(&pwd)
	rows.Close()
	if pwd != token {
		return "2", nil
	}
	return "1", nil
}

//***********************************
func getCatTable(db *sql.DB, name string) string {

	rows, err := db.Query("SELECT DISTINCT catname FROM apoio WHERE USERNAME = ? ORDER BY catname ASC", name)
	if err != nil {
		return ""
	}
	defer rows.Close()
	var categories []Category
	var cat Category
	for rows.Next() {
		err := rows.Scan(&cat.Catname)
		if err != nil {
			fmt.Println(err)
			return ""
		}

		categories = append(categories, cat)
	}
	j, err := json.Marshal(categories)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		fmt.Println(string(j))
	}
	return string(j)
}
//***********************************
func getPendingCatTable(db *sql.DB, name string) string {

	rows, err := db.Query("SELECT catname, sender FROM pending WHERE receiver = ? ORDER BY catname ASC", name)
	if err != nil {
		return ""
	}
	defer rows.Close()
	var categories []Category
	var cat Category
	for rows.Next() {
		err := rows.Scan(&cat.Catname, &cat.Sender)
		if err != nil {
			fmt.Println(err)
			return ""
		}

		categories = append(categories, cat)
	}
	j, err := json.Marshal(categories)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		fmt.Println(string(j))
	}
	return string(j)
}
//***********************************
func getShareCatTable(db *sql.DB, name string) string {

	rows, err := db.Query("SELECT catname, sender FROM accepted WHERE receiver = ? ORDER BY catname ASC", name)
	if err != nil {
		return ""
	}
	defer rows.Close()
	var categories []Category
	var cat Category
	for rows.Next() {
		err := rows.Scan(&cat.Catname, &cat.Sender)
		if err != nil {
			fmt.Println(err)
			return ""
		}

		categories = append(categories, cat)
	}
	j, err := json.Marshal(categories)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		fmt.Println(string(j))
	}
	return string(j)
}
//***********************************
func getExportCatTable(db *sql.DB, name string) string {

	rows, err := db.Query("SELECT catname FROM shcats WHERE sender = ? ORDER BY catname ASC", name)
	if err != nil {
		return ""
	}
	defer rows.Close()
	var categories []Category
	var cat Category
	for rows.Next() {
		err := rows.Scan(&cat.Catname)
		if err != nil {
			fmt.Println(err)
			return ""
		}

		categories = append(categories, cat)
	}
	j, err := json.Marshal(categories)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		fmt.Println(string(j))
	}
	return string(j)
}
//***********************************
func getTaskTable(db *sql.DB, name string, catID string) string {

	rows, err := db.Query("SELECT catname, taskname, duedate, status, upredicted, spredicted, current FROM apoio WHERE USERNAME = ? and catname = ?", name, catID)
	if err != nil {
		return ""
	}
	defer rows.Close()
	var tasks []Task
	var tsk Task
	for rows.Next() {
		err := rows.Scan(&tsk.Catname, &tsk.Taskname, &tsk.Duedate, &tsk.Status, &tsk.Upredicted, &tsk.Spredicted, &tsk.Current)
		if err != nil {
			fmt.Println(err)
			return ""
		}

		tasks = append(tasks, tsk)
	}
	j, err := json.Marshal(tasks)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
	return string(j)
}
//***********************************
func getShareTaskTable(db *sql.DB, name string, catID string, sender string) string {

	rows, err := db.Query("SELECT catname, sender, taskname, duedate, status, upredicted, spredicted, current FROM sharetasks WHERE username = ? and catname = ? AND sender = ?", name, catID, sender)
	if err != nil {
		return ""
	}
	defer rows.Close()
	var tasks []Task
	var tsk Task
	for rows.Next() {
		err := rows.Scan(&tsk.Catname, &tsk.Sender, &tsk.Taskname, &tsk.Duedate, &tsk.Status, &tsk.Upredicted, &tsk.Spredicted, &tsk.Current)
		if err != nil {
			fmt.Println(err)
			return ""
		}

		tasks = append(tasks, tsk)
	}
	j, err := json.Marshal(tasks)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
	return string(j)
}
//****************************************************************************
func getExportTaskTable(db *sql.DB, name string, catID string) string {

	rows, err := db.Query("SELECT catname, taskname, current FROM sharetasks WHERE sender = ? and catname = ? AND status = ?", name, catID, "Complete")
	if err != nil {
		return ""
	}
	defer rows.Close()
	var tasks []Task
	var tsk Task
	for rows.Next() {
		err := rows.Scan(&tsk.Catname, &tsk.Taskname, &tsk.Current)
		if err != nil {
			fmt.Println(err)
			return ""
		}

		tasks = append(tasks, tsk)
	}
	j, err := json.Marshal(tasks)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
	return string(j)
}
//***********************************

func updateTime(db *sql.DB, name string, catID string, taskID string, time string) string {
	stmt, err := db.Prepare("UPDATE apoio SET current = ? WHERE username = ? AND catname = ? AND taskname = ?")

	stmt.Exec(time, name, catID, taskID)
	stmt.Close()
	stmt2, err := db.Prepare("UPDATE apoio SET status = ? WHERE username = ? AND catname = ? AND taskname = ?")
	stmt2.Exec("Paused", name, catID, taskID)
	stmt2.Close()
	if err != nil {
		print(err)
		return "2"
	}

	return "1"
}
//***********************************

func updateShareTime(db *sql.DB, name string, sender string, catID string, taskID string, time string) string {
	stmt, err := db.Prepare("UPDATE sharetasks SET current = ? WHERE username = ? AND sender = ? AND catname = ? AND taskname = ?")

	stmt.Exec(time, name, sender, catID, taskID)
	stmt.Close()
	stmt2, err := db.Prepare("UPDATE sharetasks SET status = ? WHERE username = ? AND sender = ? AND catname = ? AND taskname = ?")
	stmt2.Exec("Paused", name, sender, catID, taskID)
	stmt2.Close()
	if err != nil {
		print(err)
		return "2"
	}

	return "1"
}
//***********************************

func finishTime(db *sql.DB, name string, catID string, taskID string, time string) string {
	var pwd string
	stmt, err := db.Prepare("UPDATE apoio SET current = ? AND status = ? WHERE username = ? AND catname = ? AND taskname = ?")
	stmt.Exec(time, name, catID, taskID)
	stmt.Close()
	stmt2, err := db.Prepare("UPDATE apoio SET status = ? WHERE username = ? AND catname = ? AND taskname = ?")
	stmt2.Exec("Complete", name, catID, taskID)
	stmt2.Close()
	findate := getTime()
	stmt3, err := db.Prepare("UPDATE apoio SET findate = ? WHERE username = ? AND catname = ? AND taskname = ?")
	stmt3.Exec(findate, name, catID, taskID)
	stmt3.Close()
	rows, err := db.Query("SELECT upredicted FROM apoio WHERE username = ? AND catname = ? AND taskname = ?", name, catID, taskID)
	rows.Next()
	rows.Scan(&pwd)
	rows.Close()
	value1, err := strconv.ParseFloat(pwd, 64)
	value2, err := strconv.ParseFloat(time, 64)
	if value1 < value2 {
		value3 := 0.00
		value3 = value2 / value1
		count := 0
		stmt4, err := db.Prepare("UPDATE count SET correct = ? WHERE username = ? AND catname = ?")
		stmt4.Exec(count, name, catID)
		stmt4.Close()
		stmt3, err := db.Prepare("UPDATE apoio SET spredicted = ? WHERE username = ? AND catname = ? AND taskname = ?")
		stmt3.Exec(value3, name, catID, taskID)
		stmt3.Close()
		if err != nil {
			print(err)
			return "2"
		}
	}
	if value1 >= value2 {
		var cnt int
		rows2, err := db.Query("SELECT correct FROM count WHERE username = ? AND catname = ?", name, catID)
		rows2.Next()
		rows2.Scan(&cnt)
		rtn := cnt + 1
		value3 := 0
		stmt4, err := db.Prepare("UPDATE count SET correct = ? WHERE username = ? AND catname = ?")
		stmt4.Exec(rtn, name, catID)
		stmt4.Close()
		stmt3, err := db.Prepare("UPDATE apoio SET spredicted = ? WHERE username = ? AND catname = ? AND taskname = ?")
		stmt3.Exec(value3, name, catID, taskID)
		stmt3.Close()
		if err != nil {
			print(err)
			return "2"
		}
		if rtn >= 5 {
			return "5"
		}
		return "1"
	}
	if err != nil {
		print(err)
		return "2"
	}

	return "1"
}
//***********************************

func finishShareTime(db *sql.DB, name string, catID string, taskID string, time string, sender string) string {
	var pwd string
	stmt, err := db.Prepare("UPDATE sharetasks SET current = ? AND status = ? WHERE username = ? AND sender = ? AND catname = ? AND taskname = ?")
	stmt.Exec(time, name, sender, catID, taskID)
	stmt.Close()
	stmt2, err := db.Prepare("UPDATE sharetasks SET status = ? WHERE username = ? AND sender = ? AND catname = ? AND taskname = ?")
	stmt2.Exec("Complete", name, sender, catID, taskID)
	stmt2.Close()
	findate := getTime()
	stmt3, err := db.Prepare("UPDATE sharetasks SET findate = ? WHERE username = ? AND sender = ? AND catname = ? AND taskname = ?")
	stmt3.Exec(findate, name, sender, catID, taskID)
	stmt3.Close()
	rows, err := db.Query("SELECT upredicted FROM sharetasks WHERE username = ? AND sender = ? AND catname = ? AND taskname = ?", name, sender, catID, taskID)
	rows.Next()
	rows.Scan(&pwd)
	rows.Close()
	value1, err := strconv.ParseFloat(pwd, 64)
	value2, err := strconv.ParseFloat(time, 64)
	if value1 < value2 {
		value3 := 0.00
		value3 = value2 / value1
		count := 0
		stmt4, err := db.Prepare("UPDATE sharecount SET correct = ? WHERE username = ? AND catname = ?")
		stmt4.Exec(count, name, catID)
		stmt4.Close()
		stmt3, err := db.Prepare("UPDATE sharetasks SET spredicted = ? WHERE username = ? AND sender = ? AND catname = ? AND taskname = ?")
		stmt3.Exec(value3, name, sender, catID, taskID)
		stmt3.Close()
		if err != nil {
			print(err)
			return "2"
		}
	}
	if value1 >= value2 {
		var cnt int
		rows2, err := db.Query("SELECT correct FROM sharecount WHERE username = ? AND sender = ? AND catname = ?", name, sender, catID)
		rows2.Next()
		rows2.Scan(&cnt)
		rtn := cnt + 1
		value3 := 0
		stmt4, err := db.Prepare("UPDATE sharecount SET correct = ? WHERE username = ? AND sender = ? AND catname = ?")
		stmt4.Exec(rtn, name, sender, catID)
		stmt4.Close()
		stmt3, err := db.Prepare("UPDATE sharetasks SET spredicted = ? WHERE username = ? AND sender = ? AND catname = ? AND taskname = ?")
		stmt3.Exec(value3, name, sender, catID, taskID)
		stmt3.Close()
		if err != nil {
			print(err)
			return "2"
		}
		if rtn >= 5 {
			return "5"
		}
		return "1"
	}
	if err != nil {
		print(err)
		return "2"
	}

	return "1"
}
//***********************************

func addTask(db *sql.DB, name string, catID string, taskID string, dueDate string, status string, upredicted string, spredicted string, current string) string {
	stmt, err := db.Prepare("INSERT INTO apoio (username, catname, taskname, duedate, findate, status, upredicted, spredicted, current) VALUES (?,?,?,?,?,?,?,?,?);")
	stmt.Exec(name, catID, taskID, dueDate, "0", status, upredicted, spredicted, current)
	stmt.Close()
	stmt2, err := db.Prepare("INSERT INTO count (username, catname, correct) VALUES(?,?,?)")
	stmt2.Exec(name, catID, "0")
	stmt2.Close()
	if err != nil {
		print(err)
		return "2"
	}

	return "1"
}
//***********************************

func addShareTask(db *sql.DB, name string, sender string, catID string, taskID string, dueDate string, status string, upredicted string, spredicted string, current string) string {
	stmt, err := db.Prepare("INSERT INTO sharetasks (username, sender, catname, taskname, duedate, findate, status, upredicted, spredicted, current) VALUES (?,?,?,?,?,?,?,?,?);")
	stmt.Exec(name, sender, catID, taskID, dueDate, "0", status, upredicted, spredicted, current)
	stmt.Close()
	stmt2, err := db.Prepare("INSERT INTO sharecount (username, sender, catname, correct) VALUES(?,?,?,?)")
	stmt2.Exec(name, sender, catID, "0")
	stmt2.Close()
	if err != nil {
		print(err)
		return "2"
	}

	return "1"
}
//***********************************

func getHint(db *sql.DB, name string, catID string, upredicted string) string {
	var pwd string
	rows, err := db.Query("SELECT avg(spredicted) FROM apoio WHERE username = ? and catname = ? and status = ? ORDER BY findate DESC LIMIT 5", name, catID, "Complete")
	if err != nil {
		return ""
	}
	rows.Next()
	rows.Scan(&pwd)
	value1, err := strconv.ParseFloat(pwd, 64)
	value2, err := strconv.ParseFloat(upredicted, 64)
	value3 := value1 * value2
	if value3 <= 0 {
		value3 = value2
	}
	s := fmt.Sprintf("%f", value3)
	return s
}

//***********************************

func addUser(db *sql.DB, name string, pwd string) string {
	rtn := "0"
	stmt, err := db.Prepare("INSERT INTO USER1 (username, password) VALUES (?,?);")
	if err != nil {
		return rtn
	}
	stmt.Exec(name, pwd)
	stmt.Close()
	return "1"
}
//***********************************

func addCat(db *sql.DB, name string, catID string) string {
	rtn := "0"
	stmt, err := db.Prepare("INSERT INTO user2 (username, catname) VALUES (?,?);")
	if err != nil {
		return rtn
	}
	stmt.Exec(name, catID)
	stmt.Close()
	return "1"
}
//***********************************

func addExportCat(db *sql.DB, name string, catID string, receiver string) string {
	rtn := "0"
	stmt, err := db.Prepare("INSERT INTO pending (sender, catname, receiver) VALUES (?,?,?);")
	if err != nil {
		return rtn
	}
	stmt.Exec(name, catID, receiver)
	stmt.Close()
	stmt2, err := db.Prepare("INSERT INTO shcats (sender, catname) VALUES (?,?);")
	if err != nil {
		return rtn
	}
	stmt2.Exec(name, catID)
	stmt2.Close()
	return "1"
}
//***********************************

func hashAndSalt(paswd []byte) string {

	hash, err := bcrypt.GenerateFromPassword(paswd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)

}

//********************************************************************

func getTime() int64 {
	now := time.Now() // current local time
	sec := now.Unix()
	return sec
}

//********************************************************************

func getDate() time.Time {
	now := time.Now() // current local time
	return now
}

//***********************************

func checkTime(db *sql.DB, username string, timestamp int64) int64 {
	var pwd int64
	rows, err := db.Query("SELECT timestamp FROM user3 WHERE USERNAME = ?", username)
	if err != nil {
		return pwd
	}
	rows.Next()
	rows.Scan(&pwd)
	rows.Close()

	timesum := timestamp - pwd

	if timesum < 86400 {
		stmt, err := db.Prepare("UPDATE user3 SET timestamp = ? WHERE username = ?")
		if err != nil {
			return 2
		}
		stmt.Exec(timestamp, username)
		stmt.Close()
		return 1
	}
	if timesum >= 86400 {
		stmt, err := db.Prepare("DELETE FROM user3 WHERE username = ?")
		if err != nil {
			return 2
		}
		stmt.Exec(username)
		stmt.Close()
		return 2
	}
	return 2
}

//***********************************
//: Reads in stuff from a connection - at most 1024 bytes

func netReading(c net.Conn) string {
	buf := make([]byte, 1024)
	nr, err := c.Read(buf)
	if err != nil {
		log.Println("Read failed")
		log.Fatal(err)
	}
	return string(buf[0:nr])
}

//***********************************

func netWrite(c net.Conn, msg string) {
	_, err := c.Write([]byte(msg))
	if err != nil {
		log.Println("Write failed")
		log.Fatal(err)
	}
}

//**************************************
func wordTrim(input string) string {
	keyRead := strings.TrimSuffix(input, "\x1F")
	return keyRead
}
