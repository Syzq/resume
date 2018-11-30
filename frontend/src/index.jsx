import React, { Component } from 'react'
import ReactDOM from 'react-dom'
import { Button, Row, Col, Input, List } from 'antd'

import axios from 'axios'

import 'antd/dist/antd.css'

export default class Index extends Component {
  constructor(props) {
    super(props)
    this.state = {
      inputValue: '',
      // todos: ['first', 'second', 'third'],
      todos: [],
      btnStatus: 'add', //默认显示添加按钮
      itemIndex: null //当前 list 的下标
    }
  }
  getAllTodos() {
    axios.get('todos').then(res => {
      if (res.data) {
        this.setState({
          todos: res.data
        })
      }
    })
  }
  componentDidMount() {
    this.getAllTodos()
    console.log(
      `
      ------------------
      < hello ! >
      ------------------
              \\   ^__^
               \\  (oo)\\_______
                  (__)\\       )\\/\\
                      ||----w |
                      ||     ||
      `
    )
  }
  render() {
    let { inputValue, btnStatus, todos } = this.state
    return (
      <div style={{ marginTop: '50px' }}>
        <Row type="flex" justify="center">
          <Col
            style={{
              width: '500px',
              margin: '0 8px 8px 0'
            }}
          >
            <Input
              placeholder="添加 TODO（按回车添加）"
              value={inputValue}
              onChange={this.inputChange}
              onPressEnter={
                btnStatus === 'add' ? this.addTodo : this.editConfirm
              }
              ref={input => {
                this.textInput = input
              }}
            />
            <List
              style={{ marginTop: '8px' }}
              bordered
              dataSource={todos}
              renderItem={(item, index) => (
                <List.Item
                  actions={[
                    <span onClick={this.editTodo.bind(this, index)}>编辑</span>,
                    <span onClick={this.delTodo.bind(this, index)}>删除</span>
                  ]}
                >
                  {item.content}
                </List.Item>
              )}
            />
          </Col>
          <Col>{this.btnStatus(btnStatus)}</Col>
        </Row>
      </div>
    )
  }
  //切换不同的按钮
  btnStatus(sts) {
    return sts === 'add' ? (
      <Button type="primary" onClick={this.addTodo}>
        添加
      </Button>
    ) : (
      <Button
        type="primary"
        onClick={this.editConfirm}
        style={{ background: 'red', borderColor: 'red' }}
      >
        修改
      </Button>
    )
  }
  // 监听输入框输入内容
  inputChange = e => {
    // console.log(e.target.value)
    this.setState({
      inputValue: e.target.value
    })
  }
  // 添加 TODO
  addTodo = e => {
    let { inputValue } = this.state
    if (inputValue === '') {
      return alert('还没有输入任何内容')
    }
    axios
      .post('add', {
        content: inputValue
      })
      .then(res => {
        this.setState({
          // todos: todos.concat(inputValue),
          inputValue: ''
        })
        this.getAllTodos()
      })
  }
  // item 列表中编辑 TODO的按钮
  editTodo = index => {
    let { todos } = this.state
    this.textInput.focus()
    this.setState({
      inputValue: todos[index].content,
      btnStatus: 'edit',
      itemIndex: index
    })
  }
  //确定修改的按钮
  editConfirm = index => {
    let { inputValue, todos, itemIndex } = this.state
    // console.log(inputValue, todos, todos[itemIndex].id)

    axios
      .put(`todos/${todos[itemIndex].id}`, {
        content: inputValue
      })
      .then(res => {
        // console.log(res)
        todos[itemIndex].content = inputValue
        this.setState({
          inputValue: '',
          btnStatus: 'add'
        })
      })
  }
  // 删除 TODO
  delTodo = index => {
    let { todos } = this.state
    // console.log(todos[index].id)
    // todos.splice(index, 1)
    // this.setState({
    //   todos
    // })
    axios.delete(`deltodo/${todos[index].id}`).then(res => {
      this.getAllTodos()
    })
  }
}

ReactDOM.render(<Index />, document.getElementById('root'))
