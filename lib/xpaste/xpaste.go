package xpaste

import "fmt"
import "runtime"
import "os/exec"


/*
This function is used to paste text on Linux using xdotool or wtype. It checks if either of these commands is available and then uses them to simulate keystrokes for the provided text.
*/
func PasteLinux(text string) error {
	if _, err := exec.LookPath("xdotool"); err == nil {
		exec.Command("xdotool", "type", "--clearmodifiers", "--delay", "0", text).Run()
	}

	if _, err := exec.LookPath("wtype"); err == nil {
		exec.Command("wtype", text).Run()
	}

	return nil
}

/*
This function is used to paste text on Windows using PowerShell. It checks if the "powershell" command is available and then uses it to simulate keystrokes for the provided text.
*/
func PasteWindows(text string) error {
	if _, err := exec.LookPath("powershell"); err == nil {
		exec.Command("powershell", "-Command", "Add-Type -AssemblyName System.Windows.Forms; [System.Windows.Forms.SendKeys]::SendWait('"+text+"')").Run()
	}

	return nil
}


/*
This function is used to paste text on macOS using AppleScript. It checks if the "osascript" command is available and then uses it to simulate keystrokes for the provided text.
*/
func PasteMacOS(text string) error {
	if _, err := exec.LookPath("osascript"); err == nil {
		exec.Command("osascript", "-e", "tell application \"System Events\" to keystroke \""+text+"\"").Run()
	}

	return nil
}


/*
This function is a wrapper that determines the operating system and calls the appropriate paste function for Linux, Windows, or macOS. It returns an error if the operating system is not supported.
*/
func PasteText(text string) error {
	switch runtime.GOOS {
	case "linux":
		return PasteLinux(text)
	case "windows":
		return PasteWindows(text)
	case "darwin":
		return PasteMacOS(text)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}
