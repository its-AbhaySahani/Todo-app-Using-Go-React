import React, { Component } from "react";
import axios from "axios";
import { Card, Header, Form, Input, Icon, Button } from "semantic-ui-react";

let endpoint = "http://localhost:9000";

class ToDoList extends Component {
  constructor(props) {
    super(props);

    this.state = {
      task: "",
      items: [],
      editTaskId: null,
    };
  }

  componentDidMount() {
    this.getTasks();
  }

  onChange = (e) => {
    this.setState({ [e.target.name]: e.target.value });
  };

  onSubmit = (e) => {
    e.preventDefault();
    if (this.state.task) {
      if (this.state.editTaskId) {
        axios
          .put(`${endpoint}/api/todo/${this.state.editTaskId}`, {
            task: this.state.task,
            done: false,
          })
          .then((res) => {
            this.getTasks();
            this.setState({ task: "", editTaskId: null });
          });
      } else {
        axios
          .post(`${endpoint}/api/todo`, {
            task: this.state.task,
            done: false,
          })
          .then((res) => {
            this.getTasks();
            this.setState({ task: "" });
          });
      }
    }
  };

  getTasks = () => {
    axios.get(`${endpoint}/api/todos`).then((res) => {
      if (res.data) {
        this.setState({ items: res.data });
      }
    });
  };

  updateTask = (id) => {
    axios
      .put(`${endpoint}/api/todo/${id}`, {
        done: true,
      })
      .then((res) => {
        this.getTasks();
      });
  };

  deleteTask = (id) => {
    axios.delete(`${endpoint}/api/todo/${id}`).then((res) => {
      this.getTasks();
    });
  };

  undoTask = (id) => {
    axios
      .put(`${endpoint}/api/todo/undo/${id}`, {
        done: false,
      })
      .then((res) => {
        this.getTasks();
      });
  };

  editTask = (id, task) => {
    this.setState({ task, editTaskId: id });
  };

  render() {
    return (
      <div>
        <div className="row">
          <Header className="header" as="h2" color="yellow">
            To Do List
          </Header>
        </div>
        <div className="row">
          <Form onSubmit={this.onSubmit}>
            <Input
              type="text"
              name="task"
              onChange={this.onChange}
              value={this.state.task}
              fluid
              placeholder="Create Task"
            />
            <Button type="submit">
              {this.state.editTaskId ? "Update Task" : "Add Task"}
            </Button>
          </Form>
        </div>
        <div className="row">
          <Card.Group>
            {this.state.items.map((item) => (
              <Card key={item.id}>
                <Card.Content>
                  <Card.Header textAlign="left">
                    <div style={{ wordWrap: "break-word" }}>{item.task}</div>
                  </Card.Header>

                  <Card.Meta textAlign="right">
                    <Icon
                      name="edit"
                      color="blue"
                      onClick={() => this.editTask(item.id, item.task)}
                    />
                    <span style={{ paddingRight: 10 }}>Edit</span>
                    <Icon
                      name="check circle"
                      color="green"
                      onClick={() => this.updateTask(item.id)}
                    />
                    <span style={{ paddingRight: 10 }}>Done</span>
                    <Icon
                      name="undo"
                      color="yellow"
                      onClick={() => this.undoTask(item.id)}
                    />
                    <span style={{ paddingRight: 10 }}>Undo</span>
                    <Icon
                      name="delete"
                      color="red"
                      onClick={() => this.deleteTask(item.id)}
                    />
                    <span>Delete</span>
                  </Card.Meta>
                </Card.Content>
              </Card>
            ))}
          </Card.Group>
        </div>
      </div>
    );
  }
}

export default ToDoList;