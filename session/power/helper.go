// SPDX-FileCopyrightText: 2018 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

package power

import (
	"github.com/godbus/dbus"
	notifications "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.notifications"
	"github.com/linuxdeepin/go-lib/dbusutil/proxy"

	// system bus
	daemon "github.com/linuxdeepin/go-dbus-factory/com.deepin.daemon.daemon"
	shutdownfront "github.com/linuxdeepin/go-dbus-factory/com.deepin.dde.shutdownfront"
	libpower "github.com/linuxdeepin/go-dbus-factory/com.deepin.system.power"
	sensorproxy "github.com/linuxdeepin/go-dbus-factory/net.hadess.sensorproxy"
	ofdbus "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.dbus"
	login1 "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.login1"

	// session bus
	display "github.com/linuxdeepin/go-dbus-factory/com.deepin.daemon.display"
	sessionwatcher "github.com/linuxdeepin/go-dbus-factory/com.deepin.daemon.sessionwatcher"
	sessionmanager "github.com/linuxdeepin/go-dbus-factory/com.deepin.sessionmanager"
	screensaver "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.screensaver"
	"github.com/linuxdeepin/go-lib/dbusutil"
	x "github.com/linuxdeepin/go-x11-client"
)

type Helper struct {
	Notifications notifications.Notifications

	Power         libpower.Power // sig
	LoginManager  login1.Manager // sig
	SensorProxy   sensorproxy.SensorProxy
	SysDBusDaemon ofdbus.DBus
	Daemon        daemon.Daemon

	SessionManager sessionmanager.SessionManager
	SessionWatcher sessionwatcher.SessionWatcher
	ShutdownFront  shutdownfront.ShutdownFront
	ScreenSaver    screensaver.ScreenSaver // sig
	Display        display.Display

	xConn *x.Conn
}

func newHelper(systemBus, sessionBus *dbus.Conn) (*Helper, error) {
	h := &Helper{}
	err := h.init(systemBus, sessionBus)
	if err != nil {
		return nil, err
	}
	return h, nil
}

func (h *Helper) init(sysBus, sessionBus *dbus.Conn) error {
	var err error

	h.Notifications = notifications.NewNotifications(sessionBus)

	h.Power = libpower.NewPower(sysBus)
	h.LoginManager = login1.NewManager(sysBus)
	h.SensorProxy = sensorproxy.NewSensorProxy(sysBus)
	h.SysDBusDaemon = ofdbus.NewDBus(sysBus)
	h.Daemon = daemon.NewDaemon(sysBus)
	h.SessionManager = sessionmanager.NewSessionManager(sessionBus)
	h.ScreenSaver = screensaver.NewScreenSaver(sessionBus)
	h.Display = display.NewDisplay(sessionBus)
	h.SessionWatcher = sessionwatcher.NewSessionWatcher(sessionBus)
	h.ShutdownFront = shutdownfront.NewShutdownFront(sessionBus)

	// init X conn
	h.xConn, err = x.NewConn()
	if err != nil {
		return err
	}
	return nil
}

func (h *Helper) initSignalExt(systemSigLoop, sessionSigLoop *dbusutil.SignalLoop) {
	// sys
	h.SysDBusDaemon.InitSignalExt(systemSigLoop, true)
	h.LoginManager.InitSignalExt(systemSigLoop, true)
	h.Power.InitSignalExt(systemSigLoop, true)
	h.SensorProxy.InitSignalExt(systemSigLoop, true)
	h.Daemon.InitSignalExt(systemSigLoop, true)
	// session
	h.ScreenSaver.InitSignalExt(sessionSigLoop, true)
	h.SessionWatcher.InitSignalExt(sessionSigLoop, true)
	h.Display.InitSignalExt(sessionSigLoop, true)
}

func (h *Helper) Destroy() {
	h.SysDBusDaemon.RemoveHandler(proxy.RemoveAllHandlers)
	h.LoginManager.RemoveHandler(proxy.RemoveAllHandlers)
	h.Power.RemoveHandler(proxy.RemoveAllHandlers)
	h.SensorProxy.RemoveHandler(proxy.RemoveAllHandlers)
	h.Daemon.RemoveHandler(proxy.RemoveAllHandlers)

	h.ScreenSaver.RemoveHandler(proxy.RemoveAllHandlers)
	h.SessionWatcher.RemoveHandler(proxy.RemoveAllHandlers)

	if h.xConn != nil {
		h.xConn.Close()
		h.xConn = nil
	}
}
