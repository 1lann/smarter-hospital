// Code generated by go-bindata.
// sources:
// navbar/navbar.tmpl
// notify/notify.tmpl
// patient-navbar/patient_navbar.tmpl
// patient-room/patient_room.tmpl
// room/room.tmpl
// DO NOT EDIT!

package views

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)
type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _navbarNavbarTmpl = []byte(`<div>

	<div class="ui top fixed menu">
		<a @click="ToggleMenu" class="on-mobile item">
			<i class="sidebar icon"></i>
		</a>

		<div class="ui container on-desktop">
			<div class="left menu">
				<span class="ui header item">
					Hospital+
				</span>
				<a v-for="item in left" :href="item.path" class="item" :class="{active: item.active}">
					{{item.name}}
				</a>
			</div>

			<div class="right menu">
				<a v-for="item in right" :href="item.path" class="item" :class="{active: item.active}">
					{{item.name}}
				</a>
			</div>
		</div>
	</div>

	<div class="ui vertical sidebar menu mobile-navbar">
		<span class="ui header item">
			Hospital+
		</span>
		<a v-for="item in left" :href="item.path" class="item" :class="{active: item.active}">
			{{item.name}}
		</a>
	</div>

</div>
`)

func navbarNavbarTmplBytes() ([]byte, error) {
	return _navbarNavbarTmpl, nil
}

