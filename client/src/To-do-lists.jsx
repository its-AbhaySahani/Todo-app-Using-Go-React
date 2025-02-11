import React, { Component } from "react";
import axios from "axios";
import { Card, Header, Form, Input, Icon, Button } from "semantic-ui-react";
import { LocalizationProvider } from "@mui/lab";
import AdapterDateFns from "@mui/lab/AdapterDateFns";
import DateTimePicker from "@mui/lab/DateTimePicker";
import TextField from "@mui/material/TextField";
import moment from "moment";

let endpoint = "http://localhost:9000";

class ToDoList extends Component {
  constructor(props) {
    super(props);

    this.state = {
      task: "",
      dateTime: new Date(),
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

  onDateTimeChange = (dateTime) => {
    this.setState({ dateTime });
  };

  onSubmit = (e) => {
    e.preventDefault();
    const token = localStorage.getItem("token");
    const formattedDate = moment.utc(this.state.dateTime).format("YYYY-MM-DD");
    const formattedTime = moment.utc(this.state.dateTime).format("HH:mm:ss");
    if (this.state.task) {
      if (this.state.editTaskId) {
        axios
          .put(
            `${endpoint}/api/todo/${this.state.editTaskId}`,
            {
              task: this.state.task,
              done: false,
              date: formattedDate,
              time: formattedTime,
            },
            {
              headers: {
                Authorization: `Bearer ${token}`,
              },
            }
          )
          .then(() => {
            this.getTasks();
            this.setState({ task: "", editTaskId: null, dateTime: new Date() });
          });
      } else {
        axios
          .post(
            `${endpoint}/api/todo`,
            {
              task: this.state.task,
              done: false,
              date: formattedDate,
              time: formattedTime,
            },
            {
              headers: {
                Authorization: `Bearer ${token}`,
              },
            }
          )
          .then(() => {
            this.getTasks();
            this.setState({ task: "", dateTime: new Date() });
          });
      }
    }
  };

  getTasks = () => {
    const token = localStorage.getItem("token");
    axios
      .get(`${endpoint}/api/todos`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })
      .then((res) => {
        if (res.data) {
          this.setState({
            items: res.data.map((item) => ({
              ...item,
              dateTime: moment.utc(`${item.date} ${item.time}`).toDate(),
            })),
          });
        }
      });
  };

  updateTask = (id) => {
    const token = localStorage.getItem("token");
    axios
      .put(
        `${endpoint}/api/todo/${id}`,
        {
          done: true,
        },
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      )
      .then(() => {
        this.getTasks();
      });
  };

  deleteTask = (id) => {
    const token = localStorage.getItem("token");
    axios
      .delete(`${endpoint}/api/todo/${id}`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })
      .then(() => {
        this.getTasks();
      });
  };

  undoTask = (id) => {
    const token = localStorage.getItem("token");
    axios
      .put(
        `${endpoint}/api/todo/undo/${id}`,
        {
          done: false,
        },
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      )
      .then(() => {
        this.getTasks();
      });
  };

  editTask = (id, task, date, time) => {
    this.setState({ task, editTaskId: id, dateTime: moment.utc(`${date} ${time}`).toDate() });
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
            <LocalizationProvider dateAdapter={AdapterDateFns}>
              <DateTimePicker
                label="Date & Time"
                value={this.state.dateTime}
                onChange={this.onDateTimeChange}
                renderInput={(params) => <TextField {...params} />}
              />
            </LocalizationProvider>
            <Button type="submit">
              {this.state.editTaskId ? "Update Task" : "Add Task"}
            </Button>
          </Form>
        </div>
        <div className="row">
          <Card.Group>
            {this.state.items.map((item) => (
              <Card key={item.id} color={item.done ? "green" : "red"}>
                <Card.Content>
                  <Card.Header textAlign="left">
                    <div style={{ wordWrap: "break-word" }}>{item.task}</div>
                  </Card.Header>
                  <Card.Meta textAlign="left">
                    <span>{moment.utc(`${item.date} ${item.time}`).format("YYYY-MM-DD HH:mm:ss")}</span>
                  </Card.Meta>
                  <Card.Meta textAlign="right">
                    <Icon
                      name="edit"
                      color="blue"
                      onClick={() => this.editTask(item.id, item.task, item.date, item.time)}
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