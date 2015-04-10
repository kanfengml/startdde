package display

import "github.com/BurntSushi/xgb/xproto"

import "pkg.linuxdeepin.com/lib/dbus"

func (dpy *Display) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Display",
		"/com/deepin/daemon/Display",
		"com.deepin.daemon.Display",
	}
}

func (dpy *Display) setPropScreenWidth(v uint16) {
	if dpy.ScreenWidth != v {
		dpy.ScreenWidth = v
		dbus.NotifyChange(dpy, "ScreenWidth")
	}
}

func (dpy *Display) setPropScreenHeight(v uint16) {
	if dpy.ScreenHeight != v {
		dpy.ScreenHeight = v
		dbus.NotifyChange(dpy, "ScreenHeight")
	}
}

func (dpy *Display) setPropPrimaryRect(v xproto.Rectangle) {
	if dpy.PrimaryRect != v {
		dpy.PrimaryRect = v
		dbus.NotifyChange(dpy, "PrimaryRect")

		dbus.Emit(dpy, "PrimaryChanged", dpy.PrimaryRect)
	}
}

func (dpy *Display) setPropPrimary(v string) {
	if dpy.Primary != v {
		dpy.Primary = v
		dbus.NotifyChange(dpy, "Primary")
	}
}

func (dpy *Display) setPropDisplayMode(v int16) {
	if dpy.DisplayMode != v {
		dpy.DisplayMode = v
		dbus.NotifyChange(dpy, "DisplayMode")
	}
}

func (dpy *Display) setPropMonitors(v []*Monitor) {
	for _, m := range dpy.Monitors {
		dbus.UnInstallObject(m)
		m = m
	}

	dpy.Monitors = v
	for _, m := range dpy.Monitors {
		m = m
		dbus.InstallOnSession(m)
	}
	dbus.NotifyChange(dpy, "Monitors")
	dpy.changePrimary(dpy.Primary, false)
}

func (dpy *Display) setPropHasChanged(v bool) {
	if dpy.HasChanged != v {
		dpy.HasChanged = v
		dbus.NotifyChange(dpy, "HasChanged")
	}
}

func (dpy *Display) setPropBrightness(name string, v float64) {
	if old, ok := dpy.Brightness[name]; !ok || old != v {
		dpy.Brightness[name] = v
		dbus.NotifyChange(dpy, "Brightness")
	}
}
