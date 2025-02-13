import React from 'react';
import { Header } from '@cred/neopop-web/lib/components';
import './Header.css';

const HeaderComp = () => {
    return (
        <div className="header-container">
            <div className="header-left">
                <div className="header-heading">CHECKMATE</div>
                <div className="header-description">List Your Tasks With CHECKMATE</div>
            </div>
            <Header
                onBackClick={() => {
                    console.log('back clicked');
                }}
            />
        </div>
    );
};

export default HeaderComp;