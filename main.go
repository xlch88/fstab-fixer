package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var colors = map[string]string{
	"1": "[0m[38;5;21m", //深蓝
	"2": "[0m[38;5;2m", //深绿
	"3": "[0m[38;5;6m", //湖蓝
	"4": "[0m[38;5;124m", //深红
	"5": "[0m[38;5;13m", //紫色
	"6": "[0m[38;5;208m", //橘色
	"7": "[0m[38;5;245m", //灰色
	"8": "[0m[38;5;239m", //深灰
	"9": "[0m[38;5;12m", //蓝
	"0": "[0m[38;5;0m", //纯黑

	"a": "[0m[38;5;46m", //亮绿
	"b": "[0m[38;5;51m", //艳青
	"c": "[0m[38;5;196m", //亮红
	"d": "[0m[38;5;200m", //品红
	"e": "[0m[38;5;226m", //亮黄
	"f": "[0m[38;5;255m", //纯白
	"p": "[0m[38;5;176m", //粉色!

	"n": "[4m", //下划线
	"l": "[1m", //粗体
	"r": "[0m", //重置

	//以下在windows不支持，linux良好
	"m": "[9m", //删除线
	"o": "[3m", //斜体

	"cc": "\u001B[0m",
}

var isDebug = false

func main() {
	fmt.Println(" " + colors["b"] + "/======================================\\")
	fmt.Println(" " + colors["b"] + "||                                    ||")
	fmt.Println(" " + colors["b"] + "||      " + colors["a"] + "fstab fixer " + colors["f"] + "| " + colors["d"] + "Version 1.0     " + colors["b"] + "||")
	fmt.Println(" " + colors["b"] + "||                                    ||")
	fmt.Println(" " + colors["b"] + "|| " + colors["p"] + "            By.Dark495 " + colors["b"] + "            ||")
	fmt.Println(" " + colors["b"] + "||                                    ||")
	fmt.Println(" " + colors["b"] + "|| " + colors["c"] + "     https://github.com/xlch88 " + colors["b"] + "    ||")
	fmt.Println(" " + colors["b"] + "||                                    ||")
	fmt.Println(" " + colors["b"] + "\\======================================/" + colors["cc"])
	fmt.Println(" " + "")

	timeoutDevices := getTimeoutDevices()

	if len(timeoutDevices) <= 0 {
		fmt.Println(colors["a"] + "All disk are fine. Everything is ok :)")
		fmt.Println(colors["a"] + "Wash and sleep, it's no fun.")
		fmt.Println(colors["2"] + "Or follow my twitter? @YueDongQwQ")
		return
	}

	fmt.Println(colors["e"] + "Ohhhhhhhhhhhhhhhhhhh !!!")
	fmt.Println(colors["e"] + "Your some disks looks like " + colors["c"] + "BOOM" + colors["e"] + " !!!")
	fmt.Println("")
	fmt.Println(colors["6"] + "List of damaged disks:")

	uuidMap := getUUIDMap()
	newFstab := string(getFstab())
	oldFstab := getFstab()

	for index, uuid := range timeoutDevices {
		fmt.Println(colors["f"]+" #", index, "- "+colors["c"]+"MountPoint "+colors["d"]+"= "+colors["f"]+"/"+uuidMap[uuid]+colors["2"]+", "+colors["b"]+"UUID "+colors["d"]+"= "+colors["f"]+uuid+colors["cc"])
		newFstab = strings.Replace(newFstab, "UUID="+uuid, "#UUID="+uuid, -1)
	}

	fmt.Println(colors["cc"])
	fmt.Println(colors["a"]+"Unmounted", len(timeoutDevices), "disks from /etc/fstab")
	fmt.Println(colors["a"] + "Old file backup to /etc/fstab.backup")

	if !isDebug {
		os.WriteFile("/etc/fstab", []byte(newFstab), 0777)
		os.WriteFile("/etc/fstab.backup", oldFstab, 0777)
	}

	fmt.Println(colors["9"] + "Please enter \"reboot\" to reboot.")
	fmt.Println(colors["cc"])
}

func getSystemdLog() string {
	var cmd *exec.Cmd
	var result []byte
	var err error

	if isDebug {
		result, _ = os.ReadFile("./test/log.txt")
	} else {
		cmd = exec.Command("/bin/sh", "-c", `journalctl -xb | grep "Timed out waiting for device"`)
		if result, err = cmd.Output(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	return string(result)
}

func getTimeoutDevices() []string {
	var uuids []string

	systemdLog := getSystemdLog()
	result := regexp.MustCompile(`Timed out waiting for device dev-disk-by\\x2duuid-(.*?)\.device.`).FindAllStringSubmatch(systemdLog, -1)
	for _, value := range result {
		uuids = append(uuids, strings.Replace(value[1], "\\x2d", "-", -1))
	}

	return uuids
}

func getFstab() []byte {
	var fstab []byte
	if isDebug {
		fstab, _ = os.ReadFile("./test/fstab.txt")
	} else {
		fstab, _ = os.ReadFile("/etc/fstab")
	}
	return fstab
}

func getUUIDMap() map[string]string {
	fstab := getFstab()
	fstabMap := make(map[string]string)

	result := regexp.MustCompile(`UUID=(.*?)\s+/(.*?)\s+(.*)`).FindAllStringSubmatch(string(fstab), -1)
	for _, value := range result {
		if len(value) < 3 {
			continue
		}
		fstabMap[value[1]] = value[2]
	}

	return fstabMap
}
