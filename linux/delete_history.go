// Install: go get golang.org/x/sys/windows/registry

package main

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"os"
	"syscall"
)

const (
	HKEY_CURRENT_USER = registry.CURRENT_USER
	KEY_ALL_ACCESS   = registry.ALL_ACCESS
)

func main() {
	if isElevated() {
		key, err := registry.OpenKey(registry.CURRENT_USER, "Software\\Microsoft\\Windows\\CurrentVersion\\Explorer\\RunMRU", registry.ALL_ACCESS)
		if err != nil {
			fmt.Printf("Error opening registry key: %v\n", err)
			os.Exit(1)
		}
		defer key.Close()

		subkeys, err := key.ReadValueNames(-1)
		if err != nil {
			fmt.Printf("Error reading registry values: %v\n", err)
			os.Exit(1)
		}

		for _, name := range subkeys {
			if err := key.DeleteValue(name); err != nil {
				fmt.Printf("Error deleting registry value %s: %v\n", name, err)
			} else {
				fmt.Printf("%s deleted\n", name)
			}
		}
	} else {
		runAsAdmin()
	}
}

func isElevated() bool {
	var adminSID *syscall.SID
	syscall.AllocateAndInitializeSid(&syscall.SECURITY_NT_AUTHORITY, 2, syscall.SECURITY_BUILTIN_DOMAIN_RID, syscall.DOMAIN_ALIAS_RID_ADMINS, 0, 0, 0, 0, 0, 0, &adminSID)
	defer syscall.FreeSid(adminSID)

	isElevated := false
	syscall.CheckTokenMembership(0, adminSID, &isElevated)
	return isElevated
}

func runAsAdmin() {
	// The code to restart the process with admin privileges goes here.
	fmt.Println("Please run this program with administrative privileges.")
}