func navbarNavbarTmpl() (*asset, error) {
	bytes, err := navbarNavbarTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "navbar/navbar.tmpl", size: 801, mode: os.FileMode(420), modTime: time.Unix(1493958209, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _notifyNotifyTmpl = []byte(`<div>
	<h2 class="header">Horizontal Card</h2>
	<div class="card horizontal">
		<div class="card-image">
			<img src="http://lorempixel.com/100/190/nature/6">
		</div>
		<div class="card-stacked">
			<div class="card-content">
				<p>I am a very simple card. I am good at containing small bits of information.</p>
			</div>
			<div class="card-action">
				<a href="#" class="waves-effect">This is a link</a>
			</div>
		</div>
	</div>
</div>
`)

func notifyNotifyTmplBytes() ([]byte, error) {
	return _notifyNotifyTmpl, nil
}

func notifyNotifyTmpl() (*asset, error) {
	bytes, err := notifyNotifyTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "notify/notify.tmpl", size: 443, mode: os.FileMode(420), modTime: time.Unix(1493904086, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _patientNavbarPatient_navbarTmpl = []byte(`<div>
	<div class="ui top fixed menu patient-navbar">
		<div class="ui header item on-mobile">
			{{name}}
		</div>

		<div class="right item on-mobile">
			<div class="ui red labeled icon button" @click="CallNurse">
				<i class="doctor icon"></i>
				Call Nurse
			</div>
		</div>

		<div class="ui container on-desktop" style="height:70px;">
			<div class="left menu">
				<div class="ui header item">
					<div class="content">
						{{time}}
						<div class="sub header">{{date}}</div>
					</div>
				</div>
				<div class="ui header item">
					<div class="content">
						{{name}}
						<div class="sub header">Room {{roomNumber}}</div>
					</div>
				</div>
			</div>

			<div class="right menu">
				<div class="item">
					<div class="ui red labeled icon button" @click="CallNurse" style="font-size:18px;">
						<i class="doctor icon"></i>
						Call Nurse
					</div>
				</div>
			</div>
		</div>
	</div>
</div>
`)

func patientNavbarPatient_navbarTmplBytes() ([]byte, error) {
	return _patientNavbarPatient_navbarTmpl, nil
}

func patientNavbarPatient_navbarTmpl() (*asset, error) {
	bytes, err := patientNavbarPatient_navbarTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "patient-navbar/patient_navbar.tmpl", size: 928, mode: os.FileMode(420), modTime: time.Unix(1494228505, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _patientRoomPatient_roomTmpl = []byte(`<div>
<patient-navbar></patient-navbar>
<div class="ui container with-patient-navbar">
	<h1 class="ui header">
		Good {{greeting}}, {{name}}.
	</h1>

	<div class="ui grid">
		<div class="eight wide computer eight wide tablet sixteen wide mobile column">
			<h2 class="ui header">
				<i class="settings icon"></i>
				<div class="content">
					Room controls
				</div>
			</h2>

			<a class="ui fluid link grey card">
				<div class="ui padded grid">
					<div class="sixteen wide column">
						<h3 class="ui header">
							<i class="grey idea icon"></i>
							<div class="content">
								Lights off
								<div class="sub header">Press to turn on lights</div>
							</div>
						</h3>
					</div>
				</div>
			</a>

			<a class="ui fluid link yellow card">
				<div class="ui padded grid">
					<div class="sixteen wide column">
						<h3 class="ui header">
							<i class="yellow idea icon"></i>
							<div class="content">
								Lights on
								<div class="sub header">Press to turn off lights</div>
							</div>
						</h3>
					</div>
				</div>
			</a>

			<a class="ui fluid link orange card">
				<div class="ui padded grid">
					<div class="sixteen wide column">
						<h3 class="ui header">
							<i class="orange sun icon"></i>
							<div class="content">
								Heating from 24&deg;C to 26&deg;C
								<div class="sub header">Change temperature and view history</div>
							</div>
						</h3>
					</div>
				</div>
			</a>

			<a class="ui fluid link blue card">
				<div class="ui padded grid">
					<div class="sixteen wide column">
						<h3 class="ui header">
							<i class="blue fa fa-snowflake-o icon"></i>
							<div class="content">
								Cooling from 26&deg;C to 24&deg;C
								<div class="sub header">Change temperature and view history</div>
							</div>
						</h3>
					</div>
				</div>
			</a>

			<h2 class="ui header">
				<i class="green plus icon"></i>
				<div class="content">
					Your health is OK
				</div>
			</h2>

			<a class="ui fluid link green card">
				<div class="ui padded grid">
					<div class="sixteen wide column">
						<h3 class="ui header">
							<i class="green heartbeat icon"></i>
							<div class="content">
								100 BPM
								<div class="sub header">View history</div>
							</div>
						</h3>
					</div>
				</div>
			</a>

			<a class="ui fluid link green card">
				<div class="ui padded grid">
					<div class="sixteen wide column">
						<h3 class="ui header">
							<i class="green hotel icon"></i>
							<div class="content">
								In bed
								<div class="sub header">View history</div>
							</div>
						</h3>
					</div>
				</div>
			</a>


			<a class="ui fluid link green card">
				<div class="ui padded grid">
					<div class="sixteen wide column">
						<h3 class="ui header">
							<i class="green theme icon"></i>
							<div class="content">
								SpO<sub>2</sub>: 90%
								<div class="sub header">View history</div>
							</div>
						</h3>
					</div>
				</div>
			</a>
		</div>
		<div class="eight wide computer eight wide tablet sixteen wide mobile column">
			<h2 class="ui header">
				<i class="calendar icon"></i>
				<div class="content">
					Your agenda
				</div>
			</h2>

			<form class="ui form" @submit.prevent="Ping">
				<div class="field">
					<label>Ping this</label>
					<input v-model="pingText" type="text" name="ping" placeholder="Ping text">
				</div>
				<button class="ui button" type="submit">Submit</button>
			</form>

		</div>
	</div>

</div>

</div>
`)

func patientRoomPatient_roomTmplBytes() ([]byte, error) {
	return _patientRoomPatient_roomTmpl, nil
}

func patientRoomPatient_roomTmpl() (*asset, error) {
	bytes, err := patientRoomPatient_roomTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "patient-room/patient_room.tmpl", size: 3520, mode: os.FileMode(420), modTime: time.Unix(1494437326, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _roomRoomTmpl = []byte(`<div>
<patient-navbar></patient-navbar>


<div class="ui container">
	<h1>Memes</h1>
</div>

</div>
`)

func roomRoomTmplBytes() ([]byte, error) {
	return _roomRoomTmpl, nil
}

func roomRoomTmpl() (*asset, error) {
	bytes, err := roomRoomTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "room/room.tmpl", size: 100, mode: os.FileMode(420), modTime: time.Unix(1493958402, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"navbar/navbar.tmpl": navbarNavbarTmpl,
	"notify/notify.tmpl": notifyNotifyTmpl,
	"patient-navbar/patient_navbar.tmpl": patientNavbarPatient_navbarTmpl,
	"patient-room/patient_room.tmpl": patientRoomPatient_roomTmpl,
	"room/room.tmpl": roomRoomTmpl,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"navbar": &bintree{nil, map[string]*bintree{
		"navbar.tmpl": &bintree{navbarNavbarTmpl, map[string]*bintree{}},
	}},
	"notify": &bintree{nil, map[string]*bintree{
		"notify.tmpl": &bintree{notifyNotifyTmpl, map[string]*bintree{}},
	}},
	"patient-navbar": &bintree{nil, map[string]*bintree{
		"patient_navbar.tmpl": &bintree{patientNavbarPatient_navbarTmpl, map[string]*bintree{}},
	}},
	"patient-room": &bintree{nil, map[string]*bintree{
		"patient_room.tmpl": &bintree{patientRoomPatient_roomTmpl, map[string]*bintree{}},
	}},
	"room": &bintree{nil, map[string]*bintree{
		"room.tmpl": &bintree{roomRoomTmpl, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

