//go:build windows

package command

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// https://github.com/golang/go/issues/17608

type job struct {
	hJob windows.Handle
}

func (j *job) Close() error {
	if j != nil && j.hJob != 0 {
		return windows.CloseHandle(j.hJob)
	}
	return nil
}

func (c *Command) beforeInitialize() error {
	var createFlags uint32 = windows.CREATE_UNICODE_ENVIRONMENT
	if c.ProcessGroup {
		createFlags |= windows.CREATE_NEW_PROCESS_GROUP
	}
	c.rawCmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: createFlags,
	}
	return nil
}

func (c *Command) delayInitialize() error {
	if !c.ProcessGroup {
		return nil
	}
	hProcess, err := windows.OpenProcess(windows.MAXIMUM_ALLOWED, false, uint32(c.rawCmd.Process.Pid))
	if err != nil {
		return err
	}
	defer windows.CloseHandle(hProcess)
	hJob, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		return err
	}
	ji := windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION{
		BasicLimitInformation: windows.JOBOBJECT_BASIC_LIMIT_INFORMATION{
			LimitFlags: windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE,
		},
	}
	if _, err = windows.SetInformationJobObject(hJob, windows.JobObjectExtendedLimitInformation, uintptr(unsafe.Pointer(&ji)), uint32(unsafe.Sizeof(ji))); err != nil {
		_ = windows.CloseHandle(hJob)
		return err
	}

	if err := windows.AssignProcessToJobObject(hJob, hProcess); err != nil {
		_ = windows.CloseHandle(hJob)
		return err
	}
	c.closeAfterWait = append(c.closeAfterWait, &job{hJob: hJob})
	return nil
}

func (c *Command) finalize() error {
	var err error
	for _, closer := range c.closeAfterWait {
		if closeErr := closer.Close(); closeErr != nil {
			err = closeErr
		}
	}
	return err
}
