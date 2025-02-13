import React, { Component } from "react";
import axios from "axios";
import { Card, Header, Form, Input, Button, Modal, Checkbox, Icon } from "semantic-ui-react"; // Import Icon here
import { LocalizationProvider } from "@mui/lab";
import AdapterDateFns from "@mui/lab/AdapterDateFns";
import DateTimePicker from "@mui/lab/DateTimePicker";
import TextField from "@mui/material/TextField";
import moment from "moment";
import Box from "./Box";
import "./Box.css";

let endpoint = "http://localhost:9000";

class ToDoList extends Component {
  constructor(props) {
    super(props);

    this.state = {
      task: "",
      description: "",
      important: false,
      dateTime: new Date(),
      items: [],
      editTaskId: null,
      modalOpen: false,
    };
  }

  componentDidMount() {
    this.getTasks();
  }

  componentDidUpdate(prevProps) {
    if (prevProps.filter !== this.props.filter) {
      this.getTasks();
    }
  }

  onChange = (e) => {
    this.setState({ [e.target.name]: e.target.value });
  };

  onCheckboxChange = (e, { checked }) => {
    this.setState({ important: checked });
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
              description: this.state.description,
              important: this.state.important,
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
            this.setState({ task: "", description: "", important: false, editTaskId: null, dateTime: new Date(), modalOpen: false });
          });
      } else {
        axios
          .post(
            `${endpoint}/api/todo`,
            {
              task: this.state.task,
              description: this.state.description,
              important: this.state.important,
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
            this.setState({ task: "", description: "", important: false, dateTime: new Date(), modalOpen: false });
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
          let items = res.data.map((item) => ({
            ...item,
            dateTime: moment.utc(`${item.date} ${item.time}`).toDate(),
          }));

          // Apply filter
          const today = moment().format("YYYY-MM-DD");
          if (this.props.filter === 'today') {
            items = items.filter(item => item.date === today);
          } else if (this.props.filter === 'important') {
            items = items.filter(item => item.important);
          } else if (this.props.filter === 'completed') {
            items = items.filter(item => item.done);
          } else if (this.props.filter === 'incomplete') {
            items = items.filter(item => !item.done);
          }

          this.setState({ items });
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

  editTask = (id, task, description, date, time, important) => {
    this.setState({ task, description, important, editTaskId: id, dateTime: moment.utc(`${date} ${time}`).toDate(), modalOpen: true });
  };

  openModal = () => {
    this.setState({ modalOpen: true });
  };

  closeModal = () => {
    this.setState({ modalOpen: false });
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
          <Card.Group>
            <Card className="box-card add-task-card" onClick={this.openModal}>
              <Card.Content>
                <Card.Header textAlign="center">
                  <Icon name="plus" size="huge" />
                </Card.Header>
              </Card.Content>
            </Card>
            {this.state.items.map((item) => (
              <Box
                key={item.id}
                item={item}
                editTask={this.editTask}
                updateTask={this.updateTask}
                undoTask={this.undoTask}
                deleteTask={this.deleteTask}
              />
            ))}
          </Card.Group>
        </div>
        <Modal open={this.state.modalOpen} onClose={this.closeModal}>
          <Modal.Header>{this.state.editTaskId ? "Update Task" : "Add Task"}</Modal.Header>
          <Modal.Content>
            <Form onSubmit={this.onSubmit}>
              <Form.Field>
                <label>Task Name</label>
                <Input
                  type="text"
                  name="task"
                  onChange={this.onChange}
                  value={this.state.task}
                  fluid
                  placeholder="Task Name"
                />
              </Form.Field>
              <Form.Field>
                <label>Description</label>
                <Input
                  type="text"
                  name="description"
                  onChange={this.onChange}
                  value={this.state.description}
                  fluid
                  placeholder="Description"
                />
              </Form.Field>
              <Form.Field>
                <label>Mark as Important</label>
                <Checkbox
                  toggle
                  checked={this.state.important}
                  onChange={this.onCheckboxChange}
                />
              </Form.Field>
              <Form.Field>
                <LocalizationProvider dateAdapter={AdapterDateFns}>
                  <DateTimePicker
                    label="Date & Time"
                    value={this.state.dateTime}
                    onChange={this.onDateTimeChange}
                    renderInput={(params) => <TextField {...params} />}
                  />
                </LocalizationProvider>
              </Form.Field>
              <Button type="submit">
                {this.state.editTaskId ? "Update Task" : "Add Task"}
              </Button>
            </Form>
          </Modal.Content>
        </Modal>
      </div>
    );
  }
}

export default ToDoList;