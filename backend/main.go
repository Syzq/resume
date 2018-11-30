package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

//Todo struct
type Todo struct {
	Id        int       `json:"id"`
	Content   string    `json:"content" `
	State     bool      `json:"state"`
	CreatedAt time.Time `json:"created_at"`
}
type Todos []Todo

type Server struct {
	db *sql.DB
}

func main() {

	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/user?charset=utf8&parseTime=true")

	checkErr(err)
	server := &Server{db: db}
	defer db.Close()

	http.HandleFunc("/add", server.addTodo)
	http.HandleFunc("/todos", server.allTodos)
	http.HandleFunc("/todos/", server.updateTodo)
	http.HandleFunc("/deltodo/", server.deleteTodo)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// add todo
func (s *Server) addTodo(w http.ResponseWriter, r *http.Request) {
	todo := &Todo{}
	log.Println("\n----获取请求的方法----", r.Method)

	if r.Method == "POST" {
		// Decode从输入流读取下一个json编码值并保存在v指向的值里
		err := json.NewDecoder(r.Body).Decode(&todo)
		// req, err := ioutil.ReadAll(r.Body)
		// json.Unmarshal([]byte(req), &todo)

		// 插入一条数据Now().Format("2006-01-02 15:04:05")
		result, err := s.db.Exec("insert into demo_todo(content,state,createdat) values(?,?,?)", todo.Content, todo.State, time.Now().Format("2006-01-02 15:04:05"))

		fmt.Printf("%v\n", todo)
		checkErr(err)
		Id64, err := result.LastInsertId()
		Id := int(Id64)
		todo = &Todo{Id: Id}
		// fmt.Println("OK", todo)

	} else {
		// json.NewEncoder(w).Encode("未知请求")
	}
}

// all todos
func (s *Server) allTodos(res http.ResponseWriter, req *http.Request) {
	var todos []*Todo

	if req.Method == "GET" {
		rows, err := s.db.Query("select id, content, state, createdat from demo_todo")
		checkErr(err)

		for rows.Next() {
			todo := &Todo{}
			// fmt.Println("%v\n", &todo.CreatedAt)
			rows.Scan(&todo.Id, &todo.Content, &todo.State, &todo.CreatedAt)
			todos = append(todos, todo)
		}

		json.NewEncoder(res).Encode(todos)
	}
}

// update todo
func (s *Server) updateTodo(res http.ResponseWriter, req *http.Request) {
	todoIndex := &Todo{}
	if req.Method == "PUT" {
		log.Println("\n----获取请求的方法----", req.Method)
		// log.Println("\n----获取URL----", req.URL.Path)

		// 正则匹配路由参数
		r, _ := regexp.Compile("\\d+")
		todoID := r.FindString(req.URL.Path)
		fmt.Println("要修改的 ID", todoID)
		err := json.NewDecoder(req.Body).Decode(todoIndex)
		checkErr(err)

		fmt.Printf("%v\n", todoIndex.Content)
		fmt.Println(todoIndex.State)
		// res, err := s.db.Exec("update demo_todo set id=?, content=?, state=?", 1, "tes2t", 0)
		// checkErr(err)
		// affect, err := res.RowsAffected()
		// checkErr(err)

		// fmt.Println("修改成功，ID 为：", affect)

		// 更新数据
		stmt, err := s.db.Prepare("UPDATE demo_todo SET content=? ,state=? where id=?")
		checkErr(err)

		res, err := stmt.Exec(todoIndex.Content, todoIndex.State, todoID)
		checkErr(err)
		 // RowsAffected返回被update、insert或delete命令影响的行数。
		affect, err := res.RowsAffected()
		checkErr(err)

		fmt.Println("修改成功：", affect)
	}
}

func (s *Server) deleteTodo(res http.ResponseWriter, req *http.Request) {
	if req.Method == "DELETE" {
		log.Println("\n----获取请求的方法----", req.Method)
		// log.Println("\n----获取请求的URL----", req.URL.Path)
		r := regexp.MustCompile("\\d+")

		todoID := r.FindString(req.URL.Path)
		fmt.Println("要删除的 ID", todoID)
		stmt, err := s.db.Prepare("delete from demo_todo where id=?")
		checkErr(err)
		_, err = stmt.Exec(todoID)
		checkErr(err)

		// affect, err := res.RowsAffected()
		// checkErr(err)
		// fmt.Println("删除成功，ID 为：", affect)
	}
}
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)

	}
}
