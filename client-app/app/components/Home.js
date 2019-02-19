// @flow
import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import { ipcRenderer } from 'electron';
import routes from '../constants/routes';
import styles from './Home.css';


type Props = {};
type State = {
  fileActive: boolean
};

export default class Home extends Component<Props, State> {
  props: Props;

  state = {
    fileActive: false
  };

  componentDidMount = ()=>{
    ipcRenderer.on("file-uploaded",(ev,arg)=>{
      alert("File "+arg+" is uploaded to go process")
    })
  }

  setActive = () => {
    this.setState({ fileActive: true });
  };

  setInactive = () => {
    this.setState({ fileActive: false });
  };

  onDragOver = e => {
    e.stopPropagation();
    e.preventDefault();
  };

  upload = e => {
    e.preventDefault();
    for (const f of e.dataTransfer.files) {
      console.log('File(s) you dragged here: ', f);
    }
    ipcRenderer.send('file',e.dataTransfer.files[0].path);
    this.setState({ fileActive: false });
  };

  render() {
    const { fileActive } = this.state;
    return (
      <div>
        <div className={styles.base} data-tid="base">
          <img
            src="https://vinci.id"
            alt="bg"
          />
        </div>
        <div
          className={styles.blur_container}
          data-tid="blur_container"
          onDragEnter={this.setActive}
          onDragLeave={this.setInactive}
          onDragExit={this.setInactive}
          onDrop={this.upload}
          onDragOver={this.onDragOver}
        >
          <h3>{fileActive ? 'Drop to upload!' : 'Drag file/folder here!'}</h3>
        </div>
      </div>
    );
  }
}
