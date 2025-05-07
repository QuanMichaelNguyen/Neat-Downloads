import React from 'react';

const Toast = ({ message }) => (
  message ? <div className="toast-popup">{message}</div> : null
);

export default Toast; 