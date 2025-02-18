import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import './RightBar.css';

const RightBar = ({ isAuthenticated, setIsAuthenticated }) => {
  const navigate = useNavigate();

  const handleLogout = () => {
    localStorage.removeItem("token");
    setIsAuthenticated(false);
    navigate("/login");
  };

  return (
    <div className="rightbar-container">
      <div className="rightbar-menu">
        <button onClick={() => navigate('/')}>Home</button>
        {isAuthenticated && <button onClick={() => navigate('/profile')}>Profile</button>}
        {isAuthenticated ? (
          <button onClick={handleLogout}>Logout</button>
        ) : (
          <>
            <button onClick={() => navigate('/login')}>Login</button>
            <button onClick={() => navigate('/register')}>Register</button>
          </>
        )}
        <Link to="/get-together">
          <button>Go to Shared Tasks</button>
        </Link>
      </div>
      <div className="rightbar-bottom-buttons">
        <button onClick={() => navigate('/create-team')}>Create Team</button>
        <button onClick={() => navigate('/my-teams')}>My Teams</button>
      </div>
    </div>
  );
};

export default RightBar;