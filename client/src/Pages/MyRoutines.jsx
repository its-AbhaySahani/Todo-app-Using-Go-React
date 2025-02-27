// client/src/Pages/MyRoutines.jsx
import React, { useEffect, useState } from "react";
import axios from "axios";
import { Card, Header, Button, Dropdown } from "semantic-ui-react";
import moment from "moment";
import Box from "../Box";
import "./MyRoutines.css";

const MyRoutines = () => {
  const [morningTasks, setMorningTasks] = useState([]);
  const [noonTasks, setNoonTasks] = useState([]);
  const [eveningTasks, setEveningTasks] = useState([]);
  const [nightTasks, setNightTasks] = useState([]);
  const [activeSection, setActiveSection] = useState("morning");
  const [selectedDay, setSelectedDay] = useState(moment().format('dddd').toLowerCase());
  const [loading, setLoading] = useState(false);

  const dayOptions = [
    { key: 'sunday', text: 'Sunday', value: 'sunday' },
    { key: 'monday', text: 'Monday', value: 'monday' },
    { key: 'tuesday', text: 'Tuesday', value: 'tuesday' },
    { key: 'wednesday', text: 'Wednesday', value: 'wednesday' },
    { key: 'thursday', text: 'Thursday', value: 'thursday' },
    { key: 'friday', text: 'Friday', value: 'friday' },
    { key: 'saturday', text: 'Saturday', value: 'saturday' },
  ];

  const fetchRoutineTasks = async (day, scheduleType) => {
    const token = localStorage.getItem("token");
    try {
      setLoading(true);
      console.log(`Fetching ${scheduleType} tasks for ${day}...`);
      
      // Set empty arrays initially before fetching new data
      if (scheduleType === "morning") setMorningTasks([]);
      if (scheduleType === "noon") setNoonTasks([]);
      if (scheduleType === "evening") setEveningTasks([]);
      if (scheduleType === "night") setNightTasks([]);
      
      const response = await axios.get(
        `http://localhost:9000/api/routine/day/${day}/${scheduleType}`,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      console.log(`Got ${scheduleType} tasks for ${day}:`, response.data);
      
      if (response.data) {
        const tasks = response.data.map(task => ({
          ...task,
          dateTime: moment.utc(`${task.date} ${task.time}`).toDate(),
        }));
        
        switch (scheduleType) {
          case "morning":
            setMorningTasks(tasks);
            break;
          case "noon":
            setNoonTasks(tasks);
            break;
          case "evening":
            setEveningTasks(tasks);
            break;
          case "night":
            setNightTasks(tasks);
            break;
          default:
            break;
        }
      }
      setLoading(false);
    } catch (error) {
      console.error(`Error fetching ${scheduleType} tasks for ${day}:`, error);
      console.error("Request URL:", `http://localhost:9000/api/routine/day/${day}/${scheduleType}`);
      
      setLoading(false);
      // Make sure to set empty arrays when there's an error
      if (scheduleType === "morning") setMorningTasks([]);
      if (scheduleType === "noon") setNoonTasks([]);
      if (scheduleType === "evening") setEveningTasks([]);
      if (scheduleType === "night") setNightTasks([]);
    }
  };

  const fetchAllTasksForDay = (day) => {
    // Clear all tasks first when day changes
    setMorningTasks([]);
    setNoonTasks([]);
    setEveningTasks([]);
    setNightTasks([]);
    
    // Then fetch new data
    fetchRoutineTasks(day, "morning");
    fetchRoutineTasks(day, "noon");
    fetchRoutineTasks(day, "evening");
    fetchRoutineTasks(day, "night");
  };

  useEffect(() => {
    fetchAllTasksForDay(selectedDay);
  }, [selectedDay]); // Re-fetch when selected day changes

  const handleDayChange = (e, { value }) => {
    setSelectedDay(value);
    // Reset all task arrays when changing days
    setMorningTasks([]);
    setNoonTasks([]);
    setEveningTasks([]);
    setNightTasks([]);
  };

  const renderTasks = (tasks) => (
    <div className="tasks-container">
      {loading ? (
        <p className="loading-message">Loading tasks...</p>
      ) : tasks.length > 0 ? (
        <Card.Group>
          {tasks.map((item) => (
            <Box key={item.id} item={item} />
          ))}
        </Card.Group>
      ) : (
        <p className="no-tasks-message">No tasks scheduled for this time.</p>
      )}
    </div>
  );

  // Get the current active tasks based on selected section
  const getActiveTasks = () => {
    switch (activeSection) {
      case "morning":
        return renderTasks(morningTasks);
      case "noon":
        return renderTasks(noonTasks);
      case "evening":
        return renderTasks(eveningTasks);
      case "night":
        return renderTasks(nightTasks);
      default:
        return renderTasks([]);
    }
  };

  return (
    <div className="my-routines">
      <div className="routines-header">
        <Header as="h2" color="yellow">
          My Daily Routines
        </Header>
        <div className="day-selector">
          <span className="day-label">Select Day: </span>
          <Dropdown
            selection
            options={dayOptions}
            value={selectedDay}
            onChange={handleDayChange}
          />
        </div>
      </div>

      <div className="time-buttons">
        <Button
          color={activeSection === "morning" ? "blue" : "grey"}
          onClick={() => setActiveSection("morning")}
        >
          Morning
        </Button>
        <Button
          color={activeSection === "noon" ? "blue" : "grey"}
          onClick={() => setActiveSection("noon")}
        >
          Noon
        </Button>
        <Button
          color={activeSection === "evening" ? "blue" : "grey"}
          onClick={() => setActiveSection("evening")}
        >
          Evening
        </Button>
        <Button
          color={activeSection === "night" ? "blue" : "grey"}
          onClick={() => setActiveSection("night")}
        >
          Night
        </Button>
      </div>

      <div className="routine-section">
        {getActiveTasks()}
      </div>
    </div>
  );
};

export default MyRoutines;