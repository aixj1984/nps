package main

import (
	"context"
	"os"
	"sync"
	"time"

	"npc-deskreen/internal/npc"
	"npc-deskreen/internal/npc/configdoc"

	"github.com/djylb/nps/lib/logs"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// NpcDeskApp is the Wails-bound NPC desktop shell.
type NpcDeskApp struct {
	ctx context.Context
	mu  sync.Mutex
	mgr *npc.Manager
}

func NewNpcDeskApp() *NpcDeskApp {
	return &NpcDeskApp{mgr: npc.NewManager("")}
}

func (a *NpcDeskApp) startup(ctx context.Context) {
	a.ctx = ctx
	logs.EnableInMemoryBuffer(0)
	logs.Init("off", "info", "", 0, 0, 0, false, false)
	logs.Info("NPC Desk 已启动")
	_ = a.mgr.LoadDefault()
	a.emitStatus()
	go a.pollStatus()
}

func (a *NpcDeskApp) shutdown(context.Context) {
	a.mgr.Stop()
}

func (a *NpcDeskApp) pollStatus() {
	ticker := time.NewTicker(800 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-a.ctx.Done():
			return
		case <-ticker.C:
			a.emitStatus()
		}
	}
}

func (a *NpcDeskApp) emitStatus() {
	if a.ctx == nil {
		return
	}
	runtime.EventsEmit(a.ctx, "npc:update", a.statusMap())
}

func (a *NpcDeskApp) statusMap() map[string]any {
	st := a.mgr.Status()
	services := make([]map[string]any, len(st.Services))
	for i, svc := range st.Services {
		services[i] = map[string]any{
			"name":   svc.Name,
			"kind":   svc.Kind,
			"mode":   svc.Mode,
			"port":   svc.Port,
			"target": svc.Target,
			"active": svc.Active,
		}
	}
	return map[string]any{
		"running":          st.Running,
		"connected":        st.Connected,
		"reconnecting":     st.Reconnecting,
		"hasFailed":        st.HasFailed,
		"lastError":        st.LastError,
		"configPath":       st.ConfigPath,
		"serverAddr":       st.ServerAddr,
		"connType":         st.ConnType,
		"vkey":             st.VKey,
		"autoReconnection":   st.AutoReconnection,
		"version":            st.Version,
		"activeForwardCount": int(st.ActiveForwardCount),
		"totalForwardCount":  int(st.TotalForwardCount),
		"recentlyForwarded":  st.RecentlyForwarded,
		"services":           services,
	}
}

func (a *NpcDeskApp) GetStatus() map[string]any {
	return a.statusMap()
}

func (a *NpcDeskApp) GetConfigContent() string {
	return a.mgr.GetRawContent()
}

func (a *NpcDeskApp) GetConfigDocument() *configdoc.Document {
	doc, err := a.mgr.GetConfigDocument()
	if err != nil {
		return &configdoc.Document{Errors: []string{err.Error()}}
	}
	return doc
}

func (a *NpcDeskApp) SaveConfigDocument(doc *configdoc.Document) error {
	if err := a.mgr.SaveConfigDocument(doc); err != nil {
		return err
	}
	a.emitStatus()
	return nil
}

func (a *NpcDeskApp) ValidateConfigDocument(doc *configdoc.Document) *configdoc.Document {
	return a.mgr.ValidateConfigDocument(doc)
}

func (a *NpcDeskApp) GetConfigPath() string {
	return a.mgr.ConfigPath()
}

func (a *NpcDeskApp) SaveConfig(content string) error {
	if err := a.mgr.SaveContent(content); err != nil {
		return err
	}
	a.emitStatus()
	return nil
}

func (a *NpcDeskApp) PickConfigFile() (string, error) {
	if a.ctx == nil {
		return "", os.ErrInvalid
	}
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择 npc.conf",
		Filters: []runtime.FileFilter{
			{DisplayName: "NPC 配置 (*.conf)", Pattern: "*.conf"},
			{DisplayName: "所有文件", Pattern: "*.*"},
		},
	})
	if err != nil || path == "" {
		return "", err
	}
	if err := a.mgr.Load(path); err != nil {
		return "", err
	}
	a.emitStatus()
	return path, nil
}

func (a *NpcDeskApp) StartClient() error {
	if err := a.mgr.Start(); err != nil {
		return err
	}
	a.emitStatus()
	return nil
}

func (a *NpcDeskApp) StopClient() {
	a.mgr.Stop()
	a.emitStatus()
}

func (a *NpcDeskApp) GetLogs() string {
	return a.mgr.Logs()
}

func (a *NpcDeskApp) ClearLogs() {
	a.mgr.ClearLogs()
}

func (a *NpcDeskApp) GetCommonSettings() npc.CommonSettings {
	return a.mgr.CommonSettings()
}

func (a *NpcDeskApp) OpenConfigFolder() error {
	return a.mgr.OpenConfigDirectory()
}
