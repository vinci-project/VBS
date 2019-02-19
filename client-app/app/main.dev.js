/* eslint global-require: off */

/**
 * This module executes inside of electron's main process. You can start
 * electron renderer process from here and communicate with the other processes
 * through IPC.
 *
 * When running `yarn build` or `yarn build-main`, this file is compiled to
 * `./app/main.prod.js` using webpack. This gives us some performance wins.
 *
 * @flow
 */
import { app, BrowserWindow, Tray, Menu,ipcMain } from 'electron';
import path from 'path';
import { autoUpdater } from 'electron-updater';
import log from 'electron-log';
import child_process from 'child_process';
import axios from 'axios';
import MenuBar from 'menubar';
import MenuBuilder from './menu';

export default class AppUpdater {
  constructor() {
    log.transports.file.level = 'info';
    autoUpdater.logger = log;
    autoUpdater.checkForUpdatesAndNotify();
  }
}

const iconPath = path.join(__dirname, 'logo.png');
let mainWindow = null;
let appIcon = null;
let GoDaemon = null;

if (process.env.NODE_ENV === 'production') {
  const sourceMapSupport = require('source-map-support');
  sourceMapSupport.install();
}

if (
  process.env.NODE_ENV === 'development' ||
  process.env.DEBUG_PROD === 'true'
) {
  require('electron-debug')();
}

const installExtensions = async () => {
  const installer = require('electron-devtools-installer');
  const forceDownload = !!process.env.UPGRADE_EXTENSIONS;
  const extensions = ['REACT_DEVELOPER_TOOLS', 'REDUX_DEVTOOLS'];

  return Promise.all(
    extensions.map(name => installer.default(installer[name], forceDownload))
  ).catch(console.log);
};

/**
 * Add event listeners...
 */

app.on('window-all-closed', () => {
  // Respect the OSX convention of having the application in memory even
  // after all windows have been closed
  if (process.platform !== 'darwin') {
    app.quit();
  }
});

// const menuURL = `file://${__dirname}/menu.html`;
// const mb = MenuBar({
//   index: menuURL,
//   icon: path.join(__dirname, '..', 'resources/logo.png')
// });

// mb.on('ready', () => {
//   console.log('Menu is ready');
//   // your app code here
// });


app.on('ready', async () => {
  if (
    process.env.NODE_ENV === 'development' ||
    process.env.DEBUG_PROD === 'true'
  ) {
    // await installExtensions();
  }

  mainWindow = new BrowserWindow({
    show: false,
    icon: iconPath,
    width: 300,
    height: 600
  });

  mainWindow.loadURL(`file://${__dirname}/app.html`);

  // @TODO: Use 'ready-to-show' event
  //        https://github.com/electron/electron/blob/master/docs/api/browser-window.md#using-ready-to-show-event
  mainWindow.webContents.on('did-finish-load', () => {
    if (!mainWindow) {
      throw new Error('"mainWindow" is not defined');
    }
    if (process.env.START_MINIMIZED) {
      mainWindow.minimize();
    } else {
      mainWindow.show();
      mainWindow.focus();
    }
  });

  mainWindow.on('closed', () => {
    mainWindow = null;
  });

  // const menuBuilder = new MenuBuilder(mainWindow);
  // menuBuilder.buildMenu();

  appIcon = new Tray(iconPath);
  const contextMenu = Menu.buildFromTemplate([
    {
      label: 'Upload',
      click: () => {}
    },
    { label: 'Quit', accelerator: 'Command+Q', selector: 'terminate:' }
  ]);
  appIcon.setToolTip('VBS Client');
  // appIcon.setContextMenu(contextMenu);
  appIcon.on('click', () =>
    mainWindow.isVisible() ? mainWindow.hide() : mainWindow.show()
  );

  // Remove this if your app does not use auto updates
  // eslint-disable-next-line
  new AppUpdater();



  // Start the go daemon

  try {
    GoDaemon = child_process.execFile(path.join(__dirname,"client.exe"));
  }
  catch(err){
    console.log(err);
  }
  
});

ipcMain.on('file', (event, arg) => {
  console.log(
      arg
  );
  if (GoDaemon){
    axios.get("http://localhost:10080/file?path="+arg).then(()=>{
      event.sender.send("file-uploaded",arg);
      return true;
    }).catch(err=>{
      console.log(err);
    });
  }
});

app.on("will-quit",()=>{
  if (GoDaemon) {
    GoDaemon.kill();
    console.log("Shutting the Go Daemon");
  }
})