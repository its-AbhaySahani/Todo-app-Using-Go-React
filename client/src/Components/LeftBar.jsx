import React from 'react';
import { Link } from 'react-router-dom';
import './LeftBar.css';

const LeftBar = ({ setFilter }) => {
  return (
    <div className="leftbar-container">
      <div className="leftbar-heading">CHECKMATE</div>
      <div className="leftbar-menu">
        <button onClick={() => setFilter('today')}>Today's Tasks</button>
        <button onClick={() => setFilter('all')}>All Tasks</button>
        <button onClick={() => setFilter('important')}>Important Tasks</button>
        <button onClick={() => setFilter('completed')}>Completed Tasks</button>
        <button onClick={() => setFilter('incomplete')}>Incomplete Tasks</button>
      </div>
    </div>
  );
};

export default LeftBar;